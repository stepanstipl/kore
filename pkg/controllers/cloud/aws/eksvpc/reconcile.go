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

package eksvpc

import (
	"context"
	"errors"
	"fmt"

	core "github.com/appvia/kore/pkg/apis/core/v1"
	eksv1alpha1 "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/appvia/kore/pkg/utils/cloud/aws"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "eksvpc.compute.kore.appvia.io"
	// ComponentVPCCreator is the name of the component for the UI
	ComponentVPCCreator = "Cluster VPC Creator"
)

func (t *eksvpcCtrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to reconcile aws vpc for an eks cluster")

	resource := &eksv1alpha1.EKSVPC{}
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
	// @step: we need to mark the cluster as pending
	if resource.Status.Conditions == nil {
		resource.Status.Conditions = core.Components{}
	}

	requeue, err := func() (bool, error) {

		logger.Debug("retrieving the vpc credentials")
		// @step: first we need to check if we have access to the credentials
		credentials, err := t.GetCredentials(ctx, resource, resource.Namespace)
		if err != nil {
			logger.WithError(err).Error("trying to retrieve cloud credentials")

			resource.Status.Conditions.SetCondition(core.Component{
				Name:    ComponentVPCCreator,
				Message: "You do not have permission to the credentials",
				Status:  core.FailureStatus,
			})

			return false, err
		}
		logger.Info("Found EKSCredential")

		client, err := aws.NewVPCClient(aws.Credentials{
			AccessKeyID:     credentials.Spec.AccessKeyID,
			SecretAccessKey: credentials.Spec.SecretAccessKey,
		}, aws.VPC{
			CidrBlock: resource.Spec.PrivateIPV4Cidr,
			Name:      resource.Name,
			Region:    resource.Spec.Region,
			Tags: map[string]string{
				"Kore": "managed",
			},
		})
		if err != nil {
			resource.Status.Conditions.SetCondition(core.Component{
				Detail:  err.Error(),
				Name:    ComponentVPCCreator,
				Message: fmt.Sprintf("Failed to create vpc as input values invalid - %s", err),
				Status:  core.FailureStatus,
			})

			return false, err
		}
		logger.Infof("Checking vpc existence %s", resource.Name)

		found, err := client.Exists()
		if err != nil {
			resource.Status.Conditions.SetCondition(core.Component{
				Detail:  err.Error(),
				Name:    ComponentVPCCreator,
				Message: "Failed to check for vpc existence",
				Status:  core.FailureStatus,
			})

			return false, err
		}

		if !found {
			status, found := resource.Status.Conditions.GetStatus(ComponentVPCCreator)
			if !found || status != core.PendingStatus {
				resource.Status.Conditions.SetCondition(core.Component{
					Name:    ComponentVPCCreator,
					Message: "Provisioning the VPC in AWS",
					Status:  core.PendingStatus,
				})
				resource.Status.Status = core.PendingStatus

				return true, nil
			}

			logger.Debug("creating or discovering a vpc in aws")
		}
		// Ensure this only reports if it exists when all resources exist or ensure update works
		// TODO: probably need a signal to only verify this every x...
		if err = client.Ensure(); err != nil {
			logger.WithError(err).Error("attempting to create or discover vpc")

			resource.Status.Conditions.SetCondition(core.Component{
				Name:    ComponentVPCCreator,
				Message: "Failed trying to provision the cluster",
				Detail:  err.Error(),
			})
			resource.Status.Status = core.FailureStatus

			return false, err
		}
		resource.Status.Infra.SecurityGroupIDs = []string{client.VPC.ControlPlaneSecurityGroupID}
		subnets := []string{}
		subnets = append(subnets, client.VPC.PublicSubnetIDs...)
		subnets = append(subnets, client.VPC.PrivateSubnetIDs...)
		resource.Status.Infra.SubnetIDs = subnets
		resource.Status.Infra.PublicIPV4EgressAddresses = client.VPC.PublicIPV4EgressAddresses

		// @step: update the state as provisioned
		resource.Status.Conditions.SetCondition(core.Component{
			Name:    ComponentVPCCreator,
			Message: "VPC has been provisioned",
			Status:  core.SuccessStatus,
		})
		resource.Status.Status = core.SuccessStatus

		return false, nil
	}()
	if err != nil {
		resource.Status.Status = core.FailureStatus
	}

	if err := t.mgr.GetClient().Status().Update(ctx, resource); err != nil {
		logger.WithError(err).Error("updating the status of eks vpc")

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
func (t *eksvpcCtrl) GetCredentials(ctx context.Context, vpc *eksv1alpha1.EKSVPC, team string) (*eksv1alpha1.EKSCredentials, error) {
	// @step: is the team permitted access to this credentials
	permitted, err := t.Teams().Team(team).Allocations().IsPermitted(ctx, vpc.Spec.Credentials)
	if err != nil {
		log.WithError(err).Error("attempting to check for permission on credentials")

		return nil, fmt.Errorf("attempting to check for permission on credentials")
	}

	if !permitted {
		log.Warn("trying to build a vpc using unallocated credentials so not permitted")

		return nil, errors.New("you do not have permissions to the eks credentials")
	}

	// @step: retrieve the credentials
	creds := &eksv1alpha1.EKSCredentials{}

	return creds, t.mgr.GetClient().Get(ctx,
		types.NamespacedName{
			Namespace: vpc.Spec.Credentials.Namespace,
			Name:      vpc.Spec.Credentials.Name,
		}, creds,
	)
}
