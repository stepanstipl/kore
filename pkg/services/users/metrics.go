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

package users

import "github.com/prometheus/client_golang/prometheus"

var (
	createCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "db_create_counter",
			Help: "A counter or the create operations in the db",
		},
	)
	deleteCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "db_delete_counter",
			Help: "A counter or the delete operations in the db",
		},
	)
	deleteLatency = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "db_delete_latency_sec",
			Help: "The latency on delete operations to the db",
		},
	)
	errorCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "db_error_counter",
			Help: "A counter of the number of errors encountered by the db",
		},
	)
	getLatency = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "db_get_latency_sec",
			Help: "The latency on get operations to the db",
		},
	)
	listLatency = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "db_list_latency_sec",
			Help: "The latency on list operations to the db",
		},
	)
	setLatency = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "db_set_latency_sec",
			Help: "The latency on set operations to the db",
		},
	)
	updateCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "db_update_counter",
			Help: "A counter or the update and add operations in the db",
		},
	)
)

func init() {
	prometheus.MustRegister(createCounter)
	prometheus.MustRegister(deleteCounter)
	prometheus.MustRegister(deleteLatency)
	prometheus.MustRegister(errorCounter)
	prometheus.MustRegister(getLatency)
	prometheus.MustRegister(listLatency)
	prometheus.MustRegister(setLatency)
	prometheus.MustRegister(updateCounter)
}
