/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package korectl

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v2"
)

var (
	loginLongDescription = `
Is used to authenticate yourself to the currently selected profile. Login
performs an oauth2 authentication flow and retrieve your identity token.

Examples:
$ korectl login    # will login and update the current profile
$ korectl login local -a http://127.0.0.1:8080  # create a profile and login
`
)

// GetLoginCommand is used to login to the api server
func GetLoginCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "login",
		Description: loginLongDescription,
		Usage:       "Authenticate yourself and retrieve a token for Appvia Kore",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "api-server,a",
				Usage: "allows you to specify the kore api server to login as `URL`",
			},
			&cli.BoolFlag{
				Name:  "force,f",
				Usage: "must be set when you want to override the api-server on an existing profile `BOOL`",
			},
			&cli.IntFlag{
				Name:  "port,p",
				Usage: "sets the local port used for redirection when authenticating",
				Value: 3001,
			},
		},

		Action: func(ctx *cli.Context) error {
			var err error

			// @step: if the api server and profile name is passed we can create a profile
			// from the name and add the server - this is effectively a inline `profile configure`
			endpoint := ctx.String("api-server")
			if endpoint != "" {
				switch {
				case !ctx.Args().Present():
					return fmt.Errorf("you must specify a profile name when logging in")

				case config.HasProfile(ctx.Args().First()) && !ctx.Bool("force"):
					return fmt.Errorf("profile name already used (note: you can use the --force option to force the update)")

				case !IsValidHostname(endpoint):
					return fmt.Errorf("invalid api server: %s", endpoint)
				}
				profile := ctx.Args().First()

				config.CreateProfile(profile, endpoint)
				config.SetCurrentProfile(profile)
			}

			// @check we have the minimum required for authentication
			if err := config.HasValidProfile(); err != nil {
				fmt.Println("Unable to authenticate:", err.Error())
				fmt.Println("You may need to reconfigure your profile via $ korectl profile configure")

				os.Exit(1)
			}

			// @step: we make done channels to signal events
			doneCh := make(chan struct{})
			errCh := make(chan error)

			var token *AuthorizationResponse

			// @step: we create a local http server to order to handle the callback
			http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
				token, err = handleLoginCallback(req, w)
				if err != nil {
					errCh <- err
					return
				}
				doneCh <- struct{}{}
			})

			serverPort := ctx.Int("port")

			// @step: we need to start the http server in the background
			go func() {
				listenAddress := fmt.Sprintf(":%d", serverPort)
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
				config.GetCurrentServer().Endpoint, serverPort)

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
			auth.OIDC = &OIDC{
				AccessToken:  token.AccessToken,
				AuthorizeURL: token.AuthorizationURL,
				ClientID:     token.ClientID,
				ClientSecret: token.ClientSecret,
				IDToken:      token.IDToken,
				RefreshToken: token.RefreshToken,
				TokenURL:     token.TokenEndpointURL,
			}

			// @step: update the local configuration on disk
			if err := config.Update(); err != nil {
				return fmt.Errorf("trying to update the client configuration: %s", err)
			}
			fmt.Println("Successfully authenticated")

			return nil
		},
	}
}

// handleLoginCallback is used to handle the callback from the api server
func handleLoginCallback(req *http.Request, resp http.ResponseWriter) (*AuthorizationResponse, error) {
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

	token := &AuthorizationResponse{}

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
