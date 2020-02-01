/*
 * Copyright (C) 2019 Appvia Ltd <info@appvia.io>
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

package controllers

import (
	"sync"
)

var (
	// controllerLock is a lock to control the controller registration
	controllerLock = &sync.RWMutex{}
	// controllerList is a list of registered controllers
	controllerList = []RegisterInterface{}
)

// Register is responsible for registering a controller
func Register(handler RegisterInterface) error {
	controllerLock.Lock()
	defer controllerLock.Unlock()

	controllerList = append(controllerList, handler)

	return nil
}

// GetControllers returns all the registered controllers
func GetControllers() []RegisterInterface {
	controllerLock.RLock()
	defer controllerLock.RUnlock()

	return controllerList
}
