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
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EKSVPCSpec defines the desired state of EKSVPC
// +k8s:openapi-gen=true
type EKSVPCSpec struct {
	// Credentials is a reference to an AWSCredentials object to use for authentication
	// +kubebuilder:validation:Required
	// +k8s:openapi-gen=false
	Credentials corev1.Ownership `json:"credentials"`
	// Cluster refers to the cluster this object belongs to
	// +kubebuilder:validation:Required
	Cluster corev1.Ownership `json:"cluster,omitempty"`
	// PrivateIPV4Cidr is the private range used for the VPC
	// +kubebuilder:validation:Required
	PrivateIPV4Cidr string `json:"privateIPV4Cidr"`
	// Region is the AWS region of the VPC and any resources created
	// +kubebuilder:validation:Required
	Region string `json:"region"`
}

// Infra defines types that cannot be specified at creation time
// These values are discovered from infrastructure AFTER a create
// It is provided as a convienece for caching values
type Infra struct {
	// VpcID is the identifier of the VPC
	VpcID string `json:"vpcID,omitempty"`
	// AvailabilityZoneIDs is the list of AZ ids
	AvailabilityZoneIDs []string `json:"availabilityZoneIDs,omitempty"`
	// AvailabilityZoneIDs is the list of AZ names
	AvailabilityZoneNames []string `json:"availabilityZoneNames,omitempty"`
	// PrivateSubnetIds is a list of subnet IDs to use for the worker nodes
	PrivateSubnetIDs []string `json:"privateSubnetIDs,omitempty"`
	// PublicSubnetIDs is a list of subnet IDs to use for resources that need a public IP (e.g. load balancers)
	PublicSubnetIDs []string `json:"publicSubnetIDs,omitempty"`
	// SecurityGroupIds is a list of security group IDs to use for a cluster
	SecurityGroupIDs []string `json:"securityGroupIDs,omitempty"`
	// PublicIPV4EgressAddresses provides the source addresses for traffic coming from the cluster
	// - can provide input for securing Kube API endpoints in managed clusters
	PublicIPV4EgressAddresses []string `json:"ipv4EgressAddresses,omitempty"`
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
