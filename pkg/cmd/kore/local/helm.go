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

package local

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/jsonutils"

	"sigs.k8s.io/yaml"
)

var (
	keypairRegex = regexp.MustCompile(`^([[:alnum:]].*)=([[:alnum:]\{\}].*)$`)
)

// GetHelmValues returns returns or prompts for the values
func (o *UpOptions) GetHelmValues(path string) (map[string]interface{}, error) {
	values := make(map[string]interface{})

	found, err := utils.FileExists(path)
	if err != nil {
		return nil, err

	}

	// @step: we retrieve the values from default or file
	switch found {
	case true:
		values, err = GetHelmValuesFromFile(path)
	default:
		values, err = GetDefaultHelmValues()
	}
	if err != nil {
		return nil, err
	}

	// @step: inject the flags if required
	if utils.Contains("version", o.FlagsChanged) {
		for _, x := range []string{"api.version", "ui.version"} {
			o.HelmValues = append(o.HelmValues, fmt.Sprintf("%s=%s", x, o.Version))
		}
	} else {
		if !found {
			for _, x := range []string{"api.version", "ui.version"} {
				o.HelmValues = append(o.HelmValues, fmt.Sprintf("%s=%s", x, o.Version))
			}
		}
	}

	// @step: marshal the values to json and apply the updates
	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(&values); err != nil {
		return nil, err
	}
	content := b.Bytes()

	for _, x := range o.HelmValues {
		e := keypairRegex.FindStringSubmatch(x)
		if len(e) < 3 {
			return nil, fmt.Errorf("invalid helm value: %q should be key=value", x)
		}
		content, err = jsonutils.SetJSONProperty(content, e[1], e[2])
		if err != nil {
			return nil, err
		}
	}
	// @step: convert the json to values for writing
	values = make(map[string]interface{})

	return values, json.NewDecoder(bytes.NewReader(content)).Decode(&values)
}

// GetHelmValuesFromFile returns the current set values
func GetHelmValuesFromFile(path string) (map[string]interface{}, error) {
	values := make(map[string]interface{})

	// @step: we read in the values.yml
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return values, yaml.Unmarshal(content, &values)
}

// GetDefaultHelmValues returns the default values required
func GetDefaultHelmValues() (map[string]interface{}, error) {
	values := DefaultHelmValues()

	a := authInfoConfig{}
	a.AuthorizeURL = os.Getenv("KORE_IDP_SERVER_URL")
	a.ClientID = os.Getenv("KORE_IDP_CLIENT_ID")
	a.ClientSecret = os.Getenv("KORE_IDP_CLIENT_SECRET")

	if a.AuthorizeURL == "" || a.ClientID == "" || a.ClientSecret == "" {
		if err := (&cmdutil.Prompts{
			&cmdutil.Prompt{Id: "Client ID", ErrMsg: "%s cannot be blank", Value: &a.ClientID},
			&cmdutil.Prompt{Id: "Client Secret", ErrMsg: "%s cannot be blank", Value: &a.ClientSecret},
			&cmdutil.Prompt{Id: "Authorization Endpoint", ErrMsg: "%s cannot be blank", Value: &a.AuthorizeURL},
		}).Collect(); err != nil {
			return nil, err
		}
	}

	values["idp"] = map[string]interface{}{
		"client_id":     a.ClientID,
		"client_secret": a.ClientSecret,
		"server_url":    a.AuthorizeURL,
	}

	return values, nil
}

// DefaultHelmValues returns the default values for the chart
func DefaultHelmValues() map[string]interface{} {
	features := []string{
		"application_services=true",
		"services=true",
	}

	return map[string]interface{}{
		"api": map[string]interface{}{
			"feature_gates": features,
			"hostPort":      10080,
			"replicas":      1,
			"serviceType":   "NodePort",
		},
		"ui": map[string]interface{}{
			"feature_gates": features,
			"hostPort":      3000,
			"replicas":      1,
			"serviceType":   "NodePort",
		},
	}
}
