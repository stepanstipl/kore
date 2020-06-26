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

package metrics

import (
	"fmt"
	"net/http"

	"github.com/appvia/kore/pkg/cmd/auth-proxy/filters"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/felixge/httpsnoop"
)

var (
	httpRequestCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "The total number http requests processed",
		},
	)
	httpCodeCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_code",
			Help: "The total number of errors in requests",
		},
		[]string{"code"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestCounter)
	prometheus.MustRegister(httpCodeCounter)
}

type metricsImpl struct{}

// New creates and returns a metric middleware
func New() filters.Middleware {
	return &metricsImpl{}
}

func (m *metricsImpl) Serve(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		httpRequestCounter.Inc()

		m := httpsnoop.CaptureMetrics(next, w, req)
		httpCodeCounter.WithLabelValues(fmt.Sprintf("%d", m.Code)).Inc()
	})
}
