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

package v1

import (
	"errors"
	"fmt"
	"strings"
)

// Status is the status of a thing
type Status string

const (
	// DeletingStatus indicates we ar deleting the resource
	DeletingStatus Status = "Deleting"
	// DeletedStatus indicates a deleted entity
	DeletedStatus Status = "Deleted"
	// DeleteFailedStatus indicates that deleting the entity failed
	DeleteFailedStatus Status = "DeleteFailed"
	// ErrorStatus indicates that a recoverable error happened
	ErrorStatus Status = "Error"
	// PendingStatus indicate we are waiting
	PendingStatus Status = "Pending"
	// SuccessStatus is a successful resource
	SuccessStatus Status = "Success"
	// FailureStatus indicates the resource has failed for one or more reasons
	FailureStatus Status = "Failure"
	// WarningStatus indicates are warning
	WarningStatus Status = "Warning"
	// Unknown is an unknown status
	Unknown Status = "Unknown"
)

func (s Status) IsFailed() bool {
	return s == FailureStatus || s == DeleteFailedStatus
}

func (s Status) OneOf(statuses ...Status) bool {
	for _, status := range statuses {
		if status == s {
			return true
		}
	}
	return false
}

// StatusAware is an interface for objects which have a status and zero or more components
type StatusAware interface {
	GetStatus() (status Status, message string)
	SetStatus(status Status)
	GetComponents() Components
}

// Condition is a reason why something failed
// +k8s:openapi-gen=true
type Condition struct {
	// Message is a human readable message
	Message string `json:"message"`
	// Detail is a actual error which might contain technical reference
	Detail string `json:"detail"`
}

// Components is a collection of components
type Components []*Component

// Component the state of a component of the resource
type Component struct {
	// Name is the name of the component
	Name string `json:"name,omitempty"`
	// Status is the status of the component
	Status Status `json:"status,omitempty"`
	// Message is a human readable message on the status of the component
	Message string `json:"message,omitempty"`
	// Detail is additional details on the error is any
	Detail string `json:"detail,omitempty"`
	// Resource is a reference to the resource
	Resource *Ownership `json:"resource,omitempty"`
}

func (c *Component) Update(status Status, message, detail string) {
	c.Status = status
	c.Message = message
	c.Detail = detail
}

func (c Component) IsFailed() bool {
	return c.Status.IsFailed()
}

func (c Component) String() string {
	if c.Message != "" && c.Detail != "" {
		return fmt.Sprintf("[%s] %s - %s: %s", c.Status, c.Name, c.Message, c.Detail)
	}

	if c.Message != "" {
		return fmt.Sprintf("[%s] %s - %s", c.Status, c.Name, c.Message)
	}

	return fmt.Sprintf("[%s] %s", c.Status, c.Name)
}

// GetStatus returns the status of a component
func (c Components) GetStatus(name string) (Status, bool) {
	comp, found := c.GetComponent(name)
	if !found {
		return Unknown, false
	}

	return comp.Status, true
}

// HasStatus returns true if any of the components has the given status
func (c Components) HasStatus(status Status) bool {
	for _, component := range c {
		if component.Status == status {
			return true
		}
	}
	return false
}

// HasStatusForAll returns true if all the components has the given status
func (c Components) HasStatusForAll(status Status) bool {
	for _, component := range c {
		if component.Status != status {
			return false
		}
	}
	return true
}

// SetStatus sets the status of a component
func (c Components) SetStatus(name string, status Status, message, detail string) {
	comp, found := c.GetComponent(name)
	if found {
		comp.Status = status
		comp.Message = message
		comp.Detail = detail
	}
}

// SetCondition sets the state of a component
func (c *Components) SetCondition(component Component) {
	item, found := c.GetComponent(component.Name)
	if !found {
		*c = append(*c, &component)

		return
	}
	item.Status = component.Status
	item.Message = component.Message
	item.Detail = component.Detail
}

// GetComponent checks if the status of a component exists
func (c Components) GetComponent(name string) (*Component, bool) {
	for _, x := range c {
		if x.Name == name {
			return x, true
		}
	}

	return nil, false
}

func (c *Components) RemoveComponent(name string) {
	for i, x := range *c {
		if x.Name == name {
			*c = append((*c)[:i], (*c)[i+1:]...)
		}
	}
}

func (c Components) Error() error {
	var messages []string
	for _, c := range c {
		if c.IsFailed() {
			messages = append(messages, c.String())
		}
	}
	if len(messages) > 0 {
		return errors.New(strings.Join(messages, ", "))
	}

	return nil
}
