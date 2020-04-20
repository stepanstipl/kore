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

package apply

import (
	"bytes"
	"fmt"

	"github.com/appvia/kore/pkg/cmd/errors"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/spf13/cobra"
)

// ApplyOptions are options to apply
type ApplyOptions struct {
	cmdutil.Factory
	// Paths is the file paths to apply
	Paths []string
	// Resource is the resource to retrieve
	Resource string
	// Team is the team name
	Team string
	// Force is used to force an operation
	Force bool
}

// NewCmdApply creates and returns the apply
func NewCmdApply(factory cmdutil.Factory) *cobra.Command {
	o := &ApplyOptions{Factory: factory}

	command := &cobra.Command{
		Use:   "apply",
		Short: "Allows you to apply one of more resources to the api",
		Run:   cmdutil.DefaultRunFunc(o),
	}

	flags := command.Flags()
	flags.StringSliceVarP(&o.Paths, "file", "f", []string{}, "path to file containing resource definition/s ('-' for stdin) `PATH`")

	return command
}

// Validate checks the inputs
func (o *ApplyOptions) Validate() error {
	if len(o.Paths) <= 0 {
		return errors.ErrNoFiles
	}

	return nil
}

// Run implements the apply action
func (o *ApplyOptions) Run() error {
	for _, file := range o.Paths {
		// @step: read in the content of the file
		content, err := utils.ReadFileOrStdin(o.Stdin(), file)
		if err != nil {
			return err
		}

		resources, err := cmdutil.ParseDocument(o, bytes.NewReader(content))
		if err != nil {
			return err
		}
		for _, x := range resources {
			// @step: check if the resource exists
			name := x.Object.GetName()
			namespace := x.Object.GetNamespace()

			// @step: create a request to check the status
			kind := x.Object.GetKind()
			groupversion := x.Object.GetObjectKind().GroupVersionKind().GroupVersion()
			request := o.Client().Resource(kind).Name(name)

			// @step: check the resource scope
			if x.Resource.IsTeamScoped() {
				// we set the team namespace to the resource namespace of team selected
				request.Team(func() string {
					if x.Object.GetNamespace() != "" {
						return x.Object.GetNamespace()
					}

					return o.Team
				}())

				switch {
				case namespace == "" && o.Team == "":
					return errors.ErrTeamMissing
				case namespace != "" && o.Team != "" && o.Team != namespace:
					return errors.NewConflictError("resource %s defines a different teams namespace", namespace)
				case namespace == "" && o.Team != "":
					x.Object.SetNamespace(o.Team)
				}
			}

			// @step: check if the resource exists already
			current := &unstructured.Unstructured{}
			current.SetGroupVersionKind(x.Object.GetObjectKind().GroupVersionKind())
			existing, err := request.Result(current).Exists()
			if err != nil {
				return err
			}

			// @step: attempt to apply the resource
			if err := request.Payload(x.Object).Result(x.Object).Update().Error(); err != nil {
				return err
			}

			// @step: if we had an exiting resource, we an use the revision to check
			// if the resource has changed
			state := "configured"
			if existing && (current.GetResourceVersion() == x.Object.GetResourceVersion()) {
				state = "no changes"
			}

			endpoint := fmt.Sprintf("%s/%s", groupversion, x.Object.GetName())

			o.Println("%s %s", endpoint, state)
		}
	}

	return nil
}
