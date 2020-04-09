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
	"fmt"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eksv1alpha1 "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/appvia/kore/pkg/utils/cloud/aws"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is responsible for deleting the aws eks cluster
func (t *eksvpcCtrl) Delete(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
		"team":      request.NamespacedName.Name,
	})
	logger.Info("attempting to delete eks vpc")

	resource := &eksv1alpha1.EKSVPC{}
	if err := t.mgr.GetClient().Get(ctx, request.NamespacedName, resource); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}
	original := resource.DeepCopy()

	finalizer := kubernetes.NewFinalizer(t.mgr.GetClient(), finalizerName)

	requeue, err := func() (bool, error) {
		creds, err := t.GetCredentials(ctx, resource, request.NamespacedName.Name)
		if err != nil {
			return false, err
		}

		// @step: TODO: first we need to check if there are any EKS clusters present
		// For now we just check the one cluster we know about...
		// @step: create a cloud client
		client, err := aws.NewVPCClient(aws.Credentials{
			AccessKeyID:     creds.Spec.AccessKeyID,
			SecretAccessKey: creds.Spec.SecretAccessKey,
		}, aws.VPC{
			CidrBlock: resource.Spec.PrivateIPV4Cidr,
			Name:      resource.Name,
			Region:    resource.Spec.Region,
		})
		if err != nil {
			return false, fmt.Errorf("invalid input details for vpc - %s", err)
		}

		found, err := client.Exists()
		if err != nil {
			return false, fmt.Errorf("checking if vpc exists - %s", err)
		}

		if found {
			eksClient := aws.NewEKSClientFromVPC(client, resource.ClusterName)

			// @step: check if the referenced CLUSTER exists and if so we wait...
			eksfound, err := eksClient.Exists()
			if err != nil {
				return false, fmt.Errorf("error checking if cluster exists: %s", err)
			}
			if eksfound {
				// We still have a CLUSTER so we can't delete this VPC yet
				// - reque

				return true, nil
			}

			// @step: lets update the status of the resource to deleting
			if resource.Status.Status != corev1.DeletingStatus {
				resource.Status.Status = corev1.DeletingStatus

				return true, nil
			}
			// We can now delete the VPC
			err = client.Delete()
			if err != nil {
				return false, err
			}
		}
		// no vpc - no reque
		return false, nil
	}()
	if err != nil {
		logger.WithError(err).Error("attempting to delete the vpc")

		return reconcile.Result{}, err
	}
	if requeue {
		if err := t.mgr.GetClient().Status().Patch(ctx, resource, client.MergeFrom(original)); err != nil {
			logger.WithError(err).Error("trying to update the resource status")

			return reconcile.Result{}, err
		}

		return reconcile.Result{Requeue: true}, nil
	}

	if err := finalizer.Remove(resource); err != nil {
		logger.WithError(err).Error("removing the finalizer")

		return reconcile.Result{}, err
	}

	logger.Info("successfully deleted the vpc")

	return reconcile.Result{}, nil
}
