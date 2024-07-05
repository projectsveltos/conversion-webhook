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

var _ = Describe("ClusterSet", func() {
	const (
		namePrefix = "conversion-"
	)

	It("Deploy ClusterSet v1alpha1 and query for v1beta1", Label("FV", "EXTENDED"), func() {
		key := randomString()
		value := randomString()

		clusterSet := &libsveltosv1alpha1.ClusterSet{
			ObjectMeta: metav1.ObjectMeta{
				Name: namePrefix + randomString(),
			},
			Spec: libsveltosv1alpha1.Spec{
				ClusterSelector: libsveltosv1alpha1.Selector(fmt.Sprintf("%s=%s", key, value)),
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

		Byf("Create ClusterSet.v1alpha1 %s", clusterSet.Name)
		Expect(k8sClient.Create(context.TODO(), clusterSet)).To(Succeed())

		Byf("Get ClusterSet.v1beta1 %s", clusterSet.Name)
		dst := &libsveltosv1beta1.ClusterSet{}
		Expect(k8sClient.Get(context.TODO(),
			types.NamespacedName{Name: clusterSet.Name},
			dst)).To(Succeed())

		Byf("Verify ClusterSet.v1beta1 selector %s", clusterSet.Name)
		Expect(len(dst.Spec.ClusterSelector.LabelSelector.MatchLabels)).To(Equal(1))
		Expect(dst.Spec.ClusterSelector.LabelSelector.MatchLabels[key]).To(Equal(value))

		currentClusterSet := &libsveltosv1alpha1.ClusterSet{}
		Expect(k8sClient.Get(context.TODO(),
			types.NamespacedName{Name: clusterSet.Name},
			currentClusterSet)).To(Succeed())

		Byf("Delete ClusterSet.v1alpha1 %s", currentClusterSet.Name)
		Expect(k8sClient.Delete(context.TODO(), currentClusterSet)).To(Succeed())
	})
})
