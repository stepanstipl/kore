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
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	aksv1alpha1 "github.com/appvia/kore/pkg/apis/aks/v1alpha1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers/helpers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/jsonutils"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2020-06-01/containerservice"
	"github.com/Azure/go-autorest/autorest/to"
	"gopkg.in/yaml.v2"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type clusterComponent struct {
	AKSCluster    *aksv1alpha1.AKS
	CACertificate *string
	ClientToken   *string
}

func newClusterComponent(aks *aksv1alpha1.AKS, caCertificate, clientToken *string) *clusterComponent {
	return &clusterComponent{
		AKSCluster:    aks,
		CACertificate: caCertificate,
		ClientToken:   clientToken,
	}
}

func (c *clusterComponent) ComponentName() string {
	return "Cluster Creator"
}

func (c *clusterComponent) Reconcile(ctx kore.Context) (reconcile.Result, error) {
	helper := helpers.NewAKSHelper(c.AKSCluster)

	client, err := helper.CreateClustersClient(ctx)
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
		if existing.Fqdn != nil {
			c.AKSCluster.Status.Endpoint = "https://" + to.String(existing.Fqdn)
		}

		if err := c.getCredentials(ctx, client); err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to get cluster credentials: %w", err)
		}

		c.AKSCluster.Status.CACertificate = *c.CACertificate

		return reconcile.Result{}, nil
	default:
		ctx.Logger().WithField("provisioningState", to.String(existing.ProvisioningState)).Debug("current state of the AKS cluster")
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}
}

func (c *clusterComponent) setProperties(properties containerservice.ManagedClusterProperties) containerservice.ManagedClusterProperties {
	var agentPoolProfiles []containerservice.ManagedClusterAgentPoolProfile
	for _, agentPoolProfile := range c.AKSCluster.Spec.AgentPoolProfiles {
		agentPoolProfiles = append(agentPoolProfiles, c.setAgentPoolProfile(agentPoolProfile, properties))
	}

	if properties.APIServerAccessProfile == nil {
		properties.APIServerAccessProfile = &containerservice.ManagedClusterAPIServerAccessProfile{}
	}
	properties.APIServerAccessProfile.AuthorizedIPRanges = to.StringSlicePtr(c.AKSCluster.Spec.AuthorizedIPRanges)
	properties.APIServerAccessProfile.EnablePrivateCluster = to.BoolPtr(c.AKSCluster.Spec.EnablePrivateCluster)

	properties.AgentPoolProfiles = &agentPoolProfiles
	properties.DNSPrefix = to.StringPtr(c.AKSCluster.Spec.DNSPrefix)
	properties.EnableRBAC = to.BoolPtr(true)
	properties.KubernetesVersion = to.StringPtr(c.AKSCluster.Spec.KubernetesVersion)

	if properties.NetworkProfile == nil {
		properties.NetworkProfile = &containerservice.NetworkProfileType{}
	}
	properties.NetworkProfile.NetworkPlugin = containerservice.NetworkPlugin(c.AKSCluster.Spec.NetworkPlugin)
	properties.NetworkProfile.NetworkPolicy = containerservice.NetworkPolicy(to.String(c.AKSCluster.Spec.NetworkPolicy))
	properties.NetworkProfile.LoadBalancerSku = "Standard"

	properties.NodeResourceGroup = to.StringPtr(nodesResourceGroupName(c.AKSCluster))
	properties.EnablePodSecurityPolicy = to.BoolPtr(c.AKSCluster.Spec.EnablePodSecurityPolicy)

	if c.AKSCluster.Spec.LinuxProfile != nil {
		if properties.LinuxProfile == nil {
			properties.LinuxProfile = &containerservice.LinuxProfile{}
		}

		if c.AKSCluster.Spec.LinuxProfile.AdminUsername != "" {
			properties.LinuxProfile.AdminUsername = to.StringPtr(c.AKSCluster.Spec.LinuxProfile.AdminUsername)
		}

		var publicKeys []containerservice.SSHPublicKey
		for _, publicKey := range c.AKSCluster.Spec.LinuxProfile.SSHPublicKeys {
			publicKeys = append(publicKeys, containerservice.SSHPublicKey{KeyData: to.StringPtr(publicKey)})
		}
		if len(publicKeys) > 0 {
			properties.LinuxProfile.SSH = &containerservice.SSHConfiguration{PublicKeys: &publicKeys}
		}
	}

	if c.AKSCluster.Spec.WindowsProfile != nil {
		if properties.WindowsProfile == nil {
			properties.WindowsProfile = &containerservice.ManagedClusterWindowsProfile{}
		}

		if c.AKSCluster.Spec.WindowsProfile.AdminUsername != "" {
			properties.WindowsProfile.AdminUsername = to.StringPtr(c.AKSCluster.Spec.WindowsProfile.AdminUsername)
		}

		if c.AKSCluster.Spec.WindowsProfile.AdminPassword != "" {
			properties.WindowsProfile.AdminPassword = to.StringPtr(c.AKSCluster.Spec.WindowsProfile.AdminPassword)
		}
	}

	return properties
}

