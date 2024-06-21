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

	"github.com/projectsveltos/event-manager/api/v1alpha1"
	"github.com/projectsveltos/event-manager/api/v1beta1"
	libsveltosv1alpha1 "github.com/projectsveltos/libsveltos/api/v1alpha1"
)

var _ = Describe("EventTrigger", func() {
	const (
		namePrefix = "conversion-"
	)

	It("Deploy EventTrigegr v1alpha1 and query for v1beta1", Label("FV", "EXTENDED"), func() {
		key := randomString()
		value := randomString()

		eventTrigger := &v1alpha1.EventTrigger{
			ObjectMeta: metav1.ObjectMeta{
				Name: namePrefix + randomString(),
			},
			Spec: v1alpha1.EventTriggerSpec{
				SourceClusterSelector: libsveltosv1alpha1.Selector(fmt.Sprintf("%s=%s", key, value)),
				EventSourceName:       randomString(),
				OneForEvent:           true,
			},
		}

		Byf("Create EventTrigger.v1alpha1 %s", eventTrigger.Name)
		Expect(k8sClient.Create(context.TODO(), eventTrigger)).To(Succeed())

		Byf("Get EventTrigger.v1beta1 %s", eventTrigger.Name)
		dst := &v1beta1.EventTrigger{}
		Expect(k8sClient.Get(context.TODO(),
			types.NamespacedName{Name: eventTrigger.Name},
			dst)).To(Succeed())

		Byf("Verify EventTrigger.v1beta1 selector %s", eventTrigger.Name)
		Expect(len(dst.Spec.SourceClusterSelector.LabelSelector.MatchLabels)).To(Equal(1))
		Expect(dst.Spec.SourceClusterSelector.LabelSelector.MatchLabels[key]).To(Equal(value))

		currentEventTrigger := &v1alpha1.EventTrigger{}
		Expect(k8sClient.Get(context.TODO(),
			types.NamespacedName{Name: eventTrigger.Name},
			currentEventTrigger)).To(Succeed())

		Byf("Delete EventTrigger.v1alpha1 %s", currentEventTrigger.Name)
		Expect(k8sClient.Delete(context.TODO(), currentEventTrigger)).To(Succeed())
	})
})
