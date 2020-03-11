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

package eks

import (
	"context"

	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	kube "github.com/appvia/kore/pkg/utils/kubernetes"

	awssess "github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	awsauth "sigs.k8s.io/aws-iam-authenticator/pkg/token"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type bootImpl struct {
	// sess is the aws session
	sess *awssess.Session
	// cluster is the gke cluster
	cluster *eks.EKS
	// client for kubernetes to access booting cluster
	client kubernetes.Interface
}

var _ controllers.Bootstrap = &bootImpl{}

// NewBootstrapClient returns a bootstrap client for EKS
func NewBootstrapClient(cluster *eks.EKS, sess *awssess.Session) (controllers.Bootstrap, error) {
	// Get AWS cloud specific credentials
	g, err := awsauth.NewGenerator(false, false)
	if err != nil {
		return nil, err
	}

	// Get EKS AWS token
	t, err := g.GetWithOptions(&awsauth.GetTokenOptions{
		ClusterID: cluster.Name,
		Session:   sess,
	})
	if err != nil {
		return nil, err
	}
	// Now a client from the token
	client, err := kube.NewFromToken(cluster.Status.Endpoint, t.Token, cluster.Status.CACertificate)
	if err != nil {
		log.WithError(err).Error("trying to create eks kubernetes client from credentials")

		return nil, err
	}

	// Now return a generic implementation from our cloud specific version
	boot := &bootImpl{
		cluster: cluster,
		sess:    sess,
		client:  client,
	}
	return controllers.NewBootstrap(boot), nil
}

func (b *bootImpl) GetLogger() *log.Entry {
	cluster := b.GetCluster()
	logger := log.WithFields(log.Fields{
		"endpoint":  cluster.Endpoint,
		"name":      cluster.Name,
		"namespace": cluster.Namespace,
		"bootstrap": "eks",
	})
	return logger
}

func (b *bootImpl) GetCluster() *controllers.BootCluster {
	return &controllers.BootCluster{
		Endpoint:  b.cluster.Status.Endpoint,
		Name:      b.cluster.Name,
		Namespace: b.cluster.Namespace,
		Client:    b.client,
	}
}

func (b *bootImpl) GetClusterObj() runtime.Object {
	return b.cluster
}

// Cloud specifics - at the moment none for AWS would be PSP's etc
func (b *bootImpl) Run(ctx context.Context, cc client.Client) error {
	logger := b.GetLogger()
	logger.Warn("AWS specifics not implimented")
	return nil
}
