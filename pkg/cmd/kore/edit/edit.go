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

package edit

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/appvia/kore/pkg/cmd/errors"
	"github.com/appvia/kore/pkg/cmd/kore/apply"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

const (
	// DefaultEditor is the default prog to use for editing
	DefaultEditor = "vi"
)

// EditOptions are the options for edit comment
type EditOptions struct {
	cmdutil.Factory
	// Name is an optional name for the resource
	Name string
	// Resource is the resource to retrieve
	Resource string
	// Team is the team name
	Team string
}

// NewCmdEdit creates and returns the edit command
func NewCmdEdit(factory cmdutil.Factory) *cobra.Command {
	o := &EditOptions{Factory: factory}

	// @step: retrieve a list of known resources
	possible, _ := factory.Resources().Names()

	command := &cobra.Command{
		Use:     "edit",
		Short:   "Allows you to edit resource in kore",
		Example: "kore edit <resource> <name> [options]",

		Run: cmdutil.DefaultRunFunc(o),

		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.BashCompDirective) {
			switch len(args) {
			case 0:
				return possible, cobra.BashCompDirectiveNoFileComp
			case 1:
				resource := o.Resources().ResolveShorthand(cmd.Flags().Arg(0))
				suggestions, err := o.Resources().LookupResourceNames(resource, cmdutil.GetTeam(cmd))
				if err != nil {
					return nil, cobra.BashCompDirectiveError
				}

				return suggestions, cobra.BashCompDirectiveNoFileComp
			}

			return nil, cobra.BashCompDirectiveNoFileComp
		},
	}

	return command
}

// Validate is used to validate the options
func (o *EditOptions) Validate() error {
	if o.Resource == "" {
		return errors.ErrMissingResource
	}
	if o.Name == "" {
		return errors.ErrMissingResourceName
	}
	// @step: lookup the resource from the cache
	resource, err := o.Resources().Lookup(utils.Pluralize(o.Resource))
	if err != nil {
		return err
	}
	// @step: if the resource if team space, lets ensure we have the team selector
	if resource.IsTeamScoped() && o.Team == "" {
		return errors.ErrTeamMissing
	}

	return nil
}

// Run implements the action
// @TODO need to add a way to quit without changing
func (o *EditOptions) Run() error {
	plural := utils.Pluralize(o.Resource)

	// @step: lookup the resource from the cache
	resource, err := o.Resources().Lookup(plural)
	if err != nil {
		return err
	}

	// @step: retrieve the resource
	object, err := o.GetResource(resource)
	if err != nil {
		return err
	}

	// @step: write the output a file - technically if were using vim we can
	// pass on the stdin
	tmpf, err := ioutil.TempFile(os.TempDir(), "kore-edit-")
	if err != nil {
		return err
	}
	defer os.Remove(tmpf.Name())

	if _, err := tmpf.Write(object); err != nil {
		return err
	}
	tmpf.Close()

	editor := exec.Command(o.GetEditor(), tmpf.Name())
	editor.Stdin = os.Stdin
	editor.Stdout = os.Stdout
	editor.Stderr = os.Stderr

	if err := editor.Run(); err != nil {
		return err
	}

	content, err := ioutil.ReadFile(tmpf.Name())
	if err != nil {
		return err
	}
	encoded, err := yaml.YAMLToJSON(content)
	if err != nil {
		return err
	}
	o.SetStdin(bytes.NewReader(encoded))

	opts := &apply.ApplyOptions{
		Factory: o,
		Paths:   []string{"-"},
		Team:    o.Team,
	}

	return cmdutil.ExecuteHandler(opts)
}

// GetResource retrieve the resource for us
func (o *EditOptions) GetResource(resource *cmdutil.Resource) ([]byte, error) {
	// @step: we need to construct the request
	request := o.Client().Resource(utils.Pluralize(o.Resource)).Name(o.Name)
	if resource.IsTeamScoped() {
		request.Team(o.Team)
	}

	// @step: we perform the get request against the api
	if err := request.Get().Error(); err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(request.Body())
	if err != nil {
		return nil, err
	}

	encoded, err := yaml.JSONToYAML(content)
	if err != nil {
		return nil, err
	}

	return encoded, nil
}

// GetEditor attempts to get the user defined editor
func (o *EditOptions) GetEditor() string {
	cmd := os.Getenv("EDITOR")
	if cmd == "" {
		return DefaultEditor
	}

	return cmd
}
