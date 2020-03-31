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
	"encoding/json"
	"strings"

	"github.com/appvia/kore/pkg/apiclient"
	"github.com/appvia/kore/pkg/apiclient/models"
	"github.com/appvia/kore/pkg/apiclient/operations"
	"github.com/appvia/kore/pkg/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func stringPrt(s string) *string {
	return &s
}

var _ = Describe("/plans", func() {
	var api *apiclient.AppviaKore
	var plan *models.V1Plan
	var planName string
	var params map[string]interface{}

	BeforeEach(func() {
		api = getApi()
		planName = "test-" + strings.ToLower(utils.Random(12))
		params = map[string]interface{}{
			"authorizedMasterNetworks": []map[string]interface{}{
				{
					"name": "default",
					"cidr": "0.0.0.0/0",
				},
			},
			"authProxyAllowedIPs":           []string{"0.0.0.0/0"},
			"defaultTeamRole":               "viewer",
			"description":                   "This is a test cluster",
			"diskSize":                      100,
			"domain":                        "testing.appvia.io",
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
			"inheritTeamMembers":            true,
			"machineType":                   "n1-standard-2",
			"maintenanceWindow":             "03:00",
			"maxSize":                       10,
			"network":                       "default",
			"region":                        "europe-west2",
			"size":                          1,
			"subnetwork":                    "default",
			"version":                       "1.14.10-gke.24",
		}
	})

	JustBeforeEach(func() {
		rawConfig, _ := json.Marshal(params)

		plan = &models.V1Plan{
			Metadata: &models.V1ObjectMeta{
				Name: planName,
			},
			Spec: &models.V1PlanSpec{
				Kind:        stringPrt("GKE"),
				Summary:     stringPrt("Test plan 1"),
				Description: stringPrt("Test plan 1"),
				Labels: map[string]string{
					"kore.appvia.io/environment": "test",
					"kore.appvia.io/kind":        "GKE",
					"kore.appvia.io/plural":      "gkes",
				},
				Configuration: apiextv1.JSON{
					Raw: rawConfig,
				},
			},
		}
	})

	AfterEach(func() {
		_, _ = api.Operations.RemovePlan(operations.NewRemovePlanParams().WithName(planName), getAuth(TestUserAdmin))
	})

	When("GET /plans", func() {
		When("user is not authenticated", func() {
			It("should return 401", func() {
				_, err := api.Operations.ListPlans(operations.NewListPlansParams(), getAuthAnon())
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(&operations.ListPlansUnauthorized{}))
			})
		})

		When("user is authenticated", func() {
			It("should return a list of all plans by default", func() {
				resp, err := api.Operations.ListPlans(operations.NewListPlansParams(), getAuth(TestUserTeam1))
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Payload.Items).To(BeAssignableToTypeOf([]*models.V1Plan{}))
				Expect(len(resp.Payload.Items)).To(BeNumerically(">", 0))
			})
		})
	})

	When("GET /plans/NAME", func() {
		When("user is not authenticated", func() {
			It("should return 401", func() {
				_, err := api.Operations.GetPlan(operations.NewGetPlanParams().WithName(planName), getAuthAnon())
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(&operations.GetPlanUnauthorized{}))
			})
		})

		When("user is authenticated", func() {
			When("the plan doesn't exist", func() {
				It("should return 404", func() {
					_, err := api.Operations.GetPlan(operations.NewGetPlanParams().WithName(planName), getAuth(TestUserTeam1))
					Expect(err).To(BeAssignableToTypeOf(&operations.GetPlanNotFound{}))
				})
			})

			When("the plan exist", func() {
				BeforeEach(func() {
					params := operations.NewUpdatePlanParams().WithName(planName).WithBody(plan)
					_, err := api.Operations.UpdatePlan(params, getAuth(TestUserAdmin))
					Expect(err).ToNot(HaveOccurred())
				})

				It("should return the plan", func() {
					resp, err := api.Operations.GetPlan(operations.NewGetPlanParams().WithName(planName), getAuth(TestUserTeam1))
					Expect(err).ToNot(HaveOccurred())
					Expect(resp.Payload.Metadata.Name).To(Equal(planName))
				})
			})
		})
	})

	When("PUT /plans/NAME", func() {
		When("user is not authenticated", func() {
			It("should return 401", func() {
				params := operations.NewUpdatePlanParams().WithName(planName).WithBody(plan)
				_, err := api.Operations.UpdatePlan(params, getAuthAnon())
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(&operations.UpdatePlanUnauthorized{}))
			})
		})

		When("user is authenticated as non-admin", func() {
			It("should return 403", func() {
				params := operations.NewUpdatePlanParams().WithName(planName).WithBody(plan)
				_, err := api.Operations.UpdatePlan(params, getAuth(TestUserTeam1))
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(&operations.UpdatePlanForbidden{}))
			})
		})

		When("user is authenticated as admin", func() {
			When("there is no plan with the same name", func() {
				It("should create one", func() {
					params := operations.NewUpdatePlanParams().WithName(planName).WithBody(plan)
					resp, err := api.Operations.UpdatePlan(params, getAuth(TestUserAdmin))
					Expect(err).ToNot(HaveOccurred())
					Expect(resp.Payload.Metadata.Name).To(Equal(planName))

				})
			})

			When("there is a plan with the same name", func() {
				BeforeEach(func() {
					params := operations.NewUpdatePlanParams().WithName(planName).WithBody(plan)
					_, err := api.Operations.UpdatePlan(params, getAuth(TestUserAdmin))
					Expect(err).ToNot(HaveOccurred())
				})

				It("should update the existing one", func() {
					plan.Spec.Description = stringPrt("Updated description")
					params := operations.NewUpdatePlanParams().WithName(planName).WithBody(plan)
					_, err := api.Operations.UpdatePlan(params, getAuth(TestUserAdmin))
					Expect(err).ToNot(HaveOccurred())

					resp, err := api.Operations.GetPlan(operations.NewGetPlanParams().WithName(planName), getAuth(TestUserAdmin))
					Expect(err).ToNot(HaveOccurred())
					Expect(*resp.Payload.Spec.Description).To(Equal("Updated description"))

				})
			})

			When("the plan contains invalid parameters", func() {
				BeforeEach(func() {
					params["authorizedMasterNetworks"] = "invalid"
				})

				It("should return 400", func() {
					params := operations.NewUpdatePlanParams().WithName(planName).WithBody(plan)
					_, err := api.Operations.UpdatePlan(params, getAuth(TestUserAdmin))
					Expect(err).To(BeAssignableToTypeOf(&operations.UpdatePlanBadRequest{}))
				})
			})
		})
	})

	When("DELETE /plans/NAME", func() {
		When("user is not authenticated", func() {
			It("should return 401", func() {
				params := operations.NewRemovePlanParams().WithName(planName)
				_, err := api.Operations.RemovePlan(params, getAuthAnon())
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(&operations.RemovePlanUnauthorized{}))
			})
		})

		When("user is authenticated as non-admin", func() {
			It("should return 403", func() {
				params := operations.NewRemovePlanParams().WithName(planName)
				_, err := api.Operations.RemovePlan(params, getAuth(TestUserTeam1))
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(&operations.RemovePlanForbidden{}))
			})
		})

		When("user is authenticated as admin", func() {
			When("the plan doesn't exist", func() {
				It("should return 404", func() {
					params := operations.NewRemovePlanParams().WithName(planName)
					_, err := api.Operations.RemovePlan(params, getAuth(TestUserAdmin))
					Expect(err).To(BeAssignableToTypeOf(&operations.RemovePlanNotFound{}))

				})
			})

			When("the plan exists", func() {
				BeforeEach(func() {
					params := operations.NewUpdatePlanParams().WithName(planName).WithBody(plan)
					_, err := api.Operations.UpdatePlan(params, getAuth(TestUserAdmin))
					Expect(err).ToNot(HaveOccurred())
				})

				It("should delete the plan", func() {
					params := operations.NewRemovePlanParams().WithName(planName)
					_, err := api.Operations.RemovePlan(params, getAuth(TestUserAdmin))
					Expect(err).ToNot(HaveOccurred())

					_, err = api.Operations.GetPlan(operations.NewGetPlanParams().WithName(planName), getAuth(TestUserAdmin))
					Expect(err).To(BeAssignableToTypeOf(&operations.GetPlanNotFound{}))

				})
			})
		})
	})
})
