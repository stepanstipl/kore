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
	"fmt"
	"strings"
	"time"

	"github.com/appvia/kore/pkg/apiclient"
	"github.com/appvia/kore/pkg/apiclient/models"
	"github.com/appvia/kore/pkg/apiclient/operations"
	"github.com/appvia/kore/pkg/apiclient/security"
	"github.com/appvia/kore/pkg/utils"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/security", func() {
	var api *apiclient.AppviaKore
	BeforeEach(func() {
		api = getApi()
	})

	Describe("/rules", func() {
		Describe("GET (ListSecurityRules)", func() {
			When("called anonymously", func() {
				It("should return 401", func() {
					_, err := api.Security.ListSecurityRules(security.NewListSecurityRulesParams(), getAuthAnon())
					Expect(err).To(HaveOccurred())
					Expect(err).To(BeAssignableToTypeOf(&security.ListSecurityRulesUnauthorized{}))
				})
			})

			When("called as user", func() {
				It("should return a list of security rules", func() {
					resp, err := api.Security.ListSecurityRules(security.NewListSecurityRulesParams(), getAuth(TestUserAdmin))
					if err != nil {
						Expect(err).ToNot(HaveOccurred())
					}
					Expect(*&resp.Payload.Items).To(BeAssignableToTypeOf([]*models.V1SecurityRule{}))
					Expect(len(resp.Payload.Items)).To(BeNumerically(">", 0))
					for _, rule := range resp.Payload.Items {
						Expect(rule.Spec).ToNot(BeNil())
						if rule.Spec.Code == "AUTHIP-01" {
							Expect(rule.Spec.Name).To(Equal("Auth Proxy IP Ranges"))
							Expect(rule.Spec.Description).ToNot(Equal(""))
							Expect(rule.Spec.AppliesTo).To(Equal([]string{"Plan"}))
						}
					}
				})
			})
		})

		Describe("GET /{code} (GetSecurityRule)", func() {
			When("called anonymously", func() {
				It("should return 401", func() {
					_, err := api.Security.GetSecurityRule(security.NewGetSecurityRuleParams().WithCode("AUTHIP-01"), getAuthAnon())
					Expect(err).To(HaveOccurred())
					Expect(err).To(BeAssignableToTypeOf(&security.GetSecurityRuleUnauthorized{}))
				})
			})

			When("called as user", func() {
				It("should return a security rule if it exists", func() {
					resp, err := api.Security.GetSecurityRule(security.NewGetSecurityRuleParams().WithCode("AUTHIP-01"), getAuth(TestUserAdmin))
					if err != nil {
						Expect(err).ToNot(HaveOccurred())
					}
					Expect(resp.Payload).To(BeAssignableToTypeOf(&models.V1SecurityRule{}))
					Expect(resp.Payload.Spec.Name).To(Equal("Auth Proxy IP Ranges"))
					Expect(resp.Payload.Spec.Description).ToNot(Equal(""))
					Expect(resp.Payload.Spec.AppliesTo).To(Equal([]string{"Plan"}))
				})

				It("should return 404 if the code does not exist", func() {
					_, err := api.Security.GetSecurityRule(security.NewGetSecurityRuleParams().WithCode("NONEXIST-01234"), getAuth(TestUserAdmin))
					Expect(err).To(HaveOccurred())
					Expect(err).To(BeAssignableToTypeOf(&security.GetSecurityRuleNotFound{}))
				})
			})
		})
	})

	Describe("/scans", func() {
		var plan *models.V1Plan
		var planName string
		var params map[string]interface{}

		BeforeEach(func() {
			planName = "secscantest-" + strings.ToLower(utils.Random(12))
			params = map[string]interface{}{
				"authorizedMasterNetworks": []map[string]interface{}{
					{
						"name": "default",
						"cidr": "0.0.0.0/0",
					},
				},
				"authProxyAllowedIPs":           []string{"0.0.0.0/0"},
				"defaultTeamRole":               "view",
				"description":                   "This is a test cluster",
				"diskSize":                      100,
				"domain":                        "testing.appvia.io",
				"enableDefaultTrafficBlock":     false,
				"enableHTTPLoadBalancer":        true,
				"enableHorizontalPodAutoscaler": true,
				"enableIstio":                   false,
				"enablePrivateEndpoint":         false,
				"enablePrivateNetwork":          false,
				"enableShieldedNodes":           true,
				"enableStackDriverLogging":      true,
				"enableStackDriverMetrics":      true,
				"inheritTeamMembers":            true,
				"maintenanceWindow":             "03:00",
				"network":                       "default",
				"region":                        "europe-west2",
				"releaseChannel":                "",
				"version":                       "1.14.10-gke.24",
				"nodePools": []map[string]interface{}{
					{
						"name":              "compute",
						"enableAutoupgrade": true,
						"version":           "",
						"enableAutoscaler":  true,
						"enableAutorepair":  true,
						"minSize":           1,
						"maxSize":           10,
						"size":              1,
						"maxPodsPerNode":    110,
						"machineType":       "n1-standard-2",
						"imageType":         "COS",
						"diskSize":          100,
						"preemptible":       false,
					},
				},
			}

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
			p := operations.NewUpdatePlanParams().WithName(planName).WithBody(plan)
			_, err := api.Operations.UpdatePlan(p, getAuth(TestUserAdmin))
			Expect(err).ToNot(HaveOccurred())
			plan.Spec.Description = stringPrt("Updated description")
			params["authProxyAllowedIPs"] = []string{"1.2.3.4/16"}
			rawConfig, _ = json.Marshal(params)
			plan.Spec.Configuration = apiextv1.JSON{Raw: rawConfig}
			p = operations.NewUpdatePlanParams().WithName(planName).WithBody(plan)
			_, err = api.Operations.UpdatePlan(p, getAuth(TestUserAdmin))
			Expect(err).ToNot(HaveOccurred())

			// Might take a little while for the scans to be produced as they're done by a background controller
			params := security.NewListSecurityScansForResourceParams().
				WithGroup("config.kore.appvia.io").
				WithVersion("v1").
				WithKind("Plan").
				WithNamespace("kore").
				WithName(planName)
			Eventually(func() int {
				resp, err := api.Security.ListSecurityScansForResource(params, getAuth(TestUserAdmin))
				if err != nil {
					return 0
				}
				return len(resp.Payload.Items)
			}, time.Second*10, time.Millisecond*200).Should(Equal(2))
		})

		AfterEach(func() {
			_, _ = api.Operations.RemovePlan(operations.NewRemovePlanParams().WithName(planName), getAuth(TestUserAdmin))
		})

		Describe("GET (ListSecurityScans)", func() {

			When("called anonymously", func() {
				It("should return 401", func() {
					_, err := api.Security.ListSecurityScans(security.NewListSecurityScansParams(), getAuthAnon())
					Expect(err).To(HaveOccurred())
					Expect(err).To(BeAssignableToTypeOf(&security.ListSecurityScansUnauthorized{}))
				})
			})

			When("called as a non-admin", func() {
				It("should return 403", func() {
					_, err := api.Security.ListSecurityScans(security.NewListSecurityScansParams(), getAuth(TestUserTeam1))
					Expect(err).To(HaveOccurred())
					Expect(err).To(BeAssignableToTypeOf(&security.ListSecurityScansForbidden{}))
				})
			})

			When("called as admin", func() {
				It("should return a list of security scans without rule results populated", func() {
					resp, err := api.Security.ListSecurityScans(security.NewListSecurityScansParams(), getAuth(TestUserAdmin))
					if err != nil {
						Expect(err).ToNot(HaveOccurred())
					}
					Expect(*&resp.Payload.Items).To(BeAssignableToTypeOf([]*models.V1SecurityScanResult{}))
					for _, scan := range resp.Payload.Items {
						Expect(len(scan.Spec.Results)).To(Equal(0))
					}
				})
			})

			When("called without latestOnly set", func() {
				It("should return a list of security scans with a null ArchivedAt date", func() {
					resp, err := api.Security.ListSecurityScans(security.NewListSecurityScansParams(), getAuth(TestUserAdmin))
					if err != nil {
						Expect(err).ToNot(HaveOccurred())
					}
					for _, scan := range resp.Payload.Items {
						Expect(scan.Spec.ArchivedAt).To(Equal(""))
					}
				})
			})

			When("called with latestOnly set false", func() {
				It("should return a list of security scans including archived scans", func() {
					f := false
					resp, err := api.Security.ListSecurityScans(security.NewListSecurityScansParams().WithLatestOnly(&f), getAuth(TestUserAdmin))
					Expect(err).ToNot(HaveOccurred())
					found := false
					for _, scan := range resp.Payload.Items {
						if scan.Spec.ArchivedAt != "" {
							found = true
						}
					}
					Expect(found).To(BeTrue())
				})
			})

		})

		Describe("GET /{group}/{version}/{kind}/{namespace}/{name} (GetSecurityScanForResource)", func() {
			It("should return the latest security scan for the resource", func() {
				params := security.NewGetSecurityScanForResourceParams().
					WithGroup("config.kore.appvia.io").
					WithVersion("v1").
					WithKind("Plan").
					WithNamespace("kore").
					WithName(planName)
				resp, err := api.Security.GetSecurityScanForResource(params, getAuth(TestUserAdmin))
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Payload.Spec.ArchivedAt).To(Equal(""))
				Expect(len(resp.Payload.Spec.Results)).ToNot(Equal(0))
				// TODO: check other fields.
			})
		})

		Describe("GET /{group}/{version}/{kind}/{namespace}/{name}/history (ListSecurityScansForResource)", func() {
			It("should return the all security scans for the resource", func() {
				params := security.NewListSecurityScansForResourceParams().
					WithGroup("config.kore.appvia.io").
					WithVersion("v1").
					WithKind("Plan").
					WithNamespace("kore").
					WithName(planName)
				resp, err := api.Security.ListSecurityScansForResource(params, getAuth(TestUserAdmin))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(resp.Payload.Items)).To(Equal(2))
				fmt.Printf("%+v", resp.Payload.Items[0].Spec)
				fmt.Printf("%+v", resp.Payload.Items[1].Spec)
				Expect(resp.Payload.Items[0].Spec.ArchivedAt).ToNot(Equal(""))
				Expect(resp.Payload.Items[1].Spec.ArchivedAt).To(Equal(""))
				Expect(len(resp.Payload.Items[0].Spec.Results)).To(Equal(0))
				Expect(len(resp.Payload.Items[1].Spec.Results)).To(Equal(0))
				// TODO: check other fields.
			})
		})

		Describe("GET /{id} (GetSecurityScan)", func() {
			It("should return the security scan by ID", func() {
				resp1, _ := api.Security.ListSecurityScans(security.NewListSecurityScansParams(), getAuth(TestUserAdmin))
				id := resp1.Payload.Items[len(resp1.Payload.Items)-1].Spec.ID
				params := security.NewGetSecurityScanParams().WithID(id)
				resp, err := api.Security.GetSecurityScan(params, getAuth(TestUserAdmin))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(resp.Payload.Spec.Results)).ToNot(Equal(0))
				// TODO: check other fields.
			})
		})

		Describe("GET /overview (GetOverview)", func() {
			It("should return an overview of the current security posture", func() {
				resp, err := api.Security.GetSecurityOverview(security.NewGetSecurityOverviewParams(), getAuth(TestUserAdmin))
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Payload).To(BeAssignableToTypeOf(&models.V1SecurityOverview{}))
			})
		})
	})
})
