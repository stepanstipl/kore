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

package cluster_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/controllers/controllerstest"
	clusterctrl "github.com/appvia/kore/pkg/controllers/management/cluster"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/kore/korefakes"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	"github.com/banzaicloud/k8s-objectmatcher/patch"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("Cluster Controller", func() {

	var test *controllerstest.Test
	var controller *clusterctrl.Controller
	var provider *korefakes.FakeClusterProvider
	var compProvider, compApp *servicesv1.Service
	var servicePlans *korefakes.FakeServicePlans
	var name = types.NamespacedName{Name: "testName", Namespace: "testNamespace"}
	var reconcileResult reconcile.Result
	var reconcileErr error
	var cluster *clustersv1.Cluster
	var clusterConfig map[string]interface{}

	var createCluster = func(status clustersv1.ClusterStatus) *clustersv1.Cluster {
		return &clustersv1.Cluster{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Clusters",
				APIVersion: clustersv1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:            name.Name,
				Namespace:       name.Namespace,
				ResourceVersion: "1",
			},
			Spec: clustersv1.ClusterSpec{
				Kind: "TEST",
				Plan: "testPlan",
				Credentials: corev1.Ownership{
					Group:     "credsGroup",
					Version:   "credsVersion",
					Kind:      "TestCredentials",
					Namespace: name.Namespace,
					Name:      "testCredentials",
				},
			},
			Status: status,
		}
	}

	var givenKubernetesExists = func() {
		cluster.Status.Components.SetCondition(corev1.Component{
			Name:   "Kubernetes/testName",
			Status: corev1.SuccessStatus,
			Resource: &corev1.Ownership{
				Group:     "clusters.compute.kore.appvia.io",
				Version:   "v1",
				Kind:      "Kubernetes",
				Namespace: "testNamespace",
				Name:      "testName",
			},
		})

		kub := &clustersv1.Kubernetes{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "testName",
				Namespace:   "testNamespace",
				Annotations: map[string]string{"kore.appvia.io/readonly": "true"},
			},
			Spec: clustersv1.KubernetesSpec{
				Cluster:  cluster.Ownership(),
				Provider: corev1.MustGetOwnershipFromObject(compProvider),
			},
			Status: clustersv1.KubernetesStatus{Status: corev1.SuccessStatus},
		}

		patchAnnotator := patch.NewAnnotator("kore.appvia.io/last-applied")
		_ = patchAnnotator.SetLastAppliedAnnotation(kub)

		test.Objects = append(test.Objects, kub)
	}

	var givenKubeAppManagerExists = func() {
		cluster.Status.Components.SetCondition(corev1.Component{
			Name:   "Service/testName-" + kore.AppAppManager,
			Status: corev1.SuccessStatus,
			Resource: &corev1.Ownership{
				Group:     "services.kore.appvia.io",
				Version:   "v1",
				Kind:      "Service",
				Namespace: "testNamespace",
				Name:      "testName-" + kore.AppAppManager,
			},
		})

		service := &servicesv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "testName-" + kore.AppAppManager,
				Namespace:   "testNamespace",
				Annotations: map[string]string{"kore.appvia.io/readonly": "true"},
			},
			Status: servicesv1.ServiceStatus{Status: corev1.SuccessStatus},
		}

		patchAnnotator := patch.NewAnnotator("kore.appvia.io/last-applied")
		_ = patchAnnotator.SetLastAppliedAnnotation(service)

		test.Objects = append(test.Objects, service)
	}

	var givenHelmOperatorExists = func() {
		cluster.Status.Components.SetCondition(corev1.Component{
			Name:   "Service/testName-" + kore.AppHelmOperator,
			Status: corev1.SuccessStatus,
			Resource: &corev1.Ownership{
				Group:     "services.kore.appvia.io",
				Version:   "v1",
				Kind:      "Service",
				Namespace: "testNamespace",
				Name:      "testName-" + kore.AppHelmOperator,
			},
		})

		service := &servicesv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "testName-" + kore.AppHelmOperator,
				Namespace:   "testNamespace",
				Annotations: map[string]string{"kore.appvia.io/readonly": "true"},
			},
			Status: servicesv1.ServiceStatus{Status: corev1.SuccessStatus},
		}

		patchAnnotator := patch.NewAnnotator("kore.appvia.io/last-applied")
		_ = patchAnnotator.SetLastAppliedAnnotation(service)

		test.Objects = append(test.Objects, service)
	}

	var givenCompAppExists = func() {
		cluster.Status.Components.SetCondition(corev1.Component{
			Name:   "Service/comp-app",
			Status: corev1.SuccessStatus,
			Resource: &corev1.Ownership{
				Group:     "services.kore.appvia.io",
				Version:   "v1",
				Kind:      "Service",
				Namespace: "testNamespace",
				Name:      "comp-app",
			},
		})

		compApp.Status.Status = corev1.SuccessStatus

		compApp.Annotations = map[string]string{
			"kore.appvia.io/readonly": "true",
		}
		patchAnnotator := patch.NewAnnotator("kore.appvia.io/last-applied")
		_ = patchAnnotator.SetLastAppliedAnnotation(compApp)

		test.Objects = append(test.Objects, compApp)
	}

	var givenCompProviderExists = func() {
		cluster.Status.Components.SetCondition(corev1.Component{
			Name:   "Service/comp-provider",
			Status: corev1.SuccessStatus,
			Resource: &corev1.Ownership{
				Group:     "services.kore.appvia.io",
				Version:   "v1",
				Kind:      "Service",
				Namespace: "testNamespace",
				Name:      "comp-provider",
			},
		})

		compProvider.Status.Status = corev1.SuccessStatus
		compProvider.Annotations = map[string]string{
			"kore.appvia.io/readonly": "true",
		}
		patchAnnotator := patch.NewAnnotator("kore.appvia.io/last-applied")
		_ = patchAnnotator.SetLastAppliedAnnotation(compProvider)

		test.Objects = append(test.Objects, compProvider)
	}

	BeforeEach(func() {
		test = controllerstest.NewTest(context.Background())
		controller = clusterctrl.NewController(test.Logger)
		provider = &korefakes.FakeClusterProvider{}
		provider.TypeReturns("TEST")

		kore.RegisterClusterProvider(provider)

		compProvider = &servicesv1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: servicesv1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "comp-provider",
				Namespace: "testNamespace",
			},
		}

		compApp = &servicesv1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: servicesv1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "comp-app",
				Namespace: "testNamespace",
			},
		}

		provider.SetComponentsStub = func(k kore.Context, v *clustersv1.Cluster, components *kore.ClusterComponents) error {
			kubernetesObj := components.Find(func(comp kore.ClusterComponent) bool {
				_, ok := comp.Object.(*clustersv1.Kubernetes)
				return ok
			})

			helmOperatorComp := components.Find(func(comp kore.ClusterComponent) bool {
				return strings.HasSuffix(comp.Object.GetName(), "-"+kore.AppHelmOperator)
			})

			components.AddComponent(&kore.ClusterComponent{
				Object:     compProvider,
				IsProvider: true,
			})
			components.Add(compApp, helmOperatorComp.Object)

			kubernetesObj.Dependencies = append(kubernetesObj.Dependencies, compProvider)

			return nil
		}

		servicePlans = &korefakes.FakeServicePlans{}
		servicePlans.GetReturns(&servicesv1.ServicePlan{}, nil)

		test.Kore.ServicePlansReturns(servicePlans)

		clusters := &korefakes.FakeClusters{}
		clusters.CheckDeleteReturns(nil)

		team := &korefakes.FakeTeam{}
		team.ClustersReturns(clusters)

		teams := &korefakes.FakeTeams{}
		teams.TeamReturns(team)

		test.Kore.TeamsReturns(teams)

		clusterConfig = map[string]interface{}{}
		cluster = createCluster(clustersv1.ClusterStatus{})
		test.Objects = []kubernetes.Object{cluster}
	})

	JustBeforeEach(func() {
		test.Initialize(controller)

		if clusterConfig != nil {
			configJson, _ := json.Marshal(clusterConfig)
			cluster.Spec.Configuration = v1beta1.JSON{Raw: configJson}
		}

		reconcileResult, reconcileErr = controller.Reconcile(reconcile.Request{NamespacedName: name})
	})

	AfterEach(func() {
		test.Stop()

		kore.UnregisterClusterProvider(provider)
	})

	Context("Reconcile", func() {
		When("getting the cluster returns an error", func() {
			BeforeEach(func() {
				test.Client.GetReturnsOnCall(0, errors.New("some random error"))
			})
			It("should return an error", func() {
				Expect(reconcileErr).To(HaveOccurred())
			})
		})

		When("the cluster object does not exist", func() {
			BeforeEach(func() {
				test.Objects = []kubernetes.Object{}
			})
			It("should not requeue", func() {
				Expect(reconcileErr).ToNot(HaveOccurred())
				Expect(reconcileResult).To(Equal(reconcile.Result{}))
			})
		})

		When("the cluster tagged as a system cluster (Kore cluster)", func() {
			BeforeEach(func() {
				cluster.Annotations = map[string]string{
					kore.AnnotationSystem: kore.AnnotationValueTrue,
				}
			})

			It("should set the status to successful", func() {
				updatedCluster := &clustersv1.Cluster{}
				test.ExpectStatusPatch(0, updatedCluster)
				Expect(updatedCluster.Status.Status).To(Equal(corev1.SuccessStatus), "")
			})

			It("should not requeue", func() {
				Expect(reconcileErr).ToNot(HaveOccurred())
				Expect(reconcileResult).To(Equal(reconcile.Result{}))
			})
		})

		When("the cluster is reconciled the first time", func() {
			It("should set the finaliser", func() {
				updatedCluster := &clustersv1.Cluster{}
				test.ExpectUpdate(0, updatedCluster)
				Expect(updatedCluster.Finalizers).To(Equal([]string{"cluster.clusters.kore.appvia.io"}))
				test.ExpectRequeue(reconcileResult, reconcileErr)
			})

			It("should requeue without an error", func() {
				test.ExpectRequeue(reconcileResult, reconcileErr)
			})

			When("setting the finaliser fails", func() {
				BeforeEach(func() {
					test.Client.UpdateReturnsOnCall(0, fmt.Errorf("some error"))
				})

				It("should retry", func() {
					Expect(test.Client.UpdateCallCount()).To(Equal(1))
					Expect(reconcileErr).To(HaveOccurred())
				})
			})
		})

		When("the cluster has a finaliser", func() {
			var updatedCluster *clustersv1.Cluster

			BeforeEach(func() {
				cluster.Finalizers = []string{"cluster.clusters.kore.appvia.io"}
			})

			JustBeforeEach(func() {
				updatedCluster = &clustersv1.Cluster{}
				test.ExpectStatusPatch(0, updatedCluster)

				Expect(test.Client.DeleteCallCount()).To(Equal(0))
			})

			It("should requeue without an error", func() {
				test.ExpectRequeue(reconcileResult, reconcileErr)
			})

			It("should set the status to pending", func() {
				Expect(updatedCluster.Status.Status).To(Equal(corev1.PendingStatus), updatedCluster.Status.Message)
			})
		})

		When("the cluster has a finaliser and is in a pending state", func() {
			var updatedCluster *clustersv1.Cluster

			BeforeEach(func() {
				cluster.Finalizers = []string{"cluster.clusters.kore.appvia.io"}
				cluster.Status.Status = corev1.PendingStatus
			})

			JustBeforeEach(func() {
				updatedCluster = &clustersv1.Cluster{}
				test.ExpectStatusPatch(0, updatedCluster)

				Expect(test.Client.DeleteCallCount()).To(Equal(0))
			})

			It("should requeue without an error", func() {
				test.ExpectRequeue(reconcileResult, reconcileErr)
			})

			It("should create the first component in dependency order (comp1)", func() {
				obj := &servicesv1.Service{}
				test.ExpectCreate(0, obj)
				Expect(obj.Name).To(Equal("comp-provider"))
				Expect(obj.Annotations[kore.Label("last-applied")]).ToNot(Equal(""))

				Expect(test.Client.CreateCallCount()).To(Equal(1))
			})

			It("should add the component to the status", func() {
				statusComp, ok := updatedCluster.Status.Components.GetComponent("Service/comp-provider")
				Expect(ok).To(BeTrue(), "Service/comp-provider component is not present in status components")
				Expect(*statusComp).To(Equal(corev1.Component{
					Name:   "Service/comp-provider",
					Status: corev1.PendingStatus,
					Resource: &corev1.Ownership{
						Group:     "services.kore.appvia.io",
						Version:   "v1",
						Kind:      "Service",
						Namespace: "testNamespace",
						Name:      "comp-provider",
					},
				}))
			})
		})

		When("all Kubernetes dependencies are ready", func() {
			var updatedCluster *clustersv1.Cluster

			BeforeEach(func() {
				cluster.Finalizers = []string{"cluster.clusters.kore.appvia.io"}
				cluster.Status.Status = corev1.PendingStatus

				givenCompProviderExists()
			})

			JustBeforeEach(func() {
				updatedCluster = &clustersv1.Cluster{}
				test.ExpectStatusPatch(0, updatedCluster)

				Expect(test.Client.DeleteCallCount()).To(Equal(0))
			})

			It("should requeue without an error", func() {
				test.ExpectRequeue(reconcileResult, reconcileErr)
			})

			It("should create the Kubernetes component", func() {
				obj := &clustersv1.Kubernetes{}
				test.ExpectCreate(0, obj)

				Expect(test.Client.CreateCallCount()).To(Equal(1))
			})

			It("should add the component to the status", func() {
				statusComp, ok := updatedCluster.Status.Components.GetComponent("Kubernetes/testName")
				Expect(ok).To(BeTrue(), "Kubernetes/testName component is not present in status components")
				Expect(*statusComp).To(Equal(corev1.Component{
					Name:   "Kubernetes/testName",
					Status: corev1.PendingStatus,
					Resource: &corev1.Ownership{
						Group:     "clusters.compute.kore.appvia.io",
						Version:   "v1",
						Kind:      "Kubernetes",
						Namespace: "testNamespace",
						Name:      "testName",
					},
				}))
			})
		})

		When("Common components are ready", func() {
			var updatedCluster *clustersv1.Cluster

			BeforeEach(func() {
				cluster.Finalizers = []string{"cluster.clusters.kore.appvia.io"}
				cluster.Status.Status = corev1.PendingStatus

				givenCompProviderExists()
				givenKubernetesExists()
				givenKubeAppManagerExists()
				givenHelmOperatorExists()
			})

			JustBeforeEach(func() {
				updatedCluster = &clustersv1.Cluster{}
				test.ExpectStatusPatch(0, updatedCluster)

				Expect(test.Client.DeleteCallCount()).To(Equal(0))
			})

			It("should requeue without an error", func() {
				test.ExpectRequeue(reconcileResult, reconcileErr)
			})

			It("should create the app service component", func() {
				obj := &servicesv1.Service{}
				test.ExpectCreate(0, obj)
				Expect(obj.Name).To(Equal("comp-app"))
				Expect(test.Client.CreateCallCount()).To(Equal(1))
			})

			It("should add the component to the status", func() {
				statusComp, ok := updatedCluster.Status.Components.GetComponent("Service/comp-app")
				Expect(ok).To(BeTrue(), "Kubernetes/testName component is not present in status components")
				Expect(*statusComp).To(Equal(corev1.Component{
					Name:   "Service/comp-app",
					Status: corev1.PendingStatus,
					Resource: &corev1.Ownership{
						Group:     "services.kore.appvia.io",
						Version:   "v1",
						Kind:      "Service",
						Namespace: "testNamespace",
						Name:      "comp-app",
					},
				}))
			})

			When("the last component is ready", func() {
				BeforeEach(func() {
					givenCompAppExists()

					provider.SetProviderDataStub = func(ctx kore.Context, cluster *clustersv1.Cluster, components *kore.ClusterComponents) error {
						cluster.Status.SetProviderData("DONE")
						return nil
					}
				})

				It("should update the cluster status to successful", func() {
					Expect(updatedCluster.Status.Status).To(Equal(corev1.SuccessStatus))
				})

				It("should not requeue", func() {
					Expect(reconcileErr).ToNot(HaveOccurred())
					Expect(reconcileResult).To(Equal(reconcile.Result{}))
				})

				It("should set the provider data", func() {
					Expect(updatedCluster.Status.ProviderData).ToNot(BeNil())
					Expect(string(updatedCluster.Status.ProviderData.Raw)).To(Equal(`"DONE"`))
				})
			})
		})

		When("The cluster has been deleted", func() {
			var updatedCluster *clustersv1.Cluster

			BeforeEach(func() {
				cluster.Finalizers = []string{"cluster.clusters.kore.appvia.io"}
				cluster.DeletionTimestamp = &metav1.Time{Time: time.Now()}
			})

			JustBeforeEach(func() {
				updatedCluster = &clustersv1.Cluster{}
				test.ExpectStatusPatch(0, updatedCluster)

				Expect(test.Client.DeleteCallCount()).To(Equal(0))
			})

			It("should requeue without an error", func() {
				test.ExpectRequeue(reconcileResult, reconcileErr)
			})

			It("should set the status to deleting", func() {
				Expect(updatedCluster.Status.Status).To(Equal(corev1.DeletingStatus), updatedCluster.Status.Message)
			})
		})

		When("The cluster has been deleted and has a deleting status", func() {
			var updatedCluster *clustersv1.Cluster

			BeforeEach(func() {
				cluster.Finalizers = []string{"cluster.clusters.kore.appvia.io"}
				cluster.DeletionTimestamp = &metav1.Time{Time: time.Now()}
				cluster.Status.Status = corev1.DeletingStatus

				givenCompProviderExists()
				givenKubernetesExists()
				givenKubeAppManagerExists()
				givenHelmOperatorExists()
				givenCompAppExists()
			})

			JustBeforeEach(func() {
				updatedCluster = &clustersv1.Cluster{}
				test.ExpectStatusPatch(0, updatedCluster)
			})

			It("should requeue without an error", func() {
				test.ExpectRequeue(reconcileResult, reconcileErr)
			})

			It("should remove the last component in dependency order (comp-app)", func() {
				obj := &servicesv1.Service{}
				test.ExpectDelete(0, obj)
				Expect(obj.Name).To(Equal("comp-app"))
				Expect(test.Client.DeleteCallCount()).To(Equal(1))
			})

			It("should set the component status to deleting", func() {
				statusComp, ok := updatedCluster.Status.Components.GetComponent("Service/comp-app")
				Expect(ok).To(BeTrue(), "Service/comp-app component is not present in status components")
				Expect(statusComp.Status).To(Equal(corev1.DeletingStatus))
			})
		})

		When("All services has been deleted", func() {
			var updatedCluster *clustersv1.Cluster

			BeforeEach(func() {
				cluster.Finalizers = []string{"cluster.clusters.kore.appvia.io"}
				cluster.DeletionTimestamp = &metav1.Time{Time: time.Now()}
				cluster.Status.Status = corev1.DeletingStatus

				givenCompProviderExists()
				givenKubernetesExists()
			})

			JustBeforeEach(func() {
				updatedCluster = &clustersv1.Cluster{}
				test.ExpectStatusPatch(0, updatedCluster)
			})

			It("should requeue without an error", func() {
				test.ExpectRequeue(reconcileResult, reconcileErr)
			})

			It("should delete the Kubernetes component", func() {
				obj := &clustersv1.Kubernetes{}
				test.ExpectDelete(0, obj)
				Expect(test.Client.DeleteCallCount()).To(Equal(1))
			})
		})

		When("All components have been deleted", func() {
			var updatedCluster *clustersv1.Cluster

			BeforeEach(func() {
				cluster.Finalizers = []string{"cluster.clusters.kore.appvia.io"}
				cluster.DeletionTimestamp = &metav1.Time{Time: time.Now()}
				cluster.Status.Status = corev1.DeletingStatus
			})

			JustBeforeEach(func() {
				updatedCluster = &clustersv1.Cluster{}
				test.ExpectStatusPatch(0, updatedCluster)
			})

			It("should not requeue", func() {
				Expect(reconcileErr).ToNot(HaveOccurred())
				Expect(reconcileResult).To(Equal(reconcile.Result{}))
			})

			It("should remove the finalizer", func() {
				Expect(updatedCluster.Finalizers).To(BeEmpty())
			})
		})
	})
})
