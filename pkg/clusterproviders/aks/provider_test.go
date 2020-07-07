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

package aks_test

import (
	"encoding/json"

	"github.com/appvia/kore/pkg/clusterproviders/aks"

	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/jsonschema"
	"github.com/appvia/kore/pkg/utils/validation"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func aksSampleData() map[string]interface{} {
	return map[string]interface{}{
		"authorizedMasterNetworks": []map[string]interface{}{
			{
				"name": "default",
				"cidr": "0.0.0.0/0",
			},
		},
		"authProxyAllowedIPs": []string{"0.0.0.0/0"},
		"description":         "This is a test cluster",
		"dnsPrefix":           "test",
		"domain":              "testdomain",
		"region":              "eu-west-2",
		"version":             "1.16.10",
		"networkPlugin":       "kubenet",
		"nodePools": []map[string]interface{}{
			{
				"name":             "compute",
				"mode":             "System",
				"enableAutoscaler": false,
				"diskSize":         100,
				"size":             1,
				"imageType":        "Linux",
				"machineType":      "test",
			},
		},
	}
}

var _ = Describe("Provider", func() {
	var provider kore.ClusterProvider

	BeforeEach(func() {
		provider = aks.Provider{}
	})

	Describe("PlanJSONSchema", func() {
		It("is a valid JSON document", func() {
			var schema map[string]interface{}
			err := json.Unmarshal([]byte(provider.PlanJSONSchema()), &schema)
			var context string
			if err != nil {
				if jsonErr, ok := err.(*json.SyntaxError); ok {
					context = provider.PlanJSONSchema()[jsonErr.Offset : jsonErr.Offset+100]
				}
			}
			Expect(err).ToNot(HaveOccurred(), "error at: %s", context)
		})

		It("compiles", func() {
			Expect(func() {
				_ = jsonschema.Validate(provider.PlanJSONSchema(), "test", nil)
			}).ToNot(Panic())
		})

		It("is a valid JSON Schema Draft 07 document", func() {
			var schema map[string]interface{}
			err := json.Unmarshal([]byte(provider.PlanJSONSchema()), &schema)
			Expect(err).ToNot(HaveOccurred())

			err = jsonschema.Validate(jsonschema.MetaSchemaDraft07, "test", schema)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should validate the test data successfully by default", func() {
			data := aksSampleData()
			err := jsonschema.Validate(provider.PlanJSONSchema(), "test", data)
			Expect(err).ToNot(HaveOccurred())
		})

		DescribeTable("Validating individual parameters",
			func(field string, value interface{}, expectError bool) {
				data := aksSampleData()
				data[field] = value
				err := jsonschema.Validate(provider.PlanJSONSchema(), "test", data)
				if expectError {
					Expect(err).To(HaveOccurred())
					Expect(err.(*validation.Error).FieldErrors).To(HaveLen(1))
					Expect(err.(*validation.Error).FieldErrors[0].Field).To(HavePrefix(field))
				} else {
					Expect(err).ToNot(HaveOccurred())
				}
			},
			Entry("authorizedMasterNetworks with no elements", "authProxyAllowedIPs", []map[string]interface{}{}, true),
			Entry("authorizedMasterNetworks with invalid value", "authProxyAllowedIPs", []map[string]interface{}{{"name": "foo", "cidr": "invalid"}}, true),
			Entry("authProxyAllowedIPs with no elements", "authProxyAllowedIPs", []string{}, true),
			Entry("authProxyAllowedIPs with invalid value", "authProxyAllowedIPs", []string{"invalid"}, true),
			Entry("description is empty", "description", "", true),
			Entry("dnsPrefix is empty", "dnsPrefix", "", true),
			Entry("domain is empty", "description", "", true),
			Entry("networkPlugin is empty", "networkPlugin", "", true),
			Entry("networkPlugin is invalid", "networkPlugin", "invalid", true),
			Entry("networkPolicy is valid if empty", "networkPolicy", "", false),
			Entry("networkPolicy is invalid", "networkPolicy", "invalid", true),
			Entry("nodePools with no elements", "nodePools", []map[string]interface{}{}, true),
			Entry("region is empty", "region", "", true),
			Entry("version is empty", "version", "", true),
		)

		DescribeTable("Validating individual nodePool parameters",
			func(field string, value interface{}, expectError bool) {
				data := aksSampleData()
				data["nodePools"].([]map[string]interface{})[0][field] = value
				err := jsonschema.Validate(provider.PlanJSONSchema(), "test", data)
				if expectError {
					Expect(err).To(HaveOccurred())
					Expect(err.(*validation.Error).HasErrors()).To(BeTrue())
					for i := range err.(*validation.Error).FieldErrors {
						Expect(err.(*validation.Error).FieldErrors[i].Field).To(HavePrefix("nodePools.0." + field))
					}
				} else {
					Expect(err).ToNot(HaveOccurred())
				}
			},
			Entry("name is empty", "name", "", true),
			Entry("name is invalid (starts with number)", "name", "1abc", true),
			Entry("name is invalid (ends with non-alphanumeric)", "name", "abc-", true),
			Entry("name is invalid (non-alphanumeric)", "name", "a@bc", true),
			Entry("name is invalid (more than 40 chars)", "name", "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdea", true),
			Entry("name is valid", "name", "a1", false),
			Entry("imageType is empty", "imageType", "", true),
			Entry("machineType is empty", "machineType", "", true),
			Entry("diskSize is less than minimum", "diskSize", 29, true),
			Entry("diskSize is the minimum", "diskSize", 30, false),
			Entry("diskSize is the maximum", "diskSize", 65536, false),
			Entry("diskSize is more than maximum", "diskSize", 65537, true),
			Entry("size is not a number", "size", "not a number", true),
			Entry("size is 0", "size", 0, true),
			Entry("size is 1", "size", 1, false),

			Entry("labels with no elements", "labels", map[string]string{}, false),
			Entry("labels with empty key", "labels", map[string]string{"": "value"}, true),
			Entry("labels with invalid name", "labels", map[string]string{"!not allowed char": "value"}, true),
		)

		Context("Validating multiple parameters", func() {
			var data map[string]interface{}
			var validationErr error

			BeforeEach(func() {
				data = aksSampleData()
			})

			JustBeforeEach(func() {
				validationErr = jsonschema.Validate(provider.PlanJSONSchema(), "test", data)
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
						Expect(validationErr.Error()).To(ContainSubstring("defaultTeamRole must be one of the following: \"view\", \"edit\", \"admin\", \"cluster-admin\""))
					})
				})

				When("defaultTeamRole is invalid", func() {
					BeforeEach(func() {
						data["defaultTeamRole"] = "notavalidrole"
					})
					It("throws a validation error", func() {
						Expect(validationErr).To(HaveOccurred())
						Expect(validationErr.(*validation.Error).FieldErrors).To(HaveLen(1))
						Expect(validationErr.Error()).To(ContainSubstring("defaultTeamRole must be one of the following: \"view\", \"edit\", \"admin\", \"cluster-admin\""))
					})
				})

				When("defaultTeamRole is not empty", func() {
					BeforeEach(func() {
						data["defaultTeamRole"] = "view"
					})
					It("doesn't throw an error", func() {
						Expect(validationErr).ToNot(HaveOccurred())
					})
				})
			})
		})
	})

	Describe("DefaultPlans", func() {
		It("should return with valid plans", func() {
			for _, plan := range provider.DefaultPlans() {
				err := jsonschema.Validate(provider.PlanJSONSchema(), plan.Name, plan.Spec.Configuration)
				Expect(err).ToNot(HaveOccurred(), "%s plan is not valid: %s", plan.Name, err)

			}
		})
	})
})
