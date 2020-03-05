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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	gkeCredPath           = path.Join(localManifests, "gke-credentials.yml")
	gkeCredAllocationPath = path.Join(localManifests, "gke-allocation.yml")
	cachedAccountKeyPath  = path.Join(localManifests, "service-account-key.json")
)

func createLocalConfig(config *Config) {
	if config.CurrentProfile != "local" {
		config.SetCurrentProfile("local")
		config.Profiles = map[string]*Profile{
			"local": {
				Server:   "local",
				AuthInfo: "local",
			},
		}
	}

	if config.GetCurrentServer().Endpoint != localEndpoint {
		config.Servers = map[string]*Server{
			"local": {Endpoint: localEndpoint},
		}
		config.AuthInfos = map[string]*AuthInfo{
			"local": nil,
		}
	}
}

type authInfoConfig struct {
	ClientID, ClientSecret string
	AuthorizeURL           string
}

func newAuthInfoConfig(config *Config) *authInfoConfig {
	result := &authInfoConfig{}

	if config.AuthInfos["local"] != nil {
		result.ClientID = config.GetCurrentAuthInfo().OIDC.ClientID
		result.ClientSecret = config.GetCurrentAuthInfo().OIDC.ClientSecret
		result.AuthorizeURL = config.GetCurrentAuthInfo().OIDC.AuthorizeURL
	}

	return result
}

func (a *authInfoConfig) createPrompts() prompts {
	return prompts{
		&prompt{id: "Client ID", errMsg: "%s cannot be blank", value: &a.ClientID},
		&prompt{id: "Client Secret", errMsg: "%s cannot be blank", value: &a.ClientSecret},
		&prompt{id: "OpenID endpoint", errMsg: "%s cannot be blank", value: &a.AuthorizeURL},
	}
}

func (a *authInfoConfig) update(config *Config) {
	config.AuthInfos = map[string]*AuthInfo{
		"local": {
			Token: nil,
			OIDC: &OIDC{
				ClientID:     a.ClientID,
				ClientSecret: a.ClientSecret,
				AuthorizeURL: a.AuthorizeURL,
			},
		},
	}
}

type gcpInfoConfig struct {
	Region  string
	Project string
	KeyPath string
}

func newGcpInfoConfig() *gcpInfoConfig {
	return &gcpInfoConfig{}
}

func (g *gcpInfoConfig) load(src string) (*gcpInfoConfig, error) {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return g, nil
	}

	content, err := ioutil.ReadFile(src)
	if err != nil {
		return g, err
	}

	var cred gke.GKECredentials
	if err := yaml.NewDecoder(bytes.NewReader(content)).Decode(&cred); err != nil {
		return g, err
	}

	g.Region = cred.Spec.Region
	g.Project = cred.Spec.Project

	if len(cred.Spec.Account) > 0 {
		if _, err := os.Stat(cachedAccountKeyPath); err == nil {
			g.KeyPath = cachedAccountKeyPath
		}
	}

	return g, nil
}

func (g *gcpInfoConfig) createPrompts() prompts {
	return prompts{
		&prompt{id: "GKE Region", labelSuffix: "(e.g. europe-west2)", errMsg: "%s cannot be blank", value: &g.Region},
		&prompt{id: "GKE Project ID", errMsg: "%s cannot be blank", value: &g.Project},
		&prompt{
			id:          "GKE Service Account Key file",
			labelSuffix: "full path to service key file (will use cached if any)",
			errMsg:      "%s cannot be blank",
			value:       &g.KeyPath,
		},
	}
}

func (g *gcpInfoConfig) generateGcpInfo() error {
	// @step: ensure the directory
	if err := os.MkdirAll(localManifests, os.FileMode(0750)); err != nil {
		return err
	}

	// @step: read in the gke credentials json
	keyData, err := ioutil.ReadFile(filepath.Clean(g.KeyPath))
	if err != nil {
		return err
	}

	name := "gke"

	// @step: we need to render the credentials and the allocations for all teams
	creds := &gke.GKECredentials{
		TypeMeta: v1.TypeMeta{
			Kind:       "GKECredentials",
			APIVersion: gke.GroupVersion.String(),
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: kore.HubAdminTeam,
		},
		Spec: gke.GKECredentialsSpec{
			Region:  g.Region,
			Project: g.Project,
			Account: string(keyData),
		},
	}

	allocation := &configv1.Allocation{
		TypeMeta: v1.TypeMeta{
			Kind:       "Allocation",
			APIVersion: configv1.GroupVersion.String(),
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: kore.HubAdminTeam,
		},
		Spec: configv1.AllocationSpec{
			Name:    name,
			Summary: "Default Credentials for building a GKE Cluster",
			Resource: corev1.Ownership{
				Group:     gke.SchemeGroupVersion.Group,
				Version:   gke.SchemeGroupVersion.Version,
				Kind:      "GKECredentials",
				Namespace: kore.HubAdminTeam,
				Name:      name,
			},
			Teams: []string{"*"},
		},
	}

	// @step: write the files to the local manifests directory
	manifests := map[string]runtime.Object{
		gkeCredPath:           creds,
		gkeCredAllocationPath: allocation,
	}

	for path, object := range manifests {
		doc, err := utils.EncodeRuntimeObjectToYAML(object)
		if err != nil {

		}
		if err := ioutil.WriteFile(path, doc, os.FileMode(0640)); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(cachedAccountKeyPath, keyData, os.FileMode(0640))
}

func GetLocalConfigureSubCommand(config *Config) cli.Command {
	return cli.Command{
		Name:  "configure",
		Usage: "Configures a profile to connect to a local Kore installation.",
		Action: func(c *cli.Context) error {
			createLocalConfig(config)

			fmt.Println("What are your Identity Broker details?")
			authInfo := newAuthInfoConfig(config)
			if err := authInfo.createPrompts().collect(); err != nil {
				return err
			}
			authInfo.update(config)

			fmt.Println("What are your Google Cloud Platform details?")
			info, err := newGcpInfoConfig().load(gkeCredPath)
			if err != nil {
				return err
			}
			if err := info.createPrompts().collect(); err != nil {
				return err
			}
			if err := info.generateGcpInfo(); err != nil {
				return err
			}

			if err := config.Update(); err != nil {
				return err
			}

			fmt.Println("...Kore is now set up to run locally,")
			fmt.Println("✅ A 'local' profile has been configured in ~/.korectl/config")
			fmt.Println("✅ Generated Kubernetes CRDs are now stored in <project root>/manifests/local directory.")
			return nil
		},
	}
}
