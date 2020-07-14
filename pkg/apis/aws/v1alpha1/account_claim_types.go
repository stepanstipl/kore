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
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AccountClaimSpec defines the desired state of AccountClaim
// +k8s:openapi-gen=true
type AccountClaimSpec struct {
	// AccountName is the name of the account to create
	// +kubebuilder:validation:Required
	AccountName string `json:"accountName"`
	// Organization is the AWS organization
	// +kubebuilder:validation:Required
	Organization corev1.Ownership `json:"organization"`
}

// AccountClaimStatus defines the observed state of AWS Account
// +k8s:openapi-gen=true
type AccountClaimStatus struct {
	// CredentialRef is the reference to the credentials secret
	CredentialRef *v1.SecretReference `json:"credentialRef,omitempty"`
	// Conditions is a set of components conditions
	Conditions *corev1.Components `json:"conditions,omitempty"`
	// AccountID is the aws account id
	AccountID string `json:"accountID,omitempty"`
	// AccountRef is a reference to the underlying aws account
	AccountRef corev1.Ownership `json:"accountRef,omitempty"`
	// Status provides a overall status
	Status corev1.Status `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSAccountClaim is the Schema for the AccountClaims API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=awsaccountclaims,scope=Namespaced
type AWSAccountClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AccountClaimSpec   `json:"spec,omitempty"`
	Status AccountClaimStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSAccountClaimList contains a list of AccountClaim
type AWSAccountClaimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSAccountClaim `json:"items"`
}