func (c *clusterComponent) createOrUpdate(
	ctx kore.Context,
	client containerservice.ManagedClustersClient,
	properties *containerservice.ManagedClusterProperties,
) error {
	_, err := client.CreateOrUpdate(ctx, resourceGroupName(c.AKSCluster), c.AKSCluster.Name, containerservice.ManagedCluster{
		ManagedClusterProperties: properties,
		Identity: &containerservice.ManagedClusterIdentity{
			Type: containerservice.SystemAssigned,
		},
		Sku: &containerservice.ManagedClusterSKU{
			Name: containerservice.ManagedClusterSKUNameBasic,
			Tier: containerservice.Paid,
		},
		Location: to.StringPtr(c.AKSCluster.Spec.Location),
		Tags: map[string]*string{
			"kore.appvia.io:owner": to.StringPtr("true"),
			"kore.appvia.io:team":  to.StringPtr(c.AKSCluster.Namespace),
		},
	})
	return err
}

func (c *clusterComponent) create(ctx kore.Context, client containerservice.ManagedClustersClient) (reconcile.Result, error) {
	properties := c.setProperties(containerservice.ManagedClusterProperties{})

	if err := c.createOrUpdate(ctx, client, &properties); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to create AKS cluster: %w", err)
	}

	return reconcile.Result{}, nil
}

func (c *clusterComponent) update(
	ctx kore.Context, client containerservice.ManagedClustersClient, existing *containerservice.ManagedCluster,
) (bool, reconcile.Result, error) {
	updatedProperties := c.setProperties(*existing.ManagedClusterProperties)

	diff, err := jsonutils.Diff(*existing.ManagedClusterProperties, updatedProperties)
	if err != nil {
		return false, reconcile.Result{}, fmt.Errorf("failed to compare cluster properties: %w", err)
	}

	if bytes.Equal(diff, []byte("{}")) {
		return false, reconcile.Result{}, nil
	}

	ctx.Logger().WithField("diff", string(diff)).Debug("updating the AKS cluster")

	if err := c.createOrUpdate(ctx, client, &updatedProperties); err != nil {
		return true, reconcile.Result{}, fmt.Errorf("failed to update the AKS cluster: %w", err)
	}

	return true, reconcile.Result{}, nil
}

func (c *clusterComponent) getCredentials(ctx kore.Context, client containerservice.ManagedClustersClient) error {
	creds, err := client.ListClusterAdminCredentials(ctx, resourceGroupName(c.AKSCluster), c.AKSCluster.Name)
	if err != nil {
		return err
	}

	if len(*creds.Kubeconfigs) == 0 {
		return fmt.Errorf("the API did not return any Kubernetes configurations: %w", err)
	}

	type kubeConfig struct {
		Clusters []struct {
			Cluster struct {
				CertificateAuthorityData string `yaml:"certificate-authority-data"`
			} `yaml:"cluster"`
		} `yaml:"clusters"`
		Users []struct {
			User struct {
				ClientCertificateData string `yaml:"client-certificate-data"`
				Token                 string `yaml:"token"`
			} `yaml:"user"`
		} `yaml:"users"`
	}

	cfg := &kubeConfig{}
	if err := yaml.Unmarshal(*(*creds.Kubeconfigs)[0].Value, cfg); err != nil {
		return fmt.Errorf("failed to decode Kubernetes configuration: %w", err)
	}

	if len(cfg.Clusters) == 0 {
		return errors.New("no cluster found in Kubernetes configuration")
	}

	var caCertificate []byte
	if cfg.Clusters[0].Cluster.CertificateAuthorityData != "" {
		caCertificate, _ = base64.StdEncoding.DecodeString(cfg.Clusters[0].Cluster.CertificateAuthorityData)
	}
	if len(caCertificate) == 0 {
		return fmt.Errorf("CA certificate is missing or invalid")
	}

	*c.CACertificate = string(caCertificate)

	if len(cfg.Users) == 0 {
		return errors.New("no user found in Kubernetes configuration")
	}

	if cfg.Users[0].User.Token == "" {
		return errors.New("client token is missing")
	}

	*c.ClientToken = cfg.Users[0].User.Token

	return nil
}

