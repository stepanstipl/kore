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

package delete

import (
	"bytes"
	"fmt"

	"github.com/appvia/kore/pkg/client"
	"github.com/appvia/kore/pkg/cmd/errors"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils"

	"github.com/spf13/cobra"
)

type DeleteOptions struct {
	cmdutil.Factory
	// Name is an optional name for the resource
	Name string
	// NoWait indicates if we should wait for a resource to provision
	NoWait bool
	// Paths is a collection of files to delete from
	Paths []string
	// Resource is the resource to retrieve
	Resource string
	// Team string
	Team string
}

// NewCmdDelete creates and returns the delete command
func NewCmdDelete(factory cmdutil.Factory) *cobra.Command {
	o := &DeleteOptions{Factory: factory}

	// @step: retrieve a list of known resources
	possible, _ := factory.Resources().Names()

	command := &cobra.Command{
		Use:                   "delete",
		Aliases:               []string{"rm"},
		DisableFlagsInUseLine: true,
		Short:                 "Allows you to delete one of more resources in kore",
		Run:                   cmdutil.DefaultRunFunc(o),

		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.BashCompDirective) {
			switch len(args) {
			case 0:
				return append(possible, "admin"), cobra.BashCompDirectiveNoFileComp
			case 1:
				suggestions, err := o.Resources().LookupResourceNames(args[0], cmdutil.GetTeam(cmd))
				if err != nil {
					return nil, cobra.BashCompDirectiveError
				}

				return suggestions, cobra.BashCompDirectiveNoFileComp
			}

			return nil, cobra.BashCompDirectiveNoFileComp
		},
	}

	command.Flags().StringSliceVarP(&o.Paths, "file", "f", []string{}, "path to file containing resource definition/s ('-' for stdin) `PATH`")

	command.AddCommand(
		NewCmdDeleteAdmin(factory),
	)

	return command
}

// Run implements the action
func (o *DeleteOptions) Run() error {
	if len(o.Paths) > 0 {
		return o.DeleteByPath()
	}

	plural := utils.Pluralize(o.Resource)

	// @step: we lookup the resource type
	resource, err := o.Resources().Lookup(plural)
	if err != nil {
		return err
	}

	request := o.Client().Resource(plural).Name(o.Name)
	if resource.IsTeamScoped() {
		request.Team(o.Team)
	}

	return o.WaitForDeletion(request, o.Name, o.NoWait)
}

// Validate checks the options
func (o *DeleteOptions) Validate() error {
	if len(o.Paths) > 0 {
		return nil
	}
	if o.Resource == "" {
		return errors.ErrMissingResourceName
	}
	if o.Name == "" {
		return errors.ErrMissingResourceName
	}

	resource, err := o.Resources().Lookup(utils.Pluralize(o.Resource))
	if err != nil {
		return err
	}
	if resource.IsTeamScoped() && o.Team == "" {
		return errors.ErrTeamMissing
	}

	return nil
}

// DeleteByPath iterate and delete from the file
func (o *DeleteOptions) DeleteByPath() error {
	for _, file := range o.Paths {
		// @step: read in the content of the file
		content, err := utils.ReadFileOrStdin(o.Stdin(), file)
		if err != nil {
			return err
		}

		documents, err := cmdutil.ParseDocument(o, bytes.NewReader(content))
		if err != nil {
			return err
		}

		for _, x := range documents {
			namespace := x.Object.GetNamespace()

			// @step: check the resource scope
			if x.Resource.IsTeamScoped() {
				switch {
				case namespace == "" && o.Team == "":
					return errors.ErrTeamMissing
				case namespace != "" && o.Team != "" && o.Team != namespace:
					return errors.NewConflictError("resource %s defines a different teams namespace", namespace)
				case namespace == "" && o.Team != "":
					x.Object.SetNamespace(o.Team)
				}
			}

			// @step: build the status line
			name := x.Object.GetName()
			kind := x.Object.GetObjectKind().GroupVersionKind().Kind
			groupversion := x.Object.GetObjectKind().GroupVersionKind().GroupVersion()

			// @step: construct a request for the resource
			request := o.Client().Resource(kind).Name(name)
			if x.Resource.IsTeamScoped() {
				request.Team(o.Team)
			}

			endpoint := fmt.Sprintf("%s/%s", groupversion, name)

			// @step: attempt to delete the resource
			if err := request.Payload(x.Object).Result(x.Object).Delete().Error(); err != nil {
				if !client.IsNotFound(err) {
					return err
				}
				o.Println("%s not found", endpoint)

				continue
			}
			o.Println("%s deleted", endpoint)
		}
	}

	return nil
}
