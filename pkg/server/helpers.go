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

package server

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// makeKubernetesConfig returns a rest.Config from the options
func makeKubernetesConfig(config KubernetesAPI) (*rest.Config, error) {
	// @step: are we creating an in-cluster kubernetes client
	if config.InCluster {
		return rest.InClusterConfig()
	}

	if config.KubeConfig != "" {
		return clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	}

	return &rest.Config{
		Host:        config.MasterAPIURL,
		BearerToken: config.Token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: config.SkipTLSVerify,
		},
	}, nil
}
