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

package organization

import (
	"context"

	gcp "github.com/appvia/kore/pkg/apis/gcp/v1alpha1"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is responsible for deleting the gcp admin project
func (t *gcpCtrl) Delete(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
		"team":      request.NamespacedName.Name,
	})
	logger.Debug("attempting to delete the gcp admin project")

	// @step: first we need grab the resource from the api
	project := &gcp.Organization{}
	if err := t.mgr.GetClient().Get(ctx, request.NamespacedName, project); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
