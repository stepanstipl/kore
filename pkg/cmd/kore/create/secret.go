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
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	cmdutils "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

var (
	createSecretLongDescription = `
Provides the ability to create secrets in the kore, from files, environments
files and literals.

$ kore create secret <name> -t <team> [options]

Examples:
# Create a secret from a file
$ kore create secret gke --from-file=<key>=<filename>
`
)

// CreateSecretOptions is used to provision a team
type CreateSecretOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	// Description is a summary for the secret is used for
	Description string
	// EnvFile is a environment file to gen secret
	EnvFile string
	// File is file path to gen secret
	File string
	// Force is used to force an operation
	Force bool
	// Name is the name of the secret
	Name string
	// Literals is a collection of secret literals
	Literals []string
	// Team is the team name
	Team string
	// Type is the type of secret
	Type string
	// Username is the username to add
	Username string
}

// NewCmdCreateSecret returns the create secret command
func NewCmdCreateSecret(factory cmdutil.Factory) *cobra.Command {
	o := &CreateSecretOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "secret",
		Short:   "Creates a secret in kore",
		Example: "kore create secret <options> [-t team]",
		Long:    createSecretLongDescription,
		PreRunE: cmdutil.RequireName,
		Run:     cmdutil.DefaultRunFunc(o),
	}

	flags := command.Flags()
	flags.StringVar(&o.Type, "type", "generic", "is the type of secret you are creating `TYPE`")
	flags.StringVarP(&o.Description, "description", "d", "", "a description for this secret `DESC`")
	flags.StringSliceVar(&o.Literals, "from-literal", []string{}, "adding a literal to the secret `KEY=NAME")
	flags.StringVar(&o.File, "from-file", "", "builds a secret from the key reference `KEY=PATH`")
	flags.StringVar(&o.File, "from-env-file", "", "builds a secret from the environment file, `KEY=PATH`")

	cmdutils.MustMarkFlagRequired(command, "description")

	return command
}

// Run implements the action
func (o *CreateSecretOptions) Run() error {
	var secret *configv1.Secret
	var err error

	switch {
	case o.Team == "":
		return errors.New("you must specify a team")
	case o.EnvFile != "":
		secret, err = createSecretFromEnvironmentFile(o.EnvFile)
	case o.File != "":
		secret, err = createSecretFromFile(o.File)
	case len(o.Literals) > 0:
		secret, err = createSecretFromLiterals(o.Literals)
	default:
		return errors.New("you must choose to create from --from-env-file, --from-file or --from-literal")
	}
	if err != nil {
		return fmt.Errorf("failed to create secret: %s", err)
	}

	secret.Spec.Description = o.Description
	secret.Spec.Type = o.Type

	if !o.Force {
		if found, err := o.ClientWithTeamResource(o.Team, o.Resources().MustLookup("secret")).Name(o.Name).Exists(); err != nil {
			return err
		} else if found {
			return fmt.Errorf("%q already exists, please use --force if your sure you want to update", o.Name)
		}
	}

	return o.ClientWithTeamResource(o.Team, o.Resources().MustLookup("secret")).
		Payload(secret).
		Update().
		Error()
}

// createSecretFromLiterals creates some secret from a collection of literals
func createSecretFromLiterals(keypairs []string) (*configv1.Secret, error) {
	secret := &configv1.Secret{
		Spec: configv1.SecretSpec{Data: make(map[string]string)},
	}

	filter := regexp.MustCompile("^[a-zA-Z0-9_]*=.*$")

	for _, kv := range keypairs {
		if !filter.MatchString(kv) {
			return nil, fmt.Errorf("invalid value: %s must conform to: %s", kv, filter)
		}
		items := strings.Split(kv, "=")
		secret.Spec.Data[items[0]] = base64.StdEncoding.EncodeToString([]byte(items[1]))
	}

	return secret, nil
}

// createSecretFromFile generates a secret from the file
func createSecretFromFile(keypair string) (*configv1.Secret, error) {
	filter := regexp.MustCompile("^.*=.*$")
	if !filter.MatchString(keypair) {
		return nil, fmt.Errorf("invalid value, must be KEY=FILEPATH")
	}

	//@step: extract the values
	items := strings.Split(keypair, "=")
	key := items[0]
	path := items[1]

	// @step: reading in the file content
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("trying to read the file: %s, error: %s", path, err)
	}

	return &configv1.Secret{
		Spec: configv1.SecretSpec{
			Data: map[string]string{
				key: base64.StdEncoding.EncodeToString(content),
			},
		},
	}, nil
}

// createSecretFromEnvironmentFile generates a secret from the environment file
func createSecretFromEnvironmentFile(path string) (*configv1.Secret, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	filter := regexp.MustCompile("^[a-zA-Z0-9_]*=.*$")

	secret := &configv1.Secret{
		Spec: configv1.SecretSpec{
			Data: make(map[string]string),
		},
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		switch {
		case strings.HasPrefix(scanner.Text(), "#"):
			continue
		case strings.HasPrefix(scanner.Text(), " "):
			continue
		case scanner.Text() == "":
			continue
		case !filter.MatchString(scanner.Text()):
			return nil, fmt.Errorf("invalid format: %s, must conform to: %s", filter, scanner.Text())
		}

		e := strings.Split(scanner.Text(), "=")

		secret.Spec.Data[e[0]] = base64.StdEncoding.EncodeToString([]byte(e[1]))
	}

	return secret, nil
}
