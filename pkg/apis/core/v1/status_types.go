/*
 * Copyright (C) 2019 Rohith Jayawardene <gambol99@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package v1

// Status is the status of a thing
type Status string

const (
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
}

// HasComponent returns the status of a component
func (c *Components) HasComponent(name string) (Status, bool) {
	comp, found := c.GetStatus(name)
	if !found {
		return Unknown, false
	}

	return comp.Status, true
}

// IsHealthy checks all the conditions
func (c *Components) IsHealthy() bool {
	for _, x := range *c {
		switch x.Status {
		case FailureStatus, WarningStatus:
			return false
		}
	}

	return true
}

// SetStatus sets the status of a component
func (c *Components) SetStatus(name string, status Status) {
	comp, found := c.GetStatus(name)
	if found {
		comp.Status = status
	}
}

// SetCondition sets the state of a component
func (c *Components) SetCondition(component Component) {
	item, found := c.GetStatus(component.Name)
	if !found {
		*c = append(*c, &component)

		return
	}
	item.Status = component.Status
	item.Message = component.Message
	item.Detail = component.Detail
}

// GetStatus checks if the status of a component exists
func (c *Components) GetStatus(name string) (*Component, bool) {
	for _, x := range *c {
		if x.Name == name {
			return x, true
		}
	}

	return nil, false
}
