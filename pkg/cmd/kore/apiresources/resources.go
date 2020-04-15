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

package apiresources

import (
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils/render"

	"github.com/spf13/cobra"
)

// APIResourceOptions is used to provision a team
type APIResourceOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	// Output is the output format
	Output string
	// Headers indicates no headers on the table output
	Headers bool
}

// NewCmdAPIResources returns the create local command
func NewCmdAPIResources(factory cmdutil.Factory) *cobra.Command {
	o := &APIResourceOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "api-resources",
		Short:   "Returns a list of all the available resources in kore",
		Example: "kore api-resource",

		Run: cmdutil.DefaultRunFunc(o),
	}

	return command
}

// Run implements the action
func (o *APIResourceOptions) Run() error {
	list, err := o.Resources().List()
	if err != nil {
		return err
	}

	return render.Render().
		Writer(o.Writer()).
		Format(o.Output).
		ShowHeaders(o.Headers).
		Resource(render.FromStruct(list)).
		Printer(
			render.Column("Name", "name"),
			render.Column("Shortnames", "shortName", func(value string) string {
				if value == "Unknown" {
					return "None"
				}
				return value
			}),
			render.Column("Scope", "scope"),
			render.Column("API Group", "groupVersion"),
			render.Column("Kind", "kind"),
		).Do()
}
