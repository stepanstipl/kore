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

	eksv1alpha1 "github.com/appvia/kore/pkg/apis/eks/v1alpha1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	log "github.com/sirupsen/logrus"
)

// Client for aws EKS and EKS nodegroups
type Client struct {
	// credentials are the eks credentials
	credentials *eksv1alpha1.EKSCredentials
	// cluster is the API object used
	cluster *eksv1alpha1.EKS
	// clusterName is the eks cluster name
	clusterName string
	// Sess is the AWS session
	Sess *session.Session
	// svc is the eks service
	svc *eks.EKS
}

// NewBasicClient gets an AWS session relating to a cluster
// TODO: maybe remove after refactor of nodegroup to use clusterref?
func NewBasicClient(cred *eksv1alpha1.EKSCredentials, clusterName, region string) (*Client, error) {
	sesh, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(cred.Spec.AccessKeyID, cred.Spec.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		credentials: cred,
		clusterName: clusterName,
		Sess:        sesh,
		svc:         eks.New(sesh),
	}, err
}

// NewClient gets an AWS and cluster session with a reference to our API object
func NewClient(cred *eksv1alpha1.EKSCredentials, cluster *eksv1alpha1.EKS) (*Client, error) {
	sesh, err := session.NewSession(&aws.Config{
		Region:      aws.String(cluster.Spec.Region),
		Credentials: credentials.NewStaticCredentials(cred.Spec.AccessKeyID, cred.Spec.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		credentials: cred,
		clusterName: cluster.Name,
		cluster:     cluster,
		Sess:        sesh,
		svc:         eks.New(sesh),
	}, err
}

// Exists checks if a cluster exists
func (c *Client) Exists() (exists bool, err error) {
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

// Create creates an EKS cluster
func (c *Client) Create() (*eks.CreateClusterOutput, error) {
	output, err := c.svc.CreateCluster(c.createClusterInput())
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Debugf("unhandled aws api error %s", aerr.Error())
		} else {
			log.WithError(err).Debug("generic error")
		}
		return nil, err
	}
	return output, nil
}

// Delete Delete an EKS cluster
func (c *Client) Delete() (*eks.DeleteClusterOutput, error) {
	input := &eks.DeleteClusterInput{
		Name: &c.clusterName,
	}
	output, err := c.svc.DeleteCluster(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Debugf("unhandled aws api error %s", aerr.Error())
		} else {
			log.WithError(err).Debug("generic error")
		}
		return nil, err
	}
	return output, nil
}

// Update should migrate changes to a cluster object
func (c *Client) Update() error {
	return errors.New("not yet implimented")
}

// VerifyCredentials is responsible for verifying AWS creds
func (c *Client) VerifyCredentials() error {
	return nil
}

func (c *Client) createClusterInput() *awseks.CreateClusterInput {
	return &awseks.CreateClusterInput{
		Name:    aws.String(c.cluster.Name),
		RoleArn: aws.String(c.cluster.Spec.RoleARN),
		Version: aws.String(c.cluster.Spec.Version),
		ResourcesVpcConfig: &awseks.VpcConfigRequest{
			SecurityGroupIds: aws.StringSlice(c.cluster.Spec.SecurityGroupIDs),
			SubnetIds:        aws.StringSlice(c.cluster.Spec.SubnetIDs),
		},
	}
}

// DescribeEKS returns the AWS EKS output
func (c *Client) DescribeEKS() (*eks.Cluster, error) {
	d, err := c.svc.DescribeCluster(&awseks.DescribeClusterInput{
		Name: aws.String(c.clusterName),
	})
	if err != nil {
		return nil, err
	}
	return d.Cluster, nil
}

// DeleteNodeGroup will remove a nodegroup from a cluster
func (c *Client) DeleteNodeGroup(nodegroup *eksv1alpha1.EKSNodeGroup) error {
	_, err := c.svc.DeleteNodegroup(&eks.DeleteNodegroupInput{
		ClusterName:   &c.clusterName,
		NodegroupName: &nodegroup.Name,
	})

	return err
}

// CreateNodeGroup will create a node group for the EKS cluster
func (c *Client) CreateNodeGroup(nodegroup *eksv1alpha1.EKSNodeGroup) (err error) {
	input := &eks.CreateNodegroupInput{
		AmiType:        aws.String(nodegroup.Spec.AMIType),
		ClusterName:    aws.String(nodegroup.Spec.Cluster.Name),
		NodeRole:       aws.String(nodegroup.Spec.NodeIAMRole),
		ReleaseVersion: aws.String(nodegroup.Spec.ReleaseVersion),
		DiskSize:       aws.Int64(nodegroup.Spec.DiskSize),
		InstanceTypes:  aws.StringSlice([]string{nodegroup.Spec.InstanceType}),
		NodegroupName:  aws.String(nodegroup.Name),
		Subnets:        aws.StringSlice(nodegroup.Spec.Subnets),
		RemoteAccess: &eks.RemoteAccessConfig{
			Ec2SshKey:            aws.String(nodegroup.Spec.EC2SSHKey),
			SourceSecurityGroups: aws.StringSlice(nodegroup.Spec.SSHSourceSecurityGroups),
		},
		ScalingConfig: &eks.NodegroupScalingConfig{
			DesiredSize: aws.Int64(nodegroup.Spec.DesiredSize),
			MaxSize:     aws.Int64(nodegroup.Spec.MaxSize),
			MinSize:     aws.Int64(nodegroup.Spec.MinSize),
		},
	}
	if len(nodegroup.Spec.Tags) > 0 {
		input.Tags = aws.StringMap(nodegroup.Spec.Tags)
	}
	if len(nodegroup.Spec.Labels) > 0 {
		input.Labels = aws.StringMap(nodegroup.Spec.Labels)
	}
	_, err = c.svc.CreateNodegroup(input)
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
func (c *Client) NodeGroupExists(nodegroup *eksv1alpha1.EKSNodeGroup) (exists bool, err error) {
	_, err = c.svc.DescribeNodegroup(&eks.DescribeNodegroupInput{
		ClusterName:   &c.clusterName,
		NodegroupName: &nodegroup.Name,
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

// ListNodeGroups get a list of the nodegroups
func (c *Client) ListNodeGroups() (nodegroups []string, err error) {
	nodegroups = make([]string, 0)
	ngo, err := c.svc.ListNodegroups(&eks.ListNodegroupsInput{
		ClusterName: &c.clusterName,
	})
	if err != nil {
		return nodegroups, err
	}
	for _, ng := range ngo.Nodegroups {
		nodegroups = append(nodegroups, *ng)
	}
	return nodegroups, nil
}

// GetEKSNodeGroupStatus the status of an existing node group
func (c *Client) GetEKSNodeGroupStatus(nodegroup *eksv1alpha1.EKSNodeGroup) (status string, err error) {
	out, err := c.svc.DescribeNodegroup(&eks.DescribeNodegroupInput{
		ClusterName:   &c.clusterName,
		NodegroupName: &nodegroup.Name,
	})
	return *out.Nodegroup.Status, err
}
