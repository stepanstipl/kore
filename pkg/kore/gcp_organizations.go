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

package kore

import (
	"context"

	gcp "github.com/appvia/kore/pkg/apis/gcp/v1alpha1"
	"github.com/appvia/kore/pkg/store"

	log "github.com/sirupsen/logrus"
)

// Organizations are the gcp project orgs
type Organizations interface {
	// Delete is responsible for deleting a gke environment
	Delete(context.Context, string) (*gcp.Organization, error)
	// Get return the definition from the api
	Get(context.Context, string) (*gcp.Organization, error)
	// List returns all the gke cluster in the team
	List(context.Context) (*gcp.OrganizationList, error)
	// Update is used to update the gke cluster definition
	Update(context.Context, *gcp.Organization) (*gcp.Organization, error)
}

type gcppcl struct {
	Interface
	// team is the team namespace
	team string
}

// Update is responsible for update a org in the kore
func (h gcppcl) Update(ctx context.Context, org *gcp.Organization) (*gcp.Organization, error) {
	org.Namespace = h.team

	err := h.Store().Client().Update(ctx,
		store.UpdateOptions.To(org),
		store.UpdateOptions.WithCreate(true),
		store.UpdateOptions.WithForce(true),
	)
	if err != nil {
		log.WithError(err).Error("trying to update the gcp project org")

		return nil, err
	}

	return org, nil
}

// Delete is used to delete a project org from kore
func (h gcppcl) Delete(ctx context.Context, name string) (*gcp.Organization, error) {
	// @step: does the project even exist
	org := &gcp.Organization{}
	if err := h.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(h.team),
		store.GetOptions.InTo(org),
		store.GetOptions.WithName(name),
	); err != nil {
		log.WithError(err).Error("trying to retrieve the gcp organization")

		return nil, err
	}

	// @step: are the any project claims referring to this org
	claims := &gcp.ProjectClaimList{}

	err := h.Store().Client().List(ctx,
		store.ListOptions.InAllNamespaces(),
		store.ListOptions.InTo(claims),
	)
	if err != nil {
		log.WithError(err).Error("trying to retrieve a list of project claims")

		return nil, err
	}

	// @step: iterate the claims and ensure nothing refers to us
	for _, claim := range claims.Items {
		if claim.Spec.Organization.Namespace == org.Namespace && claim.Spec.Organization.Name == org.Name {
			return nil, NewErrNotAllowed("gcp organization already has project claims, these must be deleted first")
		}
	}

	// @tep: cool we can go already an delete this
	if err := h.Store().Client().Delete(ctx, store.DeleteOptions.From(org)); err != nil {
		log.WithError(err).Error("trying to delete the gcp organization")

		return nil, err
	}

	return org, nil
}

// Get returns the class from the kore
func (h gcppcl) Get(ctx context.Context, name string) (*gcp.Organization, error) {
	org := &gcp.Organization{}

	if found, err := h.Has(ctx, name); err != nil {
		return nil, err
	} else if !found {
		return nil, ErrNotFound
	}

	return org, h.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(h.team),
		store.GetOptions.WithName(name),
		store.GetOptions.InTo(org),
	)
}

// List returns a list of orgs
func (h gcppcl) List(ctx context.Context) (*gcp.OrganizationList, error) {
	orgs := &gcp.OrganizationList{}

	return orgs, h.Store().Client().List(ctx,
		store.ListOptions.InNamespace(h.team),
		store.ListOptions.InTo(orgs),
	)
}

// Has checks if a resource exists within an available class in the scope
func (h gcppcl) Has(ctx context.Context, name string) (bool, error) {
	return h.Store().Client().Has(ctx,
		store.HasOptions.InNamespace(h.team),
		store.HasOptions.From(&gcp.Organization{}),
		store.HasOptions.WithName(name),
	)
}
