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
	"errors"

	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	awsc "github.com/appvia/kore/pkg/controllers/cloud/aws"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/cloud/aws"
)

// KoreEKS can provide reuse for read only operations with EKS objects across our controllers
type KoreEKS struct {
	ctx kore.Context
	// EKSResource object
	resource *eks.EKS
}

// NewKoreEKS creates a reusable EKS object across multiple controllers
func NewKoreEKS(ctx kore.Context, resource *eks.EKS) *KoreEKS {
	return &KoreEKS{
		ctx:      ctx,
		resource: resource,
	}
}

// GetClusterClient returns a EKS cluster client
func (a *KoreEKS) GetClusterClient() (*aws.Client, error) {
	// @step: first we need to check if we have access to the credentials
	creds, err := awsc.GetCredentials(a.ctx, a.resource.Namespace, a.resource.Spec.Credentials)
	if err != nil {
		return nil, err
	}
	if creds == nil {
		return nil, errors.New("unbaled to access credentials - should always check before getting client")
	}

	client, err := aws.NewEKSClient(creds, a.resource)
	if err != nil {
		return nil, err
	}

	return client, nil
}
