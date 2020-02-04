/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package kore

import (
	"context"

	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/store"

	log "github.com/sirupsen/logrus"
)

/*
	These should probably be moved into some type of self registered
	handlers - which handles the resource specific elements
*/

// GKECredentials is the gke interface
type GKECredentials interface {
	// Delete is responsible for deleting a gke environment
	Delete(context.Context, string) error
	// Get return the definition from the api
	Get(context.Context, string) (*gke.GKECredentials, error)
	// List returns all the gke cluster in the team
	List(context.Context) (*gke.GKECredentialsList, error)
	// Update is used to update the gke cluster definition
	Update(context.Context, *gke.GKECredentials) (*gke.GKECredentials, error)
}

type gkeCredsImpl struct {
	*cloudImpl
	// team is the request team
	team string
}

// Delete is responsible for deleting a gke environment
func (h *gkeCredsImpl) Delete(ctx context.Context, name string) error {
	logger := log.WithFields(log.Fields{
		"name": name,
		"team": h.team,
	})
	user := authentication.MustGetIdentity(ctx)
	if !user.IsGlobalAdmin() {
		return ErrUnauthorized
	}

	creds := &gke.GKECredentials{}
	if err := h.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(h.team),
		store.GetOptions.InTo(creds),
	); err != nil {
		logger.WithError(err).Error("trying to retrieve the credentials")

		return err
	}

	// @TODO check the user us an admin in the team

	// @step: check if we have any namespaces allocated to teams

	// @TODO add an audit entry indicating the request to remove the option

	// @step: issue the request to remove the cluster
	return h.Store().Client().Delete(ctx, store.DeleteOptions.From(creds))
}

// Get return the definition from the api
func (h *gkeCredsImpl) Get(ctx context.Context, name string) (*gke.GKECredentials, error) {
	cluster := &gke.GKECredentials{}

	return cluster, h.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(h.team),
		store.GetOptions.InTo(cluster),
		store.GetOptions.WithName(name),
	)
}

// List returns all the gke cluster in the team
func (h *gkeCredsImpl) List(ctx context.Context) (*gke.GKECredentialsList, error) {
	list := &gke.GKECredentialsList{}

	return list, h.Store().Client().List(ctx,
		store.ListOptions.InNamespace(h.team),
		store.ListOptions.InTo(list),
	)
}

// Update is called to update or create a gke instance
func (h *gkeCredsImpl) Update(ctx context.Context, cluster *gke.GKECredentials) (*gke.GKECredentials, error) {
	logger := log.WithFields(log.Fields{
		"name": cluster.Name,
		"team": h.team,
	})

	// @step: update the resource in the api
	if err := h.Store().Client().Update(ctx,
		store.UpdateOptions.To(cluster),
		store.UpdateOptions.WithCreate(true),
		store.UpdateOptions.WithForce(true),
		store.UpdateOptions.WithPatch(true),
	); err != nil {
		logger.WithError(err).Error("trying to update the gke cluster")

		// @TODO update the audit

		return nil, err
	}

	return cluster, nil
}
