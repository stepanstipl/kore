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
	"context"
	"fmt"
	"time"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/utils"

	osb "github.com/kubernetes-sigs/go-open-service-broker-client/v2"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (s *Provider) pollLastOperation(
	ctx context.Context,
	logger logrus.FieldLogger,
	service *servicesv1.Service,
	plan *servicesv1.ServicePlan,
	component *corev1.Component,
) (reconcile.Result, error) {
	providerPlan, err := s.plan(service.Spec.Kind, plan.Name)
	if err != nil {
		return reconcile.Result{}, err
	}

	var operationKey *osb.OperationKey
	if service.Status.ProviderData != "" {
		o := osb.OperationKey(service.Status.ProviderData)
		operationKey = &o
	}

	logger.WithField("operation", operationKey).Debug("polling last operation from service broker")

	resp, err := s.client.PollLastOperation(&osb.LastOperationRequest{
		InstanceID:   service.Status.ProviderID,
		ServiceID:    utils.StringPtr(providerPlan.serviceID),
		PlanID:       utils.StringPtr(providerPlan.id),
		OperationKey: operationKey,
	})
	if err != nil {
		if component.Name == ComponentDeprovision && isHttpNotFound(err) {
			component.Update(corev1.SuccessStatus, "", "")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, handleError(component, "failed to poll last operation on the service broker", err)
	}

	logger.WithField("response", resp).Debug("last operation response from service broker")

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
