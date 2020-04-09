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
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// IamClient describes a aws session and Iam service
type IamClient struct {
	// sess is the AWS session
	sess *session.Session
	// svc is the iam service
	svc *iam.IAM
	// namePrefix is the common name of the objects managed
	namePrefix string
	myARN      *string
}

const (
	// Policies required for eks clusters:
	amazonEKSClusterPolicy = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
	amazonEKSServicePolicy = "arn:aws:iam::aws:policy/AmazonEKSServicePolicy"

	clusterStsTrustPolicy = `{
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

	// amazonEKSWorkerNodePolicy provides read-only access to Amazon EC2 Container Registry repositories.
	amazonEKSWorkerNodePolicy          = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
	amazonEC2ContainerRegistryReadOnly = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
	amazonEKSCNIPolicy                 = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
)

// NewIamClient will create a new IamClient
func NewIamClient(sess *session.Session, namePrefix string) *IamClient {
	return &IamClient{
		sess:       sess,
		svc:        iam.New(sess),
		namePrefix: namePrefix,
	}
}

func (i *IamClient) getMyARN() (*string, error) {
	if i.myARN == nil {
		// Get currewnt user
		guo, err := i.svc.GetUser(&iam.GetUserInput{})
		if err != nil {
			return nil, err
		}
		i.myARN = guo.User.Arn
	}
	return i.myARN, nil
}

// EnsureEksClusterRole will return the cluster role and the nodepool role
func (i *IamClient) EnsureEksClusterRole() (*iam.Role, error) {
	arn, err := i.getMyARN()
	if err != nil {
		return nil, fmt.Errorf("cannot create eks roles as error obtaining my arn - %s", err)
	}
	clusterSts := fmt.Sprintf(clusterStsTrustPolicy, *arn)
	cr, err := i.ensureRole("eks-cluster", []string{
		amazonEKSClusterPolicy,
		amazonEKSServicePolicy,
	}, clusterSts)
	if err != nil {
		return nil, err
	}
	return cr, nil
}

// EnsureEksNodePoolRole will create a nodepool eks role
func (i *IamClient) EnsureEksNodePoolRole() (*iam.Role, error) {
	npr, err := i.ensureRole("eks-nodepool", []string{
		amazonEKSWorkerNodePolicy,
		amazonEC2ContainerRegistryReadOnly,
		amazonEKSCNIPolicy,
	}, nodeStsTrustPolicy)
	if err != nil {
		return nil, err
	}
	return npr, nil
}

func (i *IamClient) ensureRole(name string, policyARNs []string, stsPolicy string) (*iam.Role, error) {
	var r *iam.Role
	rn := fmt.Sprintf("%s-%s", i.namePrefix, name)
	r, err := func(rn string) (*iam.Role, error) {
		gr, err := i.svc.GetRole(&iam.GetRoleInput{
			RoleName: &rn,
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case iam.ErrCodeNoSuchEntityException:
					// OK not found
					return nil, nil
				}
			}
			return nil, err
		}
		return gr.Role, nil
	}(rn)
	if err != nil {
		return nil, fmt.Errorf("error looking up aws iam role %s - %s", rn, err)
	}
	if r == nil {
		// create the role
		cro, err := i.svc.CreateRole(&iam.CreateRoleInput{
			AssumeRolePolicyDocument: aws.String(stsPolicy),
			Path:                     aws.String("/"),
			RoleName:                 aws.String(rn),
		})
		if err != nil {
			return nil, fmt.Errorf("error creating role %s - %s", rn, err)
		}
		r = cro.Role
	}
	// Ensure the policies are correct for the role
	lpo, err := i.svc.ListAttachedRolePolicies(&iam.ListAttachedRolePoliciesInput{
		RoleName: r.RoleName,
	})
	if err != nil {
		return nil, fmt.Errorf("error listing attached policies for role [%s] - %s", *r.RoleName, err)
	}
	for _, pa := range policyARNs {
		found := false
		for _, fp := range lpo.AttachedPolicies {
			if *fp.PolicyArn == pa {
				found = true
				break
			}
		}
		if !found {
			_, err := i.svc.AttachRolePolicy(&iam.AttachRolePolicyInput{
				PolicyArn: &pa,
				RoleName:  r.RoleName,
			})
			if err != nil {
				return nil, fmt.Errorf("error attaching policy %s to role %s - %s", pa, *r.RoleName, err)
			}
		}
	}
	return r, nil
}
