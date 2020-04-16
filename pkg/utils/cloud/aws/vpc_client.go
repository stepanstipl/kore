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
	"net"
	"sort"
	"strconv"
	"strings"

	"github.com/appvia/kore/pkg/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	AZLimit = 3
)

// VPCClient for aws VPC
type VPCClient struct {
	// credentials are the aws credentials
	credentials Credentials
	// Sess is the AWS session
	Sess *session.Session
	// The svc to access EC2 resources
	svc *ec2.EC2
	// The VPC to act on
	VPC VPC
}

// NewVPCClient gets an AWS and session with a reference to a matching VPC
func NewVPCClient(creds Credentials, vpc VPC) (*VPCClient, error) {
	sess := getNewSession(creds, vpc.Region)

	// TODO: verify the current CIDR is big enough for the required subnets given:
	// - the constants PrivateNetworkMaskSize, PublicNetworkMaskSize
	// - the expected number of AZ's to use...
	// for now it must be a /15 (assuming three az's)
	_, _, err := net.ParseCIDR(vpc.CidrBlock)
	if err != nil {
		return nil, fmt.Errorf("invalid netmask provided %s", vpc.CidrBlock)
	}
	bitStr := strings.Split(vpc.CidrBlock, "/")[1]
	bits, _ := strconv.ParseInt(bitStr, 10, 8)
	if bits < 16 {
		return nil, fmt.Errorf("vpc cidr too small to create 3x /%d and 3x /%d subnets", PrivateNetworkMaskSize, PublicNetworkMaskSize)
	}
	return &VPCClient{
		credentials: creds,
		Sess:        sess,
		VPC:         vpc,
		svc:         ec2.New(sess),
	}, nil
}

// Ensure will create or update a VPC with ALL required global resources
func (c *VPCClient) Ensure() (ready bool, _ error) {
	// Check if the VPC exists
	found, err := c.Exists()
	if err != nil {
		return false, err
	}
	// Now check it's resources global resources exist
	if !found {
		// time to create
		o, err := c.svc.CreateVpc(&ec2.CreateVpcInput{CidrBlock: aws.String(c.VPC.CidrBlock)})
		if err != nil {
			return false, fmt.Errorf("error creating a new aws vpc %s - %s", c.VPC.Name, err)
		}
		err = createTags(
			*c.svc,
			c.VPC.Name,
			*o.Vpc.VpcId,
			c.VPC.Tags,
		)
		if err != nil {
			return false, fmt.Errorf("error tagging new aws vpc %s, id %s - %s", c.VPC.Name, *o.Vpc.VpcId, err)
		}
		c.VPC.awsObj = o.Vpc
	}

	// Next ensure VPC params set - EnableDnsSupport
	_, err = c.svc.ModifyVpcAttribute(&ec2.ModifyVpcAttributeInput{
		EnableDnsSupport: &ec2.AttributeBooleanValue{
			Value: aws.Bool(true),
		},
		VpcId: c.VPC.awsObj.VpcId,
	})
	if err != nil {
		return false, err
	}
	// Next ensure VPC params set - EnableDnsHostnames
	// Only one at a time, see https://github.com/aws/aws-sdk-go/issues/415
	_, err = c.svc.ModifyVpcAttribute(&ec2.ModifyVpcAttributeInput{
		EnableDnsHostnames: &ec2.AttributeBooleanValue{
			Value: aws.Bool(true),
		},
		VpcId: c.VPC.awsObj.VpcId,
	})
	if err != nil {
		return false, err
	}

	// ensure we have an internet gateway and attach
	igw, err := EnsureInternetGateway(*c.svc, c.VPC)
	if err != nil {
		return false, err
	}

	azs, err := c.getAZs(AZLimit)
	if err != nil {
		return false, err
	}

	// First discover any public subnets or create
	// The public networks will use the very first subnets from the VPC
	vpcStartIP, _, err := net.ParseCIDR(c.VPC.CidrBlock)
	if err != nil {
		return false, err
	}
	publicSubnets, err := EnsurePublicSubnets(
		*c.svc,
		c.VPC,
		azs,
		vpcStartIP,
		*igw.InternetGatewayId,
	)
	if err != nil {
		return false, err
	}

	// Public done, save the subnet ID's
	c.VPC.PublicSubnetIDs = nil
	for _, s := range publicSubnets {
		c.VPC.PublicSubnetIDs = append(c.VPC.PublicSubnetIDs, *s.SubnetId)
	}

	// Get next network address for the internal subnets from the last of the public addresses
	// Get last networ address from public addresses:
	lastPublicSubnet := publicSubnets[len(publicSubnets)-1]
	_, lastPublicNet, err := net.ParseCIDR(*lastPublicSubnet.CidrBlock)
	if err != nil {
		return false, fmt.Errorf("bad ciddr on last aws public subnet %s - %s", *lastPublicSubnet.CidrBlock, err)
	}
	// Get next network of the private size from the last public network
	privateNet, err := utils.GetSubnetFromLast(lastPublicNet, PrivateNetworkMaskSize)
	if err != nil {
		return false, fmt.Errorf("error trying to work next subnet of size %d from %s - %s", PrivateNetworkMaskSize, *lastPublicSubnet.CidrBlock, err)
	}

	privateSubnets, natGateways, ready, err := EnsurePrivateSubnets(*c.svc, c.VPC, azs, privateNet.IP, publicSubnets)
	if err != nil || !ready {
		return ready, err
	}
	c.VPC.PrivateSubnetIDs = nil
	for _, s := range privateSubnets {
		c.VPC.PrivateSubnetIDs = append(c.VPC.PrivateSubnetIDs, *s.SubnetId)
	}

	c.VPC.PublicIPV4EgressAddresses = nil
	for _, n := range natGateways {
		for _, gwa := range n.NatGatewayAddresses {
			c.VPC.PublicIPV4EgressAddresses = append(c.VPC.PublicIPV4EgressAddresses, aws.StringValue(gwa.PublicIp))
		}
	}

	// create security group for master control plane...
	securityGroup, err := EnsureSecurityGroup(*c.svc, c.VPC, SecurityGroupTypeEKSCluster, "eks required group for allowing communication with master nodes")
	if err != nil {
		return false, fmt.Errorf("error finding or creating security group for eks master comms - %s", err)
	}

	c.VPC.ControlPlaneSecurityGroupID = *securityGroup.GroupId

	return true, nil
}

