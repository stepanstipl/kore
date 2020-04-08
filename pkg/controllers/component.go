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

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ComponentReconciler
type ComponentReconciler interface {
	Reconcile() (requeue bool, err error)
	Delete() (requeue bool, err error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Component
type Component interface {
	Name() string
	Dependencies() []string
	Reconcile() (bool, error)
	Delete() (bool, error)
}

func NewComponent(name string, dependencies []string, reconciler ComponentReconciler) Component {
	return component{
		name:         name,
		dependencies: dependencies,
		reconciler:   reconciler,
	}
}

type component struct {
	name         string
	dependencies []string
	reconciler   ComponentReconciler
}

func (r component) Name() string {
	return r.name
}

func (r component) Dependencies() []string {
	return r.dependencies
}

func (r component) Reconcile() (bool, error) {
	return r.reconciler.Reconcile()
}

func (r component) Delete() (bool, error) {
	return r.reconciler.Delete()
}
