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

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AWSOrganizationSpec defines the desired state of an AWS Organization
// +k8s:openapi-gen=true
type AWSOrganizationSpec struct {
	// SsoUser is the user who will be the organisational account owner for all accounts
	SsoUser SSOUser `json:"ssoUser"`
	// OuName is the name of the parent Organizational Unit (OU) to use for provisioning accounts
	OuName string `json:"ouName"`
	// Region is the region where control tower is enabled in the master account
	Region string `json:"region"`
	// RoleARN is the role to assume when provisioning accounts
	RoleARN string `json:"roleARN"`
	// CredentialsRef is a reference to the credentials used to provision
	// the accounts
	CredentialsRef *v1.SecretReference `json:"credentialsRef"`
}

// SSOUser describes the details required to identify an AWS SSO user to user for all accounts
type SSOUser struct {
	// Email is the unique user email address specified for the AWS SSO user
	Email string `json:"email"`
	// FirstName is the firstname(s) field for an AWS SSO user
	FirstName string `json:"firstName"`
	// LastName is the last name of an SSO user
	LastName string `json:"lastName"`
}

// AWSOrganizationStatus defines the observed state of Organization
// +k8s:openapi-gen=true
type AWSOrganizationStatus struct {
	// Conditions is a set of components conditions
	Conditions *core.Components `json:"conditions,omitempty"`
	// AccountID is the AWS Account ID used for the master account
	AccountID string `json:"accountID,omitempty"`
	// Status provides a overall status
	Status core.Status `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSOrganization is the Schema for the organization API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=awsorganizations,scope=Namespaced
type AWSOrganization struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSOrganizationSpec   `json:"spec,omitempty"`
	Status AWSOrganizationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSOrganizationList contains a list of AWSOrganization
type AWSOrganizationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSOrganization `json:"items"`
}
