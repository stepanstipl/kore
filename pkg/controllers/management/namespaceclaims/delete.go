/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 * 
 * This file is part of kore.
 * 
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * 
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * 
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package namespaceclaims

import (
	"context"

	kubev1 "github.com/appvia/kube-operator/pkg/apis/kube/v1"

	corev1 "github.com/appvia/hub-apis/pkg/apis/core/v1"
	"github.com/appvia/hub-apiserver/pkg/hub"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is responsible for removig the namespace claim any remote configuration
func (a *nsCtrl) Delete(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithFields(log.Fields{
		"name": resource.Name,
		"namespace": resource.Namespace,
	})
	logger.Debug("attempting to delete the resource")

	// @step: check if we are the current finalizer
	finalizer := kubernetes.NewFinalizer(a.mgr.GetClient(), FinalizerName)

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
