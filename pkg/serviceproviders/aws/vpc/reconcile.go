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

package vpc

import (
	"fmt"
	"time"

	"github.com/appvia/kore/pkg/kore"

	"github.com/appvia/kore/pkg/controllers"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (p Provider) Reconcile(
	ctx kore.ServiceProviderContext,
	service *servicesv1.Service,
) (reconcile.Result, error) {
	client, err := p.createVPCClient(ctx, service)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("configuration is invalid: %w", err)
	}

	// Ensure this only reports if it exists when all resources exist or ensure update works
	ready, err := client.Ensure()
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to create or update the EKS VPC: %w", err)
	}

	if !ready {
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}

	providerData := ProviderData{
		PrivateSubnetIDs:          client.VPC.PrivateSubnetIDs,
		PublicSubnetIDs:           client.VPC.PublicSubnetIDs,
		SecurityGroupIDs:          []string{client.VPC.ControlPlaneSecurityGroupID},
		PublicIPV4EgressAddresses: client.VPC.PublicIPV4EgressAddresses,
	}

	if err := service.Status.SetProviderData(providerData); err != nil {
		return reconcile.Result{}, controllers.NewCriticalError(err)
	}

	return reconcile.Result{}, nil
}
