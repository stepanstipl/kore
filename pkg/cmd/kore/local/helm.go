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
	"io/ioutil"
	"os"

	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils"

	"sigs.k8s.io/yaml"
)

// UpdateHelmValues is used to update the inline values
func (o *UpOptions) UpdateHelmValues(path string) error {
	// @step: get the current content
	current, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	// make a copy of original
	values := make([]byte, len(current))
	copy(values, current)

	// @step: iterate through any changes
	updated, err := func() ([]byte, error) {
		if utils.Contains("release", o.FlagsChanged) {
			for _, x := range []string{"api.version", "ui.version"} {
				values, err = UpdateYAML(values, x, o.Release)
				if err != nil {
					return nil, err
				}
			}
		}

		return values, nil
	}()
	if err != nil {
		return err
	}

	if !bytes.Equal(current, updated) {
		return ioutil.WriteFile(path, updated, os.FileMode(0750))
	}

	return nil
}

// GetHelmValues returns returns or prompts for the values
func (o *UpOptions) GetHelmValues(path string) (map[string]interface{}, error) {
	found, err := utils.FileExists(path)
	if err != nil {
		return nil, err

	} else if !found {
		values := o.GetDefaultHelmValues()

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
	} else {
		if err := o.UpdateHelmValues(path); err != nil {
			return nil, err
		}
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
func (o *UpOptions) GetDefaultHelmValues() map[string]interface{} {
	return map[string]interface{}{
		"api": map[string]interface{}{
			"feature_gates": []string{"services=true"},
			"hostPort":      10080,
			"replicas":      1,
			"serviceType":   "NodePort",
			"version":       o.Version,
		},
		"ui": map[string]interface{}{
			"feature_gates": []string{"services=true"},
			"hostPort":      3000,
			"replicas":      1,
			"serviceType":   "NodePort",
			"version":       o.Version,
		},
	}
}
