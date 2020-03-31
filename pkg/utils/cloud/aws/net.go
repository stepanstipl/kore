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
	"fmt"
	"net"
	"time"

	"github.com/appvia/kore/pkg/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
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

// EnsureInternetGateway ensures an internet gateway resource is ceated an attached
func EnsureInternetGateway(svc ec2.EC2, vpc VPC) (*ec2.InternetGateway, error) {
	name := vpc.Name
	tags := vpc.Tags
	ig, err := func() (*ec2.InternetGateway, error) {
		o, err := svc.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{
			Filters: getEc2TagFiltersFromNameAndTags(name, tags),
		})
		if err != nil {

			return nil, fmt.Errorf("could not retrieve aws internet gateway %s - %s", name, err)
		}

		if len(o.InternetGateways) == 1 {

			return o.InternetGateways[0], nil
		}
		if len(o.InternetGateways) > 1 {

			return nil, fmt.Errorf("can not retrieve a single internet gateway for %s as multiple returned", name)
		}
		return nil, nil
	}()
	if err != nil {
		return nil, err
	}
	if ig == nil {
		o, err := svc.CreateInternetGateway(&ec2.CreateInternetGatewayInput{})
		if err != nil {
			return nil, fmt.Errorf("can not create an aws internet gateway %s - %s", name, err)
		}
		ig = o.InternetGateway
		err = tagFromIDNameAndTags(svc, name, *ig.InternetGatewayId, tags)
		if err != nil {
			return nil, fmt.Errorf("error tagging new aws internetgateway %s, id %s - %s", name, *ig.InternetGatewayId, err)
		}
	}
	// now check if IG is attached to our VPC
	attached := false
	for _, a := range ig.Attachments {
		if *a.VpcId == *vpc.awsObj.VpcId {
			attached = true
			break
		}
	}
	if !attached {
		_, err := svc.AttachInternetGateway(&ec2.AttachInternetGatewayInput{
			InternetGatewayId: ig.InternetGatewayId,
			VpcId:             vpc.awsObj.VpcId,
		})
		if err != nil {
			return nil, fmt.Errorf("error trying to attach internet gateway %s id %s to vpc %s", name, *ig.InternetGatewayId, *vpc.awsObj.VpcId)
		}
	}

	// return the internet gateway id
	return ig, nil
}

// EnsurePublicRoutes will discover or create a public route table and return it
func EnsurePublicRoutes(svc ec2.EC2, vpc VPC, ig *ec2.InternetGateway) (*ec2.RouteTable, error) {
	name := fmt.Sprintf("%s-Public-Routes-Internet", vpc.Name)
	extraTags := map[string]string{
		"Network": "Public",
	}

	return EnsureDefaultRoute(svc, vpc, ig.InternetGatewayId, false, name, extraTags)
}

// EnsureDefaultRoute will discover or create and configure route table for a gateway
func EnsureDefaultRoute(svc ec2.EC2, vpc VPC, gwID *string, natGw bool, name string, extraTags map[string]string) (*ec2.RouteTable, error) {
	tags := vpc.getTagsCopyWith(extraTags)
	params := map[string]string{
		"vpc-id": *vpc.awsObj.VpcId,
	}

	// Discover or create route table
	rt, err := func() (*ec2.RouteTable, error) {
		o, err := svc.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
			Filters: getEc2TagFiltersFromNameTagsAndParams(name, tags, params),
		})
		if err != nil {

			return nil, fmt.Errorf("could not retrieve aws route table %s - %s", name, err)
		}

		if len(o.RouteTables) == 1 {

			return o.RouteTables[0], nil
		}
		if len(o.RouteTables) > 1 {

			return nil, fmt.Errorf("can not retrieve a single routetable for %s as multiple returned", name)
		}
		return nil, nil
	}()
	if err != nil {
		return nil, err
	}
	if rt == nil {
		o, err := svc.CreateRouteTable(&ec2.CreateRouteTableInput{
			VpcId: vpc.awsObj.VpcId,
		})
		if err != nil {
			return nil, fmt.Errorf("can not create an aws routetable %s - %s", name, err)
		}
		rt = o.RouteTable
		err = tagFromIDNameAndTags(
			svc,
			name,
			*rt.RouteTableId,
			tags)
		if err != nil {
			return nil, fmt.Errorf("error tagging new aws routetable %s, id %s - %s", name, *rt.RouteTableId, err)
		}
	}
	routeExists := false
	for _, r := range rt.Routes {
		if *r.DestinationCidrBlock == "0.0.0.0/0" {
			if natGw {
				if r.NatGatewayId != nil {
					if *r.NatGatewayId == *gwID {
						routeExists = true
						break
					}
				}
			} else {
				if r.GatewayId != nil {
					if *r.GatewayId == *gwID {
						routeExists = true
						break
					}
				}
			}
		}
	}
	if !routeExists {
		ri := &ec2.CreateRouteInput{
			DestinationCidrBlock: aws.String("0.0.0.0/0"),
			RouteTableId:         rt.RouteTableId,
		}
		if natGw {
			ri.NatGatewayId = gwID
		} else {
			ri.GatewayId = gwID
		}
		_, err := svc.CreateRoute(ri)
		if err != nil {
			return nil, fmt.Errorf("error trying to set default route %s id %s to vpc %s - %s", name, *rt.RouteTableId, *vpc.awsObj.VpcId, err)
		}
	}

	// return the Routetable
	return rt, nil
}

