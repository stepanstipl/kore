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

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	"github.com/appvia/kore/pkg/kore/assets"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/jsonschema"
	"github.com/appvia/kore/pkg/utils/validation"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
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

	delErr := c.Store().Client().Delete(ctx, store.DeleteOptions.From(original))

	if delErr == nil {
		if err := c.Security().ArchiveResourceScans(ctx, original.TypeMeta, original.ObjectMeta); err != nil {
			// Log but continue in case of errors here - the cluster IS deleted.
			log.WithError(err).Warning("error while archiving security scans for cluster")
		}
	}

	return original, delErr
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
		log.WithError(err).Error("trying to retrieve the cluster")

		return nil, err
	}

	return cluster, nil
}

// Update is used to update the cluster
func (c *clustersImpl) Update(ctx context.Context, cluster *clustersv1.Cluster) error {
	user := authentication.MustGetIdentity(ctx)
	if !user.IsMember(c.team) && !user.IsGlobalAdmin() {
		return NewErrNotAllowed("must be global admin or a team member")
	}

	existing, err := c.Get(ctx, cluster.Name)
	if err != nil && err != ErrNotFound {
		return err
	}

	if existing != nil {
		verr := validation.NewError("cluster has failed validation")
		if existing.Spec.Kind != cluster.Spec.Kind {
			verr.AddFieldErrorf("kind", validation.ReadOnly, "can not be changed after a cluster was created")
		}
		if existing.Spec.Plan != cluster.Spec.Plan {
			verr.AddFieldErrorf("plan", validation.ReadOnly, "can not be changed after a cluster was created")
		}
		if verr.HasErrors() {
			return verr
		}
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

	if err := c.validateAccounting(ctx, cluster); err != nil {
		return err
	}

	return c.Store().Client().Update(ctx,
		store.UpdateOptions.To(cluster),
		store.UpdateOptions.WithCreate(true),
	)
}

// validateAccounting is responsible for checking if accounting
func (c *clustersImpl) validateAccounting(ctx context.Context, cluster *clustersv1.Cluster) error {
	fmt.Println("KIND", cluster.Spec.Credentials.Kind)

	if cluster.Spec.Credentials.Kind != "AccountManagement" {
		return nil
	}

	// @choice: if the cluster already exist we should bypass this check
	if _, err := c.Teams().Team(c.team).Clusters().Get(ctx, cluster.Name); err != nil {
		if err != ErrNotFound {
			return err
		}
	} else {
		return nil
	}

	// @step: does this team having accounts enabled
	permitted, err := c.Teams().Team(c.team).Allocations().IsPermitted(ctx, cluster.Spec.Credentials)
	if err != nil {
		return err
	}
	if !permitted {
		return ErrNotAllowed{message: "account management is not allocated to the team"}
	}

	// @step: if the plan requested available in the account plan?
	account, err := c.Accounts().Get(ctx, cluster.Spec.Credentials.Name)
	if err != nil {
		return err
	}

	if len(account.Spec.Rules) <= 0 {
		return nil
	}

	found, list := func() (bool, []string) {
		var list []string

		for _, rule := range account.Spec.Rules {
			list = append(list, rule.Plans...)
			if utils.Contains(cluster.Spec.Plan, rule.Plans) {
				return true, nil
			}
		}

		return false, list
	}()
	if !found {
		return ErrNotAllowed{
			message: fmt.Sprintf("Plan is not permitted by accounting rules (allowed: %s)", strings.Join(list, ",")),
		}
	}

	return nil
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

	supported := []string{"EKSCredentials", "GKECredentials", "AccountManagement"}
	if !utils.Contains(alloc.Spec.Resource.Kind, supported) {
		return validation.NewError("cluster has failed validation").WithFieldErrorf(
			"credentials",
			validation.InvalidType,
			"must be %q type",
			strings.Join(supported, ","),
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
		if err := jsonschema.Validate(assets.GKEPlanSchema, cluster.Name, clusterConfig); err != nil {
			return err
		}
	case "EKS":
		if err := jsonschema.Validate(assets.EKSPlanSchema, cluster.Name, clusterConfig); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid cluster kind: %q", cluster.Spec.Kind)
	}

	editableParams, err := c.plans.GetEditablePlanParams(ctx, c.team, cluster.Spec.Kind)
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
