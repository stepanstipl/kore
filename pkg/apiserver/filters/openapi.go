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
