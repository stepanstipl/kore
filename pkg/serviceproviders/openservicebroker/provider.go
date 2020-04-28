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

package openservicebroker

import (
	"encoding/json"
	"fmt"
	"reflect"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/kore/assets"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/jsonschema"

	osb "github.com/kubernetes-sigs/go-open-service-broker-client/v2"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	DefaultPlan              = "kore-default"
	MetadataKeyConfiguration = "kore.appvia.io/configuration"
	ComponentProvision       = "Provision"
	ComponentUpdate          = "Update"
	ComponentDeprovision     = "Deprovision"
	ComponentBind            = "Bind"
	ComponentUnbind          = "Unbind"
)

var _ kore.ServiceProvider = &Provider{}

type providerService struct {
	id          string
	bindable    bool
	defaultPlan *providerPlan
	plans       map[string]providerPlan
}

type providerPlan struct {
	name              string
	id                string
	serviceID         string
	schema            string
	bindable          bool
	credentialsSchema string
}

type Provider struct {
	name         string
	client       osb.Client
	servicePlans []servicesv1.ServicePlan
	services     map[string]providerService
}

// NewProvider creates a new service provider which is backed by an Open Service Broker API compatible HTTP service
func NewProvider(name string, client osb.Client) (*Provider, error) {
	var plans []servicesv1.ServicePlan
	services := map[string]providerService{}

	catalog, err := client.GetCatalog()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch catalog from service broker: %w", err)
	}
	for _, s := range catalog.Services {
		if !kore.ResourceNameFilter.MatchString(s.Name) {
			return nil, fmt.Errorf("%q service name is invalid, must match %s", s.Name, kore.ResourceNameFilter.String())
		}

		providerService := providerService{
			id:       s.ID,
			plans:    map[string]providerPlan{},
			bindable: s.Bindable,
		}

		for _, p := range s.Plans {
			if !kore.ResourceNameFilter.MatchString(p.Name) {
				return nil, fmt.Errorf("%q plan name is invalid, must match %s", p.Name, kore.ResourceNameFilter.String())
			}

			servicePlan, err := catalogPlanToServicePlan(s, p)
			if err != nil {
				return nil, err
			}

			schema, err := getPlanSchema(p)
			if err != nil {
				return nil, err
			}

			credentialsSchema, err := getCredentialsSchema(p)
			if err != nil {
				return nil, err
			}

			providerPlan := providerPlan{
				name:              p.Name,
				id:                p.ID,
				serviceID:         s.ID,
				bindable:          utils.BoolValue(p.Bindable),
				schema:            schema,
				credentialsSchema: credentialsSchema,
			}

			if p.Name == DefaultPlan {
				providerService.defaultPlan = &providerPlan

				if schema == "" {
					return nil, fmt.Errorf("%s plan does not have a schema for provisioning", p.Name)
				}

				if credentialsSchema == "" {
					return nil, fmt.Errorf("%s plan does not have a schema for bind", p.Name)
				}
			} else {
				plans = append(plans, servicePlan)
				providerService.plans[planName(s, p)] = providerPlan
			}
		}

		services[s.Name] = providerService
	}

	return &Provider{
		name:         name,
		client:       client,
		services:     services,
		servicePlans: plans,
	}, nil
}

func (s *Provider) Name() string {
	return s.name
}

func (s *Provider) Kinds() []string {
	var kinds []string
	for kind := range s.services {
		kinds = append(kinds, kind)
	}
	return kinds
}

func (s *Provider) Plans() []servicesv1.ServicePlan {
	return s.servicePlans
}

func (s *Provider) JSONSchema(kind string, plan string) (string, error) {
	p, err := s.planWithFilter(kind, plan, func(p providerPlan) bool { return p.schema != "" })
	if err != nil {
		return "", err
	}

	return p.schema, nil
}

func (s *Provider) CredentialsJSONSchema(kind string, plan string) (string, error) {
	p, err := s.planWithFilter(kind, plan, func(p providerPlan) bool { return p.credentialsSchema != "" })
	if err != nil {
		return "", err
	}

	return p.credentialsSchema, nil
}

