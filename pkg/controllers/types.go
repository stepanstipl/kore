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
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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
	reconcile.Reconciler
}

// Interface2 is a temporary interface to introduce a new run function where the dependencies will be injected
// TODO: migrate all controllers to the new interface
type Interface2 interface {
	RegisterInterface
	Reconcile(kore.Context, reconcile.Request) (reconcile.Result, error)
	// Initialize registers dependencies and sets up watches
	Initialize(kore.Context, controller.Controller) error
}

type ManagerOptionsAware interface {
	// ManagerOptions returns the manager options
	ManagerOptions() manager.Options
}

type ControllerOptionsAware interface {
	// ControllerOptions returns the controller options
	ControllerOptions(kore.Context) controller.Options
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Client
type Client interface {
	client.Client
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Manager
type Manager interface {
	manager.Manager
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Controller
type Controller interface {
	controller.Controller
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . StatusWriter
type StatusWriter interface {
	client.StatusWriter
}

// EnsureFunc defines a method to ensure a state
type EnsureFunc func(kore.Context) (reconcile.Result, error)
