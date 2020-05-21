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

package bootstrap

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/appvia/kore/pkg/cmd/kore/bootstrap/providers"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

// DestroyOptions are the options for bringing down the cluster
type DestroyOptions struct {
	cmdutil.Factory
	// Name is an optional name for the resource
	Name string
	// Provider is the cloud provider to use
	Provider string
	// logger is a internal logger
	logger providers.Logger
}

// NewCmdBootstrapDestroy creates and returns the bootstrap destroy command
func NewCmdBootstrapDestroy(factory cmdutil.Factory) *cobra.Command {
	o := &DestroyOptions{Factory: factory, logger: newProviderLogger(factory)}

	command := &cobra.Command{
		Use:     "destroy",
		Short:   "Brings down kore on a local kubernetes cluster",
		Long:    usage,
		Example: "kore alpha bootstrap destroy <name> [options]",
		Run:     cmdutil.DefaultRunFunc(o),
	}

	flags := command.Flags()
	flags.StringVar(&o.Provider, "provider", "kind", "local kubernetes provider to use")

	return command
}

// Validate checks the options
func (o *DestroyOptions) Validate() error {
	return nil
}

// Run implements the action
func (o *DestroyOptions) Run() error {
	o.Name = ClusterName

	tasks := []TaskFunc{
		o.EnsurePreflightChecks,
		o.EnsureLocalKubernetesDeletion,
	}
	for _, x := range tasks {
		if err := x(context.TODO()); err != nil {
			return err
		}
	}

	return nil
}

// EnsurePreflightChecks is responsible for have everything moving forward
func (o *DestroyOptions) EnsurePreflightChecks(ctx context.Context) error {
	return (&Task{
		Description: "Passed preflight checks for deployment",
		Handler: func(ctx context.Context) error {
			for _, x := range []string{o.Provider, Kubectl} {
				if _, err := exec.LookPath(x); err != nil {
					return fmt.Errorf("missing binary: %s in $PATH", x)
				}
			}

			return nil
		},
	}).Run(ctx, o.Writer())
}

// EnsureLocalKubernetesDeletion is responsible for deleting the local instance
func (o *DestroyOptions) EnsureLocalKubernetesDeletion(ctx context.Context) error {
	provider, err := GetProvider(o.Factory, o.Provider)
	if err != nil {
		return err
	}

	// @step: perform the preflight checks for the provider
	if err := (&Task{
		Description: "Passed preflight checks for local cluster provider",
		Handler: func(ctx context.Context) error {
			return provider.Preflight(ctx)
		},
	}).Run(ctx, o.Writer()); err != nil {
		return err
	}

	return (&Task{
		Header:      "Attempting to delete the local kubernetes cluster",
		Description: "Removed the local kubernetes cluster",
		Handler: func(ctx context.Context) error {
			return provider.Destroy(ctx, o.Name)
		},
	}).Run(ctx, o.Writer())
}
