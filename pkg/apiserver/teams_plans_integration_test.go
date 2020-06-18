// +build integration

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

package apiserver_test

import (
	"github.com/appvia/kore/pkg/apiclient"
	"github.com/appvia/kore/pkg/apiclient/models"
	"github.com/appvia/kore/pkg/apiclient/operations"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GET /teams/{team}/plans/{plan} (GetTeamPlanDetails)", func() {
	var api *apiclient.AppviaKore
	var testTeam1 string

	BeforeEach(func() {
		api = getApi()
		testTeam1 = getTestTeam(TestTeam1).Name
	})

	When("called anonymously", func() {
		It("should return 401", func() {
			_, err := api.Operations.GetTeamPlanDetails(operations.NewGetTeamPlanDetailsParams().WithTeam(testTeam1).WithPlan("gke-development"), getAuthAnon())
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(&operations.GetTeamPlanDetailsUnauthorized{}))
		})
	})

	When("called as a non-admin", func() {
		It("should return 403 if not in the team in question", func() {
			_, err := api.Operations.GetTeamPlanDetails(operations.NewGetTeamPlanDetailsParams().WithTeam(testTeam1).WithPlan("gke-development"), getAuth(TestUserTeam2))
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(&operations.GetTeamPlanDetailsForbidden{}))
		})
		It("should return details if called by a user in the team", func() {
			resp, err := api.Operations.GetTeamPlanDetails(operations.NewGetTeamPlanDetailsParams().WithTeam(testTeam1).WithPlan("gke-development"), getAuth(TestUserTeam1))
			if err != nil {
				Expect(err).ToNot(HaveOccurred())
			}
			Expect(*&resp.Payload).To(BeAssignableToTypeOf(&models.ApiserverTeamPlan{}))
		})
	})

	When("called as admin", func() {
		It("should return plan details", func() {
			resp, err := api.Operations.GetTeamPlanDetails(operations.NewGetTeamPlanDetailsParams().WithTeam(testTeam1).WithPlan("gke-development"), getAuth(TestUserAdmin))
			if err != nil {
				Expect(err).ToNot(HaveOccurred())
			}
			Expect(*&resp.Payload).To(BeAssignableToTypeOf(&models.ApiserverTeamPlan{}))
		})
	})

	When("called for a non-existant plan", func() {
		It("should return 404", func() {
			_, err := api.Operations.GetTeamPlanDetails(operations.NewGetTeamPlanDetailsParams().WithTeam(testTeam1).WithPlan("gke-development-nonexistant"), getAuth(TestUserTeam1))
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(&operations.GetTeamPlanDetailsNotFound{}))
		})
	})

	When("a plan and team exist", func() {
		It("should return the spec of the plan", func() {
			resp, err := api.Operations.GetTeamPlanDetails(operations.NewGetTeamPlanDetailsParams().WithTeam(testTeam1).WithPlan("gke-development"), getAuth(TestUserAdmin))
			if err != nil {
				Expect(err).ToNot(HaveOccurred())
			}
			Expect(*&resp.Payload.Plan).To(BeAssignableToTypeOf(&models.V1PlanSpec{}))
			Expect(*&resp.Payload.Plan.Configuration).ToNot(BeNil())
		})

		It("should return the JSON schema relevant for the plan as a string", func() {
			resp, err := api.Operations.GetTeamPlanDetails(operations.NewGetTeamPlanDetailsParams().WithTeam(testTeam1).WithPlan("gke-development"), getAuth(TestUserAdmin))
			if err != nil {
				Expect(err).ToNot(HaveOccurred())
			}
			Expect(*&resp.Payload.Schema).ToNot(Equal(""))
		})

		It("should return the list of editable parameters for this team", func() {
			resp, err := api.Operations.GetTeamPlanDetails(operations.NewGetTeamPlanDetailsParams().WithTeam(testTeam1).WithPlan("gke-development"), getAuth(TestUserAdmin))
			if err != nil {
				Expect(err).ToNot(HaveOccurred())
			}
			Expect(*&resp.Payload.EditableParams).To(BeAssignableToTypeOf([]string{}))
			Expect(*&resp.Payload.EditableParams).ToNot(Equal([]string{}))
		})
	})
})
