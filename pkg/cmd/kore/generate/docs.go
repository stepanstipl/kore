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

package generate

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	cmdutil "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// DocOptions are the command options
type DocOptions struct {
	cmdutil.Factory
	// Directory is the place to render the markdown
	Directory string
	// Root is the root command
	Root *cobra.Command
}

const (
	fmTemplate = `---
date: %s
title: "%s"
slug: %s
url: %s
---
`
)

// NewCmdDocs disables the given feature gate
func NewCmdDocs(factory cmdutil.Factory, root *cobra.Command) *cobra.Command {
	o := &DocOptions{Factory: factory, Root: root}

	cmd := &cobra.Command{
		Use:     "docs",
		Short:   "Generates the markdown cli reference",
		Run:     cmdutil.DefaultRunFunc(o),
		Example: "kore gen docs [options]",
	}

	flags := cmd.Flags()
	flags.StringVar(&o.Directory, "directory", "generated", "directory to write the generated files `PATH`")

	return cmd
}

// Run implements the action
func (o *DocOptions) Run() error {
	if err := os.MkdirAll(o.Directory, os.FileMode(0775)); err != nil {
		return err
	}

	err := doc.GenMarkdownTreeCustom(o.Root, o.Directory,
		func(filename string) string {
			now := time.Now().Format(time.RFC3339)
			name := filepath.Base(filename)
			base := strings.TrimSuffix(name, path.Ext(name))
			url := "/commands/" + strings.ToLower(base) + "/"

			return fmt.Sprintf(fmTemplate, now, strings.Replace(base, "_", " ", -1), base, url)
		},
		func(name string) string {
			base := strings.TrimSuffix(name, path.Ext(name))

			return "/commands/" + strings.ToLower(base) + "/"
		},
	)

	return err
}

// Validate checks the options
func (o *DocOptions) Validate() error {
	return nil
}
