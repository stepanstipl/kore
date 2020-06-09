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

import "github.com/appvia/kore/pkg/serviceproviders/openservicebroker"

// ProviderConfiguration is the data type for the provider configuration
type ProviderConfiguration struct {
	ChartRepositoryType string `json:"chartRepositoryType"`
	// ChartRepository is the repository URL of the Helm chart for the aws-servicebroker
	ChartRepository string `json:"chartRepository"`
	// ChartVersion is the version of the Helm chart for the aws-servicebroker
	ChartVersion string `json:"chartVersion"`
	// chartRepositoryRef is the Git repository URL of the Helm chart for the aws-servicebroker
	ChartRepositoryRef string `json:"chartRepositoryRef"`
	// ChartRepositoryPath is the path to the chart relative to the repository root.
	ChartRepositoryPath string `json:"chartRepositoryPath"`
	// Region is the AWS region where the DynamoDB table will be created
	Region string `json:"region"`
	// TableName is the DynamoDB table name where state will be stored
	TableName string `json:"tableName"`
	// S3BucketName is the name of the S3 bucket used to store the CloudFormation templates for the service plans
	S3BucketName string `json:"s3BucketName"`
	// S3BucketRegion is the region of the S3 bucket used to store the CloudFormation templates for the service plans
	S3BucketRegion string `json:"s3BucketRegion"`
	// S3BucketKey is the path in the S3 bucket used to store the CloudFormation templates for the service plans
	S3BucketKey string `json:"s3BucketKey"`
	// AWSAccessKeyID is the AWS access key id
	AWSAccessKeyID string `json:"aws_access_key_id"`
	// AWSSecretAccessKey is the AWS secret access key
	AWSSecretAccessKey string `json:"aws_secret_access_key"`

	openservicebroker.CatalogConfiguration `json:",inline"`
}

func DefaultProviderConfiguration() *ProviderConfiguration {
	return &ProviderConfiguration{
		ChartRepositoryType: "git",
		ChartRepository:     "https://github.com/appvia/aws-servicebroker",
		ChartRepositoryPath: "packaging/helm/aws-servicebroker",
		Region:              "us-east-1",
		TableName:           "aws-service-broker",
		S3BucketName:        "awsservicebroker",
		S3BucketRegion:      "us-east-1",
		S3BucketKey:         "templates/latest/",
	}
}