// EnsureNATGateways will find or create aws nat gatways
func EnsureNATGateways(svc ec2.EC2, publicSubnets []*ec2.Subnet, vpc VPC) (map[string]*ec2.NatGateway, error) {
	natGwsByAz := map[string]*ec2.NatGateway{}
	for _, ps := range publicSubnets {
		n, err := EnsureNATGateway(svc, ps, vpc)
		if err != nil {
			return nil, err
		}
		natGwsByAz[*ps.AvailabilityZoneId] = n
	}
	return natGwsByAz, nil
}

// EnsureNATGateway will find or create a nat gateway for a vpc
func EnsureNATGateway(svc ec2.EC2, publicSubnet *ec2.Subnet, vpc VPC) (*ec2.NatGateway, error) {
	nameSuffix := "-nat-gateway-" + *publicSubnet.AvailabilityZoneId
	name := vpc.Name + nameSuffix
	tags := vpc.Tags

	// Discover NatGateway if it exists
	n, err := func() (*ec2.NatGateway, error) {
		o, err := svc.DescribeNatGateways(&ec2.DescribeNatGatewaysInput{
			Filter: getEc2TagFiltersFromNameAndTags(name, tags),
		})
		if err != nil {

			return nil, fmt.Errorf("could not retrieve aws nat gateway %s - %s", name, err)
		}
		if len(o.NatGateways) > 0 {
			count := 0
			n := &ec2.NatGateway{}
			for _, gw := range o.NatGateways {
				if *gw.State == ec2.NatGatewayStateAvailable || *gw.State == ec2.NatGatewayStatePending {
					count++
					n = gw
				}
			}
			if count == 1 {
				return n, nil
			}
		}
		return nil, nil
	}()
	if err != nil {
		return nil, err
	}
	if n != nil {
		if *n.State == ec2.NatGatewayStateDeleted || *n.State == ec2.NatGatewayStateDeleting {
			// We need this - create!
			n = nil
		}
	}
	if n == nil {
		// First check EIP exists
		eip, err := EnsureEIP(svc, vpc, vpc.Name+"-eip-for-"+nameSuffix)
		if err != nil {
			return nil, err
		}
		// Now create NAT Gateway...
		// ...and wait till ready!
		a, err := svc.CreateNatGateway(&ec2.CreateNatGatewayInput{
			AllocationId: eip.AllocationId,
			SubnetId:     publicSubnet.SubnetId,
		})
		if err != nil {
			return nil, fmt.Errorf("can not create an aws nat gateway %s - %s", name, err)
		}
		n = a.NatGateway
		err = tagFromIDNameAndTags(
			svc,
			name,
			*n.NatGatewayId,
			tags)
		if err != nil {
			return nil, fmt.Errorf("error tagging new aws nat gateway %s, id %s - %s", name, *n.NatGatewayId, err)
		}
	}
	return n, nil
}

// WaitForNatGateway will ensure gateway is ready before trying to use
func WaitForNatGateway(svc ec2.EC2, n *ec2.NatGateway) (bool, error) {
	// lets give amazon a minute to make this good...
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Minute)
	defer cancel()
	for {
		// @step: we break out or continue
		select {
		case <-ctx.Done():
			return false, errors.New("waiting for nat gateway has been cancelled")
		default:
		}
		// now check the status
		o, err := svc.DescribeNatGateways(&ec2.DescribeNatGatewaysInput{
			NatGatewayIds: []*string{n.NatGatewayId},
		})
		if err != nil {

			return false, fmt.Errorf("could not retrieve aws nat gateway %s - %s", *n.NatGatewayId, err)
		}
		if len(o.NatGateways) != 1 {

			return false, fmt.Errorf("no single nat gateway returned %d when looking up %s", len(o.NatGateways), *n.NatGatewayId)
		}
		if len(o.NatGateways) == 1 {
			if *o.NatGateways[0].State == ec2.NatGatewayStateAvailable {
				return true, nil
			}
		}

		time.Sleep(15 * time.Second)
	}
}

