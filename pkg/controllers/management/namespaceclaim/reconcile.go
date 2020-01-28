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

	core "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/hub"
	kubev1 "github.com/appvia/kore/pkg/apis/kube/v1"

	"github.com/gambol99/hub-utils/pkg/finalizers"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Update is resposible for reconciling the resource
func (r *ReconcileNamespaceClaim) Update(
	ctx context.Context,
	cl client.Client,
	cc kubernetes.Interface,
	resource *kubev1.NamespaceClaim) error {

	uid := string(resource.GetUID())

	// --- Logic ---
	// we have a client to the remote kubernetes cluster
	// we check if the team has a team namespace policy
	// we need to check the namespace is there and if not create it
	// we need to check the rolebinding exists and if not create it
	// we need to check that all the members of the team are in the binding
	// we set ourselves as the finalizer on the resource if not there already
	// we set the status of the resource to Success and the Phase is Installed
	// we sit back, relax and contain our smug smile

	team := HubLabel(resource, "team")
	workspace := HubLabel(resource, "workspace")

	rlog := log.WithValues(
		"namespace.name", resource.Spec.Name,
		"resource.name", resource.Name,
		"resource.namespace", resource.Namespace,
		"team.name", team,
		"workspace.name", workspace,
		"uid", uid)

	//
	// @step: check the namespace exists, if not create it, else update it
	//
	annotations := resource.Spec.Annotations
	if annotations == nil {
		annotations = make(map[string]string, 0)
	}
	annotations["hub.appvia.io/uid"] = uid

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        resource.Spec.Name,
			Labels:      resource.Spec.Labels,
			Annotations: annotations,
		},
	}

	err := func() error {
		if _, err := cc.CoreV1().Namespaces().Get(resource.Spec.Name, metav1.GetOptions{}); err != nil {
			if !kerrors.IsNotFound(err) {
				return fmt.Errorf("failed to check namespace exists: %s", err)
			}
			rlog.Info("creating the namespace resource in cluster")

			// else we need to create the namespace
			if _, err := cc.CoreV1().Namespaces().Create(namespace); err != nil {
				return fmt.Errorf("failed to create namespace: %s", err)
			}
		} else {
			rlog.Info("updating the namespace resource in cluster")

			if _, err := cc.CoreV1().Namespaces().Update(namespace); err != nil {
				return fmt.Errorf("failed to update namespace: %s", err)
			}
		}

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

	//
	// @step we need to check the rolebinding exists and if not create it
	//
	binding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      RoleBindingName,
			Namespace: resource.Spec.Name,
			Labels: map[string]string{
				"hub.appvia.io/team":      resource.GetLabels()[hub.Label("team")],
				"hub.appvia.io/workspace": resource.GetLabels()[hub.Label("workspace")],
			},
			Annotations: map[string]string{"hub.appvia.io/uid": uid},
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
		resource.Status.Status = core.FailureStatus
		resource.Status.Conditions = []core.Condition{{
			Detail:  err.Error(),
			Message: "failed to retrieve a list of users",
			Code:    http.StatusServiceUnavailable,
		}}

		return err
	}
	rlog.WithValues(
		"users", len(membership.Items),
	).Info("found the x members in the team")

	for _, x := range membership.Items {
		binding.Subjects = append(binding.Subjects, rbacv1.Subject{
			APIGroup: rbacv1.GroupName,
			Kind:     rbacv1.UserKind,
			Name:     x.Spec.Username,
		})
	}
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

	//
	// @step: set ourselves as the finalizer on the resource if not there already
	//
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
