package korectl

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
)

const localEndpoint string = "http://127.0.0.1:10080"

func createLocalConfig(config *Config) error {
	config.CurrentProfile = "local"

	config.Profiles = map[string]*Profile{
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
	return config.Update()
}

type authInfo struct {
	clientId       string
	clientSecret   string
	openIdEndpoint string
}

func (a *authInfo) getClientId() error {
	id, err := a.prompt("ClientID")
	if err != nil {
		return err
	}
	a.clientId = id
	return nil
}

func (a *authInfo) getClientSecret() error {
	secret, err := a.prompt("Client Secret")
	if err != nil {
		return err
	}
	a.clientSecret = secret
	return nil
}

func (a *authInfo) getOpenIdEndpoint() error {
	endpoint, err := a.prompt("OpenID endpoint")
	if err != nil {
		return err
	}
	a.openIdEndpoint = endpoint
	return nil
}

func (a *authInfo) prompt(target string) (string, error) {
	prompt := promptui.Prompt{
		Label: fmt.Sprintf("%s (for your Identity Broker)", target),
		Validate: func(in string) error {
			if len(in) == 0 {
				return fmt.Errorf("%s cannot be blank", target)
			}
			return nil
		},
	}
	return prompt.Run()
}

func (a *authInfo) collect() error {
	if err := a.getClientId(); err != nil {
		return err
	}

	if err := a.getClientSecret(); err != nil {
		return err
	}

	if err := a.getOpenIdEndpoint(); err != nil {
		return err
	}

	return nil
}

func (a *authInfo) update(config *Config) error {
	config.AuthInfos = map[string]*AuthInfo{
		"local": {
			Token: nil,
			OIDC: &OIDC{
				ClientID:     a.clientId,
				ClientSecret: a.clientSecret,
				AuthorizeURL: a.openIdEndpoint,
			},
		},
	}

	return config.Update()
}

func GetLocalCommand(config *Config) cli.Command {
	return cli.Command{
		Name:  "local",
		Usage: "Used to configure and run a local instance of Kore.",
		Action: func(c *cli.Context) error {
			fmt.Println("Let's setup Kore to run locally.")
			fmt.Println("First, we need your Identity Broker details...")
			config, err := GetOrCreateClientConfiguration()
			if err != nil {
				return err
			}

			if err := createLocalConfig(config); err != nil {
				return err
			}

			authInfo := &authInfo{}
			if err := authInfo.collect(); err != nil {
				return err
			}
			if err := authInfo.update(config); err != nil {
				return err
			}

			return nil
		},
	}
}
