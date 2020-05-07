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
	"time"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/jsonschema"
	"github.com/appvia/kore/pkg/utils/validation"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

type ServiceProviderFactory interface {
	// Type returns the service provider type
	Type() string
	// JSONSchema returns the JSON schema for the provider's configuration
	JSONSchema() string
	// CreateProvider creates the provider
	CreateProvider(ServiceProviderContext, *servicesv1.ServiceProvider) (_ ServiceProvider, complete bool, _ error)
	// TearDownProvider makes sure all provider resources are deleted
	TearDownProvider(ServiceProviderContext, *servicesv1.ServiceProvider) (complete bool, err error)
	// RequiredCredentialTypes returns with the required credential types
	RequiredCredentialTypes() []schema.GroupVersionKind
}

type ServiceProviderContext struct {
	Context context.Context
	Logger  log.FieldLogger
	Client  client.Client
}

func (s ServiceProviderContext) Deadline() (deadline time.Time, ok bool) {
	return s.Context.Deadline()
}

func (s ServiceProviderContext) Done() <-chan struct{} {
	return s.Context.Done()
}

func (s ServiceProviderContext) Err() error {
	return s.Context.Err()
}

func (s ServiceProviderContext) Value(key interface{}) interface{} {
	return s.Context.Value(key)
}

func NewServiceProviderContext(
	ctx context.Context,
	logger log.FieldLogger,
	client client.Client,
) ServiceProviderContext {
	return ServiceProviderContext{
		Context: ctx,
		Logger:  logger,
		Client:  client,
	}
}

type ServiceProvider interface {
	// Name returns a unique id for the service provider
	Name() string
	// Kinds returns a list of service kinds supported by this provider. All kinds must be unique
	Kinds() []servicesv1.ServiceKind
	// Plans returns all default service plans for this provider
	Plans() []servicesv1.ServicePlan
	// PlanJSONSchema returns the JSON schema for the given service kind and plan
	PlanJSONSchema(kind string, plan string) (string, error)
	// CredentialsJSONSchema returns the JSON schema for the credentials configuration
	CredentialsJSONSchema(kind string, plan string) (string, error)
	// RequiredCredentialTypes returns with the required credential types
	RequiredCredentialTypes(kind string) ([]schema.GroupVersionKind, error)
	// Reconcile will create or update the service
	Reconcile(ServiceProviderContext, *servicesv1.Service) (reconcile.Result, error)
	// Delete will delete the service
	Delete(ServiceProviderContext, *servicesv1.Service) (reconcile.Result, error)
	// ReconcileCredentials will create or update the service credentials
	ReconcileCredentials(ServiceProviderContext, *servicesv1.Service, *servicesv1.ServiceCredentials) (reconcile.Result, map[string]string, error)
	// DeleteCredentials will delete the service credentials
	DeleteCredentials(ServiceProviderContext, *servicesv1.Service, *servicesv1.ServiceCredentials) (reconcile.Result, error)
}

// ServiceProviders is the interface to manage service providers
type ServiceProviders interface {
	// Delete is used to delete a service provider in the kore
	Delete(context.Context, string) (*servicesv1.ServiceProvider, error)
	// Get returns the service provider
	Get(context.Context, string) (*servicesv1.ServiceProvider, error)
	// List returns the existing service providers
	List(context.Context) (*servicesv1.ServiceProviderList, error)
	// Has checks if a service provider exists
	Has(context.Context, string) (bool, error)
	// Update is responsible for updating a service provider
	Update(context.Context, *servicesv1.ServiceProvider) error
	// GetEditableProviderParams returns with the editable service provider parameters for a specific team and service kind
	GetEditableProviderParams(ctx context.Context, team string, clusterKind string) (map[string]bool, error)
	// GetProviderForKind returns a service provider for the given service kind
	GetProviderForKind(kind string) ServiceProvider
	// Register registers a new provider
	Register(ServiceProviderContext, *servicesv1.ServiceProvider) (_ ServiceProvider, complete bool, _ error)
	// Unregister removes the given provider
	Unregister(ServiceProviderContext, *servicesv1.ServiceProvider) (complete bool, _ error)
}

type serviceProvidersImpl struct {
	Interface
	providers     map[string]ServiceProvider
	providersLock sync.RWMutex
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

	if err := jsonschema.Validate(factory.JSONSchema(), "provider", provider.Spec.Configuration); err != nil {
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

// Register registers a new provider
func (p *serviceProvidersImpl) Register(ctx ServiceProviderContext, serviceProvider *servicesv1.ServiceProvider) (_ ServiceProvider, complete bool, _ error) {
	p.providersLock.Lock()
	defer p.providersLock.Unlock()

	factory, ok := serviceProviderFactories[serviceProvider.Spec.Type]
	if !ok {
		return nil, false, fmt.Errorf("service provider type %q is invalid", serviceProvider.Spec.Type)
	}

	provider, complete, err := factory.CreateProvider(ctx, serviceProvider)
	if err != nil || !complete {
		return nil, complete, err
	}

	var supportedKinds []string
	for _, kind := range provider.Kinds() {
		supportedKinds = append(supportedKinds, kind.Name)
	}

	for _, kind := range supportedKinds {
		if p, registered := p.providers[kind]; registered {
			if p.Name() != serviceProvider.Name {
				return nil, false, fmt.Errorf("service kind is already registered by an other service provider: %s", p.Name())
			}
		}
	}

	// check for removed kinds
	for kind, provider := range p.providers {
		if provider.Name() == serviceProvider.Name && !utils.Contains(kind, supportedKinds) {
			if err := p.unregisterKind(ctx, kind); err != nil {
				return nil, false, err
			}
		}
	}

	if p.providers == nil {
		p.providers = map[string]ServiceProvider{}
	}

	for _, kind := range supportedKinds {
		p.providers[kind] = provider
	}

	return provider, true, nil
}

// Unregister removes the given provider
func (p *serviceProvidersImpl) Unregister(ctx ServiceProviderContext, serviceProvider *servicesv1.ServiceProvider) (complete bool, _ error) {
	for _, kind := range serviceProvider.Status.SupportedKinds {
		if err := p.unregisterKind(ctx, kind); err != nil {
			return false, err
		}
	}

	factory, ok := serviceProviderFactories[serviceProvider.Spec.Type]
	if !ok {
		return false, fmt.Errorf("service provider type %q is invalid", serviceProvider.Spec.Type)
	}

	return factory.TearDownProvider(ctx, serviceProvider)
}

func (p *serviceProvidersImpl) unregisterKind(ctx context.Context, kind string) error {
	plans, err := p.ServicePlans().ListFiltered(ctx, func(x servicesv1.ServicePlan) bool { return x.Spec.Kind == kind })
	if err != nil {
		return fmt.Errorf("failed to get service plans: %w", err)
	}
	for _, plan := range plans.Items {
		if _, err := p.ServicePlans().Delete(ctx, plan.Name); err != nil && err != ErrNotFound {
			return fmt.Errorf("failed to delete service plan %q: %w", plan.Name, err)
		}
	}

	_, err = p.ServiceKinds().Delete(ctx, kind)
	if err != nil && err != ErrNotFound {
		return fmt.Errorf("failed to delete service kind: %w", err)
	}

	p.providersLock.Lock()
	defer p.providersLock.Unlock()

	delete(p.providers, kind)

	return nil
}

func (p *serviceProvidersImpl) GetProviderForKind(kind string) ServiceProvider {
	p.providersLock.RLock()
	defer p.providersLock.RUnlock()

	return p.providers[kind]
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
