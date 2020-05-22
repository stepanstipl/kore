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
	"errors"
	"fmt"
	"time"

	"github.com/appvia/kore/pkg/controllers"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applicationv1beta "sigs.k8s.io/application/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (p Provider) Reconcile(
	ctx kore.ServiceProviderContext,
	service *servicesv1.Service,
) (reconcile.Result, error) {
	config, err := getAppConfiguration(service)
	if err != nil {
		return reconcile.Result{}, err
	}

	if service.Spec.Cluster.Name == "" || service.Spec.Cluster.Namespace == "" || service.Spec.ClusterNamespace == "" {
		return reconcile.Result{}, controllers.NewCriticalError(errors.New("a cluster and namespace must be defined on the service"))
	}

	clusterClient, err := createClusterClient(ctx, service)
	if err != nil {
		return reconcile.Result{}, err
	}
	if clusterClient == nil {
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
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

	if err := kubernetes.EnsureNamespace(ctx, clusterClient, &v1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: v1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: service.Spec.ClusterNamespace,
		},
	}); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to create namespace: %q: %w", service.Spec.ClusterNamespace, err)
	}

	var app *applicationv1beta.Application
	if compiledResources.Application() != nil {
		app = compiledResources.Application().DeepCopy()
	}

	if app != nil {
		exists, err := kubernetes.GetIfExists(ctx, clusterClient, app)
		if err != nil {
			if !utils.IsMissingKind(err) {
				return reconcile.Result{}, fmt.Errorf("failed to get application %q: %w", app.Name, err)
			}
		}

		if exists {
			for _, condition := range app.Status.Conditions {
				if condition.Type == applicationv1beta.Ready {
					if condition.Status == "True" {
						service.Status.Status = corev1.SuccessStatus
						service.Status.Message = condition.Message
						// We will actively monitor the application status and update the service
						return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
					} else {
						service.Status.Status = corev1.PendingStatus
						service.Status.Message = condition.Message
					}
				}
			}
		}
	}

	providerData := &ProviderData{}
	if err := service.Status.GetProviderData(providerData); err != nil {
		return reconcile.Result{}, err
	}

	// Check if we need to delete any resources
	for _, existingResource := range providerData.Resources {
		found := false
		for _, r := range compiledResources {
			if existingResource.Equals(corev1.MustGetOwnershipFromObject(r)) {
				found = true
				break
			}
		}
		if !found {
			u := existingResource.ToUnstructured()
			if err := kubernetes.DeleteIfExists(ctx, clusterClient, &u); err != nil {
				return reconcile.Result{}, fmt.Errorf("failed to delete %s: %w", utils.GetUnstructuredSelfLink(&u), err)
			}
		}
	}

	updatedProviderData := ProviderData{}

	for _, resource := range compiledResources {
		if _, ok := resource.(*v1.Namespace); ok {
			continue
		}

		ctx.Logger.WithField("resource", kubernetes.MustGetRuntimeSelfLink(resource)).Debug("creating/updating application resource")
		if err := ensureResource(ctx, clusterClient, resource.DeepCopyObject()); err != nil {
			return reconcile.Result{}, err
		}

		updatedProviderData.Resources = append(updatedProviderData.Resources, corev1.MustGetOwnershipFromObject(resource))
	}

	if err := service.Status.SetProviderData(updatedProviderData); err != nil {
		return reconcile.Result{}, err
	}

	if app == nil {
		return reconcile.Result{}, nil
	}

	for _, condition := range app.Status.Conditions {
		if condition.Type == applicationv1beta.Error {
			if condition.Status == "True" {
				return reconcile.Result{}, errors.New(condition.Message)
			}
		}
	}

	return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
}
