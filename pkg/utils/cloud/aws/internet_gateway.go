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

	"github.com/aws/aws-sdk-go/service/ec2"
)

func getInternetGatewayName(vpc VPC) string {
	return vpc.Name
}

// EnsureInternetGateway ensures an Internet Gateway resource is created and attached
func EnsureInternetGateway(svc ec2.EC2, vpc VPC) (*ec2.InternetGateway, error) {
	name := getInternetGatewayName(vpc)

	ig, err := getInternetGateway(svc, name)
	if err != nil {
		return nil, err
	}

	if ig == nil {
		res, err := svc.CreateInternetGateway(&ec2.CreateInternetGatewayInput{})
		if err != nil {
			return nil, fmt.Errorf("failed to create AWS Internet Gateway %s: %w", name, err)
		}

		ig = res.InternetGateway
		err = createTags(svc, name, *ig.InternetGatewayId, vpc.Tags)
		if err != nil {
			return nil, fmt.Errorf("error tagging AWS Internet Gateway %s (%s): %w", name, *ig.InternetGatewayId, err)
		}
	}

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
			return nil, fmt.Errorf("error trying to attach Internet Gateway %s (%s) to VPC: %w", name, *ig.InternetGatewayId, err)
		}
	}

	return ig, nil
}

func getInternetGateway(svc ec2.EC2, name string) (*ec2.InternetGateway, error) {
	res, err := svc.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{
		Filters: []*ec2.Filter{getEc2TagNameFilter(name)},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve AWS Internet Gateway %s: %w", name, err)
	}

	switch len(res.InternetGateways) {
	case 0:
		return nil, nil
	case 1:
		return res.InternetGateways[0], nil
	default:
		return nil, fmt.Errorf("can not retrieve a single Internet Gateway for %s as multiple returned", name)
	}
}

// DeleteInternetGateway deletes an internet gateway if it exists
func DeleteInternetGateway(svc ec2.EC2, vpc VPC) error {
	name := getInternetGatewayName(vpc)

	ig, err := getInternetGateway(svc, name)
	if err != nil {
		return err
	}

	if ig == nil {
		return nil
	}

	if !IsKoreManaged(ig.Tags) {
		return nil
	}

	attached := false
	for _, a := range ig.Attachments {
		if *a.VpcId == *vpc.awsObj.VpcId {
			attached = true
			break
		}
	}

	if attached {
		_, err := svc.DetachInternetGateway(&ec2.DetachInternetGatewayInput{
			InternetGatewayId: ig.InternetGatewayId,
			VpcId:             vpc.awsObj.VpcId,
		})
		if err != nil {
			return fmt.Errorf("failed to detach Internet Gateway %s (%s) from VPC: %w", name, *ig.InternetGatewayId, err)
		}
	}

	_, err = svc.DeleteInternetGateway(&ec2.DeleteInternetGatewayInput{
		InternetGatewayId: ig.InternetGatewayId,
	})
	if err != nil {
		return fmt.Errorf("failed to delete AWS Internet Gateway %s (%s): %w", name, *ig.InternetGatewayId, err)
	}

	return nil

}
