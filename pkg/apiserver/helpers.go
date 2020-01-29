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

import (
	"io"
	"net/http"
	"strings"

	"github.com/appvia/kore/pkg/hub"
	log "github.com/sirupsen/logrus"

	restful "github.com/emicklei/go-restful"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
func returnNotImplemented(req *restful.Request, wr *restful.Response) {
	wr.WriteHeader(http.StatusNotImplemented)
}
*/

// newList provides an api list type
func newList() *List {
	return &List{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "List",
		},
	}
}

func makeListWithSize(size int) *List {
	l := newList()
	l.Items = make([]string, size)

	return l
}

// handleErrors is a generic wrapper for handling the error from downstream hub brigde
func handleErrors(req *restful.Request, resp *restful.Response, handler func() error) {
	if err := handler(); err != nil {
		code := http.StatusInternalServerError
		switch err {
		case hub.ErrNotFound:
			code = http.StatusNotFound
		case hub.ErrNotAllowed{}:
			code = http.StatusNotAcceptable
		case hub.ErrUnauthorized:
			code = http.StatusForbidden
		case hub.ErrRequestInvalid:
			code = http.StatusBadRequest
		case io.EOF:
			code = http.StatusBadRequest
		}
		if strings.Contains(err.Error(), "record not found") {
			code = http.StatusNotFound
		}

		e := newError(err.Error()).
			WithCode(code).
			WithVerb(req.Request.Method).
			WithURI(req.Request.RequestURI)

		if err := resp.WriteHeaderAndEntity(code, e); err != nil {
			log.WithError(err).Error("failed to respond to request")
		}
	}
}
