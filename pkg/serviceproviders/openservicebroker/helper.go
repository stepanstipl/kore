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
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/kore"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers"

	osb "github.com/kubernetes-sigs/go-open-service-broker-client/v2"
)

func getServicePlanProviderData(ctx kore.Context, service *servicesv1.Service) (ServicePlanProviderData, error) {
	plan, err := ctx.Kore().ServicePlans().Get(ctx, service.Spec.Plan)
	if err != nil {
		return ServicePlanProviderData{}, fmt.Errorf("failed to load the service plan %q: %w", service.Spec.Plan, err)
	}
	providerData := ServicePlanProviderData{}
	if err := plan.Spec.GetProviderData(&providerData); err != nil {
		return ServicePlanProviderData{}, err
	}

	if providerData.ServiceID != "" && providerData.PlanID != "" {
		return providerData, nil
	}

	kind, err := ctx.Kore().ServiceKinds().Get(ctx, service.Spec.Kind)
	if err != nil {
		return ServicePlanProviderData{}, fmt.Errorf("failed to load the service kind %q: %w", service.Spec.Kind, err)
	}

	kindProviderData := ServiceKindProviderData{}
	if err := kind.Spec.GetProviderData(&kindProviderData); err != nil {
		return ServicePlanProviderData{}, err
	}

	return ServicePlanProviderData{
		PlanID:    kindProviderData.DefaultPlanID,
		ServiceID: kindProviderData.ServiceID,
	}, nil
}

func bindingCredentialsToStringMap(credentials map[string]interface{}) (map[string]string, error) {
	res := map[string]string{}
	for k, v := range credentials {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Array, reflect.Map, reflect.Struct:
			encoded, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			res[k] = string(encoded)
		default:
			res[k] = fmt.Sprintf("%v", v)
		}
	}
	return res, nil
}

func isHttpNotFound(err error) bool {
	if httpErr, ok := err.(osb.HTTPStatusCodeError); ok {
		return httpErr.StatusCode == http.StatusNotFound || httpErr.StatusCode == http.StatusGone
	}
	return false
}

func isHttpBadRequest(err error) bool {
	if httpErr, ok := err.(osb.HTTPStatusCodeError); ok {
		return httpErr.StatusCode == http.StatusBadRequest
	}
	return false
}

func handleError(component *corev1.Component, message string, err error) error {
	component.Update(corev1.ErrorStatus, message, err.Error())

	if isHttpBadRequest(err) {
		component.Status = corev1.FailureStatus
		return controllers.NewCriticalError(fmt.Errorf("%s: %w", message, err))
	}

	return fmt.Errorf("%s: %w", message, err)
}
