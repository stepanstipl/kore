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
	"encoding/json"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	core "github.com/appvia/kore/pkg/apis/core/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EKSVPCApplier interface {
	ApplyEKSVPC(*EKSVPC)
}

// EKSVPCSpec defines the desired state of EKSVPC
// +k8s:openapi-gen=true
type EKSVPCSpec struct {
	// Credentials is a reference to an AWSCredentials object to use for authentication
	// +kubebuilder:validation:Required
	// +k8s:openapi-gen=false
	Credentials corev1.Ownership `json:"credentials"`
	// Cluster refers to the cluster this object belongs to
	// +kubebuilder:validation:Required
	Cluster corev1.Ownership `json:"cluster,omitempty"`
	// PrivateIPV4Cidr is the private range used for the VPC
	// +kubebuilder:validation:Required
	PrivateIPV4Cidr string `json:"privateIPV4Cidr"`
	// Region is the AWS region of the VPC and any resources created
	// +kubebuilder:validation:Required
	Region string `json:"region"`
}

// Infra defines types that cannot be specified at creation time
// These values are discovered from infrastructure AFTER a create
// It is provided as a convienece for caching values
type Infra struct {
	// PrivateSubnetIds is a list of subnet IDs to use for the worker nodes
	// +listType=set
	PrivateSubnetIDs []string `json:"privateSubnetIDs,omitempty"`
	// PublicSubnetIDs is a list of subnet IDs to use for resources that need a public IP (e.g. load balancers)
	// +listType=set
	PublicSubnetIDs []string `json:"publicSubnetIDs,omitempty"`
	// SecurityGroupIds is a list of security group IDs to use for a cluster
	// +listType=set
	SecurityGroupIDs []string `json:"securityGroupIDs,omitempty"`
	// PublicIPV4EgressAddresses provides the source addresses for traffic coming from the cluster
	// - can provide input for securing Kube API endpoints in managed clusters
	PublicIPV4EgressAddresses []string `json:"ipv4EgressAddresses,omitempty"`
}

// EKSVPCStatus defines the observed state of a VPC
// +k8s:openapi-gen=true
type EKSVPCStatus struct {
	// Conditions is the status of the components
	Conditions core.Components `json:"conditions,omitempty"`
	// Status provides a overall status
	Status core.Status `json:"status,omitempty"`
	// Infra provides a cache of values discovered from infrastructure
	// k8s:openapi-gen=false
	Infra Infra `json:"infra,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EKSVPC is the Schema for the eksvpc API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=eksvpcs,scope=Namespaced
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status",description="The overall status of the vpc"
type EKSVPC struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EKSVPCSpec   `json:"spec,omitempty"`
	Status EKSVPCStatus `json:"status,omitempty"`
}

// NewEKSVPC creates a new EKSVPC object
func NewEKSVPC(name, namespace string) *EKSVPC {
	return &EKSVPC{
		TypeMeta: metav1.TypeMeta{
			Kind:       "EKSVPC",
			APIVersion: GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func (e *EKSVPC) GetStatus() (corev1.Status, string) {
	return e.Status.Status, ""
}

func (e *EKSVPC) SetStatus(status corev1.Status) {
	e.Status.Status = status
}

func (e *EKSVPC) GetComponents() corev1.Components {
	return e.Status.Conditions
}

func (e *EKSVPC) ApplyClusterConfiguration(cluster *clustersv1.Cluster) error {
	if err := json.Unmarshal(cluster.Spec.Configuration.Raw, &e.Spec); err != nil {
		return err
	}

	e.Spec.Cluster = cluster.Ownership()
	e.Spec.Credentials = cluster.Spec.Credentials

	return nil
}

func (e *EKSVPC) ComponentDependencies() []string {
	return nil
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EKSVPCList contains a list of EKSVPC objects
type EKSVPCList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EKSVPC `json:"items"`
}
