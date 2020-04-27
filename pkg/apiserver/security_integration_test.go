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
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/securityscans", func() {
	var api *apiclient.AppviaKore
	var plan *models.V1Plan
	var planName string
	var params map[string]interface{}

	BeforeEach(func() {
		api = getApi()
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
		params := operations.NewUpdatePlanParams().WithName(planName).WithBody(plan)
		_, _ = api.Operations.UpdatePlan(params, getAuth(TestUserAdmin))
		plan.Spec.Description = stringPrt("Updated description")
		params = operations.NewUpdatePlanParams().WithName(planName).WithBody(plan)
		_, _ = api.Operations.UpdatePlan(params, getAuth(TestUserAdmin))
	})

	AfterEach(func() {
		_, _ = api.Operations.RemovePlan(operations.NewRemovePlanParams().WithName(planName), getAuth(TestUserAdmin))
	})

	Describe("GET (ListSecurityScans)", func() {

		When("called anonymously", func() {
			It("should return 401", func() {
				_, err := api.Operations.ListSecurityScans(operations.NewListSecurityScansParams(), getAuthAnon())
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(&operations.ListSecurityScansUnauthorized{}))
			})
		})

		When("called as a non-admin", func() {
			It("should return 403", func() {
				_, err := api.Operations.ListSecurityScans(operations.NewListSecurityScansParams(), getAuth(TestUserTeam1))
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(&operations.ListSecurityScansForbidden{}))
			})
		})

		When("called as admin", func() {
			It("should return a list of security scans without rule results populated", func() {
				resp, err := api.Operations.ListSecurityScans(operations.NewListSecurityScansParams(), getAuth(TestUserAdmin))
				if err != nil {
					Expect(err).ToNot(HaveOccurred())
				}
				Expect(*&resp.Payload.Items).To(BeAssignableToTypeOf([]*models.V1ScanResult{}))
				for _, scan := range resp.Payload.Items {
					Expect(len(scan.Spec.Results)).To(Equal(0))
				}
			})
		})

		When("called without latestOnly set", func() {
			It("should return a list of security scans with a null ArchivedAt date", func() {
				resp, err := api.Operations.ListSecurityScans(operations.NewListSecurityScansParams(), getAuth(TestUserAdmin))
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
				resp, err := api.Operations.ListSecurityScans(operations.NewListSecurityScansParams().WithLatestOnly(&f), getAuth(TestUserAdmin))
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
			params := operations.NewGetSecurityScanForResourceParams().
				WithGroup("config.kore.appvia.io").
				WithVersion("v1").
				WithKind("Plan").
				WithNamespace("kore").
				WithName(planName)
			resp, err := api.Operations.GetSecurityScanForResource(params, getAuth(TestUserAdmin))
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Payload.Spec.ArchivedAt).To(Equal(""))
			Expect(len(resp.Payload.Spec.Results)).ToNot(Equal(0))
			// TODO: check other fields.
		})
	})

	Describe("GET /{group}/{version}/{kind}/{namespace}/{name}/history (ListSecurityScansForResource)", func() {
		It("should return the all security scans for the resource", func() {
			params := operations.NewListSecurityScansForResourceParams().
				WithGroup("config.kore.appvia.io").
				WithVersion("v1").
				WithKind("Plan").
				WithNamespace("kore").
				WithName(planName)
			resp, err := api.Operations.ListSecurityScansForResource(params, getAuth(TestUserAdmin))
			Expect(err).ToNot(HaveOccurred())
			Expect(len(resp.Payload.Items)).To(Equal(2))
			Expect(resp.Payload.Items[0].Spec.ArchivedAt).ToNot(Equal(""))
			Expect(resp.Payload.Items[1].Spec.ArchivedAt).To(Equal(""))
			Expect(len(resp.Payload.Items[0].Spec.Results)).To(Equal(0))
			Expect(len(resp.Payload.Items[1].Spec.Results)).To(Equal(0))
			// TODO: check other fields.
		})
	})

	Describe("GET /{id} (GetSecurityScan)", func() {
		It("should return the security scan by ID", func() {
			resp1, _ := api.Operations.ListSecurityScans(operations.NewListSecurityScansParams(), getAuth(TestUserAdmin))
			id := resp1.Payload.Items[len(resp1.Payload.Items)-1].Spec.ID
			params := operations.NewGetSecurityScanParams().WithID(id)
			resp, err := api.Operations.GetSecurityScan(params, getAuth(TestUserAdmin))
			Expect(err).ToNot(HaveOccurred())
			Expect(len(resp.Payload.Spec.Results)).ToNot(Equal(0))
			// TODO: check other fields.
		})
	})
})
