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

package eksnodegroup

import (
	"context"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/appvia/kore/pkg/utils/cloud/aws"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// EnsureDeletionStatus makes sure the resource is set to deleting
func (n *eksNodeGroupCtrl) EnsureDeletionStatus(ctx context.Context, resource runtime.Object) (reconcile.Result, error) {
	group := resource.(*eks.EKSNodeGroup)

	if group.Status.Status != corev1.DeletingStatus {
		group.Status.Status = corev1.DeletingStatus

		return reconcile.Result{Requeue: true}, nil
	}

	return reconcile.Result{}, nil
}

// EnsureDeletion ensures the nodegroup is deleting
func (n *eksNodeGroupCtrl) EnsureDeletion(ctx context.Context, resource runtime.Object) (reconcile.Result, error) {
	group := resource.(*eks.EKSNodeGroup)

	logger := log.WithFields(log.Fields{
		"name":      group.Name,
		"namespace": group.Namespace,
	})
	logger.Debug("attempting to delete eks nodegroup")

	creds, err := n.GetCredentials(ctx, group, group.Namespace)
	if err != nil {
		return reconcile.Result{}, err
	}

	// @step: create a cloud client for us
	client, err := aws.NewBasicClient(creds, group.Spec.Cluster.Name, group.Spec.Region)
	if err != nil {
		log.WithError(err).Error("trying to create a aws client for the nodegroup")

		return reconcile.Result{}, err
	}

	// @step: check if the nodegroup exists and if so we wait or the operation or the exit
	found, err := client.NodeGroupExists(ctx, group)
	if err != nil {
		log.WithError(err).Error("trying to check if nodegroup exists")

		return reconcile.Result{}, err
	}
	if found {
		return reconcile.Result{}, client.DeleteNodeGroup(ctx, group)
	}

	return reconcile.Result{}, nil
}
