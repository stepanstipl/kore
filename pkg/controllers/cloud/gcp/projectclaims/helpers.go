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

package projectclaims

import (
	"context"

	gcp "github.com/appvia/kore/pkg/apis/gcp/v1alpha1"
	"github.com/appvia/kore/pkg/kore"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// findProjectClaims finds any references to the project
func findProjectClaims(ctx context.Context, cc client.Client, project *gcp.Project, claim *gcp.ProjectClaim) ([]gcp.ProjectClaim, error) {
	// @step: the claim has to have a reference
	err := kore.IsOwnershipValid(claim.Status.ProjectRef)
	if err != nil {
		return nil, err
	}

	// retrieve a listing of the claims
	list := &gcp.ProjectClaimList{}
	err = cc.List(ctx, list, client.InNamespace(claim.Namespace))
	if err != nil {
		return nil, err
	}

	// filter out anything owned
	var filtered []gcp.ProjectClaim
	for i := 0; i < len(list.Items); i++ {
		// we can skip ourself
		if list.Items[i].Name == claim.Name {
			continue
		}
		matched, err := kore.IsResourceOwner(project, claim.Status.ProjectRef)
		if err != nil {
			return nil, err
		}
		if matched {
			filtered = append(filtered, list.Items[i])
		}
	}

	return filtered, nil
}
