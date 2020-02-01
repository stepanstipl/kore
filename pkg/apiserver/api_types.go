/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package apiserver

import (
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
