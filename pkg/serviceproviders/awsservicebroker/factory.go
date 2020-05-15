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

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eksv1alpha1 "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/serviceproviders/openservicebroker"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	osb "github.com/kubernetes-sigs/go-open-service-broker-client/v2"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	kore.RegisterServiceProviderFactory(ProviderFactory{})
}

const (
	S3BucketTagImported  = "kore.appvia.io/initialized"
	ComponentDynamoDB    = "DynamoDB Table"
	ComponentS3Bucket    = "S3 Bucket"
	ComponentHelmRelease = "Helm Release"
	ComponentProvider    = "Provider"
)

type ProviderFactory struct{}

func (d ProviderFactory) Type() string {
	return "aws-servicebroker"
}

func (d ProviderFactory) JSONSchema() string {
	return `{
		"$id": "https://appvia.io/schemas/serviceprovider/aws-servicebroker.json",
		"$schema": "http://json-schema.org/draft-07/schema#",
		"description": "This is a custom service provider for aws-servicebroker (https://github.com/awslabs/aws-servicebroker)",
		"type": "object",
		"additionalProperties": false,
		"properties": {
			"chartRepositoryType": {
				"type": "string",
				"enum": ["git", "helm"]
			},
			"chartRepository": {
				"type": "string",
				"minLength": 1
			},
			"chartVersion": {
				"type": "string",
				"minLength": 1
			},
			"chartRepositoryRef": {
				"type": "string",
				"minLength": 1
			},
			"chartRepositoryPath": {
				"type": "string",
				"minLength": 1
			},
			"region": {
				"type": "string",
				"default": "us-east-1",
				"minLength": 1
			},
			"tableName": {
				"type": "string",
				"default": "aws-service-broker",
				"minLength": 1
			},
			"s3BucketName": {
				"type": "string",
				"default": "awsservicebroker",
				"minLength": 1
			},
			"s3BucketRegion": {
				"type": "string",
				"default": "us-east-1",
				"minLength": 1
			},
			"s3BucketKey": {
				"type": "string",
				"default": "templates/latest/"
			}
		}
	}`
}

