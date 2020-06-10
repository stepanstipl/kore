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
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"
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

	autoscalerNodeGroupAGSAccessPolicy = `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"autoscaling:DescribeAutoScalingInstances",
					"autoscaling:DescribeAutoScalingGroups",
					"autoscaling:DescribeTags",
					"autoscaling:DescribeLaunchConfigurations"
				],
				"Resource": "*"
			},
			{
				"Effect": "Allow",
				"Action": [
					"autoscaling:SetDesiredCapacity",
					"autoscaling:TerminateInstanceInAutoScalingGroup"
				],
				"Resource": "%s"
			}
		]
	}`

	autoscalerTrustPolicy = `{
		"Version": "2012-10-17",
		"Statement": [
			{
			"Effect": "Allow",
			"Principal": {
				"Federated": "arn:aws:iam::%s:oidc-provider/%s"
			},
			"Action": "sts:AssumeRoleWithWebIdentity",
				"Condition": {
					"StringEquals": {
						"%s:sub": "system:serviceaccount:kube-system:cluster-autoscaler"
					}
				}
			}
		]
	}`

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

// EnsureIRSA will enable IRSA IAM Roles for Service Accounts for an EKS cluster
func (i *IamClient) EnsureIRSA(clusterARN, oidcIssuerURL string) error {
	parsedARN, err := arn.Parse(clusterARN)
	if err != nil {
		return errors.Wrapf(err, "unexpected invalid ARN: %q", clusterARN)
	}
	switch parsedARN.Partition {
	case "aws", "aws-cn", "aws-us-gov":
	default:
		return fmt.Errorf("unknown EKS ARN: %q", clusterARN)
	}
	oidc, err := NewOpenIDConnectManager(i.svc, parsedARN.AccountID, oidcIssuerURL, parsedARN.Partition)
	if err != nil {
		return err
	}
	providerExists, err := oidc.CheckProviderExists()
	if err != nil {
		return err
	}

	if !providerExists {
		if err := oidc.CreateProvider(); err != nil {
			return err
		}
	}
	return nil
}

// GetEKSRoleName returns the name of a EKS iam role
func (i *IamClient) GetEKSRoleName(prefix string) string {
	return prefix + "-eks-cluster"
}

// GetEKSNodeGroupRoleName returns the role name
func (i *IamClient) GetEKSNodeGroupRoleName(prefix string) string {
	return prefix + "-eks-nodepool"
}

// GetEKSNodeGroupAutoscalingPolicyName returns a policy name for a nodegroup
func (i *IamClient) GetEKSNodeGroupAutoscalingPolicyName(prefix string) string {
	return prefix + "-eks-ng-autoscaling"
}

// GetNodeGroupAutoscalingPolicyArn returns a policy ARN name for a nodegroup
func (i *IamClient) GetNodeGroupAutoscalingPolicyArn(accountID, prefix string) string {
	return fmt.Sprintf("arn:aws:iam::%s:policy/%s", accountID, i.GetEKSNodeGroupAutoscalingPolicyName(prefix))
}

// GetNodeGroupAutoscalingRoleName returns an IAM Role name for the autoscaller
func (i *IamClient) GetNodeGroupAutoscalingRoleName(prefix string) string {
	return prefix + "-eks-autoscaling"
}

// DeleteEKSClutserRole is responsible for deleting the eks iam role
func (i *IamClient) DeleteEKSClutserRole(ctx context.Context, prefix string) error {
	return i.DeleteRole(ctx, i.GetEKSRoleName(prefix))
}

// DeleteEKSNodeGroupRole is responsible for removing any iam role associated to the nodegroup
func (i *IamClient) DeleteEKSNodeGroupRole(ctx context.Context, prefix string) error {
	return i.DeleteRole(ctx, i.GetEKSNodeGroupRoleName(prefix))
}

