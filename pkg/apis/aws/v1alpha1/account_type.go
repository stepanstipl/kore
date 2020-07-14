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

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AccountSpec defines the desired state of AccountClaim
// +k8s:openapi-gen=true
type AccountSpec struct {
	// AccountName is the name of the account to create. We do this internally so
	// we can easily change the account name without changing the resource name
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	AccountName string `json:"accountName"`
	// Region is the default aws region resources will be created in for this account
	// +kubebuilder:validation:Required
	Region string `json:"region"`
	// Organization is a reference to the aws organisation to use
	// +kubebuilder:validation:Required
	Organization core.Ownership `json:"organization"`
	// Labels are a set of labels on the project
	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`
}

// AccountStatus defines the observed state of an AWS Account
// +k8s:openapi-gen=true
type AccountStatus struct {
	// CredentialRef is the reference to the credentials secret
	CredentialRef *v1.SecretReference `json:"credentialRef,omitempty"`
	// AccountID is the aws account id
	AccountID string `json:"accountID,omitempty"`
	// ServiceCatalogProvisioningID is the control tower account factory, service catalog provisioning record ID. If set, creation is being tracked
	ServiceCatalogProvisioningID string `json:"serviceCatalogProvisioningID,omitempty"`
	// Status provides a overall status
	Status core.Status `json:"status,omitempty"`
	// Conditions is a set of components conditions
	Conditions *core.Components `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSAccount is the Schema for the AccountClaims API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=awsaccount,scope=Namespaced
type AWSAccount struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AccountSpec   `json:"spec,omitempty"`
	Status AccountStatus `json:"status,omitempty"`
}

// Ownership creates and returns an ownership reference
func (p *AWSAccount) Ownership() corev1.Ownership {
	return corev1.Ownership{
		Group:     GroupVersion.Group,
		Kind:      "AWSAccount",
		Name:      p.Name,
		Namespace: p.Namespace,
		Version:   GroupVersion.Version,
	}
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSAccountList contains a list of AccountClaim
type AWSAccountList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSAccount `json:"items"`
}
