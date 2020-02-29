/*
 * Copyright (C) 2019  Rohith Jayawardene <gambol99@gmail.com>
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

package namespaceclaims

import (
	"context"
	"fmt"
	"strings"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	core "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers"
	ctrl "github.com/appvia/kore/pkg/controllers/management/kubernetes"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	log "github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	// finalizerName is our finalizer name
	finalizerName = "namespaceclaims.kore.appvia.io"
)

// Reconcile is responsible for reconciling the resource
func (a *nsCtrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	// --- Logic ---
	// we have a client to the remote kubernetes cluster
	// we check if the team has a team namespace policy
	// we need to check the namespace is there and if not create it
	// we need to check the rolebinding exists and if not create it
	// we need to check that all the members of the team are in the binding
	// we set ourselves as the finalizer on the resource if not there already
	// we set the status of the resource to Success and the Phase is Installed
	// we sit back, relax and contain our smug smile

	logger := log.WithFields(log.Fields{
		"name":      request.Name,
		"namespace": request.Namespace,
	})
	logger.Debug("attempting to reconcile the nameresource claim")

	// @step: retrieve the resource from the api
	resource := &clustersv1.NamespaceClaim{}
	if err := a.mgr.GetClient().Get(ctx, request.NamespacedName, resource); err != nil {
		if !kerrors.IsNotFound(err) {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}
	original := resource.DeepCopy()

	// @step: create a finalizer for the resource
	finalizer := kubernetes.NewFinalizer(a.mgr.GetClient(), finalizerName)

	if resource.GetDeletionTimestamp() != nil {
		if finalizer.IsDeletionCandidate(resource) {
			return a.Delete(request)
		}

		return reconcile.Result{}, nil
	}

	result, err := func() (reconcile.Result, error) {
		// @step: ensure the namespace is for a cluster you own
		if resource.Spec.Cluster.Namespace != resource.Namespace {
			resource.Status.Status = core.FailureStatus
			resource.Status.Conditions = []core.Condition{{
				Detail:  "access denied",
				Message: "Cannot create namespace on cluster not owned by you",
			}}

			return reconcile.Result{}, nil
		}

		// @step: check the status of the cluster
		cluster := &clustersv1.Kubernetes{}
		if err := a.mgr.GetClient().Get(context.Background(), types.NamespacedName{
			Name:      resource.Spec.Cluster.Name,
			Namespace: resource.Spec.Cluster.Namespace,
		}, cluster); err != nil {
			if !kerrors.IsNotFound(err) {
				logger.WithError(err).Error("Trying to retrieve the cluster")

				return reconcile.Result{}, err
			}

			// @checkpoint the cluster is not available yet
			resource.Status.Status = core.PendingStatus
			resource.Status.Conditions = []core.Condition{{
				Detail:  "cluster does not exist",
				Message: "No cluster: " + resource.Spec.Cluster.Name + " exists for this team",
			}}

			// @TODO we probably need a way of escaping this loop?
			return reconcile.Result{RequeueAfter: 3 * time.Minute}, nil
		}

		// @step: check the overall status of the cluster
		if cluster.Status.Status == core.PendingStatus {
			resource.Status.Status = core.PendingStatus
			resource.Status.Conditions = []core.Condition{{
				Detail:  "cluster provisioning is still pending",
				Message: "Cluster " + resource.Spec.Cluster.Name + " is still pending",
			}}

			return reconcile.Result{RequeueAfter: 3 * time.Minute}, nil
		}

		// @step: check the provisioning status
		status, found := cluster.Status.Components.GetStatus(ctrl.ComponentClusterCreate)
		if !found {
			logger.Warn("cluster does not have a status on the provisioning yet")

			resource.Status.Status = core.PendingStatus
			resource.Status.Conditions = []core.Condition{{
				Detail:  "cluster is pending, retrying later",
				Message: "Cluster: " + resource.Spec.Cluster.Name + " is still pending",
			}}

			return reconcile.Result{RequeueAfter: 3 * time.Minute}, nil
		}
		switch status.Status {
		case core.PendingStatus:
			logger.Warn("cluster provision is not successful yet, waiting")

			resource.Status.Status = core.PendingStatus
			resource.Status.Conditions = []core.Condition{{
				Detail:  "cluster provisioning is still pending",
				Message: "Cluster " + resource.Spec.Cluster.Name + " is still pending",
			}}

			return reconcile.Result{RequeueAfter: 3 * time.Minute}, nil
		case core.SuccessStatus:
		default:
			resource.Status.Status = core.PendingStatus
			resource.Status.Conditions = []core.Condition{{
				Detail:  "cluster has failed to provision, will retry",
				Message: "Cluster " + resource.Spec.Cluster.Name + " is in a failed state",
			}}

			return reconcile.Result{RequeueAfter: 3 * time.Minute}, nil
		}

		// @step: create credentials for the cluster
		cc, err := controllers.CreateClientFromSecret(context.Background(), a.mgr.GetClient(),
			resource.Spec.Cluster.Namespace, resource.Spec.Cluster.Name)
		if err != nil {
			logger.WithError(err).Error("trying to create client from cluster secret")

			return reconcile.Result{}, err
		}

		// @step: ensure the namespace claim exists
		if err := kubernetes.EnsureNamespace(ctx, cc, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:        resource.Spec.Name,
				Labels:      resource.Spec.Labels,
				Annotations: resource.Spec.Annotations,
			},
		}); err != nil {
			logger.WithError(err).Error("trying to provision the namespace in remote cluster")

			return reconcile.Result{}, err
		}

		// @step we need to check the rolebinding exists and if not create it
		logger.Debug("ensuring the binding to the namespace admin exists")

		// @logic we use either the default namespace admin or the role specified in the spec
		roleName := ClusterRoleName
		if resource.Spec.DefaultTeamRole != "" {
			roleName = resource.Spec.DefaultTeamRole
		}

		if !resource.Spec.DisableTeamMemberInheritance {
			binding := &rbacv1.RoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:      RoleBindingName,
					Namespace: resource.Spec.Name,
					Labels:    map[string]string{kore.Label("owned"): "true"},
				},
				RoleRef: rbacv1.RoleRef{
					APIGroup: rbacv1.GroupName,
					Kind:     "ClusterRole",
					Name:     roleName,
				},
			}

			// @step: retrieve all the users in the team
			users, err := a.Teams().Team(request.Namespace).Members().List(ctx)
			if err != nil {
				logger.WithError(err).Error("trying to retrieve a list of members in the team")

				return reconcile.Result{}, err
			}

			logger.WithField(
				"users", len(users.Items),
			).Debug("found the x members in the team")

			for _, x := range users.Items {
				binding.Subjects = append(binding.Subjects, rbacv1.Subject{
					APIGroup: rbacv1.GroupName,
					Kind:     rbacv1.UserKind,
					Name:     x.Spec.Username,
				})
			}

			// @step: ensuring the binding exists
			if _, err := kubernetes.CreateOrUpdate(ctx, cc, binding); err != nil {
				logger.WithError(err).Error("trying to ensure the namespace team binding")

				return reconcile.Result{}, err
			}
		}

		if len(resource.Spec.UsersRoles) > 0 {
			logger.WithField(
				"users", len(resource.Spec.UsersRoles),
			).Debug("found the x members in user roles")

			// @step: we need to build up a bunch of bindings to roles
			bindings := make(map[string]*rbacv1.RoleBinding)

			for _, user := range resource.Spec.UsersRoles {
				for _, x := range user.Roles {
					binding, found := bindings[x]
					if !found {
						bindings[x] = &rbacv1.RoleBinding{
							ObjectMeta: metav1.ObjectMeta{
								Name:      fmt.Sprintf("kore:userroles:%s", x),
								Namespace: resource.Spec.Name,
								Labels:    map[string]string{kore.Label("owned"): "true"},
							},
							RoleRef: rbacv1.RoleRef{
								APIGroup: rbacv1.GroupName,
								Kind:     "ClusterRole",
								Name:     x,
							},
							Subjects: []rbacv1.Subject{{
								APIGroup: rbacv1.GroupName,
								Kind:     rbacv1.UserKind,
								Name:     user.Username,
							}},
						}
					} else {
						binding.Subjects = append(binding.Subjects, rbacv1.Subject{
							APIGroup: rbacv1.GroupName,
							Kind:     rbacv1.UserKind,
							Name:     user.Username,
						})
					}
				}
			}

			// @step: iterate the above and apply them
			for _, binding := range bindings {
				if _, err := kubernetes.CreateOrUpdate(ctx, cc, binding); err != nil {
					logger.WithError(err).Error("trying to ensure the binding")

					return reconcile.Result{}, err
				}
			}

			// @step: we need to remove any bindings which are not longer valid
			list := &rbacv1.RoleBindingList{}
			if err := cc.List(ctx, list, client.InNamespace(resource.Spec.Name)); err != nil {
				logger.WithError(err).Error("trying to list current role bindings")

				return reconcile.Result{}, err
			}

			for _, x := range list.Items {
				if strings.HasPrefix(x.Name, "kore:userroles:") {
					if _, found := bindings[x.Name]; !found {
						logger.WithField(
							"binding", x.Name,
						).Debug("attempting to delete the role binding as no longer referenced")

						if err := kubernetes.DeleteIfExists(ctx, cc, &x); err != nil {
							logger.WithError(err).Error("trying to delete the role binding")

							return reconcile.Result{}, err
						}
					}
				}
			}
		}

		resource.Status.Status = core.SuccessStatus
		resource.Status.Conditions = []core.Condition{}

		return reconcile.Result{}, nil
	}()
	if err != nil {
		logger.WithError(err).Error("trying to reconcile the nameresource claim")

		resource.Status.Status = core.FailureStatus
		resource.Status.Conditions = []core.Condition{{
			Message: "Failed trying to reconcile the nameresource claim",
			Detail:  err.Error(),
		}}
	} else {
		if finalizer.NeedToAdd(resource) {
			if err := finalizer.Add(resource); err != nil {
				logger.WithError(err).Error("trying to add the finalizer")

				return reconcile.Result{}, err
			}

			return reconcile.Result{Requeue: true}, nil
		}
	}

	if err := a.mgr.GetClient().Status().Patch(ctx, resource, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("trying to update the status of the resource")

		return reconcile.Result{}, err
	}

	return result, err
}
