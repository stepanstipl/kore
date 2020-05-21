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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/httputils"
	"github.com/appvia/kore/pkg/version"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/chart/loader"
)

// UpOptions are the options for bootstrapping
type UpOptions struct {
	cmdutil.Factory
	// Name is an optional name for the resource
	Name string
	// Release is tagged release to use
	Release string
	// Provider is the cloud provider to use
	Provider string
	// ContextName is the name of the kubernetes context
	ContextName string
	// EnableDeploy indicates we should deploy the application
	EnableDeploy bool
	// Force indicates we should force any changes
	Force bool
	// Wait indicates we wait for the deployment to finish
	Wait bool
	// ValuesFile is the file containing the configurable values
	ValuesFile string
	// Values a collection of values passed to the helm chart
	Values map[string]interface{}
}

// NewCmdBootstrapUp creates and returns the bootstrap up command
func NewCmdBootstrapUp(factory cmdutil.Factory) *cobra.Command {
	o := &UpOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "up",
		Short:   "Brings up kore on a local kubernetes cluster",
		Long:    usage,
		Example: "kore alpha bootstrap up <name> [options]",
		Run:     cmdutil.DefaultRunFunc(o),
	}

	flags := command.Flags()
	flags.StringVar(&o.Provider, "provider", "kind", "local kubernetes provider to use `NAME`")
	flags.StringVar(&o.Release, "release", version.Release, "chart version to use for deployment `VERSION`")
	flags.StringVar(&o.ValuesFile, "values", "values.yaml", "path to the file container helm values `PATH`")
	flags.BoolVar(&o.EnableDeploy, "enable-deploy", true, "indicates if we should deploy the kore application `BOOL`")
	flags.BoolVar(&o.Wait, "wait", true, "indicates we wait for the deployment to complete `BOOL`")
	flags.BoolVar(&o.Force, "force", false, "indicates we should force any changes `BOOL`")

	return command
}

type authInfoConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	AuthorizeURL string `yaml:"authorize_url"`
}

// Validate checks the options
func (o *UpOptions) Validate() error {
	return nil
}

// Run implements the action
func (o *UpOptions) Run() error {
	o.Name = ClusterName

	tasks := []TaskFunc{
		o.EnsurePreflightChecks,
		o.EnsureHelmValues,
		o.EnsureLocalKubernetes,
		o.EnsureKubernetesContext,
		o.EnsureKoreRelease,
		o.EnsureUP,
	}
	for _, x := range tasks {
		if err := x(context.TODO()); err != nil {
			return err
		}
	}

	o.Println("")
	o.Println("You can access the Kore portal via http://localhost:3000")
	o.Println("Configure your CLI via $ kore login -a http://localhost:10080 local")
	o.Println("")

	return nil
}

