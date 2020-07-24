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

package kubernetes

import (
	"fmt"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	schema2 "github.com/appvia/kore/pkg/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Object is a Kubernetes object
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Object
type Object interface {
	runtime.Object
	metav1.Object
}

// ObjectWithStatus is a Kubernetes object where you can set/get the status and manage the status components
type ObjectWithStatus interface {
	Object
	GetStatus() (status corev1.Status, message string)
	SetStatus(status corev1.Status, message string)
}

type ObjectWithStatusComponents interface {
	Object
	StatusComponents() *corev1.Components
}

// NewObject creates a new object given the GVK definition
func NewObject(gvk schema.GroupVersionKind) (Object, error) {
	ro, err := schema2.GetScheme().New(gvk)
	if err != nil {
		return nil, err
	}

	if o, ok := ro.(Object); ok {
		return o, nil
	}

	return nil, fmt.Errorf("%T object doesn't implement kubernetes.Object", ro)
}
