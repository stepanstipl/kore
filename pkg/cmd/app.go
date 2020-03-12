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

	orderedArgs, err := orderArgs(appWrapper{app: a.app}, args)
	if err != nil {
		return err
	}

	return a.app.Run(orderedArgs)
}

func (a *App) helpPrinterCustom(out io.Writer, templ string, data interface{}, customFuncs map[string]interface{}) {
	a.origHelpPrinterCustom(out, templ, data, customFuncs)
	if data != a.app {
		a.origHelpPrinterCustom(a.app.Writer, globalOptionsTemplate, a.app, nil)
	}
}

func orderArgs(target wrapper, args []string) ([]string, error) {
	head := []string{args[0]}
	var tail []string

	for i := 1; i < len(args); i++ {
		arg := args[i]

		if arg == "--" {
			tail = append(tail, args[i:]...)
			break
		}

		if isFlag := strings.HasPrefix(arg, "-"); isFlag {
			flagName := strings.TrimLeft(arg, "-")
			res, err := parseFlagFromArgs(target, args[i:], flagName)
			if err != nil {
				return nil, err
			}
			if len(res) > 0 {
				head = append(head, res...)
				i += len(res) - 1
			} else {
				tail = append(tail, arg)
			}
		} else {
			tail = append(tail, arg)
		}
	}

	if len(tail) > 0 {
		if cmd := target.Command(tail[0]); cmd != nil {
			var err error
			tail, err = orderArgs(commandWrapper{cmd: cmd}, tail)
			if err != nil {
				return nil, err
			}

		}
	}

	return append(head, tail...), nil
}

func parseFlagFromArgs(target wrapper, args []string, name string) ([]string, error) {
	for i := 0; i < len(target.Flags()); i++ {
		flag := target.Flags()[i]
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

type wrapper interface {
	Command(name string) *cli.Command
	Flags() []cli.Flag
}

type appWrapper struct {
	app *cli.App
}

func (a appWrapper) Command(name string) *cli.Command {
	return a.app.Command(name)
}

func (a appWrapper) Flags() []cli.Flag {
	return a.app.Flags
}

type commandWrapper struct {
	cmd *cli.Command
}

func (c commandWrapper) Command(name string) *cli.Command {
	for _, cmd := range c.cmd.Subcommands {
		for _, cmdName := range cmd.Names() {
			if cmdName == name {
				return cmd
			}
		}
	}
	return nil
}

func (c commandWrapper) Flags() []cli.Flag {
	return c.cmd.Flags
}
