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

	"github.com/appvia/kore/pkg/kore"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/utils"

	osb "github.com/kubernetes-sigs/go-open-service-broker-client/v2"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (p *Provider) pollLastOperation(
	ctx kore.ServiceProviderContext,
	service *servicesv1.Service,
	component *corev1.Component,
) (reconcile.Result, error) {
	providerPlan, err := p.plan(service)
	if err != nil {
		return reconcile.Result{}, err
	}

	providerData := ProviderData{}
	if err := service.Status.GetProviderData(&providerData); err != nil {
		return reconcile.Result{}, err
	}

	ctx.Logger.WithField("operation", providerData.Operation).Debug("polling last operation from service broker")

	resp, err := p.client.PollLastOperation(&osb.LastOperationRequest{
		InstanceID:   service.Status.ProviderID,
		ServiceID:    utils.StringPtr(providerPlan.serviceID),
		PlanID:       utils.StringPtr(providerPlan.osbPlan.ID),
		OperationKey: providerData.Operation,
	})
	if err != nil {
		if component.Name == ComponentDeprovision && isHttpNotFound(err) {
			component.Update(corev1.SuccessStatus, "", "")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, handleError(component, "failed to poll last operation on the service broker", err)
	}

	ctx.Logger.WithField("response", resp).Debug("last operation response from service broker")

	component.Message = utils.StringValue(resp.Description)

	switch resp.State {
	case osb.StateInProgress:
		return reconcile.Result{RequeueAfter: 5 * time.Second}, nil
	case osb.StateSucceeded:
		component.Status = corev1.SuccessStatus
		return reconcile.Result{}, nil
	case osb.StateFailed:
		component.Status = corev1.FailureStatus
		return reconcile.Result{}, controllers.NewCriticalError(fmt.Errorf("last operation failed on the service broker: %s", component.Message))
	default:
		return reconcile.Result{}, fmt.Errorf("invalid last operation state from service broker: %s", resp.State)
	}
}
