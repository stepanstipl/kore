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

	"k8s.io/apimachinery/pkg/runtime/schema"

	log "github.com/sirupsen/logrus"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// DefaultServiceProviders contains all registered service providers
var DefaultServiceProviders = &ServiceProviderRegistry{
	providers: map[string]ServiceProvider{},
}

type ServiceProvider interface {
	// Name returns a unique id for the service provider
	Name() string
	// Kinds returns a list of service kinds supported by this provider. All kinds must be unique
	Kinds() []string
	// Plans returns all default service plans for this provider
	Plans() []servicesv1.ServicePlan
	// JSONSchema returns the JSON schema for a service kind
	JSONSchema(kind string) string
	// CredentialsJSONSchema returns the JSON schema for the credentials configuration
	CredentialsJSONSchema(kind string) string
	// RequiredCredentialTypes returns with the required credential types
	RequiredCredentialTypes(kind string) []schema.GroupVersionKind
	// Reconcile will create or update the service
	Reconcile(context.Context, log.FieldLogger, *servicesv1.Service) (reconcile.Result, error)
	// Delete will delete the service
	Delete(context.Context, log.FieldLogger, *servicesv1.Service) (reconcile.Result, error)
}

type ServiceProviderRegistry struct {
	providers map[string]ServiceProvider
}

func (s *ServiceProviderRegistry) Register(provider ServiceProvider) {
	_, exists := s.providers[provider.Name()]
	if exists {
		panic(fmt.Errorf("service provider with name %q was already registered", provider.Name()))
	}

	s.providers[provider.Name()] = provider
}

func (s *ServiceProviderRegistry) GetProviderForKind(kind string) ServiceProvider {
	for _, provider := range s.providers {
		for _, k := range provider.Kinds() {
			if strings.EqualFold(k, kind) {
				return provider
			}
		}
	}
	return nil
}

func (s *ServiceProviderRegistry) GetAllPlans() []servicesv1.ServicePlan {
	var res []servicesv1.ServicePlan
	for _, provider := range s.providers {
		res = append(res, provider.Plans()...)
	}
	return res
}

func (s *ServiceProviderRegistry) Providers() map[string]ServiceProvider {
	res := make(map[string]ServiceProvider, len(s.providers))
	for k, v := range s.providers {
		res[k] = v
	}
	return res
}
