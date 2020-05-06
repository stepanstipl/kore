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

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils/jsonschema"
	"github.com/appvia/kore/pkg/utils/validation"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

// ServicePlans is the interface to manage service plans
type ServicePlans interface {
	// Delete is used to delete a service plan in the kore
	Delete(context.Context, string) (*servicesv1.ServicePlan, error)
	// Get returns the service plan
	Get(context.Context, string) (*servicesv1.ServicePlan, error)
	// List returns the existing service plans
	List(context.Context) (*servicesv1.ServicePlanList, error)
	// ListFiltered returns a list of service plans using the given filter.
	ListFiltered(context.Context, func(servicesv1.ServicePlan) bool) (*servicesv1.ServicePlanList, error)
	// Has checks if a service plan exists
	Has(context.Context, string) (bool, error)
	// Update is responsible for updating a service plan
	Update(context.Context, *servicesv1.ServicePlan) error
	// GetEditablePlanParams returns with the editable service plan parameters for a specific team and service kind
	GetEditablePlanParams(ctx context.Context, team string, clusterKind string) (map[string]bool, error)
}

type servicePlansImpl struct {
	Interface
}

// Update is responsible for updating a service plan
func (p servicePlansImpl) Update(ctx context.Context, plan *servicesv1.ServicePlan) error {
	if err := IsValidResourceName("plan", plan.Name); err != nil {
		return err
	}

	if !strings.HasPrefix(plan.Name, plan.Spec.Kind+"-") {
		return validation.NewError("%q failed validation", plan.Name).
			WithFieldErrorf("name", validation.InvalidValue, "must start with %s-", plan.Spec.Kind)
	}

	if plan.Namespace != HubNamespace {
		return validation.NewError("%q failed validation", plan.Name).
			WithFieldErrorf("namespace", validation.InvalidValue, "must be %q", HubNamespace)
	}

	provider := p.ServiceProviders().GetProviderForKind(plan.Spec.Kind)
	if provider == nil {
		return validation.NewError("%q failed validation", plan.Name).
			WithFieldErrorf("kind", validation.InvalidType, "%q is not a known service kind", plan.Spec.Kind)
	}

	schema, err := provider.PlanJSONSchema(plan.Spec.Kind, plan.PlanShortName())
	if err != nil {
		return err
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
		log.WithError(err).Error("failed to update a service plan in the kore")

		return err
	}

	return nil
}

// Delete is used to delete a service plan in the kore
func (p servicePlansImpl) Delete(ctx context.Context, name string) (*servicesv1.ServicePlan, error) {
	plan := &servicesv1.ServicePlan{}
	err := p.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(HubNamespace),
		store.GetOptions.InTo(plan),
		store.GetOptions.WithName(name),
	)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		log.WithError(err).Error("failed to retrieve the service plan")

		return nil, err
	}

	servicesWithPlan, err := p.getServicesWithPlan(ctx, name)
	if err != nil {
		return nil, err
	}
	if len(servicesWithPlan) > 0 {
		if len(servicesWithPlan) <= 5 {
			return nil, fmt.Errorf(
				"the service plan can not be deleted as there are %d services using it: %s",
				len(servicesWithPlan),
				strings.Join(servicesWithPlan, ", "),
			)
		}
		return nil, fmt.Errorf(
			"the service plan can not be deleted as there are %d services using it",
			len(servicesWithPlan),
		)
	}

	if err := p.Store().Client().Delete(ctx, store.DeleteOptions.From(plan)); err != nil {
		log.WithError(err).Error("failed to delete the service plan")

		return nil, err
	}

	return plan, nil
}

// Get returns the service plan
func (p servicePlansImpl) Get(ctx context.Context, name string) (*servicesv1.ServicePlan, error) {
	plan := &servicesv1.ServicePlan{}

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

// List returns the existing service plans
func (p servicePlansImpl) List(ctx context.Context) (*servicesv1.ServicePlanList, error) {
	plans := &servicesv1.ServicePlanList{}

	return plans, p.Store().Client().List(ctx,
		store.ListOptions.InNamespace(HubNamespace),
		store.ListOptions.InTo(plans),
	)
}

// ListFiltered returns a list of service plans using the given filter.
// A service plan is included if the filter function returns true
func (p servicePlansImpl) ListFiltered(ctx context.Context, filter func(plan servicesv1.ServicePlan) bool) (*servicesv1.ServicePlanList, error) {
	var res []servicesv1.ServicePlan

	servicePlansList, err := p.ServicePlans().List(ctx)
	if err != nil {
		return nil, err
	}

	for _, servicePlan := range servicePlansList.Items {
		if filter(servicePlan) {
			res = append(res, servicePlan)
		}
	}

	servicePlansList.Items = res

	return servicePlansList, nil
}

// Has checks if a service plan exists
func (p servicePlansImpl) Has(ctx context.Context, name string) (bool, error) {
	return p.Store().Client().Has(ctx,
		store.HasOptions.InNamespace(HubNamespace),
		store.HasOptions.From(&servicesv1.ServicePlan{}),
		store.HasOptions.WithName(name),
	)
}

// GetEditablePlanParams returns with the editable service plan parameters for a specific team and service kind
func (p servicePlansImpl) GetEditablePlanParams(ctx context.Context, team string, clusterKind string) (map[string]bool, error) {
	// TODO: implement this when the service plan policies are implemented
	return nil, nil
}

func (p servicePlansImpl) getServicesWithPlan(ctx context.Context, clusterName string) ([]string, error) {
	var res []string

	teamList, err := p.Teams().List(ctx)
	if err != nil {
		return nil, err
	}

	for _, team := range teamList.Items {
		servicesList, err := p.Teams().Team(team.Name).Services().List(ctx)
		if err != nil {
			return nil, err
		}
		for _, service := range servicesList.Items {
			if service.Spec.Plan == clusterName {
				res = append(res, fmt.Sprintf("%s/%s", team.Name, service.Name))
			}
		}
	}

	return res, nil
}
