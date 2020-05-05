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

package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	ekscc "github.com/appvia/kore/pkg/controllers/cloud/aws/eks"
	gkecc "github.com/appvia/kore/pkg/controllers/cloud/gcp/gke"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	rbacv1 "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "kubernetes.clusters.kore.appvia.io"
	// ComponentClusterCreate is the component name
	ComponentClusterCreate = "Cluster Provisioned"
	// ComponentAPIAuthProxy is the component name
	ComponentAPIAuthProxy = "SSO Authentication"
	// ComponentClusterAppMan is the component name for the Kore Cluster application manager
	ComponentClusterAppMan = "Kore Cluster Manager"
	// ComponentClusterUsers is the component name for Kore team users of this cluster
	ComponentClusterUsers = "Kore Cluster Users"
	// ComponentClusterRoles is the component name for inherited RBAC team roles
	ComponentClusterRoles = "Kore Cluster Roles"
	// ComponentClusterNetworkPolicy is the component name for Kubernetes network policy
	ComponentClusterNetworkPolicy = "Kubernetes Network Policy"
)

// Reconcile is the entrypoint for the reconciliation logic
func (a k8sCtrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to reconcile the kubernetes cluster")

	// @step: retrieve the type from the api
	object := &clustersv1.Kubernetes{}
	if err := a.mgr.GetClient().Get(context.Background(), request.NamespacedName, object); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	// @step: keep the original object from api
	original := object.DeepCopy()
	finalizer := kubernetes.NewFinalizer(a.mgr.GetClient(), finalizerName)

	if finalizer.IsDeletionCandidate(object) {
		return a.Delete(context.Background(), object)
	}

	team := object.Namespace
	token := &configv1.Secret{}

	result, err := func() (reconcile.Result, error) {
		object.Status.Status = corev1.PendingStatus

		logger.Debugf("retrieving the cluster credentials from secret %s/%s", object.Namespace, object.Name)

		// @step: retrieve the provider credentials secret
		account, err := controllers.GetConfigSecret(context.Background(),
			a.mgr.GetClient(),
			object.Namespace,
			object.Name)

		if err != nil {
			if !kerrors.IsNotFound(err) {
				logger.WithError(err).Error("trying to retrieve the admin token, will retry")

				return reconcile.Result{RequeueAfter: 5 * time.Minute}, nil
			}

			logger.Debugf("no credentials found from cluster %s", err)

			// it wasn't found - is the cluster provider backed?
			if !kore.IsProviderBacked(object) {
				object.Status.Status = corev1.WarningStatus
				object.Status.Components.SetCondition(corev1.Component{
					Name:    ComponentClusterCreate,
					Message: "Credentials for cluster not available yet",
					Status:  corev1.WarningStatus,
				})

				return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
			}

			// the secret doesn't appear to be available yet - it's either been pushed
			// of the provider hasn't finished
			if kore.IsProviderBacked(object) {

				// @step: check if the provider is still pending
				u, err := kore.ToUnstructuredFromOwnership(object.Spec.Provider)
				if err != nil {
					logger.WithError(err).Error("invalid group version kind in resource")

					object.Status.Status = corev1.FailureStatus
					object.Status.Components.SetCondition(corev1.Component{
						Name:    ComponentClusterCreate,
						Message: "Invalid cloud provider reference",
						Status:  corev1.FailureStatus,
					})

					return reconcile.Result{}, err
				}

				if found, err := kubernetes.GetIfExists(context.Background(), a.mgr.GetClient(), u); err != nil {
					logger.WithError(err).Error("trying to get the cloud provider resource")

					return reconcile.Result{}, err
				} else if !found {
					logger.WithError(err).Error("cloud provider resource does not exist")

					object.Status.Status = corev1.FailureStatus
					object.Status.Components.SetCondition(corev1.Component{
						Name:    ComponentClusterCreate,
						Message: "Cloud provider resource does not exist",
						Status:  corev1.FailureStatus,
					})

					return reconcile.Result{}, err
				}

				// @check if the cloud provider has failed
				if a.Config().EnableClusterProviderCheck {
					ready, err := a.CheckProviderStatus(context.Background(), object)
					if err != nil {
						logger.WithError(err).Warn("error getting cluster provider status")

						return reconcile.Result{}, err
					}
					if !ready {
						logger.Debug("cloud provider is not ready yet, waiting")

						return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
					}
				}

				object.Status.Status = corev1.PendingStatus
				object.Status.Components.SetCondition(corev1.Component{
					Name:    ComponentClusterCreate,
					Message: "Waiting for cluster to be provisioned",
					Status:  corev1.PendingStatus,
				})

				return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
			}
		}
		token = account

		// @step: create a client for the remote cluster
		client, err := kubernetes.NewRuntimeClientFromConfigSecret(token)
		if err != nil {
			logger.WithError(err).Error("trying to create client from credentials secret")

			object.Status.Components.SetCondition(corev1.Component{
				Name:    ComponentClusterCreate,
				Message: "Unable to access cluster using provided cluster credentials",
				Detail:  err.Error(),
				Status:  corev1.FailureStatus,
			})

			return reconcile.Result{}, err
		}

		// @step: ensure the kube-api proxy is deployed
		// @TODO need to move this out into something else, but for now its cool
		logger.Debug("ensure the api proxy service is provisioned")

		if err := a.EnsureAPIService(context.Background(), client, object); err != nil {
			logger.WithError(err).Error("trying to provision the api service")

			object.Status.Status = corev1.FailureStatus
			object.Status.Components.SetCondition(corev1.Component{
				Name:    ComponentAPIAuthProxy,
				Message: "Failed to provision the Service proxy",
				Detail:  err.Error(),
				Status:  corev1.FailureStatus,
			})

			return reconcile.Result{RequeueAfter: 2 * time.Minute}, err
		}

		object.Status.Status = corev1.FailureStatus
		object.Status.Components.SetCondition(corev1.Component{
			Name:    ComponentAPIAuthProxy,
			Message: "Service proxy is running and available",
			Status:  corev1.SuccessStatus,
		})

		/*
			if original.Status.Endpoint == "" {
				return reconcile.Result{Requeue: true}, nil
			}
		*/

		// @step: ensure all cluster components are deployed
		components, err := a.EnsureClusterman(context.Background(), client, object)
		if err != nil {
			logger.WithError(err).Error("trying to provision the clusterappman service")

			object.Status.Status = corev1.FailureStatus
			object.Status.Components.SetCondition(corev1.Component{
				Name:    ComponentClusterAppMan,
				Message: "Kore failed to deploy kore Cluster Manager component",
				Detail:  err.Error(),
				Status:  corev1.FailureStatus,
			})

			return reconcile.Result{RequeueAfter: 2 * time.Minute}, err
		}
		object.Status.Components.SetCondition(corev1.Component{
			Name:    ComponentClusterAppMan,
			Message: "Cluster manager component is running and available",
			Status:  corev1.SuccessStatus,
		})
		// Provide visibility of remote cluster apps
		for _, component := range *components {
			object.Status.Components.SetCondition(*component)
		}

		// @step: we start by reconcile the cluster admins if any
		if len(object.Spec.ClusterUsers) > 0 {
			logger.Debug("attempting to reconcile cluster users for the cluster")

			for name, users := range ClusterUserRolesToMap(object.Spec.ClusterUsers) {
				binding := &rbacv1.ClusterRoleBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name: "kore:clusterusers:" + name,
						Labels: map[string]string{
							kore.Label("owned"): "true",
						},
					},
					RoleRef: rbacv1.RoleRef{
						APIGroup: rbacv1.SchemeGroupVersion.Group,
						Kind:     "ClusterRole",
						Name:     name,
					},
				}
				sort.Strings(users)

				for _, user := range users {
					binding.Subjects = append(binding.Subjects, rbacv1.Subject{
						APIGroup: rbacv1.SchemeGroupVersion.Group,
						Kind:     "User",
						Name:     user,
					})
				}

				logger.WithField("name", name).Debug("ensuring the cluster role binding exists")

				// @step: ensure the binding for the role exists
				_, err = kubernetes.CreateOrUpdateManagedClusterRoleBinding(context.Background(), client, binding)
				if err != nil {
					logger.WithError(err).Error("trying to enforce the cluster role binding for cluster users")

					object.Status.Status = corev1.FailureStatus
					object.Status.Components.SetCondition(corev1.Component{
						Name:    ComponentClusterUsers,
						Message: "Failed to provision one of more cluster user roles",
						Status:  corev1.FailureStatus,
						Detail:  err.Error(),
					})

					return reconcile.Result{}, err
				}
			}

			object.Status.Status = corev1.SuccessStatus
			object.Status.Components.SetCondition(corev1.Component{
				Name:    ComponentClusterUsers,
				Message: "Cluster users have been provisioned",
				Status:  corev1.SuccessStatus,
			})
		} else {
			logger.Debug("removing any bindings for cluster users as non defined")
			if err := kubernetes.DeleteBindingsWithPrefix(context.Background(), client, "kore:clusterusers:"); err != nil {
				logger.WithError(err).Error("trying to delete any cluster user role bindings")

				return reconcile.Result{}, err
			}
		}

		// @step: check if team users are inherited
		if object.Spec.InheritTeamMembers {
			logger.Debug("attempting to reconcile the inherited users")

			if object.Spec.DefaultTeamRole == "" {
				object.Status.Components.SetCondition(corev1.Component{
					Name:    ComponentClusterRoles,
					Message: "Team inheritance enabled but no default team role defined",
					Status:  corev1.WarningStatus,
				})
				object.Status.Status = corev1.WarningStatus

			} else {
				// @step: we need to provision binding and all users
				binding := &rbacv1.ClusterRoleBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name: "kore:team:inherited",
						Labels: map[string]string{
							kore.Label("owned"): "true",
						},
					},
					RoleRef: rbacv1.RoleRef{
						APIGroup: rbacv1.SchemeGroupVersion.Group,
						Kind:     "ClusterRole",
						Name:     object.Spec.DefaultTeamRole,
					},
				}
				// @step: retrieve a list of all the members in the team
				members, err := a.Teams().Team(team).Members().List(context.Background())
				if err != nil {
					logger.WithError(err).Error("trying to retrieve members of team")

					object.Status.Components.SetCondition(corev1.Component{
						Name:    ComponentClusterRoles,
						Message: "Failed to apply all team members to inherited role",
						Detail:  err.Error(),
						Status:  corev1.WarningStatus,
					})
					object.Status.Status = corev1.WarningStatus

					return reconcile.Result{}, err
				}

				log.WithField("users", len(members.Items)).Debug("adding x members to the cluster default role")

				for _, user := range members.Items {
					if !user.Spec.Disabled {
						binding.Subjects = append(binding.Subjects, rbacv1.Subject{
							APIGroup: rbacv1.SchemeGroupVersion.Group,
							Kind:     "User",
							Name:     user.Spec.Username,
						})
					}
				}

				// @step: ensure the binding for the role exists
				_, err = kubernetes.CreateOrUpdateManagedClusterRoleBinding(context.Background(), client, binding)
				if err != nil {
					logger.WithError(err).Error("trying to enforce the cluster role binding for cluster users")

					object.Status.Status = corev1.FailureStatus
					object.Status.Components.SetCondition(corev1.Component{
						Name:    ComponentClusterRoles,
						Message: "Failed to provision one of more inherited user roles",
						Status:  corev1.FailureStatus,
						Detail:  err.Error(),
					})

					return reconcile.Result{}, err
				}

				object.Status.Components.SetCondition(corev1.Component{
					Name:    ComponentClusterRoles,
					Message: "Provision one of more inherited users on cluster",
					Status:  corev1.SuccessStatus,
				})
			}
		} else {
			logger.Debug("removing any inherited users reconcilation as inheritence is disabled")

			if err := kubernetes.DeleteIfExists(context.Background(), client, &rbacv1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: "kore:team:inherited",
				},
			}); err != nil {
				logger.WithError(err).Error("trying to delete any inherited role binding")

				return reconcile.Result{}, err
			}
		}

		// @step: is default network block enabled?
		if object.Spec.EnableDefaultTrafficBlock != nil && *object.Spec.EnableDefaultTrafficBlock {
			// @step: ensure the remote cluster has the traffic blocked
			logger.Debug("ensuring that network policies are enabled by default on all namespaces")

			object.Status.Components.SetCondition(corev1.Component{
				Name:    ComponentClusterNetworkPolicy,
				Message: "Network policy has been provisioned",
				Status:  corev1.SuccessStatus,
			})
		} else {
			logger.Debug("skipping default network policy and feature is disabled")
		}

		object.Status.APIEndpoint = token.Spec.Data["endpoint"]
		object.Status.CaCertificate = token.Spec.Data["ca.crt"]
		//object.Status.Endpoint = a.APIHostname(object)
		object.Status.Status = corev1.SuccessStatus

		object.Status.Components.SetCondition(corev1.Component{
			Name:    ComponentClusterCreate,
			Message: "Cluster has been successfully provisioned",
			Status:  corev1.SuccessStatus,
		})

		return reconcile.Result{RequeueAfter: 30 * time.Minute}, nil
	}()
	if err != nil {
		object.Status.Status = corev1.FailureStatus
	}
	if err == nil {
		// check if we need to add the finalizer
		if finalizer.NeedToAdd(object) {
			if err := finalizer.Add(object); err != nil {
				logger.WithError(err).Error("trying to add ourself as a finalizer")

				return reconcile.Result{}, err
			}

			return reconcile.Result{Requeue: true}, nil
		}
	}

	// @step: the resource has been reconcile, update the status
	if err := a.mgr.GetClient().Status().Patch(context.Background(), object, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("trying to update the kubernetes status")

		return reconcile.Result{}, err
	}

	return result, nil
}

