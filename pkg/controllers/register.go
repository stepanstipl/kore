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

package controllers

import (
	"log"
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

	log.Printf("adding controller %s", handler.Name())
	controllerList = append(controllerList, handler)

	return nil
}

// GetControllers returns all the registered controllers
func GetControllers() []RegisterInterface {
	controllerLock.RLock()
	defer controllerLock.RUnlock()

	return controllerList
}