func (d ProviderFactory) CreateProvider(ctx kore.ServiceProviderContext, serviceProvider *servicesv1.ServiceProvider) (_ kore.ServiceProvider, complete bool, _ error) {
	var config = DefaultProviderConfiguration()

	if err := serviceProvider.Spec.GetConfiguration(config); err != nil {
		return nil, false, fmt.Errorf("failed to process aws-servicebroker configuration: %w", err)
	}

	awsAccessKeyID, awsSecretAccessKey, err := getCredentials(ctx, serviceProvider)
	if err != nil {
		return nil, false, err
	}

	cfg := aws.NewConfig().
		WithCredentials(credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, "")).
		WithRegion(config.Region)

	sess := session.Must(session.NewSession(cfg))

	if err := d.ensureDynamoDBTable(sess, config); err != nil {
		serviceProvider.Status.Components.SetCondition(corev1.Component{
			Name:    ComponentDynamoDB,
			Status:  corev1.ErrorStatus,
			Message: "Failed to create or update DynamoDB table",
			Detail:  err.Error(),
		})
		return nil, false, err
	}

	serviceProvider.Status.Components.SetCondition(corev1.Component{Name: ComponentDynamoDB, Status: corev1.SuccessStatus})

	if config.S3BucketName != "awsservicebroker" {
		if err := d.ensureS3Bucket(sess, config); err != nil {
			serviceProvider.Status.Components.SetCondition(corev1.Component{
				Name:    ComponentS3Bucket,
				Status:  corev1.ErrorStatus,
				Message: "Failed to create or update S3 bucket",
				Detail:  err.Error(),
			})
			return nil, false, err
		}
	}

	serviceProvider.Status.Components.SetCondition(corev1.Component{Name: ComponentS3Bucket, Status: corev1.SuccessStatus})

	clientSecret, certsSecret, complete, err := d.ensureHelmRelease(ctx, config, serviceProvider.Name, awsAccessKeyID, awsSecretAccessKey)
	if err != nil {
		serviceProvider.Status.Components.SetCondition(corev1.Component{
			Name:    ComponentHelmRelease,
			Status:  corev1.ErrorStatus,
			Message: "Failed to deploy the aws-servicebroker Helm chart",
			Detail:  err.Error(),
		})
		return nil, false, err
	}

	if !complete {
		serviceProvider.Status.Components.SetCondition(corev1.Component{Name: ComponentHelmRelease, Status: corev1.PendingStatus})
		return nil, false, nil
	}

	serviceProvider.Status.Components.SetCondition(corev1.Component{Name: ComponentHelmRelease, Status: corev1.SuccessStatus})

	osbConfig := osb.DefaultClientConfiguration()
	osbConfig.URL = fmt.Sprintf("https://%s-aws-servicebroker.kore-serviceprovider-%s.svc", serviceProvider.Name, serviceProvider.Name)
	osbConfig.AuthConfig = &osb.AuthConfig{
		BearerConfig: &osb.BearerConfig{Token: string(clientSecret.Data["token"])},
	}
	osbConfig.CAData = certsSecret.Data["ca.crt"]

	osbClient, err := osb.NewClient(osbConfig)
	if err != nil {
		return nil, false, err
	}

	provider, err := openservicebroker.NewProvider(serviceProvider.Name, osbClient)
	if err != nil {
		serviceProvider.Status.Components.SetCondition(corev1.Component{
			Name:    ComponentProvider,
			Status:  corev1.ErrorStatus,
			Message: "Failed to initialize provider",
			Detail:  err.Error(),
		})
		return nil, false, err
	}

	serviceProvider.Status.Components.SetCondition(corev1.Component{Name: ComponentProvider, Status: corev1.SuccessStatus})

	return provider, true, nil
}

func (d ProviderFactory) TearDownProvider(ctx kore.ServiceProviderContext, serviceProvider *servicesv1.ServiceProvider) (complete bool, _ error) {
	var config = DefaultProviderConfiguration()

	if err := serviceProvider.Spec.GetConfiguration(config); err != nil {
		return false, fmt.Errorf("failed to process aws-servicebroker configuration: %w", err)
	}

	done, err := d.deleteHelmRelease(ctx, serviceProvider.Name)
	if err != nil {
		serviceProvider.Status.Components.SetCondition(corev1.Component{
			Name:    ComponentHelmRelease,
			Status:  corev1.ErrorStatus,
			Message: "Failed to delete Helm deployment",
			Detail:  err.Error(),
		})
		return false, err
	}
	if !done {
		return false, nil
	}

	serviceProvider.Status.Components.SetCondition(corev1.Component{Name: ComponentHelmRelease, Status: corev1.DeletedStatus})

	awsAccessKeyID, awsSecretAccessKey, err := getCredentials(ctx, serviceProvider)
	if err != nil {
		return false, err
	}

	cfg := aws.NewConfig().
		WithCredentials(credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, "")).
		WithRegion(config.Region)

	sess := session.Must(session.NewSession(cfg))

	if err := d.deleteDynamoDBTable(sess, config); err != nil {
		serviceProvider.Status.Components.SetCondition(corev1.Component{
			Name:    ComponentDynamoDB,
			Status:  corev1.ErrorStatus,
			Message: "Failed to delete DynamoDB table",
			Detail:  err.Error(),
		})
		return false, err
	}

	serviceProvider.Status.Components.SetCondition(corev1.Component{Name: ComponentDynamoDB, Status: corev1.DeletedStatus})

	return true, nil
}

func (d ProviderFactory) RequiredCredentialTypes() []schema.GroupVersionKind {
	return []schema.GroupVersionKind{
		{
			Group:   eksv1alpha1.GroupVersion.Group,
			Version: eksv1alpha1.GroupVersion.Version,
			Kind:    "EKSCredentials",
		},
	}
}