package korectl

import (
	"fmt"

	"github.com/urfave/cli"
)

const localEndpoint string = "http://127.0.0.1:10080"

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

func GetLocalCommand(config *Config) cli.Command {
	return cli.Command{
		Name:  "local",
		Usage: "Used to configure and run a local instance of Kore.",
		Action: func(c *cli.Context) error {
			fmt.Println("Hello from local.")
			config, err := GetOrCreateClientConfiguration()
			if err != nil {
				return err
			}

			return createLocalConfig(config)
		},
	}
}
