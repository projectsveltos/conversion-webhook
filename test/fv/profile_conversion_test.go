/*
Copyright 2024. projectsveltos.io. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fv_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	configv1alpha1 "github.com/projectsveltos/addon-controller/api/v1alpha1"
	configv1beta1 "github.com/projectsveltos/addon-controller/api/v1beta1"
	libsveltosv1alpha1 "github.com/projectsveltos/libsveltos/api/v1alpha1"
	libsveltosv1beta1 "github.com/projectsveltos/libsveltos/api/v1beta1"
)

var _ = Describe("Profile", func() {
	const (
		namePrefix = "conversion-"
	)

	It("Deploy Profile v1alpha1 and query for v1beta1", Label("FV", "EXTENDED"), func() {
		key1 := randomString()
		value1 := randomString()
		key2 := randomString()
		value2 := randomString()

		profile := configv1alpha1.Profile{
			ObjectMeta: metav1.ObjectMeta{
				Name:      randomString(),
				Namespace: namePrefix + randomString(),
			},
			Spec: configv1alpha1.Spec{
				ClusterSelector: libsveltosv1alpha1.Selector(fmt.Sprintf("%s=%s, %s=%s", key1, value1, key2, value2)),
				HelmCharts: []configv1alpha1.HelmChart{
					{
						RepositoryURL:    randomString(),
						RepositoryName:   randomString(),
						ChartName:        randomString(),
						ChartVersion:     randomString(),
						ReleaseName:      randomString(),
						ReleaseNamespace: randomString(),
						Values:           randomString(),
						Options: &configv1alpha1.HelmOptions{
							Labels: map[string]string{
								randomString(): randomString(),
							},
						},
					},
					{
						RepositoryURL:    randomString(),
						RepositoryName:   randomString(),
						ChartName:        randomString(),
						ChartVersion:     randomString(),
						ReleaseName:      randomString(),
						ReleaseNamespace: randomString(),
						Values:           randomString(),
						Options: &configv1alpha1.HelmOptions{
							Labels: map[string]string{
								randomString(): randomString(),
							},
						},
					},
				},
				PolicyRefs: []configv1alpha1.PolicyRef{
					{
						Kind:      string(libsveltosv1beta1.ConfigMapReferencedResourceKind),
						Namespace: randomString(),
						Name:      randomString(),
					},
					{
						Kind:      string(libsveltosv1beta1.ConfigMapReferencedResourceKind),
						Namespace: randomString(),
						Name:      randomString(),
					},
				},
				KustomizationRefs: []configv1alpha1.KustomizationRef{
					{
						Namespace: randomString(),
						Name:      randomString(),
						Kind:      string(libsveltosv1alpha1.SecretReferencedResourceKind),
					},
					{
						Namespace: randomString(),
						Name:      randomString(),
						Kind:      string(libsveltosv1alpha1.SecretReferencedResourceKind),
					},
				},
				ClusterRefs: []corev1.ObjectReference{
					{
						Kind:      libsveltosv1alpha1.SveltosClusterKind,
						Namespace: randomString(),
						Name:      randomString(),
					},
					{
						Kind:      libsveltosv1alpha1.SveltosClusterKind,
						Namespace: randomString(),
						Name:      randomString(),
					},
				},
			},
		}

		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: profile.Namespace,
			},
		}

		Byf("Create Namespace %s", ns.Name)
		Expect(k8sClient.Create(context.TODO(), ns)).To(Succeed())

		Byf("Create Profile.v1alpha1 %s", profile.Name)
		Expect(k8sClient.Create(context.TODO(), &profile)).To(Succeed())

		Byf("Get Profile.v1beta1 %s/%s", profile.Namespace, profile.Name)
		dst := &configv1beta1.Profile{}
		Expect(k8sClient.Get(context.TODO(),
			types.NamespacedName{Namespace: profile.Namespace, Name: profile.Name},
			dst)).To(Succeed())

		Byf("Verify Profile.v1beta1 selector %s", profile.Name)
		Expect(len(dst.Spec.ClusterSelector.LabelSelector.MatchLabels)).To(Equal(2))
		Expect(dst.Spec.ClusterSelector.LabelSelector.MatchLabels[key1]).To(Equal(value1))
		Expect(dst.Spec.ClusterSelector.LabelSelector.MatchLabels[key2]).To(Equal(value2))

		currentProfile := &configv1alpha1.Profile{}
		Expect(k8sClient.Get(context.TODO(),
			types.NamespacedName{Namespace: profile.Namespace, Name: profile.Name},
			currentProfile)).To(Succeed())

		Byf("Delete Profile.v1alpha1 %s/%s", currentProfile.Namespace, currentProfile.Name)
		Expect(k8sClient.Delete(context.TODO(), currentProfile)).To(Succeed())
	})
})
