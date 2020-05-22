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

package profiles

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	cmdutil "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
	"github.com/tidwall/sjson"
)

var (
	setLongDescription = `
Sets an individual value in a kore configuration file

Paths are a dot delimited name where each token represents either an attribute
name or a map key i.e. profiles.local.server

Examples:
Set the server value of the local profile
kore profile set profiles.local.server https://1.2.3.4

Set the default team for the profile
kore profile set current.team myteam
`
)

type SetOptions struct {
	cmdutil.Factory
	// Path is the jsonpath for the value
	Path string
	// Value is the value to set
	Value string
}

// NewCmdProfilesSet creates and returns the profile set command
func NewCmdProfilesSet(factory cmdutil.Factory) *cobra.Command {
	o := &SetOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "set",
		Short:   "allows you to set various options within the current selected profile",
		Example: "kore profile use <name>",
		Long:    setLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			o.Path = cmd.Flags().Arg(0)
			o.Value = cmd.Flags().Arg(1)

			o.CheckError(o.Validate())
			o.CheckError(o.Run())
		},

		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if cmd.Flags().NArg() == 0 {
				return []string{"current.team"}, cobra.ShellCompDirectiveNoFileComp
			}

			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return command
}

// Validate checks the options
func (o *SetOptions) Validate() error {
	if o.Path == "" {
		return errors.New("no json path defined")
	}

	return nil
}

// Run implements the actions
func (o *SetOptions) Run() error {
	config := o.Config()

	// @step: we replace current if there with the path of the profile
	if strings.HasPrefix(o.Path, "current.") {
		o.Path = "profiles." + config.CurrentProfile + "." + strings.TrimPrefix(o.Path, "current.")
	}

	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(config); err != nil {
		return err
	}

	value, err := sjson.Set(b.String(), o.Path, o.Value)
	if err != nil {
		return err
	}
	if value == "" {
		return fmt.Errorf("Path %q was not found configuration", o.Path)
	}
	o.Println("Successfully updated the parameter: %q in configuration", o.Path)

	if err := json.NewDecoder(strings.NewReader(value)).Decode(config); err != nil {
		return err
	}

	return o.UpdateConfig()
}
