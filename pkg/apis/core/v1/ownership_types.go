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
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

// Ownership indicates the ownership of a resource
// +k8s:openapi-gen=true
type Ownership struct {
	// Group is the api group
	Group string `json:"group"`
	// Version is the group version
	Version string `json:"version"`
	// Kind is the name of the resource under the group
	Kind string `json:"kind"`
	// Namespace is the location of the object
	Namespace string `json:"namespace"`
	// Name is name of the resource
	Name string `json:"name"`
}

func GetOwnershipFromObject(obj runtime.Object) (Ownership, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return Ownership{}, err
	}
	return Ownership{
		Group:     obj.GetObjectKind().GroupVersionKind().Group,
		Version:   obj.GetObjectKind().GroupVersionKind().Version,
		Kind:      obj.GetObjectKind().GroupVersionKind().Kind,
		Namespace: accessor.GetNamespace(),
		Name:      accessor.GetName(),
	}, nil
}

func MustGetOwnershipFromObject(obj runtime.Object) Ownership {
	ownership, err := GetOwnershipFromObject(obj)
	if err != nil {
		panic(err)
	}
	return ownership
}

func (o Ownership) IsSameType(o2 Ownership) bool {
	return strings.EqualFold(o.Group, o2.Group) &&
		strings.EqualFold(o.Version, o2.Version) &&
		strings.EqualFold(o.Kind, o2.Kind) &&
		strings.EqualFold(o.Namespace, o2.Namespace)
}

func (o Ownership) Equals(o2 Ownership) bool {
	return o.IsSameType(o2) && o.Name == o2.Name
}

func (o Ownership) HasGroupVersionKind(gvk schema.GroupVersionKind) bool {
	return strings.EqualFold(gvk.Group, o.Group) && strings.EqualFold(gvk.Version, o.Version) && strings.EqualFold(gvk.Kind, o.Kind)
}

func (o Ownership) NamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Name:      o.Name,
		Namespace: o.Namespace,
	}
}

func (o Ownership) ToUnstructured() unstructured.Unstructured {
	u := unstructured.Unstructured{}
	u.SetGroupVersionKind(o.GroupVersionKind())
	u.SetNamespace(o.Namespace)
	u.SetName(o.Name)
	return u
}

func (o Ownership) GroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   o.Group,
		Version: o.Version,
		Kind:    o.Kind,
	}
}
