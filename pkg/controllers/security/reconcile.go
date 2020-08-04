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

package security

import (
	"github.com/appvia/kore/pkg/kore"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconcile is responsible for handling the scanning of a kind
func (c *Controller) Reconcile(ctx kore.Context, request reconcile.Request) (reconcile.Result, error) {
	ctx.Logger().Debug("attempting to reconcile the security scans")

	// @step: retrieve the object from the api
	t := c.srckind.Type.DeepCopyObject()
	if err := ctx.Client().Get(ctx, request.NamespacedName, t); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		ctx.Logger().WithError(err).Error("trying to retrieve from api")

		return reconcile.Result{}, err
	}

	meta, err := kubernetes.GetMeta(t)
	if err != nil {
		return reconcile.Result{}, err
	}

	if !meta.GetDeletionTimestamp().IsZero() {
		return c.Delete(ctx, t)
	}

	switch c.kind {
	case "Cluster":
		err = ctx.Kore().Security().ScanCluster(ctx, ctx.Client(), t.(*clustersv1.Cluster))
	case "Plan":
		err = ctx.Kore().Security().ScanPlan(ctx, ctx.Client(), t.(*configv1.Plan))
	}
	if err != nil {
		ctx.Logger().WithError(err).Error("trying to run security scan")

		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