// CheckProviderStatus checks the status of the provider behind the cluster
func (a k8sCtrl) CheckProviderStatus(ctx context.Context, resource *clustersv1.Kubernetes) (bool, error) {
	logger := log.WithFields(log.Fields{
		"name":      resource.Name,
		"namespace": resource.Namespace,
		"provider":  resource.Spec.Provider.Kind,
	})
	logger.Debug("checking the status of the cloud provider")

	key := types.NamespacedName{
		Namespace: resource.Spec.Provider.Namespace,
		Name:      resource.Spec.Provider.Name,
	}
	switch resource.Spec.Provider.Kind {
	case "EKS":
		p := &eks.EKS{}

		if err := a.mgr.GetClient().Get(ctx, key, p); err != nil {
			logger.WithError(err).Error("trying to retrieve the eks cluster from api")
		}

		if p.Status.Conditions == nil {
			logger.Debugf("Cluster %s does not have a status yet", resource.Name)

			return false, nil
		}

		// @check if we have a provider status for provisioning yet
		status, found := p.Status.Conditions.GetComponent(ekscc.ComponentClusterCreator)
		if !found {
			logger.Warn("eks cluster has no status, throwing an error until we know")

			return false, nil
		}

		if status.Status == corev1.FailureStatus {
			message := status.Message
			if message == "" {
				message = "EKS Cluster has failed to provision correctly"
			}

			resource.Status.Components.SetCondition(corev1.Component{
				Detail:  status.Detail,
				Name:    ComponentClusterCreate,
				Message: message,
				Status:  corev1.FailureStatus,
			})
			resource.Status.Status = corev1.FailureStatus

			return false, errors.New("cloud provider is in a failed state")
		}

		return true, nil
	case "GKE":
		p := &gke.GKE{}
		if err := a.mgr.GetClient().Get(ctx, key, p); err != nil {
			logger.WithError(err).Error("trying to retrieve the gke cluster from api")
		}

		if p.Status.Status == corev1.FailureStatus {
			resource.Status.Components.SetCondition(corev1.Component{
				Name:    ComponentClusterCreate,
				Message: "GKE Cluster is in a failing state",
				Status:  corev1.FailureStatus,
			})
			resource.Status.Status = corev1.FailureStatus

			return false, fmt.Errorf("gke cluster is in a failing state")
		}

		if p.Status.Conditions == nil {
			logger.Debugf("Cluster %s does not have a status yet", resource.Name)
			return false, nil
		}

		// @check if we have a provider status for provisioning yet
		status, found := p.Status.Conditions.GetComponent(gkecc.ComponentClusterCreator)
		if !found {
			logger.Warn("gke cluster has no status, throwing an error until we know")

			return false, nil
		}

		switch status.Status {
		case corev1.SuccessStatus:
			return true, nil
		case corev1.FailureStatus:
			message := status.Message
			if message == "" {
				message = "GKE Cluster has failed to provision correctly"
			}

			resource.Status.Components.SetCondition(corev1.Component{
				Detail:  status.Detail,
				Name:    ComponentClusterCreate,
				Message: message,
				Status:  corev1.FailureStatus,
			})
			resource.Status.Status = corev1.FailureStatus

			return false, errors.New("cloud provider is in a failed state")

		case corev1.PendingStatus:
			resource.Status.Components.SetCondition(corev1.Component{
				Name:    ComponentClusterCreate,
				Message: "GKE Cluster is a pending state",
				Status:  corev1.PendingStatus,
			})
			logger.Debug("cloud provider is in a pending state")

			return false, nil

		default:
			resource.Status.Components.SetCondition(corev1.Component{
				Name:    ComponentClusterCreate,
				Message: "GKE Cluster is an unknown state",
				Status:  corev1.PendingStatus,
			})
			logger.Debug("cloud provider is in a unknown state")

			return false, nil
		}
	}

	return false, fmt.Errorf("cluster built with an unknown provider: %s", resource.Spec.Provider.Kind)
}
