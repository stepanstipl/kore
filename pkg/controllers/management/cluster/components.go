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
	"errors"
	"fmt"

	"github.com/appvia/kore/pkg/schema"
	"github.com/appvia/kore/pkg/utils"

	"github.com/heimdalr/dag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// NewComponenets creates and returns a components
func NewComponents() (*Components, error) {
	graph := dag.NewDAG()
	root := &Vertex{ID: "root"}

	if err := graph.AddVertex(root); err != nil {
		return nil, err
	}

	return &Components{graph: graph, root: root}, nil
}

// Add is responsible for adding a new vertex to the graph
func (c *Components) Add(object runtime.Object) *Vertex {
	mo, ok := object.(metav1.Object)
	if !ok {
		panic("trying to add an object which does not implement metav1.Object")
	}

	gvk, found, err := schema.GetGroupKindVersion(object)
	if err != nil || !found {
		panic("trying to find the gvk of the resource " + err.Error())
	}

	v := &Vertex{
		ID: fmt.Sprintf("%s/%s/%s/%s/%s",
			gvk.Group,
			gvk.Version,
			gvk.Kind,
			mo.GetNamespace(),
			mo.GetName(),
		),
		Name:   fmt.Sprintf("%s/%s", gvk.Kind, mo.GetName()),
		Object: object,
	}

	if err := c.graph.AddVertex(v); err != nil {
		panic("trying to add the object to the graph")
	}

	if c.graph.GetSize() <= 0 {
		c.graph.AddEdge(c.root, v)
	}

	return v
}

// Get returns of vertex of a particular type
func (c *Components) Get(kind interface{}) ([]*Vertex, error) {
	list, err := c.Walk()
	if err != nil {
		return nil, err
	}

	var filtered []*Vertex
	for i := 0; i < len(list); i++ {
		if utils.IsEqualType(list[i].Object, kind) {
			filtered = append(filtered, list[i])
		}
	}

	return filtered, nil
}

// GetFirst returns the first found
func (c *Components) GetFirst(kind interface{}) (*Vertex, error) {
	list, err := c.Get(kind)
	if err != nil || len(list) <= 0 {
		return nil, errors.New("resource not found")
	}

	return list[0], nil
}

// Edge is used to connect to vertices
func (c *Components) Edge(from, to *Vertex) {
	if err := c.graph.AddEdge(from, to); err != nil {
		panic("trying to add edge to the graph")
	}
}

// Walk is used to walk the components in dependencies order
func (c *Components) Walk() ([]*Vertex, error) {
	var list []*Vertex

	ch, _, err := c.graph.DescendantsWalker(c.root)
	if err != nil {
		return nil, err
	}
	for v := range ch {
		list = append(list, v.(*Vertex))
	}

	return list, nil
}

// WalkFunc calls a callback on the vertex
func (c *Components) WalkFunc(callback func(*Vertex) (bool, error)) error {
	list, err := c.Walk()
	if err != nil {
		return err
	}
	for i := 0; i < len(list); i++ {
		m, err := callback(list[i])
		if err != nil {
			return err
		}
		if !m {
			return nil
		}
	}

	return nil
}

// InverseWalkFunc walks in the reverse order
func (c *Components) InverseWalkFunc(callback func(*Vertex) (bool, error)) error {
	list, err := c.Walk()
	if err != nil {
		return err
	}
	for i := len(list) - 1; i >= 0; i-- {
		m, err := callback(list[i])
		if err != nil {
			return err
		}
		if !m {
			return nil
		}
	}

	return nil
}
