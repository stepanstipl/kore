/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of hub-apiserver.
 *
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package hubctl

import (
	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// Populate kubeconfig is used to populate the users kubeconfig
func PopulateKubeconfig(clusters *clustersv1.KubernetesList, kubeconfig string, config *Config) error {
	cfg, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return err
	}

	cfg.AuthInfos["kore"] = &api.AuthInfo{
		AuthProvider: &api.AuthProviderConfig{
			Name: "oidc",
			Config: map[string]string{
				"access-token":   config.Credentials.AccessToken,
				"client-id":      config.Credentials.ClientID,
				"client-secret":  config.Credentials.ClientSecret,
				"id-token":       config.Credentials.IDToken,
				"idp-issuer-url": config.AuthorizeURL,
				"refresh-token":  config.Credentials.RefreshToken,
			},
		},
	}

	for _, x := range clusters.Items {
		if x.Status.Endpoint == "" {
			log.Warnf("skipping cluster: %s as it does not have an endpoint yet", x.Name)
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
