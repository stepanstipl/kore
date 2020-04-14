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
	"errors"
	"fmt"

	core "github.com/appvia/kore/pkg/apis/core/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	aws "github.com/appvia/kore/pkg/utils/cloud/aws"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "eksnodegroup.compute.kore.appvia.io"
	// ComponentClusterNodegroupCreator is the name of the component for the UI
	ComponentClusterNodegroupCreator = "Cluster Nodegroup Creator"
)

// Reconcile is responsible for reconciling the eks nodegroup
func (n *ctrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to reconcile aws eks cluster nodegroup")

	resource := &eks.EKSNodeGroup{}
	if err := n.mgr.GetClient().Get(ctx, request.NamespacedName, resource); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}
	original := resource.DeepCopyObject()
	finalizer := kubernetes.NewFinalizer(n.mgr.GetClient(), finalizerName)

	// @step: are we deleting the resource
	if finalizer.IsDeletionCandidate(resource) {
		return n.Delete(request)
	}
	// @step: we need to mark the cluster as pending
	if resource.Status.Conditions == nil {
		resource.Status.Conditions = core.Components{}
	}

	result, err := func() (reconcile.Result, error) {
		// @step: add the finalizer if required
		if finalizer.NeedToAdd(resource) {
			if err := finalizer.Add(resource); err != nil {
				logger.WithError(err).Error("trying to add finalizer from eks resource")

				return reconcile.Result{}, err
			}

			return reconcile.Result{Requeue: true}, nil
		}

		// @step: retrieve the cloud credentials for the aws account
		credentials, err := n.GetCredentials(ctx, resource, resource.Namespace)
		if err != nil {
			resource.Status.Conditions.SetCondition(corev1.Component{
				Name:    ComponentClusterNodegroupCreator,
				Message: "You do not have permission to the credentials",
				Status:  corev1.FailureStatus,
			})

			return reconcile.Result{}, err
		}

		// @step: retrieve the eke client for us
		client, err := n.GetClusterClient(ctx, resource)
		if err != nil {
			logger.WithError(err).Error("trying to create eks cluster client")

			return reconcile.Result{}, err
		}

		ensure := []controllers.EnsureFunc{
			n.EnsureNodeGroupIsPending(resource),
			n.EnsureClusterReady(resource),
			n.EnsureNodeRole(resource, credentials),
			n.EnsureNodeGroup(client, resource),
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
	if err := n.mgr.GetClient().Status().Patch(ctx, resource, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("updating the status of eks nodegroup")

		return reconcile.Result{}, err
	}

	return result, err
}

// GetClusterClient returns a EKS cluster client
func (n *ctrl) GetClusterClient(ctx context.Context, resource *eks.EKSNodeGroup) (*aws.Client, error) {
	credentials, err := n.GetCredentials(ctx, resource, resource.Namespace)
	if err != nil {
		resource.Status.Conditions.SetCondition(corev1.Component{
			Name:    ComponentClusterNodegroupCreator,
			Message: "You do not have permission to the credentials",
			Status:  corev1.FailureStatus,
		})

		return nil, err
	}

	client, err := aws.NewBasicClient(credentials, resource.Spec.Cluster.Name, resource.Spec.Region)
	if err != nil {
		resource.Status.Conditions.SetCondition(corev1.Component{
			Detail:  err.Error(),
			Name:    ComponentClusterNodegroupCreator,
			Message: "Failed to create EKS client, please check credentials",
			Status:  corev1.FailureStatus,
		})

		return nil, err
	}

	return client, nil
}

// GetCredentials returns the cloud credential
func (n *ctrl) GetCredentials(ctx context.Context, ng *eks.EKSNodeGroup, team string) (*eks.EKSCredentials, error) {
	// @step: is the team permitted access to this credentials
	permitted, err := n.Teams().Team(team).Allocations().IsPermitted(ctx, ng.Spec.Credentials)
	if err != nil {
		log.WithError(err).Error("attempting to check for permission on credentials")

		return nil, fmt.Errorf("attempting to check for permission on credentials")
	}

	if !permitted {
		log.Warn("trying to build eks cluster unallocated permissions")

		return nil, errors.New("you do not have permissions to the eks credentials")
	}

	// @step: retrieve the credentials
	creds := &eks.EKSCredentials{}

	return creds, n.mgr.GetClient().Get(
		ctx,
		types.NamespacedName{
			Namespace: ng.Spec.Credentials.Namespace,
			Name:      ng.Spec.Credentials.Name,
		},
		creds,
	)
}
