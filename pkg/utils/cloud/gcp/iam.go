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

	"github.com/appvia/kore/pkg/utils"

	resourcemanager "google.golang.org/api/cloudresourcemanager/v1beta1"
	"google.golang.org/api/option"
)

// MaxChunkSize is the largest number of permissions that can be checked in one request
const MaxChunkSize = 100

// CheckServiceAccountPermissions is responsible for checking the service account the permissions
func CheckServiceAccountPermissions(
	ctx context.Context,
	project string,
	serviceAccountKey string,
	list []string,
) (bool, []string, error) {
	options := option.WithCredentialsJSON([]byte(serviceAccountKey))

	// @step: we create a client from the service account first
	client, err := resourcemanager.NewService(context.Background(), options)
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
