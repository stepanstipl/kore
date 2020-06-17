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

	"github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// options returns an options method
type optionsFunc func(*hubImpl)

func WithUsersService(v persistence.Interface) optionsFunc {
	return func(h *hubImpl) {
		h.persistenceMgr = v
	}
}

// DeleteOptions controls how to delete objects
type DeleteOptions struct {
	SkipCheck      bool
	IgnoreReadOnly bool
	Cascade        bool
}

func (d DeleteOptions) StoreOptions() []store.DeleteOptionFunc {
	var opts []store.DeleteOptionFunc
	if d.Cascade {
		opts = append(opts, store.DeleteOptions.PropagationPolicy(metav1.DeletePropagationForeground))
	}
	return opts
}

func (d DeleteOptions) Check(object kubernetes.Object, f func(o ...DeleteOptionFunc) error) error {
	if d.SkipCheck {
		return nil
	}

	if !d.IgnoreReadOnly && IsReadOnlyResource(object) {
		return NewErrNotAllowed(fmt.Sprintf("%q is read-only and can not be deleted", object.GetName()))
	}

	return f(d.Copy)
}

func (d DeleteOptions) Copy(opts *DeleteOptions) {
	*opts = d
}

// ResolveDeleteOptions will apply all delete option modifiers and returns with the delete options
func ResolveDeleteOptions(d []DeleteOptionFunc) *DeleteOptions {
	res := &DeleteOptions{}

	for _, f := range d {
		f(res)
	}

	return res
}

// DeleteOptionFunc is a delete option modifier function
type DeleteOptionFunc func(opts *DeleteOptions)

// DeleteOptionIgnoreReadOnly controls whether we have to ignore the read only annotation on the object
func DeleteOptionIgnoreReadOnly(ignoreReadonly bool) DeleteOptionFunc {
	return func(opts *DeleteOptions) {
		opts.IgnoreReadOnly = ignoreReadonly
	}
}

// DeleteOptionCascade controls whether we should delete all objects owned by the target object
func DeleteOptionCascade(cascade bool) DeleteOptionFunc {
	return func(opts *DeleteOptions) {
		opts.Cascade = cascade
	}
}

// DeleteOptionSkipCheck controls whether the delete checks should be skipped
func DeleteOptionSkipCheck(skipCheck bool) DeleteOptionFunc {
	return func(opts *DeleteOptions) {
		opts.SkipCheck = skipCheck
	}
}
