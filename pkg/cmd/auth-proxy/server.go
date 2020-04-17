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
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/armon/go-proxyproto"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// authImpl implements the authentication proxy
type authImpl struct {
	sync.RWMutex
	logger          log.FieldLogger
	config          Config
	verifiers       []Verifier
	stopCh          chan struct{}
	token           string
	allowedNetworks []*net.IPNet
	addr            string
}

// New creates and returns a new authentication proxy
func New(
	logger log.FieldLogger,
	config Config,
	verifiers []Verifier,
) (Interface, error) {
	if err := config.IsValid(); err != nil {
		return nil, err
	}

	var allowedNetworks []*net.IPNet
	for _, cidr := range config.AllowedIPs {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, fmt.Errorf("invalid CIDR notation: %q", cidr)
		}
		allowedNetworks = append(allowedNetworks, network)
	}

	content, err := ioutil.ReadFile(config.UpstreamAuthorizationToken)
	if err != nil {
		return nil, err
	}
	token := strings.TrimSuffix(string(content), "\n")

	return &authImpl{
		logger:          logger,
		config:          config,
		stopCh:          make(chan struct{}),
		token:           token,
		verifiers:       verifiers,
		allowedNetworks: allowedNetworks,
	}, nil
}

// Run is called to start the proxy up
func (a *authImpl) Run(ctx context.Context) error {
	router := httprouter.New()
	origin, err := url.Parse(a.config.UpstreamURL)
	if err != nil {
		return err
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(origin)

	reverseProxy.Director = func(req *http.Request) {
		req.Header.Set("Host", origin.Host)
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Origin-Host", origin.Host)

		req.URL.Scheme = origin.Scheme
		req.URL.Host = origin.Host

		httpRequestCounter.Inc()
	}

	reverseProxy.ModifyResponse = func(resp *http.Response) error {
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			httpErrorCounter.Inc()
		}

		return nil
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
			// @step: handle a simple health check
			if req.URL.Path == "/ready" {
				resp.WriteHeader(http.StatusOK)
				_, _ = resp.Write([]byte("OK\n"))

				return
			}

			// @step: handle the ip filtering if required
			var remoteIP net.IP
			if host, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
				remoteIP = net.ParseIP(host)
			}
			if remoteIP == nil {
				a.logger.WithField("remote_address", req.RemoteAddr).Warnf("invalid remote address, access forbidden")

				return
			}
			if !a.isIPAllowed(remoteIP) {
				a.logger.WithField("remote_address", req.RemoteAddr).Warnf("access forbidden")
				resp.WriteHeader(http.StatusForbidden)
				_, _ = resp.Write([]byte("Forbidden\n"))

				return
			}

			// @step: attempt to verify the request against any of the verifiers
			allowed := func() bool {
				for _, x := range a.verifiers {
					matched, err := x.Admit(req)
					if err != nil {
						a.logger.WithError(err).Error("trying to verifier the inbound request")

						continue
					}
					if matched {
						return true
					}
				}

				return false
			}()
			if !allowed {
				authFailureCounter.Inc()

				a.logger.WithError(err).Error("trying to verify the inbound request")
				resp.WriteHeader(http.StatusForbidden)

				return
			}

			reverseProxy.ServeHTTP(resp, req)
		})
	}

	listener, err := net.Listen("tcp", a.config.Listen)
	if err != nil {
		return err
	}
	// @note: by enabling proxy protocol we wil have to change the health checks to tcp
	if a.config.EnableProxyProtocol {
		listener = &proxyproto.Listener{Listener: listener}
	}
	a.addr = listener.Addr().String()

	hs := &http.Server{Addr: listener.Addr().String(), Handler: router}

	go func() {
		a.logger.WithFields(log.Fields{
			"addr": hs.Addr,
		}).Info("starting the auth proxy service")

		switch a.config.HasTLS() {
		case true:
			if err := hs.ServeTLS(listener, a.config.TLSCert, a.config.TLSKey); err != nil && err != http.ErrServerClosed {
				a.logger.WithError(err).Fatal("trying to start the http server")
			}
		default:
			if err := hs.Serve(listener); err != nil && err != http.ErrServerClosed {
				a.logger.WithError(err).Fatal("trying to start the http server")
			}
		}
	}()

	msl, err := net.Listen("tcp", a.config.MetricsListen)
	if err != nil {
		return err
	}
	ms := &http.Server{Addr: msl.Addr().String(), Handler: promhttp.Handler()}

	go func() {
		a.logger.WithFields(log.Fields{
			"addr": ms.Addr,
		}).Info("starting the auth proxy metrics http server")

		if err := ms.Serve(msl); err != nil && err != http.ErrServerClosed {
			if err != http.ErrServerClosed {
				a.logger.WithError(err).Fatal("trying to start the metrics http server")
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

// Stop is called to halt the proxy
func (a *authImpl) Stop() error {
	a.logger.Info("shutting down the http services")
	a.stopCh <- struct{}{}

	return nil
}

func (a *authImpl) isIPAllowed(ip net.IP) bool {
	if ip == nil {
		return false
	}
	for _, network := range a.allowedNetworks {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

func (a *authImpl) Addr() string {
	return a.addr
}
