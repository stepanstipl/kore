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

package identity

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	// locker provides a spindle on the slice
	locker = &sync.RWMutex{}
	// plugins is a list of registered plugins
	plugins []Plugin
)

// Register adds an authentication plugin into the api
func Register(plugin Plugin) error {
	locker.Lock()
	defer locker.Unlock()

	log.WithFields(log.Fields{
		"plugin": plugin.Name(),
	}).Info("registering the authentication plugin")

	plugins = append(plugins, plugin)

	return nil
}

// GetPlugins returns the list of registered plugins
func GetPlugins() []Plugin {
	locker.RLock()
	defer locker.RUnlock()

	return plugins
}
