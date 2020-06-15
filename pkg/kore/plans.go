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
	"fmt"
	"strings"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	"github.com/appvia/kore/pkg/kore/assets"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils/jsonschema"
	"github.com/appvia/kore/pkg/utils/validation"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

// Plans is the interface to the class plans
type Plans interface {
	// Delete is used to delete a plan in the kore
	Delete(context.Context, string) (*configv1.Plan, error)
	// Get returns the class from the kore
	Get(context.Context, string) (*configv1.Plan, error)
	// List returns a list of plans
	List(context.Context) (*configv1.PlanList, error)
	// Has checks if a resource exists within an available class in the scope
	Has(context.Context, string) (bool, error)
	// Update is responsible for update a plan in the kore
	Update(ctx context.Context, plan *configv1.Plan, ignoreReadonly bool) error
	// GetEditablePlanParams returns with the editable plan parameters for a specific team and cluster kind
	GetEditablePlanParams(ctx context.Context, team string, clusterKind string) ([]string, error)
}

type plansImpl struct {
	Interface
}

// Update is responsible for update a plan in the kore
func (p plansImpl) Update(ctx context.Context, plan *configv1.Plan, ignoreReadonly bool) error {
	plan.Namespace = HubNamespace

	user := authentication.MustGetIdentity(ctx)
	if !user.IsGlobalAdmin() {
		log.WithField("user", user.Username()).Warn("trying to update a plan without permissions")

		return ErrUnauthorized
	}

	if !ignoreReadonly {
		original, err := p.Get(ctx, plan.Name)
		if err != nil && err != ErrNotFound {
			return err
		}

		if original != nil && original.Annotations[AnnotationReadOnly] == AnnotationValueTrue {
			return validation.NewError("the plan can not be updated").WithFieldError(validation.FieldRoot, validation.ReadOnly, "plan is read-only")
		}
		if plan.Annotations[AnnotationReadOnly] == AnnotationValueTrue {
			return validation.NewError("the plan can not be updated").WithFieldError(validation.FieldRoot, validation.ReadOnly, "read-only flag can not be set")
		}
	}

	schema, err := assets.GetClusterSchema(plan.Spec.Kind)
	if err != nil {
		return validation.NewError("cluster failed validation").
			WithFieldError("kind", validation.InvalidType, err.Error())
	}

	if err := jsonschema.Validate(schema, "plan", plan.Spec.Configuration); err != nil {
		return err
	}

	err = p.Store().Client().Update(ctx,
		store.UpdateOptions.To(plan),
		store.UpdateOptions.WithCreate(true),
		store.UpdateOptions.WithForce(true),
	)
	if err != nil {
		log.WithError(err).Error("trying to update a plan in the kore")

		return err
	}

	return nil
}

// Delete is used to delete a plan in the kore
func (p plansImpl) Delete(ctx context.Context, name string) (*configv1.Plan, error) {
	user := authentication.MustGetIdentity(ctx)
	if !user.IsGlobalAdmin() {
		log.WithField("user", user.Username()).Warn("trying to delete a plan without permission")

		return nil, ErrUnauthorized
	}

	plan := &configv1.Plan{}
	err := p.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(HubNamespace),
		store.GetOptions.InTo(plan),
		store.GetOptions.WithName(name),
	)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		log.WithError(err).Error("trying to retrieve plan in the kore")

		return nil, err
	}

	if plan.Annotations[AnnotationReadOnly] == AnnotationValueTrue {
		return nil, validation.NewError("the plan can not be deleted").WithFieldError(validation.FieldRoot, validation.ReadOnly, "plan is read-only")
	}

	clustersWithPlan, err := p.getClustersWithPlan(ctx, name)
	if err != nil {
		return nil, err
	}
	if len(clustersWithPlan) > 0 {
		if len(clustersWithPlan) <= 5 {
			return nil, fmt.Errorf(
				"the plan can not be deleted as there are %d clusters using it: %s",
				len(clustersWithPlan),
				strings.Join(clustersWithPlan, ", "),
			)
		}
		return nil, fmt.Errorf(
			"the plan can not be deleted as there are %d clusters using it",
			len(clustersWithPlan),
		)
	}

	return plan, p.Store().Client().Delete(ctx, store.DeleteOptions.From(plan))
}

// Get returns the class from the kore
func (p plansImpl) Get(ctx context.Context, name string) (*configv1.Plan, error) {
	plan := &configv1.Plan{}

	if found, err := p.Has(ctx, name); err != nil {
		return nil, err
	} else if !found {
		return nil, ErrNotFound
	}

	return plan, p.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(HubNamespace),
		store.GetOptions.WithName(name),
		store.GetOptions.InTo(plan),
	)
}

// List returns a list of classes
func (p plansImpl) List(ctx context.Context) (*configv1.PlanList, error) {
	planList := &configv1.PlanList{}

	err := p.Store().Client().List(ctx,
		store.ListOptions.InNamespace(HubNamespace),
		store.ListOptions.InTo(planList),
	)
	if err != nil {
		return nil, err
	}

	res := []configv1.Plan{}
	for _, plan := range planList.Items {
		if plan.Annotations[AnnotationSystem] == AnnotationValueTrue {
			continue
		}
		res = append(res, plan)
	}
	planList.Items = res

	return planList, nil
}

// Has checks if a resource exists within an available class in the scope
func (p plansImpl) Has(ctx context.Context, name string) (bool, error) {
	return p.Store().Client().Has(ctx,
		store.HasOptions.InNamespace(HubNamespace),
		store.HasOptions.From(&configv1.Plan{}),
		store.HasOptions.WithName(name),
	)
}

func (p plansImpl) GetEditablePlanParams(ctx context.Context, team string, clusterKind string) ([]string, error) {
	editableParams := map[string]bool{}
	planPolicyAllocations, err := p.Teams().Team(team).Allocations().ListAllocationsByType(
		ctx, "config.kore.appvia.io", "v1", "PlanPolicy",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load plan policies assigned to the team: %s", err)
	}

	for _, alloc := range planPolicyAllocations.Items {
		planPolicy, err := p.PlanPolicies().Get(ctx, alloc.Spec.Resource.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to load plan policy: %s", alloc.Spec.Resource.Name)
		}
		if planPolicy.Spec.Kind != clusterKind {
			continue
		}
		for _, property := range planPolicy.Spec.Properties {
			switch {
			case property.DisallowUpdate:
				editableParams[property.Name] = false
			case property.AllowUpdate:
				if _, isSet := editableParams[property.Name]; !isSet {
					editableParams[property.Name] = true
				}
			}
		}
	}

	var res []string
	for p, editable := range editableParams {
		if editable {
			res = append(res, p)
		}
	}
	return res, nil
}

func (p plansImpl) getClustersWithPlan(ctx context.Context, clusterName string) ([]string, error) {
	var res []string

	teamList, err := p.Teams().List(ctx)
	if err != nil {
		return nil, err
	}

	for _, team := range teamList.Items {
		clusterList, err := p.Teams().Team(team.Name).Clusters().List(ctx)
		if err != nil {
			return nil, err
		}
		for _, cluster := range clusterList.Items {
			if cluster.Spec.Plan == clusterName {
				res = append(res, fmt.Sprintf("%s/%s", team.Name, cluster.Name))
			}
		}
	}

	return res, nil
}
