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

	appsv1 "k8s.io/api/apps/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CreateOrUpdateDeployment does what is says on the tin
func CreateOrUpdateDeployment(ctx context.Context, cc client.Client, d *appsv1.Deployment) (*appsv1.Deployment, error) {
	if err := cc.Create(ctx, d); err != nil {
		if !kerrors.IsAlreadyExists(err) {
			return nil, err
		}

		key := types.NamespacedName{
			Namespace: d.Namespace,
			Name:      d.Name,
		}
		current := d.DeepCopy()
		if err := cc.Get(ctx, key, current); err != nil {
			return nil, err
		}

		d.SetResourceVersion(current.GetResourceVersion())
		d.SetGeneration(current.GetGeneration())

		return d, cc.Update(ctx, d)
	}

	return d, nil
}
