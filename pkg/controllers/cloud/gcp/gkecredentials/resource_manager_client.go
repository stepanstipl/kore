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

package gkecredentials

import (
	"context"

	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/utils"

	resourcemanager "google.golang.org/api/cloudresourcemanager/v1beta1"
	"google.golang.org/api/option"
)

// MaxChunkSize is the largest number of permissions that can be checked in one request
const MaxChunkSize = 100

// CrmClient is a permissions client
type CrmClient struct {
	crm         *resourcemanager.Service
	credentials *gke.GKECredentials
}

// NewClient creates and returns a permissions verifier
func NewClient(credentials *gke.GKECredentials) (*CrmClient, error) {
	options := []option.ClientOption{option.WithCredentialsJSON([]byte(credentials.Spec.Account))}

	crm, err := resourcemanager.NewService(context.Background(), options...)
	if err != nil {
		return nil, err
	}

	return &CrmClient{crm: crm, credentials: credentials}, nil
}

// HasRequiredPermissions tests whether a serviceaccount has the required permissions for cluster manager
func (c *CrmClient) HasRequiredPermissions() (bool, error) {
	permissions := utils.ChunkBy(requiredPermissions(), MaxChunkSize)
	for _, chunk := range permissions {
		allFound, err := c.hasPermissions(chunk)
		if err != nil {
			return false, err
		}
		if !allFound {
			return false, nil
		}
	}

	return true, nil
}

// hadPermission checks if we have all the permissions passed
func (c *CrmClient) hasPermissions(permissions []string) (bool, error) {
	request := &resourcemanager.TestIamPermissionsRequest{
		Permissions: permissions,
	}

	resp, err := c.crm.Projects.TestIamPermissions(c.credentials.Spec.Project, request).Do()
	if err != nil {
		return false, err
	}

	return len(resp.Permissions) == len(permissions), nil
}

func requiredPermissions() []string {
	return []string{
		"container.clusterRoleBindings.create",
		"container.clusterRoleBindings.get",
		"container.clusterRoles.bind",
		"container.clusterRoles.create",
		"container.clusters.create",
		"container.clusters.delete",
		"container.clusters.getCredentials",
		"container.clusters.list",
		"container.operations.get",
		"container.operations.list",
		"container.podSecurityPolicies.create",
		"container.secrets.get",
		"container.serviceAccounts.create",
		"container.serviceAccounts.get",
		"iam.serviceAccounts.actAs",
		"iam.serviceAccounts.get",
		"iam.serviceAccounts.list",
		"resourcemanager.projects.get",
	}
}
