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
	// AgentPoolProfiles is the set of node pools for this cluster.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	AgentPoolProfiles []AgentPoolProfile `json:"agentPoolProfiles"`
	// AuthorizedIPRanges are IP ranges to whitelist for incoming traffic to the API servers
	AuthorizedIPRanges []string `json:"authorizedIPRanges,omitempty"`
	// Cluster refers to the cluster this object belongs to
	// +kubebuilder:validation:Required
	Cluster corev1.Ownership `json:"cluster,omitempty"`
	// Credentials is a reference to the AKS credentials object to use
	// +kubebuilder:validation:Required
	// +k8s:openapi-gen=false
	Credentials corev1.Ownership `json:"credentials"`
	// Description provides a short summary / description of the cluster.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Description string `json:"description"`
	// DNSPrefix is the DNS prefix for the cluster
	// Must contain between 3 and 45 characters, and can contain only letters, numbers, and hyphens.
	// It must start with a letter and must end with a letter or a number.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	DNSPrefix string `json:"dnsPrefix"`
	// EnablePodSecurityPolicy indicates whether Pod Security Policies should be enabled
	// Note that this also requires role based access control to be enabled.
	// This feature is currently in preview and PodSecurityPolicyPreview for namespace Microsoft.ContainerService must be enabled.
	EnablePodSecurityPolicy bool `json:"enablePodSecurityPolicy,omitempty"`
	// EnablePrivateCluster controls whether the Kubernetes API is only exposed on the private network
	EnablePrivateCluster bool `json:"enablePrivateCluster,omitempty"`
	// KubernetesVersion is the Kubernetes version
	KubernetesVersion string `json:"kubernetesVersion,omitempty"`
	// LinuxProfile is the configuration for Linux VMs
	LinuxProfile *LinuxProfile `json:"linuxProfile,omitempty"`
	// Location is the location where the AKS cluster should be created
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Location string `json:"location"`
	// NetworkPlugin is the network plugin to use for networking. "azure" or "kubenet"
	// +kubebuilder:validation:Enum=azure;kubenet
	// +kubebuilder:validation:Required
	NetworkPlugin string `json:"networkPlugin"`
	// NetworkPolicy is the network policy to use for networking. "", "azure" or "calico"
	// +kubebuilder:validation:Enum=azure;calico
	// +kubebuilder:validation:Optional
	NetworkPolicy *string `json:"networkPolicy,omitempty"`
	// WindowsProfile is the configuration for Windows VMs
	WindowsProfile *WindowsProfile `json:"windowsProfile,omitempty"`
	// Tags is a collection of metadata tags to apply to the Azure resources which make up this cluster
	// +kubebuilder:validation:Optional
	Tags map[string]string `json:"tags,omitempty"`
}

// AgentPoolProfile represents a node pool within a GKE cluster
type AgentPoolProfile struct {
	// Name provides a descriptive name for this node pool - must be unique within cluster
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
	// Mode Type of the node pool.
	// System node pools serve the primary purpose of hosting critical system pods such as CoreDNS and tunnelfront.
	// User node pools serve the primary purpose of hosting your application pods.
	Mode string `json:"mode"`
	// EnableAutoScaling indicates if the node pool should be configured with
	// autoscaling turned on
	// +kubebuilder:validation:Optional
	EnableAutoScaling bool `json:"enableAutoScaling,omitempty"`
	// NodeImageVersion is the initial kubernetes version which the node group should be
	// configured with.
	// +kubebuilder:validation:Optional
	NodeImageVersion string `json:"nodeImageVersion,omitempty"`
	// Count is the number of nodes
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Required
	Count int64 `json:"count"`
	// MinCount assuming the autoscaler is enabled this is the maximum number
	// nodes permitted
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Required
	MinCount int64 `json:"minCount"`
	// MaxCount assuming the autoscaler is enabled this is the maximum number
	// nodes permitted
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Required
	MaxCount int64 `json:"maxCount"`
	// MaxPods controls how many pods can be scheduled onto each node in this pool
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Optional
	MaxPods int64 `json:"maxPods,omitempty"`
	// VMSize controls the type of nodes used in this node pool
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	VMSize string `json:"vmSize"`
	// OsType controls the operating system image of nodes used in this node pool
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Enum=Linux;Windows
	OsType string `json:"osType"`
	// OsDiskSizeGB is the size of the disk used by the compute nodes.
	// +kubebuilder:validation:Minimum=10
	// +kubebuilder:validation:Required
	OsDiskSizeGB int64 `json:"osDiskSizeGB"`
	// NodeLabels is a set of labels to help Kubernetes workloads find this group
	// +kubebuilder:validation:Optional
	NodeLabels map[string]string `json:"nodeLabels,omitempty"`
	// NodeTaints are a collection of kubernetes taints applied to the node on provisioning
	// +kubebuilder:validation:Optional
	NodeTaints []NodeTaint `json:"nodeTaints,omitempty"`
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
