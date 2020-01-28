/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of hub-apiserver.
 *
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package controllers

import (
	"context"

	"github.com/appvia/kore/pkg/utils/kubernetes"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// ReconcileHandler is a wrapper the a controller handler
type ReconcileHandler struct {
	// HandlerFunc handles the reconcilation request
	HandlerFunc func(reconcile.Request) (reconcile.Result, error)
}

// Reconcile wraps the caller
func (r *ReconcileHandler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	if r.HandlerFunc != nil {
		return r.HandlerFunc(request)
	}

	return reconcile.Result{}, nil
}

// CreateClientFromSecret is used to retrieve the secret and create a runtime client
func CreateClientFromSecret(ctx context.Context, cc client.Client, name, namespace string) (client.Client, error) {
	// @step: retrieve the credentials for the cluster
	credentials := &v1.Secret{}
	if err := cc.Get(ctx, types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, credentials); err != nil {
		return nil, err
	}

	return kubernetes.NewRuntimeClientFromSecret(credentials)
}

// NewController is used to create and return a controller
func NewController(name string, mgr manager.Manager, src source.Source, fn reconcile.Reconciler) (controller.Controller, error) {
	ctrl, err := controller.New(name, mgr, controller.Options{
		MaxConcurrentReconciles: 10,
		Reconciler:              fn,
	})
	if err != nil {
		return nil, err
	}

	// @step: setup watches for the resources
	if err := ctrl.Watch(src,
		&handler.EnqueueRequestForObject{},
		&predicate.GenerationChangedPredicate{},
	); err != nil {
		return nil, err
	}

	return ctrl, nil
}
