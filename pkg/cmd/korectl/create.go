package korectl

import (
	"github.com/urfave/cli"
)

func GetCreateCommand(config *Config) cli.Command {
	return cli.Command{
		Name:  "create",
		Usage: "creates various objects",

		Subcommands: []cli.Command{
			GetCreateTeamCommand(config),
		},
	}
}
