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

package kubernetes

import (
	"os"
	"path"

	"github.com/appvia/kore/pkg/utils"
)

// GetKubeConfigPath return the path for kubeconig
func GetKubeConfigPath() string {
	return func() string {
		p := os.ExpandEnv(os.Getenv("$KUBECONFIG"))
		if p != "" {
			return p
		}

		return os.ExpandEnv(path.Join("${HOME}", ".kube", "config"))
	}()
}

// GetOrCreateKubeConfig is used to retrieve the kubeconfig path
func GetOrCreateKubeConfig() (string, error) {
	_, err := utils.EnsureFileExists(GetKubeConfigPath())
	if err != nil {
		return "", err
	}

	return GetKubeConfigPath(), nil
}
