/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package gke

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"

	log "github.com/sirupsen/logrus"
	compute "google.golang.org/api/compute/v0.beta"
	container "google.golang.org/api/container/v1beta1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
)

type gkeClient struct {
	// cm is the compute client
	cm *compute.Service
	// ce in the container engine client
	ce *container.Service
	// credentials are the gke credentials
	credentials *gke.GKECredentials
	// cluster is the gke cluster
	cluster *gke.GKE
}

// NewClient returns a gcp client for us
func NewClient(credentials *gke.GKECredentials, cluster *gke.GKE) (*gkeClient, error) {
	options := []option.ClientOption{option.WithCredentialsJSON([]byte(credentials.Spec.Account))}

	cm, err := compute.NewService(context.Background(), options...)
	if err != nil {
		return nil, err
	}

	ce, err := container.NewService(context.Background(), options...)
	if err != nil {
		return nil, err
	}

	return &gkeClient{
		cm:          cm,
		ce:          ce,
		credentials: credentials,
		cluster:     cluster,
	}, nil
}

// Delete attempts to delete the cluster from gke
func (g *gkeClient) Delete(ctx context.Context) error {
	logger := log.WithFields(log.Fields{
		"cluster":   g.cluster.Name,
		"namespace": g.cluster.Namespace,
		"project":   g.credentials.Spec.Project,
		"region":    g.credentials.Spec.Region,
	})
	logger.Info("attempting to delete the cluster fomr gcp")

	found, err := g.Exists()
	if err != nil {
		logger.WithError(err).Error("trying to check for the cluster")

		return err
	}
	if !found {
		return nil
	}

	cluster, _, err := g.GetCluster()
	if err != nil {
		logger.WithError(err).Error("trying to retrieve the cluster")

		return err
	}
	path := fmt.Sprintf("projects/%s/locations/%s/clusters/%s",
		g.credentials.Spec.Project,
		g.credentials.Spec.Region,
		cluster.Name)

	// @step: check for any ongoing operation
	id, found, err := g.FindOperation(ctx, "DELETE_CLUSTER", "clusters", cluster.Name)
	if err != nil {
		logger.WithError(err).Error("trying to check for current operations")

		return err
	}
	if !found {
		operation, err := g.ce.Projects.Locations.Clusters.Delete(path).Do()
		if err != nil {
			logger.WithError(err).Error("trying to delete the cluster")

			return err
		}
		id = operation.Name
	}

	logger.Info("waiting for the operation to complete or fail")

	if err := g.WaitOnOperation(ctx, id, 30*time.Second, 10*time.Minute); err != nil {
		logger.WithError(err).Error("trying to wait for operaion to complete")

		return err
	}
	logger.Info("gke cluster has been deleted")

	return nil
}

// Create is used to create the cluster in gcp
func (g *gkeClient) Create(ctx context.Context) (*container.Cluster, error) {
	logger := log.WithFields(log.Fields{
		"cluster":   g.cluster.Name,
		"namespace": g.cluster.Namespace,
		"project":   g.credentials.Spec.Project,
		"region":    g.credentials.Spec.Region,
	})
	logger.Info("attempting to create the gke cluster")

	// @step: we create the definitions
	def, err := g.CreateDefinition()
	if err != nil {
		logger.WithError(err).Error("attempting to create the cluster definition")

		return nil, err
	}

	// @step: looking for any ongoing operation
	id, found, err := g.FindOperation(ctx, "CREATE_CLUSTER", "clusters", g.cluster.Name)
	if err != nil {
		return nil, err
	}
	if !found {
		// @step: we request the cluster
		ticket, err := g.CreateCluster(ctx, def)
		if err != nil {
			logger.WithError(err).Error("attempting to request the cluster")

			return nil, err
		}
		id = ticket.Name
	}

	// @step: wait for the google to finish
	interval := time.Duration(10) * time.Second
	timeout := time.Duration(20) * time.Minute

	// @step: we wait for it to finish
	if err := g.WaitOnOperation(ctx, id, interval, timeout); err != nil {
		logger.WithError(err).Error("attempting to wait for operation to complete")

		return nil, err
	}

	// @step: retrieve the state of the cluster via api
	gc, _, err := g.GetCluster()
	if err != nil {
		logger.WithError(err).Error("retrieving gke cluster details")

		return nil, err
	}

	return gc, nil
}

