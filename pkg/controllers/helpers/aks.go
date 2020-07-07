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

package helpers

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2020-06-01/containerservice"
	"github.com/Azure/go-autorest/autorest/azure/auth"

	aksv1alpha1 "github.com/appvia/kore/pkg/apis/aks/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"
)

type azureCredentials struct {
	subscriptionID string
	tenantID       string
	clientID       string
	clientSecret   string
}

// AKSHelper is a helper to manage AKS clusters
type AKSHelper struct {
	aksCluster *aksv1alpha1.AKS
}

// NewAKSHelper creates a an AKS cluster helper
func NewAKSHelper(aksCluster *aksv1alpha1.AKS) *AKSHelper {
	return &AKSHelper{
		aksCluster: aksCluster,
	}
}

// CreateClusterClient creates and AKS Cluster API client
func (a *AKSHelper) CreateClusterClient(ctx kore.Context) (containerservice.ManagedClustersClient, error) {
	// @step: first we need to check if we have access to the credentials
	creds, err := a.GetCredentials(ctx)
	if err != nil {
		return containerservice.ManagedClustersClient{}, err
	}

	config := auth.NewClientCredentialsConfig(creds.clientID, creds.clientSecret, creds.tenantID)
	authorizer, err := config.Authorizer()
	if err != nil {
		return containerservice.ManagedClustersClient{}, err
	}

	client := containerservice.NewManagedClustersClient(creds.subscriptionID)
	client.Authorizer = authorizer

	return client, nil
}

// GetCredentials returns the cloud credential
func (a *AKSHelper) GetCredentials(ctx kore.Context) (*azureCredentials, error) {
	// @step: is the team permitted access to this credentials
	permitted, err := ctx.Kore().Teams().Team(a.aksCluster.Namespace).Allocations().IsPermitted(ctx, a.aksCluster.Spec.Credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to get AKS credentials: %w", err)
	}

	if !permitted {
		return nil, fmt.Errorf("%q credentials can not be used by team %q", a.aksCluster.Spec.Credentials.Name, a.aksCluster.Namespace)
	}

	aksCreds := &aksv1alpha1.AKSCredentials{
		ObjectMeta: metav1.ObjectMeta{
			Name:      a.aksCluster.Spec.Credentials.Name,
			Namespace: a.aksCluster.Spec.Credentials.Namespace,
		},
	}
	found, err := kubernetes.GetIfExists(ctx, ctx.Client(), aksCreds)
	if err != nil {
		return nil, fmt.Errorf("failed to get AKS credentials: %w", err)
	}
	if !found {
		return nil, fmt.Errorf("AKS credentials %q not found", a.aksCluster.Spec.Credentials.Name)
	}

	// @step: we need to grab the secret
	secret, err := controllers.GetDecodedSecret(ctx, ctx.Client(), aksCreds.Spec.CredentialsRef)
	if err != nil {
		return nil, fmt.Errorf("failed to get credential secrets %q: %w", aksCreds.Spec.CredentialsRef.Name, err)
	}

	return &azureCredentials{
		subscriptionID: aksCreds.Spec.SubscriptionID,
		tenantID:       aksCreds.Spec.TenantID,
		clientID:       aksCreds.Spec.ClientID,
		clientSecret:   secret.Spec.Data["client_secret"],
	}, nil
}
