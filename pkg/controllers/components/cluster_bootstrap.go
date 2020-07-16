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

package components

import (
	"fmt"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ClusterBootstrap is a controller component for bootstrapping a cluster
type ClusterBootstrap struct {
	bootstrapCreator func(kore.Context) (controllers.Bootstrap, error)
}

// NewClusterBootstrap creates a new bootstrap component
func NewClusterBootstrap(bootstrapCreator func(kore.Context) (controllers.Bootstrap, error)) ClusterBootstrap {
	return ClusterBootstrap{bootstrapCreator: bootstrapCreator}
}

func (c ClusterBootstrap) ComponentName() string {
	return "Cluster Bootstrap"
}

func (c ClusterBootstrap) SetComponent(component *corev1.Component) {
}

func (c ClusterBootstrap) Reconcile(ctx kore.Context) (reconcile.Result, error) {
	bootstrap, err := c.bootstrapCreator(ctx)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to bootstrap the cluster: %w", err)
	}

	// TODO: we should move the business logic from Run here
	if err := controllers.NewBootstrap(bootstrap).Run(ctx, ctx.Client()); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to bootstrap the cluster: %w", err)
	}

	return reconcile.Result{}, nil
}

func (c ClusterBootstrap) Delete(ctx kore.Context) (reconcile.Result, error) {
	return reconcile.Result{}, nil
}
