/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
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
