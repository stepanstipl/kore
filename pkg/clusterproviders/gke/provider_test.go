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

package gke_test

import (
	"encoding/json"

	"github.com/appvia/kore/pkg/clusterproviders/gke"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/jsonschema"
	"github.com/appvia/kore/pkg/utils/validation"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func gkeSampleData() map[string]interface{} {
	return map[string]interface{}{
		"authorizedMasterNetworks": []map[string]interface{}{
			{
				"name": "default",
				"cidr": "0.0.0.0/0",
			},
		},
		"authProxyAllowedIPs":           []string{"0.0.0.0/0"},
		"description":                   "This is a test cluster",
		"domain":                        "testdomain",
		"enableDefaultTrafficBlock":     false,
		"enableHTTPLoadBalancer":        true,
		"enableHorizontalPodAutoscaler": true,
		"enableIstio":                   false,
		"enablePrivateEndpoint":         false,
		"enablePrivateNetwork":          false,
		"enableShieldedNodes":           true,
		"enableStackDriverLogging":      true,
		"enableStackDriverMetrics":      true,
		"nodePools": []map[string]interface{}{
			{
				"name":              "compute",
				"enableAutoupgrade": true,
				"enableAutoscaler":  false,
				"enableAutorepair":  true,
				"diskSize":          100,
				"size":              1,
				"imageType":         "COS",
				"machineType":       "n1-standard-2",
				"version":           "1.14.10-gke.24",
			},
		},
		"inheritTeamMembers": false,
		"maintenanceWindow":  "03:00",
		"region":             "europe-west2",
		"version":            "1.14.10-gke.24",
		"releaseChannel":     "",
	}
}

func gkeSampleDataDeprecated() map[string]interface{} {
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

var _ = Describe("Provider", func() {
	var provider kore.ClusterProvider

	BeforeEach(func() {
		provider = gke.Provider{}
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
			data := gkeSampleData()
			err := jsonschema.Validate(provider.PlanJSONSchema(), "test", data)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should validate the deprecated data successfully by default", func() {
			data := gkeSampleDataDeprecated()
			err := jsonschema.Validate(provider.PlanJSONSchema(), "test", data)
			Expect(err).ToNot(HaveOccurred())
		})

		DescribeTable("Validating individual parameters",
			func(field string, value interface{}, expectError bool) {
				data := gkeSampleData()
				data[field] = value
				err := jsonschema.Validate(provider.PlanJSONSchema(), "test", data)
				if expectError {
					Expect(err).To(HaveOccurred())
					Expect(err.(*validation.Error).FieldErrors[0].Field).To(HavePrefix(field))
				} else {
					Expect(err).ToNot(HaveOccurred())
				}
			},
			Entry("description is empty", "description", "", true),
			Entry("domain is empty", "domain", "", true),
			Entry("region is empty", "region", "", true),
			Entry("version is empty", "version", "", true),
			Entry("version is invalid", "version", "horse", true),
			Entry("version is invalid - 1.x.y-abc.z", "version", "1.15.1-abc.1", true),
			Entry("version is valid - current", "version", "-", false),
			Entry("version is valid - latest", "version", "latest", false),
			Entry("version is valid - 1.x", "version", "1.15", false),
			Entry("version is valid - 1.x.y", "version", "1.15.1", false),
			Entry("version is valid - 1.x.y-gke.z", "version", "1.15.1-gke.1", false),
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
			Entry("nodePools with no elements", "nodePools", []map[string]interface{}{}, true),
		)

		DescribeTable("Validating individual nodePool parameters",
			func(field string, value interface{}, expectError bool) {
				data := gkeSampleData()
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
			Entry("machineType does not match machine type pattern", "machineType", "horse", true),
			Entry("diskSize is less than minimum", "diskSize", 9, true),
			Entry("diskSize is the minimum", "diskSize", 10, false),
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
				data = gkeSampleData()
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
						Expect(validationErr.(*validation.Error).FieldErrors).To(HaveLen(2))
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

				When("defaultTeamRole is not empty", func() {
					BeforeEach(func() {
						data["defaultTeamRole"] = "cluster-admin"
					})
					It("doesn't throw an error", func() {
						Expect(validationErr).ToNot(HaveOccurred())
					})
				})
			})

			When("node pool auto-scale is true", func() {
				BeforeEach(func() {
					data["nodePools"].([]map[string]interface{})[0]["enableAutoscaler"] = true
				})

				When("minSize and maxSize are not set", func() {
					It("throws a validation error", func() {
						Expect(validationErr).To(HaveOccurred())
						Expect(validationErr.Error()).To(ContainSubstring("nodePools.0: minSize is required"))
						Expect(validationErr.Error()).To(ContainSubstring("nodePools.0: maxSize is required"))
					})
				})

				When("just minSize is set", func() {
					BeforeEach(func() {
						data["nodePools"].([]map[string]interface{})[0]["minSize"] = 1
					})
					It("throws a validation error", func() {
						Expect(validationErr).To(HaveOccurred())
						Expect(validationErr.Error()).To(ContainSubstring("nodePools.0: maxSize is required"))
					})
				})

				When("just maxSize is set", func() {
					BeforeEach(func() {
						data["nodePools"].([]map[string]interface{})[0]["maxSize"] = 10
					})
					It("throws a validation error", func() {
						Expect(validationErr).To(HaveOccurred())
						Expect(validationErr.Error()).To(ContainSubstring("nodePools.0: minSize is required"))
					})
				})

				When("minSize and maxSize are both set", func() {
					BeforeEach(func() {
						data["nodePools"].([]map[string]interface{})[0]["minSize"] = 1
						data["nodePools"].([]map[string]interface{})[0]["maxSize"] = 10
					})
					It("doesn't throw an error", func() {
						Expect(validationErr).ToNot(HaveOccurred())
					})
				})
			})

			When("releaseChannel is set", func() {
				BeforeEach(func() {
					data["releaseChannel"] = "REGULAR"
					data["version"] = ""
					data["nodePools"].([]map[string]interface{})[0]["version"] = ""
					data["nodePools"].([]map[string]interface{})[0]["enableAutoupgrade"] = true
				})

				When("version and node pool version are not set and node pool auto-upgrade is true", func() {
					It("doesn't throw an error", func() {
						Expect(validationErr).ToNot(HaveOccurred())
					})
				})

				When("version is set", func() {
					BeforeEach(func() {
						data["version"] = "1.14"
					})

					It("throws a validation error", func() {
						Expect(validationErr).To(HaveOccurred())
						Expect(validationErr.Error()).To(ContainSubstring("version: version does not match: \"\""))
					})
				})

				When("nodePool version is set", func() {
					BeforeEach(func() {
						data["nodePools"].([]map[string]interface{})[0]["version"] = "1.14"
					})

					It("throws a validation error", func() {
						Expect(validationErr).To(HaveOccurred())
						Expect(validationErr.Error()).To(ContainSubstring("nodePools.0.version: nodePools.0.version does not match: \"\""))
					})
				})

				When("nodePool auto-upgrade is disabled", func() {
					BeforeEach(func() {
						data["nodePools"].([]map[string]interface{})[0]["enableAutoupgrade"] = false
					})

					It("throws a validation error", func() {
						Expect(validationErr).To(HaveOccurred())
						Expect(validationErr.Error()).To(ContainSubstring("nodePools.0.enableAutoupgrade: nodePools.0.enableAutoupgrade does not match: true"))
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
