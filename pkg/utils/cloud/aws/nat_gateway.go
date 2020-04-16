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
	"time"

	"github.com/appvia/kore/pkg/utils"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/ec2"
)

// EIPTargetNATGateway is the EIP target name for an NAT gateway
const EIPTargetNATGateway = "nat-gateway"

func getNATGatewayName(vpc VPC, az string) string {
	return fmt.Sprintf("%s-%s", vpc.Name, az)
}

// EnsureNATGateway will ensure a NAT Gateway exist for a given subnet
// It also provisions an Elastic Ip for the gateway
func EnsureNATGateway(svc ec2.EC2, vpc VPC, subnet *ec2.Subnet) (_ *ec2.NatGateway, ready bool, _ error) {
	name := getNATGatewayName(vpc, *subnet.AvailabilityZoneId)

	// Discover NatGateway if it exists
	natGateway, err := getNATGateway(svc, name, ec2.NatGatewayStateAvailable, ec2.NatGatewayStatePending)
	if err != nil {
		return nil, false, err
	}

	if natGateway == nil {
		// First check EIP exists
		eip, err := EnsureEIP(svc, vpc, EIPTargetNATGateway, *subnet.AvailabilityZoneId)
		if err != nil {
			return nil, false, err
		}

		// Now create NAT Gateway...
		res, err := svc.CreateNatGateway(&ec2.CreateNatGatewayInput{
			AllocationId: eip.AllocationId,
			SubnetId:     subnet.SubnetId,
			TagSpecifications: []*ec2.TagSpecification{
				{
					ResourceType: aws.String(ec2.ResourceTypeNatgateway),
					Tags:         createEC2TagsWithName(name, vpc.Tags),
				},
			},
		})
		if err != nil {
			return nil, false, fmt.Errorf("can not create the AWS NAT Gateway %s - %w", name, err)
		}
		natGateway = res.NatGateway

		// As the AWS API is eventually consistent, we have to wait for the NAT gateway to exist, so we don't
		// create additional ones accidentally. This won't wait for an available state.
		if err := waitForNATGateway(svc, *natGateway.NatGatewayId); err != nil {
			return nil, false, fmt.Errorf("failed to get the AWS NAT Gateway %s (%s) - %w", name, *res.NatGateway.NatGatewayId, err)
		}
	}

	return natGateway, *natGateway.State == ec2.NatGatewayStateAvailable, nil
}

func getNATGateway(svc ec2.EC2, name string, states ...string) (*ec2.NatGateway, error) {
	filters := []*ec2.Filter{
		getEc2TagNameFilter(name),
	}
	if len(states) > 0 {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("state"),
			Values: aws.StringSlice(states),
		})
	}
	res, err := svc.DescribeNatGateways(&ec2.DescribeNatGatewaysInput{
		Filter: filters,
	})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve AWS NAT Gateway %s: %w", name, err)
	}

	switch len(res.NatGateways) {
	case 0:
		return nil, nil
	case 1:
		return res.NatGateways[0], nil
	default:
		return nil, fmt.Errorf("more than one AWS NAT Gateway was found with name %q", name)
	}
}

func waitForNATGateway(svc ec2.EC2, gatewayID string) error {
	err := utils.RetryWithTimeout(context.Background(), 5*time.Minute, 5*time.Second, func() (finished bool, _ error) {
		res, err := svc.DescribeNatGateways(&ec2.DescribeNatGatewaysInput{
			NatGatewayIds: []*string{aws.String(gatewayID)},
		})
		if err != nil { // Ignore any errors while retrying
			return false, nil
		}
		return len(res.NatGateways) == 1, nil
	})

	if err == utils.ErrCancelled {
		_, err = svc.DescribeNatGateways(&ec2.DescribeNatGatewaysInput{
			NatGatewayIds: []*string{aws.String(gatewayID)},
		})
	}
	return err
}

// DeleteNATGateway deletes the NAT gateway if it exists and releases the Elastic IP
func DeleteNATGateway(svc ec2.EC2, vpc VPC, az string) (ready bool, _ error) {
	name := getNATGatewayName(vpc, az)

	natGateway, err := getNATGateway(svc, name)
	if err != nil {
		return false, fmt.Errorf("error getting the NAT gateway %s: %w", name, err)
	}

	if natGateway != nil && !IsKoreManaged(natGateway.Tags) {
		return true, nil
	}

	if natGateway != nil && *natGateway.State == ec2.NatGatewayStateDeleted {
		// We delete the tags only on a best-effort basis, the NAT gateway will go away in a short time either way
		_ = deleteTags(svc, name, *natGateway.NatGatewayId, vpc.Tags)
	}

	if natGateway == nil || *natGateway.State == ec2.NatGatewayStateDeleted {
		err := DeleteEIP(svc, vpc, EIPTargetNATGateway, az)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	if *natGateway.State == ec2.NatGatewayStateDeleting {
		return false, nil
	}

	_, err = svc.DeleteNatGateway(&ec2.DeleteNatGatewayInput{
		NatGatewayId: natGateway.NatGatewayId,
	})
	if err != nil {
		return false, fmt.Errorf("failed to delete NAT gateway %s (%s): %w", name, *natGateway.NatGatewayId, err)
	}

	return false, nil
}
