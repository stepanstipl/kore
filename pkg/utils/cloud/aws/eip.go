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

func getEIPName(vpc VPC, resourceType, az string) string {
	return fmt.Sprintf("%s-%s-%s", vpc.Name, resourceType, az)
}

// EnsureEIP will make sure an Elastic IP exists
func EnsureEIP(svc ec2.EC2, vpc VPC, target, az string) (*ec2.Address, error) {
	name := getEIPName(vpc, target, az)

	address, err := getEIP(svc, vpc, name)
	if err != nil {
		return nil, err
	}

	if address == nil {
		// Create the EIP...
		aao, err := svc.AllocateAddress(&ec2.AllocateAddressInput{
			Domain: aws.String("vpc"),
		})
		if err != nil {
			return nil, fmt.Errorf("can not create an aws public IP address (EIP) %s - %s", name, err)
		}
		address = &ec2.Address{
			AllocationId: aao.AllocationId,
			PublicIp:     aao.PublicIp,
		}
		err = createTags(
			svc,
			name,
			*aao.AllocationId,
			vpc.Tags)
		if err != nil {
			return nil, fmt.Errorf("error tagging new aws public IP address (EIP) %s, id %s - %s", name, *aao.AllocationId, err)
		}
	}
	return address, nil
}

// DeleteEIP will delete an Elastic IP if it exists
func DeleteEIP(svc ec2.EC2, vpc VPC, target, az string) error {
	name := getEIPName(vpc, target, az)
	eip, err := getEIP(svc, vpc, name)
	if err != nil {
		return fmt.Errorf("failed to get existing AWS Elastic IP %s: %w", name, err)
	}

	if eip == nil {
		return nil
	}

	if !IsKoreManaged(eip.Tags) {
		return nil
	}

	_, err = svc.ReleaseAddress(&ec2.ReleaseAddressInput{
		AllocationId: eip.AllocationId,
	})
	if err != nil {
		return fmt.Errorf("failed to release AWS Elastic IP %s: %w", name, err)
	}

	return nil
}

func getEIP(svc ec2.EC2, vpc VPC, name string) (*ec2.Address, error) {
	address, err := svc.DescribeAddresses(&ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{getEc2TagNameFilter(name)},
	})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve AWS Elastic IP %s - %s", name, err)
	}

	if len(address.Addresses) == 1 {
		return address.Addresses[0], nil
	}

	if len(address.Addresses) > 1 {
		return nil, fmt.Errorf("can not retrieve a single AWS Elastic IP for %s as multiple returned", name)
	}

	return nil, nil
}
