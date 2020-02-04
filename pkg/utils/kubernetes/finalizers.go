/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Finalizable is an kubernetes resource api object that supports finalizers
type Finalizable interface {
	DeepCopyObject() runtime.Object
	GetFinalizers() []string
	SetFinalizers(finalizers []string)
	GetDeletionTimestamp() *metav1.Time
}

// Finalizer manages the finalizers for resources in kubernetes
type Finalizer struct {
	driver client.Client
	value  string
}

// NewFinalizer constructs a new finalizer manager
func NewFinalizer(driver client.Client, finalizerValue string) *Finalizer {
	return &Finalizer{
		driver: driver,
		value:  finalizerValue,
	}
}

// Add adds a finalizer to an object
func (c *Finalizer) Add(resource Finalizable) error {
	finalizers := append(resource.GetFinalizers(), c.value)
	resource.SetFinalizers(finalizers)

	return c.driver.Update(context.Background(), resource.DeepCopyObject())
}

// Remove removes a finalizer from an object
func (c *Finalizer) Remove(resource Finalizable) error {
	finalizers := resource.GetFinalizers()
	for idx, finalizer := range finalizers {
		if finalizer == c.value {
			finalizers = append(finalizers[:idx], finalizers[idx+1:]...)
			break
		}
	}
	resource.SetFinalizers(finalizers)

	return c.driver.Update(context.Background(), resource.DeepCopyObject())
}

// IsDeletionCandidate checks if the resource is a candidate for deletion
func (c *Finalizer) IsDeletionCandidate(resource Finalizable) bool {
	return resource.GetDeletionTimestamp() != nil && c.getIndex(resource) != -1
}

// NeedToAdd checks if the resource should have but does not have the finalizer
func (c *Finalizer) NeedToAdd(resource Finalizable) bool {
	return resource.GetDeletionTimestamp() == nil && c.getIndex(resource) == -1
}

func (c *Finalizer) getIndex(resource Finalizable) int {
	for i, v := range resource.GetFinalizers() {
		if v == c.value {
			return i
		}
	}
	return -1
}
