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

const (
	SecurityGroupTypeEKSCluster = "eks-cluster"
)

func getSecurityGroupName(vpc VPC, name string) string {
	return fmt.Sprintf("%s-%s", vpc.Name, name)
}

// EnsureSecurityGroup will make sure the given security group exists
func EnsureSecurityGroup(svc ec2.EC2, vpc VPC, name, description string) (*ec2.SecurityGroup, error) {
	name = getSecurityGroupName(vpc, name)

	sg, err := getSecurityGroup(svc, name)
	if err != nil {
		return nil, err
	}

	if sg == nil {
		res, err := svc.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
			VpcId:       vpc.awsObj.VpcId,
			Description: aws.String(description),
			GroupName:   aws.String(name),
		})
		if err != nil {
			return nil, fmt.Errorf("can not create an AWS Security Group %s: %w", name, err)
		}

		err = createTags(
			svc,
			name,
			*res.GroupId,
			vpc.Tags)
		if err != nil {
			return nil, fmt.Errorf("error tagging AWS Security Group %s (%s): %w", name, *sg.GroupId, err)
		}

		// As the AWS API is eventually consistent, we have to wait for the security group to exist, so we don't
		// try to create additional ones accidentally.
		sg, err = waitForSecurityGroup(svc, *res.GroupId)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve AWS Security Groups %s (%s): %w", name, *res.GroupId, err)
		}
	}

	return sg, nil
}

func getSecurityGroup(svc ec2.EC2, name string) (*ec2.SecurityGroup, error) {
	res, err := svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{getEc2TagNameFilter(name)},
	})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve AWS Security Groups %s: %w", name, err)
	}

	if len(res.SecurityGroups) > 0 {
		return res.SecurityGroups[0], nil
	}

	return nil, nil
}

func waitForSecurityGroup(svc ec2.EC2, id string) (*ec2.SecurityGroup, error) {
	var sg *ec2.SecurityGroup
	err := utils.RetryWithTimeout(context.Background(), 5*time.Minute, 5*time.Second, func() (finished bool, _ error) {
		res, err := svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
			GroupIds: []*string{aws.String(id)},
		})
		if err != nil {
			return false, nil
		}
		if len(res.SecurityGroups) == 1 {
			sg = res.SecurityGroups[0]
			return true, nil
		}
		return false, nil
	})
	if err == utils.ErrCancelled {
		_, err = svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
			GroupIds: []*string{aws.String(id)},
		})
	}
	return sg, err
}

// DeleteSecurityGroup will delete the security group if it exists
func DeleteSecurityGroup(svc ec2.EC2, vpc VPC, name string) error {
	name = getSecurityGroupName(vpc, name)

	sg, err := getSecurityGroup(svc, name)
	if err != nil {
		return err
	}

	if sg == nil {
		return nil
	}

	_, err = svc.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
		GroupId: sg.GroupId,
	})
	if err != nil {
		return fmt.Errorf("failed to delete AWS Security Group %s (%s): %w", name, *sg.GroupId, err)
	}

	return nil
}