// Update is called to update the cluster
func (g *gkeClient) Update(ctx context.Context) error {
	logger := log.WithFields(log.Fields{
		"name": g.cluster.Name,
	})
	logger.Info("attempting to update the cluster")

	_, err := g.CreateUpdateDefinition()
	if err != nil {
		log.WithError(err).Error("creating the update request")

		return err
	}

	return nil
}

// CreateUpdateDefinition returns a cluster update definition
func (g *gkeClient) CreateUpdateDefinition() (*container.UpdateClusterRequest, error) {
	return &container.UpdateClusterRequest{
		ProjectId: g.credentials.Spec.Project,
	}, nil
}

// CreateDefinition returns a cluster definition
func (g *gkeClient) CreateDefinition() (*container.CreateClusterRequest, error) {
	// @step: retrieve a list of location to place this gke cluster
	locations, err := g.Locations()
	if err != nil {
		return nil, err
	}

	cluster := g.cluster

	resource := &container.Cluster{
		Name:                  cluster.Name,
		Description:           cluster.Spec.Description,
		InitialClusterVersion: cluster.Spec.Version,

		AddonsConfig: &container.AddonsConfig{
			CloudRunConfig: &container.CloudRunConfig{
				Disabled: true,
			},
			IstioConfig: &container.IstioConfig{
				Auth:     "AUTH_NONE",
				Disabled: !cluster.Spec.EnableIstio,
			},
			HttpLoadBalancing: &container.HttpLoadBalancing{
				Disabled: !cluster.Spec.EnableHTTPLoadBalancer,
			},
			HorizontalPodAutoscaling: &container.HorizontalPodAutoscaling{
				Disabled: !cluster.Spec.EnableAutoscaler,
			},
			KubernetesDashboard: &container.KubernetesDashboard{
				Disabled: true,
			},
			NetworkPolicyConfig: &container.NetworkPolicyConfig{
				Disabled: false,
			},
		},

		BinaryAuthorization:     &container.BinaryAuthorization{Enabled: false},
		LegacyAbac:              &container.LegacyAbac{Enabled: false},
		NetworkPolicy:           &container.NetworkPolicy{Enabled: true, Provider: "CALICO"},
		PodSecurityPolicyConfig: &container.PodSecurityPolicyConfig{Enabled: true},
		Locations:               locations,
		ShieldedNodes:           &container.ShieldedNodes{Enabled: cluster.Spec.EnableShieldedNodes},

		MaintenancePolicy: &container.MaintenancePolicy{
			Window: &container.MaintenanceWindow{
				DailyMaintenanceWindow: &container.DailyMaintenanceWindow{
					StartTime: cluster.Spec.MaintenanceWindow,
				},
			},
		},

		MasterAuth: &container.MasterAuth{
			ClientCertificateConfig: &container.ClientCertificateConfig{
				IssueClientCertificate: false,
			},
		},

		IpAllocationPolicy: &container.IPAllocationPolicy{
			ClusterIpv4CidrBlock: cluster.Spec.ClusterIPV4Cidr,
			CreateSubnetwork:     false,
			ServicesIpv4Cidr:     cluster.Spec.ServicesIPV4Cidr,
			SubnetworkName:       cluster.Spec.Subnetwork,
			UseIpAliases:         true,
		},

		MonitoringService: func() string {
			if cluster.Spec.EnableStackDriverLogging {
				return "monitoring.googleapis.com/kubernetes"
			}
			return ""
		}(),
		LoggingService: func() string {
			if cluster.Spec.EnableStackDriverMetrics {
				return "logging.googleapis.com/kubernetes"
			}
			return ""
		}(),

		NodePools: []*container.NodePool{
			{
				Name: "compute",
				Autoscaling: &container.NodePoolAutoscaling{
					Autoprovisioned: false,
					Enabled:         cluster.Spec.EnableAutoscaler,
					MaxNodeCount:    cluster.Spec.MaxSize,
					MinNodeCount:    1,
				},
				Config: &container.NodeConfig{
					DiskSizeGb:  cluster.Spec.DiskSize,
					ImageType:   cluster.Spec.ImageType,
					MachineType: cluster.Spec.MachineType,
					OauthScopes: []string{
						"https://www.googleapis.com/auth/compute",
						"https://www.googleapis.com/auth/devstorage.read_only",
						"https://www.googleapis.com/auth/logging.write",
						"https://www.googleapis.com/auth/monitoring",
					},
					Preemptible: false,
					Tags:        []string{cluster.Name},
				},
				InitialNodeCount: cluster.Spec.Size,
				Locations:        locations,
				Management: &container.NodeManagement{
					AutoRepair:  cluster.Spec.EnableAutorepair,
					AutoUpgrade: cluster.Spec.EnableAutoupgrade,
				},
				MaxPodsConstraint: &container.MaxPodsConstraint{
					MaxPodsPerNode: 110,
				},
				Version: cluster.Spec.Version,
			},
		},
	}

	if cluster.Spec.EnablePrivateNetwork {
		resource.PrivateCluster = true
		resource.PrivateClusterConfig = &container.PrivateClusterConfig{
			EnablePrivateEndpoint: false,
			EnablePrivateNodes:    true,
			MasterIpv4CidrBlock:   cluster.Spec.MasterIPV4Cidr,
		}

		if len(cluster.Spec.AuthorizedMasterNetworks) > 0 {
			resource.MasterAuthorizedNetworksConfig = &container.MasterAuthorizedNetworksConfig{
				CidrBlocks: []*container.CidrBlock{},
				Enabled:    true,
			}
			for _, x := range cluster.Spec.AuthorizedMasterNetworks {
				resource.MasterAuthorizedNetworksConfig.CidrBlocks = append(resource.MasterAuthorizedNetworksConfig.CidrBlocks, &container.CidrBlock{
					CidrBlock:   x.CIDR,
					DisplayName: x.Name,
				})
			}
		} else {
			resource.MasterAuthorizedNetworksConfig = &container.MasterAuthorizedNetworksConfig{
				CidrBlocks: []*container.CidrBlock{},
				Enabled:    false,
			}
		}
	}

	// @step: fill in the master auth if required
	if cluster.Spec.AuthorizedMasterNetworks == nil {
		cluster.Spec.AuthorizedMasterNetworks = make([]*gke.AuthorizedNetwork, 0)
	}

	return &container.CreateClusterRequest{
		ProjectId: g.credentials.Spec.Project,
		Cluster:   resource,
	}, nil
}

