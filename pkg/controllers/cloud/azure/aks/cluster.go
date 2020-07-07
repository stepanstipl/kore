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

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2020-06-01/containerservice"
	"github.com/Azure/go-autorest/autorest/to"
	aksv1alpha1 "github.com/appvia/kore/pkg/apis/aks/v1alpha1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers/helpers"
	"github.com/appvia/kore/pkg/kore"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ClusterComponent struct {
	AKSCluster *aksv1alpha1.AKS
}

func (c ClusterComponent) Reconcile(ctx kore.Context) (reconcile.Result, error) {
	helper := helpers.NewAKSHelper(c.AKSCluster)

	client, err := helper.CreateClusterClient(ctx)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to create AKS API client: %w", err)
	}

	existing, err := c.getClusterIfExists(ctx, client)
	if err != nil {
		return reconcile.Result{}, err
	}

	if existing == nil {
		res, err := c.create(ctx, client)
		if err != nil || res.Requeue || res.RequeueAfter > 0 {
			return res, err
		}

		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	} else {
		updated, res, err := c.update(ctx, client, existing)

		if err != nil || res.Requeue || res.RequeueAfter > 0 {
			return res, err
		}

		if updated {
			return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
		}
	}

	switch to.String(existing.ProvisioningState) {
	case "Succeeded":
		return reconcile.Result{}, nil
	default:
		ctx.Logger().WithField("provisioningState", to.String(existing.ProvisioningState)).Debug("current state of the AKS cluster")
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}
}

func (c ClusterComponent) create(ctx kore.Context, client containerservice.ManagedClustersClient) (reconcile.Result, error) {
	var agentPoolProfiles []containerservice.ManagedClusterAgentPoolProfile
	for _, nodePool := range c.AKSCluster.Spec.NodePools {
		agentPoolProfiles = append(agentPoolProfiles, c.createAgentPoolProfile(nodePool))
	}

	properties := &containerservice.ManagedClusterProperties{
		APIServerAccessProfile: &containerservice.ManagedClusterAPIServerAccessProfile{
			AuthorizedIPRanges:   to.StringSlicePtr(c.AKSCluster.Spec.APIServerAuthorizedIPRanges),
			EnablePrivateCluster: to.BoolPtr(c.AKSCluster.Spec.PrivateClusterEnabled),
		},
		AgentPoolProfiles: &agentPoolProfiles,
		DNSPrefix:         to.StringPtr(c.AKSCluster.Spec.DNSPrefix),
		EnableRBAC:        to.BoolPtr(true),
		KubernetesVersion: to.StringPtr(c.AKSCluster.Spec.Version),
		NetworkProfile: &containerservice.NetworkProfileType{
			NetworkPlugin: containerservice.NetworkPlugin(c.AKSCluster.Spec.NetworkPlugin),
			NetworkPolicy: containerservice.NetworkPolicy(c.AKSCluster.Spec.NetworkPolicy),
		},
		NodeResourceGroup:       to.StringPtr(c.resourceGroup()),
		EnablePodSecurityPolicy: to.BoolPtr(c.AKSCluster.Spec.EnablePodSecurityPolicy),
	}

	if c.AKSCluster.Spec.LinuxProfile != nil {
		var publicKeys []containerservice.SSHPublicKey
		for _, publicKey := range c.AKSCluster.Spec.LinuxProfile.SSHPublicKeys {
			publicKeys = append(publicKeys, containerservice.SSHPublicKey{KeyData: to.StringPtr(publicKey)})
		}

		if c.AKSCluster.Spec.LinuxProfile.AdminUsername != "" || len(publicKeys) > 0 {
			properties.LinuxProfile = &containerservice.LinuxProfile{}
			if c.AKSCluster.Spec.LinuxProfile.AdminUsername != "" {
				properties.LinuxProfile.AdminUsername = to.StringPtr(c.AKSCluster.Spec.LinuxProfile.AdminUsername)
			}
			if len(publicKeys) > 0 {
				properties.LinuxProfile.SSH = &containerservice.SSHConfiguration{PublicKeys: &publicKeys}
			}
		}
	}

	if c.AKSCluster.Spec.WindowsProfile != nil {
		if c.AKSCluster.Spec.WindowsProfile.AdminUsername != "" || c.AKSCluster.Spec.WindowsProfile.AdminPassword != "" {
			properties.WindowsProfile = &containerservice.ManagedClusterWindowsProfile{}
			if c.AKSCluster.Spec.WindowsProfile.AdminUsername != "" {
				properties.WindowsProfile.AdminUsername = to.StringPtr(c.AKSCluster.Spec.WindowsProfile.AdminUsername)
			}

			if c.AKSCluster.Spec.WindowsProfile.AdminPassword != "" {
				properties.WindowsProfile.AdminPassword = to.StringPtr(c.AKSCluster.Spec.WindowsProfile.AdminPassword)
			}
		}
	}

	_, err := client.CreateOrUpdate(ctx, c.resourceGroup(), c.AKSCluster.Name, containerservice.ManagedCluster{
		ManagedClusterProperties: properties,
		Identity: &containerservice.ManagedClusterIdentity{
			Type: containerservice.SystemAssigned,
		},
		Sku: &containerservice.ManagedClusterSKU{
			Name: containerservice.ManagedClusterSKUNameBasic,
			Tier: containerservice.Paid,
		},
		Location: to.StringPtr(c.AKSCluster.Spec.Region),
		Tags: map[string]*string{
			kore.Label("owner"): to.StringPtr("true"),
			kore.Label("team"):  to.StringPtr(c.AKSCluster.Namespace),
		},
	})

	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to create AKS cluster: %w", err)
	}

	return reconcile.Result{}, nil
}

