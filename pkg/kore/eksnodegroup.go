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
	"github.com/appvia/kore/pkg/services/users"
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
	// team is the request team
	team string
}

// Delete is responsible for deleting a eks nodegroup
func (n *eksNGImpl) Delete(ctx context.Context, name string) error {
	logger := log.WithFields(log.Fields{
		"cluster": name,
		"team":    n.team,
	})
	user := authentication.MustGetIdentity(ctx)

	nodegroup := &eks.EKSNodeGroup{}
	if err := n.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(n.team),
		store.GetOptions.InTo(nodegroup),
	); err != nil {
		logger.WithError(err).Error("trying to retrieve the nodegroup")

		return err
	}
	if nodegroup.Namespace != n.team {
		n.Audit().Record(ctx,
			users.Resource("EKSNodeGroup"),
			users.Team(n.team),
			users.User(user.Username()),
		).Event("user attempting to delete the nodegroup from kore")

		logger.Warn("attempting to delete a nodegroup from another team")

		return NewErrNotAllowed("you cannot delete a nodegroup from another team")
	}
	// @TODO check the user us an admin in the team

	// @TODO add an audit entry indicating the request to remove the option

	// @step: issue the request to remove the cluster
	return n.Store().Client().Delete(ctx, store.DeleteOptions.From(nodegroup))
}

// Get return the definition from the api
func (n *eksNGImpl) Get(ctx context.Context, name string) (*eks.EKSNodeGroup, error) {
	nodegroup := &eks.EKSNodeGroup{}

	err := n.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(n.team),
		store.GetOptions.InTo(nodegroup),
		store.GetOptions.WithName(name),
	)
	if err != nil {
		log.WithError(err).Error("trying to retrieve the nodegroup")

		return nil, err
	}

	return nodegroup, nil
}

// List returns all the gke cluster in the team
func (n *eksNGImpl) List(ctx context.Context) (*eks.EKSNodeGroupList, error) {
	list := &eks.EKSNodeGroupList{}

	return list, n.Store().Client().List(ctx,
		store.ListOptions.InNamespace(n.team),
		store.ListOptions.InTo(list),
	)
}

// Update is called to update or create a gke instance
func (n *eksNGImpl) Update(ctx context.Context, nodegroup *eks.EKSNodeGroup) (*eks.EKSNodeGroup, error) {
	logger := log.WithFields(log.Fields{
		"name": nodegroup.Name,
		"team": n.team,
	})

	// @step: update the resource in the api
	if err := n.Store().Client().Update(ctx,
		store.UpdateOptions.To(nodegroup),
		store.UpdateOptions.WithCreate(true),
		store.UpdateOptions.WithForce(true),
		store.UpdateOptions.WithPatch(true),
	); err != nil {
		logger.WithError(err).Error("trying to update the eks cluster")

		// @TODO update the audit

		return nil, err
	}

	return nodegroup, nil
}
