/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of hub-apiserver.
 *
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package kubernetes

import (
	"context"
	"errors"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DeleteIfExists removes the resource if it exists
func DeleteIfExists(ctx context.Context, cc client.Client, object runtime.Object) error {
	if err := cc.Delete(ctx, object); err != nil {
		if !kerrors.IsNotFound(err) {
			return err
		}

		return nil
	}

	return nil
}

// CreateOrUpdate is responsible for updating a resouce
func CreateOrUpdate(ctx context.Context, cc client.Client, object runtime.Object) (runtime.Object, error) {
	key, err := client.ObjectKeyFromObject(object)
	if err != nil {
		return nil, err
	}

	// @step: we first try and create the role
	if err := cc.Create(ctx, object); err != nil {
		if !kerrors.IsAlreadyExists(err) {
			return nil, err
		}
		// @step: we need to retrieve the current one
		original := object.DeepCopyObject()

		if err := cc.Get(ctx, key, original); err != nil {
			return nil, err
		}

		nobj, ok := object.(metav1.Object)
		if !ok {
			return nil, errors.New("object does not implement the metav1.Object")
		}
		old, ok := original.(metav1.Object)
		if !ok {
			return nil, errors.New("original object does not implement the metav1.Object")
		}
		nobj.SetResourceVersion(old.GetResourceVersion())

		return object, cc.Patch(ctx, object, client.MergeFrom(original))
	}

	return object, nil
}
