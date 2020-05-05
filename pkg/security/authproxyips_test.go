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
	"encoding/json"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	securityv1 "github.com/appvia/kore/pkg/apis/security/v1"
	"github.com/appvia/kore/pkg/security"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getConfig(ipRanges []string) *configv1.Plan {
	config := map[string]interface{}{
		"authProxyAllowedIPs": ipRanges,
	}
	configBytes, err := json.Marshal(config)
	Expect(err).ToNot(HaveOccurred())
	return &configv1.Plan{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "config.kore.appvia.io/v1",
			Kind:       "Plan",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "test1",
		},
		Spec: configv1.PlanSpec{
			Configuration: apiextv1.JSON{
				Raw: configBytes,
			},
		},
	}
}

var _ = Describe("AuthProxyIPsRule", func() {
	scanner := security.NewEmpty()

	BeforeSuite(func() {
		scanner.RegisterRule(&security.AuthProxyIPRangeRule{})
	})

	Describe("CheckPlan", func() {

		When("called with a plan which contains no IP ranges", func() {
			It("should return a warning", func() {
				result := scanner.ScanPlan(getConfig([]string{}))
				Expect(result.Spec.OverallStatus).To(Equal(securityv1.Warning))
				Expect(result.Spec.Results[0].Status).To(Equal(securityv1.Warning))
			})
		})

		When("called with a plan which contains IP ranges", func() {
			It("should return a warning if an IP range specifies a wider CIDR than /16", func() {
				result := scanner.ScanPlan(getConfig([]string{"1.2.3.4/15", "2.3.4.5/16"}))
				Expect(result.Spec.OverallStatus).To(Equal(securityv1.Warning))
				Expect(result.Spec.Results[0].Status).To(Equal(securityv1.Warning))
			})

			It("should return compliant if all IP ranges specify /16 or narrower ranges", func() {
				result := scanner.ScanPlan(getConfig([]string{"1.2.3.4/16", "2.3.4.5/16"}))
				Expect(result.Spec.OverallStatus).To(Equal(securityv1.Compliant))
				Expect(result.Spec.Results[0].Status).To(Equal(securityv1.Compliant))
			})
		})

	})
})
