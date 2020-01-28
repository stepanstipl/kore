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
