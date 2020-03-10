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
