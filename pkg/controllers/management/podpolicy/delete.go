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

package podpolicy

import (
	"context"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is responsible for deleting any pop which were created
func (a pspCtrl) Delete(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithFields(log.Fields{
		"name":      request.Namespace,
		"namespace": request.Name,
	})
	logger.Debug("attempting to delete the object")

	// @step: retrieve the type from the api
	policy := &clustersv1.ManagedPodSecurityPolicy{}
	if err := a.mgr.GetClient().Get(context.Background(), request.NamespacedName, policy); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	// @step: create a finalizer and check if we are deleting
	finalizer := kubernetes.NewFinalizer(a.mgr.GetClient(), finalizerName)

	if err := finalizer.Remove(policy); err != nil {
		log.WithError(err).Error("trying to remove the pod security policy")

		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
