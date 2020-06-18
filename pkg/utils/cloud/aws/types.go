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
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Credentials defines a SET of AWS credentials
type Credentials struct {
	// SecretAccessKey is the AWS Secret Access Key
	SecretAccessKey string
	// AccessKeyID is the AWS Access Key ID
	AccessKeyID string
	// AccountID is the AWS account these credentials reside within
	AccountID string
}

// VPC is the input for creating an aws VPC client
type VPC struct {
	// CidrBlock is the private network address range for any private subnects
	CidrBlock string
	// Name is the VPC name in aws
	Name string
	// Region is the AWS region of the VPC
	Region string
	// Tags - how to find resources
	Tags map[string]string
	// Cache of aws VPC
	awsObj *ec2.Vpc
}

// ASGDetails is the information extracted from an ARN from an autoscaling group
type ASGDetails struct {
	Name string
	ID   string
	ARN  arn.ARN
}

// NodeGroupAutoScaler is the input for creating IAM roles and policies
type NodeGroupAutoScaler struct {
	NodeGroupName  string
	AutoScalingARN string
}

// VPCResult contains all relevant information about the created VPC
type VPCResult struct {
	VPC                         *ec2.Vpc
	PublicSubnets               []*ec2.Subnet
	PrivateSubnets              []*ec2.Subnet
	NATGateways                 []*ec2.NatGateway
	ControlPlaneSecurityGroupID string
}
