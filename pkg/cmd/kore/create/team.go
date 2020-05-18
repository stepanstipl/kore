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

package create

import (
	"fmt"

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils/render"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateTeamOptions is used to provision a team
type CreateTeamOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	// Description is a description for it
	Description string
	// DryRun indicates we only dryrun the resources
	DryRun bool
	// Name is the name of the team
	Name string
	// NoWait indicates if we should wait for a resource to provision
	NoWait bool
}

// NewCmdCreateTeam returns the create team command
func NewCmdCreateTeam(factory cmdutil.Factory) *cobra.Command {
	o := &CreateTeamOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "team",
		Aliases: []string{"teams"},
		Short:   "Creates a team in kore, adding yourself as the team admin",
		Example: "kore create team <name>",
		PreRunE: cmdutil.RequireName,
		Run:     cmdutil.DefaultRunFunc(o),
	}

	flags := command.Flags()
	flags.StringVarP(&o.Description, "description", "d", "", "the description of the team")
	flags.BoolVar(&o.DryRun, "dry-run", false, "shows the resource but does not apply or create (defaults: false)")

	return command
}

// Run is responsible for creating the team
func (o CreateTeamOptions) Run() error {
	found, err := o.ClientWithResource(o.Resources().MustLookup("team")).Name(o.Name).Exists()
	if err != nil {
		return err
	}
	if found {
		return fmt.Errorf("%q already exists, please edit instead", o.Name)
	}

	// @step: create the resource from options
	team := &orgv1.Team{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Team",
			APIVersion: orgv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: o.Name,
		},
		Spec: orgv1.TeamSpec{
			Summary:     o.Name,
			Description: o.Description,
		},
	}

	if o.DryRun {
		return render.Render().
			Writer(o.Writer()).
			Format(render.FormatYAML).
			Resource(render.FromStruct(team)).
			Do()
	}

	return o.WaitForCreation(
		o.ClientWithResource(o.Resources().MustLookup("team")).
			Name(o.Name).
			Payload(team).
			Result(&orgv1.Team{}),
		o.NoWait,
	)
}
