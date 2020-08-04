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

package kore

import (
	"os"

	"github.com/appvia/kore/pkg/controllers/management/cluster"
	log "github.com/sirupsen/logrus"

	"github.com/appvia/kore/pkg/schema"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	upLongDescription = `
Allows to you retrieve the resources from the kore api. The command format
is <resource> [name]. When the optional name is not provided we will return
a full listing of all the <resource>s from the API. Examples of resource types
are users, teams, gkes, clusters amongst a few.

You can list all the available resource via $ kore api-resources

Though for a better experience all the resource are autocompletes for you.
Take a look at $ kore completion for details
`
	upExamples = `
# List users:
$ kore get users

#Get information about a specific user:
$ kore get user admin [-o yaml]
`
)

// UpOptions the are the options for a get command
type UpOptions struct {
	cmdutil.Factory

	// Paths is the manifest paths to apply
	Paths []string
}

// NewCmdUp creates and returns the up command
func NewCmdUp(factory cmdutil.Factory) *cobra.Command {
	o := &UpOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "up",
		Long:    upExamples,
		Example: "up examples",
		Run:     cmdutil.DefaultRunFunc(o),
	}

	flags := command.Flags()
	flags.StringSliceVarP(&o.Paths, "file", "f", []string{}, "path to file containing resource definition/s ('-' for stdin) `PATH`")

	return command
}

// Validate is used to validate the options
func (o *UpOptions) Validate() error {
	return nil
}

// Run implements the action
func (o *UpOptions) Run() error {
	if len(o.Paths) == 0 {
		path, err := o.findFile("kore.yml", "kore.yaml")
		if err != nil {
			return err
		}
		if path != "" {
			o.Paths = []string{path}
		}
	}

	if len(o.Paths) == 0 {
		path, err := runInit()
		if err != nil {
			return err
		}
		o.Paths = []string{path}
	}

	var objects []runtime.Object
	for _, file := range o.Paths {
		// @step: read in the content of the file
		content, err := utils.ReadFileOrStdin(o.Stdin(), file)
		if err != nil {
			return err
		}

		fileObjects := &kubernetes.Objects{}
		if err := fileObjects.UnmarshalYAML(content); err != nil {
			return err
		}

		objects = append(objects, *fileObjects...)
	}

	client := fake.NewFakeClientWithScheme(schema.GetScheme(), objects...)

	controller := cluster.NewController(log.StandardLogger())
	controller.Initialize()

	return nil
}

func (o *UpOptions) findFile(paths ...string) (string, error) {
	var fi os.FileInfo
	var err error

	for _, path := range paths {
		fi, err = os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", err
		}

		return fi.Name(), nil
	}

	return "", nil
}
