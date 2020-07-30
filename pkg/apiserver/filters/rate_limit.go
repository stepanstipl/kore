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
	"net/http"
	"strconv"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateConfig is the configuration
type RateConfig struct {
	// Period is the limiting period
	Period time.Duration
	// Limit is the limit per period
	Limit int64
}

// NewRateLimiter returns a rate limiting middleware using a memory store
func NewRateLimiter(limit RateConfig) restful.FilterFunction {
	store := memory.NewStore()
	throttle := limiter.New(store, limiter.Rate{
		Period: limit.Period,
		Limit:  limit.Limit,
	})

	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		limiterCtx, err := throttle.Get(req.Request.Context(), req.Request.RemoteAddr)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)

			return
		}

		h := resp.Header()
		h.Set("X-RateLimit-Limit", strconv.FormatInt(limiterCtx.Limit, 10))
		h.Set("X-RateLimit-Remaining", strconv.FormatInt(limiterCtx.Remaining, 10))
		h.Set("X-RateLimit-Reset", strconv.FormatInt(limiterCtx.Reset, 10))

		if limiterCtx.Reached {
			resp.WriteHeader(http.StatusTooManyRequests)

			return
		}

		chain.ProcessFilter(req, resp)
	}
}
