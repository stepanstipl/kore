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

package kubernetes

import (
	"fmt"
	"strings"
)

// DependencyResolver will resolve all Kubernetes object dependencies
// It will look for non-existing dependencies and cycles
//
// It uses Tarjan's strongly connected components algorithm to detect cycles.
type DependencyResolver struct {
	nodes  map[string]*graphNode
	index  int
	s      graphNodeStack
	result []Object
	cycles [][]string
}

// NewDependencyResolver creates a new dependency resolver
func NewDependencyResolver() *DependencyResolver {
	return &DependencyResolver{
		nodes: make(map[string]*graphNode),
	}
}

// AddNode adds a new node to the dependency graph
func (d *DependencyResolver) AddNode(o Object, dependencies ...Object) {
	id := d.getObjectID(o)

	if _, exists := d.nodes[id]; exists {
		panic(fmt.Errorf("%q node was already added", id))
	}

	d.nodes[id] = &graphNode{
		ID:           id,
		Object:       o,
		Dependencies: dependencies,
		Index:        -1,
	}
}

// Resolve will resolve the dependency graph and returns with the reverse topological ordering of the DAG
// formed by the Kubernetes objects
// It will throw an error of any circular references are found in the dependency graph
func (d *DependencyResolver) Resolve() ([]Object, error) {
	for _, v := range d.nodes {
		if v.Index == -1 {
			if err := d.strongConnect(v); err != nil {
				return nil, err
			}
		}
	}

	return d.result, nil
}

// strongConnect will find all the strongly connected components in the dependency graph based on Tarjan's algorithm
func (d *DependencyResolver) strongConnect(v *graphNode) error {
	v.Index = d.index
	v.LowLink = d.index
	d.index++
	d.s.Push(v)
	v.OnStack = true
	for _, dep := range v.Dependencies {
		w, ok := d.nodes[d.getObjectID(dep)]
		if !ok {
			return fmt.Errorf("dependency not found: %q", d.getObjectID(dep))
		}
		if v.ID == w.ID {
			d.cycles = append(d.cycles, []string{v.ID})
		}
		if w.Index == -1 {
			if err := d.strongConnect(w); err != nil {
				return err
			}
			v.LowLink = d.min(v.LowLink, w.LowLink)
		} else if w.OnStack {
			v.LowLink = d.min(v.LowLink, w.Index)
		}
	}
	if v.LowLink == v.Index {
		var component []string
		for {
			w := d.s.Pop()
			w.OnStack = false
			component = append(component, w.ID)
			if w.ID == v.ID {
				break
			}
		}
		if len(component) > 1 {
			return fmt.Errorf("circular reference detected: %s", strings.Join(component, ", "))
		}
		d.result = append(d.result, d.nodes[component[0]].Object)
	}

	return nil
}

func (d *DependencyResolver) getObjectID(o Object) string {
	gvk := o.GetObjectKind().GroupVersionKind()
	return fmt.Sprintf("%s/%s/%s/%s/%s", gvk.Group, gvk.Version, strings.ToLower(gvk.Kind), o.GetNamespace(), o.GetName())
}

func (d *DependencyResolver) min(i1 int, i2 int) int {
	if i1 <= i2 {
		return i1
	}
	return i2
}

type graphNode struct {
	ID           string
	Object       Object
	Dependencies []Object
	Index        int
	LowLink      int
	OnStack      bool
}

type graphNodeStack []*graphNode

func (s *graphNodeStack) Push(n *graphNode) {
	*s = append(*s, n)
}

func (s *graphNodeStack) Pop() *graphNode {
	l := len(*s)
	n := (*s)[l-1]
	*s = (*s)[:l-1]
	return n
}
