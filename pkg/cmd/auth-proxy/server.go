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

package authproxy

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/appvia/kore/pkg/utils"
	"github.com/coreos/go-oidc"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// authImpl implements the authentication proxy
type authImpl struct {
	sync.RWMutex
	// config is the configuration for the service
	config Config
	// verifier is the rsa
	verifier *oidc.IDTokenVerifier
	// stopCh is the stop channel
	stopCh chan struct{}
	// token is the upstream token
	token string
}

// New creates and returns a new authentication proxy
func New(config Config) (Interface, error) {
	if err := config.IsValid(); err != nil {
		return nil, err
	}

	var verifier *oidc.IDTokenVerifier

	options := &oidc.Config{
		ClientID:          config.ClientID,
		SkipClientIDCheck: true,
		SkipExpiryCheck:   false,
	}

	content, err := ioutil.ReadFile(config.UpstreamAuthorizationToken)
	if err != nil {
		return nil, err
	}
	token := strings.TrimSuffix(string(content), "\n")

	// @step: do we have a signing ca?
	if config.SigningCA != "" {
		log.WithField(
			"signing_ca", config.SigningCA,
		).Info("using the signing certificate to verifier the requests")

		keyset, err := newStaticKeySet(config.SigningCA)
		if err != nil {
			return nil, err
		}

		verifier = oidc.NewVerifier(config.ClientID, keyset, options)
	}
	if config.DiscoveryURL != "" {
		log.WithField(
			"discovery-url", config.DiscoveryURL,
		).Info("using the discovery endpoint to verifier the requests")

		provider, err := oidc.NewProvider(context.Background(), config.DiscoveryURL)
		if err != nil {
			log.WithError(err).Error("trying to retrieve provider details")

			return nil, err
		}

		verifier = provider.Verifier(options)
	}

	return &authImpl{
		config:   config,
		stopCh:   make(chan struct{}),
		token:    string(token),
		verifier: verifier,
	}, nil
}

// Run is called to start the proxy up
func (a *authImpl) Run(ctx context.Context) error {
	// @step: start the http service
	router := httprouter.New()
	origin, err := url.Parse(a.config.UpstreamURL)
	if err != nil {
		return err
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(origin)

	reverseProxy.Director = func(req *http.Request) {
		user := req.Context().Value(Key{}).(string)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.token))
		req.Header.Add("Impersonate-User", user)
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)

		req.URL.Scheme = origin.Scheme
		req.URL.Host = origin.Host
	}
	reverseProxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		},
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
	}

	for _, method := range AllMethods {
		router.Handle(method, "/*catchall", func(resp http.ResponseWriter, req *http.Request, p httprouter.Params) {
			var username string

			err := func() error {
				// @step: extract the token from the request
				bearer, found := utils.GetBearerToken(req.Header.Get("Authorization"))
				if !found {
					return errors.New("no authorization token")
				}

				// @step: parse and extract the identity
				raw, err := a.verifier.Verify(req.Context(), bearer)
				if err != nil {
					return err
				}

				// @step: extract the username if any
				claims, err := utils.NewClaimsFromToken(raw)
				if err != nil {
					return err
				}

				user, found := claims.GetUserClaim(a.config.UserClaims...)
				if !found {
					return errors.New("no username found in the identity token")
				}
				username = user

				return nil
			}()
			if err != nil {
				log.WithError(err).Error("trying to verify the inbound request")
				resp.WriteHeader(http.StatusForbidden)

				return
			}

			// @step: add the username to the context
			ctx := context.WithValue(req.Context(), Key{}, username)

			reverseProxy.ServeHTTP(resp, req.WithContext(ctx))
		})
	}

	hs := &http.Server{Addr: a.config.Listen, Handler: router}

	go func() {
		log.WithField(
			"listen", a.config.Listen,
		).Info("starting the auth proxy service")

		switch a.config.HasTLS() {
		case true:
			if err := hs.ListenAndServeTLS(a.config.TLSCert, a.config.TLSKey); err != nil {
				log.WithError(err).Fatal("trying to start the http server")
			}
		default:
			if err := hs.ListenAndServe(); err != nil {
				log.WithError(err).Fatal("trying to start the http server")
			}
		}
	}()

	go func() {
		<-a.stopCh
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_ = hs.Shutdown(ctx)
	}()

	return nil
}

// Stop is called to halt the proxy
func (a *authImpl) Stop() error {
	log.Info("shutting down the http services")
	a.stopCh <- struct{}{}

	return nil
}
