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

package components

import (
	"fmt"

	"github.com/appvia/kore/pkg/controllers"

	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ controllers.ComponentWithDeletionStatus = Finalizer{}

// Finalizer is a controller component for adding/removing a finalizer
type Finalizer struct {
	name   string
	object kubernetes.Object
}

// NewFinalizer creates the finalizer controller component
func NewFinalizer(name string, object kubernetes.Object) Finalizer {
	return Finalizer{name: name, object: object}
}

func (f Finalizer) Reconcile(ctx kore.Context) (reconcile.Result, error) {
	finalizer := kubernetes.NewFinalizer(ctx.Client(), f.name)

	if finalizer.NeedToAdd(f.object) {
		if err := finalizer.Add(f.object); err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to add finalizer %q to %s", f.name, kubernetes.MustGetRuntimeSelfLink(f.object))
		}

		return reconcile.Result{Requeue: true}, nil
	}

	return reconcile.Result{}, nil
}

func (f Finalizer) Delete(ctx kore.Context) (reconcile.Result, error) {
	finalizer := kubernetes.NewFinalizer(ctx.Client(), f.name)

	if err := finalizer.Remove(f.object); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (f Finalizer) IsDeleted(ctx kore.Context) (bool, error) {
	finalizer := kubernetes.NewFinalizer(ctx.Client(), f.name)
	return finalizer.NeedToAdd(f.object), nil
}
