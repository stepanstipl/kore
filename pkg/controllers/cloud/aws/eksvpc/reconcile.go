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
	"time"

	core "github.com/appvia/kore/pkg/apis/core/v1"
	eksv1alpha1 "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/utils/cloud/aws"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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

	if finalizer.NeedToAdd(resource) {
		logger.Info("adding our finalizer to the team resource")

		if err := finalizer.Add(resource); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{Requeue: true}, nil
	}

	// @step: we need to mark the cluster as pending
	if resource.Status.Conditions == nil {
		resource.Status.Conditions = core.Components{}
	}

	result, err := func() (reconcile.Result, error) {
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

			return reconcile.Result{}, err
		}
		logger.Info("Found EKSCredential")

		client, err := aws.NewVPCClient(*credentials, aws.VPC{
			CidrBlock: resource.Spec.PrivateIPV4Cidr,
			Name:      resource.Name,
			Region:    resource.Spec.Region,
			Tags: map[string]string{
				aws.TagKoreManaged: "true",
			},
		})
		if err != nil {
			resource.Status.Conditions.SetCondition(core.Component{
				Detail:  err.Error(),
				Name:    ComponentVPCCreator,
				Message: fmt.Sprintf("Failed to create vpc as input values invalid - %s", err),
				Status:  core.FailureStatus,
			})

			return reconcile.Result{}, err
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

			return reconcile.Result{}, err
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

				return reconcile.Result{Requeue: true}, nil
			}

			logger.Debug("creating or discovering a vpc in aws")
		}
		// Ensure this only reports if it exists when all resources exist or ensure update works
		ready, err := client.Ensure()
		if err != nil {
			logger.WithError(err).Error("failed to create or update the EKS VPC")

			resource.Status.Conditions.SetCondition(core.Component{
				Name:    ComponentVPCCreator,
				Message: "Failed trying to provision the EKS VPC",
				Detail:  err.Error(),
			})
		}

		if err != nil || !ready {
			return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
		}

		resource.Status.Infra.SecurityGroupIDs = []string{client.VPC.ControlPlaneSecurityGroupID}
		resource.Status.Infra.PublicSubnetIDs = client.VPC.PublicSubnetIDs
		resource.Status.Infra.PrivateSubnetIDs = client.VPC.PrivateSubnetIDs
		resource.Status.Infra.PublicIPV4EgressAddresses = client.VPC.PublicIPV4EgressAddresses

		// @step: update the state as provisioned
		resource.Status.Conditions.SetCondition(core.Component{
			Name:    ComponentVPCCreator,
			Message: "VPC has been provisioned",
			Status:  core.SuccessStatus,
		})
		resource.Status.Status = core.SuccessStatus

		return reconcile.Result{}, nil
	}()
	if err != nil {
		resource.Status.Status = core.FailureStatus
	}

	if err := t.mgr.GetClient().Status().Update(ctx, resource); err != nil {
		logger.WithError(err).Error("updating the status of eks vpc")

		return reconcile.Result{}, err
	}

	return result, err
}

// GetCredentials returns the cloud credential
func (t *eksvpcCtrl) GetCredentials(ctx context.Context, vpc *eksv1alpha1.EKSVPC, team string) (*aws.Credentials, error) {
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
	creds := &eksv1alpha1.EKSCredentials{
		ObjectMeta: metav1.ObjectMeta{
			Name:      vpc.Spec.Credentials.Name,
			Namespace: vpc.Spec.Credentials.Namespace,
		},
	}
	found, err := kubernetes.GetIfExists(ctx, t.mgr.GetClient(), creds)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("eks credentials: (%s/%s) not found", vpc.Spec.Credentials.Namespace, vpc.Spec.Credentials.Name)
	}

	// for backwards-compatibility, use the creds set on the EKSCredentials resource, if they exist
	if creds.Spec.SecretAccessKey != "" && creds.Spec.AccessKeyID != "" {
		return &aws.Credentials{
			AccountID:       creds.Spec.AccountID,
			AccessKeyID:     creds.Spec.AccessKeyID,
			SecretAccessKey: creds.Spec.SecretAccessKey,
		}, nil
	}

	// @step: we need to grab the secret
	secret, err := controllers.GetDecodedSecret(ctx, t.mgr.GetClient(), creds.Spec.CredentialsRef)
	if err != nil {
		return nil, err
	}

	return &aws.Credentials{
		AccountID:       creds.Spec.AccountID,
		AccessKeyID:     secret.Spec.Data["access_key_id"],
		SecretAccessKey: secret.Spec.Data["access_secret_key"],
	}, nil
}
