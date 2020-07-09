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
	"strings"
	"text/template"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/jsonschema"
	"github.com/appvia/kore/pkg/utils/jsonutils"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	"github.com/appvia/kore/pkg/utils/validation"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

// ServicePlanDetails contains information about a service plan in the given team/cluster etc. context
type ServicePlanDetails struct {
	servicesv1.ServicePlanSpec `json:",inline"`
	EditableParams             []string `json:"editableParams"`
}

// ServicePlans is the interface to manage service plans
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ServicePlans
type ServicePlans interface {
	// CheckDelete verifies whether the service plan can be deleted
	CheckDelete(context.Context, *servicesv1.ServicePlan, ...DeleteOptionFunc) error
	// Delete is used to delete a service plan in the kore
	Delete(context.Context, string, ...DeleteOptionFunc) (*servicesv1.ServicePlan, error)
	// Get returns the service plan
	Get(context.Context, string) (*servicesv1.ServicePlan, error)
	// GetDetails returns information about a service plan in the given team/cluster etc. context
	GetDetails(ctx context.Context, name, team, clusterName string) (ServicePlanDetails, error)
	// List returns the existing service plans
	// The optional filter functions can be used to include items only for which all functions return true
	List(context.Context, ...func(plan servicesv1.ServicePlan) bool) (*servicesv1.ServicePlanList, error)
	// Has checks if a service plan exists
	Has(context.Context, string) (bool, error)
	// Update is responsible for updating a service plan
	Update(ctx context.Context, plan *servicesv1.ServicePlan, ignoreReadonly bool) error
}

type servicePlansImpl struct {
	Interface
}

