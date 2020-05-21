/**
 * Copyright 2020 Appvia Ltd <info@appvia.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package security_test

import (
	"context"
	"encoding/json"
	"io/ioutil"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	securityv1 "github.com/appvia/kore/pkg/apis/security/v1"
	"github.com/appvia/kore/pkg/security"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

var _ = Describe("GKE Autoscaling Rule", func() {
	scanner := security.NewEmpty()

	BeforeEach(func() {
		scanner.RegisterRule(&security.GKEAutoscaling{})
	})

	renderPlan := func(values map[string]interface{}) *configv1.Plan {
		c, _ := json.Marshal(values)

		return &configv1.Plan{
			TypeMeta: metav1.TypeMeta{
				APIVersion: configv1.GroupVersion.String(),
				Kind:       "Plan",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "test",
				Name:      "test1",
			},
			Spec: configv1.PlanSpec{
				Configuration: apiextv1.JSON{Raw: c},
			},
		}
	}

	Describe("CheckPlan", func() {
		When("called with a plan which not related to gke", func() {
			It("should return an ignore", func() {
				p := renderPlan(map[string]interface{}{})
				p.Spec.Kind = "EKS"
				result := scanner.ScanPlan(context.Background(), fake.NewFakeClient(), p)

				Expect(result.Spec.OverallStatus).To(Equal(securityv1.Compliant))
				Expect(len(result.Spec.Results)).To(Equal(0))
			})
		})

		Context("GKE Plan", func() {
			plan := renderPlan(map[string]interface{}{})
			plan.Spec.Kind = "GKE"

			When("called with GKE kind but no enableAutoscaler value", func() {
				It("should not have a result", func() {
					result := scanner.ScanPlan(context.Background(), fake.NewFakeClient(), plan)
					Expect(result.Spec.OverallStatus).To(Equal(securityv1.Compliant))
					Expect(len(result.Spec.Results)).To(Equal(0))
				})
			})

			When("called with GKE kind and enableAutoscaler is false", func() {
				It("should be non-compliant", func() {
					plan.Spec.Configuration.Raw = []byte(`{"nodePools":[{"enableAutoscaler":false}]}`)

					result := scanner.ScanPlan(context.Background(), fake.NewFakeClient(), plan)
					Expect(result.Spec.OverallStatus).To(Equal(securityv1.Warning))
				})
			})

			When("called with GKE kind and enableAutoscaler is true", func() {
				It("should be compliant", func() {
					plan.Spec.Configuration.Raw = []byte(`{"nodePools":[{"enableAutoscaler":true}]}`)

					result := scanner.ScanPlan(context.Background(), fake.NewFakeClient(), plan)
					Expect(result.Spec.OverallStatus).To(Equal(securityv1.Compliant))
				})
			})
		})
	})
})
