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
	"fmt"
	"strings"

	"github.com/appvia/kore/pkg/kore/assets"
	"github.com/appvia/kore/pkg/utils/jsonschema"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Plans", func() {
	It("All plans should be valid", func() {
		plans := assets.GetDefaultPlans()

		for _, plan := range plans {
			var schema string
			switch strings.ToLower(plan.Spec.Kind) {
			case "gke":
				schema = assets.GKEPlanSchema
			case "eks":
				schema = assets.EKSPlanSchema
			default:
				Fail(fmt.Sprintf("unknown plan kind: %q", plan.Spec.Kind))
			}
			err := jsonschema.Validate(schema, plan.Name, plan.Spec.Configuration.Raw)
			Expect(err).ToNot(HaveOccurred(), "%s plan is not valid: %s", plan.Name, err)

		}
	})

})
