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

package kore

import (
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

// NewCmdCompletion creates and returns the shell completion command
func NewCmdCompletion(factory cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "completion",
		DisableFlagsInUseLine: true,
		Short:                 "Provides the autocomplete output so you can source into your shell",
		Example:               "kore completion <shell>",
		Run:                   cmdutil.RunHelp,
	}

	cmd.AddCommand(NewCmdCompletionsBash(factory))
	cmd.AddCommand(NewCmdCompletionsZsh(factory))

	return cmd
}

// NewCmdCompletionsBash returns the bash completion command
func NewCmdCompletionsBash(factory cmdutil.Factory) *cobra.Command {
	command := &cobra.Command{
		Use:     "bash",
		Short:   "generate the bash command auto-completion code",
		Example: "source <(kore completion bash)",

		Run: func(cmd *cobra.Command, args []string) {
			if err := root.GenBashCompletion(factory.Writer()); err != nil {
				panic(err)
			}
		},
	}

	return command
}

// NewCmdCompletionsZsh returns the zsh completion comman
func NewCmdCompletionsZsh(factory cmdutil.Factory) *cobra.Command {
	command := &cobra.Command{
		Use:     "zsh",
		Short:   "generate the zsh command auto-completion code",
		Example: "source <(kore completion zsh)",

		Run: func(cmd *cobra.Command, args []string) {
			if err := root.GenZshCompletion(factory.Writer()); err != nil {
				panic(err)
			}
		},
	}

	return command
}
