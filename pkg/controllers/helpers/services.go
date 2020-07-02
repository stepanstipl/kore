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

package helpers

import (
	"errors"
	"fmt"
	"sort"
	"time"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/schema"
	"github.com/appvia/kore/pkg/serviceproviders/application"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func ApplyServicePlanToAppService(ctx kore.Context, service *servicesv1.Service, planName string, values map[string]interface{}) error {
	servicePlan, err := ctx.Kore().ServicePlans().Get(ctx, planName)
	if err != nil {
		return fmt.Errorf("failed to get service plan %q: %w", planName, err)
	}

	service.Spec.Kind = servicePlan.Spec.Kind
	service.Spec.Plan = servicePlan.Name

	switch servicePlan.Spec.Kind {
	case application.ServiceKindApp:
		config := &application.AppConfiguration{}
		if err := servicePlan.Spec.GetConfiguration(config); err != nil {
			return err
		}
		config.Values = values
		return service.Spec.SetConfiguration(config)

	case application.ServiceKindHelmApp:
		config := &application.HelmAppConfiguration{}
		if err := servicePlan.Spec.GetConfiguration(config); err != nil {
			return err
		}
		config.Values = values
		return service.Spec.SetConfiguration(config)
	default:
		panic(fmt.Errorf("method called with invalid service kind: %s", servicePlan.Spec.Kind))
	}
}

// EnsureServices will create or update services and return reconciliation info
func EnsureServices(ctx kore.Context, services []servicesv1.Service, owner runtime.Object, components *corev1.Components) (reconcile.Result, error) {
	sortedServices := servicesv1.PriorityServiceSlice(make([]servicesv1.Service, 0, len(services)))
	for _, s := range services {
		sortedServices = append(sortedServices, s)
	}
	sort.Sort(sortedServices)

	for _, service := range sortedServices {
		gvk, found, err := schema.GetGroupKindVersion(&service)
		if err != nil || !found {
			panic(errors.New("resource GVK not found for service objects"))
		}
		service.GetObjectKind().SetGroupVersionKind(gvk)

		result, err := EnsureService(
			ctx,
			service.DeepCopy(),
			owner,
			components,
		)
		if err != nil {
			components.SetStatus("Service/"+service.Name, corev1.ErrorStatus, err.Error(), "")
			return reconcile.Result{}, err
		}
		if result.Requeue || result.RequeueAfter > 0 {
			return result, nil
		}
	}

	// Delete the removed services
	for _, comp := range *components {
		if comp.Resource == nil {
			continue
		}

		serviceExists := func() bool {
			for _, service := range sortedServices {
				if comp.Resource.Equals(corev1.MustGetOwnershipFromObject(&service)) {
					return true
				}
			}
			return false
		}()
		if serviceExists {
			continue
		}

		service := &servicesv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      comp.Resource.Name,
				Namespace: comp.Resource.Namespace,
				Annotations: map[string]string{
					kore.AnnotationOwner: kubernetes.MustGetRuntimeSelfLink(owner),
				},
			},
		}

		res, err := DeleteService(ctx, service, owner, components)
		if err != nil || res.Requeue || res.RequeueAfter > 0 {
			return res, err
		}

		components.RemoveComponent("Service/" + service.Name)
	}

	return reconcile.Result{}, nil
}

// EnsureService will create or update a service and return reconciliation info
func EnsureService(ctx kore.Context, original *servicesv1.Service, owner runtime.Object, components *corev1.Components) (reconcile.Result, error) {
	resource := corev1.MustGetOwnershipFromObject(original)
	components.SetCondition(corev1.Component{
		Name:     "Service/" + original.Name,
		Status:   corev1.PendingStatus,
		Message:  "",
		Detail:   "",
		Resource: &resource,
	})

	if original.Annotations == nil {
		original.Annotations = map[string]string{}
	}
	original.Annotations[kore.AnnotationOwner] = kubernetes.MustGetRuntimeSelfLink(owner)

	current := servicesv1.NewService(original.Name, original.Namespace)
	exists, err := kubernetes.GetIfExists(ctx, ctx.Client(), current)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get service %q: %w", current.Name, err)
	}

	if exists {
		components.SetStatus("Service/"+current.Name, current.Status.Status, current.Status.Message, "")
	}

	updated, err := kubernetes.UpdateIfChangedSinceLastUpdate(ctx, ctx.Client(), original, current)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to update %q service: %w", original.Name, err)
	}

	if updated {
		ctx.Logger().WithField("service", original.Name).Debug("service has changed")
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}

	switch current.Status.Status {
	case corev1.SuccessStatus:
		return reconcile.Result{}, nil
	case corev1.ErrorStatus, corev1.FailureStatus:
		return reconcile.Result{}, fmt.Errorf("%q admin service has an error status", current.Name)
	default:
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}
}

// DeleteServices will remove services and return reconcile status
func DeleteServices(ctx kore.Context, team string, owner runtime.Object, components *corev1.Components) (reconcile.Result, error) {
	adminServicesList, err := ctx.Kore().Teams().Team(team).Services().List(ctx, func(service servicesv1.Service) bool {
		return service.Annotations[kore.AnnotationOwner] == kubernetes.MustGetRuntimeSelfLink(owner)
	})
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to list services: %w", err)
	}

	adminServices := servicesv1.PriorityServiceSlice(adminServicesList.Items)
	sort.Sort(sort.Reverse(adminServices))

	for _, service := range adminServices {
		components.SetStatus("Service/"+service.Name, corev1.DeletingStatus, "", "")

		result, err := DeleteService(
			ctx,
			&service,
			owner,
			components,
		)
		if err != nil {
			components.SetStatus("Service/"+service.Name, corev1.ErrorStatus, err.Error(), "")
			return reconcile.Result{}, err
		}
		if result.Requeue || result.RequeueAfter > 0 {
			return result, nil
		}
	}

	return reconcile.Result{}, nil
}

// DeleteService will remove a service and return reconcile status
func DeleteService(ctx kore.Context, service *servicesv1.Service, owner runtime.Object, components *corev1.Components) (reconcile.Result, error) {
	if service.Annotations[kore.AnnotationOwner] != kubernetes.MustGetRuntimeSelfLink(owner) {
		return reconcile.Result{}, fmt.Errorf("the service can not be deleted as it doesn't belong to %s", kubernetes.MustGetRuntimeSelfLink(owner))
	}

	service = service.DeepCopy()

	exists, err := kubernetes.GetIfExists(ctx, ctx.Client(), service)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get service %q: %w", service.Name, err)
	}

	if !exists {
		components.SetStatus("Service/"+service.Name, corev1.DeletedStatus, "", "")
		return reconcile.Result{}, nil
	}

	components.SetStatus("Service/"+service.Name, service.Status.Status, service.Status.Message, "")

	if service.DeletionTimestamp != nil {
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}

	if err := kubernetes.DeleteIfExists(ctx, ctx.Client(), service); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to delete admin service %q: %w", service.Name, err)
	}

	return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
}
