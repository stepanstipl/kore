/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package eksnodegroup

import (
	"context"
	"time"

	core "github.com/appvia/kore/pkg/apis/core/v1"
	eksv1alpha1 "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	eksctl "github.com/appvia/kore/pkg/controllers/cloud/aws/eks"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconcile controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (t *eksNodeGroupCtrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithFields(log.Fields{
		"controller": t.Name(),
	})
	logger.Info("Reconciling EKSNodeGroup")

	// Fetch the EKSNodeGroup instance
	nodegroup := &eksv1alpha1.EKSNodeGroup{}

	if err := t.mgr.GetClient().Get(context.TODO(), request.NamespacedName, nodegroup); err != nil {
		if errors.IsNotFound(err) {

			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	logger.Info("Found AWSNodeGroup CR")

	credentials := &eksv1alpha1.EKSCredentials{}
	reference := types.NamespacedName{
		Namespace: nodegroup.Spec.Use.Namespace,
		Name:      nodegroup.Spec.Use.Name,
	}

	ctx := context.Background()

	err := t.mgr.GetClient().Get(ctx, reference, credentials)
	if err != nil {

		return reconcile.Result{}, err
	}

	logger.Info("Found AWSCredential CR")
	client, err := eksctl.NewClient(credentials, nodegroup.ClusterName, nodegroup.Spec.Region)
	nodeGroupExists, err := client.NodeGroupExists(nodegroup)
	if err != nil {

		return reconcile.Result{}, err
	}

	if nodeGroupExists {
		logger.Info("Nodegroup exists")

		return reconcile.Result{}, nil
	}

	// Set status to pending
	nodegroup.Status.Status = core.PendingStatus
	if err := t.mgr.GetClient().Status().Update(ctx, nodegroup); err != nil {
		logger.Error(err, "failed to update the resource status")
		return reconcile.Result{}, err
	}

	// Create node group
	logger.Info("Creating nodegroup")
	err = client.CreateNodeGroup(nodegroup)
	if err != nil {
		logger.Error(err, "create nodegroup error")
		return reconcile.Result{}, err
	}

	// TODO - doesn't look right
	// Wait for node group to become ACTIVE
	for {
		logger.Info("Checking the status of the node group: " + nodegroup.Spec.NodeGroupName)

		nodestatus, err := client.GetEKSNodeGroupStatus(nodegroup)
		if err != nil {
			return reconcile.Result{}, err
		}

		if nodestatus == "ACTIVE" {
			logger.Info("Nodegroup active:" + nodegroup.Spec.NodeGroupName)
			// Set status to success
			nodegroup.Status.Status = core.SuccessStatus

			if err := t.mgr.GetClient().Status().Update(ctx, nodegroup); err != nil {
				logger.Error(err, "failed to update the resource status")
				return reconcile.Result{}, err
			}
			break
		}
		if nodestatus == "ERROR" {
			logger.Info("Node group has ERROR status:" + nodegroup.Spec.NodeGroupName)
			break
		}
		time.Sleep(5000 * time.Millisecond)
	}

	return reconcile.Result{}, nil
}
