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
	"sync"

	"github.com/appvia/kore/pkg/utils/configuration"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/jsonschema"
	"github.com/appvia/kore/pkg/utils/validation"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var serviceProviderFactories = map[string]ServiceProviderFactory{}
var spfLock = sync.Mutex{}

func RegisterServiceProviderFactory(factory ServiceProviderFactory) {
	spfLock.Lock()
	defer spfLock.Unlock()

	_, exists := serviceProviderFactories[factory.Type()]
	if exists {
		panic(fmt.Errorf("service provider type %q was already registered", factory.Type()))
	}

	serviceProviderFactories[factory.Type()] = factory
}

func ServiceProviderFactories() map[string]ServiceProviderFactory {
	spfLock.Lock()
	defer spfLock.Unlock()

	res := make(map[string]ServiceProviderFactory, len(serviceProviderFactories))
	for k, v := range serviceProviderFactories {
		res[k] = v
	}
	return res
}

type ServiceProviderCatalog struct {
	Plans []servicesv1.ServicePlan
	Kinds []servicesv1.ServiceKind
}

type ServiceProviderFactory interface {
	// Type returns the service provider type
	Type() string
	// JSONSchema returns the JSON schema for the provider's configuration
	JSONSchema() string
	// Create creates the provider
	Create(Context, *servicesv1.ServiceProvider) (ServiceProvider, error)
	// SetUp makes sure all provider dependencies are created
	SetUp(Context, *servicesv1.ServiceProvider) (complete bool, _ error)
	// TearDown makes sure all provider dependencies are deleted
	TearDown(Context, *servicesv1.ServiceProvider) (complete bool, _ error)
	// RequiredCredentialTypes returns with the required credential types
	RequiredCredentialTypes() []schema.GroupVersionKind
	// DefaultProviders returns with a list of providers which should be automatically installed
	DefaultProviders() []servicesv1.ServiceProvider
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ServiceProvider
type ServiceProvider interface {
	// Name returns a unique id for the service provider
	Name() string
	// AdminServices returns with the default admin services we need to install when the service provider is registered
	AdminServices() []servicesv1.Service
	// Catalog returns with the service provider catalog. It's the provider's responsibility to load the catalog if required
	Catalog(Context, *servicesv1.ServiceProvider) (ServiceProviderCatalog, error)
	// Reconcile will create or update the service
	Reconcile(Context, *servicesv1.Service) (reconcile.Result, error)
	// Delete will delete the service
	Delete(Context, *servicesv1.Service) (reconcile.Result, error)
	// ReconcileCredentials will create or update the service credentials
	ReconcileCredentials(Context, *servicesv1.Service, *servicesv1.ServiceCredentials) (reconcile.Result, map[string]string, error)
	// DeleteCredentials will delete the service credentials
	DeleteCredentials(Context, *servicesv1.Service, *servicesv1.ServiceCredentials) (reconcile.Result, error)
}

// ServiceProviders is the interface to manage service providers
type ServiceProviders interface {
	// Delete is used to delete a service provider in the kore
	Delete(context.Context, string) (*servicesv1.ServiceProvider, error)
	// Get returns the service provider
	Get(context.Context, string) (*servicesv1.ServiceProvider, error)
	// List returns the existing service providers
	List(context.Context) (*servicesv1.ServiceProviderList, error)
	// ListFiltered returns a list of service providers using the given filter.
	ListFiltered(context.Context, func(servicesv1.ServiceProvider) bool) (*servicesv1.ServiceProviderList, error)
	// Has checks if a service provider exists
	Has(context.Context, string) (bool, error)
	// Update is responsible for updating a service provider
	Update(context.Context, *servicesv1.ServiceProvider) error
	// GetEditableProviderParams returns with the editable service provider parameters for a specific team and service kind
	GetEditableProviderParams(ctx context.Context, team string, clusterKind string) (map[string]bool, error)
	// GetProviderForKind returns a service provider for the given service kind
	GetProviderForKind(ctx Context, kind string) (ServiceProvider, error)
	// Register registers a new service provider
	Register(ctx Context, serviceProvider *servicesv1.ServiceProvider) (ServiceProvider, error)
	// SetUp makes sure the provider dependencies are created
	SetUp(Context, *servicesv1.ServiceProvider) (complete bool, _ error)
	// Catalog loads the service provider catalog
	Catalog(Context, *servicesv1.ServiceProvider) (ServiceProviderCatalog, error)
	// Unregister removes the given provider
	Unregister(Context, *servicesv1.ServiceProvider) (complete bool, _ error)
}

type serviceProvidersImpl struct {
	Interface
	providers       map[string]ServiceProvider
	providersByKind map[string]ServiceProvider
	providersLock   sync.Mutex
}

// Update is responsible for updating a service provider
func (p *serviceProvidersImpl) Update(ctx context.Context, provider *servicesv1.ServiceProvider) error {
	existing, err := p.Get(ctx, provider.Name)
	if err != nil && err != ErrNotFound {
		return err
	}

	if existing != nil && existing.Spec.Type != provider.Spec.Type {
		return validation.NewError("%q failed validation", provider.Name).
			WithFieldErrorf("type", validation.ReadOnly, "can not be changed")
	}

	if err := IsValidResourceName("provider", provider.Name); err != nil {
		return err
	}

	if provider.Namespace != HubNamespace {
		return validation.NewError("%q failed validation", provider.Name).
			WithFieldErrorf("namespace", validation.InvalidValue, "must be %q", HubNamespace)
	}

	factory, ok := serviceProviderFactories[provider.Spec.Type]
	if !ok {
		return validation.NewError("%q failed validation", provider.Name).
			WithFieldErrorf("type", validation.InvalidType, "%q is not a valid service provider type", provider.Spec.Type)
	}

	config := map[string]interface{}{}
	if err := configuration.ParseObjectConfiguration(ctx, p.Store().RuntimeClient(), provider, &config); err != nil {
		return fmt.Errorf("failed to parse service provider configuration: %s", err)
	}

	if err := jsonschema.Validate(factory.JSONSchema(), "provider", config); err != nil {
		return err
	}

	if err := p.validateCredentials(ctx, provider, factory); err != nil {
		return err
	}

	err = p.Store().Client().Update(ctx,
		store.UpdateOptions.To(provider),
		store.UpdateOptions.WithCreate(true),
		store.UpdateOptions.WithForce(true),
	)
	if err != nil {
		log.WithError(err).Error("failed to update a service provider")

		return err
	}

	return nil
}

// Delete is used to delete a service provider in the kore
func (p *serviceProvidersImpl) Delete(ctx context.Context, name string) (*servicesv1.ServiceProvider, error) {
	provider := &servicesv1.ServiceProvider{}
	err := p.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(HubNamespace),
		store.GetOptions.InTo(provider),
		store.GetOptions.WithName(name),
	)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		log.WithError(err).Error("failed to retrieve the service provider")

		return nil, err
	}

	if err := p.Store().Client().Delete(ctx, store.DeleteOptions.From(provider)); err != nil {
		log.WithError(err).Error("failed to delete the service provider")

		return nil, err
	}

	return provider, nil
}

