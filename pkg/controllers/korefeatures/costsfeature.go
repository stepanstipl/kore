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

package features

import (
	"fmt"
	"strings"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	v1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/serviceproviders/application"
)

func (c *Controller) getCostsServices(koreCtx kore.Context, config *configv1.KoreFeature) ([]v1.Service, error) {
	servicePlan, err := c.kore.ServicePlans().Get(koreCtx, "app-kore-costs")
	if err != nil {
		return nil, fmt.Errorf("failed to get service plan app-kore-costs: %w", err)
	}

	cluster := corev1.Ownership{
		Group:     clustersv1.ClusterGVK.Group,
		Version:   clustersv1.ClusterGVK.Version,
		Kind:      clustersv1.ClusterGVK.Kind,
		Namespace: "kore-admin",
		Name:      "kore",
	}

	service := application.CreateSystemServiceFromPlan(
		*servicePlan,
		cluster,
		"kore-costs",
		kore.HubAdminTeam,
	)

	if config.Spec.Configuration["gcp_credentials"] != "" {
		service.Spec.ConfigurationFrom = addSecret(
			service.Spec.ConfigurationFrom,
			config.Spec.Configuration["gcp_credentials"],
			"secrets.gcp_credentials",
			"key")
	}

	if config.Spec.Configuration["aws_credentials"] != "" {
		service.Spec.ConfigurationFrom = addSecret(
			service.Spec.ConfigurationFrom,
			config.Spec.Configuration["aws_credentials"],
			"secrets.aws_access_key",
			"access_key_id")
		service.Spec.ConfigurationFrom = addSecret(
			service.Spec.ConfigurationFrom,
			config.Spec.Configuration["aws_credentials"],
			"secrets.aws_secret_key",
			"access_secret_key")
	}

	return []v1.Service{service}, nil
}

func addSecret(cfs []corev1.ConfigurationFromSource, namespacedName string, path string, key string) []corev1.ConfigurationFromSource {
	nn := strings.Split(namespacedName, "/")
	cfs = append(cfs, corev1.ConfigurationFromSource{
		Path: path,
		SecretKeyRef: &corev1.OptionalSecretKeySelector{
			SecretKeySelector: corev1.SecretKeySelector{
				Key:       key,
				Namespace: nn[0],
				Name:      nn[1],
			},
		},
	})
	return cfs
}
