/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package korectl

import (
	"errors"

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
	auth := config.GetCurrentAuthInfo()
	if auth.OIDC == nil {
		return errors.New("you must be using a context backed by a idp")
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
