package cmd

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

type App struct {
	app *cli.App
}

func NewApp(app *cli.App) *App {
	return &App{app: app}
}

func (a *App) Run(args []string) error {
	orderedArgs, err := a.orderArgs(args)
	if err != nil {
		return err
	}

	return a.app.Run(orderedArgs)
}

func (a *App) orderArgs(args []string) ([]string, error) {
	flagArgs := []string{args[0]}
	var valueArgs []string

	for i := 1; i < len(args); i++ {
		arg := args[i]

		if arg == "--" {
			if len(args) > i+1 {
				valueArgs = append(valueArgs, args[i+1:]...)
			}
			break
		}

		if isFlag := strings.HasPrefix(arg, "-"); isFlag {
			flagName := strings.TrimLeft(arg, "-")
			res, err := a.parseFlagFromArgs(args[i:], flagName)
			if err != nil {
				return nil, err
			}
			if len(res) > 0 {
				flagArgs = append(flagArgs, res...)
				i += len(res) - 1
			} else {
				valueArgs = append(valueArgs, arg)
			}
		} else {
			valueArgs = append(valueArgs, arg)
		}
	}

	return append(flagArgs, valueArgs...), nil
}

func (a *App) parseFlagFromArgs(args []string, name string) ([]string, error) {
	for i := 0; i < len(a.app.Flags); i++ {
		flag := a.app.Flags[i]
		for _, flagName := range flag.Names() {
			if name == flagName {
				if f, ok := flag.(cli.DocGenerationFlag); ok {
					if f.TakesValue() {
						if len(args) == 1 {
							return nil, fmt.Errorf("%q parameter expects a value", flagName)
						}
						return []string{args[0], args[1]}, nil
					} else {
						return []string{args[0]}, nil
					}
				} else {
					panic(fmt.Errorf("%T global flag type is not supported yet, please add it to cli.App", flag))
				}
			}
		}
	}
	return nil, nil
}
