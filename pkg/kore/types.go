/*
 * Copyright (C) 2019 Appvia Ltd <info@appvia.io>
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

package kore

import (
	"context"
	"time"

	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/services/users"
	"github.com/appvia/kore/pkg/services/users/model"
	"github.com/appvia/kore/pkg/store"
)

const (
	// HubNamespace is the default namespace for the kore
	HubNamespace = "kore"
	// HubDefaultTeam is the default team
	HubDefaultTeam = "kore-default"
	// HubAdminTeam is the default kore admin team
	HubAdminTeam = "kore-admin"
	// HubAdminUser is the default kore admin user
	HubAdminUser = "admin"
)

var (
	// Client is the default client for the kore
	Client Interface
)

// Interface is the contrat between the api and store
type Interface interface {
	// Audit returns the audit interface
	Audit() users.Audit
	// Config returns the kore configure
	Config() *Config
	// Invitations returns the invitations interface
	Invitations() Invitations
	// GetUserIdenity returns the idenity if any of the a user
	GetUserIdentity(context.Context, string) (authentication.Identity, bool, error)
	// GetUserIdenityByProvider returns the idenity if any of the a user
	GetUserIdentityByProvider(context.Context, string, string) (*model.Identity, bool, error)
	// Plans returns the plans interface
	Plans() Plans
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
	// ClientID is the client id for the openid authenticator
	ClientID string `json:"client-id,omitempty"`
	// ClientSecret is the client secret to use
	ClientSecret string `json:"client-secret,omitempty"`
	// ClientScopes are additional scopes to add to the request
	ClientScopes []string `json:"client-scopes,omitempty"`
	// CertificateAuthority is the path to a CA
	CertificateAuthority string `json:"certificate-authority,omitempty"`
	// CertificateAuthorityKey is the path to the private key
	CertificateAuthorityKey string `json:"certificate-authority-key,omitempty"`
	// ClusterAppManImage is the image to use for cluster application management
	ClusterAppManImage string `json:"cluster-app-man-image,omitempty"`
	// DEX is the config required to configure dex
	DEX DEX `json:"dex,omitempty"`
	// DiscoveryURL is the openid discovery url
	DiscoveryURL string `json:"discovery-url,omitempty"`
	// EnabledClusterDeletion indicates we should delete cloud providers
	EnableClusterDeletion bool `json:"enable-cluster-deletion,omitempty"`
	// EnableClusterDeletionBlock indicates we should only delete the cluster if the cloud
	// provider is deleted
	EnableClusterDeletionBlock bool `json:"enable-cluster-deletion-block,omitempty"`
	// EnableClusterProviderCheck indicate the k8s controller should check the status of the
	// cloud provider as well
	EnableClusterProviderCheck bool `json:"enable-cluster-provider-check,omitempty"`
	// HMAC is the token used to sign things
	HMAC string `json:"hmac"`
	// PublicHubURL is the public url for the kore (the ui not the api)
	PublicHubURL string `json:"public-kore-url,omitempty"`
	// PublicAPIURL is the public url for the api
	PublicAPIURL string `json:"public-api-url,omitempty"`
	// UserClaims is collection of claims to identify the username
	UserClaims []string `json:"user-claims,omitempty"`
}
