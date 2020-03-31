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
	"encoding/json"
	"strings"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	gkev1alpha1 "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ClusterSpec defines the desired state of a cluster
// +k8s:openapi-gen=true
type ClusterSpec struct {
	// Kind refers to the cluster type (e.g. GKE, EKS)
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Kind string `json:"kind"`
	// Plan is the name of the cluster plan which was used to create this cluster
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Plan string `json:"plan"`
	// Configuration are the configuration values for this cluster
	// It will contain values from the plan + overrides by the user
	// This will provide a simple interface to calculate diffs between plan and cluster configuration
	// +kubebuilder:validation:Type=object
	Configuration apiextv1.JSON `json:"configuration"`
	// Credentials is a reference to the credentials object to use
	// +kubebuilder:validation:Required
	Credentials corev1.Ownership `json:"credentials"`
}

// ClusterStatus defines the observed state of a cluster
// +k8s:openapi-gen=true
type ClusterStatus struct {
	// APIEndpoint is the kubernetes API endpoint url
	// +kubebuilder:validation:Optional
	APIEndpoint string `json:"apiEndpoint,omitempty"`
	// CaCertificate is the base64 encoded cluster certificate
	// +kubebuilder:validation:Optional
	CaCertificate string `json:"caCertificate,omitempty"`
	// Components is a collection of component statuses
	// +kubebuilder:validation:Optional
	Components corev1.Components `json:"components,omitempty"`
	// AuthProxyEndpoint is the endpoint of the authentication proxy for this cluster
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=1
	AuthProxyEndpoint string `json:"authProxyEndpoint,omitempty"`
	// Status is the overall status of the cluster
	// +kubebuilder:validation:Optional
	Status corev1.Status `json:"status,omitempty"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Cluster is the Schema for the plans API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=clusters
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec,omitempty"`
	Status ClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterList contains a list of clusters
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}

func (c *Cluster) CreateResources(clusterUsers []ClusterUser) ([]runtime.Object, error) {
	var res []runtime.Object
	var provider runtime.Object

	switch strings.ToLower(c.Kind) {
	case "gke":
		gkeProvider := &gkev1alpha1.GKE{
			TypeMeta: metav1.TypeMeta{
				Kind:       gkev1alpha1.GroupName,
				APIVersion: gkev1alpha1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      c.Name,
				Namespace: c.Namespace,
			},
			Spec:   gkev1alpha1.GKESpec{},
			Status: gkev1alpha1.GKEStatus{},
		}
		if err := json.Unmarshal(c.Spec.Configuration.Raw, &gkeProvider.Spec); err != nil {
			return nil, err
		}
		provider = gkeProvider
	}

	res = append(res, provider)

	kubernetes := &Kubernetes{
		TypeMeta: metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       c.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.Name,
			Namespace: c.Namespace,
		},
		Spec: KubernetesSpec{
			Provider: corev1.Ownership{
				Group:     provider.GetObjectKind().GroupVersionKind().Group,
				Kind:      provider.GetObjectKind().GroupVersionKind().Kind,
				Name:      c.Name,
				Namespace: c.Namespace,
				Version:   provider.GetObjectKind().GroupVersionKind().Version,
			},
			ClusterUsers: clusterUsers,
		},
	}
	if err := json.Unmarshal(c.Spec.Configuration.Raw, kubernetes); err != nil {
		return nil, err
	}
	res = append(res, kubernetes)

	return res, nil
}