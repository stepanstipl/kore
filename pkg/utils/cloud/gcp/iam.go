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

package gcp

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/appvia/kore/pkg/utils"

	resourcemanager "google.golang.org/api/cloudresourcemanager/v1beta1"
	"google.golang.org/api/option"
)

var (
	// ErrClientEmailMissing indicates the client email is missing in service account
	ErrClientEmailMissing = errors.New("client email missing in service account json")
)

// MaxChunkSize is the largest number of permissions that can be checked in one request
const MaxChunkSize = 100

// CreateResourceManagerClientFromServiceAccount creates a resource manager client
func CreateResourceManagerClientFromServiceAccount(sa string) (*resourcemanager.Service, error) {
	options := option.WithCredentialsJSON([]byte(sa))

	return resourcemanager.NewService(context.Background(), options)
}

// GetServiceAccountFromKeyFile extract the service account name from the service account key
func GetServiceAccountFromKeyFile(sa string) (string, bool, error) {
	value, err := GetServiceAccountKeyAttribute(sa, "client_email")

	return value, len(value) > 0, err
}

// GetServiceAccountKeyAttribute decodes the service account key and extract a value
func GetServiceAccountKeyAttribute(sa, attribute string) (string, error) {
	values := make(map[string]interface{})

	if err := json.NewDecoder(strings.NewReader(sa)).Decode(&values); err != nil {
		return "", err
	}

	value, ok := values[attribute].(string)
	if !ok {
		return "", nil
	}

	return value, nil
}

// GetServiceAccountOrganizationsIDs retrieve the organizations for a service account
func GetServiceAccountOrganizationsIDs(ctx context.Context, sa string) ([]string, error) {
	client, err := CreateResourceManagerClientFromServiceAccount(sa)
	if err != nil {
		return nil, err
	}

	resp, err := client.Organizations.List().Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	var list []string

	for _, x := range resp.Organizations {
		list = append(list, x.OrganizationId)
	}

	return list, nil
}

// CheckOrganizationRoles checks the role the service account has
func CheckOrganizationRoles(ctx context.Context, id, email string, client *resourcemanager.Service) ([]string, error) {
	resp, err := client.
		Organizations.
		GetIamPolicy("organizations/"+id, &resourcemanager.GetIamPolicyRequest{}).
		Context(ctx).
		Do()
	if err != nil {
		return nil, err
	}

	member := email
	if strings.Contains(email, "gserviceaccount.com") {
		member = "serviceAccount:" + email
	}

	var list []string
	for _, x := range resp.Bindings {
		if utils.Contains(member, x.Members) {
			list = append(list, x.Role)
		}
	}

	return list, nil
}

// CheckServiceAccountPermissions is responsible for checking the service account the permissions
func CheckServiceAccountPermissions(
	ctx context.Context,
	project string,
	serviceAccountKey string,
	list []string,
) (bool, []string, error) {

	// @step: we create a client from the service account first
	client, err := CreateResourceManagerClientFromServiceAccount(serviceAccountKey)
	if err != nil {
		return false, nil, err
	}

	return CheckPermissions(ctx, project, client, list)
}

// HasPermissions checks if we have the correct permissions
func HasPermissions(
	ctx context.Context,
	project string,
	client *resourcemanager.Service,
	list []string) (bool, error) {

	success, _, err := CheckPermissions(ctx, project, client, list)
	if err != nil {
		return false, err
	}

	return success, err
}

// CheckPermissions checks if we have the correct and which ones are missing
func CheckPermissions(
	ctx context.Context,
	project string,
	client *resourcemanager.Service,
	list []string) (bool, []string, error) {

	var missing []string

	for _, chunk := range utils.ChunkBy(list, MaxChunkSize) {
		request := &resourcemanager.TestIamPermissionsRequest{
			Permissions: chunk,
		}

		resp, err := client.Projects.TestIamPermissions(project, request).Context(ctx).Do()
		if err != nil {
			return false, nil, err
		}
		if len(resp.Permissions) != len(chunk) {
			for _, x := range chunk {
				if !utils.Contains(x, resp.Permissions) {
					missing = append(missing, x)
				}
			}
		}
	}

	return (len(missing) <= 0), missing, nil
}
