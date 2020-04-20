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

package netfilter

import (
	"fmt"
	"net"
	"net/http"

	"github.com/appvia/kore/pkg/cmd/auth-proxy/filters"

	log "github.com/sirupsen/logrus"
)

// Options are the configurable for the filter
type Options struct {
	// Permitted is a collection of network ranges permitted
	Permitted []string
}

type filterImpl struct {
	Options
	ranges []*net.IPNet
}

// New creates and returns a filter for network ranges
func New(options Options) (filters.Middleware, error) {
	var list []*net.IPNet

	for _, x := range options.Permitted {
		_, network, err := net.ParseCIDR(x)
		if err != nil {
			return nil, fmt.Errorf("invalid cidr notation: %q", x)
		}
		list = append(list, network)
	}

	return &filterImpl{Options: options, ranges: list}, nil
}

// Serve handles the incoming request
func (f *filterImpl) Serve(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if len(f.ranges) > 0 && !f.isAllowed(req) {
			log.WithField(
				"address", req.RemoteAddr,
			).Debug("denying the request from address")

			w.WriteHeader(http.StatusForbidden)

			return
		}

		next.ServeHTTP(w, req)
	})
}

// isAllowed iterates and checks the address against the permitted ranges
func (f *filterImpl) isAllowed(req *http.Request) bool {
	ip := f.remoteAddress(req)
	if ip == nil {
		return false
	}

	for _, network := range f.ranges {
		if network.Contains(ip) {
			return true
		}
	}

	return false
}

// remoteAddress return the remote address from the request
func (f *filterImpl) remoteAddress(req *http.Request) net.IP {
	if host, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		return net.ParseIP(host)
	}

	return nil
}
