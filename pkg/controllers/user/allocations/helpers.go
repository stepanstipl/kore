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
