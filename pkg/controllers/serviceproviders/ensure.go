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

package serviceproviders

import (
	"fmt"

	"github.com/appvia/kore/pkg/kore"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"

	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (c *Controller) ensurePending(serviceProvider *servicesv1.ServiceProvider) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		if serviceProvider.Status.Status == "" {
			serviceProvider.Status.Status = corev1.PendingStatus
			return reconcile.Result{Requeue: true}, nil
		}

		if serviceProvider.Status.Status != corev1.PendingStatus {
			serviceProvider.Status.Status = corev1.PendingStatus
		}

		return reconcile.Result{}, nil
	}
}

func (c *Controller) ensureFinalizer(serviceProvider *servicesv1.ServiceProvider, finalizer *kubernetes.Finalizer) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		if finalizer.NeedToAdd(serviceProvider) {
			err := finalizer.Add(serviceProvider)
			if err != nil {
				return reconcile.Result{}, fmt.Errorf("failed to set the finalizer: %w", err)
			}
			return reconcile.Result{Requeue: true}, nil
		}
		return reconcile.Result{}, nil
	}
}
