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
	"sync"
)

var (
	// handlersLock
	handlersLock = &sync.RWMutex{}
	// handlers a list of handlers to register with the api
	handlers []Handler
)

// RegisterHandler is responsible for collecting the handlers
func RegisterHandler(handler Handler) {
	handlersLock.Lock()
	defer handlersLock.Unlock()

	handlers = append(handlers, handler)
}

// GetRegisteredHandlers returns a list of handlers
func GetRegisteredHandlers() []Handler {
	handlersLock.RLock()
	defer handlersLock.RUnlock()

	return handlers
}
