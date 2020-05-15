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

// ProjectSpec defines the desired state of ProjectClaim
// +k8s:openapi-gen=true
type ProjectSpec struct {
	// ProjectName is the name of the project to create. We do this internally so
	// we can easily change the project name without changing the resource name
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	ProjectName string `json:"projectName"`
	// Organization is a reference to the gcp admin project to use
	// +kubebuilder:validation:Required
	Organization core.Ownership `json:"organization"`
	// Labels are a set of labels on the project
	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`
}

// ProjectStatus defines the observed state of GCP Project
// +k8s:openapi-gen=true
type ProjectStatus struct {
	// CredentialRef is the reference to the credentials secret
	CredentialRef *v1.SecretReference `json:"credentialRef,omitempty"`
	// ProjectID is the project id
	ProjectID string `json:"projectID,omitempty"`
	// Status provides a overall status
	Status core.Status `json:"status,omitempty"`
	// Conditions is a set of components conditions
	Conditions *core.Components `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Project is the Schema for the ProjectClaims API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=projects,scope=Namespaced
type Project struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProjectSpec   `json:"spec,omitempty"`
	Status ProjectStatus `json:"status,omitempty"`
}

// Ownership creates and returns an ownership reference
func (p *Project) Ownership() corev1.Ownership {
	return corev1.Ownership{
		Group:     GroupVersion.Group,
		Kind:      "Project",
		Name:      p.Name,
		Namespace: p.Namespace,
		Version:   GroupVersion.Version,
	}
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProjectList contains a list of ProjectClaim
type ProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Project `json:"items"`
}
