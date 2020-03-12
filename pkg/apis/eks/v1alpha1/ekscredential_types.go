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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EKSCredentialsSpec defines the desired state of EKSCredential
// +k8s:openapi-gen=true
type EKSCredentialsSpec struct {
	// SecretAccessKey is the AWS Secret Access Key
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:Required
	SecretAccessKey string `json:"secretAccessKey"`
	// AccessKeyID is the AWS Access Key ID
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:Required
	AccessKeyID string `json:"accessKeyID"`
	// AccountID is the AWS account these credentials reside within
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:Required
	AccountID string `json:"accountID"`
}

// EKSCredentialsStatus defines the observed state of EKSCredential
// +k8s:openapi-gen=true
type EKSCredentialsStatus struct {
	// Conditions is a collection of potential issues
	// +listType=set
	Conditions []corev1.Condition `json:"conditions,omitempty"`
	// Verified checks that the credentials are ok and valid
	Verified *bool `json:"verified,omitempty"`
	// Status provides a overall status
	Status corev1.Status `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EKSCredentials is the Schema for the ekscredentials API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=ekscredentials,scope=Namespaced
type EKSCredentials struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EKSCredentialsSpec   `json:"spec,omitempty"`
	Status EKSCredentialsStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EKSCredentialsList contains a list of EKSCredential
type EKSCredentialsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EKSCredentials `json:"items"`
}
