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

package aks

import (
	"context"
	"fmt"

	aksv1alpha1 "github.com/appvia/kore/pkg/apis/aks/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	kube "github.com/appvia/kore/pkg/utils/kubernetes"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type bootImpl struct {
	cluster          *aksv1alpha1.AKS
	kubernetesClient kubernetes.Interface
}

var _ controllers.Bootstrap = &bootImpl{}

// NewBootstrapClient returns a bootstrap client for EKS
func NewBootstrapClient(cluster *aksv1alpha1.AKS, clientToken, caCertificate string) (controllers.Bootstrap, error) {
	kubernetesClient, err := kube.NewFromToken(cluster.Status.Endpoint, clientToken, caCertificate)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}
	return &bootImpl{
		cluster:          cluster,
		kubernetesClient: kubernetesClient,
	}, nil
}

func (b bootImpl) Run(ctx context.Context, client client.Client) error {
	return nil
}

func (b bootImpl) GetCluster() *controllers.BootCluster {
	return &controllers.BootCluster{
		Endpoint:  b.cluster.Status.Endpoint,
		Name:      b.cluster.Name,
		Namespace: b.cluster.Namespace,
		Client:    b.kubernetesClient,
	}
}

func (b bootImpl) GetClusterObj() runtime.Object {
	return b.cluster
}

func (b bootImpl) GetLogger() *log.Entry {
	logger := log.WithFields(log.Fields{
		"endpoint":  b.cluster.Status.Endpoint,
		"name":      b.cluster.Name,
		"namespace": b.cluster.Namespace,
		"bootstrap": "aks",
	})
	return logger
}
