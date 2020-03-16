package v1alpha1

import (
	core "github.com/appvia/kore/pkg/apis/core/v1"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GCPAdminProjectSpec defines the desired state of GCPAdminProject
// +k8s:openapi-gen=true
type GCPAdminProjectSpec struct {
	// ParentType is the type of parent this project has
	// Valid types are: "organization", "folder", and "project"
	// +kubebuilder:validation:Enum=organization;folder;project
	// +kubebuilder:validation:Required
	ParentType string `json:"parentType"`
	// ParentID is the type specific ID of the parent this project has
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	ParentID string `json:"parentID"`
	// BillingAccountName is the resource name of the billing account associated with the project
	// e.g. '012345-567890-ABCDEF'
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	BillingAccount string `json:"billingAccount"`
	// ServiceAccount is the name used when creating the service account
	// e.g. 'hub-admin'
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	ServiceAccount string `json:"serviceAccount"`
	// TokenRef is a reference to an ephemeral oauth token used provision the admin project
	// +kubebuilder:validation:Optional
	TokenRef *v1.SecretReference `json:"tokenRef,omitempty"`
	// CredentialsRef is a reference to the credentials used to provision provision
	// the projects - this is either created by dynamically from the oauth token or
	// provided for us
	// +kubebuilder:validation:Optional
	CredentialsRef *v1.SecretReference `json:"credentialsRef"`
}

// GCPAdminProjectStatus defines the observed state of GCPAdminProject
// +k8s:openapi-gen=true
type GCPAdminProjectStatus struct {
	// Conditions is a set of components conditions
	Conditions *core.Components `json:"conditions,omitempty"`
	// Project is the GCP project ID
	ProjectID string `json:"projectID,omitempty"`
	// Status provides a overall status
	Status core.Status `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GCPAdminProject is the Schema for the gcpadminprojects API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=gcpadminprojects,scope=Namespaced
type GCPAdminProject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GCPAdminProjectSpec   `json:"spec,omitempty"`
	Status GCPAdminProjectStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GCPAdminProjectList contains a list of GCPAdminProject
type GCPAdminProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GCPAdminProject `json:"items"`
}
