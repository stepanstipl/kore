/*
 * Copyright (C) 2019 Appvia Ltd <info@appvia.io>
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

package kore

import (
	"context"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

// Allocations is the interface to team allocations
type Allocations interface {
	// Delete is responsible for deleting an allocation
	Delete(context.Context, string) (*configv1.Allocation, error)
	// Exists check if an allocation exists
	Exists(context.Context, string) (bool, error)
	// IsPermitted checks if a resource is permitted access
	IsPermitted(context.Context, corev1.Ownership) (bool, error)
	// Get retrieves an allocation the kore
	Get(context.Context, string) (*configv1.Allocation, error)
	// GetAssigned returns an assigned allocation
	GetAssigned(context.Context, string) (*configv1.Allocation, error)
	// List returns a list of all the allocations
	List(context.Context) (*configv1.AllocationList, error)
	// ListAllocationsAssigned returns a list of all allocations shared to me
	ListAllocationsAssigned(context.Context) (*configv1.AllocationList, error)
	// Update is responsible for updating / creating an allocation
	Update(context.Context, *configv1.Allocation) error
}

// acaImpl is the allocations interface
type acaImpl struct {
	*hubImpl
	// team and namespace
	team string
}

// IsPermitted checks if a team is permitted access to a resource via an allocation
func (a acaImpl) IsPermitted(ctx context.Context, resource corev1.Ownership) (bool, error) {
	// @step: we list all allocation in the remote team
	list := &configv1.AllocationList{}

	// @step: if the namespaces are the same we can continue
	if resource.Namespace == a.team {
		log.Info("skipping the permission check as team and resource are in the same team namespace")

		return true, nil
	}

	err := a.Store().Client().List(ctx,
		store.ListOptions.InTo(list),
		store.ListOptions.InNamespace(resource.Namespace),
	)
	if err != nil {
		log.WithError(err).Error("attempting to list allocations from team")

		return false, err
	}
	// @step: iterate the allocations and check for my team name or allteams

	for _, x := range list.Items {
		// @step: does this point to the resource?
		if !IsOwn(x.Spec.Resource, resource) {
			continue
		}
		// do we have an all teams allocation?
		if utils.Contains(configv1.AllTeams, x.Spec.Teams) {
			return true, nil
		}
		// does out team exist in the allocation?
		if utils.Contains(a.team, x.Spec.Teams) {
			return true, nil
		}
	}

	return false, nil
}

// Exists check if an allocation exists
func (a acaImpl) Exists(ctx context.Context, name string) (bool, error) {
	return a.Store().Client().Has(ctx,
		store.HasOptions.From(&configv1.Allocation{}),
		store.HasOptions.InNamespace(a.team),
		store.HasOptions.WithName(name),
	)
}

// Delete is responsible for deleting an allocation
func (a acaImpl) Delete(ctx context.Context, name string) (*configv1.Allocation, error) {
	logger := log.WithFields(log.Fields{
		"name": name,
		"team": a.team,
	})
	logger.Info("deleting the allocation in team")

	// @step: check the allocation exists
	object, err := a.Get(ctx, name)
	if err != nil {
		logger.WithError(err).Error("failed to retrieve allocation")

		return nil, err
	}

	return object, a.Store().Client().Delete(ctx, store.DeleteOptions.From(object))
}

// Get retrieves an allocation the kore
func (a acaImpl) Get(ctx context.Context, name string) (*configv1.Allocation, error) {
	object := &configv1.Allocation{}

	if err := a.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(a.team),
		store.GetOptions.InTo(object),
		store.GetOptions.WithName(name),
	); err != nil {
		if kerrors.IsNotFound(err) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return object, nil
}

// GetAssigned returns an assigned allocation
func (a acaImpl) GetAssigned(ctx context.Context, name string) (*configv1.Allocation, error) {
	list, err := a.ListAllocationsAssigned(ctx)
	if err != nil {
		log.WithError(err).Error("trying to retrieve list of assigned allocations")

		return nil, err
	}

	for _, x := range list.Items {
		if x.Name == name {
			return &x, nil
		}
	}

	return nil, ErrNotFound
}

// ListAllocationsAssigned returns a list of all allocations which you have access to
func (a acaImpl) ListAllocationsAssigned(ctx context.Context) (*configv1.AllocationList, error) {
	// @step: find all in the kore
	all := &configv1.AllocationList{}

	if err := a.Store().Client().List(ctx,
		store.ListOptions.InAllNamespaces(),
		store.ListOptions.InTo(all),
	); err != nil {
		log.WithError(err).Error("trying to retrieve a list of all allocations")

		return nil, err
	}

	list := &configv1.AllocationList{}

	// @step: find anything for use or all teams
	for _, x := range all.Items {
		if utils.Contains("*", x.Spec.Teams) || utils.Contains(a.team, x.Spec.Teams) {
			list.Items = append(list.Items, x)
		}
		// add anything owned by us
		if x.Namespace == a.team {
			list.Items = append(list.Items, x)
		}
	}

	return list, nil
}

// List returns a list of all the allocations
func (a acaImpl) List(ctx context.Context) (*configv1.AllocationList, error) {
	items := &configv1.AllocationList{}

	return items, a.Store().Client().List(ctx,
		store.ListOptions.InNamespace(a.team),
		store.ListOptions.InTo(items),
	)
}

// Update is responsible for updating / creating an allocation
func (a acaImpl) Update(ctx context.Context, allocation *configv1.Allocation) error {
	logger := log.WithFields(log.Fields{
		"group":              allocation.Spec.Resource.Group,
		"kind":               allocation.Spec.Resource.Kind,
		"resource.name":      allocation.Spec.Resource.Name,
		"resource.namespace": allocation.Spec.Resource.Namespace,
		"teams":              allocation.Spec.Teams,
		"version":            allocation.Spec.Resource.Version,
	})
	logger.Info("attempting to create allocation for resource")

	// @step: ensure our namespace
	if allocation.Namespace != a.team {
		return ErrNotAllowed{message: "allocation must be within your team"}
	}

	// @step: ensure the resource exists in our namespace - though it will be
	// picked up the controller anyhow
	if allocation.Spec.Resource.Namespace != a.team {
		return ErrNotAllowed{message: "you cannot allocate a resource which you do not own"}
	}

	return a.Store().Client().Update(ctx,
		store.UpdateOptions.To(allocation),
		store.UpdateOptions.WithCreate(true),
		store.UpdateOptions.WithForce(true),
	)
}
