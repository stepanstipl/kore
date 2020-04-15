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

package apiserver

import (
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Error is a generic error returns by the api
type Error struct {
	// Code is an optional machine readable code used to describe the code
	Code int `json:"code"`
	// Detail is the actual error thrown by the upstream
	Detail string `json:"detail"`
	// Message is a human readable message related to the error
	Message string `json:"message"`
	// URI is the uri of the request
	URI string `json:"uri"`
	// Verb was the http request verb used
	Verb string `json:"verb"`
}

// Error returns the error message
func (e Error) Error() string {
	return e.Message
}

// OwnershipList is an ownership list
type OwnershipList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []corev1.Ownership `json:"items"`
}

// UnstructuredList is a unstructured list
type UnstructuredList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []interface{} `json:"items,omitempty"`
}

// List is a list of strings
type List struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []string `json:"items"`
}

// TeamPlan is an API-only struct used to represent a bunch of info about a plan
// in the context of a team so the UI doesn't have to make about ten API calls and
// work it all out itself.
type TeamPlan struct {
	Schema            string            `json:"schema"`
	ParameterEditable map[string]bool   `json:"parameterEditable"`
	Plan              configv1.PlanSpec `json:"plan,omitempty"`
}
