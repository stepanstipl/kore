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

	log "github.com/sirupsen/logrus"
)

/*
	These should probably be moved into some type of self registered
	handlers - which handles the resource specific elements
*/

// EKSNodeGroup is the eks nodegroup interface
type EKSNodeGroup interface {
	// Delete is responsible for deleting a eks nodegroup
	Delete(context.Context, string) error
	// Get return the definition from the api
	Get(context.Context, string) (*eks.EKSNodeGroup, error)
	// List returns all the gke cluster in the team
	List(context.Context) (*eks.EKSNodeGroupList, error)
	// Update is used to update the gke cluster definition
	Update(context.Context, *eks.EKSNodeGroup) (*eks.EKSNodeGroup, error)
}

type eksNGImpl struct {
	*cloudImpl
	//cluster is the request cluster
	cluster string
	// team is the request team
	team string
}

// Delete is responsible for deleting a eks nodegroup
func (h *eksNGImpl) Delete(ctx context.Context, name string) error {
	logger := log.WithFields(log.Fields{
		"name":    name,
		"team":    h.team,
		"cluster": h.cluster,
	})
	user := authentication.MustGetIdentity(ctx)
	if !user.IsGlobalAdmin() {
		return ErrUnauthorized
	}

	creds := &eks.EKSCredentials{}
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
func (h *eksCredsImpl) Get(ctx context.Context, name string) (*eks.EKSCredentials, error) {
	cluster := &eks.EKSCredentials{}

	return cluster, h.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(h.team),
		store.GetOptions.InTo(cluster),
		store.GetOptions.WithName(name),
	)
}

// List returns all the gke cluster in the team
func (h *eksCredsImpl) List(ctx context.Context) (*eks.EKSCredentialsList, error) {
	list := &eks.EKSCredentialsList{}

	return list, h.Store().Client().List(ctx,
		store.ListOptions.InNamespace(h.team),
		store.ListOptions.InTo(list),
	)
}

// Update is called to update or create a gke instance
func (h *eksCredsImpl) Update(ctx context.Context, cluster *eks.EKSCredentials) (*eks.EKSCredentials, error) {
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
		logger.WithError(err).Error("trying to update the eks cluster")

		// @TODO update the audit

		return nil, err
	}

	return cluster, nil
}
