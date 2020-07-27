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
	"time"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"
	"github.com/Azure/go-autorest/autorest/to"
	aksv1alpha1 "github.com/appvia/kore/pkg/apis/aks/v1alpha1"
	"github.com/appvia/kore/pkg/controllers/helpers"
	"github.com/appvia/kore/pkg/kore"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type resourceGroupComponent struct {
	AKSCluster        *aksv1alpha1.AKS
	ResourceGroupName string
}

func newResourceGroupComponent(aks *aksv1alpha1.AKS, resourceGroupName string) resourceGroupComponent {
	return resourceGroupComponent{
		AKSCluster:        aks,
		ResourceGroupName: resourceGroupName,
	}
}

func (c resourceGroupComponent) ComponentName() string {
	return "Resource Group"
}

func (c resourceGroupComponent) Reconcile(ctx kore.Context) (reconcile.Result, error) {
	helper := helpers.NewAKSHelper(c.AKSCluster)

	client, err := helper.CreateResourceGroupsClient(ctx)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to create AKS API client: %w", err)
	}

	existing, err := c.getResourceGroupIfExists(ctx, client)
	if err != nil {
		return reconcile.Result{}, err
	}

	if existing == nil {
		group, err := client.CreateOrUpdate(ctx, c.ResourceGroupName, resources.Group{
			Location: to.StringPtr(c.AKSCluster.Spec.Location),
			Tags:     *to.StringMapPtr(c.AKSCluster.Spec.Tags),
		})

		if err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to create resource group %s: %w", c.ResourceGroupName, err)
		}

		if to.String(group.Properties.ProvisioningState) == "Succeeded" {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}

	switch to.String(existing.Properties.ProvisioningState) {
	case "Succeeded":
		return reconcile.Result{}, nil
	default:
		ctx.Logger().WithField("provisioningState", to.String(existing.Properties.ProvisioningState)).Debug("current state of the Resource Group")

		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}
}

func (c resourceGroupComponent) Delete(ctx kore.Context) (reconcile.Result, error) {
	helper := helpers.NewAKSHelper(c.AKSCluster)

	client, err := helper.CreateResourceGroupsClient(ctx)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to create AKS API client: %w", err)
	}

	existing, err := c.getResourceGroupIfExists(ctx, client)
	if err != nil {
		return reconcile.Result{}, err
	}

	if existing == nil {
		return reconcile.Result{}, nil
	}

	ctx.Logger().WithField("provisioningState", to.String(existing.Properties.ProvisioningState)).Debug("current state of the Resource Group")

	_, err = client.Delete(ctx, c.ResourceGroupName)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to delete Resource Group: %w", err)
	}

	return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
}

func (c resourceGroupComponent) getResourceGroupIfExists(ctx kore.Context, client resources.GroupsClient) (*resources.Group, error) {
	existing, err := client.Get(ctx, c.ResourceGroupName)
	if err != nil {
		if isNotFound(existing.Response) {
			return nil, nil
		}

		return nil, fmt.Errorf("getting existing Resource Group failed: %w", err)
	}

	if existing.Properties == nil {
		return nil, fmt.Errorf("getting existing Resource Group failed: properties was empty")
	}

	return &existing, nil
}

func (c resourceGroupComponent) SetComponent(_ *corev1.Component) {
}
