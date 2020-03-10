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

package v1

import (
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InstallPlanSpec defines the desired state of Allocation
// +k8s:openapi-gen=true
type InstallPlanSpec struct {
	// Approved indicates if the update has been approved
	Approved bool `json:"approved,omitempty"`
}

// InstallPlanStatus defines the observed state of Allocation
// +k8s:openapi-gen=true
type InstallPlanStatus struct {
	// Conditions is a collection of potential issues
	// +listType=set
	Conditions []corev1.Condition `json:"conditions,omitempty"`
	// Deployed is the applciation deployment parameters
	Deployed AppDeployment `json:"deployed"`
	// Update is the incoming deployment is requiring approval
	Update AppDeployment `json:"update,omitempty"`
	// Status is the general status of the resource
	Status corev1.Status `json:"status,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InstallPlan is the Schema for the allocations API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=installplans,scope=Namespaced
type InstallPlan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstallPlanSpec   `json:"spec,omitempty"`
	Status InstallPlanStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InstallPlanList contains a list of Allocation
type InstallPlanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InstallPlan `json:"items"`
}
