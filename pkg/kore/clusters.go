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
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/jsonschema"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	"github.com/appvia/kore/pkg/utils/validation"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

const (
	// ClusterKindKore is the default Cluster Kind for the cluster Kore is running in
	ClusterKindKore = "Kore"
)

// Clusters returns the an interface for handling clusters
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Clusters
type Clusters interface {
	// CheckDelete verifies whether the cluster can be deleted
	CheckDelete(context.Context, *clustersv1.Cluster, ...DeleteOptionFunc) error
	// Delete is used to delete a cluster
	Delete(context.Context, string, ...DeleteOptionFunc) (*clustersv1.Cluster, error)
	// Get returns a specific cluster
	Get(context.Context, string) (*clustersv1.Cluster, error)
	// List returns a list of clusters we have access to
	// The optional filter functions can be used to include items only for which all functions return true
	List(context.Context, ...func(clustersv1.Cluster) bool) (*clustersv1.ClusterList, error)
	// Update is used to update the cluster
	Update(context.Context, *clustersv1.Cluster) error
}

type clustersImpl struct {
	*hubImpl
	// team is the name
	team string
}

// CheckDelete verifies whether the cluster can be deleted
func (c *clustersImpl) CheckDelete(ctx context.Context, cluster *clustersv1.Cluster, o ...DeleteOptionFunc) error {
	opts := ResolveDeleteOptions(o)

	if !opts.Cascade {
		var dependents []kubernetes.DependentReference
		services, err := c.Teams().Team(c.team).Services().List(ctx, func(s servicesv1.Service) bool { return kubernetes.HasOwnerReference(&s, cluster) })
		if err != nil {
			return fmt.Errorf("failed to list services: %w", err)
		}
		for _, item := range services.Items {
			dependents = append(dependents, kubernetes.DependentReferenceFromObject(&item))
		}

		namespaceClaims, err := c.Teams().Team(c.team).NamespaceClaims().List(ctx, func(nc clustersv1.NamespaceClaim) bool { return kubernetes.HasOwnerReference(&nc, cluster) })
		if err != nil {
			return fmt.Errorf("failed to list namespace claims: %w", err)
		}
		for _, item := range namespaceClaims.Items {
			dependents = append(dependents, kubernetes.DependentReferenceFromObject(&item))
		}

		if len(dependents) > 0 {
			return validation.ErrDependencyViolation{
				Message:    "the following objects need to be deleted first",
				Dependents: dependents,
			}
		}
	}

	return nil
}

