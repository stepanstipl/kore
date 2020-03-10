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
