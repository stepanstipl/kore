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
	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	kube "github.com/appvia/kore/pkg/utils/kubernetes"

	awssess "github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
	k8s "k8s.io/client-go/kubernetes"
)

type bootImpl struct {
	// credentials are the gke credentials
	sess *awssess.Session
	// cluster is the gke cluster
	cluster *eks.EKS
	// client is the k8s client
	client k8s.Interface
}

// NewBootstrapClient returns a bootstrap client
func NewBootstrapClient(cluster *eks.EKS, sess *awssess.Session) (*bootImpl, error) {
	// @step: retrieve the credentials for the cluster

	client, err := kube.NewEKSClient(
		cluster.Spec.Name,
		cluster.Status.Endpoint,
		cluster.Spec.RoleARN,
		cluster.Spec.Region,
		cluster.Status.CACertificate,
		sess,
	)
	if err != nil {
		log.WithError(err).Error("trying to create eks kubernetes client from credentials")

		return nil, err
	}
	return &bootImpl{
		cluster: cluster,
		client:  client,
		sess:    sess,
	}, nil
}
