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

	costsv1 "github.com/appvia/kore/pkg/apis/costs/v1beta1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// EnsureCloudInfo makes sure that cloudinfo is present or not in the specified cluster as required
func (a *Controller) EnsureCloudInfo(cost *costsv1.Cost) controllers.EnsureFunc {
	return func(c context.Context) (reconcile.Result, error) {
		ctx := kore.NewContext(c, a.logger, a.mgr.GetClient(), a.Interface)
		cloudinfo := newCloudInfo(ctx, cost)

		required, err := cloudinfo.IsRequired()
		if err != nil {
			a.logger.WithError(err).Error("Trying to determine if cloudinfo required")
			return reconcile.Result{}, err
		}
		if !required {
			a.logger.Info("Cloudinfo not required, ensuring it is NOT present.")
			return cloudinfo.Delete()
		}
		a.logger.Info("Cloudinfo required, ensuring it is present.")
		return cloudinfo.Ensure()
	}
}
