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
	"encoding/base64"
	"fmt"
	"time"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/cloud/aws"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// EnsureResourcePending ensures the resource is pending
func (t *eksCtrl) EnsureResourcePending(cluster *eks.EKS) controllers.EnsureFunc {
	return func(ctx context.Context) (reconcile.Result, error) {
		if cluster.Status.Status != corev1.PendingStatus {
			cluster.Status.Status = corev1.PendingStatus

			return reconcile.Result{Requeue: true}, nil
		}

		return reconcile.Result{}, nil
	}
}

// EnsureCluster is responsible for creating the cluster
func (t *eksCtrl) EnsureCluster(client *aws.Client, cluster *eks.EKS) controllers.EnsureFunc {
	return func(ctx context.Context) (reconcile.Result, error) {
		logger := log.WithFields(log.Fields{
			"name":      cluster.Name,
			"namespace": cluster.Namespace,
		})
		logger.Debug("attempting to ensure the eks cluster")

		// @step: check if the cluster already exists
		existing, err := client.Exists(ctx)
		if err != nil {
			cluster.Status.Conditions.SetCondition(corev1.Component{
				Detail:  err.Error(),
				Name:    ComponentClusterCreator,
				Message: "Failed to check for cluster existence",
				Status:  corev1.FailureStatus,
			})

			return reconcile.Result{}, err
		}
		if !existing {
			// @step: ensure we update the status and the component
			status, found := cluster.Status.Conditions.GetStatus(ComponentClusterCreator)
			if !found || status != corev1.PendingStatus {
				cluster.Status.Conditions.SetCondition(corev1.Component{
					Name:    ComponentClusterCreator,
					Message: "Provisioning the EKS cluster in AWS",
					Status:  corev1.PendingStatus,
				})
				cluster.Status.Status = corev1.PendingStatus

				return reconcile.Result{Requeue: true}, nil
			}

			if err := client.Create(ctx); err != nil {
				logger.WithError(err).Error("failed to create cluster")

				// The IAM role is not always available right away after creation and the API might return with the following error:
				// InvalidParameterException: Role with arn: <ARN> could not be assumed because it does not exist or the trusted entity is not correct
				// In this case we are going to retry and not throw an error
				if client.IsInvalidParameterException(err) {
					return reconcile.Result{RequeueAfter: 5 * time.Second}, nil
				}

				cluster.Status.Conditions.SetCondition(corev1.Component{
					Name:    ComponentClusterCreator,
					Message: "Failed trying to provision the cluster, will retry",
					Detail:  err.Error(),
				})

				return reconcile.Result{}, err
			}
		} else {
			// @step: we need to sure the cluster is active
			if err := client.WaitForClusterReady(ctx); err != nil {
				logger.WithError(err).Error("trying to ensure the cluster is active")

				return reconcile.Result{}, err
			}

			// TODO - client needs to manage migrations
		}

		// @step: update the state as provisioned
		cluster.Status.Conditions.SetCondition(corev1.Component{
			Name:    ComponentClusterCreator,
			Message: "Cluster has been provisioned",
			Status:  corev1.SuccessStatus,
		})

		// @step: retrieve and check the status of the cluster
		c, err := client.Describe(ctx)
		if err != nil {
			logger.WithError(err).Error("trying to retrieve the eks cluster")

			return reconcile.Result{}, err
		}

		// Active cluster
		ca, err := base64.StdEncoding.DecodeString(*c.CertificateAuthority.Data)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("invalid base64 ca data from aws for eks endpoint %s,%v", *c.Endpoint, c.CertificateAuthority.Data)
		}
		cluster.Status.CACertificate = string(ca)
		cluster.Status.Endpoint = *c.Endpoint
		cluster.Status.Status = corev1.SuccessStatus

		return reconcile.Result{}, nil
	}
}

// EnsureClusterBootstrap ensures the cluster is correctly bootstrapped
func (t *eksCtrl) EnsureClusterBootstrap(client *aws.Client, cluster *eks.EKS) controllers.EnsureFunc {
	return func(ctx context.Context) (reconcile.Result, error) {
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

		if err := boot.Run(ctx, t.mgr.GetClient()); err != nil {
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
	return func(ctx context.Context) (reconcile.Result, error) {
		logger := log.WithFields(log.Fields{
			"name":      cluster.Name,
			"namespace": cluster.Namespace,
		})
		logger.Debug("attempting to ensure the iam role for the eks cluster")

		// @step: first we need to check if we have access to the credentials
		credentials, err := t.GetCredentials(ctx, cluster, cluster.Namespace)
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
		client := aws.NewIamClient(aws.Credentials{
			AccessKeyID:     credentials.Spec.AccessKeyID,
			SecretAccessKey: credentials.Spec.SecretAccessKey,
		})

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
	return func(ctx context.Context) (reconcile.Result, error) {
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
	return func(ctx context.Context) (reconcile.Result, error) {
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
	return func(ctx context.Context) (reconcile.Result, error) {
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
		if found {
			logger.Debug("eks cluster exists, attempting to delete now")

			if err := client.Delete(ctx); err != nil {
				logger.WithError(err).Error("trying to delete the eks cluster from team")

				return reconcile.Result{}, err
			}
		}

		return reconcile.Result{}, nil
	}
}

// EnsureSecretDeletion ensure the cluster secret is removed
func (t *eksCtrl) EnsureSecretDeletion(cluster *eks.EKS) controllers.EnsureFunc {
	return func(ctx context.Context) (reconcile.Result, error) {
		// @step: we can now delete the sysadmin token now
		if err := controllers.DeleteClusterCredentialsSecret(ctx,
			t.mgr.GetClient(), cluster.Namespace, cluster.Name); err != nil {

			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}
}
