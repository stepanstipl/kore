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

// AKSCredentialsSpec defines the desired state of AKSCredentials
// +k8s:openapi-gen=true
type AKSCredentialsSpec struct {
	// SubscriptionID is the Azure Subscription ID
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	SubscriptionID string `json:"subscriptionID"`
	// TenantID is the Azure Tenant ID
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	TenantID string `json:"tenantID"`
	// ClientID is the Azure client ID
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	ClientID string `json:"clientID"`
	// CredentialsRef is a reference to the credentials used to create clusters
	// +kubebuilder:validation:Optional
	CredentialsRef *v1.SecretReference `json:"credentialsRef,omitempty"`
}

// AKSCredentialsStatus defines the observed state of AKSCredentials
// +k8s:openapi-gen=true
type AKSCredentialsStatus struct {
	// Conditions is a collection of potential issues
	Conditions []corev1.Condition `json:"conditions,omitempty"`
	// Verified checks that the credentials are ok and valid
	Verified *bool `json:"verified,omitempty"`
	// Status provides a overall status
	Status corev1.Status `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AKSCredentials are used for storing Azure credentials needed to create AKS clusters
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=akscredentials,scope=Namespaced
// +kubebuilder:printcolumn:name="Subscription ID",type="string",JSONPath=".spec.subscriptionID",description="Azure Subscription ID"
// +kubebuilder:printcolumn:name="Tenant ID",type="string",JSONPath=".spec.tenantID",description="Azure Tenant ID"
// +kubebuilder:printcolumn:name="Verified",type="string",JSONPath=".status.verified",description="Indicates is the credentials have been verified"
type AKSCredentials struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AKSCredentialsSpec   `json:"spec,omitempty"`
	Status AKSCredentialsStatus `json:"status,omitempty"`
}

func (in *AKSCredentials) GetStatus() (status corev1.Status, message string) {
	return in.Status.Status, ""
}

func (in *AKSCredentials) SetStatus(status corev1.Status, _ string) {
	in.Status.Status = status
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AKSCredentialsList contains a list of AKSCredentials
type AKSCredentialsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AKSCredentials `json:"items"`
}
