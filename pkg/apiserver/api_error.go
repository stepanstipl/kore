/*
 * Copyright (C) 2019 Rohith Jayawardene <gambol99@gmail.com>
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

import "net/http"

// newError returns a new api error
func newError(message string) *Error {
	return &Error{
		Message: message,
		Code:    http.StatusInternalServerError,
	}
}

// WithVerb sets the request verb
func (e *Error) WithVerb(v string) *Error {
	e.Verb = v

	return e
}

// WithURI sets the request uri
func (e *Error) WithURI(v string) *Error {
	e.URI = v

	return e
}

// WithDetail addes a detailed message
func (e *Error) WithDetail(v string) *Error {
	e.Detail = v

	return e
}

// WithCode adds a code to the api error
func (e *Error) WithCode(v int) *Error {
	e.Code = v

	return e
}
