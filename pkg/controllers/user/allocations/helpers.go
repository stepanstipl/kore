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

package allocations

import (
	"context"
	"errors"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

// validateAllocation is used to check the allocations is ok
func (a acCtrl) validateAllocation(ctx context.Context, allocation *configv1.Allocation) error {
	if allocation.Spec.Name == "" {
		return errors.New("no name defined")
	}
	if allocation.Spec.Summary == "" {
		return errors.New("no summary defined")
	}
	if allocation.Spec.Resource.Group == "" {
		return errors.New("no resource group")
	}
	if allocation.Spec.Resource.Version == "" {
		return errors.New("no resource version")
	}
	if allocation.Spec.Resource.Kind == "" {
		return errors.New("no resource kind")
	}
	if allocation.Spec.Resource.Namespace == "" {
		return errors.New("no resource namespace")
	}
	if allocation.Spec.Resource.Name == "" {
		return errors.New("no resource name")
	}

	// @step: check the resource exists
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   allocation.Spec.Resource.Group,
		Version: allocation.Spec.Resource.Version,
		Kind:    allocation.Spec.Resource.Kind,
	})

	if err := a.mgr.GetClient().Get(ctx, types.NamespacedName{
		Namespace: allocation.Spec.Resource.Namespace,
		Name:      allocation.Spec.Resource.Name,
	}, u); err != nil {
		if !kerrors.IsNotFound(err) {
			return err
		}

		return errors.New("resource does not exist")
	}

	return nil
}
