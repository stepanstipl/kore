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

package v1alpha1

import (
	core "github.com/appvia/kore/pkg/apis/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EKSNodeGroupSpec defines the desired state of EKSNodeGroup
// +k8s:openapi-gen=true
type EKSNodeGroupSpec struct {
	AMIType string `json:"aMIType,omitempty"`
	// +kubebuilder:validation:Required
	ClusterName string `json:"clusterName"`
	DiskSize    int64  `json:"diskSize,omitempty"`
	// +listType=set
	InstanceTypes []string          `json:"instanceTypes,omitempty"`
	Labels        map[string]string `json:"labels,omitempty"`
	NodeGroupName string            `json:"nodeGroupName"`
	// +kubebuilder:validation:Required
	NodeRole       string `json:"nodeRole"`
	ReleaseVersion string `json:"releaseVersion,omitempty"`
	RemoteAccess   string `json:"remoteAccess,omitempty"`
	DesiredSize    int64  `json:"desiredSize,omitempty"`
	// +kubebuilder:validation:Minimum=100
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
	// The Amazon EC2 SSH key that provides access for SSH communication with
	// the worker nodes in the managed node group
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html
	// +listType=set
	SourceSecurityGroups []string `json:"sourceSecurityGroups,omitempty"`
	// The security groups that are allowed SSH access (port 22) to the worker nodes
	EC2SSHKey string `json:"eC2SSHKey,omitempty"`
	// Use is a reference to an AWSCredentials object to use for authentication
	// +kubebuilder:validation:Required
	// +k8s:openapi-gen=false
	Use core.Ownership `json:"use"`
}

// EKSNodeGroupStatus defines the observed state of EKSNodeGroup
// +k8s:openapi-gen=true
type EKSNodeGroupStatus struct {
	// Status provides a overall status
	Status core.Status `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EKSNodeGroup is the Schema for the eksnodegroups API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=eksnodegroups,scope=Namespaced
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
