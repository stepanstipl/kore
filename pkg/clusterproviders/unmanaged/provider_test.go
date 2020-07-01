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

package unmanaged_test

import (
	"encoding/json"

	"github.com/appvia/kore/pkg/clusterproviders/unmanaged"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/jsonschema"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Provider", func() {
	var provider kore.ClusterProvider

	BeforeEach(func() {
		provider = unmanaged.Provider{}
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