func (c ClusterComponent) update(
	ctx kore.Context, client containerservice.ManagedClustersClient, existing *containerservice.ManagedCluster,
) (bool, reconcile.Result, error) {
	return false, reconcile.Result{}, nil
}

func (c ClusterComponent) Delete(ctx kore.Context) (reconcile.Result, error) {
	helper := helpers.NewAKSHelper(c.AKSCluster)

	client, err := helper.CreateClusterClient(ctx)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to create AKS API client: %w", err)
	}

	existing, err := c.getClusterIfExists(ctx, client)
	if err != nil {
		return reconcile.Result{}, err
	}

	if existing == nil {
		return reconcile.Result{}, nil
	}

	ctx.Logger().WithField("provisioningState", to.String(existing.ProvisioningState)).Debug("current state of the AKS cluster")

	_, err = client.Delete(ctx, c.resourceGroup(), c.AKSCluster.Name)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to delete AKS cluster: %w", err)
	}

	return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
}

func (c ClusterComponent) ComponentName() string {
	return "Cluster Creator"
}

func (c ClusterComponent) SetComponent(_ *corev1.Component) {
}

func (c ClusterComponent) getClusterIfExists(ctx kore.Context, client containerservice.ManagedClustersClient) (*containerservice.ManagedCluster, error) {
	existing, err := client.Get(ctx, c.resourceGroup(), c.AKSCluster.Name)
	if err != nil {
		return nil, fmt.Errorf("getting existing AKS cluster failed: %w", err)
	}

	if isNotFound(existing.Response) {
		return nil, nil
	}

	if existing.ManagedClusterProperties == nil {
		return nil, fmt.Errorf("getting existing AKS cluster failed: properties was empty")
	}

	return &existing, nil
}

func (c ClusterComponent) resourceGroup() string {
	return "kore-" + c.AKSCluster.Namespace
}

func (c ClusterComponent) createAgentPoolProfile(nodepool aksv1alpha1.AKSNodePool) containerservice.ManagedClusterAgentPoolProfile {
	taints := []string{}
	for _, taint := range nodepool.Taints {
		taints = append(taints, fmt.Sprintf("%s=%s:%s", taint.Key, taint.Value, taint.Effect))
	}
	return containerservice.ManagedClusterAgentPoolProfile{
		Name:              to.StringPtr(nodepool.Name),
		MinCount:          to.Int32Ptr(int32(nodepool.MinSize)),
		MaxCount:          to.Int32Ptr(int32(nodepool.MaxSize)),
		Count:             to.Int32Ptr(int32(nodepool.Size)),
		MaxPods:           to.Int32Ptr(int32(nodepool.MaxPodsPerNode)),
		VMSize:            containerservice.VMSizeTypes(nodepool.MachineType),
		OsDiskSizeGB:      to.Int32Ptr(int32(nodepool.DiskSize)),
		OsType:            containerservice.OSType(nodepool.ImageType),
		EnableAutoScaling: to.BoolPtr(nodepool.EnableAutoscaler),
		Type:              containerservice.VirtualMachineScaleSets,
		Mode:              containerservice.User,
		NodeImageVersion:  to.StringPtr(nodepool.Version),
		NodeLabels:        *to.StringMapPtr(nodepool.Labels),
		NodeTaints:        &taints,
	}
}
