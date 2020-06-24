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

package costs

import (
	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	costsv1 "github.com/appvia/kore/pkg/apis/costs/v1beta1"
	"github.com/appvia/kore/pkg/controllers/helpers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/serviceproviders/application"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type cloudinfo struct {
	config *costsv1.Cost
	ctx    kore.Context
}

func newCloudInfo(ctx kore.Context, config *costsv1.Cost) *cloudinfo {
	return &cloudinfo{
		config: config,
		ctx:    ctx,
	}
}

// IsRequired determines if we need cloudinfo in the kore cluster
func (a *cloudinfo) IsRequired() (bool, error) {
	if !a.config.Spec.Enabled {
		return false, nil
	}
	// If we have info credentials provided for any of the supported clouds, the
	// cloudinfo is required.
	if _, ok := a.config.Spec.InfoCredentials[costsv1.CostCloudProviderAmazon]; ok {
		return true, nil
	}
	if _, ok := a.config.Spec.InfoCredentials[costsv1.CostCloudProviderGoogle]; ok {
		return true, nil
	}
	// @TODO: Once Azure supported, uncomment this line:
	// if _, ok := a.config.Spec.InfoCredentials[costsv1.CostCloudProviderAzure]; ok {
	// 	return true, nil
	// }
	return false, nil
}

// Delete will remove cloudinfo if we no longer have any credentials configured.
func (a *cloudinfo) Delete() (reconcile.Result, error) {
	// @TODO: Check if all credentials have been removed and delete cloudinfo if so.
	return reconcile.Result{}, nil
}

func (a *cloudinfo) Ensure() (reconcile.Result, error) {
	// @step: Prepare configuration for the cloudinfo helm chart
	cloudInfoConfig, err := a.getCloudInfoConfig()
	if err != nil {
		return reconcile.Result{}, err
	}

	// @step: Get a reference to the kore cluster and kubernetes into which we'll deploy
	koreCluster := &clustersv1.Cluster{}
	if err := a.ctx.Client().Get(a.ctx, types.NamespacedName{Name: "kore", Namespace: "kore-admin"}, koreCluster); err != nil {
		return reconcile.Result{}, err
	}
	koreKubernetes := &clustersv1.Kubernetes{}
	if err := a.ctx.Client().Get(a.ctx, types.NamespacedName{Name: "kore", Namespace: "kore-admin"}, koreKubernetes); err != nil {
		return reconcile.Result{}, err
	}

	// @step: Prepare helm-based service for deploy
	cloudinfoService, err := helpers.GetServiceFromPlanNameAndValues(
		a.ctx,
		application.HelmAppCloudInfo,
		koreKubernetes,
		"kore-costs",
		cloudInfoConfig,
	)
	if err != nil {
		return reconcile.Result{}, err
	}

	// @step: Ensure helm chart is deployed
	return helpers.EnsureService(a.ctx, cloudinfoService, koreCluster, a.config.Status.Components)
}

func (a *cloudinfo) getCloudInfoConfig() (map[string]interface{}, error) {
	cloudInfoConfig := map[string]interface{}{
		"app": map[string]interface{}{
			"logLevel": "debug",
		},
		"image": map[string]interface{}{
			"tag": "0.12.1",
		},
		"store": map[string]interface{}{
			"redis": map[string]interface{}{
				"enabled": true,
			},
		},
		"redis": map[string]interface{}{
			"enabled": true,
		},
		"providers": map[string]interface{}{
			"google": map[string]interface{}{
				"enabled":     false,
				"credentials": "",
			},
			"amazon": map[string]interface{}{
				"enabled":   false,
				"accessKey": "",
				"secretKey": "",
			},
			"azure": map[string]interface{}{
				"enabled":        false,
				"subscriptionId": "",
				"clientId":       "",
				"clientSecret":   "",
				"tenantId":       "",
			},
		},
	}

	providers := cloudInfoConfig["providers"].(map[string]interface{})

	if creds, ok := a.config.Spec.InfoCredentials[costsv1.CostCloudProviderGoogle]; ok {
		secret, err := a.getSecret(creds)
		if err != nil {
			return nil, err
		}
		google := providers["google"].(map[string]interface{})
		google["enabled"] = true
		google["credentials"] = secret.Spec.Data["key"]
	}

	if creds, ok := a.config.Spec.InfoCredentials[costsv1.CostCloudProviderAmazon]; ok {
		secret, err := a.getSecret(creds)
		if err != nil {
			return nil, err
		}
		amazon := providers["amazon"].(map[string]interface{})
		amazon["enabled"] = true
		amazon["accessKey"] = secret.Spec.Data["access_key_id"]
		amazon["secretKey"] = secret.Spec.Data["access_secret_key"]
	}

	// @TODO: Add support for Azure credentials here once we define the layout for a kore azure secret.

	return cloudInfoConfig, nil
}

func (a *cloudinfo) getSecret(secretRef v1.SecretReference) (*configv1.Secret, error) {
	// Get credential from reference.
	secret := &configv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      a.config.Spec.InfoCredentials[costsv1.CostCloudProviderGoogle].Name,
			Namespace: a.config.Spec.InfoCredentials[costsv1.CostCloudProviderGoogle].Namespace,
		},
	}
	ref, err := client.ObjectKeyFromObject(secret)
	if err != nil {
		return nil, err
	}
	if err := a.ctx.Client().Get(a.ctx, ref, secret); err != nil {
		return nil, err
	}
	return secret, nil
}