// Get returns the service provider
func (p *serviceProvidersImpl) Get(ctx context.Context, name string) (*servicesv1.ServiceProvider, error) {
	provider := &servicesv1.ServiceProvider{}

	if found, err := p.Has(ctx, name); err != nil {
		return nil, err
	} else if !found {
		return nil, ErrNotFound
	}

	return provider, p.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(HubNamespace),
		store.GetOptions.WithName(name),
		store.GetOptions.InTo(provider),
	)
}

// List returns the existing service providers
func (p *serviceProvidersImpl) List(ctx context.Context) (*servicesv1.ServiceProviderList, error) {
	providers := &servicesv1.ServiceProviderList{}

	return providers, p.Store().Client().List(ctx,
		store.ListOptions.InNamespace(HubNamespace),
		store.ListOptions.InTo(providers),
	)
}

// ListFiltered returns a list of service providers using the given filter.
// A service provider is included if the filter function returns true
func (p *serviceProvidersImpl) ListFiltered(ctx context.Context, filter func(plan servicesv1.ServiceProvider) bool) (*servicesv1.ServiceProviderList, error) {
	var res []servicesv1.ServiceProvider

	serviceProvidersList, err := p.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, serviceProvider := range serviceProvidersList.Items {
		if filter(serviceProvider) {
			res = append(res, serviceProvider)
		}
	}

	serviceProvidersList.Items = res

	return serviceProvidersList, nil
}

// Has checks if a service provider exists
func (p *serviceProvidersImpl) Has(ctx context.Context, name string) (bool, error) {
	return p.Store().Client().Has(ctx,
		store.HasOptions.InNamespace(HubNamespace),
		store.HasOptions.From(&servicesv1.ServiceProvider{}),
		store.HasOptions.WithName(name),
	)
}

// GetEditableProviderParams returns with the editable service provider parameters for a specific team and service kind
func (p *serviceProvidersImpl) GetEditableProviderParams(ctx context.Context, team string, clusterKind string) (map[string]bool, error) {
	// TODO: implement this when the service provider policies are implemented
	return nil, nil
}

// Register registers a new service provider
func (p *serviceProvidersImpl) Register(ctx Context, serviceProvider *servicesv1.ServiceProvider) (ServiceProvider, error) {
	p.providersLock.Lock()
	defer p.providersLock.Unlock()

	return p.register(ctx, serviceProvider)
}

func (p *serviceProvidersImpl) register(ctx Context, serviceProvider *servicesv1.ServiceProvider) (ServiceProvider, error) {
	factory, ok := serviceProviderFactories[serviceProvider.Spec.Type]
	if !ok {
		log.WithField("serviceProvider", serviceProvider.Name).Warning("Service provider is not known, skipping registration")
		return nil, fmt.Errorf("service provider type is not supported: %q", serviceProvider.Spec.Type)
	}

	provider, err := factory.Create(ctx, serviceProvider)
	if err != nil {
		return nil, err
	}

	if p.providers == nil {
		p.providers = map[string]ServiceProvider{}
	}
	p.providers[serviceProvider.Name] = provider

	if p.providersByKind == nil {
		p.providersByKind = map[string]ServiceProvider{}
	}
	for _, kind := range serviceProvider.Status.SupportedKinds {
		p.providersByKind[kind] = provider
	}

	return provider, nil
}

