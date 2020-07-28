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

	aws "github.com/appvia/kore/pkg/apis/aws/v1alpha1"
	"github.com/appvia/kore/pkg/store"

	log "github.com/sirupsen/logrus"
)

// AWSOrganizations are the aws account orgs
type AWSOrganizations interface {
	// Delete is responsible for deleting a aws organisation
	Delete(context.Context, string) (*aws.AWSOrganization, error)
	// Get return the definition from the api
	Get(context.Context, string) (*aws.AWSOrganization, error)
	// List returns all the aws orgaisations in the team
	List(context.Context) (*aws.AWSOrganizationList, error)
	// Update is used to update the aws orgaisation definition
	Update(context.Context, *aws.AWSOrganization) (*aws.AWSOrganization, error)
}

type awsocl struct {
	Interface
	// team is the team namespace
	team string
}

// Update is responsible for update an aws org in kore
func (h awsocl) Update(ctx context.Context, org *aws.AWSOrganization) (*aws.AWSOrganization, error) {
	org.Namespace = h.team

	err := h.Store().Client().Update(ctx,
		store.UpdateOptions.To(org),
		store.UpdateOptions.WithCreate(true),
		store.UpdateOptions.WithForce(true),
	)
	if err != nil {
		log.WithError(err).Error("trying to update the aws account org")

		return nil, err
	}

	return org, nil
}

// Delete is used to delete a project org from kore
func (h awsocl) Delete(ctx context.Context, name string) (*aws.AWSOrganization, error) {
	// @step: does the orgaisation even exist
	org := &aws.AWSOrganization{}
	if err := h.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(h.team),
		store.GetOptions.InTo(org),
		store.GetOptions.WithName(name),
	); err != nil {
		log.WithError(err).Error("trying to retrieve the aws organization")

		return nil, err
	}

	// @step: are the any project claims referring to this org
	claims := &aws.AWSAccountClaimList{}

	err := h.Store().Client().List(ctx,
		store.ListOptions.InAllNamespaces(),
		store.ListOptions.InTo(claims),
	)
	if err != nil {
		log.WithError(err).Error("trying to retrieve a list of account claims")

		return nil, err
	}

	// @step: iterate the claims and ensure nothing refers to us
	for _, claim := range claims.Items {
		if claim.Spec.Organization.Namespace == org.Namespace && claim.Spec.Organization.Name == org.Name {
			return nil, NewErrNotAllowed("aws organization already has account claims, these must be deleted first")
		}
	}

	// @tep: cool we can go already an delete this
	if err := h.Store().Client().Delete(ctx, store.DeleteOptions.From(org)); err != nil {
		log.WithError(err).Error("trying to delete the aws organization")

		return nil, err
	}

	return org, nil
}

// Get returns the class from the kore
func (h awsocl) Get(ctx context.Context, name string) (*aws.AWSOrganization, error) {
	org := &aws.AWSOrganization{}

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
func (h awsocl) List(ctx context.Context) (*aws.AWSOrganizationList, error) {
	orgs := &aws.AWSOrganizationList{}

	return orgs, h.Store().Client().List(ctx,
		store.ListOptions.InNamespace(h.team),
		store.ListOptions.InTo(orgs),
	)
}

// Has checks if a resource exists within an available class in the scope
func (h awsocl) Has(ctx context.Context, name string) (bool, error) {
	return h.Store().Client().Has(ctx,
		store.HasOptions.InNamespace(h.team),
		store.HasOptions.From(&aws.AWSOrganization{}),
		store.HasOptions.WithName(name),
	)
}
