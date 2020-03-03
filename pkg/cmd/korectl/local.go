package korectl

import (
	"github.com/urfave/cli"
)

const localEndpoint string = "http://127.0.0.1:10080"
const localManifests string = "./manifests/local"
const localCompose string = "./hack/compose"

func GetLocalCommand(config *Config) cli.Command {
	cmd := cli.Command{
		Name:  "local",
		Usage: "Used to configure and run a local instance of Kore.",
	}
	cmd.Subcommands = append(cmd.Subcommands, GetLocalConfigureSubCommand(config))
	cmd.Subcommands = append(cmd.Subcommands, GetLocalRunSubCommands(config)...)
	return cmd
}