// GetCluster returns a cluster config
func (g *gkeClient) GetCluster() (*container.Cluster, bool, error) {
	clusters, err := g.GetClusters()
	if err != nil {
		return nil, false, err
	}
	for _, x := range clusters {
		if x.Name == g.cluster.Name {
			return x, true, nil
		}
	}

	return nil, false, nil
}

// GetClusters returns a list of clusters which are available
func (g *gkeClient) GetClusters() ([]*container.Cluster, error) {
	var list []*container.Cluster

	path := fmt.Sprintf("projects/%s/locations/%s", g.credentials.Spec.Project, g.credentials.Spec.Region)

	err := wait.ExponentialBackoff(retry.DefaultRetry, func() (done bool, err error) {
		resp, err := g.ce.Projects.Locations.Clusters.List(path).Do()
		if err != nil {
			log.Error(err, "failed to retrieve clusters")

			switch err := err.(type) {
			case *googleapi.Error:
				if err.Code == http.StatusForbidden {
					return false, err
				}
			default:
				return false, nil
			}
		}

		list = resp.Clusters

		return true, nil
	})

	return list, err
}

// CreateCluster is responsible for posting the cluster configuration
func (g *gkeClient) CreateCluster(ctx context.Context, request *container.CreateClusterRequest) (*container.Operation, error) {
	var operation *container.Operation

	path := fmt.Sprintf("projects/%s/locations/%s", g.credentials.Spec.Project, g.credentials.Spec.Region)

	if err := wait.ExponentialBackoff(retry.DefaultRetry, func() (done bool, err error) {
		resp, err := g.ce.Projects.Locations.Clusters.Create(path, request).Do()
		if err != nil {
			switch err := err.(type) {
			case *googleapi.Error:
				if err.Code == http.StatusBadRequest {
					return false, err
				}
			default:
				return false, nil
			}

			return false, nil
		}
		operation = resp

		return true, nil
	}); err != nil {
		return nil, err
	}

	return operation, nil
}

