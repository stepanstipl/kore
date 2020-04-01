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
	eksv1alpha1 "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	aws "github.com/appvia/kore/pkg/utils/cloud/aws"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	"github.com/aws/aws-sdk-go/service/eks"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "eksnodegroup.compute.kore.appvia.io"
	// ComponentClusterNodegroupCreator is the name of the component for the UI
	ComponentClusterNodegroupCreator = "Cluster Nodegroup Creator"
)

// Reconcile controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (n *eksNodeGroupCtrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to reconcile aws eks cluster node group")

	resource := &eksv1alpha1.EKSNodeGroup{}
	if err := n.mgr.GetClient().Get(ctx, request.NamespacedName, resource); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	finalizer := kubernetes.NewFinalizer(n.mgr.GetClient(), finalizerName)

	// @step: are we deleting the resource
	if finalizer.IsDeletionCandidate(resource) {
		return n.Delete(request)
	}
	// @step: we need to mark the cluster as pending
	if resource.Status.Conditions == nil {
		resource.Status.Conditions = &core.Components{}
	}

	requeue, err := func() (bool, error) {
		logger.Debug("retrieving the eks cluster credential")
		// @step: first we need to check if we have access to the credentials
		credentials, err := n.GetCredentials(ctx, resource, resource.Namespace)
		if err != nil {
			logger.WithError(err).Error("trying to retrieve cloud credentials")

			resource.Status.Conditions.SetCondition(core.Component{
				Name:    ComponentClusterNodegroupCreator,
				Message: "You do not have permission to the credentials",
				Status:  core.FailureStatus,
			})

			return false, err
		}
		logger.Info("Found EKSCredentials")

		client, err := aws.NewBasicClient(
			credentials,
			resource.Spec.ClusterName,
			resource.Spec.Region,
		)
		if err != nil {
			logger.WithError(err).Error("attempting to create the cluster client")

			resource.Status.Conditions.SetCondition(core.Component{
				Detail:  err.Error(),
				Name:    ComponentClusterNodegroupCreator,
				Message: "Failed to create EKS client, please check credentials",
				Status:  core.FailureStatus,
			})

			return false, err
		}
		logger.Info("Checking cluster nodegroup existence")
		found, err := client.NodeGroupExists(resource)
		if err != nil {
			resource.Status.Conditions.SetCondition(core.Component{
				Detail:  err.Error(),
				Name:    ComponentClusterNodegroupCreator,
				Message: "Failed to check for cluster nodegroup existence",
				Status:  core.FailureStatus,
			})

			return false, err
		}

		if !found {
			logger.Debug("Cheking cluster exists for nodegroup")

			// An ACTIVE cluster is a prerequisite for creating node groups
			// first check if CLUSTER exists
			clusterFound, err := client.Exists()
			if err != nil {
				logger.Debugf("error trying to check cluster exists %s for nodegroup %s", resource.Spec.ClusterName, resource.Name)

				resource.Status.Conditions.SetCondition(core.Component{
					Name:    ComponentClusterNodegroupCreator,
					Message: fmt.Sprintf("Waiting for cluster %s to exist in aws", resource.Spec.ClusterName),
					Status:  core.FailureStatus,
				})

				return false, err
			}
			if !clusterFound {
				logger.Debugf("waiting for cluster %s to exist", resource.Spec.ClusterName)

				resource.Status.Conditions.SetCondition(core.Component{
					Name:    ComponentClusterNodegroupCreator,
					Message: fmt.Sprintf("Waiting for cluster %s to exist in aws", resource.Spec.ClusterName),
					Status:  core.PendingStatus,
				})

				return true, nil
			}
			eksCluster, err := client.DescribeEKS()
			if err != nil {
				logger.Debugf("error trying to check status of cluster %s for nodegroup %s", resource.Spec.ClusterName, resource.Name)

				resource.Status.Conditions.SetCondition(core.Component{
					Name:    ComponentClusterNodegroupCreator,
					Message: fmt.Sprintf("Error checking for status of cluster %s", resource.Spec.ClusterName),
					Status:  core.FailureStatus,
				})

				return false, err
			}
			if *eksCluster.Status != eks.ClusterStatusActive {

				logger.Debugf("cluster status is %s", *eksCluster.Status)
				if *eksCluster.Status == eks.ClusterStatusCreating {
					resource.Status.Conditions.SetCondition(core.Component{
						Name:    ComponentClusterNodegroupCreator,
						Message: "Waiting for cluster to provision the EKS cluster nodegroup in AWS",
						Status:  core.PendingStatus,
					})
					resource.Status.Status = core.PendingStatus

					return true, nil
				}

				// Not active or creating problem
				em := fmt.Sprintf("bad status for cluster %s: %s", resource.Spec.ClusterName, *eksCluster.Status)
				resource.Status.Conditions.SetCondition(core.Component{
					Name:    ComponentClusterNodegroupCreator,
					Message: em,
					Status:  core.FailureStatus,
				})

				return false, errors.New(em)
			}

			status, found := resource.Status.Conditions.GetStatus(ComponentClusterNodegroupCreator)
			if !found || status != core.PendingStatus {
				resource.Status.Conditions.SetCondition(core.Component{
					Name:    ComponentClusterNodegroupCreator,
					Message: "Provisioning the EKS cluster nodegroup in AWS",
					Status:  core.PendingStatus,
				})
				resource.Status.Status = core.PendingStatus

				return true, nil
			}

			logger.Debug("creating a new eks cluster nodegroup in aws")
			if err := client.CreateNodeGroup(resource); err != nil {
				logger.WithError(err).Error("attempting to create cluster nodegroup")

				resource.Status.Conditions.SetCondition(core.Component{
					Name:    ComponentClusterNodegroupCreator,
					Message: "Failed trying to provision the cluster nodegroup",
					Detail:  err.Error(),
				})
				resource.Status.Status = core.FailureStatus

				return false, err
			}
		}

		// Get nodegroup status
		logger.Info("Checking the status of the node group: " + resource.Name)

		nodestatus, err := client.GetEKSNodeGroupStatus(resource)
		if err != nil {
			logger.WithError(err).Error("trying to retrieve the eks cluster nodegroup")

			return false, err
		}

		if nodestatus == awseks.NodegroupStatusCreateFailed {
			return false, fmt.Errorf("Cluster nodegroup has failed status:%s", resource.Name)
		}
		if nodestatus != awseks.NodegroupStatusActive {
			logger.Debugf("cluster nodegroup %s not ready requeing", resource.Name)

			// not ready, reque no errors
			return true, nil
		}
		logger.Info("Nodegroup active:" + resource.Name)
		// Set status to success
		// @step: update the state as provisioned
		resource.Status.Conditions.SetCondition(core.Component{
			Name:    ComponentClusterNodegroupCreator,
			Message: "Cluster nodegroup has been provisioned",
			Status:  core.SuccessStatus,
		})
		return false, nil
	}()
	if err != nil {
		resource.Status.Status = core.FailureStatus
	}

	if err := n.mgr.GetClient().Status().Update(ctx, resource); err != nil {
		logger.WithError(err).Error("updating the status of eks cluster nodegroup")

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
func (n *eksNodeGroupCtrl) GetCredentials(ctx context.Context, ng *eksv1alpha1.EKSNodeGroup, team string) (*eksv1alpha1.EKSCredentials, error) {
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
	creds := &eksv1alpha1.EKSCredentials{}

	return creds, n.mgr.GetClient().Get(
		ctx,
		types.NamespacedName{
			Namespace: ng.Spec.Credentials.Namespace,
			Name:      ng.Spec.Credentials.Name,
		},
		creds,
	)
}
