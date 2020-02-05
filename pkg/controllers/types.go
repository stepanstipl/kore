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
	"context"

	"github.com/appvia/kore/pkg/kore"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s.io/client-go/rest"
)

// NameInterface defines a requirement to implement it's name
type NameInterface interface {
	// Name is the name of the controller
	Name() string
}

type RegisterInterface interface {
	NameInterface
	// Run starts the controller
	Run(context.Context, *rest.Config, kore.Interface) error
	// Stop instructs the controller to stop
	Stop(context.Context) error
}

// Interface is the contract for a controller
type Interface interface {
	RegisterInterface
	// Reconcile is the reconcile method
	Reconcile(reconcile.Request) (reconcile.Result, error)
}
