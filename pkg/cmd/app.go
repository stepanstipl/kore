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

package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/appvia/kore/pkg/utils"

	"github.com/urfave/cli/v2"
)

func init() {
	cli.SubcommandHelpTemplate = `NAME:
   {{.HelpName}} - {{.Usage}}

USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}

DESCRIPTION:
   {{if .Description}}{{.Description}}{{else}}{{.Usage}}{{end}}{{if len .VisibleCategories}}

COMMANDS:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{else}}{{range .VisibleCommands}}
   {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`

	cli.CommandHelpTemplate = `NAME:
   {{.HelpName}} - {{.Usage}}

USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Category}}

CATEGORY:
   {{.Category}}{{end}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if .VisibleFlags}}

OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{else}}
{{end}}
`
}

var globalOptionsTemplate = `{{if .VisibleFlags}}GLOBAL OPTIONS:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}
{{end}}
`

type App struct {
	app                   *cli.App
	origHelpPrinterCustom func(io.Writer, string, interface{}, map[string]interface{})
}

func NewApp(app *cli.App) *App {
	return &App{
		app: app,
	}
}

func (a *App) Run(args []string) error {
	a.origHelpPrinterCustom = cli.HelpPrinterCustom
	cli.HelpPrinterCustom = a.helpPrinterCustom
	defer func() {
		cli.HelpPrinterCustom = a.origHelpPrinterCustom
	}()

	orderedArgs, err := a.Reorder(args)
	if err != nil {
		return err
	}

	return a.app.Run(orderedArgs)
}

// Reorder is responsible for ordering the argus
func (a *App) Reorder(args []string) ([]string, error) {
	var parent *cli.Command
	ordered := []string{}

	// @step: create three arrays used to hold the global flags
	// and the head and tail of a command
	var head, tail, global []string

	var found bool

	// @step: we iterate the arguments in the os.Args
	for i := 1; i < len(args); i++ {
		arg := args[i]

		// @step: if this a break
		if arg == "--" {
			tail = append(tail, args[i:]...)

			break
		}

		// @step: is this a flag argument?
		if strings.HasPrefix(arg, "-") {
			parsed := strings.TrimLeft(arg, "-")

			// @step: we check if the flag is global for related to the command
			flag, found := a.getGlobalFlag(parsed)
			if found {
				// we are dealing with a command flag - we need to know if it has args though
				values, err := a.getFlagValues(parsed, args[i:], flag)
				if err != nil {
					return nil, err
				}
				// add the arguments to the global
				global = append(global, values...)
				// remove the values from the index
				i += len(values) - 1
			} else {
				// we have a command flag, lets see if we can find it
				if parent == nil {
					tail = append(tail, arg)
				} else {
					found, err := func() (bool, error) {
						for _, flag := range parent.Flags {
							if utils.Contains(parsed, flag.Names()) {
								values, err := a.getFlagValues(parsed, args[i:], flag)
								if err != nil {
									return false, err
								}
								head = append(head, values...)
								i += len(values) - 1

								return true, nil
							}
						}

						return false, nil
					}()
					if err != nil {
						return nil, err
					}
					if !found {
						tail = append(tail, arg)
					}
				}
			}
		} else {
			// @step: the argument is not a flag, is it a command?
			if parent == nil {
				parent, found = a.getAppCommand(arg)
				if !found {
					tail = append(tail, arg)
				} else {
					head = append(head, arg)
				}
			} else {
				command, found := a.getSubCommand(arg, parent)
				if found {
					parent = command
					head = append(head, arg)
					ordered = append(ordered, head...)
					ordered = append(ordered, tail...)
					head = []string{}
					tail = []string{}
				} else {
					tail = append(tail, arg)
				}
			}
		}
	}

	ordered = append(ordered, head...)
	ordered = append(ordered, tail...)
	full := append(global, ordered...)

	return append([]string{args[0]}, full...), nil
}

func (a *App) getAppCommand(name string) (*cli.Command, bool) {
	for _, x := range a.app.Commands {
		if utils.Contains(name, x.Names()) {
			return x, true
		}
	}

	return nil, false
}

func (a *App) getSubCommand(name string, parent *cli.Command) (*cli.Command, bool) {
	for _, x := range parent.Subcommands {
		if utils.Contains(name, x.Names()) {
			return x, true
		}
	}

	return nil, false
}

func (a *App) getGlobalFlag(name string) (cli.Flag, bool) {
	for _, flag := range a.app.Flags {
		if utils.Contains(name, flag.Names()) {
			return flag, true
		}
	}

	return nil, false
}

// getFlagValues returns the flag and if it requires a parameter, the param
func (a *App) getFlagValues(name string, args []string, flag cli.Flag) ([]string, error) {
	// @step: we check if the flag requires a value
	if f, ok := flag.(cli.DocGenerationFlag); ok {
		if f.TakesValue() {
			if len(args) == 1 {
				return nil, fmt.Errorf("%q parameter expects a value", name)
			}

			return []string{args[0], args[1]}, nil
		}

		return []string{args[0]}, nil
	}

	panic(fmt.Errorf("%T global flag type is not supported yet, please add it to cli.App", flag))
}

func (a *App) helpPrinterCustom(out io.Writer, templ string, data interface{}, customFuncs map[string]interface{}) {
	a.origHelpPrinterCustom(out, templ, data, customFuncs)
	if data != a.app {
		a.origHelpPrinterCustom(a.app.Writer, globalOptionsTemplate, a.app, nil)
	}
}
