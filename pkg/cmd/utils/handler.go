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
	"errors"
	"fmt"
	"os"

	"github.com/appvia/kore/pkg/client"
	"github.com/appvia/kore/pkg/client/config"
	"github.com/appvia/kore/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// RunHandler is the run hanler
type RunHandler interface {
	Factory
	// Run is called to handle the action
	Run() error
	// Validate is called to verify the options
	Validate() error
}

// DefaultHandler is a default handler for factory commands
type DefaultHandler struct{}

// Validate is used to validate any options
func (d *DefaultHandler) Validate() error {
	return nil
}

// Run is use to call the action
func (d *DefaultHandler) Run() error {
	return nil
}

// DefaultRunFunc performs a default run handler
func DefaultRunFunc(o RunHandler) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		// @step: we complete the global flags
		utils.SetReflectedField("Force", GetFlagBool(cmd, "force"), o)
		utils.SetReflectedField("NoWait", GetFlagBool(cmd, "no-wait"), o)
		utils.SetReflectedField("Output", GetFlagString(cmd, "output"), o)
		utils.SetReflectedField("Headers", GetFlagBool(cmd, "show-headers"), o)

		// @step: does the caller want to know the path of the kore binary?
		if utils.HasReflectField("KorePath", o) {
			exec, err := os.Executable()
			if err != nil {
				o.CheckError(err)
			}
			utils.SetReflectedField("KorePath", exec, o)
		}

		// @step: some handler require they know the diff between flags set and defaults
		if utils.HasReflectField("FlagsChanged", o) {
			var changed []string

			cmd.Flags().Visit(func(f *pflag.Flag) {
				if f.Changed {
					changed = append(changed, f.Name)
				}
			})
			utils.SetReflectedField("FlagsChanged", changed, o)
		}

		err := utils.SetAndValidateReflectedField("Team", GetFlagString(cmd, "team"), o, checkTeam(o, o.Config()))
		o.CheckError(err)

		// @step: we can help with resource and name as well
		if utils.HasReflectField("Resource", o) && utils.HasReflectField("Name", o) {
			if cmd.Flags().Arg(0) != "" {
				resource, err := o.Resources().Lookup(cmd.Flags().Arg(0))
				o.CheckError(err)
				utils.SetReflectedField("Resource", resource.Name, o)
			}
			utils.SetReflectedField("Name", cmd.Flags().Arg(1), o)
		} else if utils.HasReflectField("Name", o) {
			utils.SetReflectedField("Name", cmd.Flags().Arg(0), o)
		}

		o.CheckError(o.Validate())
		o.CheckError(o.Run())
	}
}

// ExecuteHandler is just shorthand for chaining the method calls
func ExecuteHandler(o RunHandler) error {
	if err := o.Validate(); err != nil {
		return err
	}

	return o.Run()
}

func checkTeam(factory Factory, cfg *config.Config) func(value interface{}) error {
	return func(value interface{}) error {
		team := value.(string)
		if team == "" {
			return nil
		}

		exists, err := factory.ClientWithResource(factory.Resources().MustLookup("team")).Name(team).Exists()
		if err != nil {
			// We will ignore authentication errors as those will be handled elsewhere
			if client.IsNotAuthorized(err) {
				return nil
			}
			return err
		}
		if !exists {
			// Let's check whether the current profile has a non-existing team and remove it
			if cfg.GetProfile(factory.Client().CurrentProfile()).Team == team {
				errStr := fmt.Sprintf("team %q does not exist, please update your profile using\n$ kore profiles set current.team <EXISTING TEAM>", team)
				cfg.GetProfile(factory.Client().CurrentProfile()).Team = ""
				if err := config.UpdateConfig(cfg, config.GetClientConfigurationPath()); err != nil {
					errStr = errStr + "\n" + err.Error()
				}
				return errors.New(errStr)
			}

			return fmt.Errorf("team %q does not exist", team)
		}

		return nil
	}
}
