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

package kubeconfig

import (
	"errors"
	"fmt"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	cmderr "github.com/appvia/kore/pkg/cmd/errors"
	"github.com/appvia/kore/pkg/cmd/kore/alpha"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	"github.com/appvia/kore/pkg/utils/render"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

const (
	// KoreContextUser is the username to use in the kubeconfig
	KoreContextUser = "kore"
)

// KubeConfigOptions is used to provision a team
type KubeConfigOptions struct {
	cmdutil.Factory
	// KorePath is the location of the kore binary
	KorePath string
	// Team string
	Team string
	// Output is the output format
	Output string
	// Headers indicates no headers on the table output
	Headers bool
}

// NewCmdKubeConfig returns the create admin command
func NewCmdKubeConfig(factory cmdutil.Factory) *cobra.Command {
	o := &KubeConfigOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "kubeconfig",
		Short:   "Adds your team provisioned clusters to your kubeconfig",
		Example: "kore kubeconfig [-t team]",
		Run:     cmdutil.DefaultRunFunc(o),
	}

	return command
}

// Validate is called to validate the options
func (o *KubeConfigOptions) Validate() error {
	if o.Team == "" {
		return errors.New("no team defined")
	}

	return nil
}

// Run implements the action
func (o *KubeConfigOptions) Run() error {
	clusters := &clustersv1.ClusterList{}

	resp, err := o.ClientWithTeamResource(o.Team, o.Resources().MustLookup("cluster")).
		Result(clusters).
		Get().Do()
	if err != nil {
		return err
	}

	if len(clusters.Items) <= 0 {
		o.Println("No clusters found in team")

		return nil
	}

	path, err := kubernetes.GetOrCreateKubeConfig()
	if err != nil {
		return err
	}

	if err := o.WriteConfig(clusters, path); err != nil {
		return err
	}

	return render.Render().
		Writer(o.Writer()).
		Resource(render.FromReader(resp.Body())).
		Format(o.Output).
		ShowHeaders(o.Headers).
		Foreach("items").
		Printer(
			render.Column("Context", "metadata.name"),
			render.Column("Cluster", "metadata.name"),
			render.Column("Endpoint", "status.authProxyEndpoint"),
		).Do()
}

// WriteConfig is responsible for updating the users kubeconfig
func (o *KubeConfigOptions) WriteConfig(clusters *clustersv1.ClusterList, path string) error {
	cfg, err := clientcmd.LoadFromFile(path)
	if err != nil {
		return err
	}

	// @step: we default to assuming the authentication method is sso
	authUser := KoreContextUser

	current := o.Client().CurrentProfile()
	auth := o.Config().GetAuthInfo(current)
	switch {
	case auth.OIDC != nil:
		//
		// @TODO: this comes with the caveat that they are using the 'SAME' idp
		// across all kore instances (which might not be the case)
		//
		cfg.AuthInfos[KoreContextUser] = &api.AuthInfo{
			AuthProvider: &api.AuthProviderConfig{
				Name: "oidc",
				Config: map[string]string{
					"access-token":   auth.OIDC.AccessToken,
					"client-id":      auth.OIDC.ClientID,
					"client-secret":  auth.OIDC.ClientSecret,
					"id-token":       auth.OIDC.IDToken,
					"idp-issuer-url": auth.OIDC.AuthorizeURL,
					"refresh-token":  auth.OIDC.RefreshToken,
				},
			},
		}

	case auth.BasicAuth != nil:
		authUser = fmt.Sprintf("kore-pf-%s", current)

		cfg.AuthInfos[authUser] = &api.AuthInfo{
			Exec: &api.ExecConfig{
				APIVersion: alpha.ExecAPIVersion,
				Command:    o.KorePath,
				Args: []string{
					"--profile=" + o.Client().CurrentProfile(),
					"alpha",
					"authorize",
				},
			},
		}

	default:
		return cmderr.NewProfileInvalidError("authentication method not supported", current)
	}

	for _, x := range clusters.Items {
		if x.Status.Status == corev1.DeletingStatus {
			continue
		}
		if x.Status.AuthProxyEndpoint == "" {
			o.Println("SKIPPING CLUSTER: %s as it does not have an endpoint yet", x.Name)
			continue
		}
		// ensure we have entries
		if cfg.Clusters == nil {
			cfg.Clusters = make(map[string]*api.Cluster)
		}
		if cfg.Contexts == nil {
			cfg.Contexts = make(map[string]*api.Context)
		}

		cfg.Clusters[x.Name] = &api.Cluster{
			InsecureSkipTLSVerify: true,
			Server:                "https://" + x.Status.AuthProxyEndpoint,
		}
		cfg.Contexts[x.Name] = &api.Context{
			AuthInfo: authUser,
			Cluster:  x.Name,
		}
	}

	return clientcmd.WriteToFile(*cfg, path)
}
