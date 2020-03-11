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

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CreateOrUpdateConfigMap does what is says on the tin
func CreateOrUpdateConfigMap(ctx context.Context, cc client.Client, cm *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	if err := cc.Create(ctx, cm); err != nil {
		if !kerrors.IsAlreadyExists(err) {
			return nil, err
		}

		key := types.NamespacedName{
			Namespace: cm.Namespace,
			Name:      cm.Name,
		}
		current := cm.DeepCopy()
		if err := cc.Get(ctx, key, current); err != nil {
			return nil, err
		}

		cm.SetResourceVersion(current.GetResourceVersion())
		cm.SetGeneration(current.GetGeneration())

		return cm, cc.Update(ctx, cm)
	}

	return cm, nil
}
