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

package aws

import (
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	awsv1alpha1 "github.com/appvia/kore/pkg/apis/aws/v1alpha1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/cloud/aws"
	"github.com/appvia/kore/pkg/utils/kubernetes"
)

// GetCredentials will obtain aws credentials from all supported sources
// If AWSAccountCredentials are not ready, will return nil
func GetCredentials(ctx kore.Context, team string, credentials corev1.Ownership) (*aws.Credentials, error) {
	// @step: is the team permitted access to this credentials
	permitted, err := ctx.Kore().Teams().Team(team).Allocations().IsPermitted(ctx, credentials)
	if err != nil {
		ctx.Logger().WithError(err).Error("attempting to check for permission on credentials")

		return nil, fmt.Errorf("attempting to check for permission on credentials")
	}

	if !permitted {
		ctx.Logger().Warn("trying to use unallocated aws credential permissions")

		return nil, errors.New("you do not have permissions to the aws credentials")
	}
	// @step: retrieve the credentials
	key := types.NamespacedName{
		Namespace: credentials.Namespace,
		Name:      credentials.Name,
	}

	// @step: are we building the cluster off a project claim or static credentials
	switch credentials.Group {
	case eks.SchemeGroupVersion.Group:
		switch kind := credentials.Kind; kind {
		case "EKSCredentials":
			return GetEKSCredentials(ctx, key)
		default:
			return nil, fmt.Errorf("unknown eks credential kind: %s", kind)
		}

	case awsv1alpha1.SchemeGroupVersion.Group:
		switch kind := credentials.Kind; kind {
		case "AWSAccountClaim":
			return GetAWSAccountClaimCredential(ctx, key)
		default:
			return nil, fmt.Errorf("unknown aws credential kind: %s", kind)
		}
	}

	return nil, fmt.Errorf("unknown credentials api group specified %s with kind %s", credentials.Group, credentials.Kind)

}

// GetEKSCredentials obtains credentials from an eks credential object / secret ref
func GetEKSCredentials(ctx kore.Context, key types.NamespacedName) (*aws.Credentials, error) {
	eksCreds := &eks.EKSCredentials{
		ObjectMeta: metav1.ObjectMeta{
			Name:      key.Name,
			Namespace: key.Namespace,
		},
	}
	found, err := kubernetes.GetIfExists(ctx, ctx.Client(), eksCreds)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("eks credentials: (%s/%s) not found", key.Namespace, key.Name)
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
	secret, err := controllers.GetDecodedSecret(ctx, ctx.Client(), eksCreds.Spec.CredentialsRef)
	if err != nil {
		return nil, err
	}

	return &aws.Credentials{
		AccountID:       eksCreds.Spec.AccountID,
		AccessKeyID:     secret.Spec.Data["access_key_id"],
		SecretAccessKey: secret.Spec.Data["access_secret_key"],
	}, nil
}

// GetAWSAccountClaimCredential retrieves aws credentials from an AWSAccountClaim object
// If the claim isn't ready, will return nil
func GetAWSAccountClaimCredential(ctx kore.Context, key types.NamespacedName) (*aws.Credentials, error) {
	c := &awsv1alpha1.AWSAccountClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      key.Name,
			Namespace: key.Namespace,
		},
	}

	found, err := kubernetes.GetIfExists(ctx, ctx.Client(), c)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("aws account claim: (%s/%s) not found", key.Namespace, key.Namespace)
	}

	// @step: we need to check the status of the project
	if c.Status.Status == corev1.FailureStatus {
		return nil, errors.New("aws account is in a failed state")
	}
	if c.Status.Status == corev1.PendingStatus {

		// No credentials yet but also not an error
		return nil, nil
	}
	if c.Status.CredentialRef == nil {
		return nil, errors.New("no aws credentials reference on account claim")
	}
	if c.Status.CredentialRef.Name == "" || c.Status.CredentialRef.Namespace == "" {
		return nil, errors.New("aws account claims credentials reference is missing fields")
	}

	// @step: we need to grab the secret
	secret, err := controllers.GetDecodedSecret(ctx, ctx.Client(), c.Status.CredentialRef)
	if err != nil {
		return nil, err
	}

	return &aws.Credentials{
		AccountID:       c.Status.AccountID,
		AccessKeyID:     secret.Spec.Data["access_key_id"],
		SecretAccessKey: secret.Spec.Data["access_secret_key"],
	}, nil
}
