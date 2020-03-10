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
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"net/http"
	"net/http/httptest"

	restful "github.com/emicklei/go-restful"
)

var (
	// SwaggerChecksum provides a middleware to checksum the swagger
	SwaggerChecksum swaggerChecksum
)

// swaggerChecksum is a hack to provide a checksum on the swagger
type swaggerChecksum struct{}

// Filter is a logging filter for the api server
func (l *swaggerChecksum) Filter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	// @step: we use a recorder to get the content of the swagger.json
	recorder := httptest.NewRecorder()

	chain.ProcessFilter(req, restful.NewResponse(recorder))

	algor := req.QueryParameter("checksum")
	if algor == "" {
		resp.WriteHeader(http.StatusOK)
		_, _ = resp.Write(recorder.Body.Bytes())

		return
	}

	// @step: read the content back from the swagger
	err := func() error {
		var h hash.Hash

		switch algor {
		case "md5":
			h = md5.New()
		default:
			h = sha256.New()
		}
		if _, err := h.Write(recorder.Body.Bytes()); err != nil {
			return err
		}
		_, err := resp.Write([]byte(hex.EncodeToString(h.Sum(nil))))

		return err
	}()
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)

		return
	}
}
