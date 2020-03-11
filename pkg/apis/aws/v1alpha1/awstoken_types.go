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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AWSTokenSpec defines the desired state of AWSToken
// +k8s:openapi-gen=true
type AWSTokenSpec struct {
	// SecretAccessKey AWS Secret Access Key
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:Required
	SecretAccessKey string `json:"secretAccessKey"`
	// AccessKeyID is the AWS Access Key ID
	// +kubebuilder:validation:MinLength=12
	// +kubebuilder:validation:MaxLength=12
	// +kubebuilder:validation:Required
	AccessKeyID string `json:"accessKeyID"`
	// SessionToken is the AWS Session Token
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:Required
	SessionToken string `json:"sessionToken"`
	// AccountID is the IS for the AWS account these credentials reside within
	// +kubebuilder:validation:MinLength=12
	// +kubebuilder:validation:MaxLength=12
	// +kubebuilder:validation:Required
	AccountID string `json:"accountID"`
	// Expiration is the expiry date time of this token
	// +kubebuilder:validation:Required
	Expiration string `json:"expiration"`
}

// AWSTokenStatus defines the observed state of AWSToken
// +k8s:openapi-gen=true
type AWSTokenStatus struct {
	// Verified checks that the credentials are ok and valid
	Verified bool `json:"verified"`
	// Status provides a overall status
	Status string `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSToken is the Schema for the awstokens API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=awstokens,scope=Namespaced
type AWSToken struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSTokenSpec   `json:"spec,omitempty"`
	Status AWSTokenStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSTokenList contains a list of AWSToken
type AWSTokenList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSToken `json:"items"`
}
