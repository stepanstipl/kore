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

package podpolicy

import (
	"context"
	"fmt"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	psp "k8s.io/api/policy/v1beta1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const finalizerName = "pod-policu.clusters.kore.appvia.io"

// Reconcile is the entrypoint for the reconciliation logic
func (a pspCtrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to renconcile the managed pod seucity policy")

	// @step: retrieve the type from the api
	policy := &clustersv1.ManagedPodSecurityPolicy{}
	if err := a.mgr.GetClient().Get(context.Background(), request.NamespacedName, policy); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}
	original := policy.DeepCopy()

	// @step: create a finalizer and check if we are deleting
	finalizer := kubernetes.NewFinalizer(a.mgr.GetClient(), finalizerName)

	if finalizer.IsDeletionCandidate(policy) {
		return a.Delete(request)
	}

	logger.Debug("attempting to retrieve a list of cluster applicable")

	// @step: retrieve a list of cluster which this role applies
	list, err := a.FilterClustersBySource(ctx,
		policy.Spec.Clusters,
		policy.Spec.Teams,
		policy.Namespace)
	if err != nil {
		logger.WithError(err).Error("trying to retrieve a list of clusters")

		policy.Status.Status = corev1.FailureStatus
		policy.Status.Conditions = []corev1.Condition{{
			Detail:  err.Error(),
			Message: "Failed trying to retrieve list of clusters to apply",
		}}

		if err := a.mgr.GetClient().Status().Patch(ctx, policy, client.MergeFrom(original)); err != nil {
			logger.WithError(err).Error("trying to update the resource")

			return reconcile.Result{}, err
		}

		return reconcile.Result{}, err
	}

	logger.WithField("clusters", len(list.Items)).Debug("applying the change to x clusters")

	err = func() error {
		// @step: we iterate the clusters and apply the pod security policies
		for _, cluster := range list.Items {
			logger := logger.WithFields(log.Fields{
				"cluster": cluster.Name,
				"team":    cluster.Namespace,
			})
			logger.Debug("attempting to reconcile the managed role in cluster")

			var failed int

			_ = func() error {
				client, err := controllers.CreateClientFromSecret(ctx, a.mgr.GetClient(), cluster.Namespace, cluster.Name)
				if err != nil {
					logger.WithError(err).Error("trying to create kubernetes client")

					return err
				}
				logger.Debug("creating the managed cluster pod security policy in cluster")

				policyName := "kore.managed." + policy.Name

				if _, err := kubernetes.CreateOrUpdate(ctx, client, &psp.PodSecurityPolicy{
					ObjectMeta: metav1.ObjectMeta{
						Name: policyName,
						Labels: map[string]string{
							kore.Label("owner"): "true",
						},
					},
					Spec: policy.Spec.Policy,
				}); err != nil {
					logger.WithError(err).Error("trying to apply policy on cluster")

					return err
				}

				return nil
			}()
			policy.Status.Status = corev1.SuccessStatus
			policy.Status.Conditions = []corev1.Condition{}

			if failed > 0 {
				policy.Status.Status = corev1.FailureStatus
				policy.Status.Conditions = []corev1.Condition{{
					Message: fmt.Sprintf("Failed to apply managed pod on %d of %d clusters", failed, len(list.Items)),
				}}
			}
		}
		return nil
	}()
	if err != nil {
		logger.WithError(err).Error("trying to apply the managed security policy")

		return reconcile.Result{}, nil
	}
	if err == nil {
		if finalizer.NeedToAdd(policy) {
			if err := finalizer.Add(policy); err != nil {
				log.WithError(err).Error("trying to add the finalizer")

				return reconcile.Result{}, nil
			}
		}

		return reconcile.Result{Requeue: true}, nil
	}

	// @step: the resource has been reconcile, update the status
	if err := a.mgr.GetClient().Status().Patch(context.Background(), policy, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("trying to update the managed pod secuity policy status")

		return reconcile.Result{}, err
	}

	return reconcile.Result{RequeueAfter: 30 * time.Minute}, nil
}
