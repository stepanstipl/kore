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

package cluster

import (
	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/controllers"

	"github.com/heimdalr/dag"
	"k8s.io/apimachinery/pkg/runtime"
)

// ClusterProviderComponents is what a provider must supply
type ClusterProviderComponents interface {
	// Components is used to generate the components for a cluster
	Components(cluster *clustersv1.Cluster, components *Components) controllers.EnsureFunc
	// Complete is used to complete any components from existing before applying
	Complete(cluster *clustersv1.Cluster, components *Components) controllers.EnsureFunc
	// SetProviderData saves the provider data on the cluster
	SetProviderData(cluster *clustersv1.Cluster, components *Components) controllers.EnsureFunc
}

// Components is a wrapper for the cluster components
type Components struct {
	graph *dag.DAG
	// root the root of the graph
	root *Vertex
}

// Vertex is a node in the graph
type Vertex struct {
	// ID is the name of the vertex
	ID string
	// Name is the kind and name
	Name string
	// Object is the runtime object
	Object runtime.Object
	// DisableComponentStatus indicates the cluster does not need to track the status
	DisableComponentStatus bool
	// Exists is set if the resource exists already
	Exists bool
}

// Id returns the unique id of the resource
func (v *Vertex) Id() string {
	return v.ID
}

// String returns a string representation of the resource
func (v *Vertex) String() string {
	return v.Name
}
