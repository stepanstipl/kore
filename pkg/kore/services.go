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

	"github.com/appvia/kore/pkg/utils"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"

	"github.com/appvia/kore/pkg/utils/jsonschema"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils/validation"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

// Services returns the an interface for handling services
type Services interface {
	// Delete is used to delete a service
	Delete(context.Context, string) (*servicesv1.Service, error)
	// Get returns a specific service
	Get(context.Context, string) (*servicesv1.Service, error)
	// List returns a list of services
	List(context.Context) (*servicesv1.ServiceList, error)
	// Update is used to update a service
	Update(context.Context, *servicesv1.Service) error
}

type servicesImpl struct {
	*hubImpl
	// team is the name
	team string
}

// Delete is used to delete a service
func (s *servicesImpl) Delete(ctx context.Context, name string) (*servicesv1.Service, error) {
	// @TODO check whether the user is an admin in the team

	user := authentication.MustGetIdentity(ctx)
	if !user.IsMember(s.team) && !user.IsGlobalAdmin() {
		return nil, NewErrNotAllowed("must be global admin or a team member")
	}

	logger := log.WithFields(log.Fields{
		"service": name,
		"team":    s.team,
		"user":    user.Username(),
	})
	logger.Info("attempting to delete the service")

	original, err := s.Get(ctx, name)
	if err != nil {
		if err == ErrNotFound {
			return nil, err
		}

		logger.WithError(err).Error("failed to retrieve the service")

		return nil, err
	}

	creds, err := s.Teams().Team(s.team).ServiceCredentials().List(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to retrieve the service credentials")

		return nil, err
	}

	if creds != nil && len(creds.Items) > 0 {
		return nil, fmt.Errorf("the service can not be deleted, please delete all service credentials first")
	}

	return original, s.Store().Client().Delete(ctx, store.DeleteOptions.From(original))
}

// List returns a list of services we have access to
func (s *servicesImpl) List(ctx context.Context) (*servicesv1.ServiceList, error) {
	user := authentication.MustGetIdentity(ctx)
	if !user.IsMember(s.team) && !user.IsGlobalAdmin() {
		return nil, NewErrNotAllowed("must be global admin or a team member")
	}

	list := &servicesv1.ServiceList{}

	return list, s.Store().Client().List(ctx,
		store.ListOptions.InNamespace(s.team),
		store.ListOptions.InTo(list),
	)
}

// Get returns a specific service
func (s *servicesImpl) Get(ctx context.Context, name string) (*servicesv1.Service, error) {
	user := authentication.MustGetIdentity(ctx)
	if !user.IsMember(s.team) && !user.IsGlobalAdmin() {
		return nil, NewErrNotAllowed("must be global admin or a team member")
	}

	service := &servicesv1.Service{}

	if err := s.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(s.team),
		store.GetOptions.InTo(service),
		store.GetOptions.WithName(name),
	); err != nil {
		if kerrors.IsNotFound(err) {
			return nil, ErrNotFound
		}

		log.WithError(err).Error("failed to retrieve the service")
		return nil, err
	}

	return service, nil
}

// Update is used to update the service
func (s *servicesImpl) Update(ctx context.Context, service *servicesv1.Service) error {
	// @TODO check whether the user is an admin in the team

	user := authentication.MustGetIdentity(ctx)
	if !user.IsMember(s.team) && !user.IsGlobalAdmin() {
		return NewErrNotAllowed("must be global admin or a team member")
	}

	if err := IsValidResourceName("service", service.Name); err != nil {
		return err
	}

	if !ResourceNameFilter.MatchString(service.Spec.Kind) {
		return validation.NewError("service has failed validation").WithFieldErrorf(
			"kind",
			validation.InvalidValue,
			"must match %s",
			ResourceNameFilter.String(),
		)
	}

	if !ResourceNameFilter.MatchString(service.Spec.Plan) {
		return validation.NewError("service has failed validation").WithFieldErrorf(
			"plan",
			validation.InvalidValue,
			"must match %s",
			ResourceNameFilter.String(),
		)
	}

	existing, err := s.Get(ctx, service.Name)
	if err != nil && err != ErrNotFound {
		return err
	}

	if existing != nil {
		verr := validation.NewError("service has failed validation")
		if existing.Spec.Kind != service.Spec.Kind {
			verr.AddFieldErrorf("kind", validation.ReadOnly, "can not be changed after a service was created")
		}
		if existing.Spec.Plan != service.Spec.Plan {
			verr.AddFieldErrorf("plan", validation.ReadOnly, "can not be changed after a service was created")
		}
		if verr.HasErrors() {
			return verr
		}
	}

	if service.Namespace == "" {
		service.Namespace = s.team
	}

	if service.Namespace != s.team {
		return validation.NewError("service has failed validation").WithFieldErrorf(
			"namespace",
			validation.MustExist,
			"must be the same as the team name: %q",
			s.team,
		)
	}

	provider := s.serviceProviders.GetProviderForKind(service.Spec.Kind)
	if provider == nil {
		return validation.NewError("%q failed validation", service.Name).
			WithFieldErrorf("kind", validation.InvalidType, "%q is not a known service kind", service.Spec.Kind)
	}

	if err := s.validateConfiguration(ctx, service, provider); err != nil {
		return err
	}

	if err := s.validateCredentials(ctx, service, provider); err != nil {
		return err
	}

	return s.Store().Client().Update(ctx,
		store.UpdateOptions.To(service),
		store.UpdateOptions.WithCreate(true),
	)
}

func (s *servicesImpl) validateConfiguration(ctx context.Context, service *servicesv1.Service, provider ServiceProvider) error {
	plan, err := s.servicePlans.Get(ctx, service.Spec.Plan)
	if err != nil {
		if err == ErrNotFound {
			return validation.NewError("%q failed validation", service.Name).
				WithFieldErrorf("plan", validation.MustExist, "%q does not exist", service.Spec.Plan)
		}
		log.WithFields(log.Fields{
			"service": service.Name,
			"team":    s.team,
			"plan":    service.Spec.Plan,
		}).WithError(err).Error("failed to load service plan")

		return err
	}

	if plan.Spec.Kind != service.Spec.Kind {
		return validation.NewError("%q failed validation", service.Name).
			WithFieldErrorf("plan", validation.InvalidType, "service has kind %q, but plan has %q", service.Spec.Kind, plan.Spec.Kind)
	}

	planConfiguration := make(map[string]interface{})
	if err := plan.Spec.GetConfiguration(&planConfiguration); err != nil {
		return fmt.Errorf("failed to parse plan configuration values: %s", err)
	}

	serviceConfig := make(map[string]interface{})
	if err := service.Spec.GetConfiguration(&serviceConfig); err != nil {
		return fmt.Errorf("failed to parse service configuration values: %s", err)
	}

	schema, err := provider.PlanJSONSchema(service.Spec.Kind, service.PlanShortName())
	if err != nil {
		return err
	}

	if err := jsonschema.Validate(schema, "service", service.Spec.Configuration.Raw); err != nil {
		return err
	}

	editableParams, err := s.servicePlans.GetEditablePlanParams(ctx, s.team, service.Spec.Kind)
	if err != nil {
		return err
	}

	verr := validation.NewError("%q failed validation", service.Name)

	for paramName, paramValue := range serviceConfig {
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

func (s *servicesImpl) validateCredentials(ctx context.Context, service *servicesv1.Service, provider ServiceProvider) error {
	expectedKinds, err := provider.RequiredCredentialTypes(service.Spec.Kind)
	if err != nil {
		return err
	}

	creds := service.Spec.Credentials

	if expectedKinds == nil {
		if creds.Kind != "" {
			return validation.NewError("service has failed validation").WithFieldError(
				"credentials",
				validation.InvalidType,
				"should not be set as this service kind doesn't require credentials",
			)
		}
		return nil
	}

	found := false
	for _, gvk := range expectedKinds {
		if creds.HasGroupVersionKind(gvk) {
			found = true
			break
		}
	}

	if !found {
		var expected []string
		for _, gvk := range expectedKinds {
			expected = append(expected, utils.FormatGroupVersionKind(gvk))
		}
		return validation.NewError("service has failed validation").WithFieldErrorf(
			"credentials",
			validation.InvalidType,
			"should be one of: %s",
			strings.Join(expected, ", "),
		)
	}

	var alloc configv1.Allocation
	credentialAllocations, err := s.Teams().Team(s.team).Allocations().ListAllocationsByType(
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
		return validation.NewError("service has failed validation").WithFieldErrorf(
			"credentials",
			validation.MustExist,
			"%q does not exist or it is not assigned to the team",
			creds.Name,
		)
	}

	return nil
}