func (c *clusterComponent) Delete(ctx kore.Context) (reconcile.Result, error) {
	helper := helpers.NewAKSHelper(c.AKSCluster)

	client, err := helper.CreateClustersClient(ctx)
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

	_, err = client.Delete(ctx, resourceGroupName(c.AKSCluster), c.AKSCluster.Name)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to delete AKS cluster: %w", err)
	}

	return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
}

func (c *clusterComponent) SetComponent(_ *corev1.Component) {
}

func (c *clusterComponent) getClusterIfExists(ctx kore.Context, client containerservice.ManagedClustersClient) (*containerservice.ManagedCluster, error) {
	existing, err := client.Get(ctx, resourceGroupName(c.AKSCluster), c.AKSCluster.Name)
	if err != nil {
		if isNotFound(existing.Response) {
			return nil, nil
		}

		return nil, fmt.Errorf("getting existing AKS cluster failed: %w", err)
	}

	if existing.ManagedClusterProperties == nil {
		return nil, fmt.Errorf("getting existing AKS cluster failed: properties was empty")
	}

	return &existing, nil
}

func (c *clusterComponent) setAgentPoolProfile(
	nodepool aksv1alpha1.AgentPoolProfile, properties containerservice.ManagedClusterProperties,
) containerservice.ManagedClusterAgentPoolProfile {
	taints := []string{}
	for _, taint := range nodepool.NodeTaints {
		taints = append(taints, fmt.Sprintf("%s=%s:%s", taint.Key, taint.Value, taint.Effect))
	}

	profile := containerservice.ManagedClusterAgentPoolProfile{}
	if properties.AgentPoolProfiles != nil {
		for _, existing := range *properties.AgentPoolProfiles {
			if to.String(existing.Name) == nodepool.Name {
				profile = existing
			}
		}
	}

	profile.Name = to.StringPtr(nodepool.Name)
	profile.MinCount = to.Int32Ptr(int32(nodepool.MinCount))
	profile.MaxCount = to.Int32Ptr(int32(nodepool.MaxCount))

	// If autoscaling is enabled, we should only set the initial value
	if nodepool.EnableAutoScaling && profile.Count == nil {
		profile.Count = to.Int32Ptr(int32(nodepool.Count))
	}
	profile.VMSize = containerservice.VMSizeTypes(nodepool.VMSize)
	profile.OsDiskSizeGB = to.Int32Ptr(int32(nodepool.OsDiskSizeGB))
	profile.OsType = containerservice.OSType(nodepool.OsType)
	profile.EnableAutoScaling = to.BoolPtr(nodepool.EnableAutoScaling)
	profile.Type = containerservice.VirtualMachineScaleSets
	profile.Mode = containerservice.AgentPoolMode(nodepool.Mode)
	if nodepool.NodeImageVersion != "" {
		profile.NodeImageVersion = to.StringPtr(nodepool.NodeImageVersion)
	}
	profile.NodeLabels = *to.StringMapPtr(nodepool.NodeLabels)

	// If there are no node taints on either objects, make sure we retain the exact empty value from the existing one
	if (profile.NodeTaints != nil && len(*profile.NodeTaints) > 0) || len(taints) > 0 {
		profile.NodeTaints = &taints
	}
	if nodepool.MaxPods > 0 {
		profile.MaxPods = to.Int32Ptr(int32(nodepool.MaxPods))
	}

	return profile
}