// EnsurePreflightChecks is responsible for have everything moving forward
func (o *UpOptions) EnsurePreflightChecks(ctx context.Context) error {
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

// EnsureHelmValues is responsible for retrieve the values for helm
func (o *UpOptions) EnsureHelmValues(ctx context.Context) error {
	found, err := utils.FileExists(o.ValuesFile)
	if err != nil {
		return err
	}
	if !found {
		o.Println("First time running bootstrap, we need your IDP settings for OpenID")
	}

	o.Values, err = GetHelmValues(o.ValuesFile)
	if err != nil {
		return err
	}

	if !found {
		if err := (&Task{
			Description: fmt.Sprintf("Persisting the values to local file: %s", o.ValuesFile),
			Handler: func(ctx context.Context) error {
				content, err := utils.ToYAML(&o.Values)
				if err != nil {
					return err
				}

				return ioutil.WriteFile(o.ValuesFile, content, os.FileMode(0660))
			},
		}).Run(ctx, o.Writer()); err != nil {
			return err
		}
	}

	return nil
}

// EnsureLocalKubernetes is responsible for provisioning the local instance
func (o *UpOptions) EnsureLocalKubernetes(ctx context.Context) error {
	provider, err := GetProvider(o.Factory, o.Provider)
	if err != nil {
		return err
	}

	if err := (&Task{
		Description: "Passed preflight checks for local cluster provider",
		Handler: func(ctx context.Context) error {
			return provider.Preflight(ctx)
		},
	}).Run(ctx, o.Writer()); err != nil {
		return err
	}

	return (&Task{
		Header:      "Attempting to build the local kubernetes cluster",
		Description: "Provisioned a local kubernetes cluster",
		Handler: func(ctx context.Context) error {
			if err := provider.Create(ctx, o.Name); err != nil {
				return err
			}

			return (&Task{
				Description: "Exported the kubeconfig from provisioned cluster",
				Handler: func(ctx context.Context) error {
					name, err := provider.Export(ctx, o.Name)
					if err != nil {
						return err
					}
					o.ContextName = name

					return nil
				},
			}).Run(ctx, o.Writer())
		},
	}).Run(ctx, o.Writer())
}

// EnsureKubernetesContext is responsible for setting the kubectl context
func (o *UpOptions) EnsureKubernetesContext(ctx context.Context) error {
	return (&Task{
		Description: fmt.Sprintf("Configured the kubectl context: %s", o.ContextName),
		Handler: func(ctx context.Context) error {
			args := []string{
				"config",
				"set-context",
				o.ContextName,
			}
			path, err := exec.LookPath(Kubectl)
			if err != nil {
				return err
			}

			return exec.CommandContext(ctx, path, args...).Run()
		},
	}).Run(ctx, o.Writer())
}

// EnsureKoreRelease is responsible for deploying the release into the cluster
func (o *UpOptions) EnsureKoreRelease(ctx context.Context) error {
	if !o.EnableDeploy {
		return nil
	}
	path := o.Release

	// @step; do we need to download the chart
	// @TODO change this - the code is chaotic
	url, found, err := func() (string, bool, error) {
		if strings.HasPrefix(path, "http") {
			return o.Release, true, nil
		}
		found, err := utils.DirExists(path)
		if err != nil {
			return o.Release, false, err
		}
		if !found {
			return GetHelmReleaseURL(path), true, nil
		}

		return "", false, nil
	}()
	if err != nil {
		return err
	}
	if found {
		// store in the tmp folder for now - helms downloader is very interwined
		path = filepath.Join(os.TempDir(), filepath.Base(url))

		if exists, err := utils.FileExists(path); err != nil {
			return err

		} else if !exists {
			if err := (&Task{
				Description: "Downloading the official helm chart",
				Handler: func(ctx context.Context) error {
					return utils.DownloadFile(ctx, path, url)
				},
			}).Run(ctx, o.Writer()); err != nil {
				return err
			}
		}
	}

	return (&Task{
		Header:      fmt.Sprintf("Attempting to deploy the Kore release %s", o.Release),
		Description: "Deployed the Kore release into the cluster",
		Handler: func(ctx context.Context) error {
			chart, err := loader.Load(path)
			if err != nil {
				return err
			}
			copied := utils.CopyMap(chart.Values)

			patched, err := utils.MergeJSON(&o.Values, &copied)
			if err != nil {
				return err
			}

			values := make(map[string]interface{})
			if err := json.NewDecoder(bytes.NewReader(patched)).Decode(&values); err != nil {
				return err
			}

			return HelmDeploy(ctx, chart, values, "kore", "kore", o.ContextName)
		},
	}).Run(ctx, o.Writer())
}

// EnsureUP is responsible for checking the service is up
func (o *UpOptions) EnsureUP(ctx context.Context) error {
	if !o.EnableDeploy || !o.Wait {
		return nil
	}
	timeout := 5 * time.Minute
	interval := 2 * time.Second

	return (&Task{
		Header:      fmt.Sprintf("Waiting for deployment to rollout successfully (%s timeout)", timeout.String()),
		Description: "Successfully deployed the kore releasse to cluster",
		Handler: func(ctx context.Context) error {
			hc := httputils.DefaultHTTPClient

			err := utils.WaitUntilComplete(ctx, timeout, interval, func() (bool, error) {
				resp, err := hc.Get("http://localhost:10080/healthz")
				if err == nil && resp.StatusCode == http.StatusOK {
					return true, nil
				}

				return false, nil
			})
			if err != nil {
				return fmt.Errorf("deployment unsuccessful, please check via kubectl -n kore get po")
			}

			return nil
		},
	}).Run(ctx, o.Writer())
}
