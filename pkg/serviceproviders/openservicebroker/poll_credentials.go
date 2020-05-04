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

func (p *Provider) pollLastBindingOperation(
	ctx kore.ServiceProviderContext,
	service *servicesv1.Service,
	creds *servicesv1.ServiceCredentials,
	component *corev1.Component,
) (reconcile.Result, map[string]string, error) {
	providerPlan, err := p.plan(service.Spec.Kind, service.Spec.Plan)
	if err != nil {
		return reconcile.Result{}, nil, err
	}

	providerData := ProviderData{}
	if err := creds.Status.GetProviderData(&providerData); err != nil {
		return reconcile.Result{}, nil, err
	}

	ctx.Logger.WithField("operation", providerData.Operation).Debug("polling last bind operation from service broker")

	resp, err := p.client.PollBindingLastOperation(&osb.BindingLastOperationRequest{
		InstanceID:   service.Status.ProviderID,
		BindingID:    creds.Status.ProviderID,
		ServiceID:    utils.StringPtr(providerPlan.serviceID),
		PlanID:       utils.StringPtr(providerPlan.id),
		OperationKey: providerData.Operation,
	})
	if err != nil {
		if component.Name == ComponentUnbind && isHttpNotFound(err) {
			component.Status = corev1.SuccessStatus
			return reconcile.Result{}, nil, nil
		}
		return reconcile.Result{}, nil, handleError(component, "failed to poll last bind operation on the service broker", err)
	}

	ctx.Logger.WithField("response", resp).Debug("last bind operation response from service broker")

	component.Message = utils.StringValue(resp.Description)

	switch resp.State {
	case osb.StateInProgress:
		return reconcile.Result{RequeueAfter: 5 * time.Second}, nil, nil
	case osb.StateSucceeded:
		if component.Name == ComponentUnbind {
			component.Status = corev1.SuccessStatus
			return reconcile.Result{}, nil, nil
		}

		ctx.Logger.Debug("requesting binding details")
		resp, err := p.client.GetBinding(&osb.GetBindingRequest{
			InstanceID: service.Status.ProviderID,
			BindingID:  creds.Status.ProviderID,
		})
		if err != nil {
			return reconcile.Result{}, nil, handleError(component, "failed to get binding details from the service broker", err)
		}

		bindingCredentials, err := bindingCredentialsToStringMap(resp.Credentials)
		if err != nil {
			return reconcile.Result{}, nil, controllers.NewCriticalError(fmt.Errorf("failed to encode binding credentials from the service broker: %w", err))
		}

		component.Status = corev1.SuccessStatus
		return reconcile.Result{}, bindingCredentials, nil
	case osb.StateFailed:
		component.Status = corev1.FailureStatus
		return reconcile.Result{}, nil, controllers.NewCriticalError(fmt.Errorf("last bind operation failed on the service broker: %s", component.Message))
	default:
		return reconcile.Result{}, nil, fmt.Errorf("invalid last bind operation state from service broker: %s", resp.State)
	}
}
