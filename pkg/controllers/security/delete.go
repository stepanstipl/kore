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

package security

import (
	"github.com/appvia/kore/pkg/kore"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is responsible handling any cleaup
func (c *Controller) Delete(ctx kore.Context, resource runtime.Object) (reconcile.Result, error) {
	ctx.Logger().Debug("attempting to reconcile the security scans deletion")

	var err error
	switch c.kind {
	case "Cluster":
		o := resource.(*clustersv1.Cluster)
		err = ctx.Kore().Security().ArchiveResourceScans(ctx, o.TypeMeta, o.ObjectMeta)
	case "Plan":
		o := resource.(*configv1.Plan)
		err = ctx.Kore().Security().ArchiveResourceScans(ctx, o.TypeMeta, o.ObjectMeta)
	}
	if err != nil {
		ctx.Logger().WithError(err).Warning("error while archiving security scans")

		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
