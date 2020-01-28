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

	kubev1 "github.com/appvia/kube-operator/pkg/apis/kube/v1"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/hub"

	"github.com/gambol99/hub-utils/pkg/finalizers"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Delete is responsible for removig the namespace claim any remote configuration
func (r *ReconcileNamespaceClaim) Delete(
	ctx context.Context,
	cl client.Client,
	client kubernetes.Interface,
	resource *kubev1.NamespaceClaim) error {

	phase := resource.Status.Phase

	rlog := log.WithValues(
		"namespace.name", resource.Spec.Name,
		"resource.name", resource.Name,
		"resource.namespace", resource.Namespace,
		"resource.team", resource.GetLabels()[hub.Label("team")],
		"resource.workspace", resource.GetLabels()[hub.Label("workspace")],
	)

	// @step: check if we are the current finalizer
	finalizer := finalizers.NewFinalizer(cl, FinalizerName)
	if !finalizer.IsDeletionCandidate(resource) {
		rlog.WithValues(
			"finalizers", resource.GetFinalizers(),
		).Info("skipping finalization until others have cleaned up")

		return nil
	}

	err := func() error {
		// @step: check the current phase of the claim and if not 'CREATED' we can forgo
		if phase != PhaseInstalled {
			log.Info("skipping the finalizer as the resource was never installed")
			return nil
		}

		log.Info("deleting the namespaceclaim from the cluster")

		// @step: delete the namespace
		if err := client.CoreV1().Namespaces().Delete(resource.Spec.Name, &metav1.DeleteOptions{}); err != nil {
			if kerrors.IsNotFound(err) {
				// @logic - cool we having nothing to do then
				resource.Status.Status = metav1.StatusSuccess
				return nil
			}
			resource.Status.Conditions = []corev1.Condition{{Message: "failed to delete the namespace in cluster"}}

			return err
		}

		return nil
	}()
	if err != nil {
		resource.Status.Status = corev1.FailureStatus
		resource.Status.Conditions = []corev1.Condition{{
			Detail:  err.Error(),
			Message: "failed to delete namespaceclaim",
		}}

		return err
	}

	// @step: remove the finalizer if one and allow the resource it be deleted
	if err := finalizer.Remove(resource); err != nil {
		resource.Status.Status = corev1.FailureStatus
		resource.Status.Conditions = []corev1.Condition{{
			Detail:  err.Error(),
			Message: "failed to remove the finalizer",
		}}

		return err
	}

	return nil
}
