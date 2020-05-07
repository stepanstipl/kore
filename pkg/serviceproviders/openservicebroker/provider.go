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
	"github.com/appvia/kore/pkg/utils/jsonschema"

	osb "github.com/kubernetes-sigs/go-open-service-broker-client/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	DefaultPlan                 = "kore-default"
	MetadataKeyConfiguration    = "kore.appvia.io/configuration"
	MetadataKeyDisplayName      = "displayName"
	MetadataKeyImageURL         = "imageUrl"
	MetadataKeyDescription      = "longDescription"
	MetadataKeyDocumentationURL = "documentationUrl"
	ComponentProvision          = "Provision"
	ComponentUpdate             = "Update"
	ComponentDeprovision        = "Deprovision"
	ComponentBind               = "Bind"
	ComponentUnbind             = "Unbind"
)

var _ kore.ServiceProvider = &Provider{}

type providerService struct {
	osbService  osb.Service
	defaultPlan *providerPlan
	plans       map[string]providerPlan
}

type providerPlan struct {
	osbPlan           osb.Plan
	name              string
	serviceID         string
	schema            string
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
	for _, catalogService := range catalog.Services {
		if !kore.ResourceNameFilter.MatchString(catalogService.Name) {
			return nil, fmt.Errorf("%q service name is invalid, must match %s", catalogService.Name, kore.ResourceNameFilter.String())
		}

		service := providerService{
			osbService: catalogService,
			plans:      map[string]providerPlan{},
		}

		for _, catalogPlan := range catalogService.Plans {
			if !kore.ResourceNameFilter.MatchString(catalogPlan.Name) {
				return nil, fmt.Errorf("%q plan name is invalid, must match %s", catalogPlan.Name, kore.ResourceNameFilter.String())
			}

			servicePlan, schema, credentialsSchema, err := parseCatalogPlan(catalogService, catalogPlan)
			if err != nil {
				return nil, err
			}

			plan := providerPlan{
				osbPlan:           catalogPlan,
				name:              catalogPlan.Name,
				serviceID:         catalogService.ID,
				schema:            schema,
				credentialsSchema: credentialsSchema,
			}

			if catalogPlan.Name == DefaultPlan {
				service.defaultPlan = &plan

				if schema == "" {
					return nil, fmt.Errorf("%s-%s plan does not have a schema for provisioning", catalogService.Name, catalogPlan.Name)
				}

				if credentialsSchema == "" {
					return nil, fmt.Errorf("%s-%s plan does not have a schema for bind", catalogService.Name, catalogPlan.Name)
				}
			} else {
				plans = append(plans, *servicePlan)
				service.plans[servicePlan.Name] = plan
			}
		}

		services[catalogService.Name] = service
	}

	return &Provider{
		name:         name,
		client:       client,
		services:     services,
		servicePlans: plans,
	}, nil
}

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) Type() string {
	return "openservicebroker"
}

func (p *Provider) Kinds() []servicesv1.ServiceKind {
	var res []servicesv1.ServiceKind
	for _, service := range p.services {
		res = append(res, servicesv1.ServiceKind{
			TypeMeta: metav1.TypeMeta{
				Kind:       servicesv1.ServiceKindGVK.Kind,
				APIVersion: servicesv1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      service.osbService.Name,
				Namespace: kore.HubNamespace,
			},
			Spec: servicesv1.ServiceKindSpec{
				Summary:          service.osbService.Description,
				DisplayName:      getMetadataStringVal(service.osbService.Metadata, MetadataKeyDisplayName, ""),
				Description:      getMetadataStringVal(service.osbService.Metadata, MetadataKeyDescription, ""),
				ImageURL:         getMetadataStringVal(service.osbService.Metadata, MetadataKeyImageURL, ""),
				DocumentationURL: getMetadataStringVal(service.osbService.Metadata, MetadataKeyDocumentationURL, ""),
			},
		})
	}
	return res
}

func (p *Provider) Plans() []servicesv1.ServicePlan {
	return p.servicePlans
}

