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

package clusterapp

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ChartRepoRef is a reference to a chart whithin a Helm Chart repository
type ChartRepoRef struct {
	// RepoURL defines a valid Chart repo URL
	RepoURL string
	// ChartName is a chart name that must exist in the chart repo above
	ChartName string
	// Version is the version of the Chart to use from the repo
	Version string
}

// ChartGitRef a reference to a git repo with a chart
// NOTE this will poll for updates
type ChartGitRef struct {
	// GitRepo e.g. git@github.com:org/repo
	GitRepo string
	// Ref is the git ref to use to prevent unexpected updates
	Ref string
	// Path is the path within the git repository to find the helm Chart
	// e.g. charts/kore
	Path string
}

// ChartRef is where to pull the chart from
type ChartRef struct {
	ChartRepoRef *ChartRepoRef
	ChartGitRef  *ChartGitRef
}

// ChartApp describes a deployment by a chart and labels of resulting objects we which to monitor readiness with
type ChartApp struct {
	// ReleaseName is a name that will be used for the resulting HelmRelease and related objects
	ReleaseName string
	// Chart is a reference to a chart to install
	Chart ChartRef
	// DescriptiveName is the human readable name name for the component
	DescriptiveName string
	// DefaultNamespace is the default namespace to use to deploy chart resources into
	DefaultNamespace string
	// Values will be serialized into yaml and merged to provide input Chart values
	Values map[string]interface{}
	// SecretValues will be created as a kubernetes secret and serialized into yaml,
	// the values will then be merged with Values to provide input Chart values
	SecretValues map[string]interface{}
	// MatchLabels will allow us to find the resources
	MatchLabels map[string]string
	// Kinds are the objects we expect the chart to create and we will monitor for readiness
	Kinds []v1.GroupKind
}

// HelmSecret reference with metatdata values
type HelmSecret struct {
	Name      string
	Namespace string
	Secret    runtime.Object
	ValuesRef map[string]interface{}
}
