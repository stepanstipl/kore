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
