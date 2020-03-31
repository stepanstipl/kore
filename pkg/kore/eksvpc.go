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

	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/store"
	kerrors "k8s.io/apimachinery/pkg/api/errors"

	log "github.com/sirupsen/logrus"
)

// EKSVPC is the eks vpc interface
type EKSVPC interface {
	// Delete is responsible for deleting a eks vpc
	Delete(context.Context, string) error
	// Get return the definition from the api
	Get(context.Context, string) (*eks.EKSVPC, error)
	// List returns all the eks vpcs for the team
	List(context.Context) (*eks.EKSVPCList, error)
	// Update is used to update the eks vpc definition
	Update(context.Context, *eks.EKSVPC) (*eks.EKSVPC, error)
}

type eksVPCImpl struct {
	*cloudImpl
	// team is the request team
	team string
}

// Delete is responsible for deleting a eks environment
func (h *eksVPCImpl) Delete(ctx context.Context, name string) error {
	logger := log.WithFields(log.Fields{
		"name": name,
		"team": h.team,
	})
	authentication.MustGetIdentity(ctx)

	// @step: retrieve the cluster
	eksvpc := &eks.EKSVPC{}
	err := h.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(h.team),
		store.GetOptions.InTo(eksvpc),
		store.GetOptions.WithName(name),
	)
	if err != nil {
		logger.WithError(err).Error("trying to retrieve the eks vpc from api")

		return err
	}

	if eksvpc.Namespace != h.team {
		logger.Warn("attempting to delete an eks vpc from another team")

		return NewErrNotAllowed("you cannot delete an eks vpc from another team")
	}

	// @step: issue the request to remove the cluster
	return h.Store().Client().Delete(ctx, store.DeleteOptions.From(eksvpc))
}

// Get return the definition from the api
func (h *eksVPCImpl) Get(ctx context.Context, name string) (*eks.EKSVPC, error) {
	vpc := &eks.EKSVPC{}

	err := h.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(h.team),
		store.GetOptions.InTo(vpc),
		store.GetOptions.WithName(name),
	)
	if err != nil {
		log.WithError(err).Error("trying to retrieve the vpc")

		return nil, err
	}

	return vpc, nil
}

// List returns all the gke cluster in the team
func (h *eksVPCImpl) List(ctx context.Context) (*eks.EKSVPCList, error) {
	list := &eks.EKSVPCList{}

	return list, h.Store().Client().List(ctx,
		store.ListOptions.InNamespace(h.team),
		store.ListOptions.InTo(list),
	)
}

// Update is called to update or create a eks instance
func (h *eksVPCImpl) Update(ctx context.Context, vpc *eks.EKSVPC) (*eks.EKSVPC, error) {
	logger := log.WithFields(log.Fields{
		"name": vpc.Name,
		"team": h.team,
	})

	// @TODO perform any checks on the vpc options before processing
	vpc.Namespace = h.team

	// @TODO check the user is a admin within the team - i.e they have the permission
	// to update the cluster

	// @step: we need to check if team has access to the credentials
	permitted, err := h.Teams().Team(h.team).Allocations().IsPermitted(ctx, vpc.Spec.Credentials)
	if err != nil {
		logger.WithError(err).Error("trying to check for credentials allocation")

		return nil, err
	}
	if !permitted {
		logger.Warn("team attempting to build an eks vpc with credentials that have not been allocated")

		return nil, NewErrNotAllowed("the requested credentials have not been allocated to you")
	}

	// @step: we need to check if the update is permitted by eks
	_, err = h.Get(ctx, vpc.Name)
	if err != nil {
		if !kerrors.IsNotFound(err) {
			logger.WithError(err).Error("trying to retrieve eks vpc")

			return nil, err
		}

		// @TODO: we need to check if what they are updating is permitted

	}

	// @step: update the resource in the api
	if err := h.Store().Client().Update(ctx,
		store.UpdateOptions.To(vpc),
		store.UpdateOptions.WithCreate(true),
		store.UpdateOptions.WithForce(true),
		store.UpdateOptions.WithPatch(true),
	); err != nil {
		logger.WithError(err).Error("trying to update the eks vpc")

		// @TODO update the audit

		return nil, err
	}

	return vpc, nil
}
