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
	"encoding/json"
	"fmt"
	"reflect"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"

	"github.com/appvia/kore/pkg/utils/configuration"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/jsonschema"
	"github.com/appvia/kore/pkg/utils/validation"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

// Services returns the an interface for handling services
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Services
type Services interface {
	// Delete is used to delete a service
	Delete(context.Context, string) (*servicesv1.Service, error)
	// Get returns a specific service
	Get(context.Context, string) (*servicesv1.Service, error)
	// List returns a list of services
	List(context.Context) (*servicesv1.ServiceList, error)
	// ListFiltered returns a list of services using the given filter.
	ListFiltered(context.Context, func(servicesv1.Service) bool) (*servicesv1.ServiceList, error)
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
	logger := log.WithFields(log.Fields{
		"service": name,
		"team":    s.team,
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

	creds, err := s.Teams().Team(s.team).ServiceCredentials().List(ctx, func(s servicesv1.ServiceCredentials) bool {
		return s.Spec.Service.Equals(corev1.MustGetOwnershipFromObject(original))
	})
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
	list := &servicesv1.ServiceList{}

	return list, s.Store().Client().List(ctx,
		store.ListOptions.InNamespace(s.team),
		store.ListOptions.InTo(list),
	)
}

// ListFiltered returns a list of services using the given filter.
// A service is included if the filter function returns true
func (p servicesImpl) ListFiltered(ctx context.Context, filter func(plan servicesv1.Service) bool) (*servicesv1.ServiceList, error) {
	res := []servicesv1.Service{}

	servicesList, err := p.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, service := range servicesList.Items {
		if filter(service) {
			res = append(res, service)
		}
	}

	servicesList.Items = res

	return servicesList, nil
}

// Get returns a specific service
func (s *servicesImpl) Get(ctx context.Context, name string) (*servicesv1.Service, error) {
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

	if err := s.validateCluster(ctx, service); err != nil {
		return err
	}

	serviceKind, err := s.ServiceKinds().Get(ctx, service.Spec.Kind)
	if err != nil {
		return err
	}

	if serviceKind.Labels[Label("platform")] == "Kubernetes" {
		if service.Spec.ClusterNamespace == "" {
			return validation.NewError("service has failed validation").
				WithFieldError("clusterNamespace", validation.Required, "must be set")
		}
	} else {
		if service.Spec.ClusterNamespace != "" {
			return validation.NewError("service has failed validation").
				WithFieldError("clusterNamespace", validation.NotAllowed, "should not be set for this service type")
		}
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

	kind, err := s.serviceKinds.Get(ctx, service.Spec.Kind)
	if err != nil {
		if err == ErrNotFound {
			return validation.NewError("%q failed validation", service.Name).
				WithFieldErrorf("kind", validation.InvalidType, "%q is not a known service kind", service.Spec.Kind)
		}
		return err
	}

	if !kind.Spec.Enabled {
		return validation.NewError("%q failed validation", service.Name).
			WithFieldErrorf("kind", validation.InvalidType, "%q is not enabled", service.Spec.Kind)
	}

	if err := s.validateConfiguration(ctx, service, existing); err != nil {
		return err
	}

	return s.Store().Client().Update(ctx,
		store.UpdateOptions.To(service),
		store.UpdateOptions.WithCreate(true),
	)
}

func (s *servicesImpl) validateConfiguration(ctx context.Context, service, existing *servicesv1.Service) error {
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

	if plan.Annotations[AnnotationSystem] == "true" {
		return validation.NewError("%q failed validation", service.Name).
			WithFieldError("plan", validation.InvalidType, "system plans can not be used to create new services")
	}

	planDetails, err := s.ServicePlans().GetDetails(ctx, plan.Name, s.team, service.Spec.Cluster.Name)
	if err != nil {
		return err
	}

	planConfiguration := make(map[string]interface{})
	if err := planDetails.GetConfiguration(&planConfiguration); err != nil {
		return fmt.Errorf("failed to parse plan configuration values: %s", err)
	}

	serviceConfig := make(map[string]interface{})
	if err := configuration.ParseObjectConfiguration(ctx, s.Store().RuntimeClient(), service, &serviceConfig); err != nil {
		return fmt.Errorf("failed to parse service configuration values: %s", err)
	}

	log.WithField("serviceConfig", serviceConfig).Debug("SERVICE VALIDATE")

	if planDetails.Schema == "" && !utils.ApiExtJSONEmpty(service.Spec.Configuration) {
		if existing == nil ||
			!utils.ApiExtJSONEquals(service.Spec.Configuration, existing.Spec.Configuration) ||
			!reflect.DeepEqual(service.Spec.ConfigurationFrom, existing.Spec.ConfigurationFrom) {
			return validation.NewError("%q failed validation", service.Name).
				WithFieldErrorf(
					"configuration",
					validation.ReadOnly,
					"the service provider doesn't have a JSON schema to validate the configuration",
				)
		}
	}

	schema := &jsonschema.Schema{}
	if planDetails.Schema != "" {
		if err := jsonschema.Validate(planDetails.Schema, "service", serviceConfig); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(planDetails.Schema), schema); err != nil {
			return err
		}
	}

	verr := validation.NewError("%q failed validation", service.Name)

	for paramName, paramValue := range serviceConfig {
		schemaProperty := schema.Properties[paramName]

		// If a const value changes in the schema, we have to allow users to migrate their objects
		if schemaProperty != nil && schemaProperty.Const != nil {
			continue
		}

		if !reflect.DeepEqual(paramValue, planConfiguration[paramName]) {
			if !utils.Contains(paramName, planDetails.EditableParams) {
				verr.AddFieldErrorf(paramName, validation.ReadOnly, "can not be changed")
			}
		}
	}
	if verr.HasErrors() {
		return verr
	}

	return nil
}

func (s *servicesImpl) validateCluster(ctx context.Context, service *servicesv1.Service) error {
	if service.Spec.Cluster.Name == "" {
		return validation.NewError("%q failed validation", service.Name).WithFieldError(
			"cluster.name",
			validation.Required,
			"must be set",
		)
	}

	if service.Spec.Cluster.Namespace != service.Namespace {
		return validation.NewError("%q failed validation", service.Name).WithFieldErrorf(
			"cluster.namespace",
			validation.InvalidValue,
			"must be the same as the team name: %q",
			s.team,
		)
	}

	if !service.Spec.Cluster.HasGroupVersionKind(clustersv1.ClusterGroupVersionKind) {
		return validation.NewError("%q failed validation", service.Name).WithFieldErrorf(
			"cluster",
			validation.InvalidValue,
			"must have type of %s",
			clustersv1.ClusterGroupVersionKind,
		)
	}

	_, err := s.Teams().Team(s.team).Clusters().Get(ctx, service.Spec.Cluster.Name)
	if err != nil {
		if err == ErrNotFound {
			return validation.NewError("%q failed validation", service.Name).WithFieldErrorf(
				"cluster",
				validation.MustExist,
				"%q cluster does not exist",
				service.Spec.Cluster.Name,
			)
		}
	}
	return nil
}