// DeleteRole is responsible for deleting a role
func (i *IamClient) DeleteRole(ctx context.Context, rolename string) error {
	role, err := i.RoleExists(ctx, rolename)
	if err != nil {
		return err
	}
	if role == nil {
		return nil
	}

	list, err := i.svc.ListAttachedRolePoliciesWithContext(ctx, &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(rolename),
	})
	if err != nil {
		return err
	}
	for _, x := range list.AttachedPolicies {
		_, err := i.svc.DetachRolePolicyWithContext(ctx, &iam.DetachRolePolicyInput{
			RoleName:  aws.String(rolename),
			PolicyArn: x.PolicyArn,
		})
		if err != nil {
			return err
		}
	}

	_, err = i.svc.DeleteRoleWithContext(ctx, &iam.DeleteRoleInput{
		RoleName: aws.String(rolename),
	})

	return err
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
	roleName := i.GetEKSRoleName(prefix)

	return i.EnsureRole(ctx, roleName, policies, policy)
}

// EnsureEKSNodePoolRole will create a nodepool eks role
func (i *IamClient) EnsureEKSNodePoolRole(ctx context.Context, prefix string) (*iam.Role, error) {
	policies := []string{
		amazonEKSWorkerNodePolicy,
		amazonEC2ContainerRegistryReadOnly,
		amazonEKSCNIPolicy,
	}
	name := i.GetEKSNodeGroupRoleName(prefix)

	return i.EnsureRole(ctx, name, policies, nodeStsTrustPolicy)
}

