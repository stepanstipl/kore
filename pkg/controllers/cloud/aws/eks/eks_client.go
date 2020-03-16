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
	// credentials are the eks credentials
	credentials *eksv1alpha1.EKSCredentials
	// clusterName is the eks cluster name
	clusterName string
	// sesh is the AWS session
	sesh *session.Session
	// svc is the eks service
	svc *eks.EKS
}

// NewClient gets an AWS session
func NewClient(cred *eksv1alpha1.EKSCredentials, clusterName, region string) (*eksClient, error) {
	sesh, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(cred.Spec.AccessKeyID, cred.Spec.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}
	return &eksClient{
		credentials: cred,
		clusterName: clusterName,
		sesh:        sesh,
		svc:         eks.New(sesh),
	}, err
}

// Exists checks if a cluster exists
func (c *eksClient) Exists() (exists bool, err error) {
	_, err = c.svc.DescribeCluster(&awseks.DescribeClusterInput{
		Name: aws.String(c.clusterName),
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

TODO - this looks better - wasn't in use in appvia/eks-operator code base; will swap if required...
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
func (c *eksClient) Create(cluster *eksv1alpha1.EKS) (output *eks.CreateClusterOutput, err error) {
	output, err = c.svc.CreateCluster(c.createDefinition(cluster))
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {

			// TODO - more from appvia/eks-operator say no more!
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
		Name: &c.clusterName,
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

func (c *eksClient) createDefinition(cluster *eksv1alpha1.EKS) *awseks.CreateClusterInput {
	return &awseks.CreateClusterInput{
		Name:    aws.String(cluster.Spec.Name),
		RoleArn: aws.String(cluster.Spec.RoleARN),
		Version: aws.String(cluster.Spec.Version),
		ResourcesVpcConfig: &awseks.VpcConfigRequest{
			SecurityGroupIds: aws.StringSlice(cluster.Spec.SecurityGroupIDs),
			SubnetIds:        aws.StringSlice(cluster.Spec.SubnetIDs),
		},
	}
}

// Describe returns the AWS EKS output
func (c *eksClient) describeEKS() (output *eks.DescribeClusterOutput, err error) {
	return c.svc.DescribeCluster(&awseks.DescribeClusterInput{
		Name: aws.String(c.clusterName),
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

// CreateNodeGroup will create a node group for the EKS cluster
func (c *eksClient) CreateNodeGroup(nodegroup *eksv1alpha1.EKSNodeGroup) (err error) {
	_, err = c.svc.CreateNodegroup(&eks.CreateNodegroupInput{
		AmiType:        aws.String(nodegroup.Spec.AMIType),
		ClusterName:    aws.String(nodegroup.Spec.ClusterName),
		NodeRole:       aws.String(nodegroup.Spec.NodeRole),
		ReleaseVersion: aws.String(nodegroup.Spec.ReleaseVersion),
		DiskSize:       aws.Int64(nodegroup.Spec.DiskSize),
		InstanceTypes:  aws.StringSlice(nodegroup.Spec.InstanceTypes),
		NodegroupName:  aws.String(nodegroup.Spec.NodeGroupName),
		Subnets:        aws.StringSlice(nodegroup.Spec.Subnets),
		RemoteAccess: &eks.RemoteAccessConfig{
			Ec2SshKey:            aws.String(nodegroup.Spec.EC2SSHKey),
			SourceSecurityGroups: aws.StringSlice(nodegroup.Spec.SourceSecurityGroups),
		},
		ScalingConfig: &eks.NodegroupScalingConfig{
			DesiredSize: aws.Int64(nodegroup.Spec.DesiredSize),
			MaxSize:     aws.Int64(nodegroup.Spec.MaxSize),
			MinSize:     aws.Int64(nodegroup.Spec.MinSize),
		},
		Tags:   aws.StringMap(nodegroup.Spec.Tags),
		Labels: aws.StringMap(nodegroup.Spec.Labels),
	})
	if err != nil {
		// TODO - oh my, more from appvia/eks-operator
		if aerr, ok := err.(awserr.Error); ok {
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

// NodeGroupExists TODO - looks wrong should probably list
func (c *eksClient) NodeGroupExists(nodegroup *eksv1alpha1.EKSNodeGroup) (exists bool, err error) {
	_, err = c.svc.DescribeNodegroup(&eks.DescribeNodegroupInput{
		ClusterName:   &c.clusterName,
		NodegroupName: &nodegroup.Spec.NodeGroupName,
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
			fmt.Println(err.Error())
			return false, err
		}
	}
	return true, nil
}

// Get the status of an existing node group
func (c *eksClient) GetEKSNodeGroupStatus(nodegroup *eksv1alpha1.EKSNodeGroup) (status string, err error) {
	out, err := c.svc.DescribeNodegroup(&eks.DescribeNodegroupInput{
		ClusterName:   &c.clusterName,
		NodegroupName: &nodegroup.Spec.NodeGroupName,
	})
	return *out.Nodegroup.Status, err
}
