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
	// AKS returns the AKS interface
	AKS() AKS
	// AWS returns the aws interface
	AWS() AWS
	// AKSCredentials manages credentials used to create AKS clusters
	AKSCredentials() AKSCredentials
	// GCP returns the gcp interface
	GCP() GCP
	// GKE returns the GKE interface
	GKE() GKE
	// GKECredentials provides access to the gkes credentials
	GKECredentials() GKECredentials
	// EKS returns the EKS interface
	EKS() EKS
	// EKSVPC provides acces to the eks's VPC dependencies
	EKSVPC() EKSVPC
	// EKSCredentials provides acces to the eks's credentials
	EKSCredentials() EKSCredentials
	// EKSNodeGroup provides access to an eks nodegroup
	EKSNodeGroup() EKSNodeGroup
}

type cloudImpl struct {
	*hubImpl
	// team is the requesting team
	team string
}

// AKS retuens a AKS interface implementation
func (c *cloudImpl) AKS() AKS {
	return &aksImpl{cloudImpl: c, team: c.team}
}

// AWS returns the aws interface
func (c *cloudImpl) AWS() AWS {
	return &awsImpl{cloudImpl: c, team: c.team}
}

// AKSCredentials returns an AKSCredentials implementation
func (c *cloudImpl) AKSCredentials() AKSCredentials {
	return &aksCredsImpl{cloudImpl: c, team: c.team}
}

// GCP returns the gcp interface
func (c *cloudImpl) GCP() GCP {
	return &gcpImpl{cloudImpl: c, team: c.team}
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

// EKSVPC returns a eksvpc interface
func (c *cloudImpl) EKSVPC() EKSVPC {
	return &eksVPCImpl{cloudImpl: c, team: c.team}
}

// EKSCredentials returns a eks interface
func (c *cloudImpl) EKSCredentials() EKSCredentials {
	return &eksCredsImpl{cloudImpl: c, team: c.team}
}

func (c *cloudImpl) EKSNodeGroup() EKSNodeGroup {
	return &eksNGImpl{cloudImpl: c, team: c.team}
}
