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
	"fmt"
	"strings"

	"github.com/appvia/kore/pkg/cmd/kore/local/providers"
	"github.com/appvia/kore/pkg/cmd/kore/local/providers/kind"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

var (
	usage = `
Bootstrap provides an experimental means of bootstrapping a local Kore installation. At
present the local installation uses "kind" https://github.com/kubernetes-sigs/kind.

Unless specified otherwise it will deploy an official tagged release from Github, though
this can be overridden using the --release flag. Note the installation is performed via
helm with a local ${HOME}/.kore/values.yaml generated in the directory. If you wish to change
any of the values post installation, update the values.yaml file and re-run the 'up' command.

Note the data persistence is tied to the installation provider. For kind as long
as the container is not delete the data is kept.
`
	examples = `
# Provision a local kore instance called 'kore' (defaults to kind)
$ kore alpha local up

# By default we will download the official tagged release; you can however override this
# behaviour by using the --release and --version flags. Note, its best to leave these unless
# you know explicitly know what you need to override.

$ kore alpha local up --release ./charts/kore
$ kore alpha local up --release https://URL

# Override the version
$ kore alpha local up --version latest --release ./charts/kore

# Destroy the local installation
$ kore alpha local destroy

# To stop the local installed without deleting the data
$ kore alpha local stop

The application should be available on http://127.0.0.1:3000. You can provision the
CLI via.

$ kore login -a http://127.0.0.1:10080 local

Post the command your Kubectl context is switched to the kind installation:

$ kubectl config current-context
`
)

var (
	// GithubRelease is the link to release
	GithubRelease = "https://github.com/appvia/kore/releases/download/%s/kore-helm-chart-%s.tgz"
	// ClusterName is the name of the cluster to create
	ClusterName = "kore"
)

const (
	// Kubectl is the name of the binary
	Kubectl = "kubectl"
)

// NewCmdBootstrap creates and returns the delete command
func NewCmdBootstrap(factory cmdutil.Factory) *cobra.Command {
	command := &cobra.Command{
		Use:     "local",
		Short:   "Provides the provision of local installation for Kore for testing",
		Long:    usage,
		Example: examples,
		Run:     cmdutil.RunHelp,
	}

	command.AddCommand(
		NewCmdBootstrapDestroy(factory),
		NewCmdBootstrapUp(factory),
		NewCmdBootstrapStop(factory),
	)

	return command
}

// GetHelmReleaseURL returns the helm release for kore
func GetHelmReleaseURL(release string) string {
	if strings.HasPrefix(release, "v") {
		return fmt.Sprintf(GithubRelease, release, release)
	}

	return release
}

// GetProvider returns the provider implementation
func GetProvider(f cmdutil.Factory, name string) (providers.Interface, error) {
	switch name {
	case "kind":
		return kind.New(newProviderLogger(f))
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}

// AddProviderFlags is used to add the provider specific options to the command line
func AddProviderFlags(cmd *cobra.Command) {
	kind.AddProviderFlags(cmd)
}

// GetProviders returns a list of support local cluster providers
func GetProviders() []string {
	return []string{"kind"}
}
