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

package login

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/appvia/kore/pkg/apiserver"
	restconfig "github.com/appvia/kore/pkg/client/config"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var (
	loginLongDescription = `
Is used to authenticate yourself to the currently selected profile. Login
performs an oauth2 authentication flow and retrieve your identity token.

$ korectl login    # will login and update the current profile
$ korectl login local -a http://127.0.0.1:8080  # create a profile and login
`
)

// LoginOptions are the options for logging in
type LoginOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	// Name is used when creating a profile
	Name string
	// Endpoint is an optional endpoint
	Endpoint string
	// Force is used to force an operation
	Force bool
	// Port is the local port to use for http server
	Port int
}

// NewCmdLogin creates and returns a login command
func NewCmdLogin(factory cmdutil.Factory) *cobra.Command {
	o := &LoginOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "login",
		Short:   "Authenticate yourself and retrieve a token for Appvia Kore",
		Long:    loginLongDescription,
		Example: "kore login [-a endpoint name]",
		Run:     cmdutil.DefaultRunFunc(o),
	}

	flags := command.Flags()
	flags.StringVarP(&o.Endpoint, "api-url", "a", "", "specify the kore api server to login `URL`")
	flags.BoolVarP(&o.Force, "force", "f", false, "must be set when you want to override the api-server on an existing profile `BOOL`")
	flags.IntVarP(&o.Port, "port", "p", 3001, "sets the local port used for redirection when authenticating `PORT`")

	return command
}

// Validate is used to validate the parameters
func (o *LoginOptions) Validate() error {
	config := o.Config()

	// @step: if the api server and profile name is passed we can create a profile
	// from the name and add the server - this is effectively a inline `profile configure`
	if o.Endpoint != "" {
		switch {
		case o.Name == "":
			return fmt.Errorf("you must specify a profile name when using endpoint -a")

		case config.HasProfile(o.Name) && !o.Force:
			return fmt.Errorf("profile name already used (note: you can use the --force option to force the update)")

			/*
				case !IsValidHostname(o.Endpoint):
					return fmt.Errorf("invalid api server: %s", o.Endpoint)
			*/
		}

		config.CreateProfile(o.Name, o.Endpoint)
		config.CurrentProfile = o.Name
	}

	return nil
}

// Run performs the command action to login
func (o *LoginOptions) Run() error {
	var err error

	config := o.Config()

	// @check we have the minimum required for authentication
	if err := o.Config().HasValidProfile(); err != nil {
		o.Println("Unable to authenticate: %s", err.Error())
		o.Println("You may need to reconfigure your profile via $ korectl profile configure")

		return errors.New("invalid profile")
	}

	// @step: we make done channels to signal events
	doneCh := make(chan struct{})
	errCh := make(chan error)

	token := &apiserver.AuthorizationResponse{}

	// @step: we create a local http server to order to handle the callback
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		token, err = handleLoginCallback(req, w)
		if err != nil {
			errCh <- err
			return
		}
		doneCh <- struct{}{}
	})

	// @step: we need to start the http server in the background
	go func() {
		listenAddress := fmt.Sprintf(":%d", o.Port)
		if err := http.ListenAndServe(listenAddress, nil); err != nil {
			errCh <- fmt.Errorf("trying to start local http server: %s", err)
		}
	}()

	fmt.Printf("Attempting to authenticate to Appvia Kore: %s [%s]\n",
		config.GetCurrentServer().Endpoint,
		config.CurrentProfile,
	)

	// @step: open a browser to the to the api server
	redirectURL := fmt.Sprintf("%s/oauth/authorize?redirect_url=http://localhost:%d",
		config.GetCurrentServer().Endpoint, o.Port)

	if err := open.Run(redirectURL); err != nil {
		return fmt.Errorf("trying to open web browser, error: %s", err)
	}

	// @step: we wait for either a done or error or timeout
	select {
	case <-doneCh:
	case err := <-errCh:
		return fmt.Errorf("trying to authorize the client: %s", err)
	case <-time.After(30 * time.Second):
		return errors.New("authorization request timed out waiting to complete")
	}

	auth := config.GetCurrentAuthInfo()
	auth.OIDC = &restconfig.OIDC{
		AccessToken:  token.AccessToken,
		AuthorizeURL: token.AuthorizationURL,
		ClientID:     token.ClientID,
		ClientSecret: token.ClientSecret,
		IDToken:      token.IDToken,
		RefreshToken: token.RefreshToken,
		TokenURL:     token.TokenEndpointURL,
	}

	// @step: update the local configuration on disk
	if err := o.UpdateConfig(); err != nil {
		return fmt.Errorf("trying to update the client configuration: %s", err)
	}

	o.Println("Successfully authenticated")

	return nil
}

// handleLoginCallback is used to handle the callback from the api server
func handleLoginCallback(req *http.Request, resp http.ResponseWriter) (*apiserver.AuthorizationResponse, error) {
	// @step: check we have a token in the return
	if req.URL.RawQuery == "" {
		return nil, errors.New("no token found in the authorization request")
	}
	if !strings.HasPrefix(req.URL.RawQuery, "token=") {
		return nil, errors.New("invalid token response from apiserver")
	}
	raw := strings.TrimPrefix(req.URL.RawQuery, "token=")

	// @step: extract and decode the token
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, err
	}

	g, err := gzip.NewReader(bytes.NewReader(decoded))
	if err != nil {
		return nil, err
	}
	g.Close()

	rb, err := ioutil.ReadAll(g)
	if err != nil {
		return nil, err
	}

	token := &apiserver.AuthorizationResponse{}

	if err := json.NewDecoder(bytes.NewReader(rb)).Decode(token); err != nil {
		return nil, err
	}

	// @step: send back the html
	shutdown := `<html><body><script>window.close();</script></body></html>`
	if _, err := resp.Write([]byte(shutdown)); err != nil {
		return nil, err
	}

	return token, nil
}
