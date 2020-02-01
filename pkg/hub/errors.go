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

import "errors"

var (
	// ErrOnRevision indicates the object is older than the revision in source.
	// They need to retrieve the latest copy and retry
	ErrOnRevision = errors.New("object revision is too old")
	// ErrNotFound the object has not been found i.e. doesn't exist
	ErrNotFound = errors.New("requested object does not exist")
	// ErrUnauthorized indicates we don't have permission to the object
	ErrUnauthorized = errors.New("object forbidden, require additional permissions")
	// ErrRequestInvalid indicates an invalid request object
	ErrRequestInvalid = errors.New("invalid request")
	// ErrServerData indicates invalid data
	ErrServerData = errors.New("invalid server data")
	// ErrServerNotAvailable indicates that the remote server is unavailable
	ErrServerNotAvailable = errors.New("a remote server is unavailable")
)

// NewErrNotAllowed returns an new not allowed error
func NewErrNotAllowed(message string) error {
	return ErrNotAllowed{message: message}
}

// ErrNotAllowed indicates an error of not allowed
type ErrNotAllowed struct {
	message string
}

// Error returns the message
func (e ErrNotAllowed) Error() string {
	if e.message == "" {
		return "object has been denied"
	}

	return e.message
}
