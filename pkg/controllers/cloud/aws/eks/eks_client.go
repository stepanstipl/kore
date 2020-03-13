/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package eks

import (
	"fmt"

	eksv1alpha1 "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	log "github.com/sirupsen/logrus"
)

type eksClient struct {
	// credentials are the gke credentials
	credentials *eksv1alpha1.EKSCredentials
	// cluster is the gke cluster
	cluster *eksv1alpha1.EKS
	// sesh is the AWS session
	sesh *session.Session
	// svc is the eks service
	svc *eks.EKS
}

// NewClient gets an AWS session
func NewClient(cred *eksv1alpha1.EKSCredentials, cluster *eksv1alpha1.EKS) (*eksClient, error) {
	sesh, err := session.NewSession(&aws.Config{
		Region:      aws.String(cluster.Spec.Region),
		Credentials: credentials.NewStaticCredentials(cred.Spec.AccessKeyID, cred.Spec.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}
	return &eksClient{
		credentials: cred,
		sesh:        sesh,
		svc:         eks.New(sesh),
	}, err
}

// CheckEKSClusterExists checks if a cluster exists
func (c *eksClient) Exists() (exists bool, err error) {
	_, err = c.svc.DescribeCluster(&awseks.DescribeClusterInput{
		Name: aws.String(c.cluster.Spec.Name),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeResourceNotFoundException:
				return false, nil
			default:
				fmt.Println(aerr.Error())
				return false, err
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
			return false, err
		}
	}
	return true, nil
}

/*

Not sure - this looks better - will swap if required...
// EKSClusterExists check that a cluster exists
func EKSClusterExists(svc *eks.EKS, clusterName string) (exists bool, err error) {
	clusterList, err := svc.ListClusters(&eks.ListClustersInput{})
	if err != nil {
		return false, err
	}
	for _, i := range clusterList.Clusters {
		if clusterName == *i {
			return true, nil
		}
	}
	return false, nil
}
*/

// Create creates an EKS cluster
func (c *eksClient) Create() (output *eks.CreateClusterOutput, err error) {
	output, err = c.svc.CreateCluster(c.createDefinition())
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {

			// TODO - say no more!
			switch aerr.Code() {
			case eks.ErrCodeResourceInUseException:
				fmt.Println(eks.ErrCodeResourceInUseException, aerr.Error())
			case eks.ErrCodeResourceLimitExceededException:
				fmt.Println(eks.ErrCodeResourceLimitExceededException, aerr.Error())
			case eks.ErrCodeInvalidParameterException:
				fmt.Println(eks.ErrCodeInvalidParameterException, aerr.Error())
			case eks.ErrCodeClientException:
				fmt.Println(eks.ErrCodeClientException, aerr.Error())
			case eks.ErrCodeServerException:
				fmt.Println(eks.ErrCodeServerException, aerr.Error())
			case eks.ErrCodeServiceUnavailableException:
				fmt.Println(eks.ErrCodeServiceUnavailableException, aerr.Error())
			case eks.ErrCodeUnsupportedAvailabilityZoneException:
				fmt.Println(eks.ErrCodeUnsupportedAvailabilityZoneException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	return
}

// DeleteEKSCluster Delete an EKS cluster
func (c *eksClient) Delete() (output *eks.DeleteClusterOutput, err error) {
	input := &eks.DeleteClusterInput{
		Name: &c.cluster.Name,
	}
	output, err = c.svc.DeleteCluster(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeResourceInUseException:
				fmt.Println(eks.ErrCodeResourceInUseException, aerr.Error())
			case eks.ErrCodeResourceNotFoundException:
				fmt.Println(eks.ErrCodeResourceNotFoundException, aerr.Error())
			case eks.ErrCodeClientException:
				fmt.Println(eks.ErrCodeClientException, aerr.Error())
			case eks.ErrCodeServerException:
				fmt.Println(eks.ErrCodeServerException, aerr.Error())
			case eks.ErrCodeServiceUnavailableException:
				fmt.Println(eks.ErrCodeServiceUnavailableException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	return
}

// VerifyCredentials is responsible for verifying AWS creds
func (c *eksClient) VerifyCredentials() error {
	return nil
}

func (c *eksClient) createDefinition() *awseks.CreateClusterInput {
	return &awseks.CreateClusterInput{
		Name:    aws.String(c.cluster.Spec.Name),
		RoleArn: aws.String(c.cluster.Spec.RoleARN),
		Version: aws.String(c.cluster.Spec.Version),
		ResourcesVpcConfig: &awseks.VpcConfigRequest{
			SecurityGroupIds: aws.StringSlice(c.cluster.Spec.SecurityGroupIDs),
			SubnetIds:        aws.StringSlice(c.cluster.Spec.SubnetIDs),
		},
	}
}

// Describe returns the AWS EKS output
func (c *eksClient) describeEKS() (output *eks.DescribeClusterOutput, err error) {
	return c.svc.DescribeCluster(&awseks.DescribeClusterInput{
		Name: aws.String(c.cluster.Spec.Name),
	})
}

// GetEKSClusterStatus gest a cluster status
// TODO - shouldn't be public should return apis types...
func (c *eksClient) GetEKSClusterStatus() (status string, err error) {
	cluster, err := c.describeEKS()
	return *cluster.Cluster.Status, err
}

// ListEKSClusters lists all EKS clusters
func (c *eksClient) ListEKSClusters(input *eks.ListClustersInput) (output *eks.ListClustersOutput, err error) {
	output, err = c.svc.ListClusters(input)
	return output, err
}
