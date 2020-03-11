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

// EKSClusterSpec defines the desired state of EKSCluster
// +k8s:openapi-gen=true
type EKSClusterSpec struct {
	// Credentials is a reference to an AWSCredentials object to use
	// for authentication
	// +kubebuilder:validation:Required
	Credentials core.Ownership `json:"credentials,omitempty"`
	// RoleARN is the role arn which provides permissions to EKS.
	// +kubebuilder:validation:Optional
	RoleARN string `json:"roleARN,omitempty"`
	// +kubebuilder:validation:Optional
	Version string `json:"version,omitempty"`
	// Region is the AWS region which the EKS cluster should be provisioned.
	// +kubebuilder:validation:Required
	Region string `json:"region"`
	// SubnetID is a collection of subnet id's which the EKS cluster should
	// be attached to - if not defined we will provision on behalf of
	// +kubebuilder:validation:Optional
	// +listType=set
	SubnetID []string `json:"subnetID,omitempty"`
	// SecurityGroupID is a list of security group IDs which the EKS cluster
	// should be attached to - If not defined we will provision on behalf of
	// +kubebuilder:validation:Optional
	// +listType=set
	SecurityGroupID []string `json:"securityGroupID,omitempty"`
	// VPC is the AWS VPC Id which the EKS cluster should reside. If not defined
	// we will provision on your behalf.
	VPC string `json:"vpc,omitempty"`
}

// EKSClusterStatus defines the observed state of EKSCluster
// +k8s:openapi-gen=true
type EKSClusterStatus struct {
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

// EKSCluster is the Schema for the eksclusters API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=eksclusters,scope=Namespaced
type EKSCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EKSClusterSpec   `json:"spec,omitempty"`
	Status EKSClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EKSClusterList contains a list of EKSCluster
type EKSClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EKSCluster `json:"items"`
}
