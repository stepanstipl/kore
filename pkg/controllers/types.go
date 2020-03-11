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
	"context"

	"github.com/appvia/kore/pkg/kore"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s.io/client-go/rest"
)

// NameInterface defines a requirement to implement it's name
type NameInterface interface {
	// Name is the name of the controller
	Name() string
}

type RegisterInterface interface {
	NameInterface
	// Run starts the controller
	Run(context.Context, *rest.Config, kore.Interface) error
	// Stop instructs the controller to stop
	Stop(context.Context) error
}

// Interface is the contract for a controller
type Interface interface {
	RegisterInterface
	// Reconcile is the reconcile method
	Reconcile(reconcile.Request) (reconcile.Result, error)
}
