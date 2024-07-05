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

	libsveltosv1alpha1 "github.com/projectsveltos/libsveltos/api/v1alpha1"
	libsveltosv1beta1 "github.com/projectsveltos/libsveltos/api/v1beta1"
)

var _ = Describe("ClusterHealthCheck", func() {
	const (
		namePrefix = "conversion-"
	)

	It("Deploy ClusterHealthCheck v1alpha1 and query for v1beta1", Label("FV", "EXTENDED"), func() {
		key := randomString()
		value := randomString()

		clusterHealthCheck := &libsveltosv1alpha1.ClusterHealthCheck{
			ObjectMeta: metav1.ObjectMeta{
				Name: namePrefix + randomString(),
			},
			Spec: libsveltosv1alpha1.ClusterHealthCheckSpec{
				ClusterSelector: libsveltosv1alpha1.Selector(fmt.Sprintf("%s=%s", key, value)),
				LivenessChecks: []libsveltosv1alpha1.LivenessCheck{
					{
						Name: randomString(),
						Type: libsveltosv1alpha1.LivenessTypeHealthCheck,
						LivenessSourceRef: &corev1.ObjectReference{
							Kind:       libsveltosv1alpha1.HealthCheckKind,
							APIVersion: libsveltosv1alpha1.GroupVersion.String(),
							Name:       randomString(),
						},
					},
				},
				Notifications: []libsveltosv1alpha1.Notification{
					{
						Name: randomString(),
						Type: libsveltosv1alpha1.NotificationTypeSlack,
						NotificationRef: &corev1.ObjectReference{
							Kind:       string(libsveltosv1alpha1.ConfigMapReferencedResourceKind),
							APIVersion: "v1",
							Namespace:  randomString(),
							Name:       randomString(),
						},
					},
				},
			},
		}

		Byf("Create ClusterHealthCheck.v1alpha1 %s", clusterHealthCheck.Name)
		Expect(k8sClient.Create(context.TODO(), clusterHealthCheck)).To(Succeed())

		Byf("Get ClusterHealthCheck.v1beta1 %s", clusterHealthCheck.Name)
		dst := &libsveltosv1beta1.ClusterHealthCheck{}
		Expect(k8sClient.Get(context.TODO(),
			types.NamespacedName{Name: clusterHealthCheck.Name},
			dst)).To(Succeed())

		Byf("Verify ClusterHealthCheck.v1beta1 selector %s", clusterHealthCheck.Name)
		Expect(len(dst.Spec.ClusterSelector.LabelSelector.MatchLabels)).To(Equal(1))
		Expect(dst.Spec.ClusterSelector.LabelSelector.MatchLabels[key]).To(Equal(value))

		currentClusterHealthCheck := &libsveltosv1alpha1.ClusterHealthCheck{}
		Expect(k8sClient.Get(context.TODO(),
			types.NamespacedName{Name: clusterHealthCheck.Name},
			currentClusterHealthCheck)).To(Succeed())

		Byf("Delete ClusterHealthCheck.v1alpha1 %s", currentClusterHealthCheck.Name)
		Expect(k8sClient.Delete(context.TODO(), currentClusterHealthCheck)).To(Succeed())
	})
})
