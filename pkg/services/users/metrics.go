/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 * 
 * This file is part of hub-apiserver.
 * 
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * 
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * 
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
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