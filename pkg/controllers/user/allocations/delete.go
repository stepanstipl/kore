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

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is responsible for deleting any allocations
func (a acCtrl) Delete(ctx context.Context, object *configv1.Allocation) (reconcile.Result, error) {
	logger := log.WithFields(log.Fields{
		"resource.name":      object.Name,
		"resource.namespace": object.Namespace,
		"team":               object.Namespace,
	})
	logger.Info("attempting to remove the allocation")

	return reconcile.Result{}, nil
}
