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

package gke

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/utils"

	"github.com/hashicorp/go-version"
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
	credentials *credentials
	// cluster is the gke cluster
	cluster *gke.GKE
	// @deprecated region
	region string
}

// NodePoolOperation represents the possible operations which may be in process on a node pool.
type NodePoolOperation string

// NodePoolOperationCreating means a node pool is being created
const NodePoolOperationCreating NodePoolOperation = "Creating"

// NodePoolOperationUpdating means features of a node pool (size, auto-scale, image type, version) are being updated.
const NodePoolOperationUpdating NodePoolOperation = "Updating"

// NodePoolOperationDeleting means a node pool is being deleted
const NodePoolOperationDeleting NodePoolOperation = "Deleting"

// NodePoolOperationNone means no operation is being performed on node pools.
const NodePoolOperationNone NodePoolOperation = "None"

// NewClient returns a gcp client for us
func NewClient(credentials *credentials, cluster *gke.GKE) (*gkeClient, error) {
	options := option.WithCredentialsJSON([]byte(credentials.key))

	cm, err := compute.NewService(context.Background(), options)
	if err != nil {
		return nil, err
	}

	ce, err := container.NewService(context.Background(), options)
	if err != nil {
		return nil, err
	}

	region := credentials.region
	if region == "" {
		region = cluster.Spec.Region
	}

	// @step: Ensure any old specs without node pools are upgraded to the new nodepool spec
	if len(cluster.Spec.NodePools) == 0 {
		// Create pool spec based on the deprecated fields and previous hard-coded
		// defaults.
		cluster.Spec.NodePools = []gke.GKENodePool{
			{
				Name:              "compute",
				Version:           cluster.Spec.Version,
				EnableAutoupgrade: cluster.Spec.EnableAutoupgrade,
				EnableAutoscaler:  cluster.Spec.EnableAutoscaler,
				EnableAutorepair:  cluster.Spec.EnableAutorepair,
				MinSize:           1,
				MaxSize:           cluster.Spec.MaxSize,
				Size:              cluster.Spec.Size,
				DiskSize:          cluster.Spec.DiskSize,
				ImageType:         cluster.Spec.ImageType,
			},
		}
	}

	return &gkeClient{
		cm:          cm,
		ce:          ce,
		credentials: credentials,
		cluster:     cluster,
		region:      region,
	}, nil
}

// Delete attempts to delete the cluster from gke
func (g *gkeClient) Delete(ctx context.Context) error {
	logger := log.WithFields(log.Fields{
		"cluster":   g.cluster.Name,
		"namespace": g.cluster.Namespace,
		"project":   g.credentials.project,
		"region":    g.region,
	})
	logger.Info("attempting to delete the cluster from gcp")

	found, err := g.Exists(ctx)
	if err != nil {
		logger.WithError(err).Error("trying to check for the cluster")

		return err
	}
	if !found {
		return nil
	}

	cluster, _, err := g.GetCluster(ctx)
	if err != nil {
		logger.WithError(err).Error("trying to retrieve the cluster")

		return err
	}
	path := fmt.Sprintf("projects/%s/locations/%s/clusters/%s",
		g.credentials.project,
		g.region,
		cluster.Name)

	// @step: check for any ongoing operation
	_, found, err = g.FindOperation(ctx, "DELETE_CLUSTER", "kubernetes", cluster.Name)
	if err != nil {
		logger.WithError(err).Error("trying to check for current operations")

		return err
	}
	if !found {
		if _, err := g.ce.Projects.Locations.Clusters.Delete(path).Do(); err != nil {
			logger.WithError(err).Error("trying to delete the cluster")

			return err
		}
		logger.Debug("requested the removal of the gke cluster")
	}

	return nil
}

// Create is used to create the cluster in gcp
func (g *gkeClient) Create(ctx context.Context) (*container.Cluster, error) {
	logger := log.WithFields(log.Fields{
		"cluster":   g.cluster.Name,
		"namespace": g.cluster.Namespace,
		"project":   g.credentials.project,
		"region":    g.region,
	})
	logger.Info("attempting to create the gke cluster")

	// @step: we create the definitions
	def, err := g.CreateDefinition()
	if err != nil {
		logger.WithError(err).Error("attempting to create the cluster definition")

		return nil, err
	}

	// @step: looking for any ongoing operation
	_, found, err := g.FindOperation(ctx, "CREATE_CLUSTER", "kubernetes", g.cluster.Name)
	if err != nil {
		return nil, err
	}
	if !found {
		// @step: we request the cluster
		if _, err := g.CreateCluster(ctx, def); err != nil {
			if err != nil {
				logger.WithError(err).Error("attempting to request the cluster")

				return nil, err
			}
		}
	}

	// @step: retrieve the state of the cluster via api
	gc, _, err := g.GetCluster(ctx)
	if err != nil {
		logger.WithError(err).Error("retrieving gke cluster details")

		return nil, err
	}

	return gc, nil
}

