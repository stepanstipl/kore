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

// EKSSpec defines the desired state of EKSCluster
// +k8s:openapi-gen=true
type EKSSpec struct {
	// Name the name of the EKS cluster
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// RoleARN is the role ARN which provides permissions to EKS
	// +kubebuilder:validation:MinLength=10
	// +kubebuilder:validation:Required
	RoleARN string `json:"roleARN"`
	// Version is the Kubernetes version to use
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:Required
	Version string `json:"version,omitempty"`
	// SubnetIds is a list of subnet IDs
	// +kubebuilder:validation:Required
	Region string `json:"region"`
	// AWS region to launch this cluster within
	// +kubebuilder:validation:Required
	// +listType=set
	SubnetIDs []string `json:"subnetIDs"`
	// SecurityGroupIds is a list of security group IDs
	// +kubebuilder:validation:Required
	// +listType=set
	SecurityGroupIDs []string `json:"securityGroupIDs,omitempty"`
	// Credentials is a reference to an EKSCredentials object to use for authentication
	// +k8s:openapi-gen=false
	Credentials core.Ownership `json:"credentials"`
}

// EKSStatus defines the observed state of EKS cluster
// +k8s:openapi-gen=true
type EKSStatus struct {
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

// EKS is the Schema for the eksclusters API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=eksclusters,scope=Namespaced
type EKS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EKSSpec   `json:"spec,omitempty"`
	Status EKSStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EKSList contains a list of EKS clusters
type EKSList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EKS `json:"items"`
}
