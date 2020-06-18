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

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eksv1alpha1 "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/controllers/helpers"
	"github.com/appvia/kore/pkg/utils/cloud/aws"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "eks.compute.kore.appvia.io"
	// ComponentClusterCreator is the name of the component for the UI
	ComponentClusterCreator = "Cluster Creator"
	// ComponentClusterBootstrap is the component name for seting up cloud credentials
	ComponentClusterBootstrap = "Cluster Initialize Access"
)

func (t *eksCtrl) GetLoggerFromReq(request reconcile.Request) log.FieldLogger {
	return log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
}

func (t *eksCtrl) GetLoggerFromRes(resource *eksv1alpha1.EKS) log.FieldLogger {
	return log.WithFields(log.Fields{
		"name":      resource.Name,
		"namespace": resource.Namespace,
	})
}

// Reconcile is responsible for handling the EKS cluster
func (t *eksCtrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	logger := t.GetLoggerFromReq(request)
	logger.Debug("attempting to reconcile aws eks cluster")

	resource := &eksv1alpha1.EKS{}
	if err := t.mgr.GetClient().Get(ctx, request.NamespacedName, resource); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}
	original := resource.DeepCopyObject()
	finalizer := kubernetes.NewFinalizer(t.mgr.GetClient(), finalizerName)

	// @step: are we deleting the resource
	if finalizer.IsDeletionCandidate(resource) {
		return t.Delete(request)
	}

	if resource.Status.Conditions == nil {
		resource.Status.Conditions = corev1.Components{}
	}

	result, err := func() (reconcile.Result, error) {
		// @step: add the finalizer if require
		if finalizer.NeedToAdd(resource) {
			if err := finalizer.Add(resource); err != nil {
				logger.WithError(err).Error("trying to add finalizer from eks resource")

				return reconcile.Result{}, err
			}

			return reconcile.Result{Requeue: true}, nil
		}

		client, err := t.GetClusterClient(ctx, resource)
		if err != nil {
			logger.WithError(err).Error("trying to create eks cluster client")

			return reconcile.Result{}, err
		}

		ensure := []controllers.EnsureFunc{
			t.EnsureResourcePending(resource),
			t.EnsureClusterRoles(resource),
			t.EnsureClusterCreation(client, resource),
			t.EnsureClusterInSync(client, resource),
			t.EnsureClusterBootstrap(client, resource),
		}

		for _, handler := range ensure {
			result, err := handler(ctx)
			if err != nil {
				return reconcile.Result{}, err
			}
			if result.Requeue || result.RequeueAfter > 0 {
				return result, nil
			}
		}

		return reconcile.Result{}, nil
	}()
	if err != nil {
		logger.WithError(err).Error("trying to reconcile the eks cluster")

		resource.Status.Status = corev1.FailureStatus
	}
	// @step: we update always update the status before throwing any error
	if err := t.mgr.GetClient().Status().Patch(ctx, resource, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("updating the status of eks cluster")

		return reconcile.Result{}, err
	}

	return result, err
}

// GetClusterClient returns a EKS cluster client
func (t *eksCtrl) GetClusterClient(ctx context.Context, resource *eksv1alpha1.EKS) (*aws.Client, error) {
	e := helpers.NewKoreEKS(
		ctx,
		resource,
		t.mgr.GetClient(),
		t.Interface,
		t.GetLoggerFromRes(resource))

	c, err := e.GetClusterClient()
	if err != nil {
		resource.Status.Conditions.SetCondition(corev1.Component{
			Detail:  err.Error(),
			Name:    ComponentClusterCreator,
			Message: "Failed to create EKS client, please check credentials",
			Status:  corev1.FailureStatus,
		})
		return nil, err
	}
	return c, nil
}

// GetCredentials returns the cloud credential
func (t *eksCtrl) GetCredentials(ctx context.Context, cluster *eksv1alpha1.EKS, team string) (*aws.Credentials, error) {
	return helpers.NewKoreEKS(
		ctx,
		cluster,
		t.mgr.GetClient(),
		t.Interface,
		t.GetLoggerFromRes(cluster)).GetCredentials(cluster.Namespace)
}
