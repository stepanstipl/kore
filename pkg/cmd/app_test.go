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
