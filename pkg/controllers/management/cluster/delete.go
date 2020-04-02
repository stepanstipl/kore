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

package cluster

import (
	"context"
	"time"

	"github.com/appvia/kore/pkg/controllers"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is responsible for deleting the cluster and all it's resources
func (a Controller) Delete(ctx context.Context, cluster *clustersv1.Cluster) (reconcile.Result, error) {
	logger := a.logger.WithFields(log.Fields{
		"name":      cluster.Name,
		"namespace": cluster.Namespace,
	})
	logger.Debug("attempting to delete the cluster from the api")

	finalizer := kubernetes.NewFinalizer(a.mgr.GetClient(), finalizerName)

	if cluster.Status.Status == corev1.DeletedStatus {
		err := finalizer.Remove(cluster)
		if err != nil {
			logger.WithError(err).Error("failed to remove the finalizer from the cluster")
		}
		return reconcile.Result{}, err
	}

	err := func() *controllers.ReconcileError {
		cluster.Status.Status = corev1.DeleteStatus

		components, err := createClusterComponents(cluster)
		if err != nil {
			return controllers.NewReconcileError(err, true)
		}

		for _, c := range components {
			componentName := a.getComponentName(c)
			var err error
			var exists bool
			if exists, err = kubernetes.GetIfExists(ctx, a.mgr.GetClient(), c); err != nil {
				return controllers.NewReconcileError(err, false).Wrapf("failed to load %s component: %w")
			}

			var status corev1.Status
			var message string
			if exists {
				switch r := c.(type) {
				case *clustersv1.Kubernetes:
					if r.GetDeletionTimestamp() == nil {
						if err := a.mgr.GetClient().Delete(ctx, c); err != nil {
							return controllers.NewReconcileError(err, false).Wrapf("failed to delete %s component: %w")
						}
					}
				}
				status, message = c.GetStatus()
			} else {
				status = corev1.DeletedStatus
			}

			if status == corev1.DeleteFailedStatus && message == "" {
				if err := c.GetComponents().Error(); err != nil {
					message = err.Error()
				}
			}

			cluster.Status.Components.SetCondition(corev1.Component{
				Name:    componentName,
				Status:  status,
				Message: message,
			})
		}

		ready := cluster.Status.Components.HasStatusForAll(corev1.DeletedStatus)
		if ready {
			cluster.Status.Status = corev1.DeletedStatus
			cluster.Status.Message = "The cluster has been deleted successfully"
			return nil
		} else if cluster.Status.Components.HasStatus(corev1.DeleteFailedStatus) {
			return controllers.NewReconcileError(cluster.Status.Components.Error(), true)
		}

		return nil
	}()

	if err != nil {
		logger.WithError(err).Error("failed to reconcile the cluster")
		if err.Critical {
			cluster.Status.Status = corev1.DeleteFailedStatus
			cluster.Status.Message = err.Error()
		}
	}

	if err := a.mgr.GetClient().Status().Update(ctx, cluster); err != nil {
		logger.WithError(err).Error("failed to update the cluster status")
		return reconcile.Result{}, err
	}

	// We haven't finished yet as we have to remove the finalizer in the last loop
	if cluster.Status.Status == corev1.DeletedStatus {
		return reconcile.Result{RequeueAfter: 1 * time.Millisecond}, nil
	}

	if cluster.Status.Status == corev1.DeleteFailedStatus {
		return reconcile.Result{}, nil
	}

	return reconcile.Result{RequeueAfter: 5 * time.Second}, err
}
