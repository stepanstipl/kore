package korectl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/ghodss/yaml"
	"github.com/urfave/cli"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const localEndpoint string = "http://127.0.0.1:10080"
const localInfraDir string = "./_kore_infra"

func updateAuthInfo(config *Config, clientId, clientSecret, openIdEndpoint string) error {
	config.AuthInfos = map[string]*AuthInfo{
		"local": {
			Token: nil,
			OIDC: &OIDC{
				ClientID:     clientId,
				ClientSecret: clientSecret,
				AuthorizeURL: openIdEndpoint,
			},
		},
	}

	return nil
}

func generateGcpInfo(region, projectId, keyPath string) error {
	keyData, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	cred := gke.GKECredentials{
		TypeMeta: v1.TypeMeta{
			Kind:       "GKECredentials",
			APIVersion: "gke.compute.kore.appvia.io/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:              "gke",
			CreationTimestamp: v1.NewTime(time.Now().UTC()),
		},
		Spec: gke.GKECredentialsSpec{
			Region:  region,
			Project: projectId,
			Account: string(keyData),
		},
		Status: gke.GKECredentialsStatus{
			Status:   corev1.SuccessStatus,
			Verified: true,
		},
	}

	data, err := yaml.Marshal(cred)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(localInfraDir, os.FileMode(0750)); err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(localInfraDir, "gke-credentials.yml"), data, os.FileMode(0640))
}

func createLocalConfig(config *Config) {
	config.CurrentContext = "local"

	config.Contexts = map[string]*Context{
		"local": {
			Server:   "local",
			AuthInfo: "local",
		},
	}

	config.Servers = map[string]*Server{
		"local": {Endpoint: localEndpoint},
	}

	config.AuthInfos = map[string]*AuthInfo{
		"local": {},
	}
}

func collectAndUpdateAuthInfo(config *Config) error {
	prompts := &prompts{prompts: []*prompt{
		{id: "Client ID", errMsg: "%s cannot be blank"},
		{id: "Client Secret", errMsg: "%s cannot be blank"},
		{id: "OpenID endpoint", errMsg: "%s cannot be blank"},
	}}

	if err := prompts.collect(); err != nil {
		return err
	}

	if err := updateAuthInfo(config,
		prompts.getValue("Client ID"),
		prompts.getValue("Client Secret"),
		prompts.getValue("OpenID endpoint"),
	); err != nil {
		return err
	}

	return nil
}

func collectAndGenerateGcpInfo() error {
	prompts := &prompts{prompts: []*prompt{
		{id: "GKE Region", labelSuffix: "(e.g. europe-west2)", errMsg: "%s cannot be blank"},
		{id: "GKE Project ID", errMsg: "%s cannot be blank"},
		{
			id:          "GKE Service Account Key file",
			labelSuffix: "(full path to the file)",
			errMsg:      "%s cannot be blank",
		},
	}}

	if err := prompts.collect(); err != nil {
		return err
	}

	if err := generateGcpInfo(
		prompts.getValue("GKE Region"),
		prompts.getValue("GKE Project ID"),
		prompts.getValue("GKE Service Account Key file"),
	); err != nil {
		return err
	}

	return nil
}

func GetLocalCommand(config *Config) cli.Command {
	return cli.Command{
		Name:  "local",
		Usage: "Used to configure and run a local instance of Kore.",
		Action: func(c *cli.Context) error {
			fmt.Println("Let's setup Kore to run locally:")
			createLocalConfig(config)

			fmt.Println("What are your Identity Broker details?")
			if err := collectAndUpdateAuthInfo(config); err != nil {
				return err
			}

			fmt.Println("What are your Google Cloud Platform details?")
			if err := collectAndGenerateGcpInfo(); err != nil {
				return err
			}

			return config.Update()
		},
	}
}
