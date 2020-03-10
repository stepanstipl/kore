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
	"errors"
	"fmt"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// Populate kubeconfig is used to populate the users kubeconfig
func PopulateKubeconfig(clusters *clustersv1.KubernetesList, kubeconfig string, config *Config) error {
	cfg, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return err
	}
	auth := config.GetCurrentAuthInfo()
	if auth.OIDC == nil {
		return errors.New("you must be using a context backed by an idp")
	}

	cfg.AuthInfos["kore"] = &api.AuthInfo{
		AuthProvider: &api.AuthProviderConfig{
			Name: "oidc",
			Config: map[string]string{
				"access-token":   auth.OIDC.AccessToken,
				"client-id":      auth.OIDC.ClientID,
				"client-secret":  auth.OIDC.ClientSecret,
				"id-token":       auth.OIDC.IDToken,
				"idp-issuer-url": auth.OIDC.AuthorizeURL,
				"refresh-token":  auth.OIDC.RefreshToken,
			},
		},
	}

	for _, x := range clusters.Items {
		if x.Status.Endpoint == "" {
			fmt.Printf("skipping cluster: %s as it does not have an endpoint yet\n", x.Name)
		}
		// @step: add the endpoint
		cfg.Clusters[x.Name] = &api.Cluster{
			InsecureSkipTLSVerify: true,
			Server:                "https://" + x.Status.Endpoint,
		}

		// @step: add the context
		if _, found := cfg.Contexts[x.Name]; !found {
			cfg.Contexts[x.Name] = &api.Context{
				AuthInfo: "kore",
				Cluster:  x.Name,
			}
		}
	}

	return clientcmd.WriteToFile(*cfg, kubeconfig)
}
