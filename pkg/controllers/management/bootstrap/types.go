/*
 * Copyright (C) 2019 Rohith Jayawardene <gambol99@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package bootstrap

// BrokerOptions are setting for the service broker
type BrokerOptions struct {
	// Name is the name of the broker
	Name string `json:"name,omitempty"`
	// Password is the password
	Password string `json:"password,omitempty"`
	// Username is a username to use
	Username string `json:"username,omitempty"`
	// Database are options for the database
	Database DatabaseOptions `json:"database,omitempty"`
}

// DatabaseOptions are the database options
type DatabaseOptions struct {
	// Name is the database name
	Name string `json:"name,omitempty"`
	// Password is the database password
	Password string `json:"password,omitempty"`
}

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

// KialiOptions are options for the kiali service
type KialiOptions struct {
	// Password is the password for the service
	Password string `json:"password,omitempty"`
	// Username is the username
	Username string `json:"username,omitempty"`
}

// OperatorOptions are the options for an operator
type OperatorOptions struct {
	// Catalog is an optional field for when we divert from the catalog
	Catalog string `json:"catalog,omitempty"`
	// Package is the operator package to install
	Package string `json:"package,omitempty"`
	// InstallPlan is the desired installplan
	InstallPlan string `json:"install_plan,omitempty"`
	// Channel is the channel to use
	Channel string `json:"channel,omitempty"`
	// Label is the selector label
	Label string `json:"label,omitempty"`
	// Namespace is the namespace to use
	Namespace string `json:"namespace,omitempty"`
}

// CatalogOptions are the options for a catalog
type CatalogOptions struct {
	// Image is the image we should use
	Image string `json:"image,omitempty"`
	// GRPC is an hostname entry
	GRPC string `json:"grpc,omitempty"`
}

// NamespaceOptions is a name to create
type NamespaceOptions struct {
	// EnableIstio indicates istio should be enabled
	EnableIstio bool `json:"enable_istio,omitempty"`
	// Name is the name of the namespace
	Name string `json:"name,omitempty"`
}

// Parameters provides the context for the job parameters
type Parameters struct {
	// BootImage is the image we are using to bootstrap
	BootImage string `json:"boot_image,omitempty"`
	// Broker are setting for the service broker
	Broker BrokerOptions `json:"broker,omitempty"`
	// Catalog are parameters for OLM catalog
	Catalog CatalogOptions `json:"catalog,omitempty"`
	// Credentials are creds for the providers
	Credentials Credentials `json:"credentials,omitempty"`
	// Domain is the cluster domain
	Domain string `json:"domain,omitempty"`
	// EnableIstio indicates if istio is enabled
	EnableIstio bool `json:"enable_istio,omitempty"`
	// EnableKiali indicates kiali should be enabled
	EnableKiali bool `json:"enable_kiali,omitempty"`
	// EnableServiceBroker indicates the broker is enabled
	EnableServiceBroker bool `json:"enable_service_broker,omitempty"`
	// Kiali are options for the kiali service
	Kiali KialiOptions `json:"kiali,omitempty"`
	// Namespaces is a collection of namespaces to create
	Namespaces []NamespaceOptions `json:"namespaces,omitempty"`
	// Operators is a collection of operators to install
	Operators []OperatorOptions `json:"operators,omitempty"`
	// OperatorGroups is a collection of operatorgroups to create
	OperatorGroups []string `json:"operator_groups,omitempty"`
	// OLMVersion is the version of the OLM to install
	OLMVersion string `json:"olm_version,omitempty"`
	// Provider is the cloud provider
	Provider string `json:"provider,omitempty"`
	// StorageClass is the class to use when creating PVC's
	StorageClass string `json:"storage_class,omitempty"`
}
