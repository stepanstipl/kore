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

package local

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/appvia/kore/pkg/client/config"
	"github.com/appvia/kore/pkg/cmd/kore/local/providers"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/httputils"
	ksutils "github.com/appvia/kore/pkg/utils/kubernetes"
	"github.com/appvia/kore/pkg/version"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	// helmVersion is the version of helm to download
	helmVersion = "v3.2.1"
	// valueFilePerms is the file permissions on the values.yaml
	valueFilePerms = os.FileMode(0600)
)

// UpOptions are the options for bootstrapping
type UpOptions struct {
	cmdutil.Factory
	// BinaryPath is the directory to place downloaded binaries
	BinaryPath string
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
	// EnableSSO indicates that single-sign-on details should be prompted
	EnableSSO bool
	// DisableUI indicates we deploy without an UI
	DisableUI bool
	// DeploymentTimeout is the amount of time we will wait for deployment
	DeploymentTimeout time.Duration
	// Force indicates we should force any changes
	Force bool
	// FlagsChanged is a list of flags which changed
	FlagsChanged []string
	// HelmPath is the path to the helm binary
	HelmPath string
	// LocalAdminPassword is the password for localadmin
	LocalAdminPassword string
	// Wait indicates we wait for the deployment to finish
	Wait bool
	// ValuesFile is the file containing the configurable values
	ValuesFile string
	// Values for helm chart
	Values map[string]interface{}
	// HelmValues a collection of values passed to the helm chart
	HelmValues []string
	// Version is the release version to use
	Version string
}

// NewCmdBootstrapUp creates and returns the bootstrap up command
func NewCmdBootstrapUp(factory cmdutil.Factory) *cobra.Command {
	o := &UpOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "up",
		Short:   "Brings up kore on a local kubernetes cluster",
		Long:    usage,
		Example: "kore alpha local up <name> [options]",
		Run:     cmdutil.DefaultRunFunc(o),
	}

	flags := command.Flags()
	flags.StringVar(&o.Provider, "provider", "kind", "local kubernetes provider to use `NAME`")
	flags.StringVar(&o.Release, "release", version.Tag, "chart version to use for deployment `CHART`")
	flags.StringVar(&o.Version, "version", version.Tag, "kore version to deployment into cluster `VERSION`")
	flags.StringVar(&o.ValuesFile, "values", os.ExpandEnv(filepath.Join(utils.UserHomeDir(), ".kore", "values.yaml")), "path to the file container helm values `PATH`")
	flags.StringVar(&o.BinaryPath, "binary-path", filepath.Join(config.GetClientPath(), "build"), "path to place any downloaded binaries if requested `PATH`")
	flags.BoolVar(&o.EnableDeploy, "enable-deploy", true, "indicates if we should deploy the kore application `BOOL`")
	flags.BoolVar(&o.EnableSSO, "enable-sso", false, "indicates we want use a openid provider for authentication `BOOL`")
	flags.StringVar(&o.LocalAdminPassword, "local-admin-password", "", "the password for local admin `PASSWORD`")
	flags.BoolVar(&o.DisableUI, "disable-ui", false, "indicates the kore ui is not deployed `BOOL`")
	flags.DurationVar(&o.DeploymentTimeout, "deployment-timeout", 5*time.Minute, "amount of time to wait for a successful deployment `DURATION`")
	flags.StringSliceVar(&o.HelmValues, "set", []string{}, "a collection of path=value used to update the helm values `KEYPAIR`")
	flags.BoolVar(&o.Wait, "wait", true, "indicates we wait for the deployment to complete `BOOL`")
	flags.BoolVar(&o.Force, "force", false, "indicates we should force any changes `BOOL`")

	// @step: add the provider specific options to the command line
	AddProviderFlags(command)

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
		o.EnsureHelm,
		o.EnsureKoreRelease,
		o.EnsureUP,
	}
	for _, x := range tasks {
		if err := x(context.TODO()); err != nil {
			return err
		}
	}

	o.Println("")
	if !o.EnableSSO {
		o.Println("Access the Kore portal via http://localhost:3000 [ admin | %s ]", o.LocalAdminPassword)
		o.Println("Configure your CLI via $ kore profile configure local -a http://localhost:10080 --account basicauth")
	} else {
		o.Println("Access the Kore portal via http://localhost:3000")
		o.Println("Configure your CLI via $ kore login -a http://localhost:10080")
	}
	o.Println("")

	return nil
}

