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

package kore

import (
	"fmt"
	"sync"

	"github.com/appvia/kore/pkg/schema"

	"github.com/appvia/kore/pkg/utils/kubernetes"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
)

var clusterProviders = map[string]ClusterProvider{}
var cpLock = sync.Mutex{}

func RegisterClusterProvider(provider ClusterProvider) {
	cpLock.Lock()
	defer cpLock.Unlock()

	_, exists := clusterProviders[provider.Type()]
	if exists {
		panic(fmt.Errorf("cluster provider %q was already registered", provider.Type()))
	}

	clusterProviders[provider.Type()] = provider
}

func UnregisterClusterProvider(provider ClusterProvider) {
	cpLock.Lock()
	defer cpLock.Unlock()

	delete(clusterProviders, provider.Type())
}

func GetClusterProvider(providerType string) (ClusterProvider, bool) {
	cpLock.Lock()
	defer cpLock.Unlock()

	provider, found := clusterProviders[providerType]
	return provider, found
}

func ClusterProviders() map[string]ClusterProvider {
	cpLock.Lock()
	defer cpLock.Unlock()

	res := make(map[string]ClusterProvider, len(clusterProviders))
	for k, v := range clusterProviders {
		res[k] = v
	}
	return res
}

type ClusterComponents []*ClusterComponent

// Add adds a cluster component
func (c *ClusterComponents) Add(object kubernetes.Object, dependencies ...kubernetes.Object) {
	c.AddComponent(&ClusterComponent{
		Object:       object,
		Dependencies: dependencies,
	})
}

// AddComponent adds a cluster component
func (c *ClusterComponents) AddComponent(comp *ClusterComponent) {
	// @step: we always ensure we have a kind on the resource
	gvk, found, err := schema.GetGroupKindVersion(comp.Object)
	if err != nil || !found {
		panic(fmt.Errorf("resource GVK not found for %s (%T)", comp.Object.GetName(), comp.Object))
	}
	comp.Object.GetObjectKind().SetGroupVersionKind(gvk)

	*c = append(*c, comp)
}

// Find returns with the first components where the selector function returns true
func (c *ClusterComponents) Find(f func(comp ClusterComponent) bool) *ClusterComponent {
	for _, comp := range *c {
		if f(*comp) {
			return comp
		}
	}
	return nil
}

// Sort sorts the cluster components in a dependency order
// It returns an error if a circular reference is found
func (c *ClusterComponents) Sort() error {
	resolver := kubernetes.NewDependencyResolver()
	for _, comp := range *c {
		resolver.AddNode(comp, comp.Dependencies...)
	}
	sorted, err := resolver.Resolve()
	if err != nil {
		return err
	}

	var res ClusterComponents
	for _, comp := range sorted {
		res = append(res, comp.(*ClusterComponent))
	}

	*c = res
	return nil
}

// ClusterComponent is a Kubernetes Object with optional dependencies
type ClusterComponent struct {
	kubernetes.Object
	// Dependencies should return the dependencies for this component
	Dependencies []kubernetes.Object
	// Exists is set if the resource exists already
	Exists bool
	// IsProvider should mark a single component which is responsible for provisioning the cluster
	IsProvider bool
}

// ComponentID returns with a unique component id
func (c ClusterComponent) ComponentID() string {
	gvk := c.GetObjectKind().GroupVersionKind()
	return fmt.Sprintf("%s/%s/%s/%s/%s",
		gvk.Group,
		gvk.Version,
		gvk.Kind,
		c.GetNamespace(),
		c.GetName(),
	)
}

// ComponentName returns with the component's display name
func (c ClusterComponent) ComponentName() string {
	return fmt.Sprintf("%s/%s", c.GetObjectKind().GroupVersionKind().Kind, c.GetName())
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ClusterProvider
type ClusterProvider interface {
	// Type returns a unique type for the cluster provider
	Type() string
	// PlanJSONSchema returns the JSON schema for the plans belonging to this cluster
	PlanJSONSchema() string
	// DefaultPlans returns with the built-in default plans
	DefaultPlans() []configv1.Plan
	// DefaultPlanPolicy returns with the built-in default plan policy
	DefaultPlanPolicy() *configv1.PlanPolicy
	// SetComponents adds all povider-specific cluster components and updates dependencies if required
	SetComponents(Context, *clustersv1.Cluster, *ClusterComponents) error
	// BeforeComponentsUpdate runs after the components are loaded but before updated
	// The cluster components will be provided in dependency order
	BeforeComponentsUpdate(Context, *clustersv1.Cluster, *ClusterComponents) error
	// SetProviderData saves the provider data on the cluster
	// The cluster components will be provided in dependency order
	SetProviderData(Context, *clustersv1.Cluster, *ClusterComponents) error
}
