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
	"encoding/base64"
	"fmt"
	"time"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/cloud/aws"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// EnsureResourcePending ensures the resource is pending
func (t *eksCtrl) EnsureResourcePending(cluster *eks.EKS) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		if cluster.Status.Status != corev1.PendingStatus {
			cluster.Status.Status = corev1.PendingStatus

			return reconcile.Result{Requeue: true}, nil
		}

		return reconcile.Result{}, nil
	}
}

// EnsureClusterCreation is responsible for ensure the cluster is provision
func (t *eksCtrl) EnsureClusterCreation(client *aws.Client, cluster *eks.EKS) controllers.EnsureFunc {
	component := ComponentClusterCreator

	return func(ctx kore.Context) (reconcile.Result, error) {
		logger := log.WithFields(log.Fields{
			"name":      cluster.Name,
			"namespace": cluster.Namespace,
		})

		// @step: check if the cluster already exists
		exists, err := client.Exists(ctx)
		if err != nil {
			cluster.Status.Conditions.SetCondition(corev1.Component{
				Detail:  err.Error(),
				Name:    component,
				Message: "Failed to check for cluster existence",
				Status:  corev1.FailureStatus,
			})

			return reconcile.Result{}, err
		}

		if exists {
			return reconcile.Result{}, nil
		}

		logger.Debug("cluster does not exist, attempting to provision")

		// @step: ensure we update the status and the component
		status, found := cluster.Status.Conditions.GetStatus(component)
		if status != corev1.PendingStatus || !found {
			cluster.Status.Conditions.SetCondition(corev1.Component{
				Name:    component,
				Message: "Provisioning the EKS cluster in AWS",
				Status:  corev1.PendingStatus,
			})
			cluster.Status.Status = corev1.PendingStatus

			return reconcile.Result{Requeue: true}, nil
		}

		if err := client.Create(ctx); err != nil {
			logger.WithError(err).Error("failed to create the eks cluster")

			// The IAM role is not always available right away after creation and the API might return with the following error:
			// InvalidParameterException: Role with arn: <ARN> could not be assumed because it does not exist or the trusted entity is not correct
			// In this case we are going to retry and not throw an error
			if client.IsInvalidParameterException(err) {
				return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
			}

			cluster.Status.Conditions.SetCondition(corev1.Component{
				Name:    component,
				Message: "Failed to provision the EKS cluster",
				Detail:  err.Error(),
			})

			return reconcile.Result{}, err
		}

		return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
	}
}

// EnsureClusterInSync is responsible for ensuring the cluster is insync
func (t *eksCtrl) EnsureClusterInSync(client *aws.Client, cluster *eks.EKS) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		logger := log.WithFields(log.Fields{
			"name":      cluster.Name,
			"namespace": cluster.Namespace,
		})
		logger.Debug("attempting to check the eks cluster is in-sync")

		// @step: we retrieve the current state
		state, err := client.Describe(ctx)
		if err != nil {
			logger.WithError(err).Error("trying to describe the cluster")

			return reconcile.Result{}, err
		}
		status := utils.StringValue(state.Status)

		logger.WithField("status", status).Debug("current state of the eks cluster")

		switch status {
		case awseks.ClusterStatusActive:
			cluster.Status.Conditions.SetCondition(corev1.Component{
				Name:    ComponentClusterCreator,
				Message: "Cluster has been provisioned",
				Status:  corev1.SuccessStatus,
			})
			cluster.Status.Status = corev1.SuccessStatus

		case awseks.ClusterStatusFailed:
			cluster.Status.Conditions.SetCondition(corev1.Component{
				Name:    ComponentClusterCreator,
				Message: "EKS Cluster has failed to provision",
				Status:  corev1.FailureStatus,
			})
			cluster.Status.Status = corev1.FailureStatus

			return reconcile.Result{}, nil

		case awseks.ClusterStatusCreating, awseks.ClusterStatusUpdating:
			cluster.Status.Status = corev1.PendingStatus

			return reconcile.Result{RequeueAfter: 30 * time.Second}, nil

		default:
			return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
		}

		// @step: has the desired state drifted and if so was an update requested
		if needupdate, err := client.Update(ctx); err != nil {
			logger.WithError(err).Error("trying check or perform an update on the eks cluster")

			return reconcile.Result{}, err
		} else if needupdate {
			cluster.Status.Status = corev1.PendingStatus
			// we requeue and wait for the state to settle
			return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
		}

		// @step: ensure the eks cluster status is updated
		cadata := utils.StringValue(state.CertificateAuthority.Data)
		endpoint := utils.StringValue(state.Endpoint)

		ca, err := base64.StdEncoding.DecodeString(cadata)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("invalid base64 ca data from aws for eks endpoint %s,%v", endpoint, ca)
		}
		cluster.Status.CACertificate = string(ca)
		cluster.Status.Endpoint = endpoint
		cluster.Status.Status = corev1.SuccessStatus

		return reconcile.Result{}, nil
	}
}

