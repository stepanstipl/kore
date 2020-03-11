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

package store

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// CreateOptionFunc are options for creating
type CreateOptionFunc func(*rclient)

// CreateOptionsType is the create options
type CreateOptionsType struct{}

// CreateOptions are the create options
var CreateOptions CreateOptionsType

// From set the injection value
func (d CreateOptionsType) From(value runtime.Object) CreateOptionFunc {
	return func(r *rclient) {
		r.value = value
	}
}

// DeleteOptionFunc are options for deleting
type DeleteOptionFunc func(*rclient)

// DeleteOptionsType is the delete options
type DeleteOptionsType struct{}

// DeleteOptions is the delete options
var DeleteOptions DeleteOptionsType

// From set the injection value
func (d DeleteOptionsType) From(value runtime.Object) DeleteOptionFunc {
	return func(r *rclient) {
		r.value = value
	}
}

// DeleteAllOptionFunc are options for deleting
type DeleteAllOptionFunc func(*rclient)

// DeleteAllOptionsType is the delete options
type DeleteAllOptionsType struct{}

// DeleteAllOptions is the delete options
var DeleteAllOptions DeleteAllOptionsType

// InNamespace are the options for list operations
func (l DeleteAllOptionsType) InNamespace(value string) DeleteAllOptionFunc {
	return func(r *rclient) {
		r.index.Namespace(value)
	}
}

// WithLabel indicates the labels used
func (l DeleteAllOptionsType) WithLabel(k, v string) DeleteAllOptionFunc {
	return func(r *rclient) {
		r.index.Label(k, v)
	}
}

// From sets the value we are injecting into
func (l DeleteAllOptionsType) From(value runtime.Object) DeleteAllOptionFunc {
	return func(r *rclient) {
		r.value = value
	}
}

// HasOptionFunc are the options for a get operation
type HasOptionFunc func(*rclient)

// HasOptionsType are the options for the get request
type HasOptionsType struct{}

// HasOptions are the options for the get request
var HasOptions HasOptionsType

// From sets the value we are injecting into
func (g HasOptionsType) From(value runtime.Object) HasOptionFunc {
	return func(r *rclient) {
		r.value = value
	}
}

// WithLabel indicates matching a object label
func (g HasOptionsType) WithLabel(k, v string) HasOptionFunc {
	return func(r *rclient) {
		r.index.Label(k, v)
	}
}

// WithName are the options for list operations
func (g HasOptionsType) WithName(value string) HasOptionFunc {
	return func(r *rclient) {
		r.index.Name(value)
	}
}

// InNamespace are the options for list operations
func (g HasOptionsType) InNamespace(value string) HasOptionFunc {
	return func(r *rclient) {
		r.index.Namespace(value)
	}
}

// WithCache indicates we can use the cache
func (g HasOptionsType) WithCache(value bool) HasOptionFunc {
	return func(r *rclient) {
		r.useCache = value
	}
}

// GetOptionFunc are the options for a get operation
type GetOptionFunc func(*rclient)

// GetOptionsType are the options for the get request
type GetOptionsType struct{}

// GetOptions are the options for the get request
var GetOptions GetOptionsType

// InTo sets the value we are injecting into
func (g GetOptionsType) InTo(value runtime.Object) GetOptionFunc {
	return func(r *rclient) {
		r.value = value
	}
}

// WithName are the options for list operations
func (g GetOptionsType) WithName(value string) GetOptionFunc {
	return func(r *rclient) {
		r.index.Name(value)
	}
}

// InNamespace are the options for list operations
func (g GetOptionsType) InNamespace(value string) GetOptionFunc {
	return func(r *rclient) {
		r.index.Namespace(value)
	}
}

// WithCache indicates we can use the cache
func (g GetOptionsType) WithCache(value bool) GetOptionFunc {
	return func(r *rclient) {
		r.useCache = value
	}
}

// ListOptionFunc are the options for a list operation
type ListOptionFunc func(*rclient)

// ListOptionsType provide list options
type ListOptionsType struct{}

// ListOptions is the default type
var ListOptions ListOptionsType

// InNamespace are the options for list operations
func (l ListOptionsType) InNamespace(value string) ListOptionFunc {
	return func(r *rclient) {
		r.index.Namespace(value)
	}
}

// InAllNamespaces indicates all namespacs
func (l ListOptionsType) InAllNamespaces() ListOptionFunc {
	return func(r *rclient) {
		r.index.Namespace("")
	}
}

// InTo sets the value we are injecting into
func (l ListOptionsType) InTo(value runtime.Object) ListOptionFunc {
	return func(r *rclient) {
		r.value = value
	}
}

// WithLabel indicates matching a object label
func (l ListOptionsType) WithLabel(k, v string) ListOptionFunc {
	return func(r *rclient) {
		r.index.Label(k, v)
	}
}

// WithCache indicates we can use the cache
func (l ListOptionsType) WithCache(value bool) ListOptionFunc {
	return func(r *rclient) {
		r.useCache = value
	}
}

// UpdateOptionFunc are the options for a update operation
type UpdateOptionFunc func(*rclient)

// UpdateOptionsType provide update options
type UpdateOptionsType struct{}

// UpdateOptions is the update options
var UpdateOptions UpdateOptionsType

// From sets what the value was
func (d UpdateOptionsType) From(value runtime.Object) UpdateOptionFunc {
	return func(r *rclient) {
		r.current = value
	}
}

// WithCreate indicates we will create the resource if not found
func (d UpdateOptionsType) WithCreate(value bool) UpdateOptionFunc {
	return func(r *rclient) {
		r.withCreate = value
	}
}

// WithForce indicates the object should force apply the resource
func (d UpdateOptionsType) WithForce(value bool) UpdateOptionFunc {
	return func(r *rclient) {
		r.withForceApply = value
	}
}

// WithPatch indicates we will check if the resource exists and try and patch
func (d UpdateOptionsType) WithPatch(value bool) UpdateOptionFunc {
	return func(r *rclient) {
		r.withPatch = value
	}
}

// To sets where you want the value to be
func (d UpdateOptionsType) To(value runtime.Object) UpdateOptionFunc {
	return func(r *rclient) {
		r.value = value
	}
}

//
// Apply Logic
//

func (d CreateOptionsType) apply(r *rclient, options ...CreateOptionFunc) {
	for _, fn := range options {
		fn(r)
	}
}

func (d DeleteOptionsType) apply(r *rclient, options ...DeleteOptionFunc) {
	for _, fn := range options {
		fn(r)
	}
}

func (l DeleteAllOptionsType) apply(r *rclient, options ...DeleteAllOptionFunc) {
	for _, fn := range options {
		fn(r)
	}
}

func (g HasOptionsType) apply(r *rclient, options ...HasOptionFunc) {
	for _, fn := range options {
		fn(r)
	}
}

func (g GetOptionsType) apply(r *rclient, options ...GetOptionFunc) {
	for _, fn := range options {
		fn(r)
	}
}

func (l ListOptionsType) apply(r *rclient, options ...ListOptionFunc) {
	for _, fn := range options {
		fn(r)
	}
}

func (d UpdateOptionsType) apply(r *rclient, options ...UpdateOptionFunc) {
	for _, fn := range options {
		fn(r)
	}
}
