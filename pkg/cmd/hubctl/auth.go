/**
 * Copyright (C) 2020 Rohith Jayawardene <info@appvia.io>
 *
 * This file is part of hub-apiserver.
 *
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package hubctl

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli"
)

// @TODO Change everthing here and polish it up

func GetLoginCommand(config Config) cli.Command {
	return cli.Command{
		Name:    "authorize",
		Aliases: []string{"auth"},
		Usage:   "Used to authenticate yourself and retrieve an authentication token",

		Action: func(ctx *cli.Context) error {
			if config.Server == "" {
				return errors.New("the 'server' field in the config needs setting the kore api service url")
			}

			// @step we need to perform a oauth2 login
			doneCh := make(chan struct{})
			errCh := make(chan error)
			var token *AuthorizationResponse

			// @step: we need to open a local http server to handle the callback
			http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
				err := func() error {
					if req.URL.RawQuery == "" {
						return errors.New("no token found in the authorization request")
					}
					if !strings.HasPrefix(req.URL.RawQuery, "token=") {
						return errors.New("invalid token response from apiserver")
					}
					raw := strings.TrimPrefix(req.URL.RawQuery, "token=")

					t, err := DecodeAuthorizationResponse(strings.NewReader(raw))
					if err != nil {
						return err
					}
					token = t

					return nil
				}()
				if err != nil {
					errCh <- err
					return
				}
				go func() {
					doneCh <- struct{}{}
				}()

				// @step: close the window
				shutdown := `
				<html>
				<body>
				<script>
					window.close()
				</script>
				</body>
				</html>`

				_, _ = w.Write([]byte(shutdown))
			})

			// @step: start the local loopback server
			go func() {
				if err := http.ListenAndServe(":3001", nil); err != nil {
					fmt.Fprintf(os.Stderr, "[error] failed to open the local http server on localhost:3000, error: %s\n", err)
					os.Exit(1)
				}
			}()

			// @step: open a brower to the to the api server
			url := fmt.Sprintf("%s/oauth/authorize?redirect_url=http://localhost:3001", config.Server)
			if err := open.Run(url); err != nil {
				fmt.Fprintf(os.Stderr, "[error] failed to open web brower to %s, error: %s", url, err)
				os.Exit(1)
			}

			// @step: we wait for either a done or error or timeout
			select {
			case <-doneCh:
			case err := <-errCh:
				fmt.Fprintf(os.Stderr, "[error] unable to authorize the client: %s", err)
				os.Exit(1)
			case <-time.After(60 * time.Second):
				fmt.Fprint(os.Stderr, "[error] authorization request has timed out")
				os.Exit(1)
			}

			config.AuthorizeURL = token.AuthorizationURL
			config.TokenURL = token.TokenEndpointURL
			config.Credentials.ClientID = token.ClientID
			config.Credentials.ClientSecret = token.ClientSecret
			config.Credentials.AccessToken = token.AccessToken
			config.Credentials.RefreshToken = token.RefreshToken
			config.Credentials.IDToken = token.IDToken

			// @step: update the configuration
			if err := config.Update(); err != nil {
				return fmt.Errorf("trying to update the client configuration: %s", err)
			}
			fmt.Println("You have successfully authenticated")

			return nil
		},
	}
}

//
func DecodeAuthorizationResponse(reader io.Reader) (*AuthorizationResponse, error) {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(string(content))
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
	resp := &AuthorizationResponse{}

	if err := json.NewDecoder(bytes.NewReader(rb)).Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
