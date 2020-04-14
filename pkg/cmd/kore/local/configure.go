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
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

const (
	// LocalProfileName is the profile to create
	LocalProfileName = "local"
	// LocalEndpoint is the local endpoint to use
	LocalEndpoint = "http://127.0.0.1:10080"
)

// LocalConfigureOptions is used to provision a team
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
		Short:   "Configures a profile to connect to a local Kore installation",
		Example: "kore local configure",

		Run: cmdutil.DefaultRunFunc(o),
	}

	return command
}

// Run implements the command action
func (o *LocalConfigureOptions) Run() error {
	// @step: ensure we have a local profile
	if err := o.MakeLocalClientConfig(); err != nil {
		return err
	}

	// @step: we should prompt for the identity provider settings
	o.Println("What are your Identity Broker details?")

	return nil
}

// MakeLocalClientConfig is used to inject a local profile
func (o *LocalConfigureOptions) MakeLocalClientConfig() error {
	o.Config().CreateProfile(LocalProfileName, LocalEndpoint)

	return o.UpdateConfig()
}
