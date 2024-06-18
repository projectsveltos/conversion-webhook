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

var _ = Describe("Set", func() {
	const (
		namePrefix = "conversion-"
	)

	It("Deploy Set v1alpha1 and query for v1beta1", Label("FV", "EXTENDED"), func() {
		key1 := randomString()
		value1 := randomString()
		key2 := randomString()
		value2 := randomString()

		set := libsveltosv1alpha1.Set{
			ObjectMeta: metav1.ObjectMeta{
				Name:      randomString(),
				Namespace: namePrefix + randomString(),
			},
			Spec: libsveltosv1alpha1.Spec{
				ClusterSelector: libsveltosv1alpha1.Selector(fmt.Sprintf("%s=%s,%s=%s", key1, value1, key2, value2)),
				ClusterRefs: []corev1.ObjectReference{
					{
						Kind:      libsveltosv1alpha1.SveltosClusterKind,
						Namespace: randomString(),
						Name:      randomString(),
					},
				},
				MaxReplicas: 1,
			},
		}

		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: set.Namespace,
			},
		}

		Byf("Create Namespace %s", ns.Name)
		Expect(k8sClient.Create(context.TODO(), ns)).To(Succeed())

		Byf("Create Set.v1alpha1 %s", set.Name)
		Expect(k8sClient.Create(context.TODO(), &set)).To(Succeed())

		Byf("Get Set.v1beta1 %s/%s", set.Namespace, set.Name)
		dst := &libsveltosv1beta1.Set{}
		Expect(k8sClient.Get(context.TODO(),
			types.NamespacedName{Namespace: set.Namespace, Name: set.Name},
			dst)).To(Succeed())

		Byf("Verify Set.v1beta1 selector %s", set.Name)
		Expect(len(dst.Spec.ClusterSelector.LabelSelector.MatchLabels)).To(Equal(2))
		Expect(dst.Spec.ClusterSelector.LabelSelector.MatchLabels[key1]).To(Equal(value1))
		Expect(dst.Spec.ClusterSelector.LabelSelector.MatchLabels[key2]).To(Equal(value2))

		currentSet := &libsveltosv1alpha1.Set{}
		Expect(k8sClient.Get(context.TODO(),
			types.NamespacedName{Namespace: set.Namespace, Name: set.Name},
			currentSet)).To(Succeed())

		Byf("Delete Set.v1alpha1 %s/%s", currentSet.Namespace, currentSet.Name)
		Expect(k8sClient.Delete(context.TODO(), currentSet)).To(Succeed())
	})
})
