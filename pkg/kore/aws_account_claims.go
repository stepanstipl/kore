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
	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/appvia/kore/pkg/store"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

// AWSAccountClaims are the gcp project claims
type AWSAccountClaims interface {
	// Delete is responsible for deleting an aws account claim
	Delete(context.Context, string) (*aws.AWSAccountClaim, error)
	// Get return the definition from the api
	Get(context.Context, string) (*aws.AWSAccountClaim, error)
	// List returns all the aws account claims for the team
	List(context.Context) (*aws.AWSAccountClaimList, error)
	// Update is used to update the aws account claim definition
	Update(context.Context, *aws.AWSAccountClaim) (*aws.AWSAccountClaim, error)
}

type awsac struct {
	Interface
	// team is the team namespace
	team string
}

// Update is responsible for update a claim in kore
func (h awsac) Update(ctx context.Context, claim *aws.AWSAccountClaim) (*aws.AWSAccountClaim, error) {
	claim.Namespace = h.team

	err := h.Store().Client().Update(ctx,
		store.UpdateOptions.To(claim),
		store.UpdateOptions.WithCreate(true),
		store.UpdateOptions.WithForce(true),
	)
	if err != nil {
		log.WithError(err).Error("trying to update the aws account claim")

		return nil, err
	}

	return claim, nil
}

// Delete is used to delete an aws account claim from kore
func (h awsac) Delete(ctx context.Context, name string) (*aws.AWSAccountClaim, error) {
	claim := &aws.AWSAccountClaim{}

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
		log.WithError(err).Error("trying to retrieve the aws account claim")

		return nil, err
	}

	// @step: are the any cluster referring this project?
	list := &eks.EKSList{}
	if err := h.Store().Client().List(ctx,
		store.ListOptions.InAllNamespaces(),
		store.ListOptions.InTo(list),
	); err != nil {
		log.WithError(err).Error("trying to check if any eks clusters exist in the account")

		return nil, err
	}

	// @step: we need to iterate and look for references
	for _, cluster := range list.Items {
		// @check this matches the project claim
		if IsOwner(claim, cluster.Spec.Credentials) {
			return nil, ErrNotAllowed{message: "the aws account has provisioned clusters, these must be delete first"}
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
func (h awsac) Get(ctx context.Context, name string) (*aws.AWSAccountClaim, error) {
	claim := &aws.AWSAccountClaim{}

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
func (h awsac) List(ctx context.Context) (*aws.AWSAccountClaimList, error) {
	claims := &aws.AWSAccountClaimList{}

	return claims, h.Store().Client().List(ctx,
		store.ListOptions.InNamespace(h.team),
		store.ListOptions.InTo(claims),
	)
}

// Has checks if a resource exists within an available class in the scope
func (h awsac) Has(ctx context.Context, name string) (bool, error) {
	return h.Store().Client().Has(ctx,
		store.HasOptions.InNamespace(h.team),
		store.HasOptions.From(&aws.AWSAccountClaim{}),
		store.HasOptions.WithName(name),
	)
}
