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

import (
	"context"
	"regexp"
	"time"

	"github.com/appvia/kore/pkg/costs"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/persistence/model"
	"github.com/appvia/kore/pkg/store"
)

const (
	// HubNamespace is the default namespace for the kore
	HubNamespace = "kore"
	// HubDefaultTeam is the default team
	HubDefaultTeam = "kore-default"
	// HubAdminTeam is the default kore admin team
	HubAdminTeam = "kore-admin"
	// HubSystem is the system namespace
	HubSystem = "kore-system"
	// HubOperators is the namespace for operators
	HubOperators = "kore-operators"
	// HubAdminUser is the default kore admin user
	HubAdminUser = "admin"
)

var (
	// Client is the default client for the kore
	Client Interface
	// ResourceNameFilter is a filter that all resource names MUST comply with
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	ResourceNameFilter = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)
	// ResourceAPIFilter defines a api version filter
	ResourceAPIFilter = regexp.MustCompile(`^[a-z\/].*\/v[a-z0-9].*$`)
)

// Interface is the contract between the api and store
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Interface
type Interface interface {
	// Accounts is the accounting interface
	Accounts() Accounts
	// AlertRules() AlertRules
	AlertRules() AlertRules
	// Audit returns the audit interface
	Audit() Audit
	// Config returns the kore configure
	Config() *Config
	// CertificateAuthority returns the CA
	CertificateAuthority() []byte
	// CertificateAuthorityKey is the private key for the CA
	CertificateAuthorityKey() []byte
	// Invitations returns the invitations interface
	Invitations() Invitations
	// GetUserIdenity returns the idenity if any of the a user
	GetUserIdentity(context.Context, string, ...MetaFunc) (authentication.Identity, bool, error)
	// GetUserIdenityByProvider returns the idenity if any of the a user
	GetUserIdentityByProvider(context.Context, string, string) (*model.Identity, bool, error)
	// Plans returns the plans interface
	Plans() Plans
	// Plans returns the plans interface
	PlanPolicies() PlanPolicies
	// ServicePlans returns the interface for service plans
	ServicePlans() ServicePlans
	// ServiceKinds returns the interface for service plans
	ServiceKinds() ServiceKinds
	// ServiceProviders returns the service provider registry
	ServiceProviders() ServiceProviders
	// IDP returns the IDP interface
	IDP() IDP
	// Users returns the users interface
	Users() Users
	// Store returns the kore store
	Store() store.Store
	// Teams returns the teams interface
	Teams() Teams
	// SignedClientCertificate is used to generate a client certificate
	SignedClientCertificate(string, string) ([]byte, []byte, error)
	// SignedServerCertificate is used to generate a server certificate
	SignedServerCertificate([]string, time.Duration) ([]byte, []byte, error)
	// Security returns the security scanner
	Security() Security
	// Persist returns the access layer for the non-Kubernetes data store
	Persist() persistence.Interface
	// Config returns the config interface
	Configs() Configs
	// Costs returns the costs business logic layer
	Costs() costs.Costs
	// Features returns the kore feature control layer
	Features() KoreFeatures
}

// DEX is the configuration required to setup identity providers
type DEX struct {
	// EnableDex indicate is the dex integration is enabled
	EnabledDex bool `json:"enabled-dex,omitempty"`
	// PublicURL the url to the external root of the DEX instance
	PublicURL string `json:"publicURL"`
	// GRPCServer is the host address of the DEX grpc server
	GRPCServer string `json:"grpcServer"`
	// GRPCPort is the port of the DEX grpc server
	GRPCPort int `json:"grpcPort"`
	// GRPCCaCrt is the CA cert of the DEX grpc server
	GRPCCaCrt string `json:"grpcCaCrt"`
	// GRPCClientCrt is the client cert to use when accessing the DEX grpc server
	GRPCClientCrt string `json:"grpcClientCrt"`
	// GRPCClientKey is the client key to use when accessing the DEX grpc server
	GRPCClientKey string `json:"grpcClientKey"`
}

// Config is the configuration for the kore bridge
type Config struct {
	// AdminPass provides a required first time user password
	AdminPass string `json:"admin-pass"`
	// AdminToken is a static admin token for authentication
	AdminToken string `json:"admin-token,omitempty"`
	// Authenticators is a collection of authentication plugins to enable
	Authenticators []string `json:"authenticators,omitempty"`
	// AuthProxyImage is the image to use for oidc proxy
	AuthProxyImage string `json:"auth-proxy-image,omitempty"`
	// CertificateAuthority is the path to a CA
	CertificateAuthority string `json:"certificate-authority,omitempty"`
	// CertificateAuthorityKey is the path to the private key
	CertificateAuthorityKey string `json:"certificate-authority-key,omitempty"`
	// DEX is the config required to configure dex
	DEX DEX `json:"dex,omitempty"`
	// EnableClusterProviderCheck indicate the k8s controller should check the status of the
	// cloud provider as well
	EnableClusterProviderCheck bool `json:"enable-cluster-provider-check,omitempty"`
	// FeatureGates defines which feature gates should be enabled/disabled
	FeatureGates map[string]bool `json:"feature-gates,omitempty"`
	// HMAC is the token used to sign things
	HMAC string `json:"hmac"`
	// IDPClientID is the client id for the openid authenticator
	IDPClientID string `json:"idp-client-id,omitempty"`
	// IDPClientScopes are additional scopes to add to the request
	IDPClientScopes []string `json:"idp-client-scopes,omitempty"`
	// IDPClientSecret is the client secret to use
	IDPClientSecret string `json:"idp-client-secret,omitempty"`
	// IDPServerURL is the openid server url
	IDPServerURL string `json:"idp-server-url,omitempty"`
	// IDPUserClaims is collection of claims to identify the username
	IDPUserClaims []string `json:"idp-user-claims,omitempty"`
	// PublicHubURL is the public url for the kore (the ui not the api)
	PublicHubURL string `json:"public-kore-url,omitempty"`
	// PublicAPIURL is the public url for the api
	PublicAPIURL string `json:"public-api-url,omitempty"`
	// LocalJWTPublicKey is the public key to use to verify JWTs if using the localjwt auth plugin
	LocalJWTPublicKey string `json:"local-jwt-public-key,omitempty"`
	// Costs is the configuration for the costs engine
	Costs costs.Config `json:"costs,omitempty"`
}
