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

package openservicebroker

import (
	"fmt"
	"time"

	"github.com/appvia/kore/pkg/utils/configuration"

	"github.com/appvia/kore/pkg/kore"

	"github.com/appvia/kore/pkg/utils"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/controllers"

	"github.com/google/uuid"
	osb "github.com/kubernetes-sigs/go-open-service-broker-client/v2"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (p *Provider) Reconcile(
	ctx kore.Context,
	service *servicesv1.Service,
) (reconcile.Result, error) {
	planProviderData, err := getServicePlanProviderData(ctx, service)
	if err != nil {
		return reconcile.Result{}, err
	}

	component, _ := service.Status.Components.GetComponent(ComponentProvision)
	if component == nil {
		component = &corev1.Component{
			Name:   ComponentProvision,
			Status: corev1.Unknown,
		}
	}
	defer func() {
		service.Status.Components.SetCondition(*component)
	}()

	if component.Status == corev1.SuccessStatus {
		// Check if there was any change to the service configuration
		if service.NeedsUpdate() {
			return p.update(ctx, service)
		}

		return reconcile.Result{}, nil
	}

	if component.Status == corev1.PendingStatus {
		return p.pollLastOperation(ctx, service, component)
	}

	component.Update(corev1.PendingStatus, "", "")

	if service.Status.ProviderID == "" {
		service.Status.ProviderID = uuid.New().String()
	}

	config := map[string]interface{}{}
	if err := configuration.ParseObjectConfiguration(ctx, ctx.Client(), service, &config); err != nil {
		return reconcile.Result{}, controllers.NewCriticalError(fmt.Errorf("failed to unmarshal service configuration: %w", err))
	}

	ctx.Logger().Debug("provisioning service with service broker")

	resp, err := p.client.ProvisionInstance(&osb.ProvisionRequest{
		InstanceID:        service.Status.ProviderID,
		AcceptsIncomplete: true,
		ServiceID:         planProviderData.ServiceID,
		PlanID:            planProviderData.PlanID,
		OrganizationGUID:  "Kore",
		SpaceGUID:         service.Namespace,
		Context: map[string]interface{}{
			"team": service.Namespace,
		},
		Parameters: config,
	})
	if err != nil {
		return reconcile.Result{}, handleError(component, "failed to call provision on the service broker", err)
	}

	ctx.Logger().WithField("response", resp).Debug("provisioning response from service broker")

	if err := service.Status.SetProviderData(ProviderData{Operation: resp.OperationKey}); err != nil {
		return reconcile.Result{}, err
	}

	if !resp.Async {
		component.Update(corev1.SuccessStatus, "", "")
		return reconcile.Result{}, nil
	}

	return reconcile.Result{RequeueAfter: 5 * time.Second}, nil
}

func (p *Provider) update(
	ctx kore.Context,
	service *servicesv1.Service,
) (reconcile.Result, error) {
	planProviderData, err := getServicePlanProviderData(ctx, service)
	if err != nil {
		return reconcile.Result{}, err
	}

	component, _ := service.Status.Components.GetComponent(ComponentUpdate)
	if component == nil {
		component = &corev1.Component{
			Name:   ComponentUpdate,
			Status: corev1.Unknown,
		}
	}
	defer func() {
		service.Status.Components.SetCondition(*component)
	}()

	if component.Status == corev1.SuccessStatus {
		component.Update(corev1.Unknown, "", "")
	}

	if component.Status == corev1.PendingStatus {
		return p.pollLastOperation(ctx, service, component)
	}

	component.Update(corev1.PendingStatus, "", "")

	if service.Status.ProviderID == "" {
		service.Status.ProviderID = uuid.New().String()
	}

	config := map[string]interface{}{}
	if err := configuration.ParseObjectConfiguration(ctx, ctx.Client(), service, &config); err != nil {
		return reconcile.Result{}, controllers.NewCriticalError(fmt.Errorf("failed to unmarshal service configuration"))
	}

	ctx.Logger().Debug("updating service with service broker")

	resp, err := p.client.UpdateInstance(&osb.UpdateInstanceRequest{
		InstanceID:        service.Status.ProviderID,
		AcceptsIncomplete: true,
		ServiceID:         planProviderData.ServiceID,
		PlanID:            utils.StringPtr(planProviderData.PlanID),
		Parameters:        config,
		Context: map[string]interface{}{
			"team": service.Namespace,
		},
	})
	if err != nil {
		return reconcile.Result{}, handleError(component, "failed to call update on the service broker", err)
	}

	ctx.Logger().WithField("response", resp).Debug("update response from service broker")

	if err := service.Status.SetProviderData(ProviderData{Operation: resp.OperationKey}); err != nil {
		return reconcile.Result{}, err
	}

	if !resp.Async {
		component.Update(corev1.SuccessStatus, "", "")
		return reconcile.Result{}, nil
	}

	return reconcile.Result{RequeueAfter: 5 * time.Second}, nil

}
