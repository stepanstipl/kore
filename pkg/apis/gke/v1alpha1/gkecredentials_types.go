/*
Copyright 2019 Appvia Ltd <info@appvia.io>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GKECredentialsSpec defines the desired state of GKECredentials
// +k8s:openapi-gen=true
type GKECredentialsSpec struct {
	// Account is the credentials used to speak the GCP APIs; you create a service account
	// under the Cloud IAM within the project, adding the permissions 'Compute
	// Admin' role to the service account via IAM tab. Once done you can create
	// a key under 'Service Accounts' and copy and paste the JSON payload here.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Account string `json:"account"`
	// Project is the GCP project these credentias pretain to
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Project string `json:"project"`
	// Region is the GCP region you wish to the cluster to reside within
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Region string `json:"region"`
}

// GKECredentialsStatus defines the observed state of GKECredentials
// +k8s:openapi-gen=true
type GKECredentialsStatus struct {
	// Conditions is a collection of potential issues
	// +listType
	Conditions []corev1.Condition `json:"conditions,omitempty"`
	// Verified checks that the credentials are ok and valid
	Verified bool `json:"verified"`
	// Status provides a overall status
	Status corev1.Status `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GKECredentials is the Schema for the gkecredentials API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=gkecredentials,scope=Namespaced
// +kubebuilder:printcolumn:name="Region",type="string",JSONPath=".spec.region",description="The name of the GCP region the clusters will reside"
// +kubebuilder:printcolumn:name="Project",type="string",JSONPath=".spec.project",description="The name of the GCP project"
// +kubebuilder:printcolumn:name="Verified",type="string",JSONPath=".status.verified",description="Indicates is the credentials have been verified"
type GKECredentials struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GKECredentialsSpec   `json:"spec,omitempty"`
	Status GKECredentialsStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GKECredentialsList contains a list of GKECredentials
type GKECredentialsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GKECredentials `json:"items"`
}
