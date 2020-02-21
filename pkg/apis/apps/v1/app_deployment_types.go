/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package v1

import (
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Subscription defines a subscription
// +kubebuilder:validation:Enum=Automatic;Manual
// +kubebuilder:validation:MinLength=1
type Subscription string

const (
	// SubscriptionManual indicates the application whichs an approval
	SubscriptionManual = "Manual"
	// SubscriptionAutomatic indicates the application is upgraded automatically
	SubscriptionAutomatic = "Automatic"
)

// AppDeploymentSpec defines the desired state of Allocation
// +k8s:openapi-gen=true
type AppDeploymentSpec struct {
	// Cluster is the cluster the application should be deployed on
	// +kubebuilder:validation:Optional
	Cluster corev1.Ownership `json:"cluster"`
	// Summary is a summary of what the application is
	// +kubebuilder:validation:Required
	// +kubebuilder:vaVlidation:MinLength=1
	Summary string `json:"summary"`
	// Decription is a longer description of what the application provides
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Description string `json:"description"`
	// Package is the name of the resource being shared
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Package string `json:"package"`
	// Version is the version of the package to install
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Version string `json:"version"`
	// Source is the source of the package
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Source string `json:"source"`
	// Capabilities defines the features supported by the package
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinItems=1
	// +listType
	Capabilities []string `json:"capabilities,omitempty"`
	// Keywords keywords whuch describe the application
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +listType
	Keywords []string `json:"keywords"`
	// Vendor is the entity whom published the package
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Vendor string `json:"vendor"`
	// Official indicates if the applcation is officially published by Appvia
	// +kubebuilder:validation:Required
	Official bool `json:"official"`
	// Replaces indicates the version this replaces
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Replaces string `json:"replaces"`
	// Subscription is the nature of upgrades i.e manual or automatic
	// +kubebuilder:validation:Required
	Subscription Subscription `json:"subscription"`
	// Values are optional values suppilied to the application deployment
	// +kubebuilder:validation:Optional
	Values apiextv1.JSON `json:"values,omitempty"`
}

// AppDeploymentStatus defines the observed state of Allocation
// +k8s:openapi-gen=true
type AppDeploymentStatus struct {
	// Status is the general status of the resource
	// +kubebuilder:validation:Required
	Status corev1.Status `json:"status,omitempty"`
	// Conditions is a collection of potential issues
	// +listType
	Conditions []corev1.Condition `json:"conditions,omitempty"`
	// InstallPlan in the name of the installplan which this deployment has deployed from
	InstallPlan string `json:"installPlan,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppDeployment is the Schema for the allocations API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=appdeployments,scope=Namespaced
type AppDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppDeploymentSpec   `json:"spec,omitempty"`
	Status AppDeploymentStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppDeploymentList contains a list of Allocation
type AppDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppDeployment `json:"items"`
}
