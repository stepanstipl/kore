/*
 * Copyright (C) 2019  Rohith Jayawardene <gambol99@gmail.com>
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

package apiserver

import "github.com/prometheus/client_golang/prometheus"

var (
	apiErrorCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "api_error_total",
			Help: "The total amount of errors encountered by the api layer",
		},
	)
	apiUsersCounters = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_users_requests_total",
			Help: "The total amount of requests to the users service broken down by method",
		},
		[]string{"operation"},
	)
)

func init() {
	prometheus.MustRegister(apiUsersCounters)
	prometheus.MustRegister(apiErrorCounter)
}
