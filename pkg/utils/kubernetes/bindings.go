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
	"strings"

	rbacv1 "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CreateOrUpdateManagedClusterRoleBinding is responsible for updating a managed cluster role
func CreateOrUpdateManagedClusterRoleBinding(ctx context.Context, cc client.Client, binding *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error) {
	// @step: we first try and create the role
	if err := cc.Create(ctx, binding); err != nil {
		if !kerrors.IsAlreadyExists(err) {
			return nil, err
		}

		current := binding.DeepCopy()
		if err := cc.Get(ctx, types.NamespacedName{Name: binding.Name}, current); err != nil {
			return nil, err
		}

		if current.RoleRef.Name != binding.RoleRef.Name {
			if err := cc.Delete(ctx, current); err != nil {
				return nil, err
			}

			return CreateOrUpdateManagedClusterRoleBinding(ctx, cc, binding)
		}

		return binding, cc.Update(ctx, binding)
	}

	return binding, nil
}

// DeleteBindingsWithPrefix removes any bindings with a specific prefix
func DeleteBindingsWithPrefix(ctx context.Context, cc client.Client, prefix string) error {
	list := &rbacv1.ClusterRoleBindingList{}

	if err := cc.List(ctx, list); err != nil {
		return err
	}

	for _, binding := range list.Items {
		if strings.HasPrefix(binding.Name, prefix) {
			if err := DeleteIfExists(ctx, cc, binding.DeepCopyObject()); err != nil {
				return err
			}
		}
	}

	return nil
}
