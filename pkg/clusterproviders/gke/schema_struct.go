// Code generated by struct-gen; DO NOT EDIT.

package gke

// The networks which are allowed to connect to this cluster (e.g. via kubectl).
type AuthProxyAllowedIP string

// The networks which are allowed to access the master control plane.
type AuthorizedMasterNetwork struct {
	Cidr string `json:"cidr"`
	Name string `json:"name"`
}

// Users who should be allowed to access this cluster, will override any default role set above for these users
type ClusterUser struct {
	Roles    []Role `json:"roles"`
	Username string `json:"username"`
}

// GKE Cluster Plan Schema
type Configuration struct {

	// AuthProxyAllowedIPRanges The networks which are allowed to connect to this cluster (e.g. via kubectl).
	AuthProxyAllowedIPRanges []AuthProxyAllowedIP `json:"authProxyAllowedIPs"`
	// AuthorizedMasterNetworks The networks which are allowed to access the master control plane.
	AuthorizedMasterNetworks []AuthorizedMasterNetwork `json:"authorizedMasterNetworks"`
	// ClusterUsers Users who should be allowed to access this cluster, will override any default role set above for these users
	ClusterUsers []ClusterUser `json:"clusterUsers,omitempty"`
	// DefaultTeamRole The role that team members will have on this cluster if 'inherit team members' enabled
	DefaultTeamRole string `json:"defaultTeamRole,omitempty"`
	// Description Meaningful description of this cluster.
	Description string `json:"description"`
	// DiskSize DEPRECATED: Set disk size on node pool instead
	DiskSize int64 `json:"diskSize,omitempty"`
	// Domain The domain for this cluster.
	Domain string `json:"domain"`
	// EnableAutorepair DEPRECATED: Set auto-repair on node pool instead
	EnableAutorepair bool `json:"enableAutorepair,omitempty"`
	// EnableAutoscaler DEPRECATED: Set auto-scale on node pool instead
	EnableAutoscaler bool `json:"enableAutoscaler,omitempty"`
	// EnableAutoupgrade DEPRECATED: Set auto-upgrade on node pool instead
	EnableAutoupgrade             bool `json:"enableAutoupgrade,omitempty"`
	EnableDefaultTrafficBlock     bool `json:"enableDefaultTrafficBlock"`
	EnableHorizontalPodAutoscaler bool `json:"enableHorizontalPodAutoscaler"`
	EnableHttploadBalancer        bool `json:"enableHTTPLoadBalancer"`
	EnableIstio                   bool `json:"enableIstio"`
	EnablePrivateEndpoint         bool `json:"enablePrivateEndpoint"`
	EnablePrivateNetwork          bool `json:"enablePrivateNetwork"`
	// EnableShieldedNodes Shielded nodes provide additional verifications of the node OS and VM, with enhanced rootkit and bootkit protection applied
	EnableShieldedNodes      bool `json:"enableShieldedNodes"`
	EnableStackDriverLogging bool `json:"enableStackDriverLogging"`
	EnableStackDriverMetrics bool `json:"enableStackDriverMetrics"`
	// ImageType DEPRECATED: Set image type on node pool instead
	ImageType string `json:"imageType,omitempty"`
	// InheritTeamMembers Whether team members will all have access to this cluster by default
	InheritTeamMembers bool `json:"inheritTeamMembers"`
	// MachineType DEPRECATED: Set machine type on node pool instead
	MachineType string `json:"machineType,omitempty"`
	// MaintenanceWindow Time of day to allow maintenance operations to be performed by the cloud provider on this cluster.
	MaintenanceWindow string `json:"maintenanceWindow"`
	// MaxSize DEPRECATED: Set max size on node pool instead
	MaxSize int64 `json:"maxSize,omitempty"`
	// Network DEPRECATED: It is not supported to specify a custom network. This property will be ignored.
	Network   string     `json:"network,omitempty"`
	NodePools []NodePool `json:"nodePools,omitempty"`
	// Region Geographical location for this cluster
	Region string `json:"region"`
	// ReleaseChannel Follow a GKE release channel to control the auto-upgrade of your cluster - if set, auto-upgrade will be true on all node groups
	ReleaseChannel string `json:"releaseChannel,omitempty"`
	// Size DEPRECATED: Set size on node pool instead
	Size int64 `json:"size,omitempty"`
	// Subnetwork DEPRECATED: Unused
	Subnetwork string `json:"subnetwork,omitempty"`
	// Version Kubernetes version - must be blank if release channel specified.
	Version string `json:"version"`
}

// A set of labels to help Kubernetes workloads find this group
type Label string

type NodePool struct {

	// DiskSize The amount of storage in GiB provisioned on the nodes in this group
	DiskSize int64 `json:"diskSize"`
	// EnableAutorepair Automatically repair any failed nodes within this node pool.
	EnableAutorepair bool `json:"enableAutorepair"`
	// EnableAutoscaler Add and remove nodes automatically based on load
	EnableAutoscaler bool `json:"enableAutoscaler"`
	// EnableAutoupgrade Enable to update this node pool updated when new GKE versions are made available by GCP - must be enabled if a release channel is selected
	EnableAutoupgrade bool `json:"enableAutoupgrade"`
	// ImageType The image type used by the nodes
	ImageType string `json:"imageType"`
	// Labels A set of labels to help Kubernetes workloads find this group
	Labels map[string]Label `json:"labels,omitempty"`
	// MachineType The type of nodes used for this node pool
	MachineType string `json:"machineType"`
	// MaxPodsPerNode The maximum number of pods that can be scheduled onto each node of this pool
	MaxPodsPerNode int64 `json:"maxPodsPerNode,omitempty"`
	// MaxSize The maximum nodes this pool should contain (if auto-scale enabled)
	MaxSize int64 `json:"maxSize,omitempty"`
	// MinSize The minimum nodes this pool should contain (if auto-scale enabled)
	MinSize int64 `json:"minSize,omitempty"`
	// Name Name of this node pool. Must be unique within the cluster.
	Name string `json:"name"`
	// Preemptible Whether to use pre-emptible nodes (cheaper, but can and will be terminated at any time, use with care).
	Preemptible bool `json:"preemptible,omitempty"`
	// Size How many nodes to build when provisioning this pool - if autoscaling enabled, this will be the initial size
	Size int64 `json:"size"`
	// Taints A collection of kubernetes taints to add on the nodes.
	Taints []Taint `json:"taints,omitempty"`
	// Version Node pool version, blank to use same version as cluster (recommended); must be blank if cluster follows a release channel. Must be within 2 minor versions of the master version (e.g. for master version 1.16, this must be 1.14, 1.15 or 1.16) or 1 minor version if auto-upgrade enabled
	Version string `json:"version"`
}

type Role string

// A collection of kubernetes taints to add on the nodes.
type Taint struct {
	// Effect The chosen effect of the taint
	Effect string `json:"effect"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}