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
	"context"
	"errors"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CheckIfNamespaceExists checks if the namespace exists
func CheckIfNamespaceExists(ctx context.Context, client client.Client, name string) (bool, error) {
	ns := &v1.Namespace{}
	ns.Name = name

	return CheckIfExists(ctx, client, ns)
}

// EnsureNamespace makes sure the namespace exists
func EnsureNamespace(ctx context.Context, cc client.Client, namespace *corev1.Namespace) error {
	if namespace == nil {
		return errors.New("no namespace defined")
	}
	if err := cc.Create(ctx, namespace); err != nil {
		if !kerrors.IsAlreadyExists(err) {
			return err
		}
		original := &corev1.Namespace{}
		if err := cc.Get(ctx, types.NamespacedName{Name: namespace.Name}, original); err != nil {
			return err
		}

		return cc.Patch(ctx, namespace, client.MergeFrom(original))
	}

	return nil
}
