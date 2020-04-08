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
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/appvia/kore/pkg/utils"
)

type CompositeReconciler struct {
	logger      logrus.FieldLogger
	reconcilers map[string]Component
	finished    map[string]bool
}

func NewCompositeReconciler(logger logrus.FieldLogger) *CompositeReconciler {
	return &CompositeReconciler{
		logger:      logger,
		reconcilers: map[string]Component{},
		finished:    map[string]bool{},
	}
}

func (c *CompositeReconciler) RegisterReconciler(name string, dependencies []string, reconciler ComponentReconciler) {
	_, exists := c.reconcilers[name]
	if exists {
		panic(fmt.Errorf("%q reconciler was already registered", name))
	}
	c.reconcilers[name] = NewComponent(name, dependencies, reconciler)
}

func (c *CompositeReconciler) RegisterComponent(component Component) {
	_, exists := c.reconcilers[component.Name()]
	if exists {
		panic(fmt.Errorf("%q reconciler was already registered", component.Name()))
	}
	c.reconcilers[component.Name()] = component
}

func (c *CompositeReconciler) Reconcile() (bool, error) {
	return c.reconcile(false)
}

func (c *CompositeReconciler) Delete() (bool, error) {
	return c.reconcile(true)
}

func (c *CompositeReconciler) reconcile(isDelete bool) (bool, error) {
	var reconcileErr error

	running := 0
	for _, r := range c.reconcilers {
		if c.readyToRun(r, isDelete) {
			running++
			requeue, err := c.callReconcile(r, isDelete)
			if err != nil {
				if IsCriticalError(err) {
					return false, err
				}
				reconcileErr = utils.AppendMultiError(reconcileErr, err)
			} else if !requeue {
				c.finished[r.Name()] = true
			}
		}
	}

	unfinishedComponents := c.unfinishedComponents()
	if len(unfinishedComponents) > 0 {
		if running == 0 {
			verb := "created"
			if isDelete {
				verb = "deleted"
			}
			return false, NewCriticalError(
				fmt.Errorf("some components couldn't be %s: %s", verb, strings.Join(unfinishedComponents, ", ")),
			)
		}
		return true, reconcileErr
	}

	return false, nil
}

func (c *CompositeReconciler) callReconcile(r Component, isDelete bool) (bool, error) {
	if isDelete {
		return r.Delete()
	}
	return r.Reconcile()
}

func (c *CompositeReconciler) unfinishedComponents() []string {
	var res []string
	for _, r := range c.reconcilers {
		if !c.finished[r.Name()] {
			res = append(res, r.Name())
		}
	}
	return res
}

func (c *CompositeReconciler) readyToRun(r Component, isDelete bool) bool {
	if c.finished[r.Name()] {
		return false
	}

	if !isDelete {
		for _, dep := range r.Dependencies() {
			if !c.finished[dep] {
				return false
			}
		}
	} else {
		for _, cr := range c.reconcilers {
			for _, dep := range cr.Dependencies() {
				if dep == r.Name() && !c.finished[cr.Name()] {
					return false
				}
			}
		}
	}

	return true
}
