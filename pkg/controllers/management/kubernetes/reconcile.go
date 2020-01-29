/*
 * Copyright (C) 2019 Rohith Jayawardene <gambol99@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package kubernetes

import (
	"context"
	"sort"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/hub"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	core "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const finalizerName = "kubernetes.clusters.hub.appvia.io"

// Reconcile is the entrypoint for the reconcilation logic
func (a k8sCtrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to renconcile the kubernetes cluster")

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
	token := &core.Secret{}

	result, err := func() (reconcile.Result, error) {
		object.Status.Status = corev1.PendingStatus

		// @step: check for the kubernetes admin token
		key := types.NamespacedName{
			Name:      object.Name,
			Namespace: object.Namespace,
		}

		logger.Debug("retrieving the cluster credentials from secret")

		// @step: retrieve the admin token
		if err := a.mgr.GetClient().Get(context.Background(), key, token); err != nil {
			if !kerrors.IsNotFound(err) {
				logger.WithError(err).Error("trying to retrieve the admin token, will retry")

				return reconcile.Result{RequeueAfter: 5 * time.Minute}, nil
			}
			logger.Debug("no credentials found from cluster")

			// it wasn't found - is the cluster provider backed?
			if !hub.IsProviderBacked(object) {
				object.Status.Status = corev1.WarningStatus
				object.Status.Components.SetCondition(corev1.Component{
					Name:    "provision",
					Message: "credentials for cluster not available yet",
					Status:  corev1.WarningStatus,
				})

				return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
			}

			// the secret doesn't appear to be available yet - it's either been pushed
			// of the provider hasn't finished
			if hub.IsProviderBacked(object) {

				// @step: check if the provider is still pending
				u, err := hub.ToUnstructuredFromOwnership(object.Spec.Provider)
				if err != nil {
					logger.WithError(err).Error("invalid group version kind in resource")

					object.Status.Status = corev1.FailureStatus
					object.Status.Components.SetCondition(corev1.Component{
						Name:    "provision",
						Message: "invalid provider cloud provider reference",
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
						Name:    "provision",
						Message: "cloud provider resource does not exist",
						Status:  corev1.FailureStatus,
					})

					return reconcile.Result{}, err
				}

				object.Status.Status = corev1.FailureStatus
				object.Status.Components.SetCondition(corev1.Component{
					Name:    "provision",
					Message: "credentials are not available",
					Status:  corev1.FailureStatus,
				})

				return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
			}
		}

		// @step: create a client for the remote cluster
		client, err := kubernetes.NewRuntimeClientFromSecret(token)
		if err != nil {
			logger.WithError(err).Error("trying to create client from credentials secret")

			object.Status.Status = corev1.FailureStatus
			object.Status.Components.SetCondition(corev1.Component{
				Name:    "provision",
				Message: "unable to create client from cluster credentials",
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
				Name:    "api",
				Message: "unable to create client from cluster credentials",
				Detail:  err.Error(),
				Status:  corev1.FailureStatus,
			})

			return reconcile.Result{RequeueAfter: 2 * time.Minute}, err
		}

		object.Status.Status = corev1.FailureStatus
		object.Status.Components.SetCondition(corev1.Component{
			Name:    "api",
			Message: "api service proxy is running and available",
			Status:  corev1.SuccessStatus,
		})

		// @step: we start by reconcile the cluster admins if any
		if len(object.Spec.ClusterUsers) > 0 {
			logger.Debug("attempting to reconcile cluster users for the cluster")

			for name, users := range ClusterUserRolesToMap(object.Spec.ClusterUsers) {
				binding := &rbacv1.ClusterRoleBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name: "hub:clusterusers:" + name,
						Labels: map[string]string{
							hub.Label("owned"): "true",
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
						Name:    "cluster_users",
						Message: "failed to provision one of more cluster user roles",
						Status:  corev1.FailureStatus,
						Detail:  err.Error(),
					})

					return reconcile.Result{}, err
				}
			}

			object.Status.Status = corev1.SuccessStatus
			object.Status.Components.SetCondition(corev1.Component{
				Name:    "cluster_users",
				Message: "cluster users have been provisioned",
				Status:  corev1.SuccessStatus,
			})
		} else {
			logger.Debug("removing any bindings for cluster users as non defined")
			if err := kubernetes.DeleteBindingsWithPrefix(context.Background(), client, "hub:clusterusers:"); err != nil {
				logger.WithError(err).Error("trying to delete any cluster user role bindings")

				return reconcile.Result{}, err
			}
		}

		// @step: check if team users are inherited
		if object.Spec.InheritTeamMembers {
			logger.Debug("attempting to reconcile the inherited users")

			if object.Spec.DefaultTeamRole == "" {
				object.Status.Components.SetCondition(corev1.Component{
					Name:    "inherited_roles",
					Message: "team inheritence enabled but no default team role defined",
					Status:  corev1.WarningStatus,
				})
				object.Status.Status = corev1.WarningStatus

			} else {
				// @step: we need to provision binding and all users
				binding := &rbacv1.ClusterRoleBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name: "hub:team:inherited",
						Labels: map[string]string{
							hub.Label("owned"): "true",
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
						Name:    "inherited_roles",
						Message: "failed to apply all team members to inherited role",
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
						Name:    "inherited_roles",
						Message: "failed to provision one of more inherited user roles",
						Status:  corev1.FailureStatus,
						Detail:  err.Error(),
					})

					return reconcile.Result{}, err
				}

				object.Status.Components.SetCondition(corev1.Component{
					Name:    "inherited_roles",
					Message: "provision one of more inherited users on cluster",
					Status:  corev1.SuccessStatus,
				})
			}
		} else {
			logger.Debug("removing any inherited users reconcilation as inheritence is disabled")

			if err := kubernetes.DeleteIfExists(context.Background(), client, &rbacv1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: "hub:team:inherited",
				},
			}); err != nil {
				logger.WithError(err).Error("trying to delete any inherited role binding")

				return reconcile.Result{}, err
			}
		}

		// @step: is default network block enabled?
		if object.Spec.EnabledDefaultTrafficBlock != nil && *object.Spec.EnabledDefaultTrafficBlock {
			// @step: ensure the remote cluster has the traffic blocked
			logger.Debug("ensuring that network policies are enabled by default on all namespaces")

			object.Status.Components.SetCondition(corev1.Component{
				Name:    "network_policy",
				Message: "network policy has been provisioned",
				Status:  corev1.SuccessStatus,
			})
		} else {
			logger.Debug("skipping default network policy and feature is disabled")
		}

		object.Status.APIEndpoint = string(token.Data["endpoint"])
		object.Status.CaCertificate = string(token.Data["ca.crt"])
		//object.Status.Endpoint = a.APIHostname(object)
		object.Status.Status = corev1.SuccessStatus

		object.Status.Components.SetCondition(corev1.Component{
			Name:    "provision",
			Message: "cluster has been successfully provisioned",
			Status:  corev1.SuccessStatus,
		})

		return reconcile.Result{RequeueAfter: 30 * time.Minute}, nil
	}()
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
