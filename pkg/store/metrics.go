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

package store

import "github.com/prometheus/client_golang/prometheus"

var (
	cacheHitCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "store_cache_hit_total",
			Help: "The total number of operations which have been served from cache",
		},
	)
	cacheMissCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "store_cache_miss_total",
			Help: "The total number of operations which missed the cache",
		},
	)
	createCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "store_create_counter",
			Help: "A counter or the create operations in the store",
		},
	)
	deleteCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "store_delete_counter",
			Help: "A counter or the delete operations in the store",
		},
	)
	deleteLatency = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "store_delete_latency_sec",
			Help: "The latency on delete operations to the store",
		},
	)
	errorCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "store_error_counter",
			Help: "A counter of the number of errors encountered by the store",
		},
	)
	getLatency = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "store_get_latency_sec",
			Help: "The latency on get operations to the store",
		},
	)
	setLatency = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "store_set_latency_sec",
			Help: "The latency on set operations to the store",
		},
	)
	updateCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "store_update_counter",
			Help: "A counter or the update and add operations in the store",
		},
	)
)

func init() {
	prometheus.MustRegister(cacheHitCounter)
	prometheus.MustRegister(cacheMissCounter)
	prometheus.MustRegister(createCounter)
	prometheus.MustRegister(deleteCounter)
	prometheus.MustRegister(deleteLatency)
	prometheus.MustRegister(errorCounter)
	prometheus.MustRegister(getLatency)
	prometheus.MustRegister(setLatency)
	prometheus.MustRegister(updateCounter)
}
