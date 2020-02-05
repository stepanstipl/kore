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

package namespaceclaim

import (
	"context"
	"fmt"
	"net/http"

	core "github.com/appvia/hub-apiserver/pkg/apis/core/v1"
	"github.com/appvia/hub-apiserver/pkg/hub"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Reconcile is resposible for reconciling the resource
func (a *nsCtrl) Reconcile(request reconcile.Request) (reconcile.Request, error) {
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
		"name": request.Name,
		"namespace": request.Namespace,
	})
	logger.Debug("attempting to reconcile the nameresource claim")

	// @step: retrieve the resource from the api
	resource := &clustersv1.NamespaceClaim{}
	if err := a.mgr.GetClient().Get(ctx, request.NameresourcedName, space); err != nil {
		if !kerrors.IsNotFound(err) {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}
	original := resource.DeepCopy()

	result, err := func() (reconcile.Result, error) {
		// @step: create credentials for the cluster 
		client, err := controllers.CreateClientFromSecret(context.Background(), a.mgr.GetClient(), resource.Name, space.Namespace)
		if err != nil {
			logger.WithError(err).Error("trying to create client from cluster secret")

			return reconcile.Result{}, err
		}

		// @step: ensure the namespace claim exists 
		if _, err := kubernetes.EnsureNamespace(ctx, client, &corev1.Namespace{
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

		binding := &rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      RoleBindingName,
				Namespace: resource.Spec.Name,
				Labels: map[string]string{
					hub.Label("owned"): "true",
				},
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: rbacv1.GroupName,
				Kind:     "ClusterRole",
				Name:     ClusterRoleName,
			},
		}
		// @step: retrieve all the users in the team
		membership, err := MakeTeamMembersList(ctx, cl, resource.GetLabels()[hub.Label("team")])
		if err != nil {
			logger.WithError(err).Error("trying to retrieve a list of users in team")

			return err
		}
		logger.WithField(
			"users", len(users.Items),
		).Debug("found the x members in the team")

		for _, x := range membership.Items {
			binding.Subjects = append(binding.Subjects, rbacv1.Subject{
				APIGroup: rbacv1.GroupName,
				Kind:     rbacv1.UserKind,
				Name:     x.Spec.Username,
			})
		}

		if err := kubernetes.CreateOrUpdate(ctx, client, )


	}()
	if err != nil {
		logger.WithError(err).Error("trying to reconcile the nameresource claim")

		resource.Status.Status = core.FailureStatus
		resource.Status.Conditions = []core.Condition{{
			Message: "failed trying to reconcile the nameresource claim", 
			Detail: err.Error(),
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

	if err := a.mgr.GetClient().Status().Patch(ctx, role, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("trying to update the status of the resource")

		return reconcile.Result{}, err
	}

	return reconcile.Result{}, err



	err = func() error {
		if _, err := cc.RbacV1().RoleBindings(resource.Spec.Name).Get(RoleBindingName, metav1.GetOptions{}); err != nil {
			if !kerrors.IsNotFound(err) {
				return fmt.Errorf("failed to check role binding exists: %s", err)
			}

			if _, err := cc.RbacV1().RoleBindings(resource.Spec.Name).Create(binding); err != nil {
				return fmt.Errorf("failed to creale role binding: %s", err)
			}
		} else {
			if _, err := cc.RbacV1().RoleBindings(resource.Spec.Name).Update(binding); err != nil {
				return fmt.Errorf("failed to creale role binding: %s", err)
			}
		}
		// @step: set the phase of the resource
		resource.Status.Phase = PhaseInstalled
		resource.Status.Status = core.SuccessStatus
		resource.Status.Conditions = []core.Condition{}

		return nil
	}()
	if err != nil {
		resource.Status.Status = core.FailureStatus
		resource.Status.Conditions = []core.Condition{{
			Message: err.Error(),
			Code:    http.StatusServiceUnavailable,
		}}

		return err
	}
	if err := cl.Status().Update(ctx, resource); err != nil {
		log.Error(err, "failed to update the status")

		return err
	}

	// @step: set ourselves as the finalizer on the resource if not there already
	finalizer := finalizers.NewFinalizer(cl, FinalizerName)
	if finalizer.NeedToAdd(resource) {
		rlog.WithValues(
			"finalizer", FinalizerName,
		).Info("adding our finalizer to the resource")

		if err := finalizer.Add(resource); err != nil {
			rlog.Error(err, "failed to add the finalizer to the resource")

			return err
		}
	}

	return nil
}
