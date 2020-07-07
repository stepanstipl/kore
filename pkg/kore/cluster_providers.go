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
	"fmt"
	"sync"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
)

var clusterProviders = map[string]ClusterProvider{}
var cpLock = sync.Mutex{}

func RegisterClusterProvider(provider ClusterProvider) {
	cpLock.Lock()
	defer cpLock.Unlock()

	_, exists := clusterProviders[provider.Type()]
	if exists {
		panic(fmt.Errorf("cluster provider %q was already registered", provider.Type()))
	}

	clusterProviders[provider.Type()] = provider
}

func UnregisterClusterProvider(provider ClusterProvider) {
	cpLock.Lock()
	defer cpLock.Unlock()

	delete(clusterProviders, provider.Type())
}

func GetClusterProvider(providerType string) (ClusterProvider, bool) {
	cpLock.Lock()
	defer cpLock.Unlock()

	provider, found := clusterProviders[providerType]
	return provider, found
}

func ClusterProviders() map[string]ClusterProvider {
	cpLock.Lock()
	defer cpLock.Unlock()

	res := make(map[string]ClusterProvider, len(clusterProviders))
	for k, v := range clusterProviders {
		res[k] = v
	}
	return res
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ClusterProvider
type ClusterProvider interface {
	// Type returns a unique type for the cluster provider
	Type() string
	// PlanJSONSchema returns the JSON schema for the plans belonging to this cluster
	PlanJSONSchema() string
	// DefaultPlans returns with the built-in default plans
	DefaultPlans() []configv1.Plan
	// DefaultPlanPolicy returns with the built-in default plan policy
	DefaultPlanPolicy() *configv1.PlanPolicy
	// SetComponents adds all povider-specific cluster components and updates dependencies if required
	SetComponents(Context, *clustersv1.Cluster, *ClusterComponents) error
	// BeforeComponentsUpdate runs after the components are loaded but before updated
	// The cluster components will be provided in dependency order
	BeforeComponentsUpdate(Context, *clustersv1.Cluster, *ClusterComponents) error
	// SetProviderData saves the provider data on the cluster
	// The cluster components will be provided in dependency order
	SetProviderData(Context, *clustersv1.Cluster, *ClusterComponents) error
}
