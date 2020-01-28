/*
 * Copyright (C) 2019  Rohith Jayawardene <gambol99@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package gke

import (
	"context"
	"errors"
	"fmt"

	core "github.com/appvia/kore/pkg/apis/core/v1"
	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "gke.compute.hub.appvia.io"
)

// Reconcile is the entrypoint for the reconcilation logic
func (t *gkeCtrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})

	logger.Debug("attempting to reconcile gke cluster")

	resource := &gke.GKE{}
	if err := t.mgr.GetClient().Get(ctx, request.NamespacedName, resource); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	finalizer := kubernetes.NewFinalizer(t.mgr.GetClient(), finalizerName)

	// @step: are we deleting the resource
	if finalizer.IsDeletionCandidate(resource) {
		return t.Delete(request)
	}

	requeue, err := func() (bool, error) {
		logger.Debug("retrieving the gke cluster credential")
		// @step: first we need to check if we have access to the credentials
		creds, err := t.GetCredentials(ctx, resource, resource.Namespace)
		if err != nil {
			logger.WithError(err).Error("trying to retrieve cloud credentials")

			resource.Status.Conditions.SetCondition(core.Component{
				Name:    "provision",
				Message: "you do not have permission to the credentials",
				Status:  core.SuccessStatus,
			})

			return false, err
		}

		client, err := NewClient(creds, resource)
		if err != nil {
			logger.WithError(err).Error("attempting to create the cluster client")

			return false, err
		}

		logger.Info("checking if the cluster already exists")

		found, err := client.Exists()
		if err != nil {
			return false, err
		}

		if !found {
			// @step: we need to mark the cluster as pending
			if resource.Status.Conditions == nil {
				resource.Status.Conditions = &core.Components{}
			}

			status, found := resource.Status.Conditions.HasComponent("provision")
			if !found || status != core.PendingStatus {
				resource.Status.Conditions.SetCondition(core.Component{
					Name:    "provision",
					Message: "provisioning the cluster in google compute",
					Status:  core.PendingStatus,
				})
				resource.Status.Status = core.PendingStatus

				return true, nil
			}

			logger.Debug("creating a new gke cluster in gcp")
			if _, err = client.Create(ctx); err != nil {
				logger.WithError(err).Error("attempting to create cluster")

				resource.Status.Conditions.SetCondition(core.Component{
					Name:    "provision",
					Message: "failed trying to provision the cluster",
					Detail:  err.Error(),
				})
				resource.Status.Status = core.FailureStatus

				return false, err
			}
		} else {
			// else we are updating it
			if err := client.Update(ctx); err != nil {
				logger.WithError(err).Error("attempting to update cluster")

				return false, err
			}
		}

		// @step: retrieve the cluster spec
		cluster, found, err := client.GetCluster()
		if err != nil {
			logger.WithError(err).Error("trying to retrieve the gke cluster")

			return false, err
		}
		if !found {
			logger.Warn("gke cluster was not found")

			return true, nil
		}

		resource.Status.CACertificate = cluster.MasterAuth.ClusterCaCertificate
		resource.Status.Endpoint = fmt.Sprintf("https://%s", cluster.Endpoint)
		resource.Status.Status = core.SuccessStatus

		// @step: update the state as provisioned
		resource.Status.Conditions.SetCondition(core.Component{
			Name:    "provision",
			Message: "cluster has been provisioned",
			Status:  core.SuccessStatus,
		})

		// @step: set the bootstrap as pending if required
		status, found := resource.Status.Conditions.HasComponent("bootstrap")
		if !found || status != core.SuccessStatus {

			resource.Status.Conditions.SetCondition(core.Component{
				Name:    "bootstrap",
				Message: "bootstrapping the gke cluster",
				Status:  core.PendingStatus,
			})

			logger.Info("attempting to bootstrap the gke cluster")

			boot, err := NewBootstrapClient(resource, creds)
			if err != nil {
				logger.WithError(err).Error("trying to create bootstrap client")

				return false, err
			}
			if err := boot.Bootstrap(ctx, t.mgr.GetClient()); err != nil {
				logger.WithError(err).Error("trying to bootstrap gke cluster")

				return false, err
			}
		}

		resource.Status.Conditions.SetCondition(core.Component{
			Name:    "bootstrap",
			Message: "successfully bootstrapped the cluster",
			Status:  core.SuccessStatus,
		})

		return false, nil
	}()
	if err != nil {
		resource.Status.Status = core.FailureStatus
	}

	if err := t.mgr.GetClient().Status().Update(ctx, resource); err != nil {
		logger.WithError(err).Error("updating the status of gke cluster")

		return reconcile.Result{}, err
	}

	if err == nil {
		if finalizer.NeedToAdd(resource) {
			logger.Info("adding our finalizer to the team resource")

			if err := finalizer.Add(resource); err != nil {
				return reconcile.Result{}, err
			}

			return reconcile.Result{Requeue: true}, nil
		}
	}

	if requeue {
		return reconcile.Result{Requeue: true}, nil
	}

	return reconcile.Result{}, nil
}

// GetCredentials returns the cloud credential
func (t *gkeCtrl) GetCredentials(ctx context.Context, cluster *gke.GKE, team string) (*gke.GKECredentials, error) {
	// @step: is the team permitted access to this credentials
	permitted, err := t.Teams().Team(team).Allocations().IsPermitted(ctx, cluster.Spec.Credentials)
	if err != nil {
		log.WithError(err).Error("attempting to check for permission on credentials")

		return nil, fmt.Errorf("attempting to check for permission on credentials")
	}

	if !permitted {
		log.Warn("trying to build gke cluster unallocated permissions")

		return nil, errors.New("you do not have permissions to the gke credentials")
	}

	// @step: retrieve the credentials
	creds := &gke.GKECredentials{}

	return creds, t.mgr.GetClient().Get(ctx,
		types.NamespacedName{
			Namespace: cluster.Spec.Credentials.Namespace,
			Name:      cluster.Spec.Credentials.Name,
		}, creds,
	)
}