// Update is called to update the cluster
func (g *gkeClient) Update(ctx context.Context) (bool, error) {
	logger := log.WithFields(log.Fields{
		"name":      g.cluster.Name,
		"namespace": g.cluster.Namespace,
	})
	logger.Info("checking if the cluster requires updating")

	// @step: get the current state of the cluster
	state, found, err := g.GetCluster(ctx)
	if !found {
		return false, errors.New("cluster was not found")
	}
	if err != nil {
		return false, err
	}

	update, err := g.CreateUpdateDefinition(state)
	if err != nil {
		log.WithError(err).Error("creating the update request")

		return false, err
	}

	// @step: we check if the update request has been altered
	if utils.IsEmpty(update.Update) {
		return false, nil
	}
	logger.Debug("desired state of the cluster has drifted, attempting to update")

	path := fmt.Sprintf("projects/%s/locations/%s/clusters/%s",
		g.credentials.project,
		g.region,
		g.cluster.Name)

	_, err = g.ce.Projects.Locations.Clusters.Update(path, update).Context(ctx).Do()
	if err != nil {
		logger.WithError(err).Error("trying to update the cluster")

		return false, err
	}
	logger.Debug("successfully requested the gke cluster to update")

	return true, nil
}

// UpdateNodePools is called to add, remove and update node pools in the cluster. Returns the operation
// being performed and the name of the node pool being operated on, or NodePoolOperationNone if no
// operation being performed.
func (g *gkeClient) UpdateNodePools(ctx context.Context) (NodePoolOperation, string, error) {
	logger := log.WithFields(log.Fields{
		"name":      g.cluster.Name,
		"namespace": g.cluster.Namespace,
	})
	logger.Info("checking if the cluster node pools require updating")

	// @step: get the current state of the cluster
	state, found, err := g.GetCluster(ctx)
	if !found {
		return NodePoolOperationNone, "", errors.New("cluster was not found")
	}
	if err != nil {
		return NodePoolOperationNone, "", err
	}

	clusterPath := fmt.Sprintf("projects/%s/locations/%s/clusters/%s",
		g.credentials.project,
		g.region,
		g.cluster.Name)

	// Process node pools - additions
	for _, nodePoolSpec := range g.cluster.Spec.NodePools {
		nodePoolExists := false
		for _, n := range state.NodePools {
			if n.Name == nodePoolSpec.Name {
				nodePoolExists = true
			}
		}

		if !nodePoolExists {
			err := g.createNodePool(ctx, clusterPath, &nodePoolSpec, logger)
			if err != nil {
				return NodePoolOperationNone, "", err
			}
			return NodePoolOperationCreating, nodePoolSpec.Name, nil
		}
	}

	// Process node pools - removals and updates
	for _, nodePool := range state.NodePools {
		var nodePoolSpec *gke.GKENodePool = nil
		for _, n := range g.cluster.Spec.NodePools {
			if n.Name == nodePool.Name {
				nodePoolSpec = &n
			}
		}

		nodePoolLogger := logger.WithField("nodePool", nodePool.Name)
		nodePoolPath := fmt.Sprintf("%s/nodePools/%s", clusterPath, nodePool.Name)

		// Node pool removed from spec, request removal from GKE
		if nodePoolSpec == nil {
			err := g.deleteNodePool(ctx, nodePoolPath, nodePoolLogger)
			if err != nil {
				return NodePoolOperationNone, "", err
			}
			return NodePoolOperationDeleting, nodePool.Name, nil
		}

		// Node pool still in spec, check for updates:
		updating, err := g.updateNodePool(ctx, nodePoolPath, nodePoolSpec, nodePool, nodePoolLogger)
		if err != nil {
			return NodePoolOperationNone, "", err
		}
		if updating {
			return NodePoolOperationUpdating, nodePoolSpec.Name, nil
		}
	}

	return NodePoolOperationNone, "", nil
}

func (g *gkeClient) createNodePool(ctx context.Context, clusterPath string, nodePoolSpec *gke.GKENodePool, logger *log.Entry) error {
	// Node pool added to spec, request addition to GKE (or check if we already have)
	nodePoolLogger := logger.WithField("nodePool", nodePoolSpec.Name)

	req, err := g.CreateNodePoolDefinition(nodePoolSpec)
	if err != nil {
		nodePoolLogger.WithError(err).Error("trying to prepare create node pool definition")

		return err
	}

	// @step: Check if already in operation
	_, found, err := g.FindOperation(ctx, "CREATE_NODE_POOL", "kubernetes", g.cluster.Name)
	if err != nil {
		return err
	}
	if found {
		nodePoolLogger.Debug("node pool creation still in progress")
		return nil
	}
	nodePoolLogger.Info("Node pool added to spec, requesting addition to cluster")
	// @step: request the node pool
	_, err = g.ce.Projects.Locations.Clusters.NodePools.Create(clusterPath, req).Context(ctx).Do()
	if err != nil {
		nodePoolLogger.WithError(err).Error("trying to create node pool")

		return err
	}
	nodePoolLogger.Debug("successfully requested the gke node pool to create")
	return nil
}

func (g *gkeClient) deleteNodePool(ctx context.Context, nodePoolPath string, nodePoolLogger *log.Entry) error {
	// @step: Check if already in operation
	_, found, err := g.FindOperation(ctx, "DELETE_NODE_POOL", "kubernetes", g.cluster.Name)
	if err != nil {
		return err
	}
	if found {
		nodePoolLogger.Debug("node pool deletion already in progress, not re-requesting")
		return nil
	}
	// @step: delete the node pool
	nodePoolLogger.Info("Node pool removed from spec, requesting removal from cluster")
	_, err = g.ce.Projects.Locations.Clusters.NodePools.Delete(nodePoolPath).Context(ctx).Do()
	if err != nil {
		nodePoolLogger.WithError(err).Error("trying to delete node pool")

		return err
	}
	nodePoolLogger.Debug("Successfully requested the node pool to delete")
	return nil
}

