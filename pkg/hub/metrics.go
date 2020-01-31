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

package hub

import "github.com/prometheus/client_golang/prometheus"

var (
	errorCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "hub_error_counter",
			Help: "A counter of the number of errors encountered by the hub api",
		},
	)
)

func init() {
	prometheus.MustRegister(errorCounter)
}
