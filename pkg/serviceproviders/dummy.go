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

package serviceproviders

import (
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/kore"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func init() {
	kore.RegisterServiceProviderFactory(DummyFactory{})
}

type DummyFactory struct{}

func (d DummyFactory) Type() string {
	return "dummy"
}

func (d DummyFactory) JSONSchema() string {
	return `{
		"$id": "https://appvia.io/schemas/serviceprovider/dummy.json",
		"$schema": "http://json-schema.org/draft-07/schema#",
		"description": "Dummy service plan schema",
		"type": "object",
		"additionalProperties": false,
		"required": [
			"iAmDummy"
		],
		"properties": {
			"iAmDummy": {
				"type": "string",
				"minLength": 1
			}
		}
	}`
}

func (d DummyFactory) Create(ctx kore.Context, serviceProvider *servicesv1.ServiceProvider) (kore.ServiceProvider, error) {
	return Dummy{name: serviceProvider.Name}, nil
}

func (d DummyFactory) SetUp(_ kore.Context, _ *servicesv1.ServiceProvider) (complete bool, _ error) {
	return true, nil
}

func (d DummyFactory) TearDown(_ kore.Context, _ *servicesv1.ServiceProvider) (complete bool, _ error) {
	return true, nil
}

func (d DummyFactory) DefaultProviders() []servicesv1.ServiceProvider {
	return nil
}

var _ kore.ServiceProvider = Dummy{}

type Dummy struct {
	name string
}

func (d Dummy) Name() string {
	return d.name
}

func (d Dummy) Catalog(_ kore.Context, _ *servicesv1.ServiceProvider) (kore.ServiceProviderCatalog, error) {
	return kore.ServiceProviderCatalog{
		Plans: d.plans(),
		Kinds: d.kinds(),
	}, nil
}

func (d Dummy) kinds() []servicesv1.ServiceKind {
	return []servicesv1.ServiceKind{
		{
			TypeMeta: metav1.TypeMeta{
				Kind:       servicesv1.ServiceKindGVK.Kind,
				APIVersion: servicesv1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "dummy",
				Namespace: kore.HubNamespace,
				Labels: map[string]string{
					kore.Label("platform"): "Kore",
				},
			},
			Spec: servicesv1.ServiceKindSpec{
				DisplayName: "Dummy",
				Summary:     "Dummy service used for testing",
				Enabled:     true,
			},
		},
	}
}

func (d Dummy) plans() []servicesv1.ServicePlan {
	return []servicesv1.ServicePlan{
		{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ServicePlan",
				APIVersion: servicesv1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "dummy-default",
				Namespace: "kore",
			},
			Spec: servicesv1.ServicePlanSpec{
				Kind:          "dummy",
				Description:   "Used for testing",
				Summary:       "This is a default dummy service plan",
				Configuration: &v1beta1.JSON{Raw: []byte(`{"foo":"bar"}`)},
				Schema: `{
					"$id": "https://appvia.io/schemas/services/dummy/dummy.json",
					"$schema": "http://json-schema.org/draft-07/schema#",
					"description": "Dummy service plan schema",
					"type": "object",
					"additionalProperties": false,
					"required": [
						"foo"
					],
					"properties": {
						"foo": {
							"type": "string",
							"minLength": 1
						}
					}
				}`,
				CredentialSchema: `{
					"$id": "https://appvia.io/schemas/services/dummy/dummy-credentials.json",
					"$schema": "http://json-schema.org/draft-07/schema#",
					"description": "Dummy service plan credentials schema",
					"type": "object",
					"additionalProperties": false,
					"required": [
						"bar"
					],
					"properties": {
						"bar": {
							"type": "string",
							"minLength": 1
						}
					}
				}`,
			},
		},
	}
}

func (d Dummy) AdminServices() []servicesv1.Service {
	return nil
}

func (d Dummy) Reconcile(
	ctx kore.Context,
	service *servicesv1.Service,
) (reconcile.Result, error) {
	return reconcile.Result{}, nil
}

func (d Dummy) Delete(
	ctx kore.Context,
	service *servicesv1.Service,
) (reconcile.Result, error) {
	return reconcile.Result{}, nil
}

func (d Dummy) ReconcileCredentials(
	ctx kore.Context,
	service *servicesv1.Service,
	creds *servicesv1.ServiceCredentials,
) (reconcile.Result, map[string]string, error) {
	res := map[string]string{
		"superSecret": creds.Name + "-secret",
	}
	return reconcile.Result{}, res, nil
}

func (d Dummy) DeleteCredentials(
	ctx kore.Context,
	service *servicesv1.Service,
	creds *servicesv1.ServiceCredentials,
) (reconcile.Result, error) {
	return reconcile.Result{}, nil
}
