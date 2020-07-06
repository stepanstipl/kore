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

package kore

// AWS is the aws interface
type AWS interface {
	// ProjectClaims returns the claims interface
	AWSAccountClaims() AWSAccountClaims
	// Organizations return the organizations interface
	AWSOrganizations() AWSOrganizations
}

type awsImpl struct {
	*cloudImpl
	// team is the request team
	team string
}

// AWSAccountClaims is responsible for managing an AWS account
func (h *awsImpl) AWSAccountClaims() AWSAccountClaims {
	return &awsac{Interface: h.cloudImpl.hubImpl, team: h.team}
}

// Organizations return the organizations interface
func (h *awsImpl) AWSOrganizations() AWSOrganizations {
	return &awsocl{Interface: h.cloudImpl.hubImpl, team: h.team}
}
