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

package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/tidwall/sjson"

	"github.com/appvia/kore/pkg/utils"

	"github.com/appvia/kore/pkg/cmd/errors"
	"github.com/appvia/kore/pkg/utils/render"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"

	"github.com/spf13/cobra"
)

// RunHelp is shorthand for displaying the usage
func RunHelp(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}

// RunHelpE is shorthand for display the usage but with error return
func RunHelpE(cmd *cobra.Command, args []string) error {
	_ = cmd.Help()

	return nil
}

// GetVerbose return the verbose flag
func GetVerbose(cmd *cobra.Command) bool {
	return GetFlagBool(cmd, "verbose")
}

// GetDebug return the verbose flag
func GetDebug(cmd *cobra.Command) bool {
	return GetFlagBool(cmd, "debug")
}

// GetNoWait returns the no-wait flag
func GetNoWait(cmd *cobra.Command) bool {
	return GetFlagBool(cmd, "no-wait")
}

// GetTeam returns the team flag
func GetTeam(cmd *cobra.Command) string {
	return GetFlagString(cmd, "team")
}

// GetFlagString returns a flag string
func GetFlagString(cmd *cobra.Command, name string) string {
	v, _ := cmd.Flags().GetString(name)

	return v
}

// GetFlagBool returns a flag boolean
func GetFlagBool(cmd *cobra.Command, name string) bool {
	v, _ := cmd.Flags().GetBool(name)

	return v
}

// PreRunEFilters provides a handler to check the options
func PreRunEFilters() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return nil
	}
}

// RequireName ensures a name positional argument
func RequireName(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		_ = cmd.Help()

		return errors.ErrMissingResourceName
	}

	return nil
}

// MustMarkFlagRequired calls MarkFlagRequired and panics if it returns an error
func MustMarkFlagRequired(cmd *cobra.Command, name string) {
	if err := cmd.MarkFlagRequired(name); err != nil {
		panic(err)
	}
}

// MustRegisterFlagCompletionFunc ensures we never require a missing flag
func MustRegisterFlagCompletionFunc(
	cmd *cobra.Command,
	flagName string,
	f func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective),
) {
	if err := cmd.RegisterFlagCompletionFunc(flagName, f); err != nil {
		panic(err)
	}
}

// ParseDocument returns a collection of parsed documents and the api endpoints
func ParseDocument(f Factory, src io.Reader) ([]*ResourceDocument, error) {
	var list []*ResourceDocument

	// @step: read in the content of the file
	content, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}
	// @step: split the yaml documents up
	splitter := regexp.MustCompile("(?m)^---\n")
	documents := splitter.Split(string(content), -1)

	for _, x := range documents {
		if x == "" {
			continue
		}

		doc, err := yaml.YAMLToJSON([]byte(x))
		if err != nil {
			return nil, err
		}

		// @step: attempt to read the document into an unstructured
		u := &unstructured.Unstructured{}
		if err := u.UnmarshalJSON(doc); err != nil {
			return nil, err
		}

		// @checks: to ensure the resource is properly defined
		if u.GetName() == "" {
			return nil, errors.ErrMissingResourceName
		}
		if u.GetKind() == "" {
			return nil, errors.ErrMissingResourceKind
		}
		if u.GetAPIVersion() == "" {
			return nil, errors.ErrMissingResourceVersion
		}

		displayName := utils.GetUnstructuredSelfLink(u)

		// @step: lookup the resource from the cache
		resource, err := f.Resources().Lookup(u.GetKind())
		if err != nil {
			return nil, fmt.Errorf("%s: %w", displayName, err)
		}

		list = append(list, &ResourceDocument{Object: u, Resource: resource})
	}

	return list, nil
}

// ConvertColumnsToRender converts the resource columns to render columns
func ConvertColumnsToRender(columns []Column) []render.PrinterColumnFunc {
	list := make([]render.PrinterColumnFunc, len(columns))

	for i, c := range columns {
		switch c.Format {
		case "age":
			list[i] = render.Column(c.Name, c.Path, render.Age())
		default:
			list[i] = render.Column(c.Name, c.Path)
		}
	}

	return list
}

func PatchJSON(document string, cliValues []string) (string, error) {
	var parameterRegexp = regexp.MustCompile(`\s*=\s*`)

	params := make(map[string]string)

	for _, x := range cliValues {
		e := parameterRegexp.Split(strings.TrimSpace(x), 2)

		if len(e) != 2 || e[0] == "" || e[1] == "" {
			return "", errors.NewInvalidParamError("param", x)
		}
		params[e[0]] = e[1]
	}

	var err error
	for key, value := range params {
		if strings.HasPrefix(value, "{") || strings.HasPrefix(value, "[") {
			document, err = sjson.SetRaw(document, key, value)
			if err != nil {
				return "", err
			}
		} else {
			document, err = sjson.Set(document, key, func(v string) interface{} {
				if num, err := strconv.ParseFloat(v, 64); err == nil {
					return num
				}
				return v
			}(value))
			if err != nil {
				return "", err
			}
		}
	}

	return document, nil
}
