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
	core "github.com/appvia/kore/pkg/apis/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EKSNodeGroupSpec defines the desired state of EKSNodeGroup
// +k8s:openapi-gen=true
type EKSNodeGroupSpec struct {
	AMIType string `json:"amiType,omitempty"`
	// +kubebuilder:validation:Required
	ClusterName  string            `json:"clusterName"`
	DiskSize     int64             `json:"diskSize,omitempty"`
	InstanceType string            `json:"instanceType,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
	// +kubebuilder:validation:Required
	NodeGroupName string `json:"nodeGroupName"`
	// +kubebuilder:validation:Required
	IamNodeRole    string `json:"iamNodeRole"`
	ReleaseVersion string `json:"releaseVersion,omitempty"`
	RemoteAccess   string `json:"remoteAccess,omitempty"`
	DesiredSize    int64  `json:"desiredSize,omitempty"`
	// +kubebuilder:validation:Maximum=100
	MaxSize int64 `json:"maxSize,omitempty"`
	// +kubebuilder:validation:Minimum=0
	MinSize int64 `json:"minSize,omitempty"`
	// +kubebuilder:validation:Required
	// +listType=set
	Subnets []string `json:"subnets"`
	// The metadata to apply to the node group
	Tags map[string]string `json:"tags,omitempty"`
	// The Kubernetes version to use for your managed nodes
	Version string `json:"version,omitempty"`
	// AWS region to launch node group within, must match the region of the cluster
	Region string `json:"region"`
	// The security groups that are allowed SSH access (port 22) to the worker nodes
	// +listType=set
	SSHSourceSecurityGroups []string `json:"sshSourceSecurityGroups,omitempty"`
	// The Amazon EC2 SSH key that provides access for SSH communication with
	// the worker nodes in the managed node group
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html
	EC2SSHKey string `json:"eC2SSHKey,omitempty"`
	// Credentials is a reference to an AWSCredentials object to use for authentication
	// +kubebuilder:validation:Required
	// +k8s:openapi-gen=false
	Credentials core.Ownership `json:"credentials"`
}

// EKSNodeGroupStatus defines the observed state of EKSNodeGroup
// +k8s:openapi-gen=true
type EKSNodeGroupStatus struct {
	// Conditions is the status of the components
	Conditions *core.Components `json:"conditions,omitempty"`
	// Status provides a overall status
	Status core.Status `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EKSNodeGroup is the Schema for the eksnodegroups API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=eksnodegroups,scope=Namespaced
// +kubebuilder:printcolumn:name="Description",type="string",JSONPath=".spec.description",description="A description of the EKS cluster nodegroup"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status",description="The overall status of the cluster nodegroup"
type EKSNodeGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EKSNodeGroupSpec   `json:"spec,omitempty"`
	Status EKSNodeGroupStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EKSNodeGroupList contains a list of EKSNodeGroup
type EKSNodeGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EKSNodeGroup `json:"items"`
}
