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

package awsorganization

import (
	"errors"
	"fmt"

	awsv1alpha1 "github.com/appvia/kore/pkg/apis/aws/v1alpha1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers/helpers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/cloud/aws"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	"github.com/aws/aws-sdk-go/aws/arn"
)

// EnsureCredentials is responsible for checking the credentials exist
func EnsureCredentials(ctx kore.Context, org *awsv1alpha1.AWSOrganization, conditions *corev1.Components) (*aws.Credentials, error) {
	stage := "provision"
	secret := helpers.CreateSecretRef(org.Spec.CredentialsRef.Namespace, org.Spec.CredentialsRef.Name)

	// @step: do we have credentials for the project?
	found, err := kubernetes.GetIfExists(ctx, ctx.Client(), secret)
	if err != nil {
		ctx.Logger().WithError(err).Error("trying to check for aws credentials")

		conditions.SetCondition(corev1.Component{
			Detail:  err.Error(),
			Name:    stage,
			Message: "Checking for the AWS credentials secret",
			Status:  corev1.FailureStatus,
		})

		return nil, err
	}

	// @step: if the credentials secret does not exist
	if !found {
		ctx.Logger().Debug("aws credentials do not exist")

		conditions.SetCondition(corev1.Component{
			Detail:  "no credentials supplied",
			Name:    stage,
			Message: "Invalid aws master account credentials - no secret found",
			Status:  corev1.FailureStatus,
		})

		return nil, errors.New("no credentials supplied")
	}
	parsedARN, err := arn.Parse(org.Spec.RoleARN)
	if err != nil {
		ctx.Logger().Debug("aws role cannot be parsed")

		conditions.SetCondition(corev1.Component{
			Detail:  "the role cannot be parsed",
			Name:    stage,
			Message: "Invalid aws role specified for master account",
			Status:  corev1.FailureStatus,
		})

		return nil, fmt.Errorf("unexpected invalid arn %s - %w", org.Spec.RoleARN, err)
	}
	creds, err := helpers.GetAWSCreds(secret, parsedARN.AccountID)
	if err != nil {
		ctx.Logger().Debug("unable to retrieve credential from secret")

		conditions.SetCondition(corev1.Component{
			Detail:  err.Error(),
			Name:    stage,
			Message: "Error decoding credential secret",
			Status:  corev1.FailureStatus,
		})

		return nil, err
	}

	return creds, nil
}
