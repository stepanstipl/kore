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

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	k8scorev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (c *Controller) delete(
	ctx context.Context,
	logger log.FieldLogger,
	service *servicesv1.Service,
	serviceCreds *servicesv1.ServiceCredentials,
	finalizer *kubernetes.Finalizer,
	provider kore.ServiceProvider,
) (reconcile.Result, error) {
	logger.Debug("attempting to delete service credentials from the api")

	if serviceCreds.Status.Status == corev1.DeletedStatus {
		err := finalizer.Remove(serviceCreds)
		if err != nil {
			logger.WithError(err).Error("failed to remove the finalizer from the service credentials")
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	original := serviceCreds.DeepCopyObject()

	result, err := func() (reconcile.Result, error) {
		if !serviceCreds.Status.Status.OneOf(corev1.DeletingStatus, corev1.DeleteFailedStatus, corev1.ErrorStatus) {
			serviceCreds.Status.Status = corev1.DeletingStatus
			return reconcile.Result{Requeue: true}, nil
		}

		// @step: create client for the cluster credentials
		client, err := controllers.CreateClient(context.Background(), c.mgr.GetClient(), serviceCreds.Spec.Cluster)
		if err != nil {
			serviceCreds.Status.Components.SetCondition(corev1.Component{
				Name:    ComponentKubernetesSecret,
				Status:  corev1.DeleteFailedStatus,
				Message: "Failed to create client from cluster secret",
				Detail:  err.Error(),
			})

			return reconcile.Result{}, fmt.Errorf("failed to create client from cluster secret: %w", err)
		}

		secret := &k8scorev1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: k8scorev1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceCreds.Name,
				Namespace: serviceCreds.Spec.ClusterNamespace,
			},
		}

		if err := client.Delete(ctx, secret); err != nil {
			if !kerrors.IsNotFound(err) {
				serviceCreds.Status.Components.SetCondition(corev1.Component{
					Name:    ComponentKubernetesSecret,
					Status:  corev1.DeleteFailedStatus,
					Message: "Failed to delete Secret object from the cluster",
					Detail:  err.Error(),
				})

				return reconcile.Result{}, fmt.Errorf("failed to delete Secret object from the cluster: %w", err)
			}
		}

		serviceCreds.Status.Components.SetStatus(ComponentKubernetesSecret, corev1.DeletedStatus, "", "")

		result, err := provider.DeleteCredentials(
			kore.NewContext(ctx, logger, c.mgr.GetClient(), c),
			service, serviceCreds,
		)
		if err != nil {
			serviceCreds.Status.Components.SetCondition(corev1.Component{
				Name:    ComponentProviderSecret,
				Status:  corev1.DeleteFailedStatus,
				Message: "Failed to request secret removal from service provider",
				Detail:  err.Error(),
			})
			return reconcile.Result{}, fmt.Errorf("failed to request secret removal from service provider: %w", err)
		}
		if result.Requeue || result.RequeueAfter > 0 {
			return result, nil
		}

		serviceCreds.Status.Components.SetStatus(ComponentProviderSecret, corev1.DeletedStatus, "", "")

		return reconcile.Result{}, nil
	}()

	if err != nil {
		logger.WithError(err).Error("failed to delete the service credentials")

		serviceCreds.Status.Status = corev1.ErrorStatus
		serviceCreds.Status.Message = err.Error()

		if controllers.IsCriticalError(err) {
			serviceCreds.Status.Status = corev1.DeleteFailedStatus
		}
	}

	if err == nil && !result.Requeue && result.RequeueAfter == 0 {
		serviceCreds.Status.Status = corev1.DeletedStatus
	}

	if err := c.mgr.GetClient().Status().Patch(ctx, serviceCreds, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("failed to update the service credentials status")
		return reconcile.Result{}, err
	}

	if err != nil {
		if controllers.IsCriticalError(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// We haven't finished yet as we have to remove the finalizer in the last loop
	if serviceCreds.Status.Status == corev1.DeletedStatus {
		return reconcile.Result{Requeue: true}, nil
	}

	if result.Requeue || result.RequeueAfter > 0 {
		return result, nil
	}

	return reconcile.Result{RequeueAfter: 30 * time.Second}, err
}
