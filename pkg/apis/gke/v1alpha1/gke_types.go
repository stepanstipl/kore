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

package v1alpha1

import (
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GKESpec defines the desired state of GKE
// +k8s:openapi-gen=true
type GKESpec struct {
	// Cluster refers to the cluster this object belongs to
	// +kubebuilder:validation:Required
	Cluster corev1.Ownership `json:"cluster,omitempty"`
	// Credentials is a reference to the gke credentials object to use
	// +kubebuilder:validation:Required
	// +k8s:openapi-gen=false
	Credentials corev1.Ownership `json:"credentials"`
	// Description provides a short summary / description of the cluster.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Description string `json:"description"`
	// Version is the kubernetes version which the cluster master should be
	// configured with. '-' gives the current GKE default version, 'latest' gives most recent,
	// 1.15 would be latest 1.15.x release, 1.15.1 would be the latest 1.15.1 release, and
	// 1.15.1-gke.1 would be the exact specified version. Must be blank if following release channel.
	// +kubebuilder:validation:Optional
	Version string `json:"version"`
	// ReleaseChannel is the GKE release channel to follow, '' (to follow no channel),
	// 'STABLE' (only battle-tested releases every few months), 'REGULAR' (stable releases
	// every few weeks) or 'RAPID' (bleeding edge, not suitable for production workloads). If anything other
	// than '', Version must be blank.
	// +kubebuilder:validation:Optional
	ReleaseChannel string `json:"releaseChannel"`
	// AuthorizedMasterNetworks is a collection of authorized networks which is
	// permitted to speak to the kubernetes API, default to all if not provided.
	// +kubebuilder:validation:Optional
	AuthorizedMasterNetworks []*AuthorizedNetwork `json:"authorizedMasterNetworks"`
	// ServicesIPV4Cidr is an optional network cidr configured for the cluster
	// services
	// +kubebuilder:validation:Optional
	ServicesIPV4Cidr string `json:"servicesIPV4Cidr"`
	// Region is the gcp region you want the cluster to reside
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Region string `json:"region,omitempty"`
	// ClusterIPV4Cidr is an optional network CIDR which is used to place the
	// pod network on
	// +kubebuilder:validation:Optional
	ClusterIPV4Cidr string `json:"clusterIPV4Cidr"`
	// EnableHorizontalPodAutoscaler indicates if the cluster is configured with
	// the horizontal pod autoscaler addon. This automatically adjusts the cpu and
	// memory resources of pods in accordance with their demand. You should ensure
	// you use PodDisruptionBudgets if this is enabled.
	// +kubebuilder:validation:Optional
	EnableHorizontalPodAutoscaler bool `json:"enableHorizontalPodAutoscaler"`
	// EnableHTTPLoadBalancer indicates if the cluster should be configured with
	// the GKE ingress controller. When enabled GKE will autodiscover your
	// ingress resources and provision load balancer on your behalf.
	// +kubebuilder:validation:Optional
	EnableHTTPLoadBalancer bool `json:"enableHTTPLoadBalancer"`
	// EnableIstio indicates if the GKE Istio service mesh is deployed to the
	// cluster; this provides a more feature rich routing and instrumentation.
	// +kubebuilder:validation:Optional
	EnableIstio bool `json:"enableIstio"`
	// EnableShieldedNodes indicates we should enable the shielded nodes options in GKE.
	// This protects against a variety of attacks by hardening the underlying GKE node
	// against rootkits and bootkits.
	EnableShieldedNodes bool `json:"enableShieldedNodes"`
	// EnableStackDriverLogging indicates if Stackdriver logging should be enabled
	// for the cluster
	// +kubebuilder:validation:Optional
	EnableStackDriverLogging bool `json:"enableStackDriverLogging"`
	// EnableStackDriverMetrics indicates if Stackdriver metrics should be enabled
	// for the cluster
	// +kubebuilder:validation:Optional
	EnableStackDriverMetrics bool `json:"enableStackDriverMetrics"`
	// EnablePrivateEndpoint indicates whether the Kubernetes API should only be accessible from internal IP addresses
	// +kubebuilder:validation:Optional
	EnablePrivateEndpoint bool `json:"enablePrivateEndpoint"`
	// EnablePrivateNetwork indicates if compute nodes should have external ip
	// addresses or use private networking and a cloud-nat device.
	// +kubebuilder:validation:Optional
	EnablePrivateNetwork bool `json:"enablePrivateNetwork"`
	// MasterIPV4Cidr is network range used when private networking is enabled.
	// This is the peering subnet used to to GKE master api layer. Note, this must
	// be unique within the network.
	// +kubebuilder:validation:Optional
	MasterIPV4Cidr string `json:"masterIPV4Cidr"`
	// MaintenanceWindow is the maintenance window provided for GKE to perform
	// upgrades if enabled.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	MaintenanceWindow string `json:"maintenanceWindow"`
	// Tags is a collection of tags (resource labels) to apply to the GCP resources which make up this cluster
	// +kubebuilder:validation:Optional
	Tags map[string]string `json:"tags,omitempty"`
	// NodePools is the set of node pools for this cluster. Required unless ALL deprecated properties except subnetwork are set.
	// +kubebuilder:validation:Optional
	NodePools []GKENodePool `json:"nodePools,omitempty"`

	// DEPRECATED: Set on node group instead, this property is now ignored. Size is the number of nodes per zone which should exist in the cluster.
	// +kubebuilder:validation:Optional
	Size int64 `json:"size,omitempty"`
	// DEPRECATED: Set on node group instead, this property is now ignored. MaxSize assuming the autoscaler is enabled this is the maximum number
	// nodes permitted
	// +kubebuilder:validation:Optional
	MaxSize int64 `json:"maxSize,omitempty"`
	// DEPRECATED: Set on node group instead, this property is now ignored. DiskSize is the size of the disk used by the compute nodes.
	// +kubebuilder:validation:Optional
	DiskSize int64 `json:"diskSize,omitempty"`
	// DEPRECATED: Set on node group instead, this property is now ignored. ImageType is the operating image to use for the default compute pool.
	// +kubebuilder:validation:Optional
	ImageType string `json:"imageType,omitempty"`
	// DEPRECATED: Set on node group instead, this property is now ignored. MachineType is the machine type which the default nodes pool should use.
	// +kubebuilder:validation:Optional
	MachineType string `json:"machineType,omitempty"`
	// DEPRECATED: This was always ignored. May be re-introduced in future. Subnetwork is name of the GCP subnetwork which the cluster nodes
	// should reside -
	// +kubebuilder:validation:Optional
	Subnetwork string `json:"subnetwork,omitempty"`
	// DEPRECATED: Set on node group instead, this property is now ignored. EnableAutoscaler indicates if the cluster should be configured with
	// cluster autoscaling turned on
	// +kubebuilder:validation:Optional
	EnableAutoscaler bool `json:"enableAutoscaler,omitempty"`
	// DEPRECATED: Set on node group instead, this property is now ignored. EnableAutoUpgrade indicates if the cluster should be configured with
	// auto upgrading enabled; meaning both nodes are masters are scheduled to upgrade during your maintenance window.
	// +kubebuilder:validation:Optional
	EnableAutoupgrade bool `json:"enableAutoupgrade,omitempty"`
	// DEPRECATED: Set on node group instead, this property is now ignored. EnableAutorepair indicates if the cluster should be configured with
	// auto repair is enabled
	// +kubebuilder:validation:Optional
	EnableAutorepair bool `json:"enableAutorepair,omitempty"`
	// DEPRECATED: Not used - now projects are created automatically, always use default.
	// Network is the GCP network the cluster reside on, which have to be unique within the GCP project and created beforehand.
	// +kubebuilder:validation:Optional
	Network string `json:"network,omitempty"`
}

// AuthorizedNetwork provides a definition for the authorized networks
type AuthorizedNetwork struct {
	// Name provides a descriptive name for this network
	Name string `json:"name"`
	// CIDR is the network range associated to this network
	CIDR string `json:"cidr"`
}

// GKENodePool represents a node pool within a GKE cluster
type GKENodePool struct {
	// Name provides a descriptive name for this node pool - must be unique within cluster
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
	// EnableAutoscaler indicates if the node pool should be configured with
	// autoscaling turned on
	// +kubebuilder:validation:Optional
	EnableAutoscaler bool `json:"enableAutoscaler"`
	// EnableAutorepair indicates if the node pool should automatically repair failed nodes
	// +kubebuilder:validation:Optional
	EnableAutorepair bool `json:"enableAutorepair"`
	// Version is the initial kubernetes version which the node group should be
	// configured with. '-' gives the same version as the master, 'latest' gives most recent,
	// 1.15 would be latest 1.15.x release, 1.15.1 would be the latest 1.15.1 release, and
	// 1.15.1-gke.1 would be the exact specified version. Must be within 2 minor versions of
	// the master version (e.g. master 1.16 supports node versios 1.14-1.16). If
	// ReleaseChannel set on cluster, this must be blank.
	// +kubebuilder:validation:Optional
	Version string `json:"version"`
	// EnableAutoUpgrade indicates if the node group should be configured with autograding
	// enabled. This must be true if the cluster has ReleaseChannel set.
	// +kubebuilder:validation:Optional
	EnableAutoupgrade bool `json:"enableAutoupgrade"`
	// Size is the number of nodes per zone which should exist in the cluster. If
	// auto-scaling is enabled, this will be the initial size of the node pool.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Required
	Size int64 `json:"size"`
	// MinSize assuming the autoscaler is enabled this is the maximum number
	// nodes permitted
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Required
	MinSize int64 `json:"minSize"`
	// MaxSize assuming the autoscaler is enabled this is the maximum number
	// nodes permitted
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Required
	MaxSize int64 `json:"maxSize"`
	// MaxPodsPerNode controls how many pods can be scheduled onto each node in this pool
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Optional
	MaxPodsPerNode int64 `json:"maxPodsPerNode"`
	// MachineType controls the type of nodes used in this node pool
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	MachineType string `json:"machineType"`
	// ImageType controls the operating system image of nodes used in this node pool
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	ImageType string `json:"imageType"`
	// DiskSize is the size of the disk used by the compute nodes.
	// +kubebuilder:validation:Minimum=10
	// +kubebuilder:validation:Required
	DiskSize int64 `json:"diskSize"`
	// Preemptible controls whether to use pre-emptible nodes.
	// +kubebuilder:validation:Optional
	Preemptible bool `json:"preemptible"`
	// Labels is a set of labels to help Kubernetes workloads find this group
	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`
	// Taints are a collection of kubernetes taints applied to the node on provisioning
	// +kubebuilder:validation:Optional
	Taints []NodeTaint `json:"taints,omitempty"`
}

// NodeTaint is the structure of a taint on a nodepool
// https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/
type NodeTaint struct {
	// Key provides the key definition for this tainer
	Key string `json:"key,omitempty"`
	// Value is arbitrary value for this taint to compare
	Value string `json:"value,omitempty"`
	// Effect is desired action on the taint
	Effect string `json:"effect,omitempty"`
}

// GKEStatus defines the observed state of GKE
// +k8s:openapi-gen=true
type GKEStatus struct {
	// Conditions is the status of the components
	Conditions corev1.Components `json:"conditions,omitempty"`
	// CACertificate is the certificate for this cluster
	CACertificate string `json:"caCertificate,omitempty"`
	// Endpoint is the endpoint of the cluster
	Endpoint string `json:"endpoint,omitempty"`
	// Status provides a overall status
	Status corev1.Status `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GKE is the Schema for the gkes API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=gkes,scope=Namespaced
// +kubebuilder:printcolumn:name="Description",type="string",JSONPath=".spec.description",description="A description of the GKE cluster"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".status.endpoint",description="The endpoint of the gke cluster"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status",description="The overall status of the cluster"
type GKE struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GKESpec   `json:"spec,omitempty"`
	Status GKEStatus `json:"status,omitempty"`
}

// Ownership returns a owner reference
func (g *GKE) Ownership() corev1.Ownership {
	return corev1.Ownership{
		Group:     GroupVersion.Group,
		Version:   GroupVersion.Version,
		Kind:      "GKE",
		Namespace: g.Namespace,
		Name:      g.Name,
	}
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GKEList contains a list of GKE
type GKEList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GKE `json:"items"`
}
