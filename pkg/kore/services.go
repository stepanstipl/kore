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
func (c *servicesImpl) Delete(ctx context.Context, name string) (*servicesv1.Service, error) {
	// @TODO check whether the user is an admin in the team

	user := authentication.MustGetIdentity(ctx)
	if !user.IsMember(c.team) && !user.IsGlobalAdmin() {
		return nil, NewErrNotAllowed("must be global admin or a team member")
	}

	logger := log.WithFields(log.Fields{
		"service": name,
		"team":    c.team,
		"user":    user.Username(),
	})
	logger.Info("attempting to delete the service")

	original, err := c.Get(ctx, name)
	if err != nil {
		if err == ErrNotFound {
			return nil, err
		}

		logger.WithError(err).Error("failed to retrieve the service")

		return nil, err
	}

	return original, c.Store().Client().Delete(ctx, store.DeleteOptions.From(original))
}

// List returns a list of services we have access to
func (c *servicesImpl) List(ctx context.Context) (*servicesv1.ServiceList, error) {
	user := authentication.MustGetIdentity(ctx)
	if !user.IsMember(c.team) && !user.IsGlobalAdmin() {
		return nil, NewErrNotAllowed("must be global admin or a team member")
	}

	list := &servicesv1.ServiceList{}

	return list, c.Store().Client().List(ctx,
		store.ListOptions.InNamespace(c.team),
		store.ListOptions.InTo(list),
	)
}

// Get returns a specific service
func (c *servicesImpl) Get(ctx context.Context, name string) (*servicesv1.Service, error) {
	user := authentication.MustGetIdentity(ctx)
	if !user.IsMember(c.team) && !user.IsGlobalAdmin() {
		return nil, NewErrNotAllowed("must be global admin or a team member")
	}

	service := &servicesv1.Service{}

	if err := c.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(c.team),
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
func (c *servicesImpl) Update(ctx context.Context, service *servicesv1.Service) error {
	// @TODO check whether the user is an admin in the team

	user := authentication.MustGetIdentity(ctx)
	if !user.IsMember(c.team) && !user.IsGlobalAdmin() {
		return NewErrNotAllowed("must be global admin or a team member")
	}

	existing, err := c.Get(ctx, service.Name)
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
		service.Namespace = c.team
	}

	if service.Namespace != c.team {
		return validation.NewError("service has failed validation").WithFieldErrorf(
			"namespace",
			validation.MustExist,
			"must be the same as the team name: %q",
			c.team,
		)
	}

	if len(service.Name) > 40 {
		return validation.NewError("service has failed validation").
			WithFieldError("name", validation.MaxLength, "must be 40 characters or less")
	}

	if err := c.validateConfiguration(ctx, service); err != nil {
		return err
	}

	return c.Store().Client().Update(ctx,
		store.UpdateOptions.To(service),
		store.UpdateOptions.WithCreate(true),
	)
}

func (c *servicesImpl) validateConfiguration(ctx context.Context, service *servicesv1.Service) error {
	plan, err := c.servicePlans.Get(ctx, service.Spec.Plan)
	if err != nil {
		if err == ErrNotFound {
			return validation.NewError("%q failed validation", service.Name).
				WithFieldErrorf("plan", validation.MustExist, "%q does not exist", service.Spec.Plan)
		}
		log.WithFields(log.Fields{
			"service": service.Name,
			"team":    c.team,
			"plan":    service.Spec.Plan,
		}).WithError(err).Error("failed to load service plan")

		return err
	}

	if !strings.EqualFold(plan.Spec.Kind, service.Spec.Kind) {
		return validation.NewError("%q failed validation", service.Name).
			WithFieldErrorf("plan", validation.InvalidType, "service has service kind %q, but plan has %q", service.Spec.Kind, plan.Spec.Kind)
	}

	planConfiguration := make(map[string]interface{})
	if err := json.NewDecoder(bytes.NewReader(plan.Spec.Configuration.Raw)).Decode(&planConfiguration); err != nil {
		return fmt.Errorf("failed to parse plan configuration values: %s", err)
	}

	serviceConfig := make(map[string]interface{})
	if err := json.NewDecoder(bytes.NewReader(service.Spec.Configuration.Raw)).Decode(&serviceConfig); err != nil {
		return fmt.Errorf("failed to parse service configuration values: %s", err)
	}

	provider := c.serviceProviders.GetProviderForKind(plan.Spec.Kind)
	if provider == nil {
		return validation.NewError("%q failed validation", service.Name).
			WithFieldErrorf("kind", validation.InvalidType, "%q is not a known service kind", plan.Spec.Kind)
	}

	if err := jsonschema.Validate(provider.JSONSchema(service.Spec.Kind), "plan", service.Spec.Configuration.Raw); err != nil {
		return err
	}

	editableParams, err := c.servicePlans.GetEditablePlanParams(ctx, c.team, service.Spec.Kind)
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
