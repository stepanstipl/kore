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

package patch

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/appvia/kore/pkg/client"
	"github.com/appvia/kore/pkg/cmd/errors"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/jsonutils"

	"github.com/spf13/cobra"
	"github.com/tidwall/sjson"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	longDescription = `
Patch allows you to apply patches to the resource managed in kore to add or
remove values from a resource.
`

	patchExamples = `
# Update the size in a cluster resource
$ kore alpha patch clusters test spec.configuration.size 1 [-t <team>]

# Update the allowed subnets
$ kore alpha patch clusters test spec.configuration.authProxyAllowedIPs.-1 127.0.0.0/8 [-t team]

#Remove the value
$ kore alpha patch clusters test spec.configuration.authProxyAllowedIPs.0
`
)

// PatchOptions are the options for patch comment
type PatchOptions struct {
	cmdutil.Factory
	// Name is an optional name for the resource
	Name string
	// Resource is the resource to retrieve
	Resource string
	// Team is the team name
	Team string
	// Key is the json path to patch
	Key string
	// Value is the value to set
	Value string
}

// NewCmdPatch creates and returns the patch command
func NewCmdPatch(factory cmdutil.Factory) *cobra.Command {
	o := &PatchOptions{Factory: factory}

	// @step: retrieve a list of known resources
	possible, _ := factory.Resources().Names()

	command := &cobra.Command{
		Use:     "patch",
		Short:   "Allows you to patch resource in kore",
		Long:    longDescription,
		Example: patchExamples,

		Run: func(cmd *cobra.Command, args []string) {
			o.Key = cmd.Flags().Arg(2)
			o.Value = cmd.Flags().Arg(3)

			cmdutil.DefaultRunFunc(o)(cmd, args)
		},

		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			switch len(args) {
			case 0:
				return possible, cobra.ShellCompDirectiveNoFileComp
			case 1:
				suggestions, err := o.Resources().LookupResourceNames(cmd.Flags().Arg(0), cmdutil.GetTeam(cmd))
				if err != nil {
					return nil, cobra.ShellCompDirectiveError
				}

				return suggestions, cobra.ShellCompDirectiveNoFileComp
			}

			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return command
}

// Validate is used to validate the options
func (o *PatchOptions) Validate() error {
	if o.Resource == "" {
		return errors.ErrMissingResource
	}
	if o.Name == "" {
		return errors.ErrMissingResourceName
	}
	if o.Key == "" {
		return errors.NewInvalidParamError("key", "missing")
	}
	resource, err := o.Resources().Lookup(o.Resource)
	if err != nil {
		return err
	}
	if resource.IsTeamScoped() && o.Team == "" {
		return errors.ErrTeamMissing
	}

	return nil
}

// Run implements the action
func (o *PatchOptions) Run() error {
	u := &unstructured.Unstructured{}

	// @step: retrieve the resource from the API
	resource, err := o.Resources().Lookup(o.Resource)
	if err != nil {
		return err
	}
	var request client.RestInterface
	if resource.IsScoped(cmdutil.TeamScope) {
		request = o.ClientWithTeamResource(o.Team, resource).Result(u).Name(o.Name).Get()
	} else {
		request = o.ClientWithResource(resource).Result(u).Name(o.Name).Get()
	}
	if err := request.Error(); err != nil {
		return err
	}
	revision := u.GetResourceVersion()

	content, err := ioutil.ReadAll(request.Body())
	if err != nil {
		return err
	}

	// @step: apply the patch to the json
	update, err := func() ([]byte, error) {
		if o.Value == "" {
			return sjson.DeleteBytes(content, o.Key)
		}

		return jsonutils.SetJSONProperty(content, o.Key, o.Value)
	}()
	if err != nil {
		return err
	}

	if err := json.NewDecoder(bytes.NewReader(update)).Decode(u); err != nil {
		return err
	}

	err = request.
		Name(o.Name).
		Payload(u).
		Result(u).
		Update().
		Error()
	if err != nil {
		return err
	}

	if revision != u.GetResourceVersion() {
		o.Println("%s configured", utils.GetUnstructuredSelfLink(u))

		return nil
	}
	o.Println("%s no changes", utils.GetUnstructuredSelfLink(u))

	return nil
}
