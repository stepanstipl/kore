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

	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/services/users"
	"github.com/appvia/kore/pkg/store"
	kerrors "k8s.io/apimachinery/pkg/api/errors"

	log "github.com/sirupsen/logrus"
)

// GKE is the gke interface
type GKE interface {
	// Delete is responsible for deleting a gke environment
	Delete(context.Context, string) error
	// Get return the definition from the api
	Get(context.Context, string) (*gke.GKE, error)
	// List returns all the gke cluster in the team
	List(context.Context) (*gke.GKEList, error)
	// Update is used to update the gke cluster definition
	Update(context.Context, *gke.GKE) (*gke.GKE, error)
}

type gkeImpl struct {
	*cloudImpl
	// team is the request team
	team string
}

// Delete is responsible for deleting a gke environment
func (h *gkeImpl) Delete(ctx context.Context, name string) error {
	logger := log.WithFields(log.Fields{
		"name": name,
		"team": h.team,
	})
	user := authentication.MustGetIdentity(ctx)

	// @step: retrieve the cluster
	cluster := &gke.GKE{}
	err := h.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(h.team),
		store.GetOptions.InTo(cluster),
		store.GetOptions.WithName(name),
	)
	if err != nil {
		logger.WithError(err).Error("trying to retrieve the cluster from api")

		return err
	}

	if cluster.Namespace != h.team {
		h.Audit().Record(ctx,
			users.Resource("GKE"),
			users.Team(h.team),
			users.User(user.Username()),
		).Event("user attempting to delete the cluster from kore")

		logger.Warn("attempting to delete a cluster from another team")

		return NewErrNotAllowed("you cannot delete a cluster from another team")
	}

	// @TODO check the user us an admin in the team

	// @step: check if we have any namespaces allocated to teams

	// @TODO add an audit entry indicating the request to remove the option
	h.Audit().Record(ctx,
		users.Resource("GKE"),
		users.Team(h.team),
		users.Verb(users.AuditDelete),
		users.User(user.Username()),
	).Event("user has deleted the cluster from kore")

	// @step: issue the request to remove the cluster
	return h.Store().Client().Delete(ctx, store.DeleteOptions.From(cluster))
}

// Get return the definition from the api
func (h *gkeImpl) Get(ctx context.Context, name string) (*gke.GKE, error) {
	cluster := &gke.GKE{}

	err := h.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(h.team),
		store.GetOptions.InTo(cluster),
		store.GetOptions.WithName(name),
	)
	if err != nil {
		log.WithError(err).Error("trying to retrieve the cluster")

		return nil, err
	}

	return cluster, nil
}

// List returns all the gke cluster in the team
func (h *gkeImpl) List(ctx context.Context) (*gke.GKEList, error) {
	list := &gke.GKEList{}

	return list, h.Store().Client().List(ctx,
		store.ListOptions.InNamespace(h.team),
		store.ListOptions.InTo(list),
	)
}

// Update is called to update or create a gke instance
func (h *gkeImpl) Update(ctx context.Context, cluster *gke.GKE) (*gke.GKE, error) {
	logger := log.WithFields(log.Fields{
		"name": cluster.Name,
		"team": h.team,
	})

	// @TODO perform any checks on the cluster options before processing
	cluster.Namespace = h.team

	// @TODO check the user is a admin within the team - i.e they have the permission
	// to update the cluster

	// @step: we need to check if team has access to the credentials
	permitted, err := h.Teams().Team(h.team).Allocations().IsPermitted(ctx, cluster.Spec.Credentials)
	if err != nil {
		logger.WithError(err).Error("trying to check for credentials allocation")

		return nil, err
	}
	if !permitted {
		logger.Warn("team attempting to build cluster of credentials which have not been allocated")

		return nil, NewErrNotAllowed("the requested credentials have not been allocated to you")
	}

	// @step: inform the audit service of the change

	// @step: we need to check if the update is permitted by gke
	_, err = h.Get(ctx, cluster.Name)
	if err != nil {
		if !kerrors.IsNotFound(err) {
			logger.WithError(err).Error("trying to retrieve from cluster")

			return nil, err
		}

		// @TODO: we need to check if what they are updating is permitted

	}

	// @step: update the resource in the api
	if err := h.Store().Client().Update(ctx,
		store.UpdateOptions.To(cluster),
		store.UpdateOptions.WithCreate(true),
		store.UpdateOptions.WithForce(true),
	); err != nil {
		logger.WithError(err).Error("trying to update the gke cluster")

		// @TODO update the audit

		return nil, err
	}

	return cluster, nil
}
