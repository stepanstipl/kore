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

// Cloud returns a collection of cloud providers
type Cloud interface {
	// GKE returms the GKE interface
	GKE() GKE
	// GKECredentials provides access to the gkes credentials
	GKECredentials() GKECredentials
	// EKS returns the EKS interface
	EKS() EKS
	// EKSCredentials provides acces to the eks's credentials
	EKSCredentials() EKSCredentials
}

type cloudImpl struct {
	*hubImpl
	// team is the requesting team
	team string
}

// GKE returns a gke interface
func (c *cloudImpl) GKE() GKE {
	return &gkeImpl{cloudImpl: c, team: c.team}
}

// GKECredentials returns a gke interface
func (c *cloudImpl) GKECredentials() GKECredentials {
	return &gkeCredsImpl{cloudImpl: c, team: c.team}
}

// EKS retuens a eks interface
func (c *cloudImpl) EKS() EKS {
	return &eksImpl{cloudImpl: c, team: c.team}
}

// EKSCredentials returns a eks interface
func (c *cloudImpl) EKSCredentials() EKSCredentials {
	return &eksCredsImpl{cloudImpl: c, team: c.team}
}
