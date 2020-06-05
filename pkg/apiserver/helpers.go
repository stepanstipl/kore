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

package apiserver

import (
	"errors"
	"io"
	"net/http"

	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/validation"

	restful "github.com/emicklei/go-restful"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
func returnNotImplemented(req *restful.Request, wr *restful.Response) {
	wr.WriteHeader(http.StatusNotImplemented)
}
*/

// withStandardErrors adds the standard internal server error (500) result to the route.
func withStandardErrors(rb *restful.RouteBuilder) *restful.RouteBuilder {
	return rb.
		Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{})
}

// withValidationErrors adds the standard bad request (400) validation error result to the route.
func withValidationErrors(rb *restful.RouteBuilder) *restful.RouteBuilder {
	return rb.
		Returns(http.StatusBadRequest, "Validation error of supplied parameters/body", validation.Error{})
}

// withAuthErrors adds the standard unauthenticated (401) and forbidden (403) results to the route.
func withAuthErrors(rb *restful.RouteBuilder) *restful.RouteBuilder {
	return rb.
		Returns(http.StatusUnauthorized, "If not authenticated", nil).
		Returns(http.StatusForbidden, "If authenticated but not authorized", nil)
}

// withAllErrors is a shorthand to add all standard, validation, and auth results to the route.
func withAllErrors(rb *restful.RouteBuilder) *restful.RouteBuilder {
	return withValidationErrors(withAuthErrors(withStandardErrors(rb)))
}

// withAllNonValidationErrors is a shorthand to add all standard and auth results to the route but not validation.
func withAllNonValidationErrors(rb *restful.RouteBuilder) *restful.RouteBuilder {
	return withAuthErrors(withStandardErrors(rb))
}

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

// handleErrors is a generic wrapper for handling the error from downstream kore brigde
func handleErrors(req *restful.Request, resp *restful.Response, handler func() error) {
	if err := handler(); err != nil {
		handleError(req, resp, err)
	}
}

func handleError(req *restful.Request, resp *restful.Response, err error) {
	code := http.StatusInternalServerError
	// Simple errors have fixed values:
	switch err {
	case kore.ErrNotFound:
		code = http.StatusNotFound
	case kore.ErrUnauthorized:
		code = http.StatusForbidden
	case kore.ErrRequestInvalid:
		code = http.StatusBadRequest
	case io.EOF:
		code = http.StatusBadRequest
	}

	// Couple of errors have their own types, treat differently:
	switch err.(type) {
	case kore.ErrNotAllowed, *kore.ErrNotAllowed:
		code = http.StatusForbidden
	case validation.Error, *validation.Error:
		code = http.StatusBadRequest
	}

	if err == gorm.ErrRecordNotFound {
		code = http.StatusNotFound
		err = errors.New("resource not found")
	}
	if kerrors.IsNotFound(err) {
		code = http.StatusNotFound
	}

	writeError(req, resp, err, code)
}

func writeError(req *restful.Request, resp *restful.Response, err error, code int) {
	switch err.(type) {
	case validation.Error, *validation.Error:
		// Error can be directly serialized to json so just return that.
	default:
		err = newError(err.Error()).
			WithCode(code).
			WithVerb(req.Request.Method).
			WithURI(req.Request.RequestURI)
	}

	if err := resp.WriteHeaderAndEntity(code, err); err != nil {
		log.WithError(err).Error("failed to respond to request")
	}
}
