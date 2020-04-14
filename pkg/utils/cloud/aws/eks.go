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
	"context"
	"errors"
	"strings"
	"time"

	eksv1alpha1 "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrClusterNotFound indicates the cluster does not exist
	ErrClusterNotFound = errors.New("eks cluster not found")
	// ErrNodeGroupNotFound indicates the nodegroup does not exist
	ErrNodeGroupNotFound = errors.New("eks nodegroup not found")
	// ErrResourceBusy indicate the resource is currently busy performing an operation
	ErrResourceBusy = errors.New("resource is busy performing an operation (upgrade, creating)")
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
func NewEKSClient(cred *eksv1alpha1.EKSCredentials, cluster *eksv1alpha1.EKS) (*Client, error) {
	sesh, err := session.NewSession(&aws.Config{
		Region: aws.String(cluster.Spec.Region),
		Credentials: credentials.NewStaticCredentials(
			cred.Spec.AccessKeyID,
			cred.Spec.SecretAccessKey,
			"",
		),
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

// NewEKSClientFromVPC will create a new eks client from an VPCClient object
func NewEKSClientFromVPC(c *VPCClient, clusterName string) *Client {
	return &Client{
		clusterName: clusterName,
		Sess:        c.Sess,
		svc:         eks.New(c.Sess),
	}
}

// Exists checks if a cluster exists
func (c *Client) Exists(ctx context.Context) (exists bool, err error) {
	_, err = c.svc.DescribeClusterWithContext(ctx, &awseks.DescribeClusterInput{
		Name: aws.String(c.clusterName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeResourceNotFoundException:
				return false, nil
			default:
				return false, err
			}
		} else {
			return false, err
		}
	}

	return true, nil
}

// Create creates an EKS cluster
func (c *Client) Create(ctx context.Context) error {
	logger := log.WithFields(log.Fields{
		"name":      c.cluster.Name,
		"namespace": c.cluster.Namespace,
	})
	logger.Debug("attempting to create the eks cluster")

	// @step: we should check if the cluster already exist
	existing, err := c.Exists(ctx)
	if err != nil {
		logger.WithError(err).Error("trying to check for the eks cluster")

		return err
	}
	if !existing {
		_, err := c.svc.CreateClusterWithContext(ctx, c.createClusterInput())
		if err != nil {
			return err
		}
	}

	return c.WaitForClusterReady(ctx)
}

// WaitForClusterReady waits for the eks cluster be go ready
func (c *Client) WaitForClusterReady(ctx context.Context) error {
	logger := log.WithFields(log.Fields{
		"name":      c.cluster.Name,
		"namespace": c.cluster.Namespace,
	})
	logger.Debug("attempting to wait for eks cluster to be ready")

	// @step: check we are waiting for something that actually exist
	existing, err := c.Exists(ctx)
	if err != nil {
		logger.WithError(err).Error("trying to check if cluster exists to wait for")

		return err
	}
	if !existing {
		logger.Warn("no eks cluster to wait for, something went wrong here")

		return ErrClusterNotFound
	}

	// @step: wait for the cluster to be created
	return utils.WaitUntilComplete(ctx, 1*time.Hour, 30*time.Second, func() (bool, error) {
		resp, err := c.svc.DescribeClusterWithContext(ctx, &eks.DescribeClusterInput{
			Name: aws.String(c.cluster.Name),
		})
		if err != nil {
			logger.WithError(err).Error("trying to check for eks cluster status")

			return false, nil
		}

		if resp.Cluster == nil {
			logger.Warn("no cluster found in the describe response")

			return false, nil
		}

		switch aws.StringValue(resp.Cluster.Status) {
		case eks.ClusterStatusActive:
			logger.Debug("eks cluster is active and ready")

			return true, nil
		case eks.ClusterStatusFailed:
			return false, errors.New("cluster has failed to provision")
		case eks.ClusterStatusCreating:
			logger.Debug("cluster is still pending in provisioning")
		default:
			logger.Warnf("unknown cluster status: %s returned", aws.StringValue(resp.Cluster.Status))
		}

		return false, nil
	})
}

// Delete is responsible for deleting the eks cluster
func (c *Client) Delete(ctx context.Context) error {
	logger := log.WithFields(log.Fields{
		"name":      c.cluster.Name,
		"namespace": c.cluster.Namespace,
	})
	logger.Debug("attempting to delete the eks cluster")

	// @step: get the state of the cluster
	resp, err := c.svc.DescribeClusterWithContext(ctx, &eks.DescribeClusterInput{
		Name: aws.String(c.cluster.Name),
	})
	if err != nil {
		if !c.IsNotFound(err) {
			logger.WithError(err).Error("truing to describe the eks cluster")

			return err
		}

		return nil
	}

	// @step: if the cluster is not deleting, try and delete now
	switch aws.StringValue(resp.Cluster.Status) {
	case eks.ClusterStatusActive, eks.ClusterStatusFailed:
		if _, err := c.svc.DeleteClusterWithContext(ctx, &eks.DeleteClusterInput{
			Name: aws.String(c.clusterName),
		}); err != nil {
			log.WithError(err).Error("trying to delete eks cluster from aws")

			return err
		}
	case eks.ClusterStatusCreating:
		logger.Debug("eks cluster is still being created, cannot be deleted yet")

		return errors.New("eks is still being created, cannot delete yet")
	case eks.ClusterStatusUpdating:
		logger.Debug("eks cluster is still being created, cannot be deleted yet")

		return errors.New("eks is still being updated, cannot delete yet")
	}

	// @step: we need to wait for the cluster to be remove
	return utils.WaitUntilComplete(context.Background(), 1*time.Hour, 30*time.Second, func() (bool, error) {
		logger.Debug("checking if the eks cluster has been deleted")

		if found, err := c.Exists(ctx); err != nil {
			logger.WithError(err).Error("trying to check for eks cluster")

			return false, nil
		} else if !found {
			logger.Debug("eks cluster has been successfully removed")

			return true, nil
		}

		return false, nil
	})
}

// Update should migrate changes to a cluster object
func (c *Client) Update() error {
	return errors.New("not yet implimented")
}

// VerifyCredentials is responsible for verifying AWS creds
func (c *Client) VerifyCredentials() error {
	// TODO: see https://github.com/appvia/kore/issues/498

	return nil
}

// Describe returns the AWS EKS output
func (c *Client) Describe(ctx context.Context) (*eks.Cluster, error) {
	d, err := c.svc.DescribeClusterWithContext(ctx, &awseks.DescribeClusterInput{
		Name: aws.String(c.clusterName),
	})
	if err != nil {
		return nil, err
	}

	return d.Cluster, nil
}

// DeleteNodeGroup will remove a nodegroup from a cluster
func (c *Client) DeleteNodeGroup(ctx context.Context, group *eksv1alpha1.EKSNodeGroup) error {
	logger := log.WithFields(log.Fields{
		"name":      group.Name,
		"namespace": group.Namespace,
	})
	logger.Debug("attempting to delete the eks nodegroup")

	// @step: check the status of the nodegroup
	resp, err := c.svc.DescribeNodegroupWithContext(ctx, &eks.DescribeNodegroupInput{
		ClusterName:   aws.String(group.Spec.Cluster.Name),
		NodegroupName: aws.String(group.Name),
	})
	if err != nil {
		if !c.IsNotFound(err) {
			logger.WithError(err).Error("trying to describe the nodegroup")

			return err
		}

		return nil
	}

	// @step: we check the status if see if its already deleting
	switch aws.StringValue(resp.Nodegroup.Status) {
	case eks.NodegroupStatusActive, eks.NodegroupStatusCreateFailed, eks.NodegroupStatusDegraded:
		_, err := c.svc.DeleteNodegroupWithContext(ctx, &eks.DeleteNodegroupInput{
			ClusterName:   aws.String(group.Spec.Cluster.Name),
			NodegroupName: aws.String(group.Name),
		})

		if err != nil {
			if c.IsNotFound(err) {
				return nil
			}
			logger.WithError(err).Error("trying to delete the nodegroup")

			return err
		}
	case eks.NodegroupStatusCreating, eks.NodegroupStatusUpdating:
		log.Warn("trying to delete nodegroup, resource operating pending")

		return ErrResourceBusy
	case eks.NodegroupStatusDeleteFailed:
		log.Error("trying to delete the nodegroup")

		return errors.New("nodegroup has failed to delete, please check console")
	}

	// @step: we need to wait for the nodegroup to be removed
	return utils.WaitUntilComplete(context.Background(), 30*time.Minute, 30*time.Second, func() (bool, error) {
		logger.Debug("checking if the eks nodegroup has been deleted")

		if found, err := c.NodeGroupExists(ctx, group); err != nil {
			logger.WithError(err).Error("trying to check for eks nodegroup")

			return false, nil
		} else if !found {
			logger.Debug("eks nodegroup has been successfully removed")

			return true, nil
		}

		return false, nil
	})
}

// CreateNodeGroup will create a node group for the EKS cluster
func (c *Client) CreateNodeGroup(ctx context.Context, group *eksv1alpha1.EKSNodeGroup) error {
	// @step: check if the nodegroup exists already
	existing, err := c.NodeGroupExists(ctx, group)
	if err != nil {
		return err
	}
	if !existing {
		input := &eks.CreateNodegroupInput{
			AmiType:        aws.String(group.Spec.AMIType),
			ClusterName:    aws.String(group.Spec.Cluster.Name),
			NodeRole:       aws.String(group.Status.NodeIAMRole),
			ReleaseVersion: aws.String(group.Spec.ReleaseVersion),
			DiskSize:       aws.Int64(group.Spec.DiskSize),
			InstanceTypes:  aws.StringSlice([]string{group.Spec.InstanceType}),
			NodegroupName:  aws.String(group.Name),
			Subnets:        aws.StringSlice(group.Spec.Subnets),
			ScalingConfig: &eks.NodegroupScalingConfig{
				DesiredSize: aws.Int64(group.Spec.DesiredSize),
				MaxSize:     aws.Int64(group.Spec.MaxSize),
				MinSize:     aws.Int64(group.Spec.MinSize),
			},
		}
		if group.Spec.EC2SSHKey != "" {
			input.RemoteAccess = &eks.RemoteAccessConfig{
				Ec2SshKey:            aws.String(group.Spec.EC2SSHKey),
				SourceSecurityGroups: aws.StringSlice(group.Spec.SSHSourceSecurityGroups),
			}
		}
		if len(group.Spec.Tags) > 0 {
			input.Tags = aws.StringMap(group.Spec.Tags)
		}
		if len(group.Spec.Labels) > 0 {
			input.Labels = aws.StringMap(group.Spec.Labels)
		}

		if _, err := c.svc.CreateNodegroup(input); err != nil {
			return err
		}
	}

	return c.WaitForNodeGroupReady(ctx, group)
}

// WaitForNodeGroupReady is responsible for waiting for the nodegroup to provision or fail
func (c *Client) WaitForNodeGroupReady(ctx context.Context, group *eksv1alpha1.EKSNodeGroup) error {
	logger := log.WithFields(log.Fields{
		"name":      group.Name,
		"namespace": group.Namespace,
	})
	logger.Debug("attempting to wait for eks node group to be ready")

	// @step: check we are waiting for something that actually exist
	existing, err := c.NodeGroupExists(ctx, group)
	if err != nil {
		logger.WithError(err).Error("trying to check if nodegroup exists")

		return err
	}
	if !existing {
		logger.Warn("no nodegroup to wait for, something went wrong here")

		return ErrNodeGroupNotFound
	}

	// @step: wait for the nodegroup to be created
	return utils.WaitUntilComplete(ctx, 1*time.Hour, 30*time.Second, func() (bool, error) {
		resp, err := c.svc.DescribeNodegroupWithContext(ctx, &eks.DescribeNodegroupInput{
			ClusterName:   aws.String(group.Spec.Cluster.Name),
			NodegroupName: aws.String(group.Name),
		})
		if err != nil {
			if c.IsNotFound(err) {
				logger.Warn("eks nodegroup does not exist")

				return true, ErrNodeGroupNotFound
			}
			logger.WithError(err).Error("trying to check for eks nodegroup status")

			return false, nil
		}

		if resp.Nodegroup == nil {
			logger.Warn("no nodegroup found in the describe response")

			return false, nil
		}

		switch aws.StringValue(resp.Nodegroup.Status) {
		case eks.NodegroupStatusActive, eks.NodegroupStatusDegraded:
			logger.Debug("eks nodegroup is active and ready")
			return true, nil
		case eks.NodegroupStatusCreateFailed:
			return false, errors.New("nodegroup failed to provision")
		case eks.NodegroupStatusCreating, eks.NodegroupStatusUpdating:
			logger.Debug("nodegroup is still pending in provisioning")
		default:
			logger.Debugf("nodegroup status: %s returned", aws.StringValue(resp.Nodegroup.Status))
		}

		return false, nil
	})

}

// NodeGroupExists is responsible for checking if the nodegroup exists
func (c *Client) NodeGroupExists(ctx context.Context, nodegroup *eksv1alpha1.EKSNodeGroup) (exists bool, err error) {
	_, err = c.svc.DescribeNodegroupWithContext(ctx, &eks.DescribeNodegroupInput{
		ClusterName:   aws.String(nodegroup.Spec.Cluster.Name),
		NodegroupName: aws.String(nodegroup.Name),
	})
	if err != nil {
		if !c.IsNotFound(err) {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

// UpdateNodeGroup is responsible for checking for a drift and applying an update if required
func (c *Client) UpdateNodeGroup(ctx context.Context, group *eksv1alpha1.EKSNodeGroup) error {

	return nil
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

// createClusterInput is used to generate the EKS cluster definition
func (c *Client) createClusterInput() *awseks.CreateClusterInput {
	d := &awseks.CreateClusterInput{
		Name:    aws.String(c.cluster.Name),
		RoleArn: aws.String(c.cluster.Status.RoleARN),
		Version: aws.String(c.cluster.Spec.Version),
		ResourcesVpcConfig: &awseks.VpcConfigRequest{
			SecurityGroupIds:      aws.StringSlice(c.cluster.Spec.SecurityGroupIDs),
			SubnetIds:             aws.StringSlice(c.cluster.Spec.SubnetIDs),
			EndpointPublicAccess:  aws.Bool(true),
			EndpointPrivateAccess: aws.Bool(true),
		},
		Tags: map[string]*string{
			kore.Label("name"):  aws.String(c.cluster.Name),
			kore.Label("owned"): aws.String("true"),
			kore.Label("team"):  aws.String(c.cluster.Namespace),
		},
	}

	for _, x := range c.cluster.Spec.AuthorizedMasterNetworks {
		d.ResourcesVpcConfig.PublicAccessCidrs = append(d.ResourcesVpcConfig.PublicAccessCidrs, aws.String(x))
	}

	return d
}

func (c *Client) IsNotFound(err error) bool {
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == eks.ErrCodeResourceNotFoundException {
			return true
		}
	}

	return false
}

func (c *Client) IsInvalidParameterException(err error) bool {
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == eks.ErrCodeInvalidParameterException && strings.Contains(aerr.Message(), "does not exist") {
			return true
		}
	}

	return false
}
