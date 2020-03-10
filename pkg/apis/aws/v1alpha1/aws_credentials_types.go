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
