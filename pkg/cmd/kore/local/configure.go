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

package local

import (
	"github.com/appvia/kore/pkg/client/config"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

const (
	// LocalProfileName is the profile to create
	LocalProfileName = "local"
	// LocalEndpoint is the local endpoint to use
	LocalEndpoint = "http://127.0.0.1:10080"
)

// LocalConfigureOptions is used to configure the local environment.
type LocalConfigureOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
}

// NewCmdLocalConfigure returns the local configure command
func NewCmdLocalConfigure(factory cmdutil.Factory) *cobra.Command {
	o := &LocalConfigureOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "configure",
		Aliases: []string{"config"},
		Short:   "Configures a local Kore demonstration installation",
		Example: "kore local configure",

		Run: cmdutil.DefaultRunFunc(o),
	}

	return command
}

// Run implements the command action
func (o *LocalConfigureOptions) Run() error {
	// @step: ensure we have a local profile
	conf := o.Config()
	conf.CreateProfile(LocalProfileName, LocalEndpoint)
	conf.CurrentProfile = LocalProfileName

	// @step: prompt for the identity provider settings
	o.Println("What are your Identity Broker details?")
	authInfo := newAuthInfoConfig(conf)
	if err := authInfo.createPrompts().Collect(); err != nil {
		return err
	}
	authInfo.update(conf)

	// @step: write the config out after updating
	if err := o.UpdateConfig(); err != nil {
		return err
	}

	o.Println("...Kore is now set up to run locally,")
	o.Println("✅ A '%s' profile has been configured in %s", LocalProfileName, config.GetClientConfigurationPath())
	return nil
}

type authInfoConfig struct {
	ClientID, ClientSecret string
	AuthorizeURL           string
}

func newAuthInfoConfig(config *config.Config) *authInfoConfig {
	result := &authInfoConfig{}

	if config.AuthInfos["local"] != nil {
		result.ClientID = config.GetCurrentAuthInfo().OIDC.ClientID
		result.ClientSecret = config.GetCurrentAuthInfo().OIDC.ClientSecret
		result.AuthorizeURL = config.GetCurrentAuthInfo().OIDC.AuthorizeURL
	}

	return result
}

func (a *authInfoConfig) createPrompts() cmdutil.Prompts {
	return cmdutil.Prompts{
		&cmdutil.Prompt{Id: "Client ID", ErrMsg: "%s cannot be blank", Value: &a.ClientID},
		&cmdutil.Prompt{Id: "Client Secret", ErrMsg: "%s cannot be blank", Value: &a.ClientSecret},
		&cmdutil.Prompt{Id: "OpenID endpoint", ErrMsg: "%s cannot be blank", Value: &a.AuthorizeURL},
	}
}

func (a *authInfoConfig) update(c *config.Config) {
	c.AddAuthInfo("local", &config.AuthInfo{
		OIDC: &config.OIDC{
			ClientID:     a.ClientID,
			ClientSecret: a.ClientSecret,
			AuthorizeURL: a.AuthorizeURL,
		},
	})
}
