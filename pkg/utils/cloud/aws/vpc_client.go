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
	"strconv"
	"strings"

	"github.com/appvia/kore/pkg/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
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
func (c *VPCClient) Ensure() error {
	// Check if the VPC exists
	found, err := c.Exists()
	if err != nil {

		return err
	}
	// Now check it's resources global resources exist
	if !found {
		// time to create
		o, err := c.svc.CreateVpc(&ec2.CreateVpcInput{
			CidrBlock: aws.String(c.VPC.CidrBlock),
		})
		if err != nil {

			return fmt.Errorf("error creating a new aws vpc %s - %s", c.VPC.Name, err)
		}
		err = tagFromIDNameAndTags(
			*c.svc,
			c.VPC.Name,
			*o.Vpc.VpcId,
			c.VPC.Tags,
		)
		if err != nil {

			return fmt.Errorf("error tagging new aws vpc %s, id %s - %s", c.VPC.Name, *o.Vpc.VpcId, err)
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

		return err
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

		return err
	}

	// ensure we have an internet gateway and attach
	ig, err := EnsureInternetGateway(*c.svc, c.VPC)
	if err != nil {

		return err
	}

	publicRt, err := EnsurePublicRoutes(*c.svc, c.VPC, ig)
	if err != nil {

		return err
	}

	// find out the number of az's for this zone
	ao, err := c.svc.DescribeAvailabilityZones(nil)
	if err != nil {
		return err
	}
	azs := ao.AvailabilityZones

	// First discover any public subnets or create
	// The public networks will use the very first subnets from the VPC
	vpcStartIP, _, err := net.ParseCIDR(c.VPC.CidrBlock)
	if err != nil {
		return err
	}
	publicSubnets, err := EnsurePublicSubnets(
		*c.svc,
		azs,
		vpcStartIP,
		c.VPC,
		publicRt,
	)
	if err != nil {
		return err
	}

	// Public done, save the subnet ID's
	c.VPC.PublicSubnetIDs = []string{}
	for _, s := range publicSubnets {
		c.VPC.PublicSubnetIDs = append(c.VPC.PublicSubnetIDs, *s.SubnetId)
	}

	// Create NAT Gateways
	natGws, err := EnsureNATGateways(*c.svc, publicSubnets, c.VPC)
	if err != nil {
		return fmt.Errorf("error trying to ensure nat gateways created for %s - %s", *c.VPC.awsObj.VpcId, err)
	}

	// Get next network address for the internal subnets from the last of the public addresses
	// Get last networ address from public addresses:
	lastPublicSubnet := publicSubnets[len(publicSubnets)-1]
	_, lastPublicNet, err := net.ParseCIDR(*lastPublicSubnet.CidrBlock)
	if err != nil {
		return fmt.Errorf("bad ciddr on last aws public subnet %s - %s", *lastPublicSubnet.CidrBlock, err)
	}
	// Get next network of the private size from the last public network
	priavteNet, err := utils.GetSubnetFromLast(lastPublicNet, PrivateNetworkMaskSize)
	if err != nil {
		return fmt.Errorf("error trying to work next subnet of size %d from %s - %s", PrivateNetworkMaskSize, *lastPublicSubnet.CidrBlock, err)
	}

	privateSubnets, err := EnsurePrivateSubnets(*c.svc, azs, priavteNet.IP, c.VPC, natGws)
	if err != nil {
		return fmt.Errorf("error trying to ensure private subnets - %s", err)
	}
	c.VPC.PrivateSubnetIDs = []string{}
	for _, s := range privateSubnets {
		c.VPC.PrivateSubnetIDs = append(c.VPC.PrivateSubnetIDs, *s.SubnetId)
	}

	// create security group for master control plane...
	c.VPC.ControlPlaneSecurityGroupID, err = EnsureSecurityGroupAndGetID(*c.svc, c.VPC, fmt.Sprintf("%s-", c.VPC.Name), "eks required group for allowing communication with master nodes")
	if err != nil {
		return fmt.Errorf("error finding or creating secruity group for eks master comms - %s", err)
	}
	return nil
}

// Exists checks if a vpc exists
func (c *VPCClient) Exists() (bool, error) {
	if c.VPC.awsObj != nil {

		return true, nil
	}
	o, err := c.svc.DescribeVpcs(&ec2.DescribeVpcsInput{
		Filters: getEc2TagFiltersFromNameTagsAndParams(
			c.VPC.Name,
			c.VPC.Tags,
			c.getVPCParams(),
		),
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
func (c *VPCClient) Delete() error {
	// TODO:
	// delete NAT gateways
	// delete VPC
	// Ensure we stop if anyrthing else detected?
	return nil
}

func (c *VPCClient) getVPCParams() map[string]string {
	return map[string]string{
		"cidr": c.VPC.CidrBlock,
	}
}
