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

package apiserver

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	"github.com/coreos/go-oidc"
	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func init() {
	RegisterHandler(&loginHandler{})
}

type loginHandler struct {
	kore.Interface
	// oidcConfig is the openid configuration
	oidcConfig *oauth2.Config
	// verifier is responsible for verification of the tokens
	verifier *oidc.IDTokenVerifier
	// provider is the oidc provider
	provider *oidc.Provider
	// default handler
	DefaultHandler
}

const googleAuthURL = "https://accounts.google.com"

// Path returns the handler path
func (l *loginHandler) Path() string {
	return "oauth"
}

// Register is responsible for handling the registration
func (l *loginHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	log.WithFields(log.Fields{
		"path": l.Path(),
	}).Info("registering the login webservice with container")

	l.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(l.Path())

	if l.Config().HasOpenID() {
		provider, err := oidc.NewProvider(context.Background(), l.Config().IDPServerURL)
		if err != nil {
			log.WithError(err).Error("failed to create the openid provider")

			return nil, err
		}
		l.provider = provider

		redirect := ""
		if l.Config().PublicAPIURL != "" {
			redirect = fmt.Sprintf("%s/oauth/callback", l.Config().PublicAPIURL)
		}

		// Configure an OpenID Connect aware OAuth2 client.
		scopes := append([]string{oidc.ScopeOpenID}, l.Config().IDPClientScopes...)
		if strings.HasPrefix(l.Config().IDPServerURL, googleAuthURL) {
			// Remove offline scope if present as it causes failures on google.
			scopes = func() (ret []string) {
				for _, s := range scopes {
					if s != oidc.ScopeOfflineAccess {
						ret = append(ret, s)
					}
				}
				return
			}()
		}
		l.oidcConfig = &oauth2.Config{
			ClientID:     l.Config().IDPClientID,
			ClientSecret: l.Config().IDPClientSecret,
			RedirectURL:  redirect,
			Endpoint:     provider.Endpoint(),
			Scopes:       scopes,
		}

		l.verifier = provider.Verifier(&oidc.Config{ClientID: l.Config().IDPClientID})
	} else {
		log.Warn("no identity provider configuration has been provided")
	}

	ws.Route(
		ws.GET("/authorize").To(l.authorizerHandler).
			Doc("Used to start the authorization flow for user authentication").
			Operation("LoginAttempted").
			Param(ws.QueryParameter("redirect_url", "The rediection url, i.e. the location to redirect post").Required(true)).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/callback").To(l.callbackHandler).
			Doc("Used to handle the authorization callback from the identity provider").
			Operation("LoginCallback").
			Param(ws.QueryParameter("code", "The authorization code returned from the identity provider").Required(true)).
			Param(ws.QueryParameter("state", "The state parameter which was passed on authorization request").Required(true)).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	return ws, nil
}

// authorizerHandler is responsible for authorizing a client
func (l *loginHandler) authorizerHandler(req *restful.Request, resp *restful.Response) {
	// @step: check if the handler has been configured
	if l.oidcConfig == nil {
		resp.WriteHeader(http.StatusNotImplemented)

		return
	}

	// session in order to process on the callback
	if req.QueryParameter("redirect_url") == "" {
		resp.WriteHeader(http.StatusBadRequest)

		return
	}
	state := base64.StdEncoding.EncodeToString([]byte(req.QueryParameter("redirect_url")))

	// @step: we either taken the public url or the host header
	if l.Config().PublicAPIURL == "" {
		l.oidcConfig.RedirectURL = fmt.Sprintf("%s/%s/oauth/callback", l.Path(), l.makeHostURL(req))
	}

	// @step: redirect the user to perform the login flow
	log.WithFields(log.Fields{
		"client_ip":    req.Request.RemoteAddr,
		"redirect_url": l.oidcConfig.RedirectURL,
		"scopes":       strings.Join(l.Config().IDPClientScopes, ","),
	}).Info("providing authorization redirect to identity service")

	var authCodeURL string
	if strings.HasPrefix(l.Config().IDPServerURL, googleAuthURL) {
		// Google has different ideas about how to request a refresh token.
		authCodeURL = l.oidcConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	} else {
		authCodeURL = l.oidcConfig.AuthCodeURL(state)
	}

	http.Redirect(resp.ResponseWriter, req.Request, authCodeURL, http.StatusTemporaryRedirect)
}

// callbackHandler is responsible for handling the callback
func (l *loginHandler) callbackHandler(req *restful.Request, resp *restful.Response) {
	// @step: check if the handler has been configured
	if l.oidcConfig == nil {
		resp.WriteHeader(http.StatusNotImplemented)

		return
	}

	// @step: retrieve the parameters
	state := req.QueryParameter("state")
	code := req.QueryParameter("code")
	if state == "" || code == "" {
		resp.WriteHeader(http.StatusBadRequest)

		return
	}

	// @step: decode the state - which holds the redirect_url
	redirect, err := base64.StdEncoding.DecodeString(state)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)

		return
	}

	valid, err := validateRedirectURL(string(redirect))
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)

		return
	} else if !valid {
		resp.WriteHeader(http.StatusForbidden)

		return
	}

	// @step: we either taken the public url or the host header
	if l.Config().PublicAPIURL == "" {
		l.oidcConfig.RedirectURL = fmt.Sprintf("%s/oauth/callback", l.makeHostURL(req))
	}

	hcode, err := func() (int, error) {
		// @step: exchange the authorization code with the identity provider
		otoken, err := l.oidcConfig.Exchange(req.Request.Context(), code)
		if err != nil {
			log.WithError(err).Error("trying to exchange the token to the idp")

			return http.StatusServiceUnavailable, errors.New("exchanging authorization code")
		}

		// @step: extract the id token from oauth2 token
		rawToken, ok := otoken.Extra("id_token").(string)
		if !ok {
			return http.StatusServiceUnavailable, errors.New("response is missing identity token")
		}

		// @step: parse and verify id token payload
		token, err := l.verifier.Verify(req.Request.Context(), rawToken)
		if err != nil {
			log.WithError(err).Error("trying to verify the token from the idp")

			return http.StatusForbidden, errors.New("token failed verification")
		}

		claims, err := utils.NewClaimsFromToken(token)
		if err != nil {
			return http.StatusForbidden, errors.New("invalid id token exchanged")
		}

		// @step: extract the username and email address
		username, found := claims.GetUserClaim(l.Config().IDPUserClaims...)
		if !found {
			return http.StatusForbidden, errors.New("no user information found in token ")
		}

		email, found := claims.GetEmail()
		if !found {
			return http.StatusForbidden, errors.New("no email found in the token")
		}

		// @step: ensure this matches a user in the kore - else we create him
		if err := l.Users().EnableUser(req.Request.Context(), username, email); err != nil {
			return http.StatusInternalServerError, err
		}

		// @step: build an authorization response and hold in the cache
		res := &AuthorizationResponse{
			AccessToken:      otoken.AccessToken,
			AuthorizationURL: l.Config().IDPServerURL,
			ClientID:         l.Config().IDPClientID,
			ClientSecret:     l.Config().IDPClientSecret,
			IDToken:          rawToken,
			RefreshToken:     otoken.RefreshToken,
			TokenEndpointURL: l.provider.Endpoint().TokenURL,
		}

		encoded, err := utils.EncodeToJSON(res)
		if err != nil {
			return http.StatusInternalServerError, errors.New("encoding the authorization response")
		}

		compressed := &bytes.Buffer{}
		wr := gzip.NewWriter(compressed)
		_, _ = wr.Write(encoded)
		wr.Close()

		redirectURL := fmt.Sprintf("%s?token=%s", redirect,
			base64.StdEncoding.EncodeToString(compressed.Bytes()),
		)

		http.Redirect(resp.ResponseWriter, req.Request, redirectURL, http.StatusTemporaryRedirect)

		return http.StatusOK, nil
	}()
	if err != nil {
		log.WithError(err).Error("failed to authorize the user to kore")

		resp.WriteHeader(hcode)
	}
}

// EnableAdminsOnly indicates if we need to be an admin user
func (l *loginHandler) EnableAdminsOnly() bool {
	return false
}

// EnableAuthentication indicates if this service needs auth
func (l *loginHandler) EnableAuthentication() bool {
	return false
}

// Name returns the name of the handler
func (l *loginHandler) Name() string {
	return "login"
}

// validateRedirectURL ensures that a given redirect_url points to localhost
func validateRedirectURL(redirectURL string) (bool, error) {
	validRedirectHosts := [...]string{"127.0.0.1", "localhost"}

	u, err := url.Parse(redirectURL)
	if err != nil {
		return false, err
	}

	for _, host := range validRedirectHosts {
		if u.Hostname() == host {
			return true, nil
		}
	}

	return false, nil
}

// makeHostURL is used to retrieve the host url from the host headers
func (l *loginHandler) makeHostURL(req *restful.Request) string {
	scheme := "http"
	if req.Request.TLS != nil || req.Request.Header.Get("X-Forward-Proto") == "https" {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s", scheme, req.Request.Host)
}