// EnsureClusterBootstrap ensures the cluster is correctly bootstrapped
func (t *eksCtrl) EnsureClusterBootstrap(client *aws.Client, cluster *eks.EKS) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		logger := log.WithFields(log.Fields{
			"name":      cluster.Name,
			"namespace": cluster.Namespace,
		})
		logger.Debug("attempting to ensure the eks cluster is bootstrapped")

		// @step: set the bootstrap as pending if required
		cluster.Status.Conditions.SetCondition(corev1.Component{
			Name:    ComponentClusterBootstrap,
			Message: "Accessing the EKS cluster",
			Status:  corev1.PendingStatus,
		})

		boot, err := NewBootstrapClient(cluster, client.Sess)
		if err != nil {
			logger.WithError(err).Error("trying to create bootstrap client")

			return reconcile.Result{}, err
		}

		if err := boot.Run(ctx, ctx.Client()); err != nil {
			logger.WithError(err).Error("trying to bootstrap eks cluster")

			return reconcile.Result{}, err
		}

		cluster.Status.Conditions.SetCondition(corev1.Component{
			Name:    ComponentClusterBootstrap,
			Message: "Successfully initialized the EKS cluster",
			Status:  corev1.SuccessStatus,
		})

		return reconcile.Result{}, nil
	}
}

// EnsureClusterRoles ensures we have the cluster IAM roles
func (t *eksCtrl) EnsureClusterRoles(cluster *eks.EKS) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		logger := log.WithFields(log.Fields{
			"name":      cluster.Name,
			"namespace": cluster.Namespace,
		})
		logger.Debug("attempting to ensure the iam role for the eks cluster")

		// @step: first we need to check if we have access to the credentials
		creds, err := t.GetCredentials(ctx, cluster, cluster.Namespace)
		if err != nil {
			logger.WithError(err).Error("trying to retrieve cloud credentials")

			cluster.Status.Conditions.SetCondition(corev1.Component{
				Name:    ComponentClusterCreator,
				Message: "You do not have permission to the credentials",
				Status:  corev1.FailureStatus,
			})

			return reconcile.Result{}, err
		}

		// @step: we need to ensure the iam role for the cluster is there
		client := aws.NewIamClient(*creds)

		role, err := client.EnsureEKSClusterRole(ctx, cluster.Name)
		if err != nil {
			logger.WithError(err).Error("trying to ensure the eks iam role")

			cluster.Status.Conditions.SetCondition(corev1.Component{
				Name:    ComponentClusterCreator,
				Message: "Failed trying to provision the EKS Cluster Role",
				Detail:  err.Error(),
			})

			return reconcile.Result{}, err
		}

		cluster.Status.RoleARN = *role.Arn

		return reconcile.Result{}, nil
	}
}

// EnsureDeletionStatus ensures the resource is in a deleting state
func (t *eksCtrl) EnsureDeletionStatus(cluster *eks.EKS) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		// @step: lets update the status of the resource to deleting
		if cluster.Status.Status != corev1.DeletingStatus {
			cluster.Status.Status = corev1.DeletingStatus

			return reconcile.Result{Requeue: true}, nil
		}

		return reconcile.Result{}, nil
	}
}

// EnsureNodeGroupsDeleted ensures all nodegroup referencing me have been deleted
func (t *eksCtrl) EnsureNodeGroupsDeleted(cluster *eks.EKS) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		logger := log.WithFields(log.Fields{
			"name":      cluster.Name,
			"namespace": cluster.Namespace,
		})
		logger.Debug("ensuring all the eks nodegroups have been deleted")

		list := &eks.EKSNodeGroupList{}
		if err := t.mgr.GetClient().List(ctx, list, client.InNamespace(cluster.Namespace)); err != nil {
			logger.WithError(err).Error("trying to list the eks nodegroups")

			return reconcile.Result{}, err
		}

		found := func() bool {
			for _, x := range list.Items {
				if kore.IsOwner(cluster, x.Spec.Cluster) {
					return true
				}
			}

			return false
		}()

		// @step: if we found nodegroup we should not delete, but requeue
		if found {
			return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
		}

		return reconcile.Result{}, nil
	}
}

