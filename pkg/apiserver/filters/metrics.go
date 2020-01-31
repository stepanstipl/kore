/*
 * Copyright (C) 2019  Appvia Ltd <info@appvia.io>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package filters

import (
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
)

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
	}()

	chain.ProcessFilter(req, resp)
}
