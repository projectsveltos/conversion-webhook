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
	utilsv1alpha1 "github.com/projectsveltos/sveltosctl/api/v1alpha1"
	utilsv1beta1 "github.com/projectsveltos/sveltosctl/api/v1beta1"
)

var _ = Describe("Techsupport", func() {
	const (
		namePrefix = "conversion-"
	)

	It("Deploy Techsupport v1alpha1 and query for v1beta1", Label("FV", "EXTENDED"), func() {
		key := randomString()
		value := randomString()

		techsupport := &utilsv1alpha1.Techsupport{
			ObjectMeta: metav1.ObjectMeta{
				Name: namePrefix + randomString(),
			},
			Spec: utilsv1alpha1.TechsupportSpec{
				ClusterSelector: libsveltosv1alpha1.Selector(fmt.Sprintf("%s=%s", key, value)),
				Schedule:        "0 * * * *",
				Resources: []utilsv1alpha1.Resource{
					{
						Namespace: randomString(),
						Group:     randomString(),
						Version:   randomString(),
						Kind:      randomString(),
					},
				},
				Logs: []utilsv1alpha1.Log{
					{
						Namespace: randomString(),
					},
				},
			},
		}

		Byf("Create Techsupport.v1alpha1 %s", techsupport.Name)
		Expect(k8sClient.Create(context.TODO(), techsupport)).To(Succeed())

		Byf("Get Techsupport.v1beta1 %s", techsupport.Name)
		dst := &utilsv1beta1.Techsupport{}
		Expect(k8sClient.Get(context.TODO(),
			types.NamespacedName{Name: techsupport.Name},
			dst)).To(Succeed())

		Byf("Verify Techsupport.v1beta1 selector %s", techsupport.Name)
		Expect(len(dst.Spec.ClusterSelector.LabelSelector.MatchLabels)).To(Equal(1))
		Expect(dst.Spec.ClusterSelector.LabelSelector.MatchLabels[key]).To(Equal(value))

		currentTechsupport := &utilsv1alpha1.Techsupport{}
		Expect(k8sClient.Get(context.TODO(),
			types.NamespacedName{Name: techsupport.Name},
			currentTechsupport)).To(Succeed())

		Byf("Delete Techsupport.v1alpha1 %s", currentTechsupport.Name)
		Expect(k8sClient.Delete(context.TODO(), currentTechsupport)).To(Succeed())
	})
})
