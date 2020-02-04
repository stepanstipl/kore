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

	rbacv1 "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CreateOrUpdateManagedClusterRole is responsible for updating a managed cluster role
func CreateOrUpdateManagedClusterRole(ctx context.Context, cc client.Client, role *rbacv1.ClusterRole) (*rbacv1.ClusterRole, error) {
	// @step: we first try and create the role
	if err := cc.Create(ctx, role); err != nil {
		if !kerrors.IsAlreadyExists(err) {
			return nil, err
		}

		// @step: we need to retrieve the current one
		original := role.DeepCopy()
		if err := cc.Get(ctx, types.NamespacedName{Name: role.Name}, original); err != nil {
			return nil, err
		}

		return role, cc.Patch(ctx, role, client.MergeFrom(original))
	}

	return role, nil
}

// DeleteClusterRoleIfExists removes the clusterrole
func DeleteClusterRoleIfExists(ctx context.Context, cc client.Client, name string) error {
	if err := cc.Delete(ctx, &rbacv1.ClusterRole{ObjectMeta: metav1.ObjectMeta{Name: name}}); err != nil {
		if !kerrors.IsNotFound(err) {
			return err
		}

		return nil
	}

	return nil
}