// EnsurePreflightChecks is responsible for have everything moving forward
func (o *UpOptions) EnsurePreflightChecks(ctx context.Context) error {
	return (&Task{
		Header:      "Provisioning Kore installation",
		Description: "Passed preflight checks for installation",
		Handler: func(ctx context.Context) error {
			for _, x := range []string{Kubectl} {
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
	var err error

	o.Values, err = o.GetHelmValues(o.ValuesFile)
	if err != nil {
		return err
	}

	if err := (&Task{
		Description: fmt.Sprintf("Persisting the values to local file: %q", o.ValuesFile),
		Handler: func(ctx context.Context) error {
			content, err := utils.ToYAML(&o.Values)
			if err != nil {
				return err
			}

			return ioutil.WriteFile(o.ValuesFile, content, valueFilePerms)
		},
	}).Run(ctx, o.Writer()); err != nil {
		return err
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
			if err := provider.Create(ctx, o.Name, providers.CreateOptions{
				AskConfirmation: !o.Force,
				DisableUI:       o.DisableUI,
			}); err != nil {
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

// EnsureHelm is responsible for making sure helm binary is available
func (o *UpOptions) EnsureHelm(ctx context.Context) error {
	if !o.EnableDeploy {
		return nil
	}

	err := (&Task{
		Handler: func(ctx context.Context) error {
			// @step: can we find helm in the search path?
			path, err := exec.LookPath("helm")
			if err != nil {
				// have we downloaded it already?
				path := filepath.Join(o.BinaryPath, "helm")
				found, err := utils.FileExists(path)
				if err != nil {
					return err
				}
				if found {
					// @step: we need to check the version of helm
					args := []string{
						"version",
					}

					combined, err := exec.CommandContext(ctx, path, args...).CombinedOutput()
					if err != nil {
						return fmt.Errorf("trying to check version of helm binary: %s", combined)
					}

					if strings.Contains(string(combined), `Version:"v3`) {
						o.HelmPath = path

						return nil
					}
				}

				return o.EnsureHelmDownload(ctx)
			}
			o.HelmPath = path

			return nil
		},
	}).Run(ctx, o.Writer())

	if err != nil {
		return err
	}

	return nil
}

// EnsureHelmDownload is responsible for downloading the helm binary
func (o *UpOptions) EnsureHelmDownload(ctx context.Context) error {
	if !o.EnableDeploy {
		return nil
	}

	release := fmt.Sprintf("https://get.helm.sh/helm-%s-%s-%s.tar.gz",
		helmVersion,
		runtime.GOOS,
		runtime.GOARCH)

	// path to save the helm release
	tmpPath := filepath.Join(os.TempDir(), filepath.Base(release))
	// path to save the binary
	path := filepath.Join(o.BinaryPath, "helm")

	// @step: request permission from the user
	logger := newProviderLogger(o.Factory)
	logger.Info("Helm binary not found in $PATH")

	if o.Force {
		logger.Info("Downloading %s (%s)", release, path)
	} else {
		logger.Infof("Download %s (%s) y/N? ", release, path)

		if !utils.AskForConfirmation(os.Stdin) {
			return fmt.Errorf("missing: %q not found in $PATH", "helm")
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	if err := (&Task{
		Header:      "Attempting to download helm release from github",
		Description: fmt.Sprintf("Downloaded the helm release (%s)", path),

		Handler: func(ctx context.Context) error {
			if err := utils.DownloadFile(ctx, tmpPath, release); err != nil {
				return err
			}
			defer func() {
				_ = os.Remove(tmpPath)
			}()

			args := []string{
				"-C", filepath.Dir(path),
				"--strip-components=1",
				"-zxf",
				tmpPath,
				fmt.Sprintf("%s-%s/helm", runtime.GOOS, runtime.GOARCH),
			}

			return exec.CommandContext(ctx, "tar", args...).Run()
		},
	}).Run(ctx, o.Writer()); err != nil {
		return err
	}
	o.HelmPath = path

	return nil
}

// EnsureKoreRelease is responsible for deploying the release into the cluster
func (o *UpOptions) EnsureKoreRelease(ctx context.Context) error {
	if !o.EnableDeploy {
		return nil
	}

	return (&Task{
		Header:      fmt.Sprintf("Attempting to deploy the Kore release %q", o.Version),
		Description: "Deployed the Kore release into the cluster",
		Handler: func(ctx context.Context) error {
			logger := newProviderLogger(o.Factory)

			switch o.Release == version.Tag {
			case true:
				logger.Info("Using the official helm chart for deployment")
			default:
				logger.Info("Using the helm release: %s for deployment", o.Release)
			}

			release, err := func() (string, error) {
				if strings.HasPrefix(o.Release, "http") {
					return o.Release, nil
				}
				found, err := utils.DirExists(o.Release)
				if err != nil {
					return "", err
				}
				if found {
					return o.Release, err
				}

				return GetHelmReleaseURL(o.Release), nil
			}()
			if err != nil {
				return err
			}

			// ensure the namespace
			config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
				clientcmd.NewDefaultClientConfigLoadingRules(),
				&clientcmd.ConfigOverrides{CurrentContext: o.ContextName},
			).ClientConfig()
			if err != nil {
				return err
			}
			client, err := kubernetes.NewForConfig(config)
			if err != nil {
				return err
			}

			// @step: wait for kubernetes api
			interval := 2 * time.Second
			timeout := 60 * time.Second

			if err := ksutils.WaitOnKubeAPI(ctx, client, interval, timeout); err != nil {
				return errors.New("timed out waiting for the kubernetes api")
			}

			ns := &v1.Namespace{}
			ns.Name = "kore"

			if _, err := client.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{}); err != nil {
				if !kerrors.IsAlreadyExists(err) {
					return err
				}
			}

			args := []string{
				"upgrade",
				"--kube-context", o.ContextName,
				"--namespace", "kore",
				"--install",
				"--wait",
				"--values", o.ValuesFile,
				"kore",
				release,
			}

			combined, err := exec.CommandContext(ctx, o.HelmPath, args...).CombinedOutput()
			if err != nil {
				return fmt.Errorf("trying to deploy helm chart: %s", combined)
			}

			return nil
		},
	}).Run(ctx, o.Writer())
}

// EnsureUP is responsible for checking the service is up
func (o *UpOptions) EnsureUP(ctx context.Context) error {
	if !o.EnableDeploy || !o.Wait {
		return nil
	}
	timeout := o.DeploymentTimeout
	interval := 2 * time.Second

	return (&Task{
		Header:      fmt.Sprintf("Waiting for deployment to rollout successfully (%s timeout)", timeout.String()),
		Description: "Successfully deployed the kore release to cluster",
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