func (g *gkeClient) updateNodePool(ctx context.Context, nodePoolPath string, nodePoolSpec *gke.GKENodePool, nodePool *container.NodePool, nodePoolLogger *log.Entry) (bool, error) {
	updatePool := false
	var err error
	if g.cluster.Spec.ReleaseChannel == "" {
		// Check versions if not following a release channel (the release channel will take care of the node
		// pools if it is specified)

		targetVersion := nodePoolSpec.Version
		if targetVersion == "" {
			// A blank target version means keep in sync with cluster, so use the cluster spec version.
			targetVersion = g.cluster.Spec.Version
		}

		updatePool, err = UpgradeRequired(nodePool.Version, targetVersion)
		if err != nil {
			return false, err
		}
	}

	if !updatePool && nodePoolSpec.ImageType != nodePool.Config.ImageType {
		updatePool = true
	}

	if updatePool {
		err := g.requestNodePoolUpgrade(ctx, nodePoolPath, nodePoolSpec, nodePoolLogger)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	// Check the auto-scale config is in sync.
	if nodePoolSpec.EnableAutoscaler != nodePool.Autoscaling.Enabled || (nodePoolSpec.EnableAutoscaler && (nodePoolSpec.MinSize != nodePool.Autoscaling.MinNodeCount || nodePoolSpec.MaxSize != nodePool.Autoscaling.MaxNodeCount)) {
		err := g.setNodePoolAutoscale(ctx, nodePoolPath, nodePoolSpec, nodePool, nodePoolLogger)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	// If auto-scale is disabled, check the size is as desired.
	if !nodePoolSpec.EnableAutoscaler && !nodePool.Autoscaling.Enabled && nodePoolSpec.Size != nodePool.InitialNodeCount {
		err := g.setNodePoolSize(ctx, nodePoolPath, nodePoolSpec, nodePool, nodePoolLogger)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func (g *gkeClient) requestNodePoolUpgrade(ctx context.Context, nodePoolPath string, nodePoolSpec *gke.GKENodePool, nodePoolLogger *log.Entry) error {
	// @step: Check if already in operation
	_, found, err := g.FindOperation(ctx, "UPGRADE_NODES", "kubernetes", g.cluster.Name)
	if err != nil {
		return err
	}
	if found {
		nodePoolLogger.Debug("node pool upgrade still in progress")
		return nil
	}

	nodePoolLogger.Info("Node pool upgrade required to match spec - image type changed or version updated")

	// @step: request upgrade / image type change
	updateReq := &container.UpdateNodePoolRequest{
		NodeVersion: nodePoolSpec.Version,
		ImageType:   nodePoolSpec.ImageType,
	}
	_, err = g.ce.Projects.Locations.Clusters.NodePools.Update(nodePoolPath, updateReq).Context(ctx).Do()
	if err != nil {
		nodePoolLogger.WithError(err).Error("trying to update node pool")

		return err
	}
	nodePoolLogger.Debug("Successfully requested the node pool update")
	return nil
}

func (g *gkeClient) setNodePoolAutoscale(ctx context.Context, nodePoolPath string, nodePoolSpec *gke.GKENodePool, nodePool *container.NodePool, nodePoolLogger *log.Entry) error {
	autoScale := func() *container.NodePoolAutoscaling {
		if nodePoolSpec.EnableAutoscaler {
			return &container.NodePoolAutoscaling{
				Autoprovisioned: false,
				Enabled:         true,
				MaxNodeCount:    nodePoolSpec.MaxSize,
				MinNodeCount:    nodePoolSpec.MinSize,
			}
		}
		return &container.NodePoolAutoscaling{
			Autoprovisioned: false,
			Enabled:         false,
		}
	}()

	nodePoolLogger.WithFields(log.Fields{
		"current": nodePool.Autoscaling,
		"spec":    autoScale,
	}).Info("Node pool auto-scale configuration update required")

	_, err := g.ce.Projects.Locations.Clusters.NodePools.SetAutoscaling(nodePoolPath, &container.SetNodePoolAutoscalingRequest{Autoscaling: autoScale}).Context(ctx).Do()
	if err != nil {
		nodePoolLogger.WithError(err).Error("trying to update node pool auto-scale configuration")

		return err
	}
	nodePoolLogger.Debug("Successfully requested the node pool auto-scale configuration change")
	return nil
}

func (g *gkeClient) setNodePoolSize(ctx context.Context, nodePoolPath string, nodePoolSpec *gke.GKENodePool, nodePool *container.NodePool, nodePoolLogger *log.Entry) error {
	// @step: Check if already in operation
	_, found, err := g.FindOperation(ctx, "SET_NODE_POOL_SIZE", "kubernetes", g.cluster.Name)
	if err != nil {
		return err
	}
	if found {
		nodePoolLogger.Debug("node pool re-size still in progress")
		return nil
	}

	nodePoolLogger.WithFields(log.Fields{
		"currentSize": nodePool.InitialNodeCount,
		"specSize":    nodePoolSpec.Size,
	}).Info("Node pool size (non-auto-scale) update required")

	_, err = g.ce.Projects.Locations.Clusters.NodePools.SetSize(nodePoolPath, &container.SetNodePoolSizeRequest{NodeCount: nodePoolSpec.Size}).Context(ctx).Do()
	if err != nil {
		nodePoolLogger.WithError(err).Error("trying to update node pool size")

		return err
	}
	nodePoolLogger.Debug("Successfully requested the node pool size change")
	return nil
}

// CreateUpdateDefinition returns a cluster update definition
// @notes: so GKE will only handle one update at a time, so if the user makes a bunch of changes to the
// spec we need to return the first change, update, requeue and do the next i guess.
func (g *gkeClient) CreateUpdateDefinition(state *container.Cluster) (*container.UpdateClusterRequest, error) {
	logger := log.WithFields(log.Fields{
		"name":      g.cluster.Name,
		"namespace": g.cluster.Namespace,
	})

	request := &container.UpdateClusterRequest{
		ProjectId: g.credentials.project,
		Update:    &container.ClusterUpdate{},
	}

	u := request.Update

	// Update release channel if changed.
	if state.ReleaseChannel.Channel != g.cluster.Spec.ReleaseChannel {
		logger.WithFields(log.Fields{
			"currentChannel": state.ReleaseChannel.Channel,
			"specChannel":    g.cluster.Spec.ReleaseChannel,
		}).Info("Release channel changed")

		u.DesiredReleaseChannel = &container.ReleaseChannel{
			Channel: g.cluster.Spec.ReleaseChannel,
		}

		return request, nil
	}

	// Manual version control only possible if not following a release channel
	if g.cluster.Spec.ReleaseChannel == "" {
		// Notes: Master version is *always* auto-upgraded by google, so whatever version is specified by the
		// spec, we must only apply an upgrade if the spec requests a version AHEAD of the current master version,
		// else it will 'flap' back and forth as we request downgrades then GCP re-upgrades the master.
		upgrade, err := UpgradeRequired(state.CurrentMasterVersion, g.cluster.Spec.Version)
		if err != nil {
			return nil, err
		}

		if upgrade {
			logger.WithFields(log.Fields{
				"currVersion":    state.CurrentMasterVersion,
				"desiredVersion": g.cluster.Spec.Version,
			}).Debug("Master upgrade required")

			u.DesiredMasterVersion = g.cluster.Spec.Version

			return request, nil
		}
	}

	return request, nil
}

// UpgradeRequired compares an actual GKE version (e.g. 1.15.1-gke.9) with a desired version
// (e.g. 1.15, 1.15.1, 1.15.1-gke.9) and returns true if the desired represents a greater version
// than the current. If the desired version is blank or one of the 'magic' values specified by
// GKE (- or latest), it will always return false as these will only be used to set the initial
// versions.
func UpgradeRequired(current string, desired string) (bool, error) {
	if desired == "" || desired == "-" || desired == "latest" {
		return false, nil
	}

	desiredV, err := version.NewVersion(desired)
	if err != nil {
		return false, err
	}
	var currentV *version.Version
	if !strings.Contains(desired, "-") && strings.Contains(current, "-") {
		// Strip the GKE section from the current version as standard
		// version semantics would mean 1.15.1-gke.9 is BEHIND 1.15.1
		// whereas in GKE, the -suffix is used to denote a GKE version
		// NOT a pre-release version.
		currentV, err = version.NewVersion(current[0:strings.Index(current, "-")])
		if err != nil {
			return false, err
		}
	} else {
		currentV, err = version.NewVersion(current)
		if err != nil {
			return false, err
		}
	}

	return desiredV.GreaterThan(currentV), nil
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

		ReleaseChannel: &container.ReleaseChannel{
			Channel: func() string {
				if cluster.Spec.ReleaseChannel == "" {
					return "UNSPECIFIED"
				}
				return cluster.Spec.ReleaseChannel
			}(),
		},

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
				Disabled: !cluster.Spec.EnableHorizontalPodAutoscaler,
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
		Network:                 cluster.Spec.Network,
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
			ServicesIpv4Cidr:     cluster.Spec.ServicesIPV4Cidr,
			CreateSubnetwork:     false,
			SubnetworkName:       "default",
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
	}

	for _, nodePool := range cluster.Spec.NodePools {
		resource.NodePools = append(resource.NodePools, g.PrepareNodePoolDefinition(&nodePool, locations))
	}

	resource.PrivateClusterConfig = &container.PrivateClusterConfig{}

	if cluster.Spec.EnablePrivateNetwork {
		resource.PrivateClusterConfig.EnablePrivateNodes = true
	}

	if cluster.Spec.EnablePrivateEndpoint {
		resource.PrivateClusterConfig.EnablePrivateEndpoint = true
		resource.PrivateClusterConfig.MasterIpv4CidrBlock = cluster.Spec.MasterIPV4Cidr
	}

	if len(cluster.Spec.AuthorizedMasterNetworks) > 0 {
		var cidrBlocks []*container.CidrBlock
		for _, an := range cluster.Spec.AuthorizedMasterNetworks {
			cidrBlocks = append(cidrBlocks, &container.CidrBlock{
				CidrBlock:   an.CIDR,
				DisplayName: an.Name,
			})
		}
		resource.MasterAuthorizedNetworksConfig = &container.MasterAuthorizedNetworksConfig{
			CidrBlocks: cidrBlocks,
			Enabled:    true,
		}
	} else {
		resource.MasterAuthorizedNetworksConfig = &container.MasterAuthorizedNetworksConfig{
			Enabled: false,
		}
	}

	return &container.CreateClusterRequest{
		ProjectId: g.credentials.project,
		Cluster:   resource,
	}, nil
}

// CreateNodePoolDefinition returns a node pool definition
func (g *gkeClient) CreateNodePoolDefinition(nodePool *gke.GKENodePool) (*container.CreateNodePoolRequest, error) {
	// @step: retrieve a list of location to place this node pool
	locations, err := g.Locations()
	if err != nil {
		return nil, err
	}

	resource := g.PrepareNodePoolDefinition(nodePool, locations)
	parentPath := fmt.Sprintf("projects/%s/locations/%s/clusters/%s",
		g.credentials.project,
		g.region,
		g.cluster.Name)

	return &container.CreateNodePoolRequest{
		Parent:   parentPath,
		NodePool: resource,
	}, nil
}

// PrepareNodePoolDefinition translates our node pool spec into a GCP API node pool struct.
func (g *gkeClient) PrepareNodePoolDefinition(nodePool *gke.GKENodePool, locations []string) *container.NodePool {
	cluster := g.cluster

	autoScale := func() *container.NodePoolAutoscaling {
		if nodePool.EnableAutoscaler {
			return &container.NodePoolAutoscaling{
				Autoprovisioned: false,
				Enabled:         true,
				MaxNodeCount:    nodePool.MaxSize,
				MinNodeCount:    nodePool.MinSize,
			}
		}
		return &container.NodePoolAutoscaling{
			Autoprovisioned: false,
			Enabled:         false,
		}
	}()

	return &container.NodePool{
		Name:        nodePool.Name,
		Autoscaling: autoScale,
		Config: &container.NodeConfig{
			DiskSizeGb:  nodePool.DiskSize,
			ImageType:   nodePool.ImageType,
			MachineType: nodePool.MachineType,
			OauthScopes: []string{
				"https://www.googleapis.com/auth/compute",
				"https://www.googleapis.com/auth/devstorage.read_only",
				"https://www.googleapis.com/auth/logging.write",
				"https://www.googleapis.com/auth/monitoring",
			},
			Preemptible: nodePool.Preemptible,
			Tags:        []string{cluster.Name},
		},
		InitialNodeCount: nodePool.Size,
		Locations:        locations,
		Management: &container.NodeManagement{
			AutoRepair: nodePool.EnableAutorepair,
			AutoUpgrade: func() bool {
				// If a release channel is set, auto upgrade MUST be true.
				if cluster.Spec.ReleaseChannel != "" {
					return true
				}
				return nodePool.EnableAutoupgrade
			}(),
		},
		MaxPodsConstraint: &container.MaxPodsConstraint{
			MaxPodsPerNode: nodePool.MaxPodsPerNode,
		},
		Version: func() string {
			// If a release channel is set, the version MUST be empty.
			if cluster.Spec.ReleaseChannel != "" {
				return ""
			}
			// If blank and not following release channel, use master version
			if nodePool.Version == "" {
				return cluster.Spec.Version
			}
			return nodePool.Version
		}(),
	}
}

// GetCluster returns a cluster config
func (g *gkeClient) GetCluster(ctx context.Context) (*container.Cluster, bool, error) {
	clusters, err := g.GetClusters(ctx)
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
func (g *gkeClient) GetClusters(ctx context.Context) ([]*container.Cluster, error) {
	var list []*container.Cluster

	path := fmt.Sprintf("projects/%s/locations/%s", g.credentials.project, g.region)

	err := wait.ExponentialBackoff(retry.DefaultRetry, func() (done bool, err error) {
		resp, err := g.ce.Projects.Locations.Clusters.List(path).Context(ctx).Do()
		if err != nil {
			log.Error(err, "failed to retrieve clusters")

			switch err := err.(type) {
			case *googleapi.Error:
				if err.Code == http.StatusForbidden {
					// we definitely need to quit here - no point in retrying
					return false, err
				}

				// @step: in absence of knowing the error, we will retry and use
				// the backoff and retry to handle this
				return false, nil
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

	path := fmt.Sprintf("projects/%s/locations/%s", g.credentials.project, g.region)

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
		g.credentials.project,
		g.region,
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
	return g.cm.Networks.Get(g.credentials.project, name).Do()
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
	resp, err := g.cm.Routers.List(g.credentials.project, g.region).Do()
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

// Locations returns a list of compute locations of for particular region
func (g *gkeClient) Locations() ([]string, error) {
	resp, err := g.ce.Projects.Locations.List(fmt.Sprintf("projects/%s", g.credentials.project)).Do()
	if err != nil {
		return []string{}, err
	}
	var list []string

	prefix := fmt.Sprintf("%s-", g.region)

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
	resp, err := g.cm.Firewalls.List(g.credentials.project).Do()
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
			_, err := g.cm.Firewalls.Update(g.credentials.project, name, rule).Do()
			return err
		}
		_, err := g.cm.Firewalls.Insert(g.credentials.project, rule).Do()

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
		g.credentials.project,
		g.region,
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
		g.credentials.project, g.region)).Do()
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
func (g *gkeClient) Exists(ctx context.Context) (bool, error) {
	log.WithField("name", g.cluster.Name).Debug("checking for gke cluster existence")

	_, found, err := g.GetCluster(ctx)

	return found, err
}
