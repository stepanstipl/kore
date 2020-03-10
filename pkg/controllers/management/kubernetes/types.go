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

package kubernetes

// Credentials are creds for the cloud provider
type Credentials struct {
	// AWS are as credentials
	AWS AWSCredentials `json:"aws,omitempty"`
	// GKe are gke credentials
	GKE GKECredentials `json:"gke,omitempty"`
}

// AWSCredentials are the aws crdentials
type AWSCredentials struct {
	// AccountID is the aws account id
	AccountID string `json:"account_id,omitempty"`
	// AccessKey is the credentials id
	AccessKey string `json:"access_key,omitempty"`
	// Region is the AWS region
	Region string `json:"region,omitempty"`
	// SecretKey is the credential key
	SecretKey string `json:"secret_key,omitempty"`
}

// GKECredentials are the creds for gcp
type GKECredentials struct {
	// Account is the json service account
	Account string `json:"account,omitempty"`
}

// NamespaceOptions is a name to create
type NamespaceOptions struct {
	// Name is the name of the namespace
	Name string `json:"name,omitempty"`
}

// Parameters provides the context for clusterappman
type Parameters struct {
	// Credentials are creds for the providers
	Credentials Credentials `json:"credentials,omitempty"`
	// Domain is the cluster domain
	Domain string `json:"domain,omitempty"`
	// Namespaces is a collection of namespaces to create
	Namespaces []NamespaceOptions `json:"namespaces,omitempty"`
	// Provider is the cloud provider
	Provider string `json:"provider,omitempty"`
	// StorageClass is the class to use when creating PVC's
	StorageClass string `json:"storage_class,omitempty"`
}
