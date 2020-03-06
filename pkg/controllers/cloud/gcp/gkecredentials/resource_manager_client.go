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

type crmClient struct {
	crm         *resourcemanager.Service
	credentials *gke.GKECredentials
}

// NewClient creates and returns a permissions verifier
// @TODO move this to a common lib
func NewClient(credentials *gke.GKECredentials) (*crmClient, error) {
	options := []option.ClientOption{option.WithCredentialsJSON([]byte(credentials.Spec.Account))}

	crm, err := resourcemanager.NewService(context.Background(), options...)
	if err != nil {
		return nil, err
	}

	return &crmClient{crm: crm, credentials: credentials}, nil
}

// HasRequiredPermissions tests whether a serviceaccount has the required permissions for cluster manager
func (c *crmClient) HasRequiredPermissions() (bool, error) {
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
func (c *crmClient) hasPermissions(permissions []string) (bool, error) {
	request := &resourcemanager.TestIamPermissionsRequest{
		Permissions: permissions,
	}

	response, err := c.crm.Projects.TestIamPermissions(c.credentials.Spec.Project, request).Do()
	if err != nil {
		return false, err
	}

	return len(response.Permissions) == len(permissions), nil
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
