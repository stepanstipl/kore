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

// EKSVPCSpec defines the desired state of EKSVPC
// +k8s:openapi-gen=true
type EKSVPCSpec struct {
	// PrivateIPV4Cidr is the private range used for the VPC
	// +kubebuilder:validation:Required
	PrivateIPV4Cidr string `json:"privateIPV4Cidr"`
	// Credentials is a reference to an AWSCredentials object to use for authentication
	// +kubebuilder:validation:Required
	// +k8s:openapi-gen=false
	Credentials core.Ownership `json:"credentials"`
	// ClusterName is used to indicate a cluster to create resources for
	// - it is used to tag cluster specific resources e.g. subnet resources are tagged unique to a cluster (for ELB's)
	// - this may become an array but keeping it simple in the first iteration
	// +kubebuilder:validation:Required
	// +k8s:openapi-gen=false
	ClusterName string `json:"clusterName"`
	// Region is the AWS region of the VPC and any resources created
	// +kubebuilder:validation:Required
	Region string `json:"region"`
}

// Infra defines types that cannot be specified at creation time
// These values are discovered from infrastructure AFTER a create
// It is provided as a convienece for caching values
type Infra struct {
	// NodeIAMROle is the IAM role assumed by the worker nodes themselves
	// not directly a VPC object is is easiest to track here first
	// If we need to support segregation of nodegroups for a single cluster, move to EKSNodegroup.Status object
	NodeIAMRole string `json:"nodeIAMRole,omitempty"`
	// ClusterIAMRoleARN is the role ARN which provides permissions to create and admister an EKS cluster
	// Although not part of a VPC it is a direct pre-requisite
	// If we need to support multiple clusters in a VPC, move to EKS.Status object
	ClusterIAMRoleARN string `json:"clusterIAMRoleARN,omitempty"`
	// SubnetIds is a list of subnet IDs to use for all nodes
	// +listType=set
	SubnetIDs []string `json:"subnetIDs,omitempty"`
	// SecurityGroupIds is a list of security group IDs to use for a cluster
	// +listType=set
	SecurityGroupIDs []string `json:"securityGroupIDs,omitempty"`
}

// EKSVPCStatus defines the observed state of a VPC
// +k8s:openapi-gen=true
type EKSVPCStatus struct {
	// Conditions is the status of the components
	Conditions core.Components `json:"conditions,omitempty"`
	// Status provides a overall status
	Status core.Status `json:"status,omitempty"`
	// Infra provides a cache of values discovered from infrastructure
	// k8s:openapi-gen=false
	Infra Infra `json:"infra,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EKSVPC is the Schema for the eksvpc API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=eksvpcs,scope=Namespaced
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status",description="The overall status of the vpc"
type EKSVPC struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EKSVPCSpec   `json:"spec,omitempty"`
	Status EKSVPCStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EKSVPCList contains a list of EKSVPC objects
type EKSVPCList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EKSVPC `json:"items"`
}
