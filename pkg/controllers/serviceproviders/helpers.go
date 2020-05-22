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

package serviceproviders

import (
	"fmt"
	"sort"
	"time"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/appvia/kore/pkg/controllers"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	"github.com/banzaicloud/k8s-objectmatcher/patch"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func EnsureServices(ctx controllers.Context, services []servicesv1.Service, owner runtime.Object, components corev1.Components) (reconcile.Result, error) {
	sortedServices := servicesv1.PriorityServiceSlice(make([]servicesv1.Service, 0, len(services)))
	for _, s := range services {
		sortedServices = append(sortedServices, s)
	}
	sort.Sort(sortedServices)

	for _, service := range sortedServices {
		components.SetStatus("Service/"+service.Name, corev1.PendingStatus, "", "")

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

	return reconcile.Result{}, nil
}

func EnsureService(ctx controllers.Context, original *servicesv1.Service, owner runtime.Object, components corev1.Components) (reconcile.Result, error) {
	if original.Annotations == nil {
		original.Annotations = map[string]string{}
	}
	original.Annotations[kore.AnnotationOwner] = kubernetes.MustGetRuntimeSelfLink(owner)

	current := servicesv1.NewService(original.Name, original.Namespace)
	exists, err := kubernetes.GetIfExists(ctx, ctx.Client(), current)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get service %q: %w", current.Name, err)
	}

	patchAnnotator := patch.NewAnnotator(kore.Label("last-applied"))

	if !exists {
		if err := patchAnnotator.SetLastAppliedAnnotation(original); err != nil {
			return reconcile.Result{}, err
		}
		original.Status.Status = corev1.PendingStatus
		if err := ctx.Client().Create(ctx, original); err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to create service %q: %w", original.Name, err)
		}
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}

	components.SetStatus("Service/"+current.Name, current.Status.Status, current.Status.Message, "")

	patchResult, err := patch.NewPatchMaker(patchAnnotator).Calculate(
		current,
		original,
		patch.IgnoreStatusFields(),
		patch.IgnoreVolumeClaimTemplateTypeMetaAndStatus(),
	)
	if err != nil {
		return reconcile.Result{}, err
	}

	if patchResult.IsEmpty() {
		switch current.Status.Status {
		case corev1.SuccessStatus:
			return reconcile.Result{}, nil
		case corev1.ErrorStatus, corev1.FailureStatus:
			return reconcile.Result{}, fmt.Errorf("%q admin service has an error status", current.Name)
		default:
			return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
		}
	}

	ctx.Logger().WithField("diff", string(patchResult.Patch)).Debug("service has changed")

	if err := patchAnnotator.SetLastAppliedAnnotation(original); err != nil {
		return reconcile.Result{}, err
	}

	if _, err := kubernetes.CreateOrUpdate(ctx, ctx.Client(), original); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to update admin service %q: %w", original.Name, err)
	}

	return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
}

func DeleteServices(ctx controllers.Context, team string, owner runtime.Object, components corev1.Components) (reconcile.Result, error) {
	adminServicesList, err := ctx.Kore().Teams().Team(team).Services().ListFiltered(ctx, func(service servicesv1.Service) bool {
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

func DeleteService(ctx controllers.Context, service *servicesv1.Service, owner runtime.Object, components corev1.Components) (reconcile.Result, error) {
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
