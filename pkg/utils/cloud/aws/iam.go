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
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// IamClient describes a aws session and Iam service
type IamClient struct {
	// session is the AWS session
	session *session.Session
	// svc is the iam service
	svc *iam.IAM
}

const (
	// Policies required for eks clusters:
	amazonEKSClusterPolicy = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
	amazonEKSServicePolicy = "arn:aws:iam::aws:policy/AmazonEKSServicePolicy"
	amazonEKSFargatePolicy = "arn:aws:iam::aws:policy/AmazonEKSFargatePodExecutionRolePolicy"

	// ClusterStsTrustPolicy provides the trust policy for the EKS cluster Role
	ClusterStsTrustPolicy = `{
		"Version": "2012-10-17",
		"Statement": [
		{
			"Effect": "Allow",
			"Principal": {
				"Service": "eks.amazonaws.com"
			},
			"Action": "sts:AssumeRole"
		},
		{
			"Effect": "Allow",
			"Principal": {
				"AWS": "%s"
			},
			"Action": "sts:AssumeRole"
		}]
	}`

	nodeStsTrustPolicy = `{
		"Version": "2012-10-17",
		"Statement": [
		  {
			"Effect": "Allow",
			"Principal": {
			  "Service": "ec2.amazonaws.com"
			},
			"Action": "sts:AssumeRole"
		  }
		]
	}
	`

	fargateTrustPolicy = `
	{
		"Version": "2012-10-17",
		"Statement": [
		  {
			"Effect": "Allow",
			"Principal": {
        "Service": "eks-fargate-pods.amazonaws.com"
      },
			"Action": "sts:AssumeRole"
		  }
		]
	}
	`

	// amazonEKSWorkerNodePolicy provides read-only access to Amazon EC2 Container Registry repositories.
	amazonEKSWorkerNodePolicy          = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
	amazonEC2ContainerRegistryReadOnly = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
	amazonEKSCNIPolicy                 = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
)

// NewIamClient will create a new IamClient
func NewIamClient(credentials Credentials) *IamClient {
	session := getNewSession(credentials, "")

	return &IamClient{session: session, svc: iam.New(session)}
}

// GetARN returns the ARN from the client
func (i *IamClient) GetARN() (string, error) {
	resp, err := i.svc.GetUser(&iam.GetUserInput{})
	if err != nil {
		return "", err
	}

	return aws.StringValue(resp.User.Arn), nil
}

// EnsureEKSClusterRole will return the cluster role and the nodepool role
func (i *IamClient) EnsureEKSClusterRole(ctx context.Context, prefix string) (*iam.Role, error) {
	arn, err := i.GetARN()
	if err != nil {
		return nil, err
	}
	policy := fmt.Sprintf(ClusterStsTrustPolicy, arn)
	policies := []string{
		amazonEKSClusterPolicy,
		amazonEKSServicePolicy,
	}

	return i.EnsureRole(ctx, prefix+"-eks-cluster", policies, policy)
}

// EnsureEKSNodePoolRole will create a nodepool eks role
func (i *IamClient) EnsureEKSNodePoolRole(ctx context.Context, prefix string) (*iam.Role, error) {
	policies := []string{
		amazonEKSWorkerNodePolicy,
		amazonEC2ContainerRegistryReadOnly,
		amazonEKSCNIPolicy,
	}

	return i.EnsureRole(ctx, prefix+"-eks-nodepool", policies, nodeStsTrustPolicy)
}

// EnsureFargatPodExecutionRole creates a default role for the fargat pods to run under
func (i IamClient) EnsureFargatPodExecutionRole(ctx context.Context, prefix string) (*iam.Role, error) {
	policies := []string{
		amazonEKSFargatePolicy,
	}

	return i.EnsureRole(ctx, prefix+"-fargate", policies, fargateTrustPolicy)
}

// EnsureRole is responsible for creating a role
func (i *IamClient) EnsureRole(ctx context.Context, name string, policies []string, stsPolicy string) (*iam.Role, error) {
	role, err := i.RoleExists(ctx, name)
	if err != nil {
		return nil, err
	}
	if role != nil {
		return role, nil
	}

	// @step: the role does not exist, so we must create it
	resp, err := i.svc.CreateRoleWithContext(ctx, &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(stsPolicy),
		Path:                     aws.String("/"),
		RoleName:                 aws.String(name),
	})
	if err != nil {
		return nil, err
	}

	// @step: ensure the policies are correct for the role
	lresp, err := i.svc.ListAttachedRolePoliciesWithContext(ctx, &iam.ListAttachedRolePoliciesInput{
		RoleName: resp.Role.RoleName,
	})
	if err != nil {
		return nil, err
	}

	for _, x := range policies {
		found := func() bool {
			for _, j := range lresp.AttachedPolicies {
				if aws.StringValue(j.PolicyArn) == x {
					return true
				}
			}

			return false
		}()

		if !found {
			_, err := i.svc.AttachRolePolicyWithContext(ctx, &iam.AttachRolePolicyInput{
				PolicyArn: aws.String(x),
				RoleName:  resp.Role.RoleName,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return resp.Role, nil
}

// RoleExists checks if a IAM role exists
func (i *IamClient) RoleExists(ctx context.Context, name string) (*iam.Role, error) {
	resp, err := i.svc.GetRoleWithContext(ctx, &iam.GetRoleInput{
		RoleName: aws.String(name),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				return nil, nil
			}
		}

		return nil, err
	}

	return resp.Role, nil
}
