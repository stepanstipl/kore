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
	"time"

	"github.com/appvia/kore/pkg/kore"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"

	osb "github.com/kubernetes-sigs/go-open-service-broker-client/v2"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (p *Provider) DeleteCredentials(
	ctx kore.Context,
	service *servicesv1.Service,
	creds *servicesv1.ServiceCredentials,
) (reconcile.Result, error) {
	planProviderData, err := getServicePlanProviderData(ctx, service)
	if err != nil {
		return reconcile.Result{}, err
	}

	component, _ := creds.Status.Components.GetComponent(ComponentUnbind)
	if component == nil {
		component = &corev1.Component{
			Name:   ComponentUnbind,
			Status: corev1.Unknown,
		}
	}
	defer func() {
		creds.Status.Components.SetCondition(*component)
	}()

	if component.Status == corev1.SuccessStatus {
		return reconcile.Result{}, nil
	}

	if component.Status == corev1.PendingStatus {
		res, _, err := p.pollLastBindingOperation(ctx, service, creds, component)
		return res, err
	}

	component.Update(corev1.PendingStatus, "", "")

	ctx.Logger().Debug("calling unbind on the service broker")

	unbindRequest := &osb.UnbindRequest{
		AcceptsIncomplete: true,
		InstanceID:        service.Status.ProviderID,
		BindingID:         creds.Status.ProviderID,
		ServiceID:         planProviderData.ServiceID,
		PlanID:            planProviderData.PlanID,
	}

	resp, err := p.client.Unbind(unbindRequest)
	if err != nil && osb.IsAsyncBindingOperationsNotAllowedError(err) {
		unbindRequest.AcceptsIncomplete = false
		resp, err = p.client.Unbind(unbindRequest)
	}
	if err != nil {
		if isHttpNotFound(err) {
			component.Status = corev1.SuccessStatus
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, handleError(component, "failed to call unbind on the service broker", err)
	}

	ctx.Logger().WithField("response", resp).Debug("unbind response from service broker")

	if err := creds.Status.SetProviderData(ProviderData{Operation: resp.OperationKey}); err != nil {
		return reconcile.Result{}, err
	}

	if !resp.Async {
		component.Update(corev1.SuccessStatus, "", "")
		return reconcile.Result{}, nil
	}

	return reconcile.Result{RequeueAfter: 5 * time.Second}, nil
}