func (s *Provider) RequiredCredentialTypes(kind string) ([]schema.GroupVersionKind, error) {
	_, ok := s.services[kind]
	if !ok {
		return nil, fmt.Errorf("%q service kind is invalid", kind)
	}
	return nil, nil
}

func (s *Provider) plan(kind, name string) (providerPlan, error) {
	return s.planWithFilter(kind, name, nil)
}

func (s *Provider) planWithFilter(kind, name string, filter func(providerPlan) bool) (providerPlan, error) {
	service, ok := s.services[kind]
	if !ok {
		return providerPlan{}, fmt.Errorf("%q service kind is invalid", kind)
	}

	if name != "" {
		if plan, ok := service.plans[name]; ok {
			if filter == nil || filter(plan) {
				return plan, nil
			}
		}
	}

	if s.services[kind].defaultPlan == nil {
		return providerPlan{}, fmt.Errorf("%q service must define a %q plan", kind, DefaultPlan)
	}
	return *s.services[kind].defaultPlan, nil
}

func catalogPlanToServicePlan(service osb.Service, plan osb.Plan) (servicesv1.ServicePlan, error) {
	name := planName(service, plan)

	configuration, ok := plan.Metadata[MetadataKeyConfiguration]
	if !ok {
		return servicesv1.ServicePlan{}, fmt.Errorf("%s plan is invalid: %s key is missing from metadata", name, MetadataKeyConfiguration)
	}

	if _, ok := configuration.(map[string]interface{}); !ok {
		return servicesv1.ServicePlan{}, fmt.Errorf("%s plan has an invalid configuration, it must be an object", name)
	}

	configJSON, err := json.Marshal(configuration)
	if err != nil {
		return servicesv1.ServicePlan{}, fmt.Errorf("%s plan is invalid: %s key can not be json encoded", name, MetadataKeyConfiguration)
	}

	return servicesv1.ServicePlan{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServicePlan",
			APIVersion: servicesv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: kore.HubNamespace,
		},
		Spec: servicesv1.ServicePlanSpec{
			Kind:          service.Name,
			Summary:       fmt.Sprintf("%s service - %s plan", service.Name, plan.Name),
			Description:   plan.Description,
			Configuration: v1beta1.JSON{Raw: configJSON},
		},
	}, nil
}

func getPlanSchema(plan osb.Plan) (string, error) {
	if plan.Schemas == nil || plan.Schemas.ServiceInstance == nil || plan.Schemas.ServiceInstance.Create == nil {
		return "", nil
	}
	return parseSchema(plan.Name+" plan", plan.Schemas.ServiceInstance.Create.Parameters)
}

func getCredentialsSchema(plan osb.Plan) (string, error) {
	if plan.Schemas == nil || plan.Schemas.ServiceBinding == nil || plan.Schemas.ServiceBinding.Create == nil {
		return "", nil
	}
	return parseSchema(plan.Name+" plan", plan.Schemas.ServiceBinding.Create.Parameters)
}

func parseSchema(subject string, val interface{}) (string, error) {
	if val == nil {
		return "", nil
	}

	var schema string
	switch reflect.TypeOf(val).Kind() {
	case reflect.Struct, reflect.Map:
		schemaBytes, err := json.Marshal(val)
		if err != nil {
			return "", fmt.Errorf("%s has an invalid provisioning schema", subject)
		}
		schema = string(schemaBytes)
	case reflect.String:
		schema = val.(string)
	default:
		return "", fmt.Errorf("%s has an invalid schema", subject)
	}

	if err := jsonschema.Validate(assets.JSONSchemaDraft07, fmt.Sprintf("%s schema", subject), schema); err != nil {
		return "", err
	}
	return schema, nil
}

func planName(service osb.Service, plan osb.Plan) string {
	return fmt.Sprintf("%s-%s", service.Name, plan.Name)
}
