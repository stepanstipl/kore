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

package controllers

import (
	"time"

	"github.com/appvia/kore/pkg/hub"
	"github.com/appvia/kore/pkg/schema"

	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// SingletontControllerOptions provides a single runner
func SingletontControllerOptions(handler Interface) controller.Options {
	return controller.Options{
		MaxConcurrentReconciles: 1,
		Reconciler:              handler,
	}
}

// DefaultControllerOptions returns default options
func DefaultControllerOptions(handler Interface) controller.Options {
	return controller.Options{
		MaxConcurrentReconciles: 10,
		Reconciler:              handler,
	}
}

// DefaultManagerOptions present default options for the managers
func DefaultManagerOptions(handler NameInterface) manager.Options {
	resync := time.Minute * 10

	return manager.Options{
		LeaderElection:          true,
		LeaderElectionID:        handler.Name() + "-lck",
		LeaderElectionNamespace: hub.HubNamespace,
		MetricsBindAddress:      "0",
		Scheme:                  schema.GetScheme(),
		SyncPeriod:              &resync,
	}
}
