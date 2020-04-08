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
	"fmt"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eksv1alpha1 "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/appvia/kore/pkg/utils/cloud/aws"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is responsible for deleting the aws eks cluster
func (n *eksNodeGroupCtrl) Delete(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
		"team":      request.NamespacedName.Name,
	})
	logger.Info("attempting to delete eks cluster nodegroup")

	// @step: first we need to check if we have access to the credentials

	resource := &eksv1alpha1.EKSNodeGroup{}
	if err := n.mgr.GetClient().Get(ctx, request.NamespacedName, resource); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	original := resource.DeepCopy()

	finalizer := kubernetes.NewFinalizer(n.mgr.GetClient(), finalizerName)

	requeue, err := func() (bool, error) {
		creds, err := n.GetCredentials(ctx, resource, request.NamespacedName.Name)
		if err != nil {
			return false, err
		}

		// @step: create a cloud client for us
		client, err := aws.NewBasicClient(creds, resource.Spec.Cluster.Name, resource.Spec.Region)
		if err != nil {
			return false, err
		}

		// @step: check if the cluster exists and if so we wait or the operation or the exit
		found, err := client.NodeGroupExists(resource)
		if err != nil {
			return false, fmt.Errorf("checking if cluster nodegroup exists: %s", err)
		}

		// @step: lets update the status of the resource to deleting
		if resource.Status.Status != corev1.DeletingStatus {
			resource.Status.Status = corev1.DeletingStatus

			return true, nil
		}

		if found {
			err := client.DeleteNodeGroup(resource)
			// TODO know which errors re should retry / reque
			return false, err
		}

		return false, nil
	}()
	if err != nil {
		logger.WithError(err).Error("attempting to delete the cluster nodegroup")

		return reconcile.Result{}, err
	}
	if requeue {
		if err := n.mgr.GetClient().Status().Patch(ctx, resource, client.MergeFrom(original)); err != nil {
			logger.WithError(err).Error("trying to update the resource status")

			return reconcile.Result{}, err
		}

		return reconcile.Result{Requeue: true}, nil
	}

	if err := finalizer.Remove(resource); err != nil {
		logger.WithError(err).Error("removing the finalizer")

		return reconcile.Result{}, err
	}

	logger.Info("successfully deleted the cluster nodegroup")

	return reconcile.Result{}, nil
}
