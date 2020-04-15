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
	"github.com/appvia/kore/pkg/version"

	"github.com/spf13/cobra"
)

// NewCmdVersion creates and returns the version command
func NewCmdVersion(factory cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "version",
		DisableFlagsInUseLine: true,
		Short:                 "Prints the version of the client",

		Run: func(cmd *cobra.Command, args []string) {
			factory.Println("Client: %s (git+sha: %s)", version.Release, version.GitSHA)
		},
	}

	return cmd
}