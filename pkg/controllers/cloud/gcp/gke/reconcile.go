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

package gke

import (
	"context"
	"errors"
	"fmt"
	"time"

	config "github.com/appvia/kore/pkg/apis/config/v1"
	core "github.com/appvia/kore/pkg/apis/core/v1"
	gcp "github.com/appvia/kore/pkg/apis/gcp/v1alpha1"
	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	gcpcc "github.com/appvia/kore/pkg/controllers/cloud/gcp/projectclaim"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "gke.compute.kore.appvia.io"
	// ComponentClusterCreator is the name of the component for the UI
	ComponentClusterCreator = "Cluster Creator"
	// ComponentClusterBootstrap is the component name for seting up cloud credentials
	ComponentClusterBootstrap = "Cluster Initialize Access"
)

// Reconcile is the entrypoint for the reconciliation logic
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
	original := resource.DeepCopyObject()

	finalizer := kubernetes.NewFinalizer(t.mgr.GetClient(), finalizerName)

	// @step: are we deleting the resource
	if finalizer.IsDeletionCandidate(resource) {
		return t.Delete(request)
	}

	// @step: we need to mark the cluster as pending
	if resource.Status.Conditions == nil {
		resource.Status.Conditions = core.Components{}
	}

	result, err := func() (reconcile.Result, error) {
		logger.Debug("retrieving the gke cluster credential")
		// @step: first we need to check if we have access to the credentials
		creds, err := t.GetCredentials(ctx, resource, resource.Namespace)
		if err != nil {
			logger.WithError(err).Error("trying to retrieve cloud credentials")

			resource.Status.Conditions.SetCondition(core.Component{
				Name:    ComponentClusterCreator,
				Message: "You do not have permission to the credentials",
				Status:  core.FailureStatus,
			})

			return reconcile.Result{}, err
		}

		client, err := NewClient(creds, resource)
		if err != nil {
			logger.WithError(err).Error("attempting to create the cluster client")

			resource.Status.Conditions.SetCondition(core.Component{
				Detail:  err.Error(),
				Name:    ComponentClusterCreator,
				Message: "Failed to create GKE client, please check credentials",
				Status:  core.FailureStatus,
			})

			return reconcile.Result{}, err
		}
		logger.Info("checking if the cluster already exists")

		found, err := client.Exists(ctx)
		if err != nil {
			resource.Status.Conditions.SetCondition(core.Component{
				Detail:  err.Error(),
				Name:    ComponentClusterCreator,
				Message: "Failed to check for cluster existence",
				Status:  core.FailureStatus,
			})

			return reconcile.Result{}, err
		}

		if !found {
			status, found := resource.Status.Conditions.GetStatus(ComponentClusterCreator)
			if !found || status != core.PendingStatus {
				resource.Status.Conditions.SetCondition(core.Component{
					Name:    ComponentClusterCreator,
					Message: "Provisioning the cluster in google compute",
					Status:  core.PendingStatus,
				})
				resource.Status.Status = core.PendingStatus

				return reconcile.Result{Requeue: true}, nil
			}

			logger.Debug("creating a new gke cluster in gcp")

			if _, err = client.Create(ctx); err != nil {
				logger.WithError(err).Error("attempting to create cluster")

				resource.Status.Conditions.SetCondition(core.Component{
					Name:    ComponentClusterCreator,
					Message: "Failed trying to provision the cluster",
					Detail:  err.Error(),
					Status:  core.FailureStatus,
				})
				resource.Status.Status = core.FailureStatus

				return reconcile.Result{}, err
			}

			return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
		}

		cluster, _, err := client.GetCluster(ctx)
		if err != nil {
			logger.WithError(err).Error("trying to retrieve the cluster status")

			return reconcile.Result{}, err
		}
		logger.WithField("status", cluster.Status).Debug("the current state of the gke cluster is")

		switch cluster.Status {
		case "PROVISIONING":
			resource.Status.Status = core.PendingStatus
			return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
		case "RECONCILING":
			resource.Status.Status = core.PendingStatus
			return reconcile.Result{RequeueAfter: 60 * time.Second}, nil
		case "STOPPING":
			resource.Status.Status = core.DeletingStatus
			return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
		case "ERROR":
			resource.Status.Status = core.FailureStatus
			// @choice .. allowing it to run though here as it's in error but might require
			// an update to fix
		case "RUNNING":
		default:
			logger.Warn("cluster is in an unknown state, choosing to requeue instead")

			return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
		}

		// @step: we check the current state against the desired and see if we need to amend
		updating, err := client.Update(ctx)
		if err != nil {
			logger.WithError(err).Error("attempting to update cluster")

			return reconcile.Result{}, err
		}
		if updating {
			logger.Debug("cluster is performing a update")

			resource.Status.Status = core.PendingStatus

			return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
		}

		// @step: enable the cloud-nat is required
		if resource.Spec.EnablePrivateNetwork {
			logger.Info("cluster has private networking enabled, ensuring a cloud-nat")
			if err := client.EnableCloudNAT(); err != nil {
				logger.WithError(err).Error("trying to ensure the cloud-nat device")

				return reconcile.Result{}, err
			}

			logger.Info("cluster has private networking enabled, ensuring a api firewall rules")

			if err := client.EnableFirewallAPIServices(); err != nil {
				logger.WithError(err).Error("creating firewall rules for api extensions")

				return reconcile.Result{}, err
			}
		}

		resource.Status.CACertificate = cluster.MasterAuth.ClusterCaCertificate
		resource.Status.Endpoint = fmt.Sprintf("https://%s", cluster.Endpoint)
		resource.Status.Status = core.SuccessStatus

		// @step: update the state as provisioned
		resource.Status.Conditions.SetCondition(core.Component{
			Name:    ComponentClusterCreator,
			Message: "Cluster has been provisioned",
			Status:  core.SuccessStatus,
		})

		// @step: set the bootstrap as pending if required
		resource.Status.Conditions.SetCondition(core.Component{
			Name:    ComponentClusterBootstrap,
			Message: "Accessing the gke cluster",
			Status:  core.PendingStatus,
		})

		logger.Info("attempting to bootstrap the gke cluster")

		bc, err := newBootstrapClient(resource, creds)
		if err != nil {
			logger.WithError(err).Error("trying to create bootstrap client")

			return reconcile.Result{}, err
		}

		if err := controllers.NewBootstrap(bc).Run(ctx, t.mgr.GetClient()); err != nil {
			logger.WithError(err).Error("trying to bootstrap gke cluster")

			return reconcile.Result{}, err
		}

		resource.Status.Conditions.SetCondition(core.Component{
			Name:    ComponentClusterBootstrap,
			Message: "Successfully initialized the cluster",
			Status:  core.SuccessStatus,
		})

		return reconcile.Result{}, nil
	}()
	if err != nil {
		resource.Status.Status = core.FailureStatus
	}

	if err := t.mgr.GetClient().Status().Patch(ctx, resource, client.MergeFrom(original)); err != nil {
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

	return result, err
}

