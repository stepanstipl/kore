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
	"strings"
	"time"

	"github.com/appvia/kore/pkg/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func getEc2TagFiltersFromNameAndTags(name string, tags map[string]string) []*ec2.Filter {
	return getEc2TagFiltersFromNameTagsAndParams(name, tags, map[string]string{})
}

func getEc2TagFiltersFromNameTagsAndParams(name string, inTags, params map[string]string) []*ec2.Filter {
	// Add name
	// don't pollute the original map
	tags := copyTagsWithName(name, inTags)
	filters := []*ec2.Filter{}

	for key, value := range params {
		filters = append(filters, &ec2.Filter{
			Name: aws.String(key),
			Values: []*string{
				aws.String(value),
			},
		})
	}
	for key, value := range tags {
		filters = append(filters, &ec2.Filter{
			Name: aws.String("tag:" + key),
			Values: []*string{
				aws.String(value),
			},
		})
	}
	return filters
}

func tagFromIDNameAndTags(svc ec2.EC2, name, id string, inTags map[string]string) error {
	ec2tags := []*ec2.Tag{}

	// don't pollute the original map
	tags := copyTagsWithName(name, inTags)

	for key, value := range tags {
		ec2tags = append(ec2tags, &ec2.Tag{
			Key:   aws.String(key),
			Value: aws.String(value),
		})
	}

	input := &ec2.CreateTagsInput{
		Resources: []*string{
			aws.String(id),
		},
		Tags: ec2tags,
	}

	timeout := 5 * time.Minute
	err := utils.RetryWithTimeout(context.Background(), timeout, 2*time.Second, func() (finished bool, _ error) {
		_, err := svc.CreateTags(input)
		if err != nil {
			if awserr, ok := err.(awserr.Error); ok {
				if strings.Contains(awserr.Code(), ".NotFound") {
					return false, nil
				}
			}
			return false, err
		}
		return true, nil
	})

	if err != nil {
		if err == utils.ErrCancelled {
			// try one last time to return a real API error
			_, err = svc.CreateTags(input)
		}
		return err
	}

	return nil
}

func copyTagsWithName(name string, inTags map[string]string) map[string]string {
	tags := map[string]string{
		"Name": name,
	}
	for k, v := range inTags {
		tags[k] = v
	}
	return tags
}
