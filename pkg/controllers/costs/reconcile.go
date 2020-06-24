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

package costs

import (
	"context"
	"reflect"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	costsv1 "github.com/appvia/kore/pkg/apis/costs/v1beta1"
	"github.com/appvia/kore/pkg/controllers"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconcile ensures the state of the costs infrastructure is correct
func (a *Controller) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	// Get a reference to the cost
	cost := &costsv1.Cost{}
	if err := a.mgr.GetClient().Get(ctx, request.NamespacedName, cost); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		a.logger.WithError(err).Error("trying to retrieve cost from api")

		return reconcile.Result{}, err
	}
	original := cost.DeepCopy()
	a.logger.WithField("cost", cost.Name).Info("Reconciling cost configuration")

	result, err := func() (reconcile.Result, error) {
		return controllers.DefaultEnsureHandler.Run(ctx,
			[]controllers.EnsureFunc{
				a.EnsureCloudInfo(cost),
			},
		)
	}()

	if err != nil {
		a.logger.WithError(err).Error("trying to ensure the cost")

		if controllers.IsCriticalError(err) {
			cost.Status.Status = corev1.FailureStatus
			cost.Status.Message = err.Error()
		}
	}

	if !reflect.DeepEqual(cost, original) {
		if err := a.mgr.GetClient().Status().Patch(ctx, cost, client.MergeFrom(original)); err != nil {
			a.logger.WithError(err).Error("trying to patch the cost status")

			return reconcile.Result{}, err
		}
	}

	return result, err
}
