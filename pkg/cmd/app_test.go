package cmd_test

import (
	"testing"

	"github.com/appvia/kore/pkg/cmd"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestAppShouldHandleNoArgs(t *testing.T) {
	var appArgs []string

	app := &cli.App{
		Name: "test",
		Action: func(ctx *cli.Context) error {
			appArgs = ctx.Args().Slice()
			return nil
		},
	}
	err := cmd.NewApp(app).Run([]string{"test"})
	require.NoError(t, err)

	require.Equal(t, []string{}, appArgs)
}

func TestAppShouldHandleMultipleArguments(t *testing.T) {
	var appArgs []string

	app := &cli.App{
		Name: "test",
		Action: func(ctx *cli.Context) error {
			appArgs = ctx.Args().Slice()
			return nil
		},
	}
	err := cmd.NewApp(app).Run([]string{"test", "arg1", "arg2"})
	require.NoError(t, err)

	require.Equal(t, []string{"arg1", "arg2"}, appArgs)
}

func TestAppShouldHandleGlobalStringFlagsBeforeArgs(t *testing.T) {
	var appArgs []string
	var globalFlag string

	app := &cli.App{
		Name: "test",
		Action: func(ctx *cli.Context) error {
			appArgs = ctx.Args().Slice()
			globalFlag = ctx.String("global-flag")
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "global-flag",
			},
		},
	}
	err := cmd.NewApp(app).Run([]string{"test", "--global-flag", "foo", "arg1", "arg2"})
	require.NoError(t, err)

	require.Equal(t, []string{"arg1", "arg2"}, appArgs)
	require.Equal(t, "foo", globalFlag)
}

func TestAppShouldHandleGlobalStringFlagsAfterArgs(t *testing.T) {
	var appArgs []string
	var globalFlag string

	app := &cli.App{
		Name: "test",
		Action: func(ctx *cli.Context) error {
			appArgs = ctx.Args().Slice()
			globalFlag = ctx.String("global-flag")
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "global-flag",
			},
		},
	}
	err := cmd.NewApp(app).Run([]string{"test", "arg1", "arg2", "--global-flag", "foo"})
	require.NoError(t, err)

	require.Equal(t, []string{"arg1", "arg2"}, appArgs)
	require.Equal(t, "foo", globalFlag)
}

func TestAppShouldReturnErrorIfNoGlobalStringValueIsPresent(t *testing.T) {
	app := &cli.App{
		Name: "test",
		Action: func(ctx *cli.Context) error {
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "global-flag",
			},
		},
	}
	err := cmd.NewApp(app).Run([]string{"test", "arg1", "--global-flag"})
	require.EqualError(t, err, "\"global-flag\" parameter expects a value")
}

func TestAppShouldHandleGlobalStringFlagsBetweenArgs(t *testing.T) {
	var appArgs []string
	var globalFlag string

	app := &cli.App{
		Name: "test",
		Action: func(ctx *cli.Context) error {
			appArgs = ctx.Args().Slice()
			globalFlag = ctx.String("global-flag")
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "global-flag",
			},
		},
	}
	err := cmd.NewApp(app).Run([]string{"test", "arg1", "--global-flag", "foo", "arg2"})
	require.NoError(t, err)

	require.Equal(t, []string{"arg1", "arg2"}, appArgs)
	require.Equal(t, "foo", globalFlag)
}

func TestAppShouldHandleGlobalBoolFlagsBeforeArgs(t *testing.T) {
	var appArgs []string
	var debug bool

	app := &cli.App{
		Name: "test",
		Action: func(ctx *cli.Context) error {
			appArgs = ctx.Args().Slice()
			debug = ctx.Bool("debug")
			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "debug",
			},
		},
	}
	err := cmd.NewApp(app).Run([]string{"test", "--debug", "arg1", "arg2"})
	require.NoError(t, err)

	require.Equal(t, []string{"arg1", "arg2"}, appArgs)
	require.True(t, debug)
}

func TestAppShouldHandleGlobalBoolFlagsBetweenArgs(t *testing.T) {
	var appArgs []string
	var debug bool

	app := &cli.App{
		Name: "test",
		Action: func(ctx *cli.Context) error {
			appArgs = ctx.Args().Slice()
			debug = ctx.Bool("debug")
			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "debug",
			},
		},
	}
	err := cmd.NewApp(app).Run([]string{"test", "arg1", "--debug", "arg2"})
	require.NoError(t, err)

	require.Equal(t, []string{"arg1", "arg2"}, appArgs)
	require.True(t, debug)
}

func TestAppShouldHandleGlobalBoolFlagsAfterArgs(t *testing.T) {
	var appArgs []string
	var debug bool

	app := &cli.App{
		Name: "test",
		Action: func(ctx *cli.Context) error {
			appArgs = ctx.Args().Slice()
			debug = ctx.Bool("debug")
			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "debug",
			},
		},
	}
	err := cmd.NewApp(app).Run([]string{"test", "arg1", "arg2", "--debug"})
	require.NoError(t, err)

	require.Equal(t, []string{"arg1", "arg2"}, appArgs)
	require.True(t, debug)
}

func TestAppShouldLeaveNonGlobalFlagsAsIs(t *testing.T) {
	var appArgs []string

	app := &cli.App{
		Name: "test",
		Action: func(ctx *cli.Context) error {
			appArgs = ctx.Args().Slice()
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "global-flag",
			},
		},
	}
	err := cmd.NewApp(app).Run([]string{"test", "arg1", "arg2", "--other-flag", "foo"})
	require.NoError(t, err)

	require.Equal(t, []string{"arg1", "arg2", "--other-flag", "foo"}, appArgs)
}

func TestAppShouldIgnoreFlagsAfterDoubleDash(t *testing.T) {
	var appArgs []string
	var globalFlag string

	app := &cli.App{
		Name: "test",
		Action: func(ctx *cli.Context) error {
			appArgs = ctx.Args().Slice()
			globalFlag = ctx.String("global-flag")
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "global-flag",
			},
		},
	}
	err := cmd.NewApp(app).Run([]string{"test", "arg1", "arg2", "--", "--global-flag", "foo"})
	require.NoError(t, err)

	require.Equal(t, []string{"arg1", "arg2", "--", "--global-flag", "foo"}, appArgs)
	require.Equal(t, "", globalFlag)
}

func TestAppShouldReorderSubCommandFlags(t *testing.T) {
	var appArgs []string
	var globalFlag string
	var cmdFlag string

	app := &cli.App{
		Name: "test",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "global-flag",
			},
		},
		Commands: []*cli.Command{
			{
				Name: "cmd",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "cmd-flag",
					},
				},
				Action: func(ctx *cli.Context) error {
					appArgs = ctx.Args().Slice()
					globalFlag = ctx.String("global-flag")
					cmdFlag = ctx.String("cmd-flag")
					return nil
				},
			},
		},
	}
	err := cmd.NewApp(app).Run([]string{"test", "cmd", "--global-flag", "foo", "arg1", "--cmd-flag", "bar"})
	require.NoError(t, err)

	require.Equal(t, []string{"arg1"}, appArgs)
	require.Equal(t, "foo", globalFlag)
	require.Equal(t, "bar", cmdFlag)
}
