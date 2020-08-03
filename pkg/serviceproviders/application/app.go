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
	"fmt"
	"strconv"

	"github.com/appvia/kore/pkg/utils/kubernetes"

	"k8s.io/apimachinery/pkg/runtime"
	applicationv1beta "sigs.k8s.io/application/api/v1beta1"
)

type AppConfiguration struct {
	// Resources is the list of Kubernetes resources to deploy
	Resources Resources `json:"resources"`
	// Values are parameters for the resource templates, which can be referenced as {{ .Values.foo }}
	Values YAMLMap `json:"values,omitempty"`
	// Secrets are parameters provided using ConfigurationFrom for the resource templates, which can be referenced as {{ .Secrets.foo }}
	Secrets map[string]interface{} `json:"secrets,omitempty"`
}

func (c *AppConfiguration) CompileResources(params ResourceParams) (Resources, error) {
	var compiledResources Resources
	for _, r := range c.Resources {
		compiled, err := compileResource(r.DeepCopyObject(), params)
		if err != nil {
			return nil, fmt.Errorf("compiling resource %v failed: %w", r, err)
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
	objects := kubernetes.Objects(r)
	manifest, err := objects.MarshalYAML()
	if err != nil {
		return nil, err
	}

	return []byte(strconv.Quote(string(manifest))), nil
}

func (r *Resources) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	objects := &kubernetes.Objects{}

	if err := objects.UnmarshalYAML([]byte(unquoted)); err != nil {
		return err
	}

	*r = Resources(*objects)
	return nil
}
