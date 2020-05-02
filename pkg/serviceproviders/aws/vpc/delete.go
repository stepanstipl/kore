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

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (p Provider) Delete(
	ctx kore.ServiceProviderContext,
	service *servicesv1.Service,
) (reconcile.Result, error) {
	client, err := p.createVPCClient(ctx, service)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("configuration is invalid: %w", err)
	}

	exists, err := client.Exists()
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("checking if vpc exists - %s", err)
	}

	if exists {
		ready, err := client.Delete(ctx)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to delete the EKS VPC: %w", err)
		}
		if !ready {
			return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
		}
	}

	return reconcile.Result{}, nil
}
