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
	"errors"
	"time"

	"github.com/appvia/kore/pkg/kore"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	ekscc "github.com/appvia/kore/pkg/controllers/cloud/aws/eks"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/cloud/aws"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// EnsureNodeGroupIsPending is responsible for setting the resource to a pending state
func (n *ctrl) EnsureNodeGroupIsPending(group *eks.EKSNodeGroup) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		if group.Status.Status != corev1.PendingStatus {
			group.Status.Status = corev1.PendingStatus

			return reconcile.Result{Requeue: true}, nil
		}

		return reconcile.Result{}, nil
	}
}

// EnsureClusterReady is responsible for checking the EKS cluster is ready
func (n *ctrl) EnsureClusterReady(group *eks.EKSNodeGroup) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {

		logger := log.WithFields(log.Fields{
			"name":      group.Name,
			"namespace": group.Namespace,
		})
		logger.Debug("attempting to ensure the eks cluster is ready")

		key := types.NamespacedName{
			Name:      group.Spec.Cluster.Name,
			Namespace: group.Spec.Cluster.Namespace,
		}

		cluster := &eks.EKS{}
		if err := ctx.Client().Get(ctx, key, cluster); err != nil {
			logger.WithError(err).Error("trying to retrieve the cluster status")

			return reconcile.Result{}, err
		}

		status, found := cluster.Status.Conditions.GetStatus(ekscc.ComponentClusterCreator)
		if !found || status != corev1.SuccessStatus {
			logger.Warn("eks cluster not ready yet, we will wait")

			return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
		}

		return reconcile.Result{}, nil
	}
}

// EnsureNodeRole is responsible for ensuring the IAM role is there
func (n *ctrl) EnsureNodeRole(group *eks.EKSNodeGroup, credentials *aws.Credentials) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		logger := log.WithFields(log.Fields{
			"name":      group.Name,
			"namespace": group.Namespace,
		})
		logger.Debug("attempting to ensure the node iam role")

		client := aws.NewIamClient(*credentials)

		role, err := client.EnsureEKSNodePoolRole(ctx, group.Name)
		if err != nil {
			group.Status.Conditions.SetCondition(corev1.Component{
				Name:    ComponentClusterNodegroupCreator,
				Message: "Failed trying to provision the eks nodepool iam roles",
				Detail:  err.Error(),
			})

			return reconcile.Result{}, err
		}

		// Save the role used for this cluster
		group.Status.NodeIAMRole = *role.Arn

		return reconcile.Result{}, nil
	}
}

// EnsureNodeGroup is responsible for making sure the nodegroup is provisioned
func (n *ctrl) EnsureNodeGroup(client *aws.Client, group *eks.EKSNodeGroup) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		logger := log.WithFields(log.Fields{
			"name":      group.Name,
			"namespace": group.Namespace,
		})
		logger.Debug("attempting to ensure the eks nodegroup")

		found, err := client.NodeGroupExists(ctx, group)
		if err != nil {
			logger.WithError(err).Error("trying to describe the nodegroup")

			return reconcile.Result{}, err
		}

		if !found {
			logger.Debug("eks nodegroup does not exist, attempting to create now")

			// @step: we set the provisioning state to pending
			status, found := group.Status.Conditions.GetStatus(ComponentClusterNodegroupCreator)
			if !found || status != corev1.PendingStatus {
				group.Status.Conditions.SetCondition(corev1.Component{
					Name:    ComponentClusterNodegroupCreator,
					Message: "Provisioning the EKS cluster nodegroup in AWS",
					Status:  corev1.PendingStatus,
				})

				return reconcile.Result{Requeue: true}, nil
			}

			// @step: attempt to ensure create the eks nodegroup
			if err := client.CreateNodeGroup(ctx, group); err != nil {
				logger.WithError(err).Error("trying to create the eks nodegroup")

				return reconcile.Result{}, err
			}

			return reconcile.Result{RequeueAfter: 60 * time.Second}, nil
		}

		// @step: get the current state of the nodegroup
		state, err := client.DescribeNodeGroup(ctx, group)
		if err != nil {
			log.WithError(err).Error("trying to describe the nodegroup")

			return reconcile.Result{}, err
		}

		logger.WithField("status", utils.StringValue(state.Status)).Debug("state of the eks nodegroup at the moment")

		switch utils.StringValue(state.Status) {
		case awseks.NodegroupStatusCreateFailed:
			group.Status.Status = corev1.FailureStatus

			group.Status.Conditions.SetCondition(corev1.Component{
				Name:    ComponentClusterNodegroupCreator,
				Message: "EKS Nodegroup has failed to create",
				Status:  corev1.FailureStatus,
			})

		case awseks.NodegroupStatusCreating, awseks.NodegroupStatusUpdating:
			group.Status.Status = corev1.PendingStatus

			return reconcile.Result{RequeueAfter: 30 * time.Second}, nil

		case awseks.NodegroupStatusDeleting:
			logger.Warn("nodegroup is defined as deleting, someone has removed outside of kore")
			// not sure what to do here?? - we should probaby way and recreate it HAHAHA :-)

			return reconcile.Result{RequeueAfter: 30 * time.Second}, nil

		case awseks.NodegroupStatusActive, awseks.NodegroupStatusDegraded:
			group.Status.Status = corev1.SuccessStatus

			group.Status.Conditions.SetCondition(corev1.Component{
				Name:    ComponentClusterNodegroupCreator,
				Message: "EKS Nodegroup has been provisioned",
				Status:  corev1.SuccessStatus,
			})

			if state.Resources != nil {
				group.Status.AutoScalingGroupNames = nil
				for _, asg := range state.Resources.AutoScalingGroups {
					group.Status.AutoScalingGroupNames = append(group.Status.AutoScalingGroupNames, *asg.Name)
				}
			}
		}

		updating, err := client.UpdateNodeGroup(ctx, group)
		if err != nil {
			logger.WithError(err).Error("trying to update the eks nodegroup")

			return reconcile.Result{}, err
		}
		if updating {
			logger.Debug("performing an update on the eks nodegroup")
			group.Status.Status = corev1.PendingStatus

			return reconcile.Result{RequeueAfter: 30 * time.Second}, err
		}

		return reconcile.Result{}, nil
	}
}

