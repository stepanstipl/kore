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
	"os"
	"testing"

	"github.com/appvia/kore/pkg/utils/cloud/aws/test"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	awscreds "github.com/aws/aws-sdk-go/aws/credentials"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

// TestAssumeRole will test that we can assume an AWS role
func TestAssumeRoleFromCreds(t *testing.T) {
	test.SkipTestIfEnvNotSet("STS", t)
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)

	region := test.GetEnvAndTestLog("KORETEST_REGION", t)
	roleArn := test.GetEnvAndTestLog("KORETEST_ASSUME_ROLE_ARN", t)
	newRegion := test.GetEnvAndTestLog("KORETEST_ASSUME_ROLE_NEW_REGION", t)
	s := AssumeRoleFromCreds(
		Credentials{
			AccessKeyID:     test.AccessKeyID,
			SecretAccessKey: test.SecretAccessKey,
			AccountID:       os.Getenv("KORETEST_STS_ACCOUNT_ID"),
		},
		roleArn,
		region,
		newRegion,
	)

	_, err := s.Config.Credentials.Get()
	require.NoError(t, err)
}

func TestAssumeRoleFromSession(t *testing.T) {
	test.SkipTestIfEnvNotSet("STS", t)
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)

	region := test.GetEnvAndTestLog("KORETEST_REGION", t)
	roleArn := test.GetEnvAndTestLog("KORETEST_ASSUME_ROLE_ARN", t)
	newRegion := test.GetEnvAndTestLog("KORETEST_ASSUME_ROLE_NEW_REGION", t)

	config := aws.NewConfig().WithCredentials(
		awscreds.NewStaticCredentials(
			test.AccessKeyID,
			test.SecretAccessKey,
			"",
		),
	)
	config.Region = aws.String(region)
	opts := session.Options{
		Config:            *config,
		SharedConfigState: session.SharedConfigEnable,
	}
	s := session.Must(session.NewSessionWithOptions(opts))
	s = AssumeRoleFromSession(s, newRegion, roleArn)

	_, err := s.Config.Credentials.Get()
	require.NoError(t, err)
}
