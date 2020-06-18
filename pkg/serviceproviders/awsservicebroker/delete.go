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

package awsservicebroker

import (
	"fmt"
	"net/http"

	"github.com/appvia/kore/pkg/utils/kubernetes"

	"github.com/appvia/kore/pkg/kore"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/iam"
)

func (d ProviderFactory) deleteIAMRole(sess *session.Session, config *ProviderConfiguration) error {
	iamClient := iam.New(sess)

	roleExists := true
	role, err := iamClient.GetRole(&iam.GetRoleInput{RoleName: aws.String(config.AWSIAMRoleName)})
	if err != nil {
		if !isAWSErr(err, iam.ErrCodeNoSuchEntityException, "") && !isAWSErrRequestFailureStatusCode(err, http.StatusNotFound) {
			return fmt.Errorf("failed to get IAM role %q: %w", config.AWSIAMRoleName, err)
		}
		roleExists = false
	}

	if !roleExists {
		return nil
	}

	managed := false
	for _, tag := range role.Role.Tags {
		if aws.StringValue(tag.Key) == "kore.appvia.io/managed" && aws.StringValue(tag.Value) == "true" {
			managed = true
			break
		}
	}

	if !managed {
		return nil
	}

	_, err = iamClient.DeleteRolePolicy(&iam.DeleteRolePolicyInput{
		RoleName:   aws.String(config.AWSIAMRoleName),
		PolicyName: aws.String("Main"),
	})
	if err != nil {
		if !isAWSErr(err, iam.ErrCodeNoSuchEntityException, "") && !isAWSErrRequestFailureStatusCode(err, http.StatusNotFound) {
			return fmt.Errorf("failed to delete policy on  IAM role %q: %w", config.AWSIAMRoleName, err)
		}
	}

	_, err = iamClient.DeleteRole(&iam.DeleteRoleInput{RoleName: aws.String(config.AWSIAMRoleName)})
	if err != nil {
		return fmt.Errorf("failed to delete IAM role %q: %w", config.AWSIAMRoleName, err)
	}

	return nil
}

func (d ProviderFactory) deleteDynamoDBTable(sess *session.Session, config *ProviderConfiguration) error {
	ddbClient := dynamodb.New(sess, &aws.Config{Region: aws.String(config.Region)})

	exists := true
	_, err := ddbClient.DescribeTable(&dynamodb.DescribeTableInput{TableName: aws.String(config.TableName)})
	if err != nil {
		if !isAWSErr(err, dynamodb.ErrCodeResourceNotFoundException, "") && !isAWSErrRequestFailureStatusCode(err, http.StatusNotFound) {
			return fmt.Errorf("failed to describe DynamoDB table %q: %w", config.TableName, err)
		}
		exists = false
	}

	if !exists {
		return nil
	}

	_, err = ddbClient.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(config.TableName)})
	if err != nil {
		return fmt.Errorf("failed to delete DynamoDB table %q: %w", config.TableName, err)
	}

	return nil
}

func (d ProviderFactory) deleteHelmRelease(ctx kore.Context, name string) (done bool, _ error) {
	namespaceName := "kore-serviceprovider-" + name

	helmRelease := &unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": "helm.fluxcd.io/v1",
		"kind":       "HelmRelease",
		"metadata": map[string]interface{}{
			"name":      name,
			"namespace": namespaceName,
		},
	}}

	exists, err := kubernetes.CheckIfExists(ctx, ctx.Client(), helmRelease)
	if err != nil {
		return false, fmt.Errorf("failed to get HelmRelease %q: %w", name, err)
	}

	if exists {
		err := kubernetes.DeleteIfExists(ctx, ctx.Client(), helmRelease)
		if err != nil {
			return false, fmt.Errorf("failed to delete HelmRelease %q: %w", name, err)
		}

		return false, nil
	}

	ns := &v1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: v1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceName,
		},
	}

	if err := kubernetes.DeleteIfExists(ctx, ctx.Client(), ns); err != nil {
		return false, fmt.Errorf("failed to deleted namespace %q: %w", namespaceName, err)
	}

	return true, nil
}
