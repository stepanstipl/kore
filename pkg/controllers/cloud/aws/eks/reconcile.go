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
	"time"

	awsv1alpha1 "github.com/appvia/kore/pkg/apis/aws/v1alpha1"
	core "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	"github.com/aws/aws-sdk-go/aws"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "eks.compute.kore.appvia.io"
	// ComponentClusterCreator is the name of the component for the UI
	ComponentClusterCreator = "Cluster Creator"
	// ComponentClusterBootstrap is the component name for seting up cloud credentials
	ComponentClusterBootstrap = "Cluster Initialize Access"
)

func (t *eksCtrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to reconcile aws eks cluster")

	cluster := &awsv1alpha1.EKSCluster{}
	if err := t.mgr.GetClient().Get(ctx, request.NamespacedName, cluster); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	finalizer := kubernetes.NewFinalizer(t.mgr.GetClient(), finalizerName)

	// @step: are we deleting the resource
	if finalizer.IsDeletionCandidate(cluster) {
		return t.Delete(request)
	}
	credentials := &awsv1alpha1.AWSCredential{}

	reference := types.NamespacedName{
		Namespace: cluster.Spec.Use.Namespace,
		Name:      cluster.Spec.Use.Name,
	}

	err := t.mgr.GetClient().Get(ctx, reference, credentials)
	if err != nil {
		return reconcile.Result{}, err
	}

	logger.Info("Found AWSCredential CR")

	sesh, err := GetAWSSession(credentials, cluster.Spec.Region)

	svc, err := GetEKSService(sesh)

	logger.Info("Checking cluster existence")

	clusterExists, err := CheckEKSClusterExists(svc, &awseks.DescribeClusterInput{
		Name: aws.String(cluster.Spec.Name),
	})

	if err != nil {
		return reconcile.Result{}, err
	}

	if clusterExists {
		logger.Info("Cluster exists: " + cluster.Spec.Name)
		return reconcile.Result{}, nil
	}

	logger.Info("Creating cluster:" + cluster.Spec.Name)

	// Cluster doesnt exist, create it
	_, err = CreateEKSCluster(svc, &awseks.CreateClusterInput{
		Name:    aws.String(cluster.Spec.Name),
		RoleArn: aws.String(cluster.Spec.RoleARN),
		Version: aws.String(cluster.Spec.Version),
		ResourcesVpcConfig: &awseks.VpcConfigRequest{
			SecurityGroupIds: aws.StringSlice(cluster.Spec.SecurityGroupIDs),
			SubnetIds:        aws.StringSlice(cluster.Spec.SubnetIDs),
		},
	})

	if err != nil {
		return reconcile.Result{}, err
	}

	// Set status to pending
	cluster.Status.Status = core.PendingStatus

	if err := t.mgr.GetClient().Status().Update(ctx, cluster); err != nil {
		logger.Error(err, "failed to update the resource status")
		return reconcile.Result{}, err
	}

	// Wait for it to become ACTIVE
	for {
		log.Println("Checking the status of cluster:", cluster.Spec.Name)

		status, err := GetEKSClusterStatus(svc, &awseks.DescribeClusterInput{
			Name: aws.String(cluster.Spec.Name),
		})

		if err != nil {
			return reconcile.Result{}, err
		}

		if status == "ACTIVE" {
			log.Println("Cluster active:", cluster.Spec.Name)
			// Set status to success
			cluster.Status.Status = core.SuccessStatus

			if err := t.mgr.GetClient().Status().Update(ctx, cluster); err != nil {
				logger.Error(err, "failed to update the resource status")
				return reconcile.Result{}, err
			}
			break
		}
		if status == "ERROR" {
			log.Println("Cluster has ERROR status:", cluster.Spec.Name)
			break
		}
		time.Sleep(5000 * time.Millisecond)
	}
	return reconcile.Result{}, nil
}
