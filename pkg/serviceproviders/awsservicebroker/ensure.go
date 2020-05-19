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
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func (d ProviderFactory) ensureDynamoDBTable(sess *session.Session, config *ProviderConfiguration) error {
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
		_, err := ddbClient.CreateTable(&dynamodb.CreateTableInput{
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("id"),
					AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
				},
				{
					AttributeName: aws.String("userid"),
					AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
				},
				{
					AttributeName: aws.String("type"),
					AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("id"),
					KeyType:       aws.String(dynamodb.KeyTypeHash),
				},
				{
					AttributeName: aws.String("userid"),
					KeyType:       aws.String(dynamodb.KeyTypeRange),
				},
			},
			GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
				{
					IndexName: aws.String("type-userid-index"),
					KeySchema: []*dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String("type"),
							KeyType:       aws.String(dynamodb.KeyTypeHash),
						},
						{
							AttributeName: aws.String("userid"),
							KeyType:       aws.String(dynamodb.KeyTypeRange),
						},
					},
					Projection: &dynamodb.Projection{
						ProjectionType:   aws.String(dynamodb.ProjectionTypeInclude),
						NonKeyAttributes: aws.StringSlice([]string{"id", "userid", "type", "locked"}),
					},
					ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(5),
						WriteCapacityUnits: aws.Int64(5),
					},
				},
			},
			LocalSecondaryIndexes: nil,
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
			TableName: aws.String(config.TableName),
		})
		if err != nil {
			return fmt.Errorf("failed to create DynamoDB table: %w", err)
		}
	}

	return nil
}

