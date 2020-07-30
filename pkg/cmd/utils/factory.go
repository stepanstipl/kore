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

package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/appvia/kore/pkg/apiserver/types"
	"github.com/appvia/kore/pkg/client"
	"github.com/appvia/kore/pkg/client/config"
	"github.com/appvia/kore/pkg/cmd/errors"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/openid"

	log "github.com/sirupsen/logrus"
)

type factory struct {
	client    client.Interface
	streams   Streams
	cfg       *config.Config
	resources Resources
}

// NewFactory returns a default factory
func NewFactory(client client.Interface, streams Streams, config *config.Config) (Factory, error) {
	resources, err := newResourceManager(client, config)
	if err != nil {
		return nil, err
	}

	return &factory{
		cfg:       config,
		client:    client,
		resources: resources,
		streams:   streams,
	}, nil
}

func (f *factory) refreshToken() {
	auth := f.cfg.GetAuthInfo(f.client.CurrentProfile())
	if auth.OIDC != nil {
		// @step: has the access token expired
		token, err := utils.NewClaimsFromRawToken(auth.OIDC.IDToken)
		if err == nil {
			if token.HasExpired() {
				log.Debug("attempting to refresh id-token")

				refresh, err := openid.DefaultTokenRefresher.RefreshToken(context.Background(),
					auth.OIDC.RefreshToken,
					auth.OIDC.TokenURL,
					auth.OIDC.ClientID,
					auth.OIDC.ClientSecret,
				)
				if err == nil {
					auth.OIDC.AccessToken = refresh.AccessToken
					auth.OIDC.IDToken = refresh.IDToken
					_ = f.UpdateConfig()
					log.Debug("id-token refreshed successfully")
				} else {
					log.WithError(err).Debug("error refreshing id-token")
					log.Warn("Failed to refresh your access token, please run kore login")
				}
			}
		}
	}
}

// Client returns the underlying client
func (f *factory) Client() client.Interface {
	return f.client
}

// ClientWithEndpoint returns the api client with a specific endpoint
func (f *factory) ClientWithEndpoint(endpoint string) client.RestInterface {
	f.refreshToken()
	return f.client.Request().Endpoint(endpoint)
}

// ClientWithResource returns the api client with a specific resource
func (f *factory) ClientWithResource(resource Resource) client.RestInterface {
	f.refreshToken()
	return f.client.Request().Resource(resource.GetAPIName())
}

// ClientWithTeamResource returns the api client with a specific team resource
func (f *factory) ClientWithTeamResource(team string, resource Resource) client.RestInterface {
	f.refreshToken()
	return f.client.Request().Team(team).Resource(resource.GetAPIName())
}

// CheckError handles the cli errors for us
func (f *factory) CheckError(kerror error) {
	err := func() error {
		switch {
		case client.IsNotAuthorized(kerror):
			return errors.ErrAuthentication
		case client.IsMethodNotAllowed(kerror):
			return errors.ErrOperationNotPermitted
		case client.IsNotImplemented(kerror):
			return errors.ErrOperationNotSupported
		case errors.IsError(kerror, &errors.ErrProfileInvalid{}):
			return fmt.Errorf("invalid: %s", kerror.Error())
		}

		return kerror
	}()
	if err != nil {
		fmt.Fprintf(f.Stderr(), "Error: %s\n", err)

		os.Exit(1)
	}
}

// Whoami returns the details of who they logged in with
func (f *factory) Whoami() (*types.WhoAmI, error) {
	who := &types.WhoAmI{}

	return who, f.ClientWithEndpoint("/whoami").Result(who).Get().Error()
}

// UpdateConfig is responsible for updating the configuration
func (f *factory) UpdateConfig() error {
	return config.UpdateConfig(f.cfg, config.GetClientConfigurationPath())
}

// Config returns the factory client configuration
func (f *factory) Config() *config.Config {
	return f.cfg
}

// SetStdin allows you to set the stdin for the factory
func (f *factory) SetStdin(in io.Reader) {
	f.streams.Stdin = in
}

// Stdin return the standard input
func (f *factory) Stdin() io.Reader {
	return f.streams.Stdin
}

// Stderr returns the io.Writer for errors
func (f *factory) Stderr() io.Writer {
	return f.streams.Stderr
}

// Writer returns the io.Writer for output
func (f *factory) Writer() io.Writer {
	return f.streams.Stdout
}

// Printf writes a message to the io.Writer
func (f *factory) Printf(message string, args ...interface{}) {
	fmt.Fprintf(f.Writer(), message, args...)
}

// Println writes a message to the io.Writer
func (f *factory) Println(message string, args ...interface{}) {
	filtered := strings.TrimRight(message, "\n")

	fmt.Fprintf(f.Writer(), filtered+"\n", args...)
}

// Resources returns the resources contract
func (f *factory) Resources() Resources {
	return f.resources
}
