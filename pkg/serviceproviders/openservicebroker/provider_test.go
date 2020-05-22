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

package openservicebroker_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/appvia/kore/pkg/kore/korefakes"

	"github.com/appvia/kore/pkg/controllers/controllersfakes"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/serviceproviders/openservicebroker"
	"github.com/appvia/kore/pkg/serviceproviders/openservicebroker/openservicebrokerfakes"
	"github.com/appvia/kore/pkg/utils"

	osb "github.com/kubernetes-sigs/go-open-service-broker-client/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	ProviderName = "test-osb-provider"
	Service1Name = "service-1"
	Service1ID   = "service-1-uuid"
	Service2Name = "service-2"
	Service2ID   = "service-2-uuid"
	plan1Name    = "plan-1"
	plan1ID      = "plan-1-uuid"
	plan2Name    = "plan-2"
	plan2ID      = "plan-2-uuid"
	plan3Name    = "plan-3"
	plan3ID      = "plan-3-uuid"

	KoreServiceName            = "kore-service"
	KoreServiceID              = "kore-service-uuid"
	KoreServiceCredentialsName = "kore-service-credentials"
	KoreServiceCredentialsID   = "kore-service-credentials-uuid"
	Namespace                  = "test"
)

var Operation = osb.OperationKey("test-op")

func schema(name string) string {
	return fmt.Sprintf(`
		{
			"$id": "https://test.appvia.io/schemas/plan.json",
			"$schema": "http://json-schema.org/draft-07/schema#",
			"description": "Test plan schema",
			"type": "object",
			"additionalProperties": false,
			"required": [
				"%s-param"
			],
			"properties": {
				"%s-param": {
					"type": "string",
					"minLength": 1,
					"default": "%s-value"
				}
			}
		}`, name, name, name,
	)
}

func credsSchema(name string) string {
	return fmt.Sprintf(`
		{
			"$id": "https://test.appvia.io/schemas/credentials.json",
			"$schema": "http://json-schema.org/draft-07/schema#",
			"description": "Test plan credentials schema",
			"type": "object",
			"additionalProperties": false,
			"required": [
				"%s-credentials-param"
			],
			"properties": {
				"%s-credentials-param": {
					"type": "string",
					"minLength": 1,
					"default": "%s-credentials-value"
				}
			}
		}`, name, name, name,
	)
}

func createService(id string, name string, plans []osb.Plan) osb.Service {
	return osb.Service{
		ID:                   id,
		Name:                 name,
		Description:          name + " description",
		Bindable:             true,
		InstancesRetrievable: true,
		BindingsRetrievable:  true,
		PlanUpdatable:        utils.BoolPtr(true),
		Plans:                plans,
	}
}

func createPlan(id string, name string) osb.Plan {
	return osb.Plan{
		ID:          id,
		Name:        name,
		Description: name + " description",
		Bindable:    utils.BoolPtr(true),
		Schemas: &osb.Schemas{
			ServiceInstance: &osb.ServiceInstanceSchema{
				Create: &osb.InputParametersSchema{
					Parameters: schema(name),
				},
			},
			ServiceBinding: &osb.ServiceBindingSchema{
				Create: &osb.InputParametersSchema{
					Parameters: credsSchema(name),
				},
			},
		},
	}
}

func createProviderData(operation *osb.OperationKey) *apiextv1.JSON {
	data := openservicebroker.ProviderData{Operation: operation}
	res, _ := json.Marshal(data)
	return &apiextv1.JSON{Raw: res}
}

