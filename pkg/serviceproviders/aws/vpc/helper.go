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

	"github.com/appvia/kore/pkg/kore"

	"github.com/appvia/kore/pkg/controllers"

	eksv1alpha1 "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/utils/cloud/aws"
)

func (p Provider) createVPCClient(ctx kore.ServiceProviderContext, service *servicesv1.Service) (*aws.VPCClient, error) {
	ctx.Logger.Debug("retrieving the VPC credentials")

	credentials := &eksv1alpha1.EKSCredentials{}
	if err := ctx.Client.Get(ctx, service.Spec.Credentials.NamespacedName(), credentials); err != nil {
		return nil, fmt.Errorf("failed to retrieve EKS credentials: %w", err)
	}

	var config Configuration
	if err := service.Spec.GetConfiguration(&config); err != nil {
		return nil, controllers.NewCriticalError(fmt.Errorf("failed to unmarshal service configuration: %w", err))
	}
	if config.Name == "" {
		config.Name = service.Name
	}

	return aws.NewVPCClient(aws.Credentials{
		AccessKeyID:     credentials.Spec.AccessKeyID,
		SecretAccessKey: credentials.Spec.SecretAccessKey,
	}, aws.VPC{
		Name:        config.Name,
		Region:      config.Region,
		CidrBlock:   config.PrivateIPV4Cidr,
		SubnetCount: config.SubnetCount,
		Tags: map[string]string{
			aws.TagKoreManaged: "true",
		},
	})
}
