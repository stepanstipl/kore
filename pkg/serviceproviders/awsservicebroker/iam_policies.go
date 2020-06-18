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

package awsservicebroker

import "fmt"

func IAMRoleTrustPolicy(awsAccountID string) string {
	return fmt.Sprintf(`{
		  "Version": "2012-10-17",
		  "Statement": [
			{
			  "Effect": "Allow",
			  "Principal": {
				"AWS": "arn:aws:iam::%s:root"
			  },
			  "Action": "sts:AssumeRole",
			  "Condition": {}
			}
		  ]
		}`,
		awsAccountID,
	)
}

func IAMRolePolicy(awsAccountID, region string) string {
	return fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Action": [
						"ssm:PutParameter",
						"ssm:GetParameters",
						"ssm:GetParameter"
					],
					"Resource": [
						"arn:aws:ssm:%s:%s:parameter/asb-*",
						"arn:aws:ssm:%s:%s:parameter/Asb*"
					]
				},
				{
					"Effect": "Allow",
					"Action": [
						"athena:*",
						"cloudformation:CancelUpdateStack",
						"cloudformation:CreateStack",
						"cloudformation:DeleteStack",
						"cloudformation:DescribeStackEvents",
						"cloudformation:DescribeStacks",
						"cloudformation:UpdateStack",
						"codecommit:*",
						"cognito:*",
						"documentdb:*",
						"dynamodb:*",
						"ec2:*",
						"elasticache:*",
						"elasticsearch:*",
						"elasticmapreduce:*",
						"iam:*",
						"lex:*",
						"kinesis:*",
						"kms:*",
						"lambda:*",
						"mq:*",
						"polly:*",
						"rds:*",
						"redshift:*",
						"rekognition:*",
						"route53:*",
						"s3:*",
						"sns:*",
						"sqs:*",
						"translate:*"
					],
					"Resource": "*"
				}
			]
		}`,
		region, awsAccountID,
		region, awsAccountID,
	)
}
