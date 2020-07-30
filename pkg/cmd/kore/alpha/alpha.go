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

package alpha

import (
	"github.com/appvia/kore/pkg/cmd/kore/featuregates"
	"github.com/appvia/kore/pkg/cmd/kore/local"
	"github.com/appvia/kore/pkg/cmd/kore/patch"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

// NewCmdAlpha creates and returns the alpha command
func NewCmdAlpha(factory cmdutil.Factory) *cobra.Command {
	command := &cobra.Command{
		Use:                   "alpha",
		DisableFlagsInUseLine: true,
		Short:                 "Experimental services and operations",
		Run:                   cmdutil.RunHelp,
	}

	command.AddCommand(
		patch.NewCmdPatch(factory),
		featuregates.NewCmdFeatureGates(factory),
		local.NewCmdBootstrap(factory),
		NewCmdAlphaAuthorize(factory),
	)

	return command
}