// GetCredentials returns the cloud credential
func (t *gkeCtrl) GetCredentials(ctx context.Context, cluster *gke.GKE, team string) (*credentials, error) {
	// @step: is the team permitted access to this credentials
	permitted, err := t.Teams().Team(team).Allocations().IsPermitted(ctx, cluster.Spec.Credentials)
	if err != nil {
		log.WithError(err).Error("attempting to check for permission on credentials")

		return nil, fmt.Errorf("attempting to check for permission on credentials")
	}

	if !permitted {
		log.Warn("trying to build gke cluster unallocated permissions")

		return nil, errors.New("you do not have permissions to the credentials")
	}

	key := types.NamespacedName{
		Namespace: cluster.Spec.Credentials.Namespace,
		Name:      cluster.Spec.Credentials.Name,
	}

	// @step: are we building the cluster off a project claim or static credentials
	switch cluster.Spec.Credentials.Group {
	case gke.SchemeGroupVersion.Group:
		switch kind := cluster.Spec.Credentials.Kind; kind {
		case "GKECredentials":
			return t.GetGKECredentials(ctx, key)
		default:
			return nil, fmt.Errorf("unknown gke credential kind: %s", kind)
		}

	case gcp.SchemeGroupVersion.Group:
		switch kind := cluster.Spec.Credentials.Kind; kind {
		case "ProjectClaim":
			return t.GetProjectClaimCredentials(ctx, key)
		default:
			return nil, fmt.Errorf("unknown gcp credential kind: %s", kind)
		}
	}

	return nil, errors.New("unknown credentials api group specified")
}

