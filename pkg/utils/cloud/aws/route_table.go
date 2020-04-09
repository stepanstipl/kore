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
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	// RouteTableTypeNATGateway is for an NAT gateway
	RouteTableTypeNATGateway = "nat-gateway"
	// RouteTableTypeInternet is for an Internet gateway
	RouteTableTypeInternet = "internet"
)

func getRouteTableName(vpc VPC, rtType, az string) string {
	if az == "" {
		return fmt.Sprintf("%s-%s", vpc.Name, rtType)
	}

	return fmt.Sprintf("%s-%s-%s", vpc.Name, rtType, az)
}

// EnsureRouteTable will make sure a route table exists
func EnsureRouteTable(svc ec2.EC2, vpc VPC, rtType, az, gwID string, tags map[string]string) (*ec2.RouteTable, error) {
	name := getRouteTableName(vpc, rtType, az)

	routeTable, err := getRouteTable(svc, name)
	if err != nil {
		return nil, err
	}

	if routeTable == nil {
		res, err := svc.CreateRouteTable(&ec2.CreateRouteTableInput{
			VpcId: vpc.awsObj.VpcId,
		})
		if err != nil {
			return nil, fmt.Errorf("can not create AWS Route Table %s: %w", name, err)
		}
		routeTable = res.RouteTable

		err = createTags(
			svc,
			name,
			*routeTable.RouteTableId,
			vpc.getTagsCopyWith(tags))
		if err != nil {
			return nil, fmt.Errorf("error tagging AWS Route Table %s (%s): %w", name, *routeTable.RouteTableId, err)
		}
	}

	if err := ensureGatewayRoute(svc, name, routeTable, rtType, gwID); err != nil {
		return nil, err
	}

	return routeTable, nil
}

func ensureGatewayRoute(svc ec2.EC2, name string, routeTable *ec2.RouteTable, rtType, gwID string) error {
	routeExists := false
	for _, r := range routeTable.Routes {
		if *r.DestinationCidrBlock == "0.0.0.0/0" {
			switch rtType {
			case RouteTableTypeInternet:
				routeExists = r.GatewayId != nil && *r.GatewayId == gwID
			case RouteTableTypeNATGateway:
				routeExists = r.NatGatewayId != nil && *r.NatGatewayId == gwID
			default:
				panic(fmt.Errorf("invalid route table type: %q", rtType))
			}
		}
	}

	if !routeExists {
		createRouteInput := &ec2.CreateRouteInput{
			DestinationCidrBlock: aws.String("0.0.0.0/0"),
			RouteTableId:         routeTable.RouteTableId,
		}
		switch rtType {
		case RouteTableTypeInternet:
			createRouteInput.GatewayId = aws.String(gwID)
		case RouteTableTypeNATGateway:
			createRouteInput.NatGatewayId = aws.String(gwID)
		default:
			panic(fmt.Errorf("invalid route table type: %q", rtType))
		}
		_, err := svc.CreateRoute(createRouteInput)
		if err != nil {
			return fmt.Errorf("failed to create 0.0.0.0/0 route in route table %s (%s): %s", name, *routeTable.RouteTableId, err)
		}
	}

	return nil
}

func getRouteTable(svc ec2.EC2, name string) (*ec2.RouteTable, error) {
	res, err := svc.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{getEc2TagNameFilter(name)},
	})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve aws route table %s - %s", name, err)
	}

	switch len(res.RouteTables) {
	case 0:
		return nil, nil
	case 1:
		return res.RouteTables[0], nil
	default:
		return nil, fmt.Errorf("failed to retrieve a single Route Table with nane %q as more than one exist with the same name", name)
	}
}

// DeleteRouteTable deletes the route table if it exists
func DeleteRouteTable(svc ec2.EC2, vpc VPC, rtType, az string) error {
	name := getRouteTableName(vpc, rtType, az)
	routeTable, err := getRouteTable(svc, name)
	if err != nil {
		return err
	}

	if routeTable == nil {
		return nil
	}

	_, err = svc.DeleteRouteTable(&ec2.DeleteRouteTableInput{
		RouteTableId: routeTable.RouteTableId,
	})
	if err != nil {
		return fmt.Errorf("failed to delete AWS Route Table %s (%s): %w", name, *routeTable.RouteTableId, err)
	}

	return nil
}