// EnableCloudNAT is responsible for enabling the cloud nat device
func (g *gkeClient) EnableCloudNAT() error {
	name := "router"

	if _, found, err := g.GetRouter(name); err != nil {
		return err
	} else if !found {
		return g.EnableRouter(name, g.cluster.Spec.Network)
	}

	return nil
}

// EnableRouter is responsible for create the default router in the account
func (g *gkeClient) EnableRouter(name, network string) error {
	// @step: retrieve the network
	net, err := g.GetNetwork(network)
	if err != nil {
		return err
	}

	_, err = g.cm.Routers.Insert(
		g.credentials.Spec.Project,
		g.credentials.Spec.Region,
		&compute.Router{
			Name:        name,
			Description: "Default router created by Appvia Kore",
			Network:     net.SelfLink,
			Nats: []*compute.RouterNat{
				{
					LogConfig:                     &compute.RouterNatLogConfig{Enable: false, Filter: "ALL"},
					Name:                          "cloud-nat",
					NatIpAllocateOption:           "AUTO_ONLY",
					SourceSubnetworkIpRangesToNat: "ALL_SUBNETWORKS_ALL_IP_RANGES",
				},
			},
		},
	).Do()

	return err
}

// EnableFirewallAPIServices is responsible for creating the firewall rules
func (g *gkeClient) EnableFirewallAPIServices() error {
	if err := g.AddFirewallRule(
		fmt.Sprintf("allow-%s-masters", g.cluster.Name),
		fmt.Sprintf("Allow APIExtensions for cluster: %s", g.cluster.Name),
		g.cluster.Spec.Network,
		g.cluster.Spec.MasterIPV4Cidr,
		g.cluster.Name,
		[]string{"tcp:443,5443,8443"}); err != nil {

		return err
	}

	return nil
}

// GetNetwork returns the network
func (g *gkeClient) GetNetwork(name string) (*compute.Network, error) {
	return g.cm.Networks.Get(g.credentials.Spec.Project, name).Do()
}

// GetRouter returns a specific router if it exists
func (g *gkeClient) GetRouter(name string) (*compute.Router, bool, error) {
	list, err := g.GetRouters()
	if err != nil {
		return nil, false, err
	}
	for _, x := range list {
		if x.Name == name {
			return x, true, nil
		}
	}

	return nil, false, nil
}

