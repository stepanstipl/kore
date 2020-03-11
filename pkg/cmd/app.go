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
			valueArgs = append(valueArgs, args[i:]...)
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

func (a *App) helpPrinterCustom(out io.Writer, templ string, data interface{}, customFuncs map[string]interface{}) {
	a.origHelpPrinterCustom(out, templ, data, customFuncs)
	if data != a.app {
		a.origHelpPrinterCustom(a.app.Writer, globalOptionsTemplate, a.app, nil)
	}
}
