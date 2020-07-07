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

// KoreFeatureType identifies the feature being configured
// +k8s:openapi-gen=true
type KoreFeatureType string

// KoreFeatureCosts represents the costs feature
const KoreFeatureCosts KoreFeatureType = "kore-costs"

// KoreFeatureSpec defines the desired state of the feature
// +k8s:openapi-gen=true
type KoreFeatureSpec struct {
	// Enabled identifies if this feature is enabled or not
	// +kubebuilder:validation:Required
	Enabled bool `json:"enabled"`
	// Feature identifies which feature this is
	// +kubebuilder:validation:Required
	FeatureType KoreFeatureType `json:"featureType"`
	// Configuration represents the key-value pairs to configure this feature
	// +kubebuilder:validation:Type=object
	Configuration map[string]string `json:"configuration"`
}

// KoreFeatureStatus defines the observed status of a feature
// +k8s:openapi-gen=true
type KoreFeatureStatus struct {
	// Status is overall status of the feature
	// +kubebuilder:validation:Optional
	Status corev1.Status `json:"status,omitempty"`
	// Message is the description of the current status
	// +kubebuilder:validation:Optional
	Message string `json:"message,omitempty"`
	// Components is a collection of component statuses
	// +kubebuilder:validation:Optional
	Components corev1.Components `json:"components,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KoreFeature is the Schema for a kore feature
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=korefeatures
type KoreFeature struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KoreFeatureSpec   `json:"spec,omitempty"`
	Status KoreFeatureStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KoreFeatureList contains a list of features
type KoreFeatureList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KoreFeature `json:"items"`
}
