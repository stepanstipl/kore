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

package aks

import (
	"fmt"

	aksv1alpha1 "github.com/appvia/kore/pkg/apis/aks/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	cc "github.com/appvia/kore/pkg/controllers/components"
	"github.com/appvia/kore/pkg/kore"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "aks.compute.kore.appvia.io"
)

// Reconcile is the entrypoint for the reconciliation logic
func (c *Controller) Reconcile(ctx kore.Context, request reconcile.Request) (reconcile.Result, error) {
	ctx.Logger().Debug("attempting to reconcile the AKS cluster")

	aksCluster := &aksv1alpha1.AKS{}
	if err := ctx.Client().Get(ctx, request.NamespacedName, aksCluster); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, fmt.Errorf("failed to retrieve AKS cluster: %w", err)
	}
	original := aksCluster.DeepCopyObject()

	var caCertificate, clientToken string

	components := controllers.Components{
		cc.NewFinalizer(finalizerName, aksCluster),
		newResourceGroupComponent(aksCluster, resourceGroupName(aksCluster)),
		newClusterComponent(aksCluster, &caCertificate, &clientToken),
		cc.NewClusterBootstrap(func(k kore.Context) (controllers.Bootstrap, error) {
			return NewBootstrapClient(aksCluster, clientToken, caCertificate)
		}),
	}

	res, err := components.Reconcile(ctx, aksCluster)
	if err != nil {
		ctx.Logger().WithError(err).Error("failed to reconcile the AKS cluster")
	}

	if err := ctx.Client().Status().Patch(ctx, aksCluster, client.MergeFrom(original)); err != nil {
		ctx.Logger().WithError(err).Error("failed to update the status of the AKS cluster")

		return reconcile.Result{}, fmt.Errorf("failed to update the status of the AKS cluster: %w", err)
	}

	return res, err
}
