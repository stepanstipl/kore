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
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

// ASGClient allows interactioin with AWS ASG apis
type ASGClient struct {
	// session is the AWS session
	session *session.Session
	// svc is the iam service
	svc *autoscaling.AutoScaling
}

// NewASGClient will create a new ASGClient
func NewASGClient(session *session.Session) *ASGClient {
	return &ASGClient{session: session, svc: autoscaling.New(session)}
}

// GetASGFromName will talk to the AWS api to get an autoscaling group ARN
func (a *ASGClient) GetASGFromName(asgName string) (*autoscaling.Group, error) {
	asgo, err := a.svc.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{
			&asgName,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("unable to list autoscaling groups - %s", err)
	}
	count := len(asgo.AutoScalingGroups)
	switch {
	case count < 1:
		return nil, fmt.Errorf("cannot find autoscaling group from name %s", asgName)
	case count > 1:
		return nil, fmt.Errorf("ambiguouse autoscaling group from name %s", asgName)
	default:
	}
	return asgo.AutoScalingGroups[0], nil
}

// GetASGDetailsFromArn get the asg name and id from a ARN
func GetASGDetailsFromArn(asgARN string) (*ASGDetails, error) {
	a := &ASGDetails{}
	pa, err := arn.Parse(asgARN)
	if err != nil {
		return nil, fmt.Errorf("error parsing aws arn from autoscaling group arn %s - %s", asgARN, err)
	}
	a.ARN = pa
	items := strings.Split(pa.Resource, ":")
	if len(items) != 3 {
		return nil, fmt.Errorf("cannot parse asg resource name and id from arn resource %s", a.ARN.Resource)
	}
	a.ID = items[1]
	nameitems := strings.Split(items[2], "/")
	if len(nameitems) != 2 {
		return nil, fmt.Errorf("cannot parse asg resource name from arn resource name field %s", items[2])
	}
	a.Name = nameitems[1]
	return a, nil
}
