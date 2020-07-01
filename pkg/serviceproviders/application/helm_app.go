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

package application

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type HelmAppConfiguration struct {
	// Source is the Helm app source. Either a Git or a Helm repository
	Source HelmAppSource `json:"source,omitempty"`
	// Values are parameters for the resource templates, which can be referenced as {{ .Values.foo }}
	Values YAMLMap `json:"values,omitempty"`
	// ResourceKinds is a list of Kubernetes resource kinds to monitor
	ResourceKinds []metav1.GroupKind `json:"resourceKinds,omitempty"`
	// ResourceSelector is a label query over `resourceKinds`
	ResourceSelector *ResourceSelector `json:"resourceSelector,omitempty"`
}

type ResourceSelector struct {
	// MatchLabels is a map of {key,value} pairs.
	MatchLabels map[string]string `json:"matchLabels,omitempty"`
}

// HelmAppSource is the Helm app source. Either a Git or a Helm repository
type HelmAppSource struct {
	// GitRepository describes a Helm chart sourced from Git
	GitRepository *GitRepository `json:"git,omitempty"`
	// HelmRepository describes a Helm chart sourced from a Helm repository
	HelmRepository *HelmRepository `json:"helm,omitempty"`
}

// GitRepository describes a Helm chart sourced from Git.
type GitRepository struct {
	// Git URL is the URL of the Git repository
	URL string `json:"url"`
	// Ref is the Git branch (or other reference) to use
	Ref string `json:"ref"`
	// Path is the path to the chart relative to the repository root.
	Path string `json:"path"`
}

// HelmRepository describes a Helm chart sourced from a Helm repository.
type HelmRepository struct {
	// RepoURL is the URL of the Helm repository
	URL string `json:"url"`
	// Name is the name of the Helm chart
	Name string `json:"name"`
	// Version is the targeted Helm chart version
	Version string `json:"version"`
}
