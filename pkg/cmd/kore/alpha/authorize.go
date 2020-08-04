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

package alpha

import (
	"errors"
	"path"

	"github.com/appvia/kore/pkg/apiserver/types"
	"github.com/appvia/kore/pkg/client"
	clientcfg "github.com/appvia/kore/pkg/client/config"
	cmderr "github.com/appvia/kore/pkg/cmd/errors"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/render"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	getAuthoizeLongDescription = `
Authorize is a kubectl credentials plugin used to retrieve an authentication
token from Kore in order to gain access to your Kubernetes infrastructure.
The command swaps the authentication methods in the currently selected profile
for a authentication token for access to the clusters.
`
)

const (
	// ExecAPIVersion is the apiversion we produce
	ExecAPIVersion = "client.authentication.k8s.io/v1beta1"
)

var (
	// AuthorizedTokensPath is the path to the authorized tokens
	AuthorizedTokensPath = path.Join(utils.UserHomeDir(), ".kore", "authorized")
)

// AuthorizeOptions are the options for the authorize command
type AuthorizeOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	// AuthorizedTokenPath is the path to the tokens file
	AuthorizedTokenPath string
}

// NewCmdAlphaAuthorize creates and return the authorize command
func NewCmdAlphaAuthorize(factory cmdutil.Factory) *cobra.Command {
	o := &AuthorizeOptions{Factory: factory}

	cmd := &cobra.Command{
		Use:   "authorize",
		Long:  getAuthoizeLongDescription,
		Short: "Authorize myself and retrieve a kubernetes token",
		Run:   cmdutil.DefaultRunFunc(o),
	}

	flags := cmd.Flags()
	flags.StringVar(&o.AuthorizedTokenPath, "authorized-tokens", AuthorizedTokensPath, "path to authorized token `PATH`")

	return cmd
}

// Run implements the action
func (o *AuthorizeOptions) Run() error {
	// @step: find the local token if any
	found, err := utils.FileExists(o.AuthorizedTokenPath)
	if err != nil {
		return err
	}

	if found {
		claims, err := o.GetCurrentKubeToken(o.AuthorizedTokenPath)
		if err != nil {
			log.WithError(err).Debug("trying to retrieve current token")
		}
		if err == nil && !claims.HasExpired() {
			log.Debug("issued token is still valid, using local")

			return o.RenderKubectlCredentials(claims)
		}
		log.Debug("issues token has expired, retrieving a new one")
	}

	// @step: retrieve a token from kore
	claims, err := o.RequestKubeToken()
	if err != nil {
		return err
	}

	// @step: write the token to file
	if err := o.UpdateKubeToken(o.AuthorizedTokenPath, string(claims.RawToken)); err != nil {
		return err
	}

	return o.RenderKubectlCredentials(claims)
}

// UpdateKubeToken is used to update the token in the token configuration
func (o *AuthorizeOptions) UpdateKubeToken(path, token string) error {
	config := &clientcfg.AuthorizedTokensConfig{}

	found, err := utils.FileExists(path)
	if err != nil {
		return err
	}
	if found {
		config, err = clientcfg.ParseTokenConfigurationFromFile(path)
		if err != nil {
			return err
		}
	}
	if config.AuthInfos == nil {
		config.AuthInfos = make(map[string]*clientcfg.AuthInfo)
	}
	config.AuthInfos[o.Client().CurrentProfile()] = &clientcfg.AuthInfo{Token: &token}

	return clientcfg.WriteTokenConfigurationToFile(config, path)
}

// GetCurrentKubeToken returns the current token from file
func (o *AuthorizeOptions) GetCurrentKubeToken(path string) (*utils.JWTToken, error) {
	config, err := clientcfg.ParseTokenConfigurationFromFile(path)
	if err != nil {
		return nil, err
	}

	if config.AuthInfos == nil {
		return nil, errors.New("no authentication found")
	}

	auth, found := config.AuthInfos[o.Client().CurrentProfile()]
	if !found {
		return nil, errors.New("no profile authentication found")
	}
	if auth.Token == nil {
		return nil, errors.New("no proile token found")
	}

	return utils.NewJWTTokenFromBytes([]byte(*auth.Token))
}

// RequestKubeToken is used to request a token from kore
func (o *AuthorizeOptions) RequestKubeToken() (*utils.JWTToken, error) {
	current := o.Client().CurrentProfile()

	// @step: retrieve a token from kore - swapping our current auth
	auth := o.Config().GetAuthInfo(current)
	if auth == nil {
		return nil, cmderr.NewProfileInvalidError("no authentication details found in profile", current)
	}

	var method string
	switch {
	case auth.BasicAuth != nil:
		method = "basicauth"
	case auth.OIDC != nil:
		method = "openid"
	case auth.Token != nil:
		return nil, cmderr.NewProfileInvalidError("token authentication cannot be used to authorize", current)
	default:
		return nil, cmderr.NewProfileInvalidError("unknown authentication mode in profile", current)
	}

	token := &types.IssuedToken{}

	// @step: we need to exchange the token for a kore minted version
	err := o.ClientWithEndpoint("/issue/authorize").
		Parameters(client.QueryParameter("method", method)).
		Result(token).
		Update().
		Error()
	if err != nil {
		return nil, err
	}

	return utils.NewJWTTokenFromBytes(token.Token)
}

// RenderKubectlCredentials is used to render the credential to screen
func (o *AuthorizeOptions) RenderKubectlCredentials(token *utils.JWTToken) error {
	exp, found := token.GetExpiry()
	if !found {
		return errors.New("no expiration in token")
	}
	expires := metav1.NewTime(exp)

	status := map[string]interface{}{
		"apiVersion": ExecAPIVersion,
		"kind":       "ExecCredential",
		"status": map[string]interface{}{
			"expirationTimestamp": &expires,
			"token":               string(token.RawToken),
		},
	}

	return render.Render().
		Writer(o.Writer()).
		Resource(render.FromStruct(status)).
		Format(render.FormatJSON).
		Do()
}