// GetRouters returns all the routers in the account
func (g *gkeClient) GetRouters() ([]*compute.Router, error) {
	resp, err := g.cm.Routers.List(g.credentials.Spec.Project, g.credentials.Spec.Region).Do()
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

// Locations returns a list of compute locations of for particular region
func (g *gkeClient) Locations() ([]string, error) {
	resp, err := g.ce.Projects.Locations.List(fmt.Sprintf("projects/%s", g.credentials.Spec.Project)).Do()
	if err != nil {
		return []string{}, err
	}
	var list []string

	prefix := fmt.Sprintf("%s-", g.credentials.Spec.Region)

	for _, x := range resp.Locations {
		if strings.HasPrefix(x.Name, prefix) {
			list = append(list, x.Name)
		}
	}

	return list, nil
}

// NetworkExists checks if the network exist
func (g *gkeClient) NetworkExists(name string) (bool, error) {
	return true, nil
}

// AddFirewallRule is responsible for adding a firewall rule
func (g *gkeClient) AddFirewallRule(name, description, network, source, target string, ports []string) error {
	// @step: we need to get the self-link reference to the network
	n, err := g.GetNetwork(network)
	if err != nil {
		return err
	}

	rule := &compute.Firewall{
		Name:          name,
		Allowed:       make([]*compute.FirewallAllowed, 0),
		Description:   description,
		Direction:     "INGRESS",
		EnableLogging: false,
		Network:       n.SelfLink,
		SourceRanges:  []string{source},
		TargetTags:    []string{target},
	}
	for _, x := range ports {
		rule.Allowed = append(rule.Allowed, &compute.FirewallAllowed{
			IPProtocol: strings.Split(x, ":")[0],
			Ports:      strings.Split(strings.Split(x, ":")[1], ","),
		})
	}
	// @step: check if the rule name already exists
	resp, err := g.cm.Firewalls.List(g.credentials.Spec.Project).Do()
	if err != nil {
		return err
	}

	var found bool
	for _, x := range resp.Items {
		if x.Name == name {
			found = true
			break
		}
	}

	// @step: attempt to apply the firewall rule
	err = func() error {
		if found {
			_, err := g.cm.Firewalls.Update(g.credentials.Spec.Project, name, rule).Do()
			return err
		}
		_, err := g.cm.Firewalls.Insert(g.credentials.Spec.Project, rule).Do()

		return err
	}()
	if err != nil {
		return err
	}

	return err
}

// GetOperation is responsible for retrieving the operation and status
func (g *gkeClient) GetOperation(id string) (*container.Operation, error) {
	logger := log.WithFields(log.Fields{
		"cluster":   g.cluster.Name,
		"operation": id,
	})
	logger.Debug("retrieving the status of the operation")

	// projects/my-project/locations/my-location/operations/my-operation
	path := fmt.Sprintf("projects/%s/locations/%s/operations/%s",
		g.credentials.Spec.Project,
		g.credentials.Spec.Region,
		id)

	var o *container.Operation

	// @step: retrieve the operation
	_ = wait.ExponentialBackoff(retry.DefaultRetry, func() (done bool, err error) {

		resp, err := g.ce.Projects.Locations.Operations.Get(path).Do()
		if err != nil {
			logger.WithError(err).Error("retrieving operation status")

			return false, nil
		}
		logger.Debug("retrieved operation status")

		o = resp

		return true, nil
	})

	return o, nil
}

// FindOperation is responsible for checking for a running operation
func (g *gkeClient) FindOperation(ctx context.Context, operationType, resource, name string) (string, bool, error) {
	logger := log.WithFields(log.Fields{
		"resource": resource,
		"type":     operationType,
		"name":     name,
	})
	logger.Debug("searching for any running operations")

	resp, err := g.ce.Projects.Locations.Operations.List(fmt.Sprintf("projects/%s/locations/%s",
		g.credentials.Spec.Project, g.credentials.Spec.Region)).Do()
	if err != nil {
		logger.WithError(err).Error("trying to retrieve a list of operations")

		return "", false, err
	}
	for _, x := range resp.Operations {
		if x.OperationType == operationType {
			if strings.HasSuffix(x.TargetLink, fmt.Sprintf("%s/%s", resource, name)) {
				if x.Status == "RUNNING" {
					return x.Name, true, nil
				}
			}
		}
	}

	return "", false, nil
}

// WaitOnOperation is responsible for waiting on a operation to complete fail
func (g *gkeClient) WaitOnOperation(ctx context.Context, id string, interval, timeout time.Duration) error {
	logger := log.WithFields(log.Fields{
		"cluster":   g.cluster.Name,
		"operation": id,
	})
	logger.Info("checking the status of operation")

	return wait.PollImmediate(interval, timeout, func() (bool, error) {
		status, err := g.GetOperation(id)
		if err != nil {
			logger.WithError(err).Error("encountered error waiting on operation")

			return false, nil
		}

		if status == nil {
			return false, nil
		}
		if status.Status == "DONE" {
			return true, nil
		}

		logger.WithField("status", status.Status).Debug("waiting for operation to finish or fail")

		return false, nil
	})
}

// Exists checks if the cluster exists
func (g *gkeClient) Exists() (bool, error) {
	log.WithField("name", g.cluster.Name).Debug("checking for gke cluster existence")

	_, found, err := g.GetCluster()

	return found, err
}