func (p *serviceProvidersImpl) SetUp(ctx Context, serviceProvider *servicesv1.ServiceProvider) (complete bool, _ error) {
	spfLock.Lock()
	factory, ok := serviceProviderFactories[serviceProvider.Spec.Type]
	spfLock.Unlock()

	if !ok {
		return false, fmt.Errorf("unknown service provider type: %q", serviceProvider.Spec.Type)
	}

	return factory.SetUp(ctx, serviceProvider)
}

// Catalog loads the service provider catalog
func (p *serviceProvidersImpl) Catalog(ctx Context, serviceProvider *servicesv1.ServiceProvider) (ServiceProviderCatalog, error) {
	p.providersLock.Lock()
	defer p.providersLock.Unlock()

	provider, ok := p.providers[serviceProvider.Name]
	if !ok {
		var err error
		if provider, err = p.register(ctx, serviceProvider); err != nil {
			return ServiceProviderCatalog{}, err
		}
	}

	catalog, err := provider.Catalog(ctx, serviceProvider)
	if err != nil {
		return ServiceProviderCatalog{}, err
	}

	var supportedKinds []string
	for _, kind := range catalog.Kinds {
		supportedKinds = append(supportedKinds, kind.Name)
	}

	for _, kind := range supportedKinds {
		if rp, registered := p.providersByKind[kind]; registered {
			if rp.Name() != serviceProvider.Name {
				return ServiceProviderCatalog{}, fmt.Errorf("service kind is already registered by an other service provider: %s", rp.Name())
			}
		}
	}

	// check for removed kinds
	for kind, rp := range p.providers {
		if rp.Name() == serviceProvider.Name && !utils.Contains(kind, supportedKinds) {
			if err := p.unregisterKind(ctx, kind); err != nil {
				return ServiceProviderCatalog{}, err
			}
		}
	}

	return catalog, nil
}

// Unregister removes the given provider
func (p *serviceProvidersImpl) Unregister(ctx Context, serviceProvider *servicesv1.ServiceProvider) (complete bool, _ error) {
	p.providersLock.Lock()
	defer p.providersLock.Unlock()

	for _, kind := range serviceProvider.Status.SupportedKinds {
		if err := p.unregisterKind(ctx, kind); err != nil {
			return false, err
		}
	}

	delete(p.providers, serviceProvider.Name)

	spfLock.Lock()
	factory, ok := serviceProviderFactories[serviceProvider.Spec.Type]
	spfLock.Unlock()

	if !ok {
		return false, fmt.Errorf("unknown service provider type: %q", serviceProvider.Spec.Type)
	}

	return factory.TearDown(ctx, serviceProvider)
}

func (p *serviceProvidersImpl) unregisterKind(ctx context.Context, kind string) error {
	plans, err := p.ServicePlans().ListFiltered(ctx, func(x servicesv1.ServicePlan) bool { return x.Spec.Kind == kind })
	if err != nil {
		return fmt.Errorf("failed to get service plans: %w", err)
	}
	for _, plan := range plans.Items {
		if _, err := p.ServicePlans().Delete(ctx, plan.Name, true); err != nil && err != ErrNotFound {
			return fmt.Errorf("failed to delete service plan %q: %w", plan.Name, err)
		}
	}

	_, err = p.ServiceKinds().Delete(ctx, kind)
	if err != nil && err != ErrNotFound {
		return fmt.Errorf("failed to delete service kind: %w", err)
	}

	delete(p.providersByKind, kind)

	return nil
}

func (p *serviceProvidersImpl) GetProviderForKind(ctx Context, kind string) (ServiceProvider, error) {
	p.providersLock.Lock()
	defer p.providersLock.Unlock()

	provider, ok := p.providersByKind[kind]
	if ok {
		return provider, nil
	}

	providerList, err := p.ServiceProviders().ListFiltered(ctx, func(provider servicesv1.ServiceProvider) bool {
		return utils.Contains(kind, provider.Status.SupportedKinds)
	})
	if err != nil {
		return nil, err
	}

	if len(providerList.Items) == 0 {
		return nil, fmt.Errorf("no available service provider for kind %q", kind)
	}

	return p.register(ctx, &providerList.Items[0])
}

func (p *serviceProvidersImpl) validateCredentials(ctx context.Context, serviceProvider *servicesv1.ServiceProvider, factory ServiceProviderFactory) error {
	expectedKinds := factory.RequiredCredentialTypes()
	creds := serviceProvider.Spec.Credentials

	if expectedKinds == nil {
		if creds.Kind != "" {
			return validation.NewError("service provider has failed validation").WithFieldError(
				"credentials",
				validation.InvalidType,
				"should not be set as this service provider doesn't require credentials",
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
		return validation.NewError("service provider has failed validation").WithFieldErrorf(
			"credentials",
			validation.InvalidType,
			"should be one of: %s",
			strings.Join(expected, ", "),
		)
	}

	return nil
}
