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

package authproxy

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/appvia/kore/pkg/cmd/auth-proxy/filters"
	"github.com/appvia/kore/pkg/cmd/auth-proxy/filters/authenticate"
	"github.com/appvia/kore/pkg/cmd/auth-proxy/filters/health"
	"github.com/appvia/kore/pkg/cmd/auth-proxy/filters/metrics"
	"github.com/appvia/kore/pkg/cmd/auth-proxy/filters/netfilter"
	"github.com/appvia/kore/pkg/cmd/auth-proxy/filters/proxy"
	"github.com/appvia/kore/pkg/cmd/auth-proxy/verifiers"

	"github.com/armon/go-proxyproto"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// authImpl implements the authentication proxy
type authImpl struct {
	*Config
	verifiers []verifiers.Interface
	stopCh    chan struct{}
	handler   *httprouter.Router
}

// New creates and returns a proxy
func New(config *Config, verifiers []verifiers.Interface) (Interface, error) {
	if err := config.IsValid(); err != nil {
		return nil, err
	}

	return &authImpl{
		Config:    config,
		stopCh:    make(chan struct{}),
		verifiers: verifiers,
		handler:   httprouter.New(),
	}, nil
}

// MakeRouter is responsible for creating the http router
func (a *authImpl) MakeRouter() error {
	// @step: create the http router and middleware filters
	middleware, err := a.makeFilters()
	if err != nil {
		return err
	}

	// @step: iterate all http methods and perform a catch all
	for _, x := range AllMethods {
		a.handler.Handle(x, "/*catchall", middleware.Wrap(nil))
	}

	return nil
}

// Run is called to start the proxy up
func (a *authImpl) Run(ctx context.Context) error {
	// @step we create the listener
	listener, err := net.Listen("tcp", a.Listen)
	if err != nil {
		return err
	}
	if a.EnableProxyProtocol {
		listener = &proxyproto.Listener{Listener: listener}
	}

	// @step: create the metrics listener
	ml, err := net.Listen("tcp", a.MetricsListen)
	if err != nil {
		return err
	}

	if err := a.MakeRouter(); err != nil {
		return err
	}

	hs := &http.Server{Addr: listener.Addr().String(), Handler: a.handler}
	ms := &http.Server{Addr: ml.Addr().String(), Handler: promhttp.Handler()}

	// @step: we start the http services
	go func() {
		log.WithFields(log.Fields{
			"addr": hs.Addr,
		}).Info("starting the auth proxy service")

		switch a.HasTLS() {
		case true:
			if err := hs.ServeTLS(listener, a.TLSCert, a.TLSKey); err != nil && err != http.ErrServerClosed {
				log.WithError(err).Fatal("trying to start the http server")
			}
		default:
			if err := hs.Serve(listener); err != nil && err != http.ErrServerClosed {
				log.WithError(err).Fatal("trying to start the http server")
			}
		}
	}()

	go func() {
		log.WithFields(log.Fields{
			"addr": ms.Addr,
		}).Info("starting the auth proxy metrics http server")

		if err := ms.Serve(ml); err != nil && err != http.ErrServerClosed {
			if err != http.ErrServerClosed {
				log.WithError(err).Fatal("trying to start the metrics http server")
			}
		}
	}()

	go func() {
		<-a.stopCh
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_ = hs.Shutdown(ctx)
		_ = ms.Shutdown(ctx)
	}()

	return nil
}

//

// makeFilters is responsible for generating the filters
func (a *authImpl) makeFilters() (filters.Interface, error) {
	filter := filters.New()

	ip, err := netfilter.New(netfilter.Options{
		Permitted: a.AllowedIPs,
	})
	if err != nil {
		return nil, err
	}
	auth, err := authenticate.New(authenticate.Options{
		Verifiers: a.verifiers,
	})
	if err != nil {
		return nil, err
	}
	proxy, err := proxy.New(proxy.Options{
		Endpoint:      a.UpstreamURL,
		FlushInterval: a.FlushInterval,
	})
	if err != nil {
		return nil, err
	}
	filter.Use(health.New(), metrics.New(), ip, auth, proxy)

	return filter, nil
}

// Stop is called to halt the proxy
func (a *authImpl) Stop() error {
	log.Info("shutting down the http services")
	a.stopCh <- struct{}{}

	return nil
}
