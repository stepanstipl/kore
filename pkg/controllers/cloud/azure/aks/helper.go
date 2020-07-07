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

package aks

import (
	"net/http"

	aksv1alpha1 "github.com/appvia/kore/pkg/apis/aks/v1alpha1"

	"github.com/Azure/go-autorest/autorest"
)

func isNotFound(resp autorest.Response) bool {
	return responseHasStatusCode(resp, http.StatusNotFound)
}

func responseHasStatusCode(resp autorest.Response, statusCode int) bool {
	return resp.Response != nil && resp.Response.StatusCode == statusCode
}

func resourceGroupName(aks *aksv1alpha1.AKS) string {
	return "kore-" + aks.Name + "-" + aks.Spec.Location
}

func nodesResourceGroupName(aks *aksv1alpha1.AKS) string {
	return "kore-" + aks.Name + "-nodes-" + aks.Spec.Location
}
