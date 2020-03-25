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

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils/validation"

	log "github.com/sirupsen/logrus"
)

// Cluster returns the clusters interface
type Clusters interface {
	// Delete is used to delete a cluster from the kore
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

// Delete is used to delete a cluster from the kore
func (c *clsImpl) Delete(ctx context.Context, name string) (*clustersv1.Kubernetes, error) {
	user := authentication.MustGetIdentity(ctx)
	logger := log.WithFields(log.Fields{
		"cluster": name,
		"team":    c.team,
		"user":    user.Username(),
	})
	logger.Info("attempting to delete the cluster")

	original, err := c.Get(ctx, name)
	if err != nil {
		logger.WithError(err).Error("trying to retrieve the cluster")

		return nil, err
	}

	// @step: check if we have any namespace on the cluster
	list, err := c.Teams().Team(c.team).NamespaceClaims().List(ctx)
	if err != nil {
		logger.WithError(err).Error("trying to list any namespace claims")

		return nil, err
	}
	for _, x := range list.Items {
		if x.Spec.Cluster.Namespace == c.team && x.Spec.Cluster.Name == name {
			return nil, ErrNotAllowed{message: "cluster has allocated namespaces please delete first"}
		}
	}

	return original, c.Store().Client().Delete(ctx, store.DeleteOptions.From(original))
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
	authentication.MustGetIdentity(ctx)

	cluster.Namespace = c.team

	// @TODO wider validation of the supplied details.
	if len(cluster.Name) > 40 {
		return validation.NewError("cluster has failed validation").
			WithFieldError("cluster.name", validation.MaxLength, "name must be 40 characters or less")
	}

	return c.Store().Client().Update(ctx,
		store.UpdateOptions.To(cluster),
		store.UpdateOptions.WithCreate(true),
	)
}
