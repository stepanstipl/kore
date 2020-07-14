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

package filters

import (
	"fmt"

	restful "github.com/emicklei/go-restful"
	"github.com/prometheus/client_golang/prometheus"
)

// DefaultMetrics is the default metrics filter
var DefaultMetrics = Metrics{}

// Metrics provides metrics for the api server
type Metrics struct{}

var (
	httpRequestAverageHistogram = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "http_request_avg_sec",
			Help: "The average latency on requests to the apiserver",
		},
	)
	httpRequestCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "The total number of http request to the apiserver",
		},
	)
	httpRequestErrorCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_request_error_total",
			Help: "The total number of http requests which have not been successful",
		},
	)
	httpCodeCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_code_total",
			Help: "The total number of http broken down by http code",
		},
		[]string{"code"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestAverageHistogram)
	prometheus.MustRegister(httpRequestCounter)
	prometheus.MustRegister(httpRequestErrorCounter)
	prometheus.MustRegister(httpCodeCounter)
}

// Filter is a logging filter for the api server
func (l Metrics) Filter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	timed := prometheus.NewTimer(httpRequestAverageHistogram)

	defer func() {
		timed.ObserveDuration()
		// @step: bump the counters
		httpRequestCounter.Inc()
		// @step: was the response an error
		if resp.StatusCode() < 200 || resp.StatusCode() > 399 {
			httpRequestErrorCounter.Inc()
		}
		httpCodeCounter.WithLabelValues(fmt.Sprintf("%d", resp.StatusCode())).Inc()
	}()

	chain.ProcessFilter(req, resp)
}
