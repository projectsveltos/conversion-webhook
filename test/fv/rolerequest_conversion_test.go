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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	libsveltosv1alpha1 "github.com/projectsveltos/libsveltos/api/v1alpha1"
	libsveltosv1beta1 "github.com/projectsveltos/libsveltos/api/v1beta1"
)

var _ = Describe("RoleRequest", func() {
	const (
		namePrefix = "conversion-"
	)

	It("Deploy RoleRequest v1alpha1 and query for v1beta1", Label("FV", "EXTENDED"), func() {
		key := randomString()
		value := randomString()

		expirationSeconds := int64(600)

		roleRequest := &libsveltosv1alpha1.RoleRequest{
			ObjectMeta: metav1.ObjectMeta{
				Name: namePrefix + randomString(),
			},
			Spec: libsveltosv1alpha1.RoleRequestSpec{
				ClusterSelector: libsveltosv1alpha1.Selector(fmt.Sprintf("%s=%s", key, value)),
				RoleRefs: []libsveltosv1alpha1.PolicyRef{
					{
						Kind:      string(libsveltosv1alpha1.ConfigMapReferencedResourceKind),
						Namespace: randomString(),
						Name:      randomString(),
					},
					{
						Kind:      string(libsveltosv1alpha1.SecretReferencedResourceKind),
						Namespace: randomString(),
						Name:      randomString(),
					},
				},
				ExpirationSeconds:       &expirationSeconds,
				ServiceAccountName:      randomString(),
				ServiceAccountNamespace: randomString(),
			},
		}

		Byf("Create RoleRequest.v1alpha1 %s", roleRequest.Name)
		Expect(k8sClient.Create(context.TODO(), roleRequest)).To(Succeed())

		Byf("Get RoleRequest.v1beta1 %s", roleRequest.Name)
		dst := &libsveltosv1beta1.RoleRequest{}
		Expect(k8sClient.Get(context.TODO(),
			types.NamespacedName{Name: roleRequest.Name},
			dst)).To(Succeed())

		Byf("Verify RoleRequest.v1beta1 selector %s", roleRequest.Name)
		Expect(len(dst.Spec.ClusterSelector.LabelSelector.MatchLabels)).To(Equal(1))
		Expect(dst.Spec.ClusterSelector.LabelSelector.MatchLabels[key]).To(Equal(value))

		currentRoleRequest := &libsveltosv1alpha1.RoleRequest{}
		Expect(k8sClient.Get(context.TODO(),
			types.NamespacedName{Name: roleRequest.Name},
			currentRoleRequest)).To(Succeed())

		Byf("Delete RoleRequest.v1alpha1 %s", currentRoleRequest.Name)
		Expect(k8sClient.Delete(context.TODO(), currentRoleRequest)).To(Succeed())
	})
})
