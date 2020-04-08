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

	"github.com/appvia/kore/pkg/utils/validation"

	"github.com/appvia/kore/pkg/utils"

	"github.com/appvia/kore/pkg/apiclient"
	"github.com/appvia/kore/pkg/apiclient/models"
	"github.com/appvia/kore/pkg/apiclient/operations"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

var _ = Describe("/clusters", func() {
	var api *apiclient.AppviaKore
	var cluster *models.V1Cluster
	var clusterName, planName, teamName, namespace string
	var configuration map[string]interface{}

	var getDefaultConfiguration = func() map[string]interface{} {
		return map[string]interface{}{
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
	}

	BeforeEach(func() {
		api = getApi()
		clusterName = "test-" + strings.ToLower(utils.Random(12))
		teamName = getTestTeam(TestTeam1).Name
		namespace = teamName
		planName = "gke-development"
		configuration = getDefaultConfiguration()
	})

	JustBeforeEach(func() {
		rawValues, _ := json.Marshal(configuration)

		cluster = &models.V1Cluster{
			Metadata: &models.V1ObjectMeta{
				Name:      clusterName,
				Namespace: namespace,
			},
			Spec: &models.V1ClusterSpec{
				Kind: stringPrt("GKE"),
				Plan: stringPrt(planName),
				Configuration: apiextv1.JSON{
					Raw: rawValues,
				},
			},
		}
	})

	AfterEach(func() {
		_, _ = api.Operations.RemoveCluster(
			operations.NewRemoveClusterParams().WithName(clusterName), getAuth(TestUserAdmin),
		)
	})

	When("GET /team/TEAM/clusters", func() {
		When("user is not authenticated", func() {
			It("should return 401", func() {
				params := operations.NewListClustersParams().
					WithTeam(teamName)
				_, err := api.Operations.ListClusters(params, getAuthAnon())
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(&operations.ListClustersUnauthorized{}))
			})
		})

		When("user is authenticated", func() {
			When("the user doesn't belong to the team", func() {
				It("should return 403", func() {
					params := operations.NewListClustersParams().
						WithTeam(teamName)
					_, err := api.Operations.ListClusters(params, getAuth(TestUserTeam2))
					Expect(err).To(BeAssignableToTypeOf(&operations.ListClustersForbidden{}))
				})
			})
		})
	})

	When("GET /team/TEAM/clusters/NAME", func() {
		When("user is not authenticated", func() {
			It("should return 401", func() {
				params := operations.NewGetClusterParams().
					WithTeam(teamName).
					WithName(clusterName)
				_, err := api.Operations.GetCluster(params, getAuthAnon())
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(&operations.GetClusterUnauthorized{}))
			})
		})

		When("user is authenticated", func() {
			When("the cluster doesn't exist", func() {
				It("should return 404", func() {
					params := operations.NewGetClusterParams().
						WithTeam(teamName).
						WithName(clusterName)
					_, err := api.Operations.GetCluster(params, getAuth(TestUserTeam1))
					Expect(err).To(BeAssignableToTypeOf(&operations.GetClusterNotFound{}))
				})
			})

			When("the user doesn't belong to the team", func() {
				It("should return 403", func() {
					params := operations.NewGetClusterParams().
						WithTeam(teamName).
						WithName(clusterName)
					_, err := api.Operations.GetCluster(params, getAuth(TestUserTeam2))
					Expect(err).To(BeAssignableToTypeOf(&operations.GetClusterForbidden{}))
				})
			})
		})
	})

	When("PUT /team/TEAM/clusters/NAME", func() {
		When("user is not authenticated", func() {
			It("should return 401", func() {
				params := operations.NewUpdateClusterParams().
					WithTeam(teamName).
					WithName(clusterName).
					WithBody(cluster)
				_, err := api.Operations.UpdateCluster(params, getAuthAnon())
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(&operations.UpdateClusterUnauthorized{}))
			})
		})

		When("user is authenticated", func() {
			When("the user doesn't belong to the team", func() {
				It("should return 403", func() {
					params := operations.NewUpdateClusterParams().
						WithTeam(teamName).
						WithName(clusterName).
						WithBody(cluster)
					_, err := api.Operations.UpdateCluster(params, getAuth(TestUserTeam2))
					Expect(err).To(BeAssignableToTypeOf(&operations.UpdateClusterForbidden{}))
				})
			})

			When("the cluster contains an invalid namespace", func() {
				BeforeEach(func() {
					namespace = "not the name of the team"
				})

				It("should return 400", func() {
					params := operations.NewUpdateClusterParams().
						WithTeam(teamName).
						WithName(clusterName).
						WithBody(cluster)
					_, err := api.Operations.UpdateCluster(params, getAuth(TestUserTeam1))
					Expect(err).To(BeAssignableToTypeOf(&operations.UpdateClusterBadRequest{}))
					verr := err.(*operations.UpdateClusterBadRequest).Payload
					Expect(verr.FieldErrors).To(HaveLen(1))
					Expect(*verr.FieldErrors[0].Field).To(Equal("namespace"))
				})
			})

			When("the configuration contains an unknown parameter", func() {
				BeforeEach(func() {
					configuration["unknown-parameter"] = "foo"
				})

				It("should return 400", func() {
					params := operations.NewUpdateClusterParams().
						WithTeam(teamName).
						WithName(clusterName).
						WithBody(cluster)
					_, err := api.Operations.UpdateCluster(params, getAuth(TestUserTeam1))
					Expect(err).To(BeAssignableToTypeOf(&operations.UpdateClusterBadRequest{}))
					verr := err.(*operations.UpdateClusterBadRequest).Payload
					Expect(verr.FieldErrors).To(HaveLen(1))
					Expect(*verr.FieldErrors[0].Field).To(Equal(validation.FieldRoot))
					Expect(*verr.FieldErrors[0].Message).To(ContainSubstring("unknown-parameter"))
				})
			})

			When("the configuration contains an invalid parameter value", func() {
				BeforeEach(func() {
					configuration["diskSize"] = "this should be a number"
				})

				It("should return 400", func() {
					params := operations.NewUpdateClusterParams().
						WithTeam(teamName).
						WithName(clusterName).
						WithBody(cluster)
					_, err := api.Operations.UpdateCluster(params, getAuth(TestUserTeam1))
					Expect(err).To(BeAssignableToTypeOf(&operations.UpdateClusterBadRequest{}))
					verr := err.(*operations.UpdateClusterBadRequest).Payload
					Expect(verr.FieldErrors).To(HaveLen(1))
					Expect(*verr.FieldErrors[0].Field).To(Equal("diskSize"))
				})
			})

			When("the configuration contains a parameter which can not be modified", func() {
				BeforeEach(func() {
					configuration["enableShieldedNodes"] = false
				})

				It("should return 400", func() {
					params := operations.NewUpdateClusterParams().
						WithTeam(teamName).
						WithName(clusterName).
						WithBody(cluster)
					_, err := api.Operations.UpdateCluster(params, getAuth(TestUserTeam1))
					Expect(err).To(BeAssignableToTypeOf(&operations.UpdateClusterBadRequest{}))
					verr := err.(*operations.UpdateClusterBadRequest).Payload
					Expect(verr.FieldErrors).To(HaveLen(1))
					Expect(*verr.FieldErrors[0].Field).To(Equal("enableShieldedNodes"))
					Expect(*verr.FieldErrors[0].Message).To(Equal("can not be changed"))
				})
			})
		})
	})

	When("DELETE /team/TEAM/clusters/NAME", func() {
		When("user is not authenticated", func() {
			It("should return 401", func() {
				params := operations.NewRemoveClusterParams().
					WithTeam(teamName).
					WithName(clusterName)
				_, err := api.Operations.RemoveCluster(params, getAuthAnon())
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(&operations.RemoveClusterUnauthorized{}))
			})
		})

		When("user is authenticated", func() {
			When("the cluster doesn't exist", func() {
				It("should return 404", func() {
					params := operations.NewRemoveClusterParams().
						WithTeam(teamName).
						WithName(clusterName)
					_, err := api.Operations.RemoveCluster(params, getAuth(TestUserTeam1))
					Expect(err).To(BeAssignableToTypeOf(&operations.RemoveClusterNotFound{}))
				})
			})

			When("the user doesn't belong to the team", func() {
				It("should return 403", func() {
					params := operations.NewRemoveClusterParams().
						WithTeam(teamName).
						WithName(clusterName)
					_, err := api.Operations.RemoveCluster(params, getAuth(TestUserTeam2))
					Expect(err).To(BeAssignableToTypeOf(&operations.RemoveClusterForbidden{}))
				})
			})
		})
	})
})