// Delete is used to delete a cluster
func (c *clustersImpl) Delete(ctx context.Context, name string, o ...DeleteOptionFunc) (*clustersv1.Cluster, error) {
	opts := ResolveDeleteOptions(o)

	// @TODO check whether the user is an admin in the team

	logger := log.WithFields(log.Fields{
		"cluster": name,
		"team":    c.team,
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

	if err := opts.Check(original, func(o ...DeleteOptionFunc) error { return c.CheckDelete(ctx, original, o...) }); err != nil {
		return nil, err
	}

	if original.Labels[LabelClusterIdentifier] != "" {
		if err := c.hubImpl.Teams().Team(c.team).Assets().MarkAssetDeleted(ctx, original.Labels[LabelClusterIdentifier]); err != nil {
			return nil, fmt.Errorf("error marking asset as deleted: %w", err)
		}
	}

	return original, c.Store().Client().Delete(ctx, append(opts.StoreOptions(), store.DeleteOptions.From(original))...)
}

// List returns a list of clusters we have access to
func (c *clustersImpl) List(ctx context.Context, filters ...func(clustersv1.Cluster) bool) (*clustersv1.ClusterList, error) {
	list := &clustersv1.ClusterList{}

	err := c.Store().Client().List(ctx,
		store.ListOptions.InNamespace(c.team),
		store.ListOptions.InTo(list),
	)
	if err != nil {
		return nil, err
	}

	if len(filters) == 0 {
		return list, nil
	}

	res := []clustersv1.Cluster{}
	for _, item := range list.Items {
		if func() bool {
			for _, filter := range filters {
				if !filter(item) {
					return false
				}
			}
			return true
		}() {
			res = append(res, item)
		}
	}
	list.Items = res

	return list, nil
}

// Get returns a specific cluster
func (c *clustersImpl) Get(ctx context.Context, name string) (*clustersv1.Cluster, error) {
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
	provider, ok := GetClusterProvider(cluster.Spec.Kind)
	if !ok {
		return validation.NewError("cluster has failed validation").
			WithFieldError("kind", validation.InvalidValue, "not supported")
	}

	existing, err := c.Get(ctx, cluster.Name)
	if err != nil && err != ErrNotFound {
		return err
	}

	if existing != nil && existing.Annotations[AnnotationReadOnly] == AnnotationValueTrue {
		return validation.NewError("the cluster can not be updated").WithFieldError(validation.FieldRoot, validation.ReadOnly, "cluster is read-only")
	}
	if cluster.Annotations[AnnotationReadOnly] == AnnotationValueTrue {
		return validation.NewError("the cluster can not be updated").WithFieldError(validation.FieldRoot, validation.ReadOnly, "read-only flag can not be set")
	}

	if existing != nil {
		verr := validation.NewError("cluster has failed validation")
		if existing.Spec.Kind != cluster.Spec.Kind {
			verr.AddFieldErrorf("kind", validation.ReadOnly, "can not be changed after a cluster was created")
		}
		if existing.Spec.Plan != cluster.Spec.Plan {
			verr.AddFieldErrorf("plan", validation.ReadOnly, "can not be changed after a cluster was created")
		}
		if existing.Labels[LabelClusterIdentifier] != cluster.Labels[LabelClusterIdentifier] {
			verr.AddFieldErrorf(LabelClusterIdentifier, validation.ReadOnly, "assigned by Kore, keep existing value")
		}
		if existing.Labels[LabelTeamIdentifier] != cluster.Labels[LabelTeamIdentifier] {
			verr.AddFieldErrorf(LabelTeamIdentifier, validation.ReadOnly, "assigned by Kore, keep existing value")
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

	if err := c.validateConfiguration(ctx, provider, cluster, existing); err != nil {
		return err
	}

	if err := c.validateCredentials(ctx, cluster); err != nil {
		return err
	}

	if err := c.validateAccounting(ctx, cluster); err != nil {
		return err
	}

	identifierAssigned, err := c.ensureIdentifiers(ctx, cluster, existing == nil)
	if err != nil {
		return err
	}

	err = c.Store().Client().Update(ctx,
		store.UpdateOptions.To(cluster),
		store.UpdateOptions.WithCreate(true),
	)
	if err != nil && identifierAssigned {
		log.WithError(err).Warn("error occurred persisting update to cluster but an identifier has been assigned so marking as deleted")
		delErr := c.hubImpl.Teams().Team(c.team).Assets().MarkAssetDeleted(ctx, cluster.Labels[LabelClusterIdentifier])
		if delErr != nil {
			// warn that this has happened, but let the causing error flow up - this is
			// not a major problem and the causal error is more important to the caller
			// than this one.
			log.WithError(delErr).Warn("error marking asset as deleted after failure to persist update to cluster")
		}
	}
	return err
}

// ensureIdentifiers checks the team and cluster identifier labels are set, returning true if a cluster identifier has
// been assigned (or re-assigned)
func (c *clustersImpl) ensureIdentifiers(ctx context.Context, cluster *clustersv1.Cluster, isNewCluster bool) (bool, error) {
	if cluster.Labels == nil {
		cluster.Labels = map[string]string{}
	}
	assets := c.hubImpl.Teams().Team(c.team).Assets()

	// @step: Validate or ensure team identifier
	var err error
	valid := true
	if cluster.Labels[LabelTeamIdentifier] != "" {
		valid, err = assets.ValidateTeamIdentifier(ctx, cluster.Labels[LabelTeamIdentifier])
	} else {
		cluster.Labels[LabelTeamIdentifier], err = assets.EnsureTeamIdentifier(ctx)
	}
	if err != nil {
		return false, err
	}
	if !valid {
		return false, validation.NewError("cluster has failed validation").
			WithFieldError(LabelTeamIdentifier, validation.InvalidValue, "leave blank to have kore auto-populate or set to correct team ID")
	}

	if !isNewCluster {
		if cluster.Labels[LabelClusterIdentifier] == "" {
			// should never get here, but defensively, return a validation error
			return false, validation.NewError("cluster has failed validation").
				WithFieldError(LabelClusterIdentifier, validation.ReadOnly, "assigned by Kore, this should be set to the existing value when updating a cluster")
		}
		return false, nil
	}

	// @step: For a new cluster if an identifier has been supplied, check valid for re-use and mark it as
	// active again
	if cluster.Labels[LabelClusterIdentifier] != "" {
		valid, err := assets.ReuseAssetIdentifier(ctx, cluster.Labels[LabelClusterIdentifier], orgv1.TeamAssetTypeCluster, cluster.Name)
		if err != nil {
			return false, err
		}
		if !valid {
			return false, validation.NewError("cluster has failed validation").
				WithFieldError(LabelClusterIdentifier, validation.InvalidValue, "assigned by Kore, leave blank or set to a previous cluster ID owned by your team for new cluster")
		}
		return true, nil
	}

	// @step: Assign an identifier to this new cluster
	cluster.Labels[LabelClusterIdentifier], err = assets.GenerateAssetIdentifier(ctx, orgv1.TeamAssetTypeCluster, cluster.Name)
	if err != nil {
		return false, err
	}
	return true, nil
}

// validateAccounting is responsible for checking if accounting
func (c *clustersImpl) validateAccounting(ctx context.Context, cluster *clustersv1.Cluster) error {
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

	// @step: does this team have accounts enabled
	permitted, err := c.Teams().Team(c.team).Allocations().IsPermitted(ctx, cluster.Spec.Credentials)
	if err != nil {
		return err
	}
	if !permitted {
		return ErrNotAllowed{message: "account management is not allocated to the team"}
	}

	// @step: is the plan requested available in the account plan?
	account, err := c.Accounts().Get(ctx, cluster.Spec.Credentials.Name)
	if err != nil {
		return err
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
		supported := strings.Join(list, ",")

		return validation.NewError("cluster failed validation, plan not part of accounting rules").
			WithFieldError("plan", validation.InvalidValue, "plan not included the accounting rules (supported: "+supported+")")
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

	supported := []string{"AKSCredentials", "EKSCredentials", "GKECredentials", "AccountManagement"}
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

func (c *clustersImpl) validateConfiguration(
	ctx context.Context, provider ClusterProvider, cluster, existing *clustersv1.Cluster,
) error {
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

	if err := jsonschema.Validate(provider.PlanJSONSchema(), cluster.Name, clusterConfig); err != nil {
		return err
	}

	if err := provider.Validate(NewContext(ctx, log.StandardLogger(), c.Store().RuntimeClient(), c), cluster); err != nil {
		return err
	}

	editableParams, err := c.plans.GetEditablePlanParams(ctx, c.team, cluster.Spec.Kind)
	if err != nil {
		return err
	}

	verr := validation.NewError("%q failed validation", cluster.Name)

	for paramName, paramValue := range clusterConfig {
		if !reflect.DeepEqual(paramValue, planConfiguration[paramName]) {
			if !utils.Contains(paramName, editableParams) {
				verr.AddFieldErrorf(paramName, validation.ReadOnly, "can not be changed")
			}
		}
	}

	if existing != nil {
		if err := jsonschema.ValidateImmutableProperties(
			provider.PlanJSONSchema(),
			"cluster",
			"spec.configuration",
			existing.Spec.Configuration.Raw,
			cluster.Spec.Configuration.Raw,
		); err != nil {
			return err
		}
	}

	if verr.HasErrors() {
		return verr
	}

	return nil
}
