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

package application

import (
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (p Provider) Delete(
	ctx kore.Context,
	service *servicesv1.Service,
) (reconcile.Result, error) {
	clusterClient, err := createClusterClient(ctx, service)
	if err != nil {
		return reconcile.Result{}, err
	}
	if clusterClient == nil {
		// If we can't create a client, we can bail out
		return reconcile.Result{}, nil
	}

	config, err := getAppConfiguration(service)
	if err != nil {
		return reconcile.Result{}, err
	}

	compiledResources, err := config.CompileResources(ResourceParams{
		Release: Release{
			Name:      service.Name,
			Namespace: service.Spec.ClusterNamespace,
		},
		Values: config.Values,
	})
	if err != nil {
		return reconcile.Result{}, err
	}

	existingResources := len(compiledResources)
	for _, resource := range compiledResources {
		if _, ok := resource.(*v1.Namespace); ok {
			existingResources--
			continue
		}

		exists, err := kubernetes.CheckIfExists(ctx, clusterClient, resource)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to check %s: %w", kubernetes.MustGetRuntimeSelfLink(resource), err)
		}

		if !exists {
			existingResources--
			continue
		}

		resourceMeta, err := meta.Accessor(resource)
		if err != nil {
			return reconcile.Result{}, err
		}

		if resourceMeta.GetDeletionTimestamp() == nil {
			ctx.Logger().WithField("resource", kubernetes.MustGetRuntimeSelfLink(resource)).Debug("deleting application resource")
			if err := kubernetes.DeleteIfExists(ctx, clusterClient, resource); err != nil {
				return reconcile.Result{}, fmt.Errorf("failed to delete %s: %w", kubernetes.MustGetRuntimeSelfLink(resource), err)
			}
		}
	}

	if existingResources > 0 {
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}

	return reconcile.Result{}, nil
}
