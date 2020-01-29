/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of hub-apiserver.
 *
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package hub

import (
	"context"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/hub/authentication"
	"github.com/appvia/kore/pkg/services/audit"
	"github.com/appvia/kore/pkg/store"

	log "github.com/sirupsen/logrus"
)

// Cluster returns the clusters interface
type Clusters interface {
	// Delete is used to delete a cluster from the hub
	Delete(context.Context, string) (*clustersv1.Kubernetes, error)
	// Get returns a specific kubernetes cluster
	Get(context.Context, string) (*clustersv1.Kubernetes, error)
	// List returns a list of cluster we have access to
	List(context.Context) (*clustersv1.KubernetesList, error)
	// Update is used to update the kubernetes object
	Update(context.Context, *clustersv1.Kubernetes) error
}

type clsImpl struct {
	*hubImpl
	// team is the name
	team string
}

// Delete is used to delete a cluster from the hub
func (c *clsImpl) Delete(ctx context.Context, name string) (*clustersv1.Kubernetes, error) {
	return nil, nil
}

// List returns a list of cluster we have access to
func (c *clsImpl) List(ctx context.Context) (*clustersv1.KubernetesList, error) {
	list := &clustersv1.KubernetesList{}

	return list, c.Store().Client().List(ctx,
		store.ListOptions.InNamespace(c.team),
		store.ListOptions.InTo(list),
	)
}

// Get returns a specific kubernetes cluster
func (c *clsImpl) Get(ctx context.Context, name string) (*clustersv1.Kubernetes, error) {
	cluster := &clustersv1.Kubernetes{}

	if err := c.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(c.team),
		store.GetOptions.InTo(cluster),
		store.GetOptions.WithName(name),
	); err != nil {
		log.WithError(err).Error("trying to retrieve the cluster")

		return nil, err
	}
	cluster.APIVersion = clustersv1.GroupVersion.String()
	cluster.Kind = "Kubernetes"

	return cluster, nil
}

// Update is used to update the kubernetes object
func (c *clsImpl) Update(ctx context.Context, cluster *clustersv1.Kubernetes) error {
	// @TODO check the user is an admin in the team
	user := authentication.MustGetIdentity(ctx)

	cluster.Namespace = c.team

	// @TODO add an entity into the audit log
	_ = c.Audit().Record(ctx,
		audit.Team(c.team),
		audit.User(user.Username()),
		audit.Type(audit.Update),
		audit.Resource(c.team+"/"+cluster.Name),
	).Event("user is update the kubernetes cluster")

	return c.Store().Client().Update(ctx,
		store.UpdateOptions.To(cluster),
		store.UpdateOptions.WithCreate(true),
	)
}