// EnsureEIP will retrieve or create an EIP
func EnsureEIP(svc ec2.EC2, vpc VPC, name string) (*ec2.Address, error) {
	tags := vpc.Tags

	// Discover NatGateway if it exists
	a, err := func() (*ec2.Address, error) {
		o, err := svc.DescribeAddresses(&ec2.DescribeAddressesInput{
			Filters: getEc2TagFiltersFromNameAndTags(name, tags),
		})
		if err != nil {

			return nil, fmt.Errorf("could not retrieve aws public IP address (EIP) %s - %s", name, err)
		}

		if len(o.Addresses) == 1 {

			return o.Addresses[0], nil
		}
		if len(o.Addresses) > 1 {

			return nil, fmt.Errorf("can not retrieve a single aws public IP address (EIP) for %s as multiple returned", name)
		}
		return nil, nil
	}()
	if err != nil {
		return nil, err
	}
	if a == nil {
		// Create the EIP...
		aao, err := svc.AllocateAddress(&ec2.AllocateAddressInput{
			Domain: aws.String("vpc"),
		})
		if err != nil {
			return nil, fmt.Errorf("can not create an aws public IP address (EIP) %s - %s", name, err)
		}
		a = &ec2.Address{
			AllocationId: aao.AllocationId,
			PublicIp:     aao.PublicIp,
		}
		err = tagFromIDNameAndTags(
			svc,
			name,
			*aao.AllocationId,
			tags)
		if err != nil {
			return nil, fmt.Errorf("error tagging new aws public IP address (EIP) %s, id %s - %s", name, *aao.AllocationId, err)
		}
	}
	return a, nil
}

// EnsurePrivateSubnets will return the aws private subnets configured with routetables and routes
func EnsurePrivateSubnets(svc ec2.EC2, azs []*ec2.AvailabilityZone, startIP net.IP, vpc VPC, natGws map[string]*ec2.NatGateway) ([]*ec2.Subnet, error) {
	ipNet, err := utils.GetSubnet(startIP, PrivateNetworkMaskSize)
	if err != nil {
		return nil, fmt.Errorf("cannot find a valid network address starting %s with mask of %d - %s", startIP.String(), PrivateNetworkMaskSize, err)
	}
	subnets := []*ec2.Subnet{}
	for i, az := range azs {
		ngw, ok := natGws[*az.ZoneId]
		if !ok {
			return nil, fmt.Errorf("missing NAT Gateway for az %s", *az.ZoneId)
		}
		sn, err := ensurePrivateSubnetInAZ(svc, az, ipNet, vpc, ngw)
		if err != nil {
			return nil, err
		}
		subnets = append(subnets, sn)
		if i < len(azs) {
			// get next network
			ipNet, err = utils.GetSubnetFromLast(ipNet, PrivateNetworkMaskSize)
			if err != nil {
				return nil, fmt.Errorf(
					"cannot formulate next network from %s for subnet %s using bitmask /%d - %s",
					ipNet.String(),
					*sn.CidrBlock,
					PublicNetworkMaskSize,
					err)
			}
		}
	}
	return subnets, nil
}

// EnsureSubnet will find or create a subnet
func EnsureSubnet(svc ec2.EC2, az *ec2.AvailabilityZone, ipNet *net.IPNet, vpc VPC, moreTags map[string]string, name string) (*ec2.Subnet, error) {
	tags := vpc.getTagsCopyWith(moreTags)
	dso, err := svc.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: getEc2TagFiltersFromNameAndTags(name, tags),
	})
	if err != nil {

		return nil, fmt.Errorf("could not retrieve aws subnets %s - %s", name, err)
	}

	if len(dso.Subnets) == 1 {

		return dso.Subnets[0], nil
	}
	if len(dso.Subnets) > 1 {

		return nil, fmt.Errorf("can not retrieve a single routetable for %s as multiple returned", name)
	}
	co, err := svc.CreateSubnet(&ec2.CreateSubnetInput{
		AvailabilityZoneId: aws.String(*az.ZoneId),
		CidrBlock:          aws.String(ipNet.String()),
		VpcId:              vpc.awsObj.VpcId,
	})
	if err != nil {
		return nil, fmt.Errorf("can not create an aws subnet %s - %s", name, err)
	}
	s := co.Subnet
	err = tagFromIDNameAndTags(
		svc,
		name,
		*s.SubnetId,
		tags)
	if err != nil {
		return nil, fmt.Errorf("error tagging new aws subnet %s, id %s - %s", name, *s.SubnetId, err)
	}
	return s, nil
}

