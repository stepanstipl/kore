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
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// IamClient describes a aws session and Iam service
type IamClient struct {
	// session is the AWS session
	session *session.Session
	// svc is the iam service
	svc *iam.IAM
	// accountID is the AWS account ID
	accountID string
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

	clusterAutoscalerNodeGroupAGSAccessPolicy = `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"autoscaling:DescribeAutoScalingInstances",
					"autoscaling:DescribeAutoScalingGroups",
					"autoscaling:DescribeTags",
					"autoscaling:DescribeLaunchConfigurations",
					"ec2:DescribeLaunchTemplateVersions"
				],
				"Resource": "*"
			},
			{
				"Effect": "Allow",
				"Action": [
					"autoscaling:SetDesiredCapacity",
					"autoscaling:TerminateInstanceInAutoScalingGroup"
				],
				"Resource": "*",
				"Condition": {
					"StringEquals": { "autoscaling:ResourceTag/k8s.io/cluster-autoscaler/%s": "owned" }
				}
			}
		]
	}`

	clusterAutoscalerTrustPolicy = `{
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

// NewIamClientFromSession from current session
func NewIamClientFromSession(session *session.Session) *IamClient {
	return &IamClient{session: session, svc: iam.New(session)}
}

func (i *IamClient) GetAWSAccountID() (string, error) {
	if i.accountID != "" {
		return i.accountID, nil
	}

	stsClient := sts.New(i.session)
	identity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", fmt.Errorf("failed to get AWS caller identity: %w", err)
	}
	i.accountID = *identity.Account

	return i.accountID, nil
}

// GetARN returns the ARN from the client
func (i *IamClient) GetARN() (string, error) {
	resp, err := i.svc.GetUser(&iam.GetUserInput{})
	if err != nil {
		return "", err
	}

	return aws.StringValue(resp.User.Arn), nil
}

// GetAccountIDFromARN will return the account ID from an ARN
func (i *IamClient) GetAccountIDFromARN(resARN string) (string, error) {
	parsedARN, err := arn.Parse(resARN)
	if err != nil {
		return "", errors.Wrapf(err, "unexpected invalid ARN: %q", resARN)
	}
	return parsedARN.AccountID, nil
}

func (i *IamClient) createOpenIDConnectManager(clusterARN, oidcIssuerURL string) (*OpenIDConnectManager, error) {
	parsedARN, err := arn.Parse(clusterARN)
	if err != nil {
		return nil, errors.Wrapf(err, "unexpected invalid ARN: %q", clusterARN)
	}
	switch parsedARN.Partition {
	case "aws", "aws-cn", "aws-us-gov":
	default:
		return nil, fmt.Errorf("unknown EKS ARN: %q", clusterARN)
	}
	return NewOpenIDConnectManager(i.svc, parsedARN.AccountID, oidcIssuerURL, parsedARN.Partition)
}

// EnsureOIDCProvider will create an OIDC provider to enable IAM Roles for Service Accounts for an EKS cluster
func (i *IamClient) EnsureOIDCProvider(clusterARN, oidcIssuerURL string) error {
	oidc, err := i.createOpenIDConnectManager(clusterARN, oidcIssuerURL)
	if err != nil {
		return err
	}

	providerExists, err := oidc.CheckProviderExists()
	if err != nil {
		return err
	}

	if providerExists {
		return nil
	}

	return oidc.CreateProvider()
}

// DeleteOIDCProvider will delete the OIDC provider which was used by the EKS cluster
func (i *IamClient) DeleteOIDCProvider(clusterARN, oidcIssuerURL string) error {
	oidc, err := i.createOpenIDConnectManager(clusterARN, oidcIssuerURL)
	if err != nil {
		return err
	}

	providerExists, err := oidc.CheckProviderExists()
	if err != nil {
		return err
	}

	if !providerExists {
		return nil
	}

	return oidc.DeleteProvider()
}

// GetEKSRoleName returns the name of a EKS iam role
func (i *IamClient) GetEKSRoleName(prefix string) string {
	return prefix + "-eks-cluster"
}

// GetEKSNodeGroupRoleName returns the role name
func (i *IamClient) GetEKSNodeGroupRoleName(prefix string) string {
	return prefix + "-eks-nodepool"
}

// DeleteEKSClusterRole is responsible for deleting the eks iam role
func (i *IamClient) DeleteEKSClusterRole(ctx context.Context, prefix string) error {
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
		log.Debugf("role %s doesnot exist", rolename)
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

// EnsureIAMRoleWithEmbeddedPolicy creates an IAM role with an embedded policy
func (i *IamClient) EnsureIAMRoleWithEmbeddedPolicy(name, description, trustPolicy, rolePolicy string) (*iam.Role, error) {
	var role *iam.Role

	roleExists := true
	roleOutput, err := i.svc.GetRole(&iam.GetRoleInput{RoleName: aws.String(name)})
	if err != nil {
		if !IsAWSErr(err, iam.ErrCodeNoSuchEntityException, "") && !IsAWSErrRequestFailureStatusCode(err, http.StatusNotFound) {
			return nil, fmt.Errorf("failed to get IAM role %q: %w", name, err)
		}
		roleExists = false
	}

	if roleExists {
		role = roleOutput.Role
		managed := true
		for _, tag := range roleOutput.Role.Tags {
			if aws.StringValue(tag.Key) == "kore.appvia.io/managed" && aws.StringValue(tag.Value) == "false" {
				managed = false
				break
			}
		}

		if !managed {
			return role, nil
		}
	}

	if !roleExists {
		roleOutput, err := i.svc.CreateRole(&iam.CreateRoleInput{
			RoleName:                 aws.String(name),
			Description:              aws.String(description),
			AssumeRolePolicyDocument: aws.String(trustPolicy),
			Tags: []*iam.Tag{
				{
					Key:   aws.String("kore.appvia.io/managed"),
					Value: aws.String("true"),
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create IAM role %q: %w", name, err)
		}
		role = roleOutput.Role
	} else if aws.StringValue(roleOutput.Role.AssumeRolePolicyDocument) != trustPolicy {
		_, err := i.svc.UpdateAssumeRolePolicy(&iam.UpdateAssumeRolePolicyInput{
			RoleName:       aws.String(name),
			PolicyDocument: aws.String(trustPolicy),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update IAM role %q: %w", name, err)
		}
	}

	_, err = i.svc.PutRolePolicy(&iam.PutRolePolicyInput{
		RoleName:       aws.String(name),
		PolicyName:     aws.String("Main"),
		PolicyDocument: aws.String(rolePolicy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add policy to IAM role %q: %w", name, err)
	}

	return role, nil
}

// DeleteIAMRoleWithEmbeddedPolicy deletes an IAM role which has one embedded policy called "Main"
func (i *IamClient) DeleteIAMRoleWithEmbeddedPolicy(name string) error {
	roleExists := true
	role, err := i.svc.GetRole(&iam.GetRoleInput{RoleName: aws.String(name)})
	if err != nil {
		if !IsAWSErr(err, iam.ErrCodeNoSuchEntityException, "") && !IsAWSErrRequestFailureStatusCode(err, http.StatusNotFound) {
			return fmt.Errorf("failed to get IAM role %q: %w", name, err)
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

	_, err = i.svc.DeleteRolePolicy(&iam.DeleteRolePolicyInput{
		RoleName:   aws.String(name),
		PolicyName: aws.String("Main"),
	})
	if err != nil {
		if !IsAWSErr(err, iam.ErrCodeNoSuchEntityException, "") && !IsAWSErrRequestFailureStatusCode(err, http.StatusNotFound) {
			return fmt.Errorf("failed to delete policy on  IAM role %q: %w", name, err)
		}
	}

	_, err = i.svc.DeleteRole(&iam.DeleteRoleInput{RoleName: aws.String(name)})
	if err != nil {
		return fmt.Errorf("failed to delete IAM role %q: %w", name, err)
	}

	return nil
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
				RoleName:  role.RoleName,
			})
			if err != nil {
				return nil, err
			}
		}
	}
	// Ensure the reverse - we must remove any policies that are not specified
	for _, ap := range lresp.AttachedPolicies {
		found := func() bool {
			for _, p := range policies {
				if aws.StringValue(ap.PolicyArn) == p {
					return true
				}
			}

			return false
		}()

		if !found {
			_, err := i.svc.DetachRolePolicyWithContext(ctx, &iam.DetachRolePolicyInput{
				PolicyArn: ap.PolicyArn,
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

func (i *IamClient) getClusterAutoscalerRoleName(clusterName string) string {
	return clusterName + "-eks-autoscaling"
}

// DeleteClusterAutoscalerRole will delete the cluster autoscaling IAM role
func (i *IamClient) DeleteClusterAutoscalerRole(clusterName string) error {
	name := i.getClusterAutoscalerRoleName(clusterName)

	return i.DeleteIAMRoleWithEmbeddedPolicy(name)
}

// EnsureClusterAutoscalerRole creates an IAM role for the cluster autoscaler
func (i *IamClient) EnsureClusterAutoscalerRole(clusterName, oidcProvider string) (*iam.Role, error) {
	name := i.getClusterAutoscalerRoleName(clusterName)

	accountID, err := i.GetAWSAccountID()
	if err != nil {
		return nil, err
	}

	issuerURL, _ := url.Parse(oidcProvider)
	hostPath := issuerURL.Hostname() + issuerURL.Path

	return i.EnsureIAMRoleWithEmbeddedPolicy(
		name,
		fmt.Sprintf("IAM role for %q Kore cluster used by the autoscaler", clusterName),
		fmt.Sprintf(clusterAutoscalerTrustPolicy, accountID, hostPath, hostPath),
		fmt.Sprintf(clusterAutoscalerNodeGroupAGSAccessPolicy, clusterName),
	)
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
