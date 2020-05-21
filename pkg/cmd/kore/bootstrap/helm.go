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

package bootstrap

import (
	"io/ioutil"
	"os"

	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/version"

	"gopkg.in/yaml.v2"
)

// GetHelmValues returns returns or prompts for the values
func GetHelmValues(path string) (map[string]interface{}, error) {
	found, err := utils.FileExists(path)
	if err != nil {
		return nil, err
	}
	// @TODO we should probably check the params in the values file
	if !found {
		values := GetDefaultHelmValues()

		a := authInfoConfig{}

		if err := (&cmdutil.Prompts{
			&cmdutil.Prompt{Id: "Client ID", ErrMsg: "%s cannot be blank", Value: &a.ClientID},
			&cmdutil.Prompt{Id: "Client Secret", ErrMsg: "%s cannot be blank", Value: &a.ClientSecret},
			&cmdutil.Prompt{Id: "Authorization Endpoint", ErrMsg: "%s cannot be blank", Value: &a.AuthorizeURL},
		}).Collect(); err != nil {
			return nil, err
		}

		values["idp"] = map[string]interface{}{
			"client_id":     a.ClientID,
			"client_secret": a.ClientSecret,
			"server_url":    a.AuthorizeURL,
		}

		return values, nil
	}

	// @step: we read in the values.yml
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	values := make(map[string]interface{})
	if err := yaml.Unmarshal(content, &values); err != nil {
		return nil, err
	}

	return values, nil
}

// GetDefaultHelmValues returns the default values for the chart
func GetDefaultHelmValues() map[string]interface{} {
	return map[string]interface{}{
		"api": map[string]interface{}{
			"feature_gates": []string{"services=true"},
			"hostPort":      10080,
			"replicas":      1,
			"serviceType":   "NodePort",
			"version":       version.Release,
		},
		"ui": map[string]interface{}{
			"feature_gates": []string{"services=true"},
			"hostPort":      3000,
			"replicas":      1,
			"serviceType":   "NodePort",
			"version":       version.Release,
		},
	}
}
