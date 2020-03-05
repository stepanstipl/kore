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