// EnsureDeletion is responsible for deleting the actual cluster
func (t *eksCtrl) EnsureDeletion(cluster *eks.EKS) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		logger := log.WithFields(log.Fields{
			"name":      cluster.Name,
			"namespace": cluster.Namespace,
		})
		logger.Debug("attempting to delete the eks cluster")

		// @step: retrieve the cloud credentials
		creds, err := t.GetCredentials(ctx, cluster, cluster.Namespace)
		if err != nil {
			return reconcile.Result{}, err
		}

		// @step: create a cloud client for us
		client, err := aws.NewEKSClient(creds, cluster)
		if err != nil {
			logger.WithError(err).Error("trying to create eks client")

			return reconcile.Result{}, err
		}

		// @step: check if the cluster exists
		found, err := client.Exists(ctx)
		if err != nil {
			logger.WithError(err).Error("trying to check if eks cluster exists")

			return reconcile.Result{}, err
		}
		if !found {
			// we can exis the loop here else we need to requeue or error
			return reconcile.Result{}, nil
		}

		logger.Debug("eks cluster exists, attempting to delete now")

		// @step: get the current state of the cluster
		state, err := client.Describe(ctx)
		if err != nil {
			logger.WithError(err).Error("trying to describe the cluster")

			return reconcile.Result{}, err
		}

		status := utils.StringValue(state.Status)
		logger.WithField("status", status).Debug("current state of the eks cluster")

		// @step: if the cluster is not deleting, try and delete now
		switch status {
		case awseks.ClusterStatusActive, awseks.ClusterStatusFailed:
			if err := client.Delete(ctx); err != nil {
				logger.WithError(err).Error("trying to delete the eks cluster from team")

				return reconcile.Result{}, err
			}
			cluster.Status.Status = corev1.DeletingStatus
			return reconcile.Result{RequeueAfter: 30 * time.Second}, nil

		case awseks.ClusterStatusCreating, awseks.ClusterStatusUpdating:
			cluster.Status.Status = corev1.PendingStatus
			return reconcile.Result{RequeueAfter: 30 * time.Second}, nil

		case awseks.ClusterStatusDeleting:
			cluster.Status.Status = corev1.DeletingStatus
			return reconcile.Result{RequeueAfter: 30 * time.Second}, nil

		default:
			return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
		}
	}
}

// EnsureRoleDeletion is responsible for deleting the IAM roles
func (t *eksCtrl) EnsureRoleDeletion(cluster *eks.EKS) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		logger := log.WithFields(log.Fields{
			"name":      cluster.Name,
			"namespace": cluster.Namespace,
		})
		logger.Debug("attempting to delete the eks cluster role")

		credentials, err := t.GetCredentials(ctx, cluster, cluster.Namespace)
		if err != nil {
			return reconcile.Result{}, err
		}
		client := aws.NewIamClient(*credentials)

		if err := client.DeleteEKSClutserRole(ctx, cluster.Name); err != nil {
			logger.WithError(err).Error("trying to delete the eks cluster role")

			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}
}

// EnsureSecretDeletion ensure the cluster secret is removed
func (t *eksCtrl) EnsureSecretDeletion(cluster *eks.EKS) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		// @step: we can now delete the sysadmin token now
		if err := controllers.DeleteClusterCredentialsSecret(ctx,
			t.mgr.GetClient(), cluster.Namespace, cluster.Name); err != nil {

			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}
}

// EnsureFinalizerRemoved removes the finalizer now
func (t *eksCtrl) EnsureFinalizerRemoved(cluster *eks.EKS) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		finalizer := kubernetes.NewFinalizer(ctx.Client(), finalizerName)
		if finalizer.IsDeletionCandidate(cluster) {
			if err := finalizer.Remove(cluster); err != nil {
				return reconcile.Result{}, err
			}
		}

		return reconcile.Result{}, nil
	}
}