var _ = Describe("Provider", func() {
	var client *openservicebrokerfakes.FakeClient
	var provider *openservicebroker.Provider
	var ctx context.Context
	var serviceProviderCtx kore.Context
	var cancel context.CancelFunc
	var logger *log.Logger
	var serviceProvider *servicesv1.ServiceProvider
	var service *servicesv1.Service
	var serviceCreds *servicesv1.ServiceCredentials
	var reconcileResult reconcile.Result
	var reconcileErr error
	var providerConfig openservicebroker.ProviderConfiguration
	var servicePlans *korefakes.FakeServicePlans
	var serviceKinds *korefakes.FakeServiceKinds

	var expectToNotRequeue = func() {
		Expect(reconcileErr).ToNot(HaveOccurred())
		Expect(reconcileResult).To(Equal(reconcile.Result{}))
	}

	var expectToRequeue = func() {
		Expect(reconcileErr).ToNot(HaveOccurred())
		if !reconcileResult.Requeue && reconcileResult.RequeueAfter == 0 {
			Fail("expected to requeue")
		}
	}

	var expectError = func(msg string) {
		Expect(reconcileErr).To(HaveOccurred())
		Expect(reconcileErr.Error()).To(ContainSubstring(msg))
	}

	var expectCriticalError = func(msg string) {
		Expect(reconcileErr).To(HaveOccurred())
		Expect(reconcileErr.Error()).To(ContainSubstring(msg))
		if !controllers.IsCriticalError(reconcileErr) {
			Fail(fmt.Sprintf("was expecting critical error, got %v", reconcileErr))
		}
	}

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
		logger = log.StandardLogger()
		logger.Out = GinkgoWriter

		providerConfig = openservicebroker.ProviderConfiguration{}

		serviceProvider = servicesv1.NewServiceProvider("test", "testns")

		servicePlans = &korefakes.FakeServicePlans{}
		serviceKinds = &korefakes.FakeServiceKinds{}
		korei := &korefakes.FakeInterface{}
		korei.ServicePlansReturns(servicePlans)
		korei.ServiceKindsReturns(serviceKinds)

		serviceProviderCtx = kore.NewContext(ctx, logger, &controllersfakes.FakeClient{}, korei)

		client = &openservicebrokerfakes.FakeClient{}
		client.GetCatalogReturns(&osb.CatalogResponse{
			Services: []osb.Service{
				createService(Service1ID, Service1Name, []osb.Plan{
					createPlan(plan1ID, plan1Name),
					createPlan(plan2ID, plan2Name),
				}),
				createService(Service2ID, Service2Name, []osb.Plan{
					createPlan(plan3ID, plan3Name),
				}),
			},
		}, nil)

		service = servicesv1.NewService(KoreServiceName, Namespace)
		service.Spec = servicesv1.ServiceSpec{
			Kind: Service1Name,
			Plan: Service1Name + "-" + plan1Name,
		}
		service.Spec.Configuration = &apiextv1.JSON{Raw: []byte(`{"serviceParam":"foo"}`)}

		serviceCreds = servicesv1.NewServiceCredentials(KoreServiceCredentialsName, Namespace)
		serviceCreds.Spec = servicesv1.ServiceCredentialsSpec{
			Kind:    Service1Name,
			Service: service.Ownership(),
		}
		serviceCreds.Spec.Configuration = &apiextv1.JSON{Raw: []byte(`{"credsParam":"foo"}`)}

		reconcileResult = reconcile.Result{}
		reconcileErr = nil

		servicePlans.GetReturns(&servicesv1.ServicePlan{
			ObjectMeta: metav1.ObjectMeta{
				Name: Service1Name + "-" + plan1Name,
			},
			Spec: servicesv1.ServicePlanSpec{
				ProviderData: &apiextv1.JSON{Raw: []byte(`{"planID":"plan-1-uuid","serviceID":"service-1-uuid"}`)},
			},
		}, nil)
	})

	AfterEach(func() {
		cancel()
	})

	JustBeforeEach(func() {
		provider = openservicebroker.NewProvider(ProviderName, providerConfig, client)
	})

	When("creating a new provider", func() {
		var catalogErr error
		var catalog kore.ServiceProviderCatalog

		JustBeforeEach(func() {
			catalog, catalogErr = provider.Catalog(serviceProviderCtx, serviceProvider)
		})

		It("should fetch the catalog successfully", func() {
			Expect(catalogErr).ToNot(HaveOccurred())
		})

		It("should create service plans from the catalog plans", func() {
			Expect(catalog.Plans[0]).To(Equal(servicesv1.ServicePlan{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ServicePlan",
					APIVersion: servicesv1.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "service-1-plan-1",
					Namespace: kore.HubNamespace,
				},
				Spec: servicesv1.ServicePlanSpec{
					Kind:             "service-1",
					Description:      "plan-1 description",
					Summary:          "service-1 service - plan-1 plan",
					Schema:           schema("plan-1"),
					CredentialSchema: credsSchema("plan-1"),
					Configuration:    &apiextv1.JSON{Raw: []byte(`{"plan-1-param":"plan-1-value"}`)},
					ProviderData:     &apiextv1.JSON{Raw: []byte(`{"planID":"plan-1-uuid","serviceID":"service-1-uuid"}`)},
				},
			}))
			Expect(catalog.Plans[1]).To(Equal(servicesv1.ServicePlan{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ServicePlan",
					APIVersion: servicesv1.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "service-1-plan-2",
					Namespace: kore.HubNamespace,
				},
				Spec: servicesv1.ServicePlanSpec{
					Kind:             "service-1",
					Description:      "plan-2 description",
					Summary:          "service-1 service - plan-2 plan",
					Schema:           schema("plan-2"),
					CredentialSchema: credsSchema("plan-2"),
					Configuration:    &apiextv1.JSON{Raw: []byte(`{"plan-2-param":"plan-2-value"}`)},
					ProviderData:     &apiextv1.JSON{Raw: []byte(`{"planID":"plan-2-uuid","serviceID":"service-1-uuid"}`)},
				},
			}))
			Expect(catalog.Plans[2]).To(Equal(servicesv1.ServicePlan{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ServicePlan",
					APIVersion: servicesv1.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "service-2-plan-3",
					Namespace: kore.HubNamespace,
				},
				Spec: servicesv1.ServicePlanSpec{
					Kind:             "service-2",
					Description:      "plan-3 description",
					Summary:          "service-2 service - plan-3 plan",
					Schema:           schema("plan-3"),
					CredentialSchema: credsSchema("plan-3"),
					Configuration:    &apiextv1.JSON{Raw: []byte(`{"plan-3-param":"plan-3-value"}`)},
					ProviderData:     &apiextv1.JSON{Raw: []byte(`{"planID":"plan-3-uuid","serviceID":"service-2-uuid"}`)},
				},
			}))
		})

		When("the catalog endpoint returns an error", func() {
			BeforeEach(func() {
				client.GetCatalogReturns(nil, fmt.Errorf("some error"))
			})

			It("should error", func() {
				Expect(catalogErr).To(MatchError("failed to fetch catalog from service broker: some error"))
			})
		})

		When("a default plan doesn't have a schema", func() {
			BeforeEach(func() {
				providerConfig.DefaultPlans = []string{Service1Name + "-" + plan1Name}
				plan := createPlan(plan1ID, plan1Name)
				plan.Schemas.ServiceInstance = nil
				service := createService(Service1ID, Service1Name, []osb.Plan{plan})
				client.GetCatalogReturns(&osb.CatalogResponse{
					Services: []osb.Service{service},
				}, nil)
			})

			It("should error", func() {
				Expect(catalogErr).To(MatchError("service-1-plan-1 plan does not have a schema for provisioning"))
			})
		})

		When("a default plan doesn't have a credentials schema", func() {
			BeforeEach(func() {
				providerConfig.DefaultPlans = []string{Service1Name + "-" + plan1Name}
				plan := createPlan(plan1ID, plan1Name)
				plan.Schemas.ServiceBinding = nil
				service := createService(Service1ID, Service1Name, []osb.Plan{plan})
				client.GetCatalogReturns(&osb.CatalogResponse{
					Services: []osb.Service{service},
				}, nil)
			})

			It("should error", func() {
				Expect(catalogErr).To(MatchError("service-1-plan-1 plan does not have a schema for bind"))
			})

			When("it is allowed", func() {
				BeforeEach(func() {
					providerConfig.AllowEmptyCredentialSchema = true
				})

				It("should not error", func() {
					Expect(catalogErr).ToNot(HaveOccurred())
				})
			})
		})

		When("there are multiple default plans for the same service", func() {
			BeforeEach(func() {
				providerConfig.DefaultPlans = []string{Service1Name + "-" + plan1Name, Service1Name + "-" + plan2Name}
			})

			It("should error", func() {
				Expect(catalogErr).To(MatchError("there are multiple default plans for the same service: service-1"))
			})
		})

		When("a plan has a configuration key in the metadata", func() {
			BeforeEach(func() {
				plan := createPlan(plan1ID, plan1Name)
				plan.Metadata = map[string]interface{}{
					openservicebroker.MetadataKeyConfiguration: map[string]interface{}{
						"extraKey": "extraValue",
					},
				}
				service := createService(Service1ID, Service1Name, []osb.Plan{plan})
				client.GetCatalogReturns(&osb.CatalogResponse{
					Services: []osb.Service{service},
				}, nil)
			})

			It("should not error", func() {
				Expect(catalogErr).ToNot(HaveOccurred())
			})

			It("should create service plans from the catalog plans", func() {
				catalog, err := provider.Catalog(serviceProviderCtx, serviceProvider)
				Expect(err).ToNot(HaveOccurred())
				config := map[string]interface{}{}
				_ = catalog.Plans[0].Spec.GetConfiguration(&config)
				Expect(config).To(HaveKeyWithValue("plan-1-param", "plan-1-value"))
				Expect(config).To(HaveKeyWithValue("extraKey", "extraValue"))
			})
		})

		When("a plan has an invalid configuration key in the metadata", func() {
			BeforeEach(func() {
				plan := createPlan(plan1ID, plan1Name)
				plan.Metadata = map[string]interface{}{
					openservicebroker.MetadataKeyConfiguration: "invalid",
				}
				service := createService(Service1ID, Service1Name, []osb.Plan{plan})
				client.GetCatalogReturns(&osb.CatalogResponse{
					Services: []osb.Service{service},
				}, nil)
			})

			It("should error", func() {
				Expect(catalogErr).To(MatchError("service-1-plan-1 plan has an invalid configuration, it must be an object"))
			})
		})
	})

	Context("Reconcile", func() {
		JustBeforeEach(func() {
			reconcileResult, reconcileErr = provider.Reconcile(serviceProviderCtx, service)
		})

		When("the service does not exist", func() {

			BeforeEach(func() {
				client.ProvisionInstanceReturns(&osb.ProvisionResponse{}, nil)
			})

			It("should call provision", func() {
				Expect(client.ProvisionInstanceCallCount()).To(Equal(1))

				req := client.ProvisionInstanceArgsForCall(0)
				Expect(req.AcceptsIncomplete).To(BeTrue())
				Expect(req.InstanceID).ToNot(BeEmpty(), "service instance id is empty")
				Expect(req.ServiceID).To(Equal(Service1ID))
				Expect(req.PlanID).To(Equal(plan1ID))
				Expect(json.Marshal(req.Parameters)).To(Equal(service.Spec.Configuration.Raw))
			})

			When("the response has async=false", func() {
				BeforeEach(func() {
					client.ProvisionInstanceReturns(&osb.ProvisionResponse{
						Async: false,
					}, nil)
				})

				It("should set the component status to success", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentProvision)
					Expect(component.Status).To(Equal(corev1.SuccessStatus))
				})

				It("should not requeue", func() {
					expectToNotRequeue()
				})
			})

			When("the response has async=true", func() {
				BeforeEach(func() {
					client.ProvisionInstanceReturns(&osb.ProvisionResponse{
						Async:        true,
						OperationKey: &Operation,
					}, nil)
				})

				It("should leave the component status as pending", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentProvision)
					Expect(component.Status).To(Equal(corev1.PendingStatus))
				})

				It("should save the operation data", func() {
					Expect(service.Status.ProviderData).To(Equal(createProviderData(&Operation)))
				})

				It("should requeue", func() {
					expectToRequeue()
				})
			})

			When("the response is an error", func() {
				BeforeEach(func() {
					client.ProvisionInstanceReturns(nil, osb.HTTPStatusCodeError{
						StatusCode:   http.StatusInternalServerError,
						ErrorMessage: utils.StringPtr("some error"),
					})
				})

				It("should set the component status to error", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentProvision)
					Expect(component.Status).To(Equal(corev1.ErrorStatus))
					Expect(component.Message).To(Equal("failed to call provision on the service broker"))
					Expect(component.Detail).To(ContainSubstring("some error"))
				})

				It("should return the error", func() {
					expectError("some error")
				})
			})

			When("the response is bad request", func() {
				BeforeEach(func() {
					client.ProvisionInstanceReturns(nil, osb.HTTPStatusCodeError{
						StatusCode:   http.StatusBadRequest,
						ErrorMessage: utils.StringPtr("some error"),
					})
				})

				It("should set the component status to failure", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentProvision)
					Expect(component.Status).To(Equal(corev1.FailureStatus))
					Expect(component.Message).To(Equal("failed to call provision on the service broker"))
					Expect(component.Detail).To(ContainSubstring("some error"))
				})

				It("should return the error", func() {
					expectCriticalError("some error")
				})
			})

		})

		When("the service creation is pending with async=true", func() {

			BeforeEach(func() {
				service.Status.Components.SetCondition(corev1.Component{
					Name:   openservicebroker.ComponentProvision,
					Status: corev1.PendingStatus,
				})
				service.Status.ProviderID = KoreServiceID
				service.Status.ProviderData = createProviderData(&Operation)

				client.PollLastOperationReturns(&osb.LastOperationResponse{
					State: osb.StateInProgress,
				}, nil)
			})

			It("should poll the last operation", func() {
				Expect(client.PollLastOperationCallCount()).To(Equal(1))

				req := client.PollLastOperationArgsForCall(0)
				Expect(req.InstanceID).To(Equal(KoreServiceID))
				Expect(*req.ServiceID).To(Equal(Service1ID))
				Expect(*req.PlanID).To(Equal(plan1ID))
				Expect(req.OperationKey).To(Equal(&Operation))
			})

			When("the status is succeeded", func() {
				BeforeEach(func() {
					client.PollLastOperationReturns(&osb.LastOperationResponse{
						State: osb.StateSucceeded,
					}, nil)
				})

				It("should set the component status to success", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentProvision)
					Expect(component.Status).To(Equal(corev1.SuccessStatus))
				})

				It("should not requeue", func() {
					expectToNotRequeue()
				})
			})

			When("the status is failed", func() {
				BeforeEach(func() {
					client.PollLastOperationReturns(&osb.LastOperationResponse{
						State:       osb.StateFailed,
						Description: utils.StringPtr("some error"),
					}, nil)
				})

				It("should set the component status to failure", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentProvision)
					Expect(component.Status).To(Equal(corev1.FailureStatus))
				})

				It("should return an error", func() {
					expectCriticalError("last operation failed on the service broker: some error")
				})
			})

			When("the status is in progress", func() {
				BeforeEach(func() {
					client.PollLastOperationReturns(&osb.LastOperationResponse{
						State: osb.StateInProgress,
					}, nil)
				})

				It("should leave the status in pending", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentProvision)
					Expect(component.Status).To(Equal(corev1.PendingStatus))
				})

				It("should requeue", func() {
					expectToRequeue()
				})
			})

			When("the response is an error", func() {
				BeforeEach(func() {
					client.PollLastOperationReturns(nil, osb.HTTPStatusCodeError{
						StatusCode:   http.StatusInternalServerError,
						ErrorMessage: utils.StringPtr("some error"),
					})
				})

				It("should set the component status to error", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentProvision)
					Expect(component.Status).To(Equal(corev1.ErrorStatus))
					Expect(component.Message).To(Equal("failed to poll last operation on the service broker"))
					Expect(component.Detail).To(ContainSubstring("some error"))
				})

				It("should return the error", func() {
					expectError("some error")
				})
			})

			When("the response is bad request", func() {
				BeforeEach(func() {
					client.PollLastOperationReturns(nil, osb.HTTPStatusCodeError{
						StatusCode:   http.StatusBadRequest,
						ErrorMessage: utils.StringPtr("some error"),
					})
				})

				It("should set the component status to failure", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentProvision)
					Expect(component.Status).To(Equal(corev1.FailureStatus))
					Expect(component.Message).To(Equal("failed to poll last operation on the service broker"))
					Expect(component.Detail).To(ContainSubstring("some error"))
				})

				It("should return the error", func() {
					expectCriticalError("some error")
				})
			})

		})

		When("the service was successfully created", func() {

			BeforeEach(func() {
				service.Status.Components.SetCondition(corev1.Component{
					Name:   openservicebroker.ComponentProvision,
					Status: corev1.SuccessStatus,
				})
				service.Status.ProviderID = KoreServiceID

				client.UpdateInstanceReturns(&osb.UpdateInstanceResponse{}, nil)
			})

			When("service plan and configuration did not change", func() {
				BeforeEach(func() {
					service.Status.Plan = service.Spec.Plan
					service.Status.Configuration = service.Spec.Configuration
				})

				It("should not call update", func() {
					Expect(client.UpdateInstanceCallCount()).To(Equal(0))
				})

				It("should not requeue", func() {
					expectToNotRequeue()
				})
			})

			When("the configuration did change", func() {
				BeforeEach(func() {
					service.Spec.Configuration = &apiextv1.JSON{Raw: []byte(`{"foo":"bar"}`)}
				})

				It("should call update", func() {
					Expect(client.UpdateInstanceCallCount()).To(Equal(1))

					req := client.UpdateInstanceArgsForCall(0)
					Expect(req.AcceptsIncomplete).To(BeTrue())
					Expect(req.InstanceID).To(Equal(KoreServiceID))
					Expect(req.ServiceID).To(Equal(Service1ID))
					Expect(json.Marshal(req.Parameters)).To(Equal(service.Spec.Configuration.Raw))
				})

				When("the response has async=false", func() {
					BeforeEach(func() {
						client.UpdateInstanceReturns(&osb.UpdateInstanceResponse{
							Async: false,
						}, nil)
					})

					It("should set the component status to success", func() {
						component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentUpdate)
						Expect(component.Status).To(Equal(corev1.SuccessStatus))
					})

					It("should not requeue", func() {
						expectToNotRequeue()
					})
				})

				When("the response has async=true", func() {
					BeforeEach(func() {
						client.UpdateInstanceReturns(&osb.UpdateInstanceResponse{
							Async:        true,
							OperationKey: &Operation,
						}, nil)
					})

					It("should leave the component status as pending", func() {
						component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentUpdate)
						Expect(component.Status).To(Equal(corev1.PendingStatus))
					})

					It("should save the operation data", func() {
						Expect(service.Status.ProviderData).To(Equal(createProviderData(&Operation)))
					})

					It("should requeue", func() {
						expectToRequeue()
					})
				})

				When("the response is an error", func() {
					BeforeEach(func() {
						client.UpdateInstanceReturns(nil, osb.HTTPStatusCodeError{
							StatusCode:   http.StatusInternalServerError,
							ErrorMessage: utils.StringPtr("some error"),
						})
					})

					It("should set the component status to error", func() {
						component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentUpdate)
						Expect(component.Status).To(Equal(corev1.ErrorStatus))
						Expect(component.Message).To(Equal("failed to call update on the service broker"))
						Expect(component.Detail).To(ContainSubstring("some error"))
					})

					It("should return the error", func() {
						expectError("some error")
					})
				})

				When("the response is bad request", func() {
					BeforeEach(func() {
						client.UpdateInstanceReturns(nil, osb.HTTPStatusCodeError{
							StatusCode:   http.StatusBadRequest,
							ErrorMessage: utils.StringPtr("some error"),
						})
					})

					It("should set the component status to failure", func() {
						component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentUpdate)
						Expect(component.Status).To(Equal(corev1.FailureStatus))
						Expect(component.Message).To(Equal("failed to call update on the service broker"))
						Expect(component.Detail).To(ContainSubstring("some error"))
					})

					It("should return the error", func() {
						expectCriticalError("some error")
					})
				})

			})

		})

		When("the service update is pending with async=true", func() {

			BeforeEach(func() {
				service.Status.Components.SetCondition(corev1.Component{
					Name:   openservicebroker.ComponentProvision,
					Status: corev1.SuccessStatus,
				})
				service.Status.Components.SetCondition(corev1.Component{
					Name:   openservicebroker.ComponentUpdate,
					Status: corev1.PendingStatus,
				})
				service.Status.ProviderID = KoreServiceID
				service.Status.ProviderData = createProviderData(&Operation)

				client.PollLastOperationReturns(&osb.LastOperationResponse{
					State: osb.StateInProgress,
				}, nil)
			})

			It("should poll the last operation", func() {
				Expect(client.PollLastOperationCallCount()).To(Equal(1))

				req := client.PollLastOperationArgsForCall(0)
				Expect(req.InstanceID).To(Equal(KoreServiceID))
				Expect(*req.ServiceID).To(Equal(Service1ID))
				Expect(*req.PlanID).To(Equal(plan1ID))
				Expect(req.OperationKey).To(Equal(&Operation))
			})

			When("the status is succeeded", func() {
				BeforeEach(func() {
					client.PollLastOperationReturns(&osb.LastOperationResponse{
						State: osb.StateSucceeded,
					}, nil)
				})

				It("should set the component status to success", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentUpdate)
					Expect(component.Status).To(Equal(corev1.SuccessStatus))
				})

				It("should not requeue", func() {
					expectToNotRequeue()
				})
			})

			When("the status is failed", func() {
				BeforeEach(func() {
					client.PollLastOperationReturns(&osb.LastOperationResponse{
						State:       osb.StateFailed,
						Description: utils.StringPtr("some error"),
					}, nil)
				})

				It("should set the component status to failure", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentUpdate)
					Expect(component.Status).To(Equal(corev1.FailureStatus))
				})

				It("should return an error", func() {
					expectCriticalError("last operation failed on the service broker: some error")
				})
			})

			When("the status is in progress", func() {
				BeforeEach(func() {
					client.PollLastOperationReturns(&osb.LastOperationResponse{
						State: osb.StateInProgress,
					}, nil)
				})

				It("should leave the status in pending", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentUpdate)
					Expect(component.Status).To(Equal(corev1.PendingStatus))
				})

				It("should requeue", func() {
					expectToRequeue()
				})
			})

			When("the response is an error", func() {
				BeforeEach(func() {
					client.PollLastOperationReturns(nil, osb.HTTPStatusCodeError{
						StatusCode:   http.StatusInternalServerError,
						ErrorMessage: utils.StringPtr("some error"),
					})
				})

				It("should set the component status to error", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentUpdate)
					Expect(component.Status).To(Equal(corev1.ErrorStatus))
					Expect(component.Message).To(Equal("failed to poll last operation on the service broker"))
					Expect(component.Detail).To(ContainSubstring("some error"))
				})

				It("should return the error", func() {
					expectError("some error")
				})
			})

			When("the response is bad request", func() {
				BeforeEach(func() {
					client.PollLastOperationReturns(nil, osb.HTTPStatusCodeError{
						StatusCode:   http.StatusBadRequest,
						ErrorMessage: utils.StringPtr("some error"),
					})
				})

				It("should set the component status to failure", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentUpdate)
					Expect(component.Status).To(Equal(corev1.FailureStatus))
					Expect(component.Message).To(Equal("failed to poll last operation on the service broker"))
					Expect(component.Detail).To(ContainSubstring("some error"))
				})

				It("should return the error", func() {
					expectCriticalError("some error")
				})
			})

		})

	})

	Context("Delete", func() {
		BeforeEach(func() {
			service.Status.ProviderID = KoreServiceID
		})

		JustBeforeEach(func() {
			reconcileResult, reconcileErr = provider.Delete(serviceProviderCtx, service)
		})

		When("the service exists", func() {

			BeforeEach(func() {
				client.DeprovisionInstanceReturns(&osb.DeprovisionResponse{}, nil)
			})

			It("should call provision", func() {
				Expect(client.DeprovisionInstanceCallCount()).To(Equal(1))

				req := client.DeprovisionInstanceArgsForCall(0)
				Expect(req.AcceptsIncomplete).To(BeTrue())
				Expect(req.InstanceID).To(Equal(KoreServiceID))
				Expect(req.ServiceID).To(Equal(Service1ID))
				Expect(req.PlanID).To(Equal(plan1ID))
			})

			When("the response has async=false", func() {
				BeforeEach(func() {
					client.DeprovisionInstanceReturns(&osb.DeprovisionResponse{
						Async: false,
					}, nil)
				})

				It("should set the component status to success", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentDeprovision)
					Expect(component.Status).To(Equal(corev1.SuccessStatus))
				})

				It("should not requeue", func() {
					expectToNotRequeue()
				})
			})

			When("the response has async=true", func() {
				BeforeEach(func() {
					client.DeprovisionInstanceReturns(&osb.DeprovisionResponse{
						Async:        true,
						OperationKey: &Operation,
					}, nil)
				})

				It("should leave the component status as pending", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentDeprovision)
					Expect(component.Status).To(Equal(corev1.PendingStatus))
				})

				It("should save the operation data", func() {
					Expect(service.Status.ProviderData).To(Equal(createProviderData(&Operation)))
				})

				It("should requeue", func() {
					expectToRequeue()
				})
			})

			When("the response is an error", func() {
				BeforeEach(func() {
					client.DeprovisionInstanceReturns(nil, osb.HTTPStatusCodeError{
						StatusCode:   http.StatusInternalServerError,
						ErrorMessage: utils.StringPtr("some error"),
					})
				})

				It("should set the component status to error", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentDeprovision)
					Expect(component.Status).To(Equal(corev1.ErrorStatus))
					Expect(component.Message).To(Equal("failed to call deprovision on the service broker"))
					Expect(component.Detail).To(ContainSubstring("some error"))
				})

				It("should return the error", func() {
					expectError("some error")
				})
			})

			When("the response is bad request", func() {
				BeforeEach(func() {
					client.DeprovisionInstanceReturns(nil, osb.HTTPStatusCodeError{
						StatusCode:   http.StatusBadRequest,
						ErrorMessage: utils.StringPtr("some error"),
					})
				})

				It("should set the component status to failure", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentDeprovision)
					Expect(component.Status).To(Equal(corev1.FailureStatus))
					Expect(component.Message).To(Equal("failed to call deprovision on the service broker"))
					Expect(component.Detail).To(ContainSubstring("some error"))
				})

				It("should return the error", func() {
					expectCriticalError("some error")
				})
			})

		})

		When("the service deletion is pending with async=true", func() {

			BeforeEach(func() {
				service.Status.Components.SetCondition(corev1.Component{
					Name:   openservicebroker.ComponentDeprovision,
					Status: corev1.PendingStatus,
				})
				service.Status.ProviderID = KoreServiceID
				service.Status.ProviderData = createProviderData(&Operation)

				client.PollLastOperationReturns(&osb.LastOperationResponse{
					State: osb.StateInProgress,
				}, nil)
			})

			It("should poll the last operation", func() {
				Expect(client.PollLastOperationCallCount()).To(Equal(1))

				req := client.PollLastOperationArgsForCall(0)
				Expect(req.InstanceID).To(Equal(KoreServiceID))
				Expect(*req.ServiceID).To(Equal(Service1ID))
				Expect(*req.PlanID).To(Equal(plan1ID))
				Expect(req.OperationKey).To(Equal(&Operation))
			})

			When("the status is succeeded", func() {
				BeforeEach(func() {
					client.PollLastOperationReturns(&osb.LastOperationResponse{
						State: osb.StateSucceeded,
					}, nil)
				})

				It("should set the component status to success", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentDeprovision)
					Expect(component.Status).To(Equal(corev1.SuccessStatus))
				})

				It("should not requeue", func() {
					expectToNotRequeue()
				})
			})

			When("the status is failed", func() {
				BeforeEach(func() {
					client.PollLastOperationReturns(&osb.LastOperationResponse{
						State:       osb.StateFailed,
						Description: utils.StringPtr("some error"),
					}, nil)
				})

				It("should set the component status to failure", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentDeprovision)
					Expect(component.Status).To(Equal(corev1.FailureStatus))
				})

				It("should return an error", func() {
					expectCriticalError("last operation failed on the service broker: some error")
				})
			})

			When("the status is in progress", func() {
				BeforeEach(func() {
					client.PollLastOperationReturns(&osb.LastOperationResponse{
						State: osb.StateInProgress,
					}, nil)
				})

				It("should leave the status in pending", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentDeprovision)
					Expect(component.Status).To(Equal(corev1.PendingStatus))
				})

				It("should requeue", func() {
					expectToRequeue()
				})
			})

			When("the response is an error", func() {
				BeforeEach(func() {
					client.PollLastOperationReturns(nil, osb.HTTPStatusCodeError{
						StatusCode:   http.StatusInternalServerError,
						ErrorMessage: utils.StringPtr("some error"),
					})
				})

				It("should set the component status to error", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentDeprovision)
					Expect(component.Status).To(Equal(corev1.ErrorStatus))
					Expect(component.Message).To(Equal("failed to poll last operation on the service broker"))
					Expect(component.Detail).To(ContainSubstring("some error"))
				})

				It("should return the error", func() {
					expectError("some error")
				})
			})

			When("the response is not found", func() {
				BeforeEach(func() {
					client.PollLastOperationReturns(nil, osb.HTTPStatusCodeError{
						StatusCode: http.StatusNotFound,
					})
				})

				It("should set the component status to success", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentDeprovision)
					Expect(component.Status).To(Equal(corev1.SuccessStatus))
				})

				It("should not requeue", func() {
					expectToNotRequeue()
				})
			})

			When("the response is bad request", func() {
				BeforeEach(func() {
					client.PollLastOperationReturns(nil, osb.HTTPStatusCodeError{
						StatusCode:   http.StatusBadRequest,
						ErrorMessage: utils.StringPtr("some error"),
					})
				})

				It("should set the component status to failure", func() {
					component, _ := service.Status.Components.GetComponent(openservicebroker.ComponentDeprovision)
					Expect(component.Status).To(Equal(corev1.FailureStatus))
					Expect(component.Message).To(Equal("failed to poll last operation on the service broker"))
					Expect(component.Detail).To(ContainSubstring("some error"))
				})

				It("should return the error", func() {
					expectCriticalError("some error")
				})
			})

		})
	})

	Context("ReconcileCredentials", func() {
		var secrets map[string]string

		BeforeEach(func() {
			service.Status.ProviderID = KoreServiceID
		})

		JustBeforeEach(func() {
			reconcileResult, secrets, reconcileErr = provider.ReconcileCredentials(serviceProviderCtx, service, serviceCreds)
		})

		When("the service credentials do not exist", func() {

			BeforeEach(func() {
				client.BindReturns(&osb.BindResponse{}, nil)
			})

			It("should call bind", func() {
				Expect(client.BindCallCount()).To(Equal(1))

				req := client.BindArgsForCall(0)
				Expect(req.AcceptsIncomplete).To(BeTrue())
				Expect(req.BindingID).ToNot(BeEmpty(), "binding id is empty")
				Expect(req.InstanceID).To(Equal(KoreServiceID))
				Expect(req.ServiceID).To(Equal(Service1ID))
				Expect(req.PlanID).To(Equal(plan1ID))
				Expect(json.Marshal(req.Parameters)).To(Equal(serviceCreds.Spec.Configuration.Raw))
			})

			When("the response has async=false", func() {
				BeforeEach(func() {
					client.BindReturns(&osb.BindResponse{
						Async:       false,
						Credentials: map[string]interface{}{"secret": 42},
					}, nil)
				})

				It("should set the component status to success", func() {
					component, _ := serviceCreds.Status.Components.GetComponent(openservicebroker.ComponentBind)
					Expect(component.Status).To(Equal(corev1.SuccessStatus))
				})

				It("should return the secret", func() {
					Expect(secrets).To(Equal(map[string]string{"secret": "42"}))
				})

				It("should not requeue", func() {
					expectToNotRequeue()
				})
			})

			When("the response has async=true", func() {
				BeforeEach(func() {
					client.BindReturns(&osb.BindResponse{
						Async:        true,
						OperationKey: &Operation,
					}, nil)
				})

				It("should leave the component status as pending", func() {
					component, _ := serviceCreds.Status.Components.GetComponent(openservicebroker.ComponentBind)
					Expect(component.Status).To(Equal(corev1.PendingStatus))
				})

				It("should save the operation data", func() {
					Expect(serviceCreds.Status.ProviderData).To(Equal(createProviderData(&Operation)))
				})

				It("should requeue", func() {
					expectToRequeue()
				})
			})

			When("the response is an error", func() {
				BeforeEach(func() {
					client.BindReturns(nil, osb.HTTPStatusCodeError{
						StatusCode:   http.StatusInternalServerError,
						ErrorMessage: utils.StringPtr("some error"),
					})
				})

				It("should set the component status to error", func() {
					component, _ := serviceCreds.Status.Components.GetComponent(openservicebroker.ComponentBind)
					Expect(component.Status).To(Equal(corev1.ErrorStatus))
					Expect(component.Message).To(Equal("failed to call bind on the service broker"))
					Expect(component.Detail).To(ContainSubstring("some error"))
				})

				It("should return the error", func() {
					expectError("some error")
				})
			})

			When("the response is bad request", func() {
				BeforeEach(func() {
					client.BindReturns(nil, osb.HTTPStatusCodeError{
						StatusCode:   http.StatusBadRequest,
						ErrorMessage: utils.StringPtr("some error"),
					})
				})

				It("should set the component status to failure", func() {
					component, _ := serviceCreds.Status.Components.GetComponent(openservicebroker.ComponentBind)
					Expect(component.Status).To(Equal(corev1.FailureStatus))
					Expect(component.Message).To(Equal("failed to call bind on the service broker"))
					Expect(component.Detail).To(ContainSubstring("some error"))
				})

				It("should return the error", func() {
					expectCriticalError("some error")
				})
			})

		})

		When("the service credentials are pending with async=true", func() {

			BeforeEach(func() {
				service.Status.ProviderID = KoreServiceID

				serviceCreds.Status.Components.SetCondition(corev1.Component{
					Name:   openservicebroker.ComponentBind,
					Status: corev1.PendingStatus,
				})
				serviceCreds.Status.ProviderID = KoreServiceCredentialsID
				serviceCreds.Status.ProviderData = createProviderData(&Operation)

				client.PollBindingLastOperationReturns(&osb.LastOperationResponse{
					State: osb.StateInProgress,
				}, nil)
			})

			It("should poll the last bind operation", func() {
				Expect(client.PollBindingLastOperationCallCount()).To(Equal(1))

				req := client.PollBindingLastOperationArgsForCall(0)
				Expect(req.InstanceID).To(Equal(KoreServiceID))
				Expect(req.BindingID).To(Equal(KoreServiceCredentialsID))
				Expect(*req.ServiceID).To(Equal(Service1ID))
				Expect(*req.PlanID).To(Equal(plan1ID))
				Expect(req.OperationKey).To(Equal(&Operation))
			})

			When("the status is succeeded", func() {
				BeforeEach(func() {
					client.PollBindingLastOperationReturns(&osb.LastOperationResponse{
						State: osb.StateSucceeded,
					}, nil)

					client.GetBindingReturns(&osb.GetBindingResponse{}, nil)
				})

				It("should get the binding details", func() {
					Expect(client.GetBindingCallCount()).To(Equal(1))
					req := client.GetBindingArgsForCall(0)
					Expect(req.BindingID).To(Equal(KoreServiceCredentialsID))
					Expect(req.InstanceID).To(Equal(KoreServiceID))
				})

				When("the binding details are returned successfully", func() {
					BeforeEach(func() {
						client.GetBindingReturns(&osb.GetBindingResponse{
							Credentials: map[string]interface{}{"secret": 42},
						}, nil)
					})

					It("should set the component status to success", func() {
						component, _ := serviceCreds.Status.Components.GetComponent(openservicebroker.ComponentBind)
						Expect(component.Status).To(Equal(corev1.SuccessStatus))
					})

					It("should return the secret", func() {
						Expect(secrets).To(Equal(map[string]string{"secret": "42"}))
					})

					It("should not requeue", func() {
						expectToNotRequeue()
					})
				})

				When("the response is an error", func() {
					BeforeEach(func() {
						client.GetBindingReturns(nil, osb.HTTPStatusCodeError{
							StatusCode:   http.StatusInternalServerError,
							ErrorMessage: utils.StringPtr("some error"),
						})
					})

					It("should set the component status to error", func() {
						component, _ := serviceCreds.Status.Components.GetComponent(openservicebroker.ComponentBind)
						Expect(component.Status).To(Equal(corev1.ErrorStatus))
						Expect(component.Message).To(Equal("failed to get binding details from the service broker"))
						Expect(component.Detail).To(ContainSubstring("some error"))
					})

					It("should return the error", func() {
						expectError("some error")
					})
				})

				When("the response is bad request", func() {
					BeforeEach(func() {
						client.GetBindingReturns(nil, osb.HTTPStatusCodeError{
							StatusCode:   http.StatusBadRequest,
							ErrorMessage: utils.StringPtr("some error"),
						})
					})

					It("should set the component status to failure", func() {
						component, _ := serviceCreds.Status.Components.GetComponent(openservicebroker.ComponentBind)
						Expect(component.Status).To(Equal(corev1.FailureStatus))
						Expect(component.Message).To(Equal("failed to get binding details from the service broker"))
						Expect(component.Detail).To(ContainSubstring("some error"))
					})

					It("should return the error", func() {
						expectCriticalError("some error")
					})
				})
			})

			When("the status is failed", func() {
				BeforeEach(func() {
					client.PollBindingLastOperationReturns(&osb.LastOperationResponse{
						State:       osb.StateFailed,
						Description: utils.StringPtr("some error"),
					}, nil)
				})

				It("should set the bind status to failure", func() {
					component, _ := serviceCreds.Status.Components.GetComponent(openservicebroker.ComponentBind)
					Expect(component.Status).To(Equal(corev1.FailureStatus))
				})

				It("should return an error", func() {
					expectCriticalError("last bind operation failed on the service broker: some error")
				})
			})

			When("the status is in progress", func() {
				BeforeEach(func() {
					client.PollBindingLastOperationReturns(&osb.LastOperationResponse{
						State: osb.StateInProgress,
					}, nil)
				})

				It("should leave the status in pending", func() {
					component, _ := serviceCreds.Status.Components.GetComponent(openservicebroker.ComponentBind)
					Expect(component.Status).To(Equal(corev1.PendingStatus))
				})

				It("should requeue", func() {
					expectToRequeue()
				})
			})

			When("the response is an error", func() {
				BeforeEach(func() {
					client.PollBindingLastOperationReturns(nil, osb.HTTPStatusCodeError{
						StatusCode:   http.StatusInternalServerError,
						ErrorMessage: utils.StringPtr("some error"),
					})
				})

				It("should set the component status to error", func() {
					component, _ := serviceCreds.Status.Components.GetComponent(openservicebroker.ComponentBind)
					Expect(component.Status).To(Equal(corev1.ErrorStatus))
					Expect(component.Message).To(Equal("failed to poll last bind operation on the service broker"))
					Expect(component.Detail).To(ContainSubstring("some error"))
				})

				It("should return the error", func() {
					expectError("some error")
				})
			})

			When("the response is bad request", func() {
				BeforeEach(func() {
					client.PollBindingLastOperationReturns(nil, osb.HTTPStatusCodeError{
						StatusCode:   http.StatusBadRequest,
						ErrorMessage: utils.StringPtr("some error"),
					})
				})

				It("should set the component status to failure", func() {
					component, _ := serviceCreds.Status.Components.GetComponent(openservicebroker.ComponentBind)
					Expect(component.Status).To(Equal(corev1.FailureStatus))
					Expect(component.Message).To(Equal("failed to poll last bind operation on the service broker"))
					Expect(component.Detail).To(ContainSubstring("some error"))
				})

				It("should return the error", func() {
					expectCriticalError("some error")
				})
			})

		})
	})

	Context("DeleteCredentials", func() {
		BeforeEach(func() {
			service.Status.ProviderID = KoreServiceID
			serviceCreds.Status.ProviderID = KoreServiceCredentialsID
		})

		JustBeforeEach(func() {
			reconcileResult, reconcileErr = provider.DeleteCredentials(serviceProviderCtx, service, serviceCreds)
		})

		When("the service credentials exist", func() {

			BeforeEach(func() {
				client.UnbindReturns(&osb.UnbindResponse{}, nil)
			})

			It("should call unbind", func() {
				Expect(client.UnbindCallCount()).To(Equal(1))

				req := client.UnbindArgsForCall(0)
				Expect(req.AcceptsIncomplete).To(BeTrue())
				Expect(req.BindingID).To(Equal(KoreServiceCredentialsID))
				Expect(req.InstanceID).To(Equal(KoreServiceID))
				Expect(req.ServiceID).To(Equal(Service1ID))
				Expect(req.PlanID).To(Equal(plan1ID))
			})

			When("the response has async=false", func() {
				BeforeEach(func() {
					client.UnbindReturns(&osb.UnbindResponse{
						Async: false,
					}, nil)
				})

				It("should set the component status to success", func() {
					component, _ := serviceCreds.Status.Components.GetComponent(openservicebroker.ComponentUnbind)
					Expect(component.Status).To(Equal(corev1.SuccessStatus))
				})

				It("should not requeue", func() {
					expectToNotRequeue()
				})
			})

			When("the response has async=true", func() {
				BeforeEach(func() {
					client.UnbindReturns(&osb.UnbindResponse{
						Async:        true,
						OperationKey: &Operation,
					}, nil)
				})

				It("should leave the component status as pending", func() {
					component, _ := serviceCreds.Status.Components.GetComponent(openservicebroker.ComponentUnbind)
					Expect(component.Status).To(Equal(corev1.PendingStatus))
				})

				It("should save the operation data", func() {
					Expect(serviceCreds.Status.ProviderData).To(Equal(createProviderData(&Operation)))
				})

				It("should requeue", func() {
					expectToRequeue()
				})
			})

			When("the response is an error", func() {
				BeforeEach(func() {
					client.UnbindReturns(nil, osb.HTTPStatusCodeError{
						StatusCode:   http.StatusInternalServerError,
						ErrorMessage: utils.StringPtr("some error"),
					})
				})

				It("should set the component status to error", func() {
					component, _ := serviceCreds.Status.Components.GetComponent(openservicebroker.ComponentUnbind)
					Expect(component.Status).To(Equal(corev1.ErrorStatus))
					Expect(component.Message).To(Equal("failed to call unbind on the service broker"))
					Expect(component.Detail).To(ContainSubstring("some error"))
				})

				It("should return the error", func() {
					expectError("some error")
				})
			})

			When("the response is bad request", func() {
				BeforeEach(func() {
					client.UnbindReturns(nil, osb.HTTPStatusCodeError{
						StatusCode:   http.StatusBadRequest,
						ErrorMessage: utils.StringPtr("some error"),
					})
				})

				It("should set the component status to failure", func() {
					component, _ := serviceCreds.Status.Components.GetComponent(openservicebroker.ComponentUnbind)
					Expect(component.Status).To(Equal(corev1.FailureStatus))
					Expect(component.Message).To(Equal("failed to call unbind on the service broker"))
					Expect(component.Detail).To(ContainSubstring("some error"))
				})

				It("should return the error", func() {
					expectCriticalError("some error")
				})
			})

		})

		When("the service credentials deletion is pending with async=true", func() {

			BeforeEach(func() {
				service.Status.ProviderID = KoreServiceID

				serviceCreds.Status.Components.SetCondition(corev1.Component{
					Name:   openservicebroker.ComponentUnbind,
					Status: corev1.PendingStatus,
				})
				serviceCreds.Status.ProviderID = KoreServiceCredentialsID
				serviceCreds.Status.ProviderData = createProviderData(&Operation)

				client.PollBindingLastOperationReturns(&osb.LastOperationResponse{
					State: osb.StateInProgress,
				}, nil)
			})

			It("should poll the last operation", func() {
				Expect(client.PollBindingLastOperationCallCount()).To(Equal(1))

				req := client.PollBindingLastOperationArgsForCall(0)
				Expect(req.InstanceID).To(Equal(KoreServiceID))
				Expect(*req.ServiceID).To(Equal(Service1ID))
				Expect(*req.PlanID).To(Equal(plan1ID))
				Expect(req.OperationKey).To(Equal(&Operation))
			})

			When("the status is succeeded", func() {
				BeforeEach(func() {
					client.PollBindingLastOperationReturns(&osb.LastOperationResponse{
						State: osb.StateSucceeded,
					}, nil)
				})

				It("should set the component status to success", func() {
					component, _ := serviceCreds.Status.Components.GetComponent(openservicebroker.ComponentUnbind)
					Expect(component.Status).To(Equal(corev1.SuccessStatus))
				})

				It("should not requeue", func() {
					expectToNotRequeue()
				})
			})

			When("the status is failed", func() {
				BeforeEach(func() {
					client.PollBindingLastOperationReturns(&osb.LastOperationResponse{
						State:       osb.StateFailed,
						Description: utils.StringPtr("some error"),
					}, nil)
				})

				It("should set the unbind status to failure", func() {
					component, _ := serviceCreds.Status.Components.GetComponent(openservicebroker.ComponentUnbind)
					Expect(component.Status).To(Equal(corev1.FailureStatus))
				})

				It("should return an error", func() {
					expectCriticalError("last bind operation failed on the service broker: some error")
				})
			})

			When("the status is in progress", func() {
				BeforeEach(func() {
					client.PollBindingLastOperationReturns(&osb.LastOperationResponse{
						State: osb.StateInProgress,
					}, nil)
				})

				It("should leave the status in pending", func() {
					component, _ := serviceCreds.Status.Components.GetComponent(openservicebroker.ComponentUnbind)
					Expect(component.Status).To(Equal(corev1.PendingStatus))
				})

				It("should requeue", func() {
					expectToRequeue()
				})
			})

			When("the response is an error", func() {
				BeforeEach(func() {
					client.PollBindingLastOperationReturns(nil, osb.HTTPStatusCodeError{
						StatusCode:   http.StatusInternalServerError,
						ErrorMessage: utils.StringPtr("some error"),
					})
				})

				It("should set the component status to error", func() {
					component, _ := serviceCreds.Status.Components.GetComponent(openservicebroker.ComponentUnbind)
					Expect(component.Status).To(Equal(corev1.ErrorStatus))
					Expect(component.Message).To(Equal("failed to poll last bind operation on the service broker"))
					Expect(component.Detail).To(ContainSubstring("some error"))
				})

				It("should return the error", func() {
					expectError("some error")
				})
			})

			When("the response is not found", func() {
				BeforeEach(func() {
					client.PollBindingLastOperationReturns(nil, osb.HTTPStatusCodeError{
						StatusCode: http.StatusNotFound,
					})
				})

				It("should set the component status to success", func() {
					component, _ := serviceCreds.Status.Components.GetComponent(openservicebroker.ComponentUnbind)
					Expect(component.Status).To(Equal(corev1.SuccessStatus))
				})

				It("should not requeue", func() {
					expectToNotRequeue()
				})
			})

			When("the response is bad request", func() {
				BeforeEach(func() {
					client.PollBindingLastOperationReturns(nil, osb.HTTPStatusCodeError{
						StatusCode:   http.StatusBadRequest,
						ErrorMessage: utils.StringPtr("some error"),
					})
				})

				It("should set the component status to failure", func() {
					component, _ := serviceCreds.Status.Components.GetComponent(openservicebroker.ComponentUnbind)
					Expect(component.Status).To(Equal(corev1.FailureStatus))
					Expect(component.Message).To(Equal("failed to poll last bind operation on the service broker"))
					Expect(component.Detail).To(ContainSubstring("some error"))
				})

				It("should return the error", func() {
					expectCriticalError("some error")
				})
			})

		})
	})

})
