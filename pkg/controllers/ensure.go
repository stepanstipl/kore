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

package controllers

import (
	"github.com/appvia/kore/pkg/kore"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	DefaultEnsureHandler = EnsureRunner{}
)

// EnsureRunner provides a wrapper for running ensure funcs
type EnsureRunner struct{}

// Run is a generic handler for running the ensure methods
func (e *EnsureRunner) Run(ctx kore.Context, ensures []EnsureFunc) (reconcile.Result, error) {
	for _, x := range ensures {
		result, err := x(ctx)
		if err != nil {
			return reconcile.Result{}, err
		}
		if result.Requeue || result.RequeueAfter > 0 {
			return result, nil
		}
	}

	return reconcile.Result{}, nil
}