func ensurePrivateSubnetInAZ(svc ec2.EC2, az *ec2.AvailabilityZone, ipNet *net.IPNet, vpc VPC, natGw *ec2.NatGateway) (*ec2.Subnet, error) {
	// Discover or create private network
	name := fmt.Sprintf("%s-PrivateSubNet-%s", vpc.Name, *az.ZoneId)

	extraTags := map[string]string{
		// Private Subnets need to be discoverable from aws kubernetes provider for internal LB's
		"kubernetes.io/role/internal-elb": "1",
		"Network":                         "Private",
	}

	s, err := EnsureSubnet(svc, az, ipNet, vpc, extraTags, name)
	if err != nil {
		return nil, err
	}

	natReady, err := WaitForNatGateway(svc, natGw)
	if err != nil {
		return nil, fmt.Errorf("error checking for status of nat gateway %s - %s", *natGw.NatGatewayId, err)
	}
	if !natReady {
		return nil, fmt.Errorf("timeout awaiting for status of nat gateway %s", *natGw.NatGatewayId)
	}
	// create any route tables with default entries first
	rt, err := EnsureDefaultRoute(
		svc,
		vpc,
		natGw.NatGatewayId,
		true,
		fmt.Sprintf("%s-PrivateRouteTable-%s", vpc.Name, *az.ZoneId),
		map[string]string{
			"Network": "Private",
		},
	)
	if err != nil {
		return nil, err
	}

	// Discover or create any route table associations
	rtaFound := false
	for _, rta := range rt.Associations {
		if *rta.SubnetId == *s.SubnetId {
			rtaFound = true
			break
		}
	}
	if !rtaFound {
		_, err := svc.AssociateRouteTable(&ec2.AssociateRouteTableInput{
			RouteTableId: rt.RouteTableId,
			SubnetId:     s.SubnetId,
		})
		if err != nil {
			return nil, err
		}
	}
	// return the subnet
	return s, nil
}

// EnsurePublicSubnets will:
// - discover or create public subnets
// - configure with routes
func EnsurePublicSubnets(svc ec2.EC2, azs []*ec2.AvailabilityZone, startIP net.IP, vpc VPC, rt *ec2.RouteTable) ([]*ec2.Subnet, error) {
	ipNet, err := utils.GetSubnet(startIP, PublicNetworkMaskSize)
	if err != nil {
		return nil, fmt.Errorf("cannot find a valid network address starting %s with mask of %d - %s", startIP.String(), PublicNetworkMaskSize, err)
	}
	subnets := []*ec2.Subnet{}
	for i, az := range azs {
		log.Debugf("will create subnet using cidr:%s in az:%s", ipNet.String(), *az.GroupName)
		sn, err := ensurePublicSubnetInAZ(svc, az, ipNet, vpc, rt)
		if err != nil {
			return nil, err
		}
		subnets = append(subnets, sn)
		if i < len(azs) {
			// get next network
			ipNet, err = utils.GetSubnetFromLast(ipNet, PublicNetworkMaskSize)
			if err != nil {
				return nil, fmt.Errorf(
					"cannot formulate next network from %s for subnet %s using bitmask /%d - %s",
					ipNet.String(),
					*sn.CidrBlock,
					PrivateNetworkMaskSize,
					err)
			}
		}
	}
	return subnets, nil
}

func ensurePublicSubnetInAZ(svc ec2.EC2, az *ec2.AvailabilityZone, ipNet *net.IPNet, vpc VPC, rt *ec2.RouteTable) (*ec2.Subnet, error) {
	// Discover or create public network
	name := fmt.Sprintf("%s-PublicSubnet-%s", vpc.Name, *az.ZoneId)

	extraTags := map[string]string{
		// Public Subnets need to be discoverable for ELB's from the aws kubernets cloud provider
		"kubernetes.io/role/elb": "1",
		"Network":                "Public",
	}

	s, err := EnsureSubnet(svc, az, ipNet, vpc, extraTags, name)
	if err != nil {
		return nil, err
	}
	rtaFound := false
	for _, rta := range rt.Associations {
		if rta.SubnetId == s.SubnetId {
			rtaFound = true
			break
		}
	}
	if !rtaFound {
		_, err := svc.AssociateRouteTable(&ec2.AssociateRouteTableInput{
			RouteTableId: rt.RouteTableId,
			SubnetId:     s.SubnetId,
		})
		if err != nil {
			return nil, err
		}
	}

	// return the subnet
	return s, nil
}
