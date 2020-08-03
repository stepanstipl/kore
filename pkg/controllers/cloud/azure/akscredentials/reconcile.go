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

package akscredentials

import (
	"context"
	"fmt"

	"github.com/appvia/kore/pkg/utils"

	aksv1alpha1 "github.com/appvia/kore/pkg/apis/aks/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	cc "github.com/appvia/kore/pkg/controllers/components"
	"github.com/appvia/kore/pkg/kore"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "aks.compute.kore.appvia.io"
)

// Reconcile is the entrypoint for the reconciliation logic
func (c *Controller) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	ctx := kore.NewContext(context.Background(), logger, c.client, c)

	logger.Debug("attempting to reconcile the AKS cluster")

	aksCredentials := &aksv1alpha1.AKSCredentials{}
	if err := c.client.Get(ctx, request.NamespacedName, aksCredentials); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, fmt.Errorf("failed to retrieve AKS cluster: %w", err)
	}
	original := aksCredentials.DeepCopyObject()

	components := controllers.Components{
		cc.NewFinalizer(finalizerName, aksCredentials),
	}

	res, err := components.Reconcile(ctx, aksCredentials)
	if err != nil {
		logger.WithError(err).Error("failed to reconcile the AKS cluster")
	}

	// TODO: implement AKS credentials verification
	aksCredentials.Status.Verified = utils.BoolPtr(true)

	if err := c.client.Status().Patch(ctx, aksCredentials, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("failed to update the status of the AKS cluster")

		return reconcile.Result{}, fmt.Errorf("failed to update the status of the AKS cluster: %w", err)
	}

	return res, err
}
