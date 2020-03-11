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

package credentials

import (
	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
)

// MaxChunkSize is the largest number of permissions that can be checked in one request
const MaxChunkSize = 100

type awsClient struct {
	credentials *eks.EKSCredentials
}

// NewClient creates and returns a permissions verifier
func NewClient(credentials *eks.EKSCredentials) (*awsClient, error) {
	awsClient := &awsClient{
		credentials: credentials,
	}

	return awsClient, nil
}

// HasRequiredPermissions tests whether the IAM roles are correct for creating a cluster
func (c *awsClient) HasRequiredPermissions() (bool, error) {
	// TODO work out AWS equivalent of IAM API verification
	return true, nil
}
