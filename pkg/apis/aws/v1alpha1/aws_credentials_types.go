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

// AWSCredentialsSpec defines the desired state of AWSCredential
// +k8s:openapi-gen=true
type AWSCredentialsSpec struct {
	// AccessKeyID is the AWS access key credentials
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	AccessKeyID string `json:"accessKeyID,omitempty"`
	// AccountID is the AWS account these credentials reside
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	AccountID string `json:"accountID,omitempty"`
	// SecretAccessKey is the AWS secret key credentials containing the permissions
	// to provision EKS
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	SecretAccessKey string `json:"secretAccessKey,omitempty"`
}

// AWSCredentialsStatus defines the observed state of AWSCredential
// +k8s:openapi-gen=true
type AWSCredentialsStatus struct {
	// Verified checks that the credentials are ok and valid
	Verified bool `json:"verified"`
	// Status provides a overall status
	Status core.Status `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSCredential is the Schema for the awscredentials API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=awscredentials,scope=Namespaced
// +kubebuilder:printcolumn:name="AccountID",type="string",JSONPath=".spec.accountId",description="The AWS account ID for the credentials"
// +kubebuilder:printcolumn:name="AccessKeyID",type="string",JSONPath=".spec.accessKeyId",description="The AWS access key we are using"
// +kubebuilder:printcolumn:name="Verified",type="string",JSONPath=".status.verified",description="Indicates if the credentials have been verified"
type AWSCredentials struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSCredentialsSpec   `json:"spec,omitempty"`
	Status AWSCredentialsStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSCredentialsList contains a list of AWSCredentials
type AWSCredentialsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSCredentials `json:"items"`
}
