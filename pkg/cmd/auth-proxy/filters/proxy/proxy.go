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

package proxy

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/appvia/kore/pkg/cmd/auth-proxy/filters"
	"github.com/appvia/kore/pkg/version"

	log "github.com/sirupsen/logrus"
)

// Options are the configurable for the filter
type Options struct {
	// Endpoint is the destination to proxy
	Endpoint string
	// FlushInterval is the flush interval for reverse proxy
	FlushInterval time.Duration
}

type pxyImpl struct {
	// upstream is the reverse proxy
	upstream *httputil.ReverseProxy
}

// Serve is responsible for implementing the filter
func (p *pxyImpl) Serve(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("X-Proxy-Version", version.Release)

		p.upstream.ServeHTTP(w, req)
	})
}

// New creates and returns a reverse proxy
func New(options Options) (filters.Middleware, error) {
	if options.Endpoint == "" {
		return nil, errors.New("no endpoint")
	}

	log.WithFields(log.Fields{
		"endpoint": options.Endpoint,
		"flush":    options.FlushInterval,
	}).Debug("using the endpoint reverse proxy")

	origin, err := url.Parse(options.Endpoint)
	if err != nil {
		return nil, err
	}
	if origin.Host == "" || origin.Scheme == "" {
		return nil, errors.New("invalid url")
	}

	rv := httputil.NewSingleHostReverseProxy(origin)
	rv.Director = func(req *http.Request) {
		req.Header.Set("Host", origin.Host)
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Origin-Host", origin.Host)

		req.URL.Scheme = origin.Scheme
		req.URL.Host = origin.Host
	}
	rv.FlushInterval = options.FlushInterval
	rv.ModifyResponse = func(resp *http.Response) error {
		return nil
	}

	rv.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		},
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
	}

	return &pxyImpl{upstream: rv}, nil
}
