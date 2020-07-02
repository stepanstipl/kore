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

	"github.com/appvia/kore/pkg/schema"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

// ClusterComponent is a Kubernetes Object with optional dependencies
type ClusterComponent struct {
	kubernetes.Object
	original kubernetes.Object
	// Dependencies should return the dependencies for this component
	Dependencies []kubernetes.Object
	// IsProvider should mark a single component which is responsible for provisioning the cluster
	IsProvider bool
}

// ComponentID returns with a unique component id
func (c *ClusterComponent) ComponentID() string {
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
func (c *ClusterComponent) ComponentName() string {
	return fmt.Sprintf("%s/%s", c.GetObjectKind().GroupVersionKind().Kind, c.GetName())
}

func (c *ClusterComponent) Load(ctx Context) error {
	obj, err := schema.GetScheme().New(c.Object.GetObjectKind().GroupVersionKind())
	if err != nil {
		return err
	}

	key := types.NamespacedName{Namespace: c.Object.GetNamespace(), Name: c.Object.GetName()}
	if err := ctx.Client().Get(ctx, key, obj); err != nil {
		if kerrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// The runtime client doesn't set the GVK on the result object
	obj.GetObjectKind().SetGroupVersionKind(c.Object.GetObjectKind().GroupVersionKind())

	c.original = obj.(kubernetes.Object)
	c.Object = obj.DeepCopyObject().(kubernetes.Object)
	return nil
}

func (c *ClusterComponent) Update(ctx Context) (bool, error) {
	return kubernetes.UpdateIfChangedSinceLastUpdate(ctx, ctx.Client(), c.Object, c.original)
}

// Exists returns true if the component exists in Kubernetes
func (c *ClusterComponent) Exists() bool {
	return c.original != nil
}

// Original returns the existing object
func (c *ClusterComponent) Original() kubernetes.Object {
	return c.original
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