// GetGKECredentials is responsible for pulling the gke credentias type
func (t *gkeCtrl) GetGKECredentials(ctx context.Context, key types.NamespacedName) (*credentials, error) {
	c := &gke.GKECredentials{
		ObjectMeta: metav1.ObjectMeta{
			Name:      key.Name,
			Namespace: key.Namespace,
		},
	}
	found, err := kubernetes.GetIfExists(ctx, t.mgr.GetClient(), c)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("gke credentials: (%s/%s) not found", key.Namespace, key.Namespace)
	}

	if c.Status.Verified == nil || !*c.Status.Verified {
		return nil, errors.New("gke credentials have failed validation, please check credentials")
	}

	// for backwards-compatibility, use the key (Account) set on the GKECredentials resource, if it exists
	if c.Spec.Account != "" {
		return &credentials{
			key:        c.Spec.Account,
			project_id: c.Spec.Project,
			project:    c.Spec.Project,
			region:     c.Spec.Region,
		}, nil
	}

	// @step: we need to grab the secret
	secret := &config.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.Spec.CredentialsRef.Name,
			Namespace: c.Spec.CredentialsRef.Namespace,
		},
	}

	found, err = kubernetes.GetIfExists(ctx, t.mgr.GetClient(), secret)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("gke credentials secret is missing")
	}

	// @step: ensure the secret is decoded before using
	if err := secret.Decode(); err != nil {
		return nil, err
	}

	return &credentials{
		key:        secret.Spec.Data["key"],
		project_id: c.Spec.Project,
		project:    c.Spec.Project,
		region:     c.Spec.Region,
	}, nil
}

// GetProjectClaimCredentials is responsible for retrieving the project claim secret
// probably needs to be moved into a common lib
func (t *gkeCtrl) GetProjectClaimCredentials(ctx context.Context, key types.NamespacedName) (*credentials, error) {
	c := &gcp.ProjectClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      key.Name,
			Namespace: key.Namespace,
		},
	}

	found, err := kubernetes.GetIfExists(ctx, t.mgr.GetClient(), c)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("gcp project claim: (%s/%s) not found", key.Namespace, key.Namespace)
	}

	// @step: we need to check the status of the project
	if c.Status.Status == core.FailureStatus {
		return nil, errors.New("gcp project is in a failed state")
	}
	if c.Status.CredentialRef == nil {
		return nil, errors.New("no gcp credentials reference on project claim")
	}
	if c.Status.CredentialRef.Name == "" || c.Status.CredentialRef.Namespace == "" {
		return nil, errors.New("gcp project claims credentials reference is missing fields")
	}

	// @step: we need to grab the secret
	secret := &config.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.Status.CredentialRef.Name,
			Namespace: c.Status.CredentialRef.Namespace,
		},
	}

	found, err = kubernetes.GetIfExists(ctx, t.mgr.GetClient(), secret)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("gcp project secret is missing")
	}

	// @step: ensure the secret is decoded before using
	if err := secret.Decode(); err != nil {
		return nil, err
	}

	return &credentials{
		key:        secret.Spec.Data[gcpcc.ServiceAccountKey],
		project_id: secret.Spec.Data[gcpcc.ProjectIDKey],
		project:    secret.Spec.Data[gcpcc.ProjectNameKey],
	}, nil
}
