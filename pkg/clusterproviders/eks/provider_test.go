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

package eks_test

import (
	"encoding/json"

	"github.com/appvia/kore/pkg/clusterproviders/eks"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/jsonschema"
	"github.com/appvia/kore/pkg/utils/validation"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func eksSampleData() map[string]interface{} {
	return map[string]interface{}{
		"authorizedMasterNetworks":  []string{"0.0.0.0/0"},
		"authProxyAllowedIPs":       []string{"0.0.0.0/0"},
		"description":               "This is a test cluster",
		"domain":                    "testdomain",
		"enableDefaultTrafficBlock": false,
		"inheritTeamMembers":        false,
		"privateIPV4Cidr":           "10.0.0.0/16",
		"region":                    "eu-west-2",
		"version":                   "1.15",
		"nodeGroups": []map[string]interface{}{
			{
				"instanceType":     "t3.medium",
				"diskSize":         10,
				"name":             "group1",
				"enableAutoscaler": false,
				"desiredSize":      1,
				"minSize":          1,
				"maxSize":          10,
				"tags": map[string]string{
					"tag1": "value1",
				},
			},
		},
	}
}

var _ = Describe("Provider", func() {
	var provider kore.ClusterProvider

	BeforeEach(func() {
		provider = eks.Provider{}
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
			data := eksSampleData()
			err := jsonschema.Validate(provider.PlanJSONSchema(), "test", data)
			Expect(err).ToNot(HaveOccurred())
		})

		DescribeTable("Validating individual parameters",
			func(field string, value interface{}, expectError bool) {
				data := eksSampleData()
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
			Entry("authProxyAllowedIPs with no elements", "authProxyAllowedIPs", []string{}, true),
			Entry("authProxyAllowedIPs with invalid value", "authProxyAllowedIPs", []string{"invalid"}, true),
			Entry("description is empty", "description", "", true),
			Entry("domain is empty", "description", "", true),
			Entry("nodeGroups with no elements", "nodeGroups", []map[string]interface{}{}, true),
			Entry("privateIPV4Cidr is empty", "privateIPV4Cidr", "", true),
			Entry("privateIPV4Cidr is invalid", "privateIPV4Cidr", "non-cidr", true),
			Entry("region is empty", "region", "", true),
			Entry("version is empty", "version", "", true),
		)

		DescribeTable("Validating individual nodeGroup parameters",
			func(field string, value interface{}, expectError bool) {
				data := eksSampleData()
				data["nodeGroups"].([]map[string]interface{})[0][field] = value
				err := jsonschema.Validate(provider.PlanJSONSchema(), "test", data)
				if expectError {
					Expect(err).To(HaveOccurred())
					Expect(err.(*validation.Error).HasErrors()).To(BeTrue())
					for i := range err.(*validation.Error).FieldErrors {
						Expect(err.(*validation.Error).FieldErrors[i].Field).To(HavePrefix("nodeGroups.0." + field))
					}
				} else {
					Expect(err).ToNot(HaveOccurred())
				}
			},
			Entry("amiType is empty", "amiType", "", true),
			Entry("instanceType is empty", "instanceType", "", true),
			Entry("desiredSize is not a number", "desiredSize", "not a number", true),
			Entry("desiredSize is 0", "desiredSize", 0, true),
			Entry("diskSize is not a number", "diskSize", "not a number", true),
			Entry("diskSize is empty", "instanceType", "", true),
			Entry("labels with no elements", "labels", map[string]string{}, false),
			Entry("labels with empty key", "labels", map[string]string{"": "value"}, true),
			Entry("labels with invalid name", "labels", map[string]string{"!not allowed char": "value"}, true),
			Entry("name is empty", "name", "", true),
			Entry("minSize is not a number", "minSize", "not a number", true),
			Entry("minSize is 0", "minSize", 0, true),
			Entry("maxSize is not a number", "maxSize", "not a number", true),
			Entry("maxSize is 0", "maxSize", 0, true),
			Entry("releaseVersion is empty", "releaseVersion", "", false),
			Entry("releaseVersion is invalid", "releaseVersion", "1.15.", true),
			Entry("releaseVersion is valid", "releaseVersion", "1.15.8-20200102", false),
			Entry("sshSourceSecurityGroups with no elements", "sshSourceSecurityGroups", []string{}, false),
			Entry("sshSourceSecurityGroups with empty value", "sshSourceSecurityGroups", []string{""}, true),
			Entry("tags with no elements", "tags", map[string]string{}, false),
			Entry("tags with empty key", "tags", map[string]string{"": "value"}, true),
			Entry("tags with invalid name", "tags", map[string]string{"!not allowed char": "value"}, true),
		)

		Context("Validating multiple parameters", func() {
			var data map[string]interface{}
			var validationErr error

			BeforeEach(func() {
				data = eksSampleData()
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
