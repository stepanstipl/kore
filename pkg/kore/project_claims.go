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
	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/store"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

// ProjectClaims are the gcp project claims
type ProjectClaims interface {
	// Delete is responsible for deleting a gke environment
	Delete(context.Context, string) (*gcp.ProjectClaim, error)
	// Get return the definition from the api
	Get(context.Context, string) (*gcp.ProjectClaim, error)
	// List returns all the gke cluster in the team
	List(context.Context) (*gcp.ProjectClaimList, error)
	// Update is used to update the gke cluster definition
	Update(context.Context, *gcp.ProjectClaim) (*gcp.ProjectClaim, error)
}

type gcppc struct {
	Interface
	// team is the team namespace
	team string
}

// Update is responsible for update a claim in kore
func (h gcppc) Update(ctx context.Context, claim *gcp.ProjectClaim) (*gcp.ProjectClaim, error) {
	claim.Namespace = h.team

	err := h.Store().Client().Update(ctx,
		store.UpdateOptions.To(claim),
		store.UpdateOptions.WithCreate(true),
		store.UpdateOptions.WithForce(true),
	)
	if err != nil {
		log.WithError(err).Error("trying to update the gcp project claim")

		return nil, err
	}

	return claim, nil
}

// Delete is used to delete a project claim from kore
func (h gcppc) Delete(ctx context.Context, name string) (*gcp.ProjectClaim, error) {
	claim := &gcp.ProjectClaim{}

	// @step: retrieve the claim from the api
	err := h.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(h.team),
		store.GetOptions.WithName(name),
		store.GetOptions.InTo(claim),
	)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		log.WithError(err).Error("trying to retrieve the gcp project claim")

		return nil, err
	}

	// @step: are the any cluster referring this project?
	list := &gke.GKEList{}
	if err := h.Store().Client().List(ctx,
		store.ListOptions.InAllNamespaces(),
		store.ListOptions.InTo(list),
	); err != nil {
		log.WithError(err).Error("trying to check if any gke clusters exist in project")

		return nil, err
	}

	// @step: we need to iterate and look for references
	for _, cluster := range list.Items {
		// @check this matches the project claim
		if IsOwner(claim, cluster.Spec.Credentials) {
			return nil, ErrNotAllowed{message: "the gcp project has provisioned clusters, these must be delete first"}
		}
	}

	// @step: else we can delete the project
	if err := h.Store().Client().Delete(ctx, store.DeleteOptions.From(claim)); err != nil {
		log.WithError(err).Error("trying to delete the claim from kore")

		return nil, err
	}

	return claim, nil
}

// Get returns the class from the kore
func (h gcppc) Get(ctx context.Context, name string) (*gcp.ProjectClaim, error) {
	claim := &gcp.ProjectClaim{}

	if found, err := h.Has(ctx, name); err != nil {
		return nil, err
	} else if !found {
		return nil, ErrNotFound
	}

	return claim, h.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(h.team),
		store.GetOptions.WithName(name),
		store.GetOptions.InTo(claim),
	)
}

// List returns a list of claims
func (h gcppc) List(ctx context.Context) (*gcp.ProjectClaimList, error) {
	claims := &gcp.ProjectClaimList{}

	return claims, h.Store().Client().List(ctx,
		store.ListOptions.InNamespace(h.team),
		store.ListOptions.InTo(claims),
	)
}

// Has checks if a resource exists within an available class in the scope
func (h gcppc) Has(ctx context.Context, name string) (bool, error) {
	return h.Store().Client().Has(ctx,
		store.HasOptions.InNamespace(h.team),
		store.HasOptions.From(&gcp.ProjectClaim{}),
		store.HasOptions.WithName(name),
	)
}
