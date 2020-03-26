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

package korectl

import (
	"fmt"

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"

	"github.com/urfave/cli/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetCreateTeamCommand returns the create team command
func GetCreateTeamCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:      "team",
		Aliases:   []string{"teams"},
		Usage:     "Creates a team",
		ArgsUsage: "TEAM",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "description",
				Usage:    "The description of the team",
				Required: false,
			},
		},

		Before: func(ctx *cli.Context) error {
			if !ctx.Args().Present() {
				return fmt.Errorf("the team should have a name")
			}

			return nil
		},

		Action: func(ctx *cli.Context) error {
			name := ctx.Args().First()
			description := ctx.String("description")
			kind := "team"
			wait := ctx.Bool("wait")

			// @step: check if the resource exist already
			if found, err := ResourceExists(config, kind, name); err != nil {
				return err
			} else if found {
				return fmt.Errorf("%q already exists, please edit instead", name)
			}

			// @step: create the resource from options
			team := &orgv1.Team{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Team",
					APIVersion: orgv1.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
				},
				Spec: orgv1.TeamSpec{
					Summary:     name,
					Description: description,
				},
			}

			req, _, err := NewRequestForResource(config, "", kind, name)
			if err != nil {
				return err
			}
			req.WithRuntimeObject(team)

			if err := req.Update(); err != nil {
				return err
			}

			return WaitForResourceCheck(ctx.Context, config, "", kind, name, wait)
		},
	}
}
