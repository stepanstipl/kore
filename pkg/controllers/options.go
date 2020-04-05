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
	"time"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/schema"

	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// DefaultControllerOptions returns default options
func DefaultControllerOptions(reconciler reconcile.Reconciler) controller.Options {
	return controller.Options{
		MaxConcurrentReconciles: 10,
		Reconciler:              reconciler,
	}
}

// DefaultManagerOptions present default options for the managers
func DefaultManagerOptions(handler NameInterface) manager.Options {
	resync := time.Minute * 10

	return manager.Options{
		LeaderElection:          true,
		LeaderElectionID:        handler.Name() + "-lck",
		LeaderElectionNamespace: kore.HubNamespace,
		MetricsBindAddress:      "0",
		Scheme:                  schema.GetScheme(),
		SyncPeriod:              &resync,
	}
}
