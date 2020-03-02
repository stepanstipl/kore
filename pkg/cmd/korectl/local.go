package korectl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const localEndpoint string = "http://127.0.0.1:10080"
const localManifests string = "./manifests/local"

var (
	gkeCredPath          = path.Join(localManifests, "gke-credentials.yml")
	cachedAccountKeyPath = path.Join(localManifests, "service-account-key.json")
)

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

	if err := os.MkdirAll(localManifests, os.FileMode(0750)); err != nil {
		return err
	}

	if err := ioutil.WriteFile(gkeCredPath, data, os.FileMode(0640)); err != nil {
		return err
	}

	return ioutil.WriteFile(cachedAccountKeyPath, keyData, os.FileMode(0640))
}

func createLocalConfig(config *Config) {
	if config.CurrentProfile != "local" {
		config.CurrentProfile = "local"
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

func collectAndUpdateAuthInfo(config *Config) error {
	var defaultClientId, defaultClientSecret, defaultAuthorizeURL string
	if config.AuthInfos["local"] != nil {
		defaultClientId = config.GetCurrentAuthInfo().OIDC.ClientID
		defaultClientSecret = config.GetCurrentAuthInfo().OIDC.ClientSecret
		defaultAuthorizeURL = config.GetCurrentAuthInfo().OIDC.AuthorizeURL
	}

	prompts := prompts{
		&prompt{id: "Client ID", errMsg: "%s cannot be blank", value: defaultClientId},
		&prompt{id: "Client Secret", errMsg: "%s cannot be blank", value: defaultClientSecret},
		&prompt{id: "OpenID endpoint", errMsg: "%s cannot be blank", value: defaultAuthorizeURL},
	}

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

func getCurrentGCPInfo() (gke.GKECredentialsSpec, error) {
	content, err := ioutil.ReadFile(gkeCredPath)
	if err != nil {
		return gke.GKECredentialsSpec{}, err
	}

	var cred gke.GKECredentials
	if err := yaml.NewDecoder(bytes.NewReader(content)).Decode(&cred); err != nil {
		return gke.GKECredentialsSpec{}, err
	}

	return cred.Spec, nil
}

func createAccountKeyPrompt(spec gke.GKECredentialsSpec) *prompt {
	var defaultVal, labelSuffix = "", "(full path to new file)"

	if len(spec.Account) > 0 {
		if _, err := os.Stat(cachedAccountKeyPath); err == nil {
			defaultVal = cachedAccountKeyPath
			labelSuffix = "(existing service key data from cached file)"
		}
	}

	return &prompt{
		id:          "GKE Service Account Key file",
		labelSuffix: labelSuffix,
		errMsg:      "%s cannot be blank",
		value:       defaultVal,
	}
}

func collectAndGenerateGcpInfo() error {
	current, _ := getCurrentGCPInfo()

	prompts := prompts{
		&prompt{id: "GKE Region", labelSuffix: "(e.g. europe-west2)", errMsg: "%s cannot be blank", value: current.Region},
		&prompt{id: "GKE Project ID", errMsg: "%s cannot be blank", value: current.Project},
		createAccountKeyPrompt(current),
	}

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
		Subcommands: []cli.Command{
			{
				Name:  "configure",
				Usage: "Used to configure a local instance of Kore.",
				Action: func(c *cli.Context) error {
					createLocalConfig(config)

					fmt.Println("What are your Identity Broker details?")
					if err := collectAndUpdateAuthInfo(config); err != nil {
						return err
					}

					fmt.Println("What are your Google Cloud Platform details?")
					if err := collectAndGenerateGcpInfo(); err != nil {
						return err
					}

					if err := config.Update(); err != nil {
						return err
					}

					fmt.Println("...Kore is now set up to run locally,")
					fmt.Println("✅ A 'local' profile has been configured in ~/.korectl/config")
					fmt.Println("✅ Generated Kubernetes CRDs are now stored in <project root>/manifests/local directory. ")
					return nil
				},
			},
		},
	}
}
