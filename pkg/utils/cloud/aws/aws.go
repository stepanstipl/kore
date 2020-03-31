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

	version "github.com/appvia/kore/pkg/version"
	"github.com/aws/aws-sdk-go/aws"
	awscreds "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
)

func getNewSession(creds Credentials, region string) *session.Session {

	logger := log.New()
	config := aws.NewConfig().WithCredentials(
		awscreds.NewStaticCredentials(creds.AccessKeyID, creds.SecretAccessKey, ""),
	)
	config.Region = &region

	config = request.WithRetryer(config, newLoggingRetryer())
	if logger.IsLevelEnabled(log.DebugLevel) {
		config = config.WithLogLevel(aws.LogDebug |
			aws.LogDebugWithHTTPBody |
			aws.LogDebugWithRequestRetries |
			aws.LogDebugWithRequestErrors |
			aws.LogDebugWithEventStreamBody)
		config = config.WithLogger(aws.LoggerFunc(func(args ...interface{}) {
			logger.Debug(fmt.Sprintln(args...))
		}))
	}

	// Create the options for the session
	opts := session.Options{
		Config:            *config,
		SharedConfigState: session.SharedConfigEnable,
	}

	s := session.Must(session.NewSessionWithOptions(opts))
	s.Handlers.Build.PushFrontNamed(request.NamedHandler{
		Name: "appviaKoreAwsAgent",
		Fn: request.MakeAddToUserAgentHandler(
			"kore", version.Version()),
	})

	return s
}
