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

package application

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	koreschema "github.com/appvia/kore/pkg/schema"
	"k8s.io/apimachinery/pkg/runtime"
	applicationv1beta "sigs.k8s.io/application/api/v1beta1"
	"sigs.k8s.io/yaml"
)

type AppConfiguration struct {
	// Resources is the list of Kubernetes resources to deploy
	Resources Resources `json:"resources"`
	// Values are parameters for the resource templates, which can be referenced as {{ .Values.foo }}
	Values YAMLMap `json:"values"`
}

func (c *AppConfiguration) CompileResources(params ResourceParams) (Resources, error) {
	var compiledResources Resources
	for _, r := range c.Resources {
		compiled, err := compileResource(r.DeepCopyObject(), params)
		if err != nil {
			return nil, err
		}
		compiledResources = append(compiledResources, compiled)
	}
	return compiledResources, nil
}

type Resources []runtime.Object

func (r Resources) Application() *applicationv1beta.Application {
	for _, res := range r {
		if app, ok := res.(*applicationv1beta.Application); ok {
			return app
		}
	}
	return nil
}

func (r Resources) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 16384))
	for _, obj := range r {
		jsonData, err := json.Marshal(obj)
		if err != nil {
			return nil, err
		}
		yamlData, err := yaml.JSONToYAML(jsonData)
		if err != nil {
			return nil, err
		}
		buf.WriteString("---\n")
		buf.Write(yamlData)
		buf.WriteRune('\n')
	}
	return []byte(strconv.Quote(buf.String())), nil
}

func (r *Resources) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	raw, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	documents := regexp.MustCompile("(?m)^---\n").Split(raw, -1)

	var objects []runtime.Object

	for _, document := range documents {
		if strings.TrimSpace(document) == "" {
			continue
		}

		obj, err := koreschema.DecodeYAML([]byte(document))
		if err != nil {
			return err
		}

		objects = append(objects, obj)
	}

	*r = objects
	return nil
}