// EnsureRole is responsible for creating a role
func (i *IamClient) EnsureRole(ctx context.Context, name string, policies []string, stsPolicy string) (*iam.Role, error) {
	role, err := i.RoleExists(ctx, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		// @step: the role does not exist, so we must create it
		resp, err := i.svc.CreateRoleWithContext(ctx, &iam.CreateRoleInput{
			AssumeRolePolicyDocument: aws.String(stsPolicy),
			Path:                     aws.String("/"),
			RoleName:                 aws.String(name),
		})
		if err != nil {
			return nil, err
		}
		role = resp.Role
	}
	// @step: ensure the policies are correct for the role
	lresp, err := i.svc.ListAttachedRolePoliciesWithContext(ctx, &iam.ListAttachedRolePoliciesInput{
		RoleName: role.RoleName,
	})
	if err != nil {
		return nil, err
	}

	for _, x := range policies {
		fmt.Printf("checking for policy %s - ", x)
		found := func() bool {
			for _, j := range lresp.AttachedPolicies {
				if aws.StringValue(j.PolicyArn) == x {
					return true
				}
			}

			return false
		}()

		if !found {
			fmt.Printf("attaching policy %s to %s - ", x, *role.RoleName)
			_, err := i.svc.AttachRolePolicyWithContext(ctx, &iam.AttachRolePolicyInput{
				PolicyArn: aws.String(x),
				RoleName:  role.RoleName,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return role, nil
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

// DeleteClusterAutoscalerRoleAndPolicies will remove all IAM objetcs relating to this clusters autoscaler
func (i *IamClient) DeleteClusterAutoscalerRoleAndPolicies(ctx context.Context, clusterName string, ngas []NodeGroupAutoScaler) error {
	ar := i.GetNodeGroupAutoscalingRoleName(clusterName)

	// Delete role first
	err := i.DeleteRole(ctx, ar)
	if err != nil {
		return fmt.Errorf("unable to delete role %s for cluster autoscaler - %s", ar, err)
	}
	// Delete unused policies
	for _, nga := range ngas {
		err := i.DeleteNodeGroupAutoscalingPolicy(ctx, nga.NodeGroupName, nga.AutoScalingARN)
		if err != nil {
			return fmt.Errorf("unable to delete policy for nodegroup %s with asg %s", nga.NodeGroupName, nga.AutoScalingARN)
		}
	}
	return nil
}

// EnsureClusterAutoscalingRoleAndPolicies creates an IAM role for the Autoscaling workload and the relevant policies
func (i *IamClient) EnsureClusterAutoscalingRoleAndPolicies(ctx context.Context, clusterName, accountID, oidcProvider string, ngas []NodeGroupAutoScaler) (*iam.Role, error) {
	polARNs := []string{}

	for _, nga := range ngas {
		pol, err := i.EnsureNodeGroupAutoscalingPolicy(ctx, nga.NodeGroupName, nga.AutoScalingARN)
		if err != nil {
			return nil, fmt.Errorf("unable to create or update policy for nodegroup %s with asg ARN %s - %s", nga.NodeGroupName, nga.AutoScalingARN, err)
		}
		polARNs = append(polARNs, *pol.Arn)
	}
	// create role with policies and sts document.
	stsPol := fmt.Sprintf(autoscalerTrustPolicy, accountID, oidcProvider, oidcProvider)

	roleName := i.GetNodeGroupAutoscalingRoleName(clusterName)
	role, err := i.EnsureRole(ctx, roleName, polARNs, stsPol)
	if err != nil {
		return nil, fmt.Errorf("cannot create role %s for cluster %s required for autoscaling group - %s", roleName, clusterName, err)
	}
	return role, nil
}

// DeleteNodeGroupAutoscalingPolicy will delete a nodegroup Autoscaling policy
func (i *IamClient) DeleteNodeGroupAutoscalingPolicy(ctx context.Context, nodeGroupName, asgArn string) error {
	// Get ASG details from Auto Scaling Group ARN
	asg, err := GetASGDetailsFromArn(asgArn)
	if err != nil {
		return err
	}
	policyARN := i.GetNodeGroupAutoscalingPolicyArn(asg.ARN.AccountID, nodeGroupName)
	exists, _, err := i.PolicyExists(policyARN)
	if err != nil {
		return err
	}
	if exists {
		i.svc.DeletePolicy(&iam.DeletePolicyInput{
			PolicyArn: aws.String(policyARN),
		})
		if err != nil {
			return fmt.Errorf("unable to delete policy %s for nodepgroup %s - %s", policyARN, nodeGroupName, err)
		}
	}
	return nil
}

// EnsureNodeGroupAutoscalingPolicy will create or update an access policy for a specific nodegroup
func (i *IamClient) EnsureNodeGroupAutoscalingPolicy(ctx context.Context, nodeGroupName, asgArn string) (*iam.Policy, error) {

	// Get ASG details from Auto Scaling Group ARN
	asg, err := GetASGDetailsFromArn(asgArn)
	if err != nil {
		return nil, err
	}
	policyARN := i.GetNodeGroupAutoscalingPolicyArn(asg.ARN.AccountID, nodeGroupName)
	exists, p, err := i.PolicyExists(policyARN)
	if err != nil {
		return nil, err
	}
	policyName := i.GetEKSNodeGroupAutoscalingPolicyName(nodeGroupName)
	if !exists {
		// create the required policy for this nodegroup
		po, err := i.svc.CreatePolicy(&iam.CreatePolicyInput{
			Description:    aws.String(fmt.Sprintf("A policy to enable autoscaling for the nodegroup %s with the autoscaling name %s", nodeGroupName, asg.Name)),
			PolicyDocument: aws.String(fmt.Sprintf(autoscalerNodeGroupAGSAccessPolicy, asgArn)),
			PolicyName:     aws.String(policyName),
		})
		if err != nil {
			return nil, fmt.Errorf("error creating policy %s - %s", policyName, err)
		}
		p = po.Policy
	}
	return p, nil
}

// PolicyExists will determine if an AWS IAM policy exists and return the policy if it does
func (i *IamClient) PolicyExists(policyARN string) (bool, *iam.Policy, error) {
	// Get IAM policy
	po, err := i.svc.GetPolicy(&iam.GetPolicyInput{
		PolicyArn: aws.String(policyARN),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				return false, nil, nil
			}
		}
		return false, nil, err
	}
	return true, po.Policy, nil
}
