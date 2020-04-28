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
	"time"

	"github.com/appvia/kore/pkg/controllers/controllerstest"

	"k8s.io/apimachinery/pkg/runtime"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eksv1alpha1 "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/appvia/kore/pkg/controllers/management/cluster"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("Cluster Controller", func() {

	var test *controllerstest.Test
	var controller *cluster.Controller

	BeforeEach(func() {
		test = controllerstest.NewTest(context.Background())
		controller = cluster.NewController(test.Logger)
	})

	JustBeforeEach(func() {
		test.Run(controller)
	})

	AfterEach(func() {
		test.Stop()
	})

	PContext("Reconcile", func() {
		var reconcileResult reconcile.Result
		var reconcileErr error
		var kind string
		var cluster *clustersv1.Cluster
		var name = types.NamespacedName{Name: "testName", Namespace: "testNamespace"}
		var clusterConfig map[string]interface{}

		var createCluster = func(kind string, status clustersv1.ClusterStatus) *clustersv1.Cluster {
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
					Kind: kind,
					Plan: "testPlan",
					Credentials: corev1.Ownership{
						Group:     "credsGroup",
						Version:   "credsVersion",
						Kind:      kind + "Credentials",
						Namespace: name.Namespace,
						Name:      "testCredentials",
					},
				},
				Status: status,
			}
		}

		BeforeEach(func() {
			clusterConfig = map[string]interface{}{}
			cluster = createCluster(kind, clustersv1.ClusterStatus{})
		})

		JustBeforeEach(func() {
			if clusterConfig != nil {
				configJson, _ := json.Marshal(clusterConfig)
				cluster.Spec.Configuration = v1beta1.JSON{Raw: configJson}
			}
			reconcileResult, reconcileErr = controller.Reconcile(reconcile.Request{NamespacedName: name})
		})

		When("getting the cluster returns an error", func() {
			BeforeEach(func() {
				test.Client.GetReturnsOnCall(0, errors.New("some random error"))
			})
			It("should requeue", func() {
				Expect(reconcileErr).To(HaveOccurred())
			})
		})

		When("the cluster object does not exist", func() {
			It("should not requeue", func() {
				Expect(reconcileErr).ToNot(HaveOccurred())
				Expect(reconcileResult).To(Equal(reconcile.Result{}))
			})
		})

		Context("EKS cluster", func() {
			var creds *eksv1alpha1.EKSCredentials

			BeforeEach(func() {
				kind = "EKS"
				clusterConfig = map[string]interface{}{
					"nodeGroups": []map[string]string{
						{"name": "ng1"},
						{"name": "ng2"},
					},
				}
				creds = &eksv1alpha1.EKSCredentials{
					TypeMeta: metav1.TypeMeta{
						Kind:       "EKSCredentials",
						APIVersion: eksv1alpha1.GroupVersion.String(),
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "testCredentials",
						Namespace: name.Namespace,
					},
					Spec:   eksv1alpha1.EKSCredentialsSpec{},
					Status: eksv1alpha1.EKSCredentialsStatus{},
				}
				test.Objects = []runtime.Object{cluster, creds}
			})

			When("the cluster is reconciled the first time", func() {
				It("should set the finaliser", func() {
					updatedCluster := &clustersv1.Cluster{}
					test.ExpectUpdate(0, updatedCluster)
					Expect(updatedCluster.Finalizers).To(Equal([]string{"cluster.clusters.kore.appvia.io"}))
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

				It("should set the status to pending", func() {
					Expect(updatedCluster.Status.Status).To(Equal(corev1.PendingStatus), updatedCluster.Status.Message)
				})

				It("should set the status to pending", func() {
					Expect(updatedCluster.Status.Status).To(Equal(corev1.PendingStatus), updatedCluster.Status.Message)
				})

				It("should create the component resources", func() {
					eksVPC := &eksv1alpha1.EKSVPC{}
					test.ExpectCreate(0, eksVPC)
					Expect(eksVPC.Spec.Cluster).To(Equal(cluster.Ownership()))

					Expect(test.Client.CreateCallCount()).To(Equal(1))
				})

				It("should requeue and not error", func() {
					test.ExpectRequeue(reconcileResult, reconcileErr)
				})

				When("the components can not be created", func() {
					BeforeEach(func() {
						clusterConfig = nil
					})

					It("should fail and not requeue", func() {
						updatedCluster := &clustersv1.Cluster{}
						test.ExpectStatusPatch(0, updatedCluster)
						Expect(updatedCluster.Status.Status).To(Equal(corev1.FailureStatus))

						Expect(reconcileErr).ToNot(HaveOccurred())
					})
				})

				When("creating a component fails", func() {
					BeforeEach(func() {
						test.Client.CreateReturnsOnCall(0, errors.New("some random error"))
					})

					It("should requeue", func() {
						Expect(test.Client.CreateCallCount()).To(Equal(1))
						_, obj, _ := test.Client.CreateArgsForCall(0)
						Expect(obj).To(BeAssignableToTypeOf(&eksv1alpha1.EKSVPC{}))

						Expect(reconcileErr).To(HaveOccurred())
					})
				})

				When("getting a component fails", func() {
					BeforeEach(func() {
						eksVPC := &eksv1alpha1.EKSVPC{ObjectMeta: metav1.ObjectMeta{Name: name.Name, Namespace: name.Namespace}}
						eksVPC.Labels = map[string]string{controllerstest.LabelGetError: "some error"}
						test.Objects = append(test.Objects, eksVPC)
					})

					It("should requeue", func() {
						Expect(reconcileErr).To(HaveOccurred())
					})
				})

				When("updating an existing component fails", func() {
					BeforeEach(func() {
						eksvpc := &eksv1alpha1.EKSVPC{ObjectMeta: metav1.ObjectMeta{Name: name.Name, Namespace: name.Namespace}}
						eksvpc.Status.Status = corev1.PendingStatus
						test.Client.PatchReturnsOnCall(0, errors.New("some random error"))
						test.Objects = append(test.Objects, eksvpc)
					})

					It("should requeue", func() {
						updatedEKSVPC := &eksv1alpha1.EKSVPC{}
						test.ExpectPatch(0, updatedEKSVPC)
						Expect(test.Client.PatchCallCount()).To(Equal(1))

						Expect(reconcileErr).To(HaveOccurred())
					})
				})

				When("applying the cluster configuration on an existing component fails", func() {
					BeforeEach(func() {
						clusterConfig = nil
						test.Objects = append(test.Objects, &eksv1alpha1.EKSVPC{ObjectMeta: metav1.ObjectMeta{Name: name.Name, Namespace: name.Namespace}})
					})

					It("should fail and not requeue", func() {
						updatedCluster := &clustersv1.Cluster{}
						test.ExpectStatusPatch(0, updatedCluster)
						Expect(updatedCluster.Status.Status).To(Equal(corev1.FailureStatus))

						Expect(reconcileErr).ToNot(HaveOccurred())
					})
				})

				When("the EKSVPC component is complete", func() {
					BeforeEach(func() {
						eksVPC := &eksv1alpha1.EKSVPC{ObjectMeta: metav1.ObjectMeta{Name: name.Name, Namespace: name.Namespace}}
						eksVPC.Labels = map[string]string{"cluster.clusters.kore.appvia.io/ResourceVersion": cluster.ResourceVersion}
						eksVPC.Status.Status = corev1.SuccessStatus
						test.Objects = append(test.Objects, eksVPC)
					})

					It("should update the component status on the cluster", func() {
						updatedCluster := &clustersv1.Cluster{}
						test.ExpectStatusPatch(0, updatedCluster)

						component, _ := updatedCluster.Status.Components.GetComponent("EKSVPC/testName")
						Expect(component.Status).To(Equal(corev1.SuccessStatus))
					})

					It("should create the EKS resources", func() {
						eks := &eksv1alpha1.EKS{}
						test.ExpectCreate(0, eks)
						Expect(eks.Spec.Cluster).To(Equal(cluster.Ownership()))

						Expect(test.Client.CreateCallCount()).To(Equal(1))
					})

					When("the EKS component is complete", func() {
						BeforeEach(func() {
							eks := &eksv1alpha1.EKS{ObjectMeta: metav1.ObjectMeta{Name: name.Name, Namespace: name.Namespace}}
							eks.Labels = map[string]string{"cluster.clusters.kore.appvia.io/ResourceVersion": cluster.ResourceVersion}
							eks.Status.Status = corev1.SuccessStatus
							test.Objects = append(test.Objects, eks)
						})

						It("should update the component status on the cluster", func() {
							updatedCluster := &clustersv1.Cluster{}
							test.ExpectStatusPatch(0, updatedCluster)

							component, _ := updatedCluster.Status.Components.GetComponent("EKS/testName")
							Expect(component.Status).To(Equal(corev1.SuccessStatus))
						})

						It("should create the EKS nodegroups", func() {
							eksNodeGroup1 := &eksv1alpha1.EKSNodeGroup{}
							test.ExpectCreate(0, eksNodeGroup1)
							Expect(eksNodeGroup1.Spec.Cluster).To(Equal(cluster.Ownership()))

							eksNodeGroup2 := &eksv1alpha1.EKSNodeGroup{}
							test.ExpectCreate(1, eksNodeGroup2)
							Expect(eksNodeGroup2.Spec.Cluster).To(Equal(cluster.Ownership()))

							Expect(test.Client.CreateCallCount()).To(Equal(2))
						})

						When("the EKS Node groups are complete", func() {
							BeforeEach(func() {
								eksNodeGroup1 := &eksv1alpha1.EKSNodeGroup{ObjectMeta: metav1.ObjectMeta{Name: name.Name + "-ng1", Namespace: name.Namespace}}
								eksNodeGroup1.Labels = map[string]string{"cluster.clusters.kore.appvia.io/ResourceVersion": cluster.ResourceVersion}
								eksNodeGroup1.Status.Status = corev1.SuccessStatus

								eksNodeGroup2 := &eksv1alpha1.EKSNodeGroup{ObjectMeta: metav1.ObjectMeta{Name: name.Name + "-ng2", Namespace: name.Namespace}}
								eksNodeGroup2.Labels = map[string]string{"cluster.clusters.kore.appvia.io/ResourceVersion": cluster.ResourceVersion}
								eksNodeGroup2.Status.Status = corev1.SuccessStatus
								test.Objects = append(test.Objects, eksNodeGroup1, eksNodeGroup2)
							})

							It("should update the component statuses on the cluster", func() {
								updatedCluster := &clustersv1.Cluster{}
								test.ExpectStatusPatch(0, updatedCluster)

								component1, _ := updatedCluster.Status.Components.GetComponent("EKSNodeGroup/testName-ng1")
								Expect(component1.Status).To(Equal(corev1.SuccessStatus))
								component2, _ := updatedCluster.Status.Components.GetComponent("EKSNodeGroup/testName-ng2")
								Expect(component2.Status).To(Equal(corev1.SuccessStatus))
							})

							It("should create the Kubernetes resource", func() {
								kubernetes := &clustersv1.Kubernetes{}
								test.ExpectCreate(0, kubernetes)
								Expect(kubernetes.Spec.Cluster).To(Equal(cluster.Ownership()))

								Expect(test.Client.CreateCallCount()).To(Equal(1))
							})

							When("Kubernetes is complete", func() {
								var kubernetes *clustersv1.Kubernetes

								BeforeEach(func() {
									kubernetes = clustersv1.NewKubernetes(name.Name, name.Namespace)
									kubernetes.Status.Status = corev1.SuccessStatus
									kubernetes.Labels = map[string]string{"cluster.clusters.kore.appvia.io/ResourceVersion": cluster.ResourceVersion}
									test.Objects = append(test.Objects, kubernetes)
								})

								It("should update the status on the cluster", func() {
									updatedCluster := &clustersv1.Cluster{}
									test.ExpectStatusPatch(0, updatedCluster)

									Expect(updatedCluster.Status.Status).To(Equal(corev1.SuccessStatus))
								})

								It("should set the CA certificate on the cluster", func() {
									Expect(updatedCluster.Status.CaCertificate).To(Equal(kubernetes.Status.CaCertificate))
								})

								It("should set the API endpoint on the cluster", func() {
									Expect(updatedCluster.Status.APIEndpoint).To(Equal(kubernetes.Status.APIEndpoint))
								})

								It("should set the Auth proxy endpoint on the cluster", func() {
									Expect(updatedCluster.Status.AuthProxyEndpoint).To(Equal(kubernetes.Status.Endpoint))
								})

								It("should finish the reconciliation", func() {
									Expect(reconcileErr).ToNot(HaveOccurred())
									Expect(reconcileResult).To(Equal(reconcile.Result{}))
								})
							})
						})
					})
				})

				When("a component failed", func() {
					BeforeEach(func() {
						eks := &eksv1alpha1.EKS{ObjectMeta: metav1.ObjectMeta{Name: name.Name, Namespace: name.Namespace}}
						eks.Status.Status = corev1.FailureStatus
						eks.Status.Conditions = corev1.Components{
							&corev1.Component{
								Name:    "some component",
								Status:  corev1.FailureStatus,
								Message: "some error",
								Detail:  "some error detail",
							},
						}
						eks.Labels = map[string]string{"cluster.clusters.kore.appvia.io/ResourceVersion": cluster.ResourceVersion}
						test.Objects = append(test.Objects, eks)
					})

					It("should update the status", func() {
						updatedCluster := &clustersv1.Cluster{}
						test.ExpectStatusPatch(0, updatedCluster)
						Expect(updatedCluster.Status.Status).To(Equal(corev1.FailureStatus))
						component, _ := updatedCluster.Status.Components.GetComponent("EKS/testName")
						Expect(component.Status).To(Equal(corev1.FailureStatus))
						Expect(component.Message).To(Equal("[Failure] some component - some error: some error detail"))
					})

					It("should fail and not requeue", func() {
						Expect(reconcileErr).ToNot(HaveOccurred())
					})
				})

				When("updating the status fails", func() {
					BeforeEach(func() {
						test.StatusClient.PatchReturnsOnCall(0, fmt.Errorf("some error"))
					})
					It("should retry", func() {
						Expect(test.StatusClient.PatchCallCount()).To(Equal(1))
						Expect(reconcileErr).To(HaveOccurred())
					})
				})
			})

			When("a cluster was successfully provisioned", func() {
				var kubernetes *clustersv1.Kubernetes
				var eks *eksv1alpha1.EKS
				var ng1, ng2 *eksv1alpha1.EKSNodeGroup
				var eksvpc *eksv1alpha1.EKSVPC

				BeforeEach(func() {
					cluster.Finalizers = []string{"cluster.clusters.kore.appvia.io"}
					cluster.Status.Components = corev1.Components{
						{Name: "EKS/testName", Status: corev1.SuccessStatus},
						{Name: "EKSVPC/testName", Status: corev1.SuccessStatus},
						{Name: "EKSNodeGroup/testName-ng1", Status: corev1.SuccessStatus},
						{Name: "EKSNodeGroup/testName-ng2", Status: corev1.SuccessStatus},
						{Name: "Kubernetes/testName", Status: corev1.SuccessStatus},
					}

					kubernetes = clustersv1.NewKubernetes(name.Name, name.Namespace)
					kubernetes.Status.Status = corev1.SuccessStatus
					eks = &eksv1alpha1.EKS{ObjectMeta: metav1.ObjectMeta{Name: name.Name, Namespace: name.Namespace}}
					eks.Status.Status = corev1.SuccessStatus
					eksvpc := &eksv1alpha1.EKSVPC{ObjectMeta: metav1.ObjectMeta{Name: name.Name, Namespace: name.Namespace}}
					eksvpc.Status.Status = corev1.SuccessStatus
					ng1 = &eksv1alpha1.EKSNodeGroup{ObjectMeta: metav1.ObjectMeta{Name: name.Name + "-ng1", Namespace: name.Namespace}}
					ng1.Status.Status = corev1.SuccessStatus
					ng2 = &eksv1alpha1.EKSNodeGroup{ObjectMeta: metav1.ObjectMeta{Name: name.Name + "-ng2", Namespace: name.Namespace}}
					ng2.Status.Status = corev1.SuccessStatus
					test.Objects = append(test.Objects, kubernetes, eks, eksvpc, ng1, ng2)
				})

				When("an EKS Node group was removed from the config", func() {
					BeforeEach(func() {
						cluster.Finalizers = []string{"cluster.clusters.kore.appvia.io"}
						clusterConfig = map[string]interface{}{
							"nodeGroups": []map[string]string{
								{"name": "ng1"},
							},
						}
					})

					It("should delete the removed node group", func() {
						deletedEKSNodeGroup := &eksv1alpha1.EKSNodeGroup{}
						test.ExpectDelete(0, deletedEKSNodeGroup)
						Expect(test.Client.DeleteCallCount()).To(Equal(1))

						Expect(deletedEKSNodeGroup.Name).To(Equal("testName-ng2"))
					})

					It("should requeue", func() {
						test.ExpectRequeue(reconcileResult, reconcileErr)
					})

					When("the EKS node group was deleted", func() {
						BeforeEach(func() {
							test.Objects = []runtime.Object{cluster, kubernetes, eks, eksvpc, ng1}
						})

						It("should remove the component", func() {
							updatedCluster := &clustersv1.Cluster{}
							test.ExpectStatusPatch(0, updatedCluster)
							_, exists := updatedCluster.Status.Components.GetComponent("EKSNodeGroup/testName-ng2")
							Expect(exists).To(BeFalse())
						})

						It("should have a success status", func() {
							updatedCluster := &clustersv1.Cluster{}
							test.ExpectStatusPatch(0, updatedCluster)
							Expect(updatedCluster.Status.Status).To(Equal(corev1.SuccessStatus))
						})
					})
				})

				When("the cluster has a deletion timestamp", func() {
					BeforeEach(func() {
						cluster.Finalizers = []string{"cluster.clusters.kore.appvia.io"}
						cluster.DeletionTimestamp = &metav1.Time{Time: time.Now()}
					})

					It("should delete the Kubernetes component first", func() {
						Expect(test.Client.DeleteCallCount()).To(Equal(1))
						deletedObj := &clustersv1.Kubernetes{}
						test.ExpectDelete(0, deletedObj)
					})

					It("should requeue and not error", func() {
						test.ExpectRequeue(reconcileResult, reconcileErr)
					})

					When("deleting the Kubernetes component fails", func() {
						BeforeEach(func() {
							test.Client.DeleteReturnsOnCall(0, errors.New("some error"))
						})

						It("should requeue", func() {
							Expect(reconcileErr).To(HaveOccurred())
						})
					})

					When("getting a component fails", func() {
						BeforeEach(func() {
							eks := &eksv1alpha1.EKS{ObjectMeta: metav1.ObjectMeta{Name: name.Name, Namespace: name.Namespace}}
							eks.Labels = map[string]string{controllerstest.LabelGetError: "some error"}
							test.Objects = []runtime.Object{cluster, eks}
						})

						It("should requeue", func() {
							Expect(reconcileErr).To(HaveOccurred())
						})
					})

					When("a component can not be deleted", func() {
						BeforeEach(func() {
							kubernetes := clustersv1.NewKubernetes(name.Name, name.Namespace)
							kubernetes.Status.Status = corev1.DeleteFailedStatus
							kubernetes.Status.Components = corev1.Components{
								{Name: "testComponent", Status: corev1.DeleteFailedStatus, Message: "testMessage", Detail: "testDetail"},
							}
							kubernetes.DeletionTimestamp = &metav1.Time{Time: time.Now()}
							test.Objects = []runtime.Object{cluster, kubernetes}
						})

						It("should set the status to delete failed", func() {
							updatedCluster := &clustersv1.Cluster{}
							test.ExpectStatusPatch(0, updatedCluster)
							Expect(updatedCluster.Status.Status).To(Equal(corev1.DeleteFailedStatus))
						})

						It("should not requeue", func() {
							Expect(reconcileErr).ToNot(HaveOccurred())
						})
					})

					When("all components were deleted", func() {
						BeforeEach(func() {
							test.Objects = []runtime.Object{cluster}
						})

						It("should set the status to deleted", func() {
							updatedCluster := &clustersv1.Cluster{}
							test.ExpectStatusPatch(0, updatedCluster)
							Expect(updatedCluster.Status.Status).To(Equal(corev1.DeletedStatus))
						})

						It("should requeue", func() {
							test.ExpectRequeue(reconcileResult, reconcileErr)
						})
					})

					When("the cluster has deleted status", func() {
						BeforeEach(func() {
							cluster.Status.Status = corev1.DeletedStatus
							test.Objects = []runtime.Object{cluster}
						})

						It("should remove the finalizer", func() {
							updatedCluster := &clustersv1.Cluster{}
							test.ExpectUpdate(0, updatedCluster)
							Expect(updatedCluster.Finalizers).To(BeEmpty())
						})

						It("should not requeue", func() {
							Expect(reconcileErr).ToNot(HaveOccurred())
							Expect(reconcileResult).To(Equal(reconcile.Result{}))
						})
					})
				})
			})
		})
	})
})
