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

// EnsureSecurityGroupAndGetID will find or create a security group and return it's id
func EnsureSecurityGroupAndGetID(svc ec2.EC2, vpc VPC, name, description string) (string, error) {
	tags := vpc.Tags
	sgo, err := svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: getEc2TagFiltersFromNameAndTags(name, tags),
	})
	if err != nil {

		return "", fmt.Errorf("could not retrieve aws security groups %s - %s", name, err)
	}

	if len(sgo.SecurityGroups) == 1 {

		return *sgo.SecurityGroups[0].GroupId, nil
	}
	if len(sgo.SecurityGroups) > 1 {

		return "", fmt.Errorf("can not retrieve a single security group for %s as multiple returned", name)
	}
	csgo, err := svc.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		VpcId:       vpc.awsObj.VpcId,
		Description: aws.String(description),
		GroupName:   aws.String(name),
	})
	if err != nil {
		return "", fmt.Errorf("can not create an aws security group %s - %s", name, err)
	}
	err = tagFromIDNameAndTags(
		svc,
		name,
		*csgo.GroupId,
		tags)
	if err != nil {
		return "", fmt.Errorf("error tagging new aws subnet %s, id %s - %s", name, *csgo.GroupId, err)
	}

	return *csgo.GroupId, nil
}
