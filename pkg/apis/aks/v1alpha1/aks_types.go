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

// AKSSpec defines the desired state of an AKS cluster
// +k8s:openapi-gen=true
type AKSSpec struct {
	// APIServerAuthorizedIPRanges are IP ranges to whitelist for incoming traffic to the API servers
	// +listType=set
	APIServerAuthorizedIPRanges []string `json:"apiServerAuthorizedIPRanges,omitempty"`
	// Cluster refers to the cluster this object belongs to
	// +kubebuilder:validation:Required
	Cluster corev1.Ownership `json:"cluster,omitempty" jsonschema:"-"`
	// Credentials is a reference to the AKS credentials object to use
	// +kubebuilder:validation:Required
	// +k8s:openapi-gen=false
	Credentials corev1.Ownership `json:"credentials" jsonschema:"-"`
	// Description provides a short summary / description of the cluster.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Description string `json:"description"`
	// DNSPrefix is the DNS prefix for the cluster
	// Must contain between 3 and 45 characters, and can contain only letters, numbers, and hyphens.
	// It must start with a letter and must end with a letter or a number.
	// +kubebuilder:validation:Required
	DNSPrefix string `json:"dnsPrefix"`
	// EnablePodSecurityPolicy indicates whether Pod Security Policies should be enabled
	// Note that this also requires role based access control to be enabled.
	// This feature is currently in preview and PodSecurityPolicyPreview for namespace Microsoft.ContainerService must be enabled.
	EnablePodSecurityPolicy bool `json:"enablePodSecurityPolicy,omitempty"`
	// Version is the Kubernetes version
	Version string `json:"version,omitempty"`
	// LinuxProfile is the configuration for Linux VMs
	LinuxProfile *LinuxProfile `json:"linuxProfile,omitempty" jsonschema:"enum=azure,enum=kubenet"`
	// NetworkPlugin is the network plugin to use for networking. "azure" or "kubenet"
	// +kubebuilder:validation:Enum=azure;kubernetes
	NetworkPlugin string `json:"networkPlugin" jsonschema:"enum=azure,enum=calico"`
	// NetworkPolicy is the network policy to use for networking. "azure" or "calico"
	// +kubebuilder:validation:Enum=azure;calico
	NetworkPolicy string `json:"networkPolicy"`
	// NodePools is the set of node pools for this cluster. Required unless ALL deprecated properties except subnetwork are set.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +listType=set
	NodePools []AKSNodePool `json:"nodePools"`
	// PrivateClusterEnabled controls whether the Kubernetes API is only exposed on the private network
	PrivateClusterEnabled bool `json:"privateClusterEnabled,omitempty"`
	// Region is the location where the AKS cluster should be created
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Region string `json:"region"`
	// WindowsProfile is the configuration for Windows VMs
	WindowsProfile *WindowsProfile `json:"windowsProfile,omitempty"`
}

// AKSNodePool represents a node pool within a GKE cluster
type AKSNodePool struct {
	// Name provides a descriptive name for this node pool - must be unique within cluster
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
	// EnableAutoscaler indicates if the node pool should be configured with
	// autoscaling turned on
	// +kubebuilder:validation:Optional
	EnableAutoscaler bool `json:"enableAutoscaler,omitempty"`
	// Version is the initial kubernetes version which the node group should be
	// configured with. '-' gives the same version as the master, 'latest' gives most recent,
	// 1.15 would be latest 1.15.x release, 1.15.1 would be the latest 1.15.1 release, and
	// 1.15.1-gke.1 would be the exact specified version. Must be within 2 minor versions of
	// the master version (e.g. master 1.16 supports node versios 1.14-1.16).
	// +kubebuilder:validation:Optional
	Version string `json:"version,omitempty"`
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
	MaxPodsPerNode int64 `json:"maxPodsPerNode,omitempty"`
	// MachineType controls the type of nodes used in this node pool
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	MachineType string `json:"machineType"`
	// ImageType controls the operating system image of nodes used in this node pool
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	ImageType string `json:"imageType" jsonschema:"enum=Linux,enum=Windows"`
	// DiskSize is the size of the disk used by the compute nodes.
	// +kubebuilder:validation:Minimum=10
	// +kubebuilder:validation:Required
	DiskSize int64 `json:"diskSize"`
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

// LinuxProfile is the configuration for Linux VMs
// +k8s:openapi-gen=true
type LinuxProfile struct {
	// AdminUsername is the admin username for Linux VMs
	AdminUsername string `json:"adminUsername"`
	// SSHPublicKeys is a list of public SSH keys to allow to connect to the Linux VMs
	// +listType=set
	SSHPublicKeys []string `json:"sshPublicKeys"`
}

// WindowsProfile is the configuration for Windows VMs
// +k8s:openapi-gen=true
type WindowsProfile struct {
	// AdminUsername is the admin username for Windows VMs
	AdminUsername string `json:"adminUsername"`
	// AdminPassword is the admin password for Windows VMs
	AdminPassword string `json:"adminPassword"`
}

// AKSStatus defines the observed state of an AKS cluster
// +k8s:openapi-gen=true
type AKSStatus struct {
	// Components is the status of the components
	Components corev1.Components `json:"components,omitempty"`
	// CACertificate is the certificate for this cluster
	CACertificate string `json:"caCertificate,omitempty"`
	// Endpoint is the endpoint of the cluster
	Endpoint string `json:"endpoint,omitempty"`
	// Status provides the overall status
	Status corev1.Status `json:"status,omitempty"`
	// Message is the status message
	Message string `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AKS is the schema for an AKS cluster object
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=aks,scope=Namespaced
// +kubebuilder:printcolumn:name="Description",type="string",JSONPath=".spec.description",description="A description of the AKS cluster"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".status.endpoint",description="The endpoint of the AKS cluster"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status",description="The overall status of AKS cluster"
type AKS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AKSSpec   `json:"spec,omitempty"`
	Status AKSStatus `json:"status,omitempty"`
}

func (a *AKS) GetStatus() (status corev1.Status, message string) {
	return a.Status.Status, a.Status.Message
}

func (a *AKS) SetStatus(status corev1.Status, message string) {
	a.Status.Status = status
	a.Status.Message = message
}

func (a *AKS) StatusComponents() *corev1.Components {
	return &a.Status.Components
}

// Ownership returns a owner reference
func (a *AKS) Ownership() corev1.Ownership {
	return corev1.Ownership{
		//Group:     GroupVersion.Group,
		//Version:   GroupVersion.Version,
		Kind:      "AKS",
		Namespace: a.Namespace,
		Name:      a.Name,
	}
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AKSList contains a list of AKS Cluster objects
type AKSList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AKS `json:"items"`
}
