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
	"time"

	"github.com/appvia/kore/pkg/kore/authentication"

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
