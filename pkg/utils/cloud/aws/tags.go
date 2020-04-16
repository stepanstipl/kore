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

const (
	// TagKoreManaged is used for tagging resources managed by Kore
	TagKoreManaged = "kore.appvia.io/managed"
)

func getEc2TagNameFilter(name string) *ec2.Filter {
	return &ec2.Filter{
		Name: aws.String("tag:Name"),
		Values: []*string{
			aws.String(name),
		},
	}
}

func createTags(svc ec2.EC2, name, id string, tags map[string]string) error {
	input := &ec2.CreateTagsInput{
		Resources: []*string{
			aws.String(id),
		},
		Tags: createEC2TagsWithName(name, tags),
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

func deleteTags(svc ec2.EC2, name, id string, tags map[string]string) error {
	input := &ec2.DeleteTagsInput{
		Resources: []*string{
			aws.String(id),
		},
		Tags: createEC2TagsWithName(name, tags),
	}

	timeout := 5 * time.Minute
	err := utils.RetryWithTimeout(context.Background(), timeout, 2*time.Second, func() (finished bool, _ error) {
		_, err := svc.DeleteTags(input)
		if err != nil {
			if awserr, ok := err.(awserr.Error); ok {
				if strings.Contains(awserr.Code(), ".NotFound") {
					return true, nil
				}
			}
			return false, err
		}
		return true, nil
	})

	if err != nil {
		if err == utils.ErrCancelled {
			// try one last time to return a real API error
			_, err = svc.DeleteTags(input)
		}
		return err
	}

	return nil
}

func createEC2TagsWithName(name string, tags map[string]string) []*ec2.Tag {
	res := []*ec2.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String(name),
		},
	}
	for key, value := range tags {
		res = append(res, &ec2.Tag{
			Key:   aws.String(key),
			Value: aws.String(value),
		})
	}
	return res
}

// IsKoreManaged will return true if a specific tag is present to signal the AWS resource is managed by Kore
func IsKoreManaged(tags []*ec2.Tag) bool {
	for _, tag := range tags {
		if *tag.Key == TagKoreManaged && *tag.Value == "true" {
			return true
		}
	}
	return false
}
