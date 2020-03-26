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
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	"gopkg.in/yaml.v2"

	"github.com/urfave/cli/v2"
)

var (
	createSecretLongDescription = `
Provides the ability to create secrets in the kore, from files, environments
files and literals.

 $ korectl create secret <name> -t <team> [options]

Examples:
 # Create a secret from a file
 $ korectl create secret gke --from-file=<key>=<filename>
`
)

// GetCreateSecretCommand creates and returns the create secret command
func GetCreateSecretCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "secret",
		Aliases:     []string{"secrets"},
		Description: formatLongDescription(createSecretLongDescription),
		Usage:       "Creates a secret in kore",
		ArgsUsage:   "<name> [options]",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "description",
				Aliases:  []string{"d"},
				Usage:    "A description for this secret `DESC`",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "type",
				Usage: "indicates the type of secret you are generating `NAME`",
				Value: "generic",
			},
			&cli.StringSliceFlag{
				Name:  "from-literal",
				Usage: "adding a literal to the secret `KEY=NAME`",
			},
			&cli.StringFlag{
				Name:  "from-file",
				Usage: "builds a secret from the key reference `KEY=NAME`",
			},
			&cli.StringFlag{
				Name:  "from-env-file",
				Usage: "builds a secret from the environment file, format NAME=VALUE `PATH`",
			},
			&cli.BoolFlag{
				Name:  "force",
				Usage: "overwrite the secret if it already exists (defaults: false) `BOOL`",
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "generate the cluster specification but does not apply `BOOL`",
			},
		},

		Before: func(ctx *cli.Context) error {
			if !ctx.Args().Present() {
				return errors.New("the secret should have a name")
			}

			return nil
		},

		Action: func(ctx *cli.Context) error {
			name := ctx.Args().First()
			kind := "secert"
			team := ctx.String("team")
			force := ctx.Bool("force")
			nowait := ctx.Bool("no-wait")
			dryrun := ctx.Bool("dry-run")

			var secret *configv1.Secret
			var err error

			switch {
			case team == "":
				return errors.New("you must specify a team")
			case ctx.String("from-env-file") != "":
				secret, err = createSecretFromEnvironmentFile(ctx.String("from-env-file"))
			case ctx.String("from-file") != "":
				secret, err = createSecretFromFile(ctx.String("from-file"))
			case len(ctx.StringSlice("from-literal")) > 0:
				secret, err = createSecretFromLiterals(ctx.StringSlice("from-literal"))
			default:
				return errors.New("you must choose to create from --from-env-file, --from-file or --from-literal")
			}
			if err != nil {
				return fmt.Errorf("failed to create secret: %s", err)
			}

			secret.Spec.Description = ctx.String("description")
			secret.Spec.Type = ctx.String("type")

			if dryrun {
				return yaml.NewEncoder(os.Stdout).Encode(secret)
			}

			if !force {
				if found, err := TeamResourceExists(config, team, kind, name); err != nil {
					return err
				} else if found {
					return fmt.Errorf("%q already exists, please use --force if your sure you want to update", name)
				}
			}

			if err := CreateTeamResource(config, team, kind, name, secret); err != nil {
				return err
			}

			return WaitForResourceCheck(context.Background(), config, team, kind, name, nowait)
		},
	}
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
