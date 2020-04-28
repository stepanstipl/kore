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

package servicecredentials

import (
	"context"
	"fmt"
	"time"

	"github.com/appvia/kore/pkg/kore"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	log "github.com/sirupsen/logrus"

	"github.com/appvia/kore/pkg/utils/kubernetes"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/controllers"
	k8scorev1 "k8s.io/api/core/v1"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (c *Controller) ensurePending(logger log.FieldLogger, serviceCreds *servicesv1.ServiceCredentials) controllers.EnsureFunc {
	return func(ctx context.Context) (reconcile.Result, error) {
		if serviceCreds.Status.Status == "" {
			c.resetStatus(serviceCreds)
			return reconcile.Result{Requeue: true}, nil
		}

		if serviceCreds.Status.Status != corev1.PendingStatus {
			c.resetStatus(serviceCreds)
		}

		return reconcile.Result{}, nil
	}
}

func (c *Controller) resetStatus(serviceCreds *servicesv1.ServiceCredentials) {
	serviceCreds.Status.Status = corev1.PendingStatus
	serviceCreds.Status.Components = corev1.Components{
		{
			Name:    ComponentProviderSecret,
			Status:  corev1.PendingStatus,
			Message: "",
			Detail:  "",
		},
		{
			Name:    ComponentKubernetesSecret,
			Status:  corev1.PendingStatus,
			Message: "",
			Detail:  "",
		},
	}
}

func (c *Controller) ensureFinalizer(logger log.FieldLogger, serviceCreds *servicesv1.ServiceCredentials, finalizer *kubernetes.Finalizer) controllers.EnsureFunc {
	return func(ctx context.Context) (reconcile.Result, error) {
		if finalizer.NeedToAdd(serviceCreds) {
			err := finalizer.Add(serviceCreds)
			if err != nil {
				return reconcile.Result{}, fmt.Errorf("failed to set the finalizer: %w", err)
			}
			return reconcile.Result{Requeue: true}, nil
		}
		return reconcile.Result{}, nil
	}
}

func (c *Controller) EnsureActiveCluster(logger log.FieldLogger, serviceCreds *servicesv1.ServiceCredentials) controllers.EnsureFunc {
	return func(ctx context.Context) (reconcile.Result, error) {
		// @step: check the status of the cluster
		cluster := &clustersv1.Cluster{}
		if err := c.mgr.GetClient().Get(context.Background(), types.NamespacedName{
			Name:      serviceCreds.Spec.Cluster.Name,
			Namespace: serviceCreds.Spec.Cluster.Namespace,
		}, cluster); err != nil {
			if kerrors.IsNotFound(err) {
				logger.Debug("cluster does not exist")
				return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
			}
			return reconcile.Result{}, err
		}

		// @step: check the overall status of the cluster
		if cluster.Status.Status != corev1.SuccessStatus {
			return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
		}
		return reconcile.Result{}, nil
	}
}

func (c *Controller) ensureSecret(
	logger log.FieldLogger,
	service *servicesv1.Service,
	plan *servicesv1.ServicePlan,
	serviceCreds *servicesv1.ServiceCredentials,
	provider kore.ServiceProvider) controllers.EnsureFunc {
	return func(ctx context.Context) (reconcile.Result, error) {
		// @step: create client for the cluster credentials
		client, err := controllers.CreateClientFromSecret(context.Background(), c.mgr.GetClient(),
			serviceCreds.Spec.Cluster.Namespace, serviceCreds.Spec.Cluster.Name)
		if err != nil {
			serviceCreds.Status.Components.SetCondition(corev1.Component{
				Name:    ComponentKubernetesSecret,
				Status:  corev1.FailureStatus,
				Message: "failed to create client from cluster secret",
				Detail:  err.Error(),
			})

			return reconcile.Result{}, fmt.Errorf("failed to create client from cluster secret: %w", err)
		}

		exists, err := kubernetes.HasSecret(ctx, client, serviceCreds.Spec.ClusterNamespace, serviceCreds.Name)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to get secret from cluster: %w", err)
		}

		if exists {
			return reconcile.Result{}, nil
		}

		result, credentials, err := provider.ReconcileCredentials(ctx, logger, service, plan, serviceCreds)
		if err != nil {
			serviceCreds.Status.Components.SetCondition(corev1.Component{
				Name:    ComponentProviderSecret,
				Status:  corev1.FailureStatus,
				Message: "failed to request secret from service provider",
				Detail:  err.Error(),
			})

			return reconcile.Result{}, fmt.Errorf("failed to request secret from service provider: %w", err)
		}
		if result.Requeue || result.RequeueAfter > 0 {
			return result, nil
		}

		serviceCreds.Status.Components.SetStatus(ComponentProviderSecret, corev1.SuccessStatus, "", "")

		secret := &k8scorev1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: k8scorev1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceCreds.Name,
				Namespace: serviceCreds.Spec.ClusterNamespace,
			},
			Type:       "generic",
			StringData: credentials,
		}

		if _, err := kubernetes.CreateOrUpdateSecret(ctx, client, secret); err != nil {
			serviceCreds.Status.Components.SetCondition(corev1.Component{
				Name:    ComponentKubernetesSecret,
				Status:  corev1.FailureStatus,
				Message: "failed to create Secret object in cluster",
				Detail:  err.Error(),
			})

			return reconcile.Result{}, fmt.Errorf("failed to create Secret object in cluster: %w", err)
		}

		serviceCreds.Status.Components.SetStatus(ComponentKubernetesSecret, corev1.SuccessStatus, "", "")

		return reconcile.Result{}, nil
	}
}
