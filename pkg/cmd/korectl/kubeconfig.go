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
	"os"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/urfave/cli/v2"
)

var kubeconfigUser = "kore"

func GetKubeconfigCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:    "kubeconfig",
		Aliases: []string{"kconfig"},
		Usage:   "Adds your team provisioned clusters to your kubeconfig",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "The name of the integration to retrieve `NAME`",
			},
		},
		Action: func(ctx *cli.Context) error {
			return DoKubeconfig(config, ctx.String("team"))
		},
	}
}

func DoKubeconfig(config *Config, team string) error {
	list := &clustersv1.KubernetesList{}

	if err := GetTeamResourceList(config, team, "kubernetes", list); err != nil {
		return err
	}

	if len(list.Items) <= 0 {
		fmt.Println("No clusters found in this team's namespace")
		return nil
	}

	kubeconfig, err := GetOrCreateKubeConfig()
	if err != nil {
		return err
	}

	if err := WriteKubeconfig(list, kubeconfig, config); err != nil {
		return err
	}

	return newKubeconfigResultPrinter(team, list).Print(os.Stdout)
}

// WriteKubeconfig writes kubeconfig to the user's kubeconfig
func WriteKubeconfig(clusters *clustersv1.KubernetesList, kubeconfig string, config *Config) error {
	cfg, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return err
	}
	auth := config.GetCurrentAuthInfo()
	if auth.OIDC == nil {
		return errors.New("you must be using a context backed by an idp")
	}

	cfg.AuthInfos[kubeconfigUser] = &api.AuthInfo{
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
			fmt.Printf("SKIPPING CLUSTER: %s as it does not have an endpoint yet\n", x.Name)
			continue
		}
		// @step: add the endpoint
		cfg.Clusters[x.Name] = &api.Cluster{
			InsecureSkipTLSVerify: true,
			Server:                "https://" + x.Status.Endpoint,
		}

		// @step: add the context
		if _, found := cfg.Contexts[x.Name]; !found {
			cfg.Contexts[x.Name] = &api.Context{
				AuthInfo: kubeconfigUser,
				Cluster:  x.Name,
			}
		}
	}

	return clientcmd.WriteToFile(*cfg, kubeconfig)
}

func newKubeconfigResultPrinter(team string, clusters *clustersv1.KubernetesList) printer {
	var (
		header  string
		columns = []string{"Context", "Cluster"}
		lines   [][]string
	)

	for _, item := range clusters.Items {
		if len(item.Status.Endpoint) < 1 {
			continue
		}
		lines = append(lines, []string{item.Name, item.Name})
	}

	if len(lines) < 1 {
		header = fmt.Sprintf("Successfully added the [%s] user to your kubeconfig", kubeconfigUser)
	} else {
		header = fmt.Sprintf("Successfully added team [%s] provisioned clusters to your kubeconfig", team)
	}

	return table{
		header,
		columns,
		lines,
	}
}