// EnsureDeletionStatus makes sure the resource is set to deleting
func (n *ctrl) EnsureDeletionStatus(group *eks.EKSNodeGroup) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		if group.Status.Status != corev1.DeletingStatus {
			group.Status.Status = corev1.DeletingStatus

			return reconcile.Result{Requeue: true}, nil
		}

		return reconcile.Result{}, nil
	}
}

// EnsureDeletion ensures the nodegroup is deleting
func (n *ctrl) EnsureDeletion(group *eks.EKSNodeGroup) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		logger := log.WithFields(log.Fields{
			"name":      group.Name,
			"namespace": group.Namespace,
		})
		logger.Debug("attempting to delete eks nodegroup")

		creds, err := n.GetCredentials(ctx, group, group.Namespace)
		if err != nil {
			return reconcile.Result{}, err
		}

		// @step: create a cloud client for us
		client, err := aws.NewBasicClient(creds, group.Spec.Cluster.Name, group.Spec.Region)
		if err != nil {
			log.WithError(err).Error("trying to create a aws client for the nodegroup")

			return reconcile.Result{}, err
		}

		// @step: check if the nodegroup exists and if so we wait or the operation or the exit
		found, err := client.NodeGroupExists(ctx, group)
		if err != nil {
			log.WithError(err).Error("trying to check if nodegroup exists")

			return reconcile.Result{}, err
		}
		if !found {
			return reconcile.Result{}, nil
		}

		state, err := client.DescribeNodeGroup(ctx, group)
		if err != nil {
			logger.WithError(err).Error("trying to describe the nodegroup state")

			return reconcile.Result{}, err
		}

		status := utils.StringValue(state.Status)

		switch status {

		case awseks.NodegroupStatusActive, awseks.NodegroupStatusCreateFailed, awseks.NodegroupStatusDegraded:
			if err := client.DeleteNodeGroup(ctx, group); err != nil {
				logger.WithError(err).Error("trying to delete the nodegroup")

				return reconcile.Result{}, err
			}
			return reconcile.Result{RequeueAfter: 30 * time.Second}, nil

		case awseks.NodegroupStatusCreating, awseks.NodegroupStatusUpdating, awseks.NodegroupStatusDeleting:
			return reconcile.Result{RequeueAfter: 30 * time.Second}, nil

		case awseks.NodegroupStatusDeleteFailed:
			return reconcile.Result{}, errors.New("nodegroup has failed to delete, please check console")
		}

		return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
	}
}

// EnsureRoleDeletion ensures the nodegroup role is removed
func (n *ctrl) EnsureRoleDeletion(group *eks.EKSNodeGroup) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		logger := log.WithFields(log.Fields{
			"name":      group.Name,
			"namespace": group.Namespace,
		})
		logger.Debug("attempting to delete eks nodegroup role")

		credentials, err := n.GetCredentials(ctx, group, group.Namespace)
		if err != nil {
			return reconcile.Result{}, err
		}
		client := aws.NewIamClient(*credentials)

		if err := client.DeleteEKSNodeGroupRole(ctx, group.Name); err != nil {
			logger.WithError(err).Error("trying to remove the nodegroup role")

			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}
}

func (n *ctrl) EnsureRemoveFinalizer(group *eks.EKSNodeGroup) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		logger := log.WithFields(log.Fields{
			"name":      group.Name,
			"namespace": group.Namespace,
		})
		logger.Debug("attempting to delete eks nodegroup finalizer")

		finalizer := kubernetes.NewFinalizer(ctx.Client(), finalizerName)
		if err := finalizer.Remove(group); err != nil {
			logger.WithError(err).Error("removing the finalizer from eks resource")

			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}
}
