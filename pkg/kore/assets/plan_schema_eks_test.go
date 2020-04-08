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

func eksSampleData() map[string]interface{} {
	return map[string]interface{}{
		"authProxyAllowedIPs":       []string{"0.0.0.0/0"},
		"description":               "This is a test cluster",
		"domain":                    "testdomain",
		"enableDefaultTrafficBlock": false,
		"inheritTeamMembers":        false,
		"region":                    "eu-west-2",
		"roleARN":                   "foo",
		"securityGroupIDs":          []string{"sg1", "sg2"},
		"subnetIDs":                 []string{"sn1", "sn2"},
		"version":                   "1.2.3",
		"nodeGroups": []map[string]interface{}{
			{
				"instanceType": "t3.medium",
				"eC2SSHKey":    "kore",
				"region":       "eu-west-2",
				"diskSize":     10,
				"name":         "group1",
				"nodeIAMRole":  "fooo",
				"desiredSize":  1,
				"minSize":      1,
				"maxSize":      10,
				"subnets":      []string{"xxx", "xxx"},
				"tags": map[string]string{
					"tag1": "value1",
				},
			},
		},
	}
}

var _ = Describe("EKSPlanSchema", func() {

	Context("The schema document", func() {
		It("is a valid JSON document", func() {
			var schema map[string]interface{}
			err := json.Unmarshal([]byte(assets.EKSPlanSchema), &schema)
			var context string
			if err != nil {
				if jsonErr, ok := err.(*json.SyntaxError); ok {
					context = assets.EKSPlanSchema[jsonErr.Offset : jsonErr.Offset+100]
				}
			}
			Expect(err).ToNot(HaveOccurred(), "error at: %s", context)
		})

		It("compiles", func() {
			Expect(func() {
				_ = jsonschema.Validate(assets.EKSPlanSchema, "test", nil)
			}).ToNot(Panic())
		})

		It("is a valid JSON Schema Draft 07 document", func() {
			var schema map[string]interface{}
			err := json.Unmarshal([]byte(assets.EKSPlanSchema), &schema)
			Expect(err).ToNot(HaveOccurred())

			err = jsonschema.Validate(assets.JSONSchemaDraft07, "test", schema)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	It("should validate the test data successfully by default", func() {
		data := eksSampleData()
		err := jsonschema.Validate(assets.EKSPlanSchema, "test", data)
		Expect(err).ToNot(HaveOccurred())
	})

	DescribeTable("Validating individual parameters",
		func(field string, value interface{}, expectError bool) {
			data := eksSampleData()
			data[field] = value
			err := jsonschema.Validate(assets.EKSPlanSchema, "test", data)
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
		Entry("region is empty", "region", "", true),
		Entry("roleARN is empty", "roleARN", "", true),
		Entry("securityGroupIDs with no elements", "securityGroupIDs", []string{}, true),
		Entry("securityGroupIDs with empty value", "securityGroupIDs", []string{""}, true),
		Entry("subnetIDs with no elements", "subnetIDs", []string{}, true),
		Entry("subnetIDs with empty value", "subnetIDs", []string{""}, true),
		Entry("version is empty", "version", "", true),
	)

	DescribeTable("Validating individual nodeGroup parameters",
		func(field string, value interface{}, expectError bool) {
			data := eksSampleData()
			data["nodeGroups"].([]map[string]interface{})[0][field] = value
			err := jsonschema.Validate(assets.EKSPlanSchema, "test", data)
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
		Entry("eC2SSHKey is empty", "eC2SSHKey", "", true),
		Entry("region is empty", "region", "", true),
		Entry("desiredSize is not a number", "desiredSize", "not a number", true),
		Entry("desiredSize is 0", "desiredSize", 0, true),
		Entry("diskSize is not a number", "diskSize", "not a number", true),
		Entry("diskSize is empty", "instanceType", "", true),
		Entry("labels with no elements", "labels", map[string]string{}, false),
		Entry("labels with empty key", "labels", map[string]string{"": "value"}, true),
		Entry("labels with invalid name", "labels", map[string]string{"!not allowed char": "value"}, true),
		Entry("name is empty", "name", "", true),
		Entry("nodeIAMRole is empty", "nodeIAMRole", "", true),
		Entry("minSize is not a number", "minSize", "not a number", true),
		Entry("minSize is 0", "minSize", 0, true),
		Entry("maxSize is not a number", "maxSize", "not a number", true),
		Entry("maxSize is 0", "maxSize", 0, true),
		Entry("releaseVersion is empty", "releaseVersion", "", true),
		Entry("sshSourceSecurityGroups with no elements", "sshSourceSecurityGroups", []string{}, false),
		Entry("sshSourceSecurityGroups with empty value", "sshSourceSecurityGroups", []string{""}, true),
		Entry("subnets with no elements", "subnets", []string{}, true),
		Entry("subnets with empty value", "subnets", []string{""}, true),
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
			validationErr = jsonschema.Validate(assets.EKSPlanSchema, "test", data)
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