// Exists checks if a vpc exists
func (c *VPCClient) Exists() (bool, error) {
	if c.VPC.awsObj != nil {

		return true, nil
	}
	o, err := c.svc.DescribeVpcs(&ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{getEc2TagNameFilter(c.VPC.Name)},
	})
	if err != nil {

		return false, err
	}
	if len(o.Vpcs) == 1 {
		// Cache the VPC
		c.VPC.awsObj = o.Vpcs[0]

		return true, nil
	}
	if len(o.Vpcs) > 1 {

		return false, fmt.Errorf("Multiple matching VPCs")
	}

	return false, nil
}

// Delete will clear up all VPC resources
// Currently noop
func (c *VPCClient) Delete() (ready bool, _ error) {
	exists, err := c.Exists()
	if err != nil {
		return false, err
	}

	if !exists {
		return true, nil
	}

	if !IsKoreManaged(c.VPC.awsObj.Tags) {
		return true, nil
	}

	azs, err := c.getAZs(AZLimit)
	if err != nil {
		return false, err
	}

	ready, err = DeletePrivateSubnets(*c.svc, c.VPC, azs)
	if err != nil || !ready {
		return ready, err
	}

	if err := DeletePublicSubnets(*c.svc, c.VPC, azs); err != nil {
		return false, err
	}

	if err := DeleteInternetGateway(*c.svc, c.VPC); err != nil {
		return false, err
	}

	if err := DeleteSecurityGroup(*c.svc, c.VPC, SecurityGroupTypeEKSCluster); err != nil {
		return false, err
	}

	_, err = c.svc.DeleteVpc(&ec2.DeleteVpcInput{
		VpcId: c.VPC.awsObj.VpcId,
	})
	if err != nil {
		return false, fmt.Errorf("failed to delete VPC %s: %w", c.VPC.Name, err)
	}

	return true, nil
}

func (c *VPCClient) getAZs(limit int) ([]string, error) {
	res, err := c.svc.DescribeAvailabilityZones(&ec2.DescribeAvailabilityZonesInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to get availability zones: %w", err)
	}

	var azs []string
	for _, az := range res.AvailabilityZones {
		azs = append(azs, *az.ZoneId)
	}

	sort.Strings(azs)

	return azs[0:limit], nil

}
