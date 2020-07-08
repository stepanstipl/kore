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
	"regexp"
	"strings"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	"github.com/appvia/kore/pkg/cmd/errors"
	cmdutils "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/kore"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spf13/cobra"
)

var createConfigLongDescription = `
Provides the ability to create a named configuration object, which contains one or more key-value pairs.
`

var createConfigExamples = `
Create a configuration object with a single key-value pair:

$ kore create config myconfig --keyval param1=value1

Create a configuration object with multiple key-value pairs:

$ kore create config myconfig --keyval param1=value1 --keyval param2=value2
`

// CreateConfigOptions is used to provision a team
type CreateConfigOptions struct {
	cmdutils.Factory
	cmdutils.DefaultHandler
	// Name is the config name
	Name string
	// KeyVals is the list of keys provided for a name
	KeyVals []string
}

// NewCmdCreateConfig returns the create user command
func NewCmdCreateConfig(factory cmdutils.Factory) *cobra.Command {
	o := &CreateConfigOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "config <NAME>",
		Short:   "Adds new config to kore",
		Long:    createConfigLongDescription,
		Example: createConfigExamples,
		PreRunE: cmdutils.RequireName,
		Run:     cmdutils.DefaultRunFunc(o),
	}

	flags := command.Flags()

	flags.StringSliceVar(&o.KeyVals, "keyval", []string{}, "Config key value pair `KEY=VALUE`")

	cmdutils.MustMarkFlagRequired(command, "keyval")

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

	match := regexp.MustCompile("^.+?=.*$")
	for _, x := range o.KeyVals {
		if !match.MatchString(x) {
			return errors.NewInvalidParamError("config", x)
		}
	}

	pairs := make(map[string]string)
	for _, pair := range o.KeyVals {
		items := strings.SplitN(pair, "=", 2)
		pairs[items[0]] = items[1]
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

	return o.ClientWithResource(o.Resources().MustLookup("config")).
		Name(o.Name).
		Payload(config).
		Update().
		Error()
}
