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
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Client is a raw runtime client
type Client interface {
	// Create is responsible for creating an object
	Create(context.Context, ...CreateOptionFunc) error
	// Delete is responsible for deleting a resource from the store
	Delete(context.Context, ...DeleteOptionFunc) error
	// DeleteAll is responsible for deleting all objects of a type
	DeleteAll(context.Context, ...DeleteAllOptionFunc) error
	// Get returns an object from the store or index
	Get(context.Context, ...GetOptionFunc) error
	// Has checks if a source exists
	Has(context.Context, ...HasOptionFunc) (bool, error)
	// List is responsible for retrieving a list of object from api or index
	List(context.Context, ...ListOptionFunc) error
	// Update is responsible for performing an update / creation
	Update(context.Context, ...UpdateOptionFunc) error
}

// Store is the contract to the store
type Store interface {
	// AddEventListener is responsible for adding a consumer of store events
	AddEventListener(*Listener) error
	// APIVersion returns a query interface
	APIVersion(string) Interface
	// Client provides the client runtime interface to the api
	Client() Client
	// GetFactories returns a list of factories
	GetFactories() []informers.SharedInformerFactory
	// Kind return a client request scoped to the kind
	Kind(string) Interface
	// Namespace returns operations for a namespace
	Namespace(string) Interface
	// RuntimeClient returns with a runtime client
	RuntimeClient() client.Client
	// Stop releases the resources
	Stop() error
	// WatchResource adds a resource type to the store and listens for events
	WatchResource(string) error
}

// EventType indicates the type of event
type EventType int

const (
	// Created indicates the resource was created
	Created EventType = iota
	// Deleted indicates the resources was remove
	Deleted
	// Updated indicates the resource was updated
	Updated
)

// Listener defines an upstream listener for
type Listener struct {
	EventHandlers
	// Resources is a collection of resources you wish to be notified about
	Resources []string `json:"resources"`
}

// EventHandlers is contract to receive change notify
type EventHandlers interface {
	// Created is called when a object is created / added
	Created(metav1.Object)
	// Deleted is called when an object is removed
	Deleted(metav1.Object)
	// Updated is called the object is updated
	Updated(metav1.Object, metav1.Object)
}

// EventHandlerFuncs is a handler for the above interface
type EventHandlerFuncs struct {
	// Created is called when a object is created / added
	CreatedHandlerFunc func(metav1.Object)
	// Deleted is called when an object is removed
	DeletedHandlerFunc func(metav1.Object)
	// Updated is called the object is updated
	UpdatedHandlerFunc func(metav1.Object, metav1.Object)
}

// Created is called whena resouce of the type is created
func (l *EventHandlerFuncs) Created(o metav1.Object) {
	if l.CreatedHandlerFunc != nil {
		l.CreatedHandlerFunc(o)
	}
}

// Deleted is called when an object is removed
func (l *EventHandlerFuncs) Deleted(o metav1.Object) {
	if l.DeletedHandlerFunc != nil {
		l.DeletedHandlerFunc(o)
	}
}

// Updated is called the object is updated
func (l *EventHandlerFuncs) Updated(before metav1.Object, after metav1.Object) {
	if l.UpdatedHandlerFunc != nil {
		l.UpdatedHandlerFunc(before, after)
	}
}

// Event defines are to callback the origin
type Event struct {
	// Type indicates the event type
	Type EventType `json:"type"`
	// Before was the object before
	Before metav1.Object `json:"before"`
	// After was the object after
	After metav1.Object `json:"after"`
	// Version is the resource version
	Version schema.GroupVersionResource `json:"version"`
}

// Query is the query filter
type Query struct {
	// APIVersion is the api version of the resource
	APIVersion string `json:"apiversion,omitempty"`
	// Kind is the api resource kind
	Kind string `json:"kind,omitempty"`
	// Labels is a collection of labels
	Labels map[string]string `json:"labels,omitempty"`
	// Name is the name of the resource
	Name string `json:"name,omitempty"`
	// Namespace is scoped namespace
	Namespace string `json:"namespace,omitempty"`
	// Version is the object version
	Version string `json:"version,omitempty"`
	// UID is the uid of the object
	UID string `json:"uid,omitempty"`
}

// Interface defines the query interface for the store
type Interface interface {
	// APIVersion is the api group
	APIVersion(string) Interface
	// Client returns the client interface
	Client() Client
	// Delete removes a object from the store
	Delete(string) error
	// Set sets and objects in the store
	Set(string, metav1.Object) error
	// Has checks if the resource exists in the store
	Has(string) (bool, error)
	// Get retrieves a resource from the store
	Get(string) (metav1.Object, error)
	// Label add a label filter
	Label(string, string) Interface
	// List retrieves a list of resources from the store
	List() ([]metav1.Object, error)
	// Kind adds the api kind type to the request
	Kind(string) Interface
	// Namespace is used to set the namespace
	Namespace(string) Interface
}
