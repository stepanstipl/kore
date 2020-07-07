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
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

// IsAWSErr returns true if the error is an AWS error
func IsAWSErr(err error, code string, message string) bool {
	var awsErr awserr.Error
	if errors.As(err, &awsErr) {
		return awsErr.Code() == code && strings.Contains(awsErr.Message(), message)
	}
	return false
}

// IsAWSErrRequestFailureStatusCode check the whether the error is an AWS error and has the given HTTP status code
func IsAWSErrRequestFailureStatusCode(err error, statusCode int) bool {
	var awsErr awserr.RequestFailure
	if errors.As(err, &awsErr) {
		return awsErr.StatusCode() == statusCode
	}
	return false
}