func (d ProviderFactory) ensureS3Bucket(sess *session.Session, config *ProviderConfiguration) error {
	s3Client := s3.New(sess, &aws.Config{Region: aws.String(config.S3BucketRegion)})

	exists := true
	_, err := s3Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(config.S3BucketName),
	})
	if err != nil {
		if !isAWSErr(err, s3.ErrCodeNoSuchBucket, "") && !isAWSErrRequestFailureStatusCode(err, http.StatusNotFound) {
			return fmt.Errorf("failed to get S3 bucket %s: %w", config.S3BucketName, err)
		}
		exists = false
	}

	if exists {
		tagging, err := s3Client.GetBucketTagging(&s3.GetBucketTaggingInput{
			Bucket: aws.String(config.S3BucketName),
		})
		if err != nil {
			if !isAWSErrRequestFailureStatusCode(err, http.StatusNotFound) {
				return fmt.Errorf("failed to get tag for S3 bucket %s: %w", config.S3BucketName, err)
			}
		}
		if tagging != nil {
			for _, tag := range tagging.TagSet {
				if *tag.Key == S3BucketTagImported {
					return nil
				}
			}
		}
	}

	if !exists {
		_, err := s3Client.CreateBucket(&s3.CreateBucketInput{
			ACL:    aws.String(s3.BucketCannedACLPrivate),
			Bucket: aws.String(config.S3BucketName),
		})
		if err != nil {
			return fmt.Errorf("failed to create S3 bucket %q: %w", config.S3BucketName, err)
		}
	}

	s3ClientOrig := s3.New(sess, &aws.Config{Region: aws.String("us-east-1")})

	listObjectsRes, err := s3ClientOrig.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String("awsservicebroker"),
		Prefix: aws.String(config.S3BucketKey),
	})
	if err != nil {
		return fmt.Errorf("failed to list objects from s3://awsservicebroker")
	}

	for _, obj := range listObjectsRes.Contents {
		err := func() error {
			targetObjRes, err := s3Client.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(config.S3BucketName),
				Key:    obj.Key,
			})
			if err != nil && !isAWSErr(err, s3.ErrCodeNoSuchKey, "") && !isAWSErrRequestFailureStatusCode(err, http.StatusNotFound) {
				return fmt.Errorf("failed to download S3 file %q, %w", *obj.Key, err)
			}

			if targetObjRes != nil && aws.Int64Value(targetObjRes.ContentLength) > 0 {
				return nil
			}

			sourceObjectRes, err := s3ClientOrig.GetObject(&s3.GetObjectInput{
				Bucket: aws.String("awsservicebroker"),
				Key:    obj.Key,
			})
			if err != nil {
				return fmt.Errorf("failed to download S3 file %s, %w", *obj.Key, err)
			}
			defer sourceObjectRes.Body.Close()

			content, err := ioutil.ReadAll(sourceObjectRes.Body)
			if err != nil {
				return fmt.Errorf("failed to download S3 file %s, %w", *obj.Key, err)
			}

			_, err = s3Client.PutObject(&s3.PutObjectInput{
				Bucket: aws.String(config.S3BucketName),
				Key:    obj.Key,
				Body:   bytes.NewReader(content),
			})

			if err != nil {
				return fmt.Errorf("failed to upload S3 file %s, %w", *obj.Key, err)
			}

			return nil
		}()

		if err != nil {
			return err
		}
	}

	_, err = s3Client.PutBucketTagging(&s3.PutBucketTaggingInput{
		Bucket: aws.String(config.S3BucketName),
		Tagging: &s3.Tagging{
			TagSet: []*s3.Tag{
				{
					Key:   aws.String(S3BucketTagImported),
					Value: aws.String("true"),
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to tag S3 bucket %q: %w", config.S3BucketName, err)
	}

	return nil
}

func (d ProviderFactory) ensureHelmRelease(
	ctx kore.ServiceProviderContext,
	config *ProviderConfiguration,
	name string,
	awsAccessKeyID string,
	awsSecretAccessKey string,
) (clientSecret *v1.Secret, certsSecret *v1.Secret, complete bool, _ error) {
	namespaceName := "kore-serviceprovider-" + name

	client := ctx.Client

	ns := &v1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: v1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceName,
		},
	}

	if err := kubernetes.EnsureNamespace(ctx, client, ns); err != nil {
		return nil, nil, false, fmt.Errorf("failed to create namespace %q: %w", namespaceName, err)
	}

	secretName := name + "-credentials"

	secret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespaceName,
		},
		Type: "Opaque",
		Data: map[string][]byte{
			"accesskeyid": []byte(awsAccessKeyID),
			"secretkey":   []byte(awsSecretAccessKey),
		},
	}

	_, err := kubernetes.CreateOrUpdate(ctx, client, secret)
	if err != nil {
		return nil, nil, false, fmt.Errorf("failed to create secret %q: %w", secretName, err)
	}

	var chart map[string]interface{}
	switch config.ChartRepositoryType {
	case "helm":
		chart = map[string]interface{}{
			"repository": config.ChartRepository,
			"version":    config.ChartVersion,
			"name":       "aws-servicebroker",
		}
	case "git":
		chart = map[string]interface{}{
			"git":  config.ChartRepository,
			"path": config.ChartRepositoryPath,
		}
		if config.ChartRepositoryRef != "" {
			chart["ref"] = config.ChartRepositoryRef
		}
	default:
		return nil, nil, false, fmt.Errorf("invalid chart repository type: %q", config.ChartRepositoryType)
	}

	helmRelease := &unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": "helm.fluxcd.io/v1",
		"kind":       "HelmRelease",
		"metadata": map[string]interface{}{
			"name":      name,
			"namespace": namespaceName,
		},
		"spec": map[string]interface{}{
			"releaseName": name,
			"chart":       chart,
			"values": map[string]interface{}{
				"deployClusterServiceBroker":    false,
				"deployNamespacedServiceBroker": false,
				"aws": map[string]interface{}{
					"region":         config.Region,
					"tablename":      config.TableName,
					"s3region":       config.S3BucketRegion,
					"bucket":         config.S3BucketName,
					"key":            config.S3BucketKey,
					"existingSecret": secretName,
				},
			},
			"forceUpgrade": true,
		},
	}}

	if _, err := kubernetes.CreateOrUpdate(ctx, client, helmRelease); err != nil {
		return nil, nil, false, fmt.Errorf("failed to create Helm release %q: %w", name, err)
	}

	clientSecret, err = getServiceAccountToken(ctx, client, namespaceName, name+"-aws-servicebroker-client")
	if err != nil || clientSecret == nil {
		return nil, nil, false, err
	}

	certsSecret, err = getSecret(ctx, client, namespaceName, name+"-aws-servicebroker-cert")
	if err != nil || certsSecret == nil {
		return nil, nil, false, err
	}

	return clientSecret, certsSecret, true, nil
}
