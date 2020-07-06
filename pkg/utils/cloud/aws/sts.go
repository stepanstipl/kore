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
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

// AssumeRoleFromCreds wraps AssumeRoleFromSession but gets a session first
func AssumeRoleFromCreds(c Credentials, roleARN, region, newRegion string) *session.Session {
	s := getNewSession(Credentials{
		AccessKeyID:     c.AccessKeyID,
		SecretAccessKey: c.SecretAccessKey,
		AccountID:       c.AccountID,
	}, region)

	return AssumeRoleFromSession(s, newRegion, roleARN)
}

// AssumeRoleFromSession will use STS to obtain a new session with the identity of the role specified
func AssumeRoleFromSession(s *session.Session, region, roleArn string) *session.Session {
	newSession := s.Copy()
	// Update the region before we switch...
	newSession.Config.Region = &region
	newSession.Config.Credentials = stscreds.NewCredentials(newSession, roleArn, func(o *stscreds.AssumeRoleProvider) {
		rsn := "kore-session"
		if rsn != "" {
			o.RoleSessionName = rsn
		}
	})

	return newSession
}
