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
	"strings"

	"github.com/appvia/kore/pkg/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	// SubnetTypePrivate is for private subnets with NAT gateways for Internet access
	SubnetTypePrivate = "private"
	// Public subnet is for public subnets with direct Internet access through an Internet Gateway
	SubnetTypePublic = "public"
)

const (
	// PublicNetworkMaskSize is the size of each Public subnet specified as a bit mask
	// This is only used for creation and will not affect detection of existing subnets
	// TODO this should come from an input scheme to detect drift
	PublicNetworkMaskSize = 19
	// PrivateNetworkMaskSize is the size of each private subnet specified as a bit mask
	// This is only used for creation and will not affect detection of existing subnets
	// TODO this should come from an input scheme to detect drift
	PrivateNetworkMaskSize = 19
)

func getSubnetName(vpc VPC, subnetType, az string) string {
	return fmt.Sprintf("%s-%s-%s", vpc.Name, subnetType, az)
}

// EnsureSubnet will make sure the given subnet exists
func EnsureSubnet(svc ec2.EC2, vpc VPC, subnetType, az, cidrBlock string, tags map[string]string) (*ec2.Subnet, error) {
	name := getSubnetName(vpc, subnetType, az)

	subnet, err := getSubnet(svc, name)
	if err != nil {
		return nil, err
	}

	if subnet == nil {
		res, err := svc.CreateSubnet(&ec2.CreateSubnetInput{
			AvailabilityZoneId: aws.String(az),
			CidrBlock:          aws.String(cidrBlock),
			VpcId:              vpc.awsObj.VpcId,
		})

		if err != nil {
			return nil, fmt.Errorf("can not create AWS Subnet %s: %w", name, err)
		}

		subnet = res.Subnet

		err = createTags(
			svc,
			name,
			*subnet.SubnetId,
			vpc.getTagsCopyWith(tags))
		if err != nil {
			return nil, fmt.Errorf("error tagging AWS subnet %s (%s): %w", name, *subnet.SubnetId, err)
		}
	}

	return subnet, nil
}

func getSubnet(svc ec2.EC2, name string) (*ec2.Subnet, error) {
	res, err := svc.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{getEc2TagNameFilter(name)},
	})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve AWS subnet %s: %w", name, err)
	}

	switch len(res.Subnets) {
	case 0:
		return nil, nil
	case 1:
		return res.Subnets[0], nil
	default:
		return nil, fmt.Errorf("failed to retrieve subnet with name %q as multiple subnets have the same name", name)
	}
}

// DeleteSubnet will delete the subnet if it exists
func DeleteSubnet(svc ec2.EC2, vpc VPC, subnetType, az string) error {
	name := getSubnetName(vpc, subnetType, az)

	subnet, err := getSubnet(svc, name)
	if err != nil {
		return err
	}

	if subnet == nil {
		return nil
	}

	if !IsKoreManaged(subnet.Tags) {
		return nil
	}

	_, err = svc.DeleteSubnet(&ec2.DeleteSubnetInput{
		SubnetId: subnet.SubnetId,
	})
	if err != nil {
		return fmt.Errorf("failed to delete AWS Subnet %s (%s): %w", name, *subnet.SubnetId, err)
	}

	return nil
}

func ensureSubnetRoute(svc ec2.EC2, subnet *ec2.Subnet, rt *ec2.RouteTable) error {
	found := false
	for _, rta := range rt.Associations {
		if *rta.SubnetId == *subnet.SubnetId {
			found = true
			break
		}
	}
	if !found {
		_, err := svc.AssociateRouteTable(&ec2.AssociateRouteTableInput{
			RouteTableId: rt.RouteTableId,
			SubnetId:     subnet.SubnetId,
		})
		if err != nil {
			return fmt.Errorf("failed to associate route table %s with subnet %s: %w", *rt.RouteTableId, *subnet.SubnetId, err)
		}
	}

	return nil
}

// EnsurePrivateSubnet will make sure the given private subnet exists
// It also provisions and NAT Gateway and the required route table.
func EnsurePrivateSubnet(svc ec2.EC2, vpc VPC, az, cidrBlock string, publicSubnet *ec2.Subnet) (_ *ec2.Subnet, _ *ec2.NatGateway, ready bool, _ error) {
	// Private Subnets need to be discoverable from aws kubernetes provider for internal LB's
	subnetTags := map[string]string{
		"kubernetes.io/role/internal-elb": "1",
		"Network":                         strings.Title(SubnetTypePrivate),
	}

	subnet, err := EnsureSubnet(svc, vpc, SubnetTypePrivate, az, cidrBlock, subnetTags)
	if err != nil {
		return nil, nil, false, err
	}

	natGateway, ready, err := EnsureNATGateway(svc, vpc, publicSubnet)
	if err != nil {
		return nil, nil, false, err
	}

	if !ready {
		return nil, nil, false, nil
	}

	routeTableTags := map[string]string{
		"Network": strings.Title(SubnetTypePrivate),
	}

	routeTable, err := EnsureRouteTable(svc, vpc, RouteTableTypeNATGateway, az, *natGateway.NatGatewayId, routeTableTags)
	if err != nil {
		return nil, nil, false, err
	}

	if err := ensureSubnetRoute(svc, subnet, routeTable); err != nil {
		return nil, nil, false, err
	}

	return subnet, natGateway, true, nil
}

