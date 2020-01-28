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

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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