// Update is responsible for updating a service plan
func (p servicePlansImpl) Update(ctx context.Context, plan *servicesv1.ServicePlan, ignoreReadonly bool) error {
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

	existing, err := p.Get(ctx, plan.Name)
	if err != nil && err != ErrNotFound {
		return fmt.Errorf("failed to get plan %q: %w", plan.Name, err)
	}

	if !ignoreReadonly {
		if existing != nil && existing.Annotations[AnnotationReadOnly] == AnnotationValueTrue {
			return validation.NewError("the plan can not be updated").WithFieldError(validation.FieldRoot, validation.ReadOnly, "plan is read-only")
		}
		if plan.Annotations[AnnotationReadOnly] == AnnotationValueTrue {
			return validation.NewError("the plan can not be updated").WithFieldError(validation.FieldRoot, validation.ReadOnly, "read-only flag can not be set")
		}
	}

	if existing != nil {
		verr := validation.NewError("%q failed validation", plan.Name)
		if existing.Spec.Kind != plan.Spec.Kind {
			verr.AddFieldErrorf("kind", validation.ReadOnly, "can not be changed after the service plan was created")
		}
		if verr.HasErrors() {
			return verr
		}
	}

	kind, err := p.ServiceKinds().Get(ctx, plan.Spec.Kind)
	if err != nil {
		return fmt.Errorf("failed to get service kind %q: %w", plan.Spec.Kind, err)
	}

	schema := plan.Spec.Schema
	if schema == "" {
		schema = kind.Spec.Schema
	}

	if schema == "" && !utils.ApiExtJSONEmpty(plan.Spec.Configuration) {
		if existing == nil || !utils.ApiExtJSONEquals(plan.Spec.Configuration, existing.Spec.Configuration) {
			return validation.NewError("%q failed validation", plan.Name).
				WithFieldErrorf(
					"configuration",
					validation.ReadOnly,
					"the service provider doesn't have a JSON schema to validate the configuration",
				)
		}
	}

	if schema != "" {
		if err := jsonschema.Validate(schema, "plan", plan.Spec.Configuration); err != nil {
			return err
		}
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

// CheckDelete verifies whether the service plan can be deleted
func (p servicePlansImpl) CheckDelete(ctx context.Context, servicePlan *servicesv1.ServicePlan, o ...DeleteOptionFunc) error {
	opts := ResolveDeleteOptions(o)

	if !opts.Cascade {
		var dependents []kubernetes.DependentReference

		teamList, err := p.Teams().List(ctx)
		if err != nil {
			return fmt.Errorf("failed to list teams: %w", err)
		}

		for _, team := range teamList.Items {
			services, err := p.Teams().Team(team.Name).Services().List(ctx, func(s servicesv1.Service) bool {
				return s.Spec.Plan == servicePlan.Name
			})
			if err != nil {
				return fmt.Errorf("failed to list services: %w", err)
			}
			for _, item := range services.Items {
				dependents = append(dependents, kubernetes.DependentReferenceFromObject(&item))
			}
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

// Delete is used to delete a service plan in the kore
func (p servicePlansImpl) Delete(ctx context.Context, name string, o ...DeleteOptionFunc) (*servicesv1.ServicePlan, error) {
	opts := ResolveDeleteOptions(o)

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

	if err := opts.Check(plan, func(o ...DeleteOptionFunc) error { return p.CheckDelete(ctx, plan, o...) }); err != nil {
		return nil, err
	}

	if err := p.Store().Client().Delete(ctx, append(opts.StoreOptions(), store.DeleteOptions.From(plan))...); err != nil {
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
func (p servicePlansImpl) List(ctx context.Context, filters ...func(servicesv1.ServicePlan) bool) (*servicesv1.ServicePlanList, error) {
	list := &servicesv1.ServicePlanList{}

	err := p.Store().Client().List(ctx,
		store.ListOptions.InNamespace(HubNamespace),
		store.ListOptions.InTo(list),
	)
	if err != nil {
		return nil, err
	}

	if len(filters) == 0 {
		return list, nil
	}

	res := []servicesv1.ServicePlan{}
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

// Has checks if a service plan exists
func (p servicePlansImpl) Has(ctx context.Context, name string) (bool, error) {
	return p.Store().Client().Has(ctx,
		store.HasOptions.InNamespace(HubNamespace),
		store.HasOptions.From(&servicesv1.ServicePlan{}),
		store.HasOptions.WithName(name),
	)
}

// GetDetails returns information about a service plan in the given team/cluster etc. context
func (p servicePlansImpl) GetDetails(ctx context.Context, name, team, clusterName string) (ServicePlanDetails, error) {
	plan, err := p.Get(ctx, name)
	if err != nil {
		return ServicePlanDetails{}, err
	}

	kind, err := p.ServiceKinds().Get(ctx, plan.Spec.Kind)
	if err != nil {
		return ServicePlanDetails{}, err
	}

	if plan.Spec.Schema == "" {
		plan.Spec.Schema = kind.Spec.Schema
	}

	if plan.Spec.CredentialSchema == "" {
		plan.Spec.CredentialSchema = kind.Spec.CredentialSchema
	}

	if team != "" && clusterName != "" {
		cluster, err := p.Teams().Team(team).Clusters().Get(ctx, clusterName)
		if err != nil {
			return ServicePlanDetails{}, err
		}

		if plan.Spec.Schema != "" {
			var err error
			if plan.Spec.Schema, err = p.compileTemplate(plan.Spec.Schema, cluster); err != nil {
				return ServicePlanDetails{}, err
			}
		}

		if plan.Spec.CredentialSchema != "" {
			var err error
			if plan.Spec.CredentialSchema, err = p.compileTemplate(plan.Spec.CredentialSchema, cluster); err != nil {
				return ServicePlanDetails{}, err
			}
		}

		if plan.Spec.Configuration != nil {
			compiledConfigBytes, err := p.compileTemplate(string(plan.Spec.Configuration.Raw), cluster)
			if err != nil {
				return ServicePlanDetails{}, err
			}
			plan.Spec.Configuration.Raw = []byte(compiledConfigBytes)
		}
	}

	var editableParams []string
	if plan.Spec.Schema != "" {
		schema := &jsonschema.Schema{}
		if err := json.Unmarshal([]byte(plan.Spec.Schema), schema); err != nil {
			return ServicePlanDetails{}, err
		}
		for name, prop := range schema.Properties {
			if prop.Const != nil {
				continue
			}

			editableParams = append(editableParams, name)
		}
	}

	return ServicePlanDetails{
		ServicePlanSpec: plan.Spec,
		EditableParams:  editableParams,
	}, nil
}

// GetSchemaForCluster returns the service plan schema generated for the given cluster
func (p servicePlansImpl) compileTemplate(content string, cluster *clustersv1.Cluster) (string, error) {
	tmpl, err := template.New("content").Parse(content)
	if err != nil {
		return "", err
	}
	tmpl.Option("missingkey=error")

	clusterObj, err := jsonutils.ToMap(cluster)
	if err != nil {
		return "", fmt.Errorf("failed to encode cluster: %w", err)
	}

	tmplBuf := bytes.NewBuffer(make([]byte, 0, 16384))
	params := map[string]interface{}{
		"cluster": clusterObj,
	}
	if err := tmpl.Execute(tmplBuf, params); err != nil {
		return "", err
	}

	return tmplBuf.String(), nil
}
