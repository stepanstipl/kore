/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
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

// @TODO Quick and dirty - the while this needs to be polished up

// Register is responsible for handling the registration
func (l *loginHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	log.WithFields(log.Fields{
		"path": "login",
	}).Info("registering the login webservice with container")

	l.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path("oauth")

	if l.Config().DiscoveryURL != "" {
		provider, err := oidc.NewProvider(context.Background(), l.Config().DiscoveryURL)
		if err != nil {
			log.WithError(err).Error("failed to create the openid provider")

			return nil, err
		}
		l.provider = provider

		// Configure an OpenID Connect aware OAuth2 client.
		l.oidcConfig = &oauth2.Config{
			ClientID:     l.Config().ClientID,
			ClientSecret: l.Config().ClientSecret,
			Endpoint:     provider.Endpoint(),
			Scopes:       append([]string{oidc.ScopeOpenID}, l.Config().ClientScopes...),
		}

		l.verifier = provider.Verifier(&oidc.Config{ClientID: l.Config().ClientID})
	} else {
		log.Warn("no identity provider configuration has been provided")
	}

	ws.Route(
		ws.GET("/authorize").To(l.authorizerHandler).
			Param(ws.QueryParameter("redirect_url", "The rediection url, i.e. the location to redirect post").Required(true)).
			Doc("Used to start the authorization flow for user authentication").
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/callback").To(l.callbackHandler).
			Doc("Used to handle the authorization callback from the identity provider").
			Param(ws.QueryParameter("code", "The authorization code returned from the identity provider").Required(true)).
			Param(ws.QueryParameter("state", "The state parameter which was passed on authorization request").Required(true)).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
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
	redirectURL := l.Config().PublicAPIURL
	if redirectURL == "" {
		scheme := "http"
		if req.Request.TLS != nil || req.Request.Header.Get("X-Forward-Proto") == "https" {
			scheme = "https"
		}
		redirectURL = fmt.Sprintf("%s://%s", scheme, req.Request.Host)
	}
	l.oidcConfig.RedirectURL = fmt.Sprintf("%s/oauth/callback", redirectURL)

	// @step: redirect the user to perform the login flow
	log.WithFields(log.Fields{
		"client_ip":    req.Request.RemoteAddr,
		"redirect_url": l.oidcConfig.RedirectURL,
		"scopes":       l.Config().ClientScopes,
	}).Info("providing authorization redirect to identity service")

	http.Redirect(resp.ResponseWriter, req.Request, l.oidcConfig.AuthCodeURL(state), http.StatusTemporaryRedirect)
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
		username, found := claims.GetUserClaim(l.Config().UserClaims...)
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
			AuthorizationURL: l.Config().DiscoveryURL,
			ClientID:         l.Config().ClientID,
			ClientSecret:     l.Config().ClientSecret,
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
