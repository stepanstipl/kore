package korectl

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
)

const localEndpoint string = "http://127.0.0.1:10080"

type prompt struct {
	id     string
	label  string
	errMsg string
	value  string
}

func (p *prompt) do() error {
	runner := promptui.Prompt{
		Label: p.id + " " + p.label,
		//Label: "A label",
		Validate: func(in string) error {
			if len(in) == 0 {
				return fmt.Errorf(p.errMsg, p.id)
			}
			return nil
		},
	}

	gathered, err := runner.Run()
	if err != nil {
		return err
	}

	p.value = gathered
	return nil
}

type prompts struct {
	prompts []*prompt
}

func (p *prompts) collect() error {
	for _, p := range p.prompts {
		if err := p.do(); err != nil {
			return err
		}
	}
	return nil
}

func (p *prompts) getValue(id string) string {
	for _, p := range p.prompts {
		if p.id == id {
			return p.value
		}
	}
	return ""
}

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

	return config.Update()
}

func createLocalConfig(config *Config) error {
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
	return config.Update()
}

func collectAndUpdateAuthInfo(config *Config) error {
	prompts := &prompts{prompts: []*prompt{
		&prompt{
			id:     "Client ID",
			label:  "%s (for your Identity Broker)",
			errMsg: "%s cannot be blank",
		},
		&prompt{
			id:     "Client Secret",
			label:  "%s (for your Identity Broker)",
			errMsg: "%s cannot be blank",
		},
		&prompt{
			id:     "OpenID endpoint",
			label:  "%s (for your Identity Broker)",
			errMsg: "%s cannot be blank",
		},
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

func GetLocalCommand(config *Config) cli.Command {
	return cli.Command{
		Name:  "local",
		Usage: "Used to configure and run a local instance of Kore.",
		Action: func(c *cli.Context) error {
			fmt.Printf("config: [%+v]", config)
			fmt.Println("Let's setup Kore to run locally.")

			if err := createLocalConfig(config); err != nil {
				return err
			}

			fmt.Println("We'll start by gathering your Identity Broker details...")
			if err := collectAndUpdateAuthInfo(config); err != nil {
				return err
			}

			return nil
		},
	}
}
