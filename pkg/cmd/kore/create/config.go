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

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	cmdutils "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/render"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spf13/cobra"
)

var (
	configLongDesciption = `
Provides the ability to create a key value pair in kore's config.

# Create the key value pair
$ kore create config <name> -k <key> -v <value>
# Create multiple key value pairs
$ kore create config <name> -k <key> -v <value> -k <key2> -v <value2>
`
)

// CreateConfigOptions is used to provision a team
type CreateConfigOptions struct {
	cmdutils.Factory
	cmdutils.DefaultHandler
	// Name is the username to add
	Name string
	// Keys is the list of keys provided for a name
	Keys []string
	// Keys is the list of values provided for a name
	Values []string
	// DryRun indicates we only dryrun the resources
	DryRun bool
}

// NewCmdCreateConfig returns the create user command
func NewCmdCreateConfig(factory cmdutils.Factory) *cobra.Command {
	o := &CreateConfigOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "config",
		Short:   "Adds new config to kore",
		Long:    configLongDesciption,
		Example: "kore create config <name> -k <key> -v <value>",
		PreRunE: cmdutils.RequireName,
		Run:     cmdutils.DefaultRunFunc(o),
	}

	flags := command.Flags()

	flags.StringSliceVarP(&o.Keys, "key", "k", []string{}, "key to compliment a name")
	flags.StringSliceVarP(&o.Values, "value", "v", []string{}, "value to compliment a name")
	flags.BoolVar(&o.DryRun, "dry-run", false, "shows the resource but does not apply or create (defaults: false)")

	cmdutils.MustMarkFlagRequired(command, "key")
	cmdutils.MustMarkFlagRequired(command, "value")

	return command
}

// Run implements the action
func (o *CreateConfigOptions) Run() error {
	found, err := o.ClientWithResource(o.Resources().MustLookup("configs")).Name(o.Name).Exists()
	if err != nil {
		return err
	}
	if found {
		return fmt.Errorf("%q already exists, please edit instead", o.Name)
	}

	pairs := make(map[string]string)
	if len(o.Keys) == len(o.Values) {
		for i := range o.Keys {
			pairs[o.Keys[i]] = o.Values[i]
		}
	} else {
		return fmt.Errorf("Invalid key values pairs inputted")
	}
	config := &configv1.Config{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Config",
			APIVersion: configv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.Name,
			Namespace: kore.HubNamespace,
		},
		Spec: configv1.ConfigSpec{
			Values: pairs,
		},
	}

	if o.DryRun {
		return render.Render().
			Writer(o.Writer()).
			Format(render.FormatYAML).
			Resource(render.FromStruct(config)).
			Do()
	}

	return o.ClientWithResource(o.Resources().MustLookup("config")).
		Name(o.Name).
		Payload(config).
		Update().
		Error()
}
