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

package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Finalizable is an kubernetes resource api object that supports finalizers
type Finalizable interface {
	runtime.Object
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

func (c *Finalizer) set(resource Finalizable) {
	finalizers := append(resource.GetFinalizers(), c.value)
	resource.SetFinalizers(finalizers)
}

func (c *Finalizer) unset(resource Finalizable) {
	finalizers := resource.GetFinalizers()
	for idx, finalizer := range finalizers {
		if finalizer == c.value {
			finalizers = append(finalizers[:idx], finalizers[idx+1:]...)
			break
		}
	}
	resource.SetFinalizers(finalizers)
}

// Add adds a finalizer to an object
func (c *Finalizer) Add(resource Finalizable) error {
	c.set(resource)
	return c.driver.Update(context.Background(), resource.DeepCopyObject())
}

func (c *Finalizer) AddIfNotSet(ctx context.Context, resource Finalizable) (bool, error) {
	if c.NeedToAdd(resource) {
		c.set(resource)
		original := resource.DeepCopyObject()
		return true, c.driver.Patch(ctx, resource, client.MergeFrom(original))
	}
	return false, nil
}

// Remove removes a finalizer from an object
func (c *Finalizer) Remove(resource Finalizable) error {
	c.unset(resource)
	return c.driver.Update(context.Background(), resource.DeepCopyObject())
}

func (c *Finalizer) RemovePatch(ctx context.Context, resource Finalizable) (bool, error) {
	if c.IsDeletionCandidate(resource) {
		c.unset(resource)
		original := resource.DeepCopyObject()
		return true, c.driver.Patch(ctx, resource, client.MergeFrom(original))
	}
	return false, nil
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