// DeletePrivateSubnet will delete the given private subnet if it exists
func DeletePrivateSubnet(svc ec2.EC2, vpc VPC, az string) (ready bool, _ error) {
	ready, err := DeleteNATGateway(svc, vpc, az)
	if err != nil || !ready {
		return ready, err
	}

	if err := DeleteSubnet(svc, vpc, SubnetTypePrivate, az); err != nil {
		return false, err
	}

	if err := DeleteRouteTable(svc, vpc, RouteTableTypeNATGateway, az); err != nil {
		return false, err
	}

	return true, nil
}

// EnsurePublicSubnet will make sure the given public subnet exists
func EnsurePublicSubnet(svc ec2.EC2, vpc VPC, az, cidrBlock string, routeTable *ec2.RouteTable) (*ec2.Subnet, error) {
	// Public Subnets need to be discoverable for ELB's from the aws kubernets cloud provider
	subnetTags := map[string]string{
		"kubernetes.io/role/elb": "1",
		"Network":                strings.Title(SubnetTypePublic),
	}

	subnet, err := EnsureSubnet(svc, vpc, SubnetTypePublic, az, cidrBlock, subnetTags)
	if err != nil {
		return nil, err
	}

	if err := ensureSubnetRoute(svc, subnet, routeTable); err != nil {
		return nil, err
	}

	return subnet, nil
}

// DeletePublicSubnet will delete the public subnet if it exists
func DeletePublicSubnet(svc ec2.EC2, vpc VPC, az string) error {
	return DeleteSubnet(svc, vpc, SubnetTypePublic, az)
}

// EnsurePublicSubnets will create all public subnets
// It also creates the route table to set up the routes to the Internet Gateway
func EnsurePublicSubnets(svc ec2.EC2, vpc VPC, azs []string, startIP net.IP, igwID string) ([]*ec2.Subnet, error) {
	routeTableTags := map[string]string{
		"Network": strings.Title(SubnetTypePublic),
	}

	routeTable, err := EnsureRouteTable(svc, vpc, RouteTableTypeInternet, "", igwID, routeTableTags)
	if err != nil {
		return nil, err
	}

	ipNet, err := utils.GetSubnet(startIP, PublicNetworkMaskSize)
	if err != nil {
		return nil, fmt.Errorf("cannot find a valid network address starting %s with mask of %d: %w", startIP.String(), PublicNetworkMaskSize, err)
	}

	var subnets []*ec2.Subnet
	for _, az := range azs {
		subnet, err := EnsurePublicSubnet(svc, vpc, az, ipNet.String(), routeTable)
		if err != nil {
			return nil, err
		}
		subnets = append(subnets, subnet)

		ipNet, err = utils.GetSubnetFromLast(ipNet, PublicNetworkMaskSize)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot formulate next network from %s for subnet %s using bitmask /%d: %w",
				ipNet.String(),
				*subnet.CidrBlock,
				PrivateNetworkMaskSize,
				err)
		}
	}
	return subnets, nil
}

// DeletePublicSubnets will make sure all public subnets are deleted
func DeletePublicSubnets(svc ec2.EC2, vpc VPC, azs []string) error {
	for _, az := range azs {
		if err := DeletePublicSubnet(svc, vpc, az); err != nil {
			return err
		}
	}

	if err := DeleteRouteTable(svc, vpc, RouteTableTypeInternet, ""); err != nil {
		return err
	}

	return nil
}

// EnsurePrivateSubnets will make sure all private subnets exist
func EnsurePrivateSubnets(svc ec2.EC2, vpc VPC, azs []string, startIP net.IP, publicSubnets []*ec2.Subnet) (_ []*ec2.Subnet, _ []*ec2.NatGateway, ready bool, _ error) {
	ipNet, err := utils.GetSubnet(startIP, PrivateNetworkMaskSize)
	if err != nil {
		return nil, nil, false, fmt.Errorf(
			"cannot find a valid network address starting %s with mask of %d: %w",
			startIP.String(),
			PrivateNetworkMaskSize,
			err,
		)
	}

	var subnets []*ec2.Subnet
	var natGateways []*ec2.NatGateway
	allReady := true
	for _, az := range azs {
		publicSubnet, found := getSubnetByAZ(publicSubnets, az)
		if !found {
			return nil, nil, false, fmt.Errorf("no public subnet found for availability zone %s", az)
		}

		subnet, natGateway, ready, err := EnsurePrivateSubnet(svc, vpc, az, ipNet.String(), publicSubnet)
		if err != nil {
			return nil, nil, false, err
		}
		allReady = allReady && ready
		subnets = append(subnets, subnet)
		natGateways = append(natGateways, natGateway)

		ipNet, err = utils.GetSubnetFromLast(ipNet, PrivateNetworkMaskSize)
		if err != nil {
			return nil, nil, false, fmt.Errorf(
				"cannot formulate next network from %s for subnet %s using bitmask /%d: %w",
				ipNet.String(),
				*subnet.CidrBlock,
				PrivateNetworkMaskSize,
				err)
		}
	}
	return subnets, natGateways, allReady, nil
}

// DeletePrivateSubnets will make sure all private subnets are deleted
func DeletePrivateSubnets(svc ec2.EC2, vpc VPC, azs []string) (ready bool, _ error) {
	allReady := true
	for _, az := range azs {
		ready, err := DeletePrivateSubnet(svc, vpc, az)
		if err != nil {
			return false, err
		}
		allReady = allReady && ready
	}
	if !allReady {
		return false, nil
	}

	return true, nil
}

func getSubnetByAZ(subnets []*ec2.Subnet, az string) (*ec2.Subnet, bool) {
	for _, subnet := range subnets {
		if *subnet.AvailabilityZoneId == az {
			return subnet, true
		}
	}
	return nil, false
}
