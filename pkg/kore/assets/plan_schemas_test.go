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

package assets_test

import (
	"encoding/json"

	"github.com/appvia/kore/pkg/utils/validation"

	"github.com/appvia/kore/pkg/kore/assets"
	"github.com/appvia/kore/pkg/utils/jsonschema"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func sampleData() map[string]interface{} {
	return map[string]interface{}{
		"authorizedMasterNetworks": []map[string]interface{}{
			{
				"name": "default",
				"cidr": "0.0.0.0/0",
			},
		},
		"authProxyAllowedIPs":           []string{"0.0.0.0/0"},
		"description":                   "This is a test cluster",
		"diskSize":                      100,
		"domain":                        "testdomain",
		"enableAutoupgrade":             true,
		"enableAutorepair":              true,
		"enableAutoscaler":              true,
		"enableDefaultTrafficBlock":     false,
		"enableHTTPLoadBalancer":        true,
		"enableHorizontalPodAutoscaler": true,
		"enableIstio":                   false,
		"enablePrivateEndpoint":         false,
		"enablePrivateNetwork":          false,
		"enableShieldedNodes":           true,
		"enableStackDriverLogging":      true,
		"enableStackDriverMetrics":      true,
		"imageType":                     "COS",
		"inheritTeamMembers":            false,
		"machineType":                   "n1-standard-2",
		"maintenanceWindow":             "03:00",
		"maxSize":                       10,
		"network":                       "default",
		"region":                        "europe-west2",
		"size":                          1,
		"subnetwork":                    "default",
		"version":                       "1.14.10-gke.24",
	}
}

var nonEmptyFields = []string{
	"description", "domain", "imageType", "machineType", "network", "region", "subnetwork", "version",
}

var _ = Describe("GKEPlanSchema", func() {

	Context("The schema document", func() {
		It("is a valid JSON document", func() {
			var schema map[string]interface{}
			err := json.Unmarshal([]byte(assets.GKEPlanSchema), &schema)
			var context string
			if err != nil {
				if jsonErr, ok := err.(*json.SyntaxError); ok {
					context = assets.GKEPlanSchema[jsonErr.Offset : jsonErr.Offset+100]
				}
			}
			Expect(err).ToNot(HaveOccurred(), "error at: %s", context)
		})

		It("compiles", func() {
			Expect(func() {
				_ = jsonschema.Validate(assets.GKEPlanSchema, "test", nil)
			}).ToNot(Panic())
		})

		It("is a valid JSON Schema Draft 07 document", func() {
			var schema map[string]interface{}
			err := json.Unmarshal([]byte(assets.GKEPlanSchema), &schema)
			Expect(err).ToNot(HaveOccurred())

			err = jsonschema.Validate(assets.JSONSchemaDraft07, "test", schema)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	It("should validate the test data successfully by default", func() {
		data := sampleData()
		err := jsonschema.Validate(assets.GKEPlanSchema, "test", data)
		Expect(err).ToNot(HaveOccurred())
	})

	DescribeTable("Validating individual parameters",
		func(field string, value interface{}, expectError bool) {
			data := sampleData()
			data[field] = value
			err := jsonschema.Validate(assets.GKEPlanSchema, "test", data)
			if expectError {
				Expect(err).To(HaveOccurred())
				Expect(err.(*validation.Error).FieldErrors).To(HaveLen(1))
				Expect(err.(*validation.Error).FieldErrors[0].Field).To(HavePrefix(field))
			} else {
				Expect(err).ToNot(HaveOccurred())
			}
		},
		Entry("description is empty", "description", "", true),
		Entry("domain is empty", "description", "", true),
		Entry("imageType is empty", "description", "", true),
		Entry("machineType is empty", "description", "", true),
		Entry("network is empty", "description", "", true),
		Entry("region is empty", "description", "", true),
		Entry("subnetwork is empty", "description", "", true),
		Entry("version is empty", "description", "", true),
		Entry("diskSize is less than minimum", "diskSize", 9, true),
		Entry("diskSize is the minimum", "diskSize", 10, false),
		Entry("diskSize is the maximum", "diskSize", 65536, false),
		Entry("diskSize is more than maximum", "diskSize", 65537, true),
		Entry("authorizedMasterNetworks with no elements", "authorizedMasterNetworks", []map[string]interface{}{}, true),
		Entry("authorizedMasterNetworks with missing cidr", "authorizedMasterNetworks", []map[string]interface{}{
			{"name": "xx"},
		}, true),
		Entry("authorizedMasterNetworks with missing name", "authorizedMasterNetworks", []map[string]interface{}{
			{"cidr": "1.2.3.4/32"},
		}, true),
		Entry("authorizedMasterNetworks with invalid cidr", "authorizedMasterNetworks", []map[string]interface{}{
			{"name": "xx", "cidr": "invalid"},
		}, true),
		Entry("authProxyAllowedIPs with no elements", "authProxyAllowedIPs", []string{}, true),
		Entry("authProxyAllowedIPs with invalid value", "authProxyAllowedIPs", []string{"invalid"}, true),
	)

	Context("Validating multiple parameters", func() {
		var data map[string]interface{}
		var validationErr error

		BeforeEach(func() {
			data = sampleData()
		})

		JustBeforeEach(func() {
			validationErr = jsonschema.Validate(assets.GKEPlanSchema, "test", data)
		})

		When("inheritTeamMembers is true", func() {
			BeforeEach(func() {
				data["inheritTeamMembers"] = true
			})

			When("defaultTeamRole is not set", func() {
				It("throws a validation error", func() {
					Expect(validationErr).To(HaveOccurred())
					Expect(validationErr.(*validation.Error).FieldErrors).To(HaveLen(1))
					Expect(validationErr.Error()).To(ContainSubstring("defaultTeamRole is required"))
				})
			})

			When("defaultTeamRole is empty", func() {
				BeforeEach(func() {
					data["defaultTeamRole"] = ""
				})
				It("throws a validation error", func() {
					Expect(validationErr).To(HaveOccurred())
					Expect(validationErr.(*validation.Error).FieldErrors).To(HaveLen(1))
					Expect(validationErr.Error()).To(ContainSubstring("defaultTeamRole: String length must be greater than or equal to 1"))
				})
			})

			When("defaultTeamRole is not empty", func() {
				BeforeEach(func() {
					data["defaultTeamRole"] = "some role"
				})
				It("doesn't throw an error", func() {
					Expect(validationErr).ToNot(HaveOccurred())
				})
			})
		})
	})

})
