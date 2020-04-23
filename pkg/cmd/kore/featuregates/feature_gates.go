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

package featuregates

import (
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils/render"
	"github.com/spf13/cobra"
)

type FeatureGatesOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	// Headers indicates no headers on the table output
	Headers bool
	// Output is the output format
	Output string
}

// NewCmdFeatureGates creates and returns the profile list command
func NewCmdFeatureGates(factory cmdutil.Factory) *cobra.Command {
	o := &FeatureGatesOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "feature-gates",
		Short:   "List all feature gates",
		Run:     cmdutil.DefaultRunFunc(o),
		Example: "kore alpha feature-gates",
	}

	command.AddCommand(
		NewCmdEnabledFeatureGate(factory),
		NewCmdDisableFeatureGate(factory),
	)

	return command
}

// Run implements the action
func (o *FeatureGatesOptions) Run() error {
	type featureGate struct {
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
	}
	var list []featureGate

	for k, v := range o.Config().FeatureGates {
		list = append(list, featureGate{
			Name:    k,
			Enabled: v,
		})
	}

	return render.Render().
		Writer(o.Writer()).
		Format(o.Output).
		ShowHeaders(o.Headers).
		Resource(render.FromStruct(&list)).
		Printer(
			render.Column("Feature gate", "name"),
			render.Column("Enabled", "enabled"),
		).Do()
}
