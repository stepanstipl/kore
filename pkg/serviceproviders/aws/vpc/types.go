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

package vpc

type Configuration struct {
	// Name is the VPC name. If empty it defaults to the service's name
	Name string `json:"name"`
	// PrivateIPV4Cidr is the private range used for the VPC
	PrivateIPV4Cidr string `json:"privateIPV4Cidr"`
	// Region is the AWS region of the VPC and any resources created
	Region string `json:"region"`
	// SubnetCount is the maximum number of subnets of each subnet type
	SubnetCount int `json:subnetCount`
}

type ProviderData struct {
	// PrivateSubnetIds is a list of subnet IDs to use for the worker nodes
	PrivateSubnetIDs []string `json:"privateSubnetIDs,omitempty"`
	// PublicSubnetIDs is a list of subnet IDs to use for resources that need a public IP (e.g. load balancers)
	PublicSubnetIDs []string `json:"publicSubnetIDs,omitempty"`
	// SecurityGroupIds is a list of security group IDs to use for a cluster
	SecurityGroupIDs []string `json:"securityGroupIDs,omitempty"`
	// PublicIPV4EgressAddresses provides the source addresses for traffic coming from the cluster
	// - can provide input for securing Kube API endpoints in managed clusters
	PublicIPV4EgressAddresses []string `json:"ipv4EgressAddresses,omitempty"`
}
