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
	"errors"
	"fmt"

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

	Context("Reconcile", func() {
		var reconcileResult reconcile.Result
		var reconcileErr error
		var cluster *clustersv1.Cluster
		var name = types.NamespacedName{Name: "testName", Namespace: "testNamespace"}

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
					Kind:          kind,
					Plan:          "testPlan",
					Configuration: v1beta1.JSON{Raw: []byte(`{"nodeGroups":[{"name":"ng1"}, {"name":"ng2"}]}`)},
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

		JustBeforeEach(func() {
			reconcileResult, reconcileErr = controller.Reconcile(reconcile.Request{NamespacedName: name})
		})

		When("getting the cluster returns an error", func() {
			BeforeEach(func() {
				test.Client.GetReturnsOnCall(0, errors.New("some random error"))
			})
			It("should requeue", func() {
				Expect(reconcileErr).ToNot(HaveOccurred())
				Expect(reconcileResult).To(Equal(reconcile.Result{Requeue: true}))
			})
		})

		When("the cluster object does not exist", func() {
			It("should not requeue", func() {
				Expect(reconcileErr).To(HaveOccurred())
				Expect(reconcileResult).To(Equal(reconcile.Result{}))
			})
		})

		Context("EKS cluster", func() {
			var creds *eksv1alpha1.EKSCredentials

			BeforeEach(func() {
				cluster = createCluster("EKS", clustersv1.ClusterStatus{})
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
					patchedCluster := &clustersv1.Cluster{}
					test.ExpectPatch(0, patchedCluster)
					Expect(patchedCluster.Finalizers).To(Equal([]string{"cluster.clusters.kore.appvia.io"}))
					test.ExpectRequeue(reconcileResult, reconcileErr)
				})

				When("setting the finaliser fails", func() {
					BeforeEach(func() {
						test.Client.PatchReturnsOnCall(0, fmt.Errorf("some error"))
					})

					It("should retry", func() {
						Expect(test.Client.PatchCallCount()).To(Equal(1))
						test.ExpectRequeue(reconcileResult, reconcileErr)
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
					test.ExpectStatusUpdate(0, updatedCluster)
				})

				It("should set the status to pending", func() {
					Expect(updatedCluster.Status.Status).To(Equal(corev1.PendingStatus), updatedCluster.Status.Message)
				})

				It("should set the components with pending status", func() {
					Expect(updatedCluster.Status.Components).To(Equal(corev1.Components{
						{Name: "EKS/testName", Status: corev1.PendingStatus},
						{Name: "EKSNodeGroup/testName-ng1", Status: corev1.PendingStatus},
						{Name: "EKSNodeGroup/testName-ng2", Status: corev1.PendingStatus},
						{Name: "Kubernetes/testName", Status: corev1.PendingStatus},
					}))
				})

				It("should set the status to pending", func() {
					Expect(updatedCluster.Status.Status).To(Equal(corev1.PendingStatus), updatedCluster.Status.Message)
				})

				It("should create the component resources", func() {
					eksCluster := &eksv1alpha1.EKS{}
					test.ExpectCreate(0, eksCluster)
					Expect(eksCluster.Spec.Cluster).To(Equal(cluster.Ownership()))

					ng1 := &eksv1alpha1.EKSNodeGroup{}
					test.ExpectCreate(1, ng1)
					Expect(ng1.Spec.Cluster).To(Equal(cluster.Ownership()))

					ng2 := &eksv1alpha1.EKSNodeGroup{}
					test.ExpectCreate(2, ng2)
					Expect(ng1.Spec.Cluster).To(Equal(cluster.Ownership()))

					kubernetes := &clustersv1.Kubernetes{}
					test.ExpectCreate(3, kubernetes)
					Expect(kubernetes.Spec.Cluster).To(Equal(cluster.Ownership()))
				})

				When("the components can not be created", func() {
					BeforeEach(func() {
						cluster.Spec.Configuration = v1beta1.JSON{}
					})

					It("should fail and not requeue", func() {
						Expect(reconcileErr).To(HaveOccurred())
					})
				})

				When("creating a component fails", func() {
					BeforeEach(func() {
						test.Client.CreateReturnsOnCall(0, errors.New("some random error"))
					})

					It("should requeue", func() {
						Expect(test.Client.CreateCallCount()).To(Equal(1))
						_, obj, _ := test.Client.CreateArgsForCall(0)
						Expect(obj).To(BeAssignableToTypeOf(&eksv1alpha1.EKS{}))

						test.ExpectRequeue(reconcileResult, reconcileErr)
					})
				})

				When("getting a component fails", func() {
					BeforeEach(func() {
						eks := eksv1alpha1.NewEKS(name.Name, name.Namespace)
						eks.Labels = map[string]string{controllerstest.LabelGetError: "some error"}
						test.Objects = append(test.Objects, eks)
					})

					It("should requeue", func() {
						test.ExpectRequeue(reconcileResult, reconcileErr)
					})
				})

				When("updating an existing component fails", func() {
					BeforeEach(func() {
						test.Client.UpdateReturnsOnCall(0, errors.New("some random error"))
						test.Objects = append(test.Objects, eksv1alpha1.NewEKS(name.Name, name.Namespace))
					})

					It("should requeue", func() {
						Expect(test.Client.UpdateCallCount()).To(Equal(1))
						_, obj, _ := test.Client.UpdateArgsForCall(0)
						Expect(obj).To(BeAssignableToTypeOf(&eksv1alpha1.EKS{}))

						test.ExpectRequeue(reconcileResult, reconcileErr)
					})
				})

				When("applying the cluster configuration on an existing component fails", func() {
					BeforeEach(func() {
						cluster.Spec.Configuration = v1beta1.JSON{}
						test.Objects = append(test.Objects, eksv1alpha1.NewEKS(name.Name, name.Namespace))
					})

					It("should not requeue", func() {
						Expect(reconcileErr).To(HaveOccurred())
					})
				})

				When("a component is complete", func() {
					BeforeEach(func() {
						eks := eksv1alpha1.NewEKS(name.Name, name.Namespace)
						eks.Labels = map[string]string{"cluster.clusters.kore.appvia.io/ResourceVersion": cluster.ResourceVersion}
						eks.Status.Status = corev1.SuccessStatus
						test.Objects = append(test.Objects, eks)
					})

					It("should update the component status on the cluster", func() {
						updatedCluster := &clustersv1.Cluster{}
						test.ExpectStatusUpdate(0, updatedCluster)

						component, _ := updatedCluster.Status.Components.GetComponent("EKS/testName")
						Expect(component.Status).To(Equal(corev1.SuccessStatus))
					})
				})

				When("the kubernetes object is complete", func() {
					var kubernetes *clustersv1.Kubernetes
					var updatedCluster *clustersv1.Cluster

					BeforeEach(func() {
						kubernetes = clustersv1.NewKubernetes(name.Name, name.Namespace)
						kubernetes.Status.Status = corev1.SuccessStatus
						kubernetes.Status.CaCertificate = "testCaCert"
						kubernetes.Status.Endpoint = "testEndpoint"
						kubernetes.Status.APIEndpoint = "testAPIEndpoint"
						kubernetes.Labels = map[string]string{"cluster.clusters.kore.appvia.io/ResourceVersion": cluster.ResourceVersion}
						test.Objects = append(test.Objects, kubernetes)
					})

					JustBeforeEach(func() {
						updatedCluster = &clustersv1.Cluster{}
						test.ExpectStatusUpdate(0, updatedCluster)
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
				})

				When("all components are complete", func() {
					BeforeEach(func() {
						eks := eksv1alpha1.NewEKS(name.Name, name.Namespace)
						eks.Status.Status = corev1.SuccessStatus
						eks.Labels = map[string]string{"cluster.clusters.kore.appvia.io/ResourceVersion": cluster.ResourceVersion}
						ng1 := eksv1alpha1.NewEKSNodeGroup(name.Name+"-ng1", name.Namespace)
						ng1.Status.Status = corev1.SuccessStatus
						ng1.Labels = map[string]string{"cluster.clusters.kore.appvia.io/ResourceVersion": cluster.ResourceVersion}
						ng2 := eksv1alpha1.NewEKSNodeGroup(name.Name+"-ng2", name.Namespace)
						ng2.Status.Status = corev1.SuccessStatus
						ng2.Labels = map[string]string{"cluster.clusters.kore.appvia.io/ResourceVersion": cluster.ResourceVersion}
						kubernetes := clustersv1.NewKubernetes(name.Name, name.Namespace)
						kubernetes.Status.Status = corev1.SuccessStatus
						kubernetes.Labels = map[string]string{"cluster.clusters.kore.appvia.io/ResourceVersion": cluster.ResourceVersion}
						test.Objects = append(test.Objects, eks, ng1, ng2, kubernetes)
					})

					It("should update the status on the cluster", func() {
						updatedCluster := &clustersv1.Cluster{}
						test.ExpectStatusUpdate(0, updatedCluster)

						Expect(updatedCluster.Status.Status).To(Equal(corev1.SuccessStatus))
					})

					It("should finish the reconciliation", func() {
						Expect(reconcileErr).ToNot(HaveOccurred())
						Expect(reconcileResult).To(Equal(reconcile.Result{}))
					})
				})

				When("a component failed", func() {
					BeforeEach(func() {
						eks := eksv1alpha1.NewEKS(name.Name, name.Namespace)
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
						test.ExpectStatusUpdate(0, updatedCluster)
						Expect(updatedCluster.Status.Status).To(Equal(corev1.FailureStatus))
						component, _ := updatedCluster.Status.Components.GetComponent("EKS/testName")
						Expect(component.Status).To(Equal(corev1.FailureStatus))
						Expect(component.Message).To(Equal("[Failure] some component - some error: some error detail"))
					})

					It("should fail and not requeue", func() {
						Expect(reconcileErr).To(HaveOccurred())
					})
				})

				When("updating the status fails", func() {
					BeforeEach(func() {
						test.StatusClient.UpdateReturnsOnCall(0, fmt.Errorf("some error"))
					})
					It("should retry", func() {
						Expect(test.StatusClient.UpdateCallCount()).To(Equal(1))
						test.ExpectRequeue(reconcileResult, reconcileErr)
					})
				})
			})
		})
	})
})
