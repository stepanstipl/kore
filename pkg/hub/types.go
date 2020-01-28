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

package hub

import (
	"context"
	"time"

	"github.com/appvia/kore/pkg/hub/authentication"
	"github.com/appvia/kore/pkg/services/audit"
	"github.com/appvia/kore/pkg/store"
)

const (
	// HubNamespace is the default namespace for the hub
	HubNamespace = "hub"
	// HubOperatorsNamespace is the namespace where operators live
	HubOperatorsNamespace = "hub-operators"
	// HubDefaultTeam is the default team
	HubDefaultTeam = "hub-default"
	// HubAdminTeam is the default hub admin team
	HubAdminTeam = "hub-admin"
	// HubAdminUser is the default hub admin user
	HubAdminUser = "admin"
)

var (
	// Client is the default client for the hub
	Client Interface
)

// Interface is the contrat between the api and store
type Interface interface {
	// Audit returns the audit interface
	Audit() audit.Interface
	// Config returns the hub configure
	Config() *Config
	// Invitations returns the invitations interface
	Invitations() Invitations
	// GetUserIdenity returns the idenity if any of the a user
	GetUserIdentity(context.Context, string) (authentication.Identity, bool, error)
	// Plans returns the plans interface
	Plans() Plans
	// IDP returns the IDP interface
	IDP() IDP
	// Users returns the users interface
	Users() Users
	// Store returns the hub store
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
	// GRPCClientKey is the client key to use when acceccing the DEX grpc server
	GRPCClientKey string `json:"grpcClientKey"`
}

// Config is the configuration for the hub bridge
type Config struct {
	// AdminPass provides a required first time user password
	AdminPass string `json:"admin-pass"`
	// AdminToken is a static admin token for authentication
	AdminToken string `json:"admin-token,omitempty"`
	// Authenticators is a collection of authentication plugins to enable
	Authenticators []string `json:"authenticators,omitempty"`
	// ClientID is the client id for the openid authenticator
	ClientID string `json:"client-id,omitempty"`
	// ClientSecret is the client secret to use
	ClientSecret string `json:"client-secret,omitempty"`
	// CertificateAuthority is the path to a CA
	CertificateAuthority string `json:"certificate-authority,omitempty"`
	// CertificateAuthorityKey is the path to the private key
	CertificateAuthorityKey string `json:"certificate-authority-key,omitempty"`
	// DEX is the config required to configure dex
	DEX DEX `json:"dex"`
	// DiscoveryURL is the openid discovery url
	DiscoveryURL string `json:"discovery-url,omitempty"`
	// HMAC is the token used to sign things
	HMAC string `json:"hmac"`
	// PublicHubURL is the public url for the hub (the ui not the api)
	PublicHubURL string `json:"public-hub-url,omitempty"`
	// PublicAPIURL is the public url for the api
	PublicAPIURL string `json:"public-api-url,omitempty"`
	// UserClaims is collection of claims to identify the username
	UserClaims []string `json:"user-claims,omitempty"`
}
