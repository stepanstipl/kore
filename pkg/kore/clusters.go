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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	kerrors "k8s.io/apimachinery/pkg/api/errors"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"

	"github.com/appvia/kore/pkg/kore/assets"
	"github.com/appvia/kore/pkg/utils/jsonschema"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils/validation"

	log "github.com/sirupsen/logrus"
)

// Clusters returns the an interface for handling clusters
type Clusters interface {
	// Delete is used to delete a cluster
	Delete(context.Context, string) (*clustersv1.Cluster, error)
	// Get returns a specific cluster
	Get(context.Context, string) (*clustersv1.Cluster, error)
	// List returns a list of clusters we have access to
	List(context.Context) (*clustersv1.ClusterList, error)
	// Update is used to update the cluster
	Update(context.Context, *clustersv1.Cluster) error
}

type clustersImpl struct {
	*hubImpl
	// team is the name
	team string
}

// Delete is used to delete a cluster
func (c *clustersImpl) Delete(ctx context.Context, name string) (*clustersv1.Cluster, error) {
	// @TODO check whether the user is an admin in the team

	user := authentication.MustGetIdentity(ctx)
	if !user.IsMember(c.team) && !user.IsGlobalAdmin() {
		return nil, NewErrNotAllowed("must be global admin or a team member")
	}

	logger := log.WithFields(log.Fields{
		"cluster": name,
		"team":    c.team,
		"user":    user.Username(),
	})
	logger.Info("attempting to delete the cluster")

	original, err := c.Get(ctx, name)
	if err != nil {
		if err == ErrNotFound {
			return nil, err
		}

		logger.WithError(err).Error("failed to retrieve the cluster")

		return nil, err
	}

	return original, c.Store().Client().Delete(ctx, store.DeleteOptions.From(original))
}

// List returns a list of clusters we have access to
func (c *clustersImpl) List(ctx context.Context) (*clustersv1.ClusterList, error) {
	user := authentication.MustGetIdentity(ctx)
	if !user.IsMember(c.team) && !user.IsGlobalAdmin() {
		return nil, NewErrNotAllowed("must be global admin or a team member")
	}

	list := &clustersv1.ClusterList{}

	return list, c.Store().Client().List(ctx,
		store.ListOptions.InNamespace(c.team),
		store.ListOptions.InTo(list),
	)
}

// Get returns a specific cluster
func (c *clustersImpl) Get(ctx context.Context, name string) (*clustersv1.Cluster, error) {
	user := authentication.MustGetIdentity(ctx)
	if !user.IsMember(c.team) && !user.IsGlobalAdmin() {
		return nil, NewErrNotAllowed("must be global admin or a team member")
	}

	cluster := &clustersv1.Cluster{}

	if err := c.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(c.team),
		store.GetOptions.InTo(cluster),
		store.GetOptions.WithName(name),
	); err != nil {
		if kerrors.IsNotFound(err) {
			return nil, ErrNotFound
		}

		log.WithError(err).Error("failed to retrieve the cluster")
		return nil, err
	}

	return cluster, nil
}

// Update is used to update the cluster
func (c *clustersImpl) Update(ctx context.Context, cluster *clustersv1.Cluster) error {
	// @TODO check whether the user is an admin in the team

	user := authentication.MustGetIdentity(ctx)
	if !user.IsMember(c.team) && !user.IsGlobalAdmin() {
		return NewErrNotAllowed("must be global admin or a team member")
	}

	if cluster.Namespace == "" {
		cluster.Namespace = c.team
	}

	if cluster.Namespace != c.team {
		return validation.NewError("cluster has failed validation").WithFieldErrorf(
			"namespace",
			validation.MustExist,
			"must be the same as the team name: %q",
			c.team,
		)
	}

	if len(cluster.Name) > 40 {
		return validation.NewError("cluster has failed validation").
			WithFieldError("name", validation.MaxLength, "must be 40 characters or less")
	}

	if err := c.validateConfiguration(ctx, cluster); err != nil {
		return err
	}

	if err := c.validateCredentials(ctx, cluster); err != nil {
		return err
	}

	return c.Store().Client().Update(ctx,
		store.UpdateOptions.To(cluster),
		store.UpdateOptions.WithCreate(true),
	)
}

func (c *clustersImpl) validateCredentials(ctx context.Context, cluster *clustersv1.Cluster) error {
	creds := cluster.Spec.Credentials
	var alloc configv1.Allocation
	credentialAllocations, err := c.Teams().Team(c.team).Allocations().ListAllocationsByType(
		ctx, creds.Group, creds.Version, creds.Kind,
	)
	if err != nil {
		return err
	}
	for _, a := range credentialAllocations.Items {
		if a.Spec.Resource.Name == creds.Name {
			alloc = a
			break
		}
	}
	if alloc.Name == "" {
		return validation.NewError("cluster has failed validation").WithFieldErrorf(
			"credentials",
			validation.MustExist,
			"%q does not exist or it is not assigned to the team",
			creds.Name,
		)
	}

	expectedKind := fmt.Sprintf("%sCredentials", cluster.Spec.Kind)
	if !strings.EqualFold(alloc.Spec.Resource.Kind, expectedKind) {
		return validation.NewError("cluster has failed validation").WithFieldErrorf(
			"credentials",
			validation.InvalidType,
			"must be %q type",
			expectedKind,
		)
	}

	return nil
}

func (c *clustersImpl) validateConfiguration(ctx context.Context, cluster *clustersv1.Cluster) error {
	plan, err := c.plans.Get(ctx, cluster.Spec.Plan)
	if err != nil {
		if err == ErrNotFound {
			return validation.NewError("%q failed validation", cluster.Name).
				WithFieldErrorf("plan", validation.MustExist, "%q does not exist", cluster.Spec.Plan)
		}
		log.WithFields(log.Fields{
			"cluster": cluster.Name,
			"team":    c.team,
			"plan":    cluster.Spec.Plan,
		}).WithError(err).Error("failed to load plan")

		return err
	}

	planConfiguration := make(map[string]interface{})
	if err := json.NewDecoder(bytes.NewReader(plan.Spec.Configuration.Raw)).Decode(&planConfiguration); err != nil {
		return fmt.Errorf("failed to parse plan configuration values: %s", err)
	}

	clusterConfig := make(map[string]interface{})
	if err := json.NewDecoder(bytes.NewReader(cluster.Spec.Configuration.Raw)).Decode(&clusterConfig); err != nil {
		return fmt.Errorf("failed to parse cluster configuration values: %s", err)
	}

	switch cluster.Spec.Kind {
	case "GKE":
		if err := jsonschema.Validate(assets.GKEPlanSchema, "plan", clusterConfig); err != nil {
			return err
		}
	case "EKS":
		// TODO: add the EKS Plan schema and validate the plan parameters
	}

	editableParams, err := c.plans.GetEditablePlanParams(ctx, c.team)
	if err != nil {
		return err
	}

	verr := validation.NewError("%q failed validation", cluster.Name)

	for paramName, paramValue := range clusterConfig {
		if !reflect.DeepEqual(paramValue, planConfiguration[paramName]) {
			if !editableParams[paramName] {
				verr.AddFieldErrorf(paramName, validation.ReadOnly, "can not be changed")
			}
		}
	}
	if verr.HasErrors() {
		return verr
	}

	return nil
}
