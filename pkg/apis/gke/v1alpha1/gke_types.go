/*
Copyright 2019 Appvia Ltd <info@appvia.io>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	core "github.com/appvia/kore/pkg/apis/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GKESpec defines the desired state of GKE
// +k8s:openapi-gen=true
type GKESpec struct {
	// GKECredentials is a reference to the gke credentials object to use
	// +kubebuilder:validation:Required
	// +k8s:openapi-gen=false
	Credentials core.Ownership `json:"credentials"`
	// Description provides a short summary / description of the cluster.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Description string `json:"description"`
	// Version is the initial kubernetes version which the cluster should be
	// configured with.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Version string `json:"version"`
	// Size is the number of nodes per zone which should exist in the cluster.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Required
	Size int64 `json:"size"`
	// MaxSize assuming the autoscaler is enabled this is the maximum number
	// nodes permitted
	// +kubebuilder:validation:Minimum=2
	// +kubebuilder:validation:Required
	MaxSize int64 `json:"maxSize"`
	// DiskSize is the size of the disk used by the compute nodes.
	// +kubebuilder:validation:Minimum=100
	// +kubebuilder:validation:Required
	DiskSize int64 `json:"diskSize"`
	// ImageType is the operating image to use for the default compute pool.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	ImageType string `json:"imageType"`
	// MachineType is the machine type which the default nodes pool should use.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	MachineType string `json:"machineType"`
	// AuthorizedMasterNetworks is a collection of authorized networks which is
	// permitted to speak to the kubernetes API, default to all if not provided.
	// +kubebuilder:validation:Optional
	// +listType=set
	AuthorizedMasterNetworks []*AuthorizedNetwork `json:"authorizedMasterNetworks"`
	// Network is the GCP network the cluster reside on, which have
	// to be unique within the GCP project and created beforehand.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Network string `json:"network"`
	// Subnetwork is name of the GCP subnetwork which the cluster nodes
	// should reside
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Subnetwork string `json:"subnetwork"`
	// ServicesIPV4Cidr is an optional network cidr configured for the cluster
	// services
	// +kubebuilder:validation:Optional
	ServicesIPV4Cidr string `json:"servicesIPV4Cidr"`
	// ClusterIPV4Cidr is an optional network CIDR which is used to place the
	// pod network on
	// +kubebuilder:validation:Optional
	ClusterIPV4Cidr string `json:"clusterIPV4Cidr"`
	// EnableAutorepair indicates if the cluster should be configured with
	// auto repair is enabled
	// +kubebuilder:validation:Optional
	EnableAutorepair bool `json:"enableAutorepair"`
	// EnableAutoscaler indicates if the cluster should be configured with
	// cluster autoscaling turned on
	// +kubebuilder:validation:Optional
	EnableAutoscaler bool `json:"enableAutoscaler"`
	// EnableAutoUpgrade indicates if the cluster should be configured with
	// autograding enabled; meaning both nodes are masters are autoscated scheduled
	// to upgrade during your maintenance window.
	// +kubebuilder:validation:Optional
	EnableAutoupgrade bool `json:"enableAutoupgrade"`
	// EnableHorizontalPodAutoscaler indicates if the cluster is configured with
	// the horizontal pod autoscaler addon. This automatically adjusts the cpu and
	// memory resources of pods in accordances with their demand. You should ensure
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
	// EnableSheildedNodes indicates we should enable the sheilds nodes options in GKE.
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
	// Tags is a collection of tags related to the cluster type
	// +kubebuilder:validation:Optional
	Tags map[string]string `json:"tags,omitempty"`
}

// AuthorizedNetwork provides a definition for the authorized networks
type AuthorizedNetwork struct {
	// Name provides a descriptive name for this network
	Name string `json:"name"`
	// CIDR is the network range associated to this network
	CIDR string `json:"cidr"`
}

// GKEStatus defines the observed state of GKE
// +k8s:openapi-gen=true
type GKEStatus struct {
	// Conditions is the status of the components
	Conditions *core.Components `json:"conditions,omitempty"`
	// CACertificate is the certificate for this cluster
	CACertificate string `json:"caCertificate,omitempty"`
	// Endpoint is the endpoint of the cluster
	Endpoint string `json:"endpoint,omitempty"`
	// Status provides a overall status
	Status core.Status `json:"status,omitempty"`
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GKEList contains a list of GKE
type GKEList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GKE `json:"items"`
}
