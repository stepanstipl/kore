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

package filters

import (
	"time"

	"github.com/appvia/kore/pkg/hub/authentication"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

var (
	// DefaultLogging is the default logging filter
	DefaultLogging = Logging{}
)

// Logging provides a generic http logger
type Logging struct{}

// Filter is a logging filter for the api server
func (l *Logging) Filter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	start := time.Now()

	user := req.Request.Context().Value(authentication.ContextKey{})

	defer func() {
		fields := log.Fields{
			"code":   resp.StatusCode(),
			"ip":     req.Request.RemoteAddr,
			"method": req.Request.Method,
			"query":  req.Request.URL.RawQuery,
			"time":   time.Since(start).String(),
			"uri":    req.Request.RequestURI,
		}
		if user != nil {
			fields["user"] = user.(authentication.Identity).Username()
		}

		log.WithFields(fields).Info("logging incoming http request")
	}()

	chain.ProcessFilter(req, resp)
}
