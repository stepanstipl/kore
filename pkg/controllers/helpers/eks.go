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
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/cloud/aws"
	"github.com/appvia/kore/pkg/utils/kubernetes"
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
	creds, err := a.GetCredentials(a.resource.Namespace)
	if err != nil {
		return nil, err
	}

	client, err := aws.NewEKSClient(creds, a.resource)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// GetCredentials returns the cloud credential
func (a *KoreEKS) GetCredentials(team string) (*aws.Credentials, error) {
	// @step: is the team permitted access to this credentials
	permitted, err := a.ctx.Kore().Teams().Team(team).Allocations().IsPermitted(a.ctx, a.resource.Spec.Credentials)
	if err != nil {
		a.ctx.Logger().WithError(err).Error("attempting to check for permission on credentials")

		return nil, fmt.Errorf("attempting to check for permission on credentials")
	}

	if !permitted {
		a.ctx.Logger().Warn("trying to build eks cluster unallocated permissions")

		return nil, errors.New("you do not have permissions to the eks credentials")
	}

	// @step: retrieve the credentials
	eksCreds := &eks.EKSCredentials{
		ObjectMeta: metav1.ObjectMeta{
			Name:      a.resource.Spec.Credentials.Name,
			Namespace: a.resource.Spec.Credentials.Namespace,
		},
	}
	found, err := kubernetes.GetIfExists(a.ctx, a.ctx.Client(), eksCreds)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("eks credentials: (%s/%s) not found", a.resource.Spec.Credentials.Namespace, a.resource.Spec.Credentials.Name)
	}

	// for backwards-compatibility, use the creds set on the EKSCredentials resource, if they exist
	if eksCreds.Spec.SecretAccessKey != "" && eksCreds.Spec.AccessKeyID != "" {
		return &aws.Credentials{
			AccountID:       eksCreds.Spec.AccountID,
			AccessKeyID:     eksCreds.Spec.AccessKeyID,
			SecretAccessKey: eksCreds.Spec.SecretAccessKey,
		}, nil
	}

	// @step: we need to grab the secret
	secret, err := controllers.GetDecodedSecret(a.ctx, a.ctx.Client(), eksCreds.Spec.CredentialsRef)
	if err != nil {
		return nil, err
	}

	return &aws.Credentials{
		AccountID:       eksCreds.Spec.AccountID,
		AccessKeyID:     secret.Spec.Data["access_key_id"],
		SecretAccessKey: secret.Spec.Data["access_secret_key"],
	}, nil
}