func (p *Provider) PlanJSONSchema(kind string, planName string) (string, error) {
	plan, found, err := p.planWithFilter(kind, planName, func(p providerPlan) bool { return p.schema != "" })
	if err != nil {
		return "", err
	}
	if !found {
		return "", nil
	}

	return plan.schema, nil

}

func (p *Provider) CredentialsJSONSchema(kind string, planName string) (string, error) {
	plan, found, err := p.planWithFilter(kind, planName, func(p providerPlan) bool { return p.credentialsSchema != "" })
	if err != nil {
		return "", err
	}
	if !found {
		return "", nil
	}

	return plan.credentialsSchema, nil
}

func (p *Provider) RequiredCredentialTypes(kind string) ([]schema.GroupVersionKind, error) {
	_, ok := p.services[kind]
	if !ok {
		return nil, fmt.Errorf("%q service kind is invalid", kind)
	}
	return nil, nil
}

func (p *Provider) plan(service *servicesv1.Service) (providerPlan, error) {
	res, found, err := p.planWithFilter(service.Spec.Kind, service.PlanShortName(), nil)
	if err != nil {
		return providerPlan{}, err
	}
	if !found {
		return providerPlan{}, fmt.Errorf("%q service must define a %q plan to create and use custom plans", service.Spec.Kind, DefaultPlan)
	}

	return res, nil
}

func (p *Provider) planWithFilter(kind, planName string, filter func(providerPlan) bool) (providerPlan, bool, error) {
	service, ok := p.services[kind]
	if !ok {
		return providerPlan{}, false, fmt.Errorf("%q service kind is invalid", kind)
	}

	if planName != "" {
		if plan, ok := service.plans[planName]; ok {
			if filter == nil || filter(plan) {
				return plan, true, nil
			}
		}
	}

	if service.defaultPlan != nil {
		if filter == nil || filter(*service.defaultPlan) {
			return *service.defaultPlan, true, nil
		}
	}

	return providerPlan{}, false, nil
}

func parseCatalogPlan(service osb.Service, plan osb.Plan) (*servicesv1.ServicePlan, string, string, error) {
	schemaStr, err := getPlanSchema(plan)
	if err != nil {
		return nil, "", "", err
	}

	credentialsSchemaStr, err := getCredentialsSchema(plan)
	if err != nil {
		return nil, "", "", err
	}

	res := &servicesv1.ServicePlan{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServicePlan",
			APIVersion: servicesv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      plan.Name,
			Namespace: kore.HubNamespace,
		},
		Spec: servicesv1.ServicePlanSpec{
			Kind:        service.Name,
			Summary:     fmt.Sprintf("%s service - %s plan", service.Name, plan.Name),
			Description: plan.Description,
		},
	}

	configuration := map[string]interface{}{}

	if rawConfiguration, hasConfig := plan.Metadata[MetadataKeyConfiguration]; hasConfig {
		var ok bool
		configuration, ok = rawConfiguration.(map[string]interface{})
		if !ok {
			return nil, "", "", fmt.Errorf("%s-%s plan has an invalid configuration, it must be an object", service.Name, plan.Name)
		}
	}

	if schemaStr != "" {
		schema := &jsonschema.Schema{}
		if err := json.Unmarshal([]byte(schemaStr), schema); err != nil {
			return nil, "", "", fmt.Errorf("failed to unmarshal JSON schema: %w", err)
		}

		for name, prop := range schema.Properties {
			if _, isSet := configuration[name]; !isSet {
				defaultValue, err := prop.ParseDefault()
				if err != nil {
					return nil, "", "", fmt.Errorf("invalid default value %v in JSON schema: %w", prop.Default, err)
				}
				if defaultValue != nil {
					configuration[name] = defaultValue
				}
			}
		}
	}

	if err := res.Spec.SetConfiguration(configuration); err != nil {
		return nil, "", "", err
	}

	return res, schemaStr, credentialsSchemaStr, nil
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

func getMetadataStringVal(metadata map[string]interface{}, key, def string) string {
	val, ok := metadata[key]
	if ok {
		if strVal, ok := val.(string); ok && strVal != "" {
			return strVal
		}
	}

	return def
}
