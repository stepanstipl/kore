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

package kubernetes

import (
	"fmt"
	"net/url"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/appvia/kore/pkg/controllers/helpers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/serviceproviders/application"
	awsutils "github.com/appvia/kore/pkg/utils/cloud/aws"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	// ComponentAWSAutoscaler is the component for managing the AWS EKS Autoscaler and dependencies
	ComponentAWSAutoscaler = "Kubernetes Autoscaler for EKS"
	// ComponentAWSAutoscalerDescription is the description for the AWS Autoscaler component
	ComponentAWSAutoscalerDescription = "Scales the number of nodes up and down as demand changes in nodegroups where this is enabled"
)

type awsAutoScaler struct {
	eks         eks.EKS
	cluster     *clustersv1.Cluster
	kubeCluster *clustersv1.Kubernetes
	OIDCUrl     *url.URL
	nodeGroups  *[]eks.EKSNodeGroup
	ctx         kore.Context
}

// newAWSAutoscalerIfEnabled will return an autoscaler componenrt if enabled
func newAwsAutoscaler(ctx kore.Context, eksCluster eks.EKS, kubeCluster *clustersv1.Kubernetes) (*awsAutoScaler, error) {
	return &awsAutoScaler{
		eks:         eksCluster,
		kubeCluster: kubeCluster,
		ctx:         ctx,
	}, nil
}

// IsRequired will determine if the autoscaler is required for any nodegroups
func (a *awsAutoScaler) IsRequired() (bool, error) {
	if a.nodeGroups == nil {
		err := a.updateAutoscalingNodeGroups()
		if err != nil {
			return false, err
		}
	}

	// Discover if any nodes have autoscaling enabled from the resources
	return (len(*a.nodeGroups) > 0), nil
}

// Delete will tidy up all AWS resources and workloads
func (a *awsAutoScaler) Delete() (reconcile.Result, error) {
	if a.nodeGroups == nil {
		err := a.updateAutoscalingNodeGroups()
		if err != nil {

			return reconcile.Result{}, err
		}
	}
	a.kubeCluster.Status.Components.SetCondition(corev1.Component{
		Name:    "Service/" + application.HelmAppClusterAutoscaler,
		Message: "Autoscaler and AWS roles are being deleted",
		Status:  corev1.DeletingStatus,
	})

	if len(*a.nodeGroups) < 1 {

		return reconcile.Result{}, fmt.Errorf("No autoscaling groups with autoscaling enabled sepecifed for cluster %s", a.cluster.Name)
	}
	// Get an the kore EKS AWS client
	koreEKS := helpers.NewKoreEKS(
		a.ctx,
		&a.eks,
		a.ctx.Client(),
		a.ctx.Kore(),
		a.ctx.Logger())

	awsClient, err := koreEKS.GetClusterClient()
	if err != nil {

		return reconcile.Result{}, err
	}
	// Get the cluster details
	awsEks, err := awsClient.Describe(a.ctx)
	if err != nil {

		return reconcile.Result{}, fmt.Errorf("error getting cluster details for %s - %s", a.cluster.Name, err)
	}

	// Now enable the OIDC provider for the cluster
	awsIAM := awsutils.NewIamClientFromSession(awsClient.Sess)
	accountID, err := awsIAM.GetAccountNameFromARN(*awsEks.Arn)
	if err != nil {

		return reconcile.Result{}, fmt.Errorf("error getting account ID from AWS EKS cluster - %s", err)
	}
	ngNames := []string{}
	for _, ng := range *a.nodeGroups {
		awsNg, err := awsClient.DescribeNodeGroup(a.ctx, &ng)
		if err != nil {

			return reconcile.Result{}, fmt.Errorf("unable to obtain nodegroup information for nodegroup %s - %s", ng.Name, err)
		}
		ngNames = append(ngNames, *awsNg.NodegroupName)
	}

	if err := awsIAM.DeleteClusterAutoscalerRoleAndPolicies(a.ctx, a.cluster.Name, ngNames, accountID); err != nil {

		return reconcile.Result{}, fmt.Errorf("unable to delete AWS IAM roles for autoscaling in cluster %s - %s", a.eks.ClusterName, err)
	}
	return reconcile.Result{}, nil
}

// Ensure will create all AWS dependencies and then deploy the autoscaler
func (a *awsAutoScaler) Ensure() (reconcile.Result, error) {
	if a.nodeGroups == nil {
		err := a.updateAutoscalingNodeGroups()
		if err != nil {
			return reconcile.Result{}, err
		}
	}
	a.kubeCluster.Status.Components.SetCondition(corev1.Component{
		Name:    "Service/" + application.HelmAppClusterAutoscaler,
		Message: "Autoscaler is being provisioned",
		Status:  corev1.PendingStatus,
	})

	if len(*a.nodeGroups) < 1 {

		return reconcile.Result{}, fmt.Errorf("No autoscaling groups with autoscaling enabled sepecifed for cluster %s", a.cluster.Name)
	}
	// Get an the kore EKS AWS client
	koreEKS := helpers.NewKoreEKS(
		a.ctx,
		&a.eks,
		a.ctx.Client(),
		a.ctx.Kore(),
		a.ctx.Logger())

	awsClient, err := koreEKS.GetClusterClient()
	if err != nil {

		return reconcile.Result{}, err
	}
	// Get the cluster details
	awsEks, err := awsClient.Describe(a.ctx)
	if err != nil {

		return reconcile.Result{}, fmt.Errorf("error getting cluster details for %s - %s", a.cluster.Name, err)
	}

	// Now enable the OIDC provider for the cluster
	awsIAM := awsutils.NewIamClientFromSession(awsClient.Sess)
	awsASG := awsutils.NewASGClient(awsClient.Sess)
	if err := awsIAM.EnsureIRSA(*awsEks.Arn, *awsEks.Identity.Oidc.Issuer); err != nil {

		return reconcile.Result{}, fmt.Errorf("error setting up identity for cluster - %s", err)
	}
	accountID, err := awsIAM.GetAccountNameFromARN(*awsEks.Arn)
	if err != nil {

		return reconcile.Result{}, fmt.Errorf("error getting account ID from AWS EKS cluster - %s", err)
	}
	// get the Autoscaling Group ARN's for the relevant nodegroups
	nags := []awsutils.NodeGroupAutoScaler{}
	asgValues := []map[string]interface{}{}
	for _, ng := range *a.nodeGroups {
		awsNg, err := awsClient.DescribeNodeGroup(a.ctx, &ng)
		if err != nil {

			return reconcile.Result{}, fmt.Errorf("unable to obtain nodegroup information for nodegroup %s - %s", ng.Name, err)
		}

		for _, asgName := range awsNg.Resources.AutoScalingGroups {
			asg, err := awsASG.GetASGFromName(*asgName.Name)
			if err != nil {
				return reconcile.Result{}, fmt.Errorf("unable to obtain aws autoscaling group information from atoscaling group name - %s", err)
			}
			// required for setting up kubernetes autoscaling deployment
			asgValues = append(asgValues, map[string]interface{}{
				"minSize": *asg.MinSize,
				"maxSize": *asg.MaxSize,
				"name":    *asg.AutoScalingGroupName,
			})

			// required for setting up IAM
			nags = append(nags, awsutils.NodeGroupAutoScaler{
				AutoScalingARN: *asg.AutoScalingGroupARN,
				NodeGroupName:  *awsNg.NodegroupName,
			})
		}
	}
	issuerURL, _ := url.Parse(*awsEks.Identity.Oidc.Issuer)
	oidcIssuer := issuerURL.Hostname() + issuerURL.Path

	asgRole, err := awsIAM.EnsureClusterAutoscalingRoleAndPolicies(
		a.ctx,
		a.eks.Name,
		accountID,
		oidcIssuer,
		nags,
	)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("unable to create all the roles required for autoscaling for %s - %s", *awsEks.Name, err)
	}
	// TODO: select a planwith the correct version
	awsAutoScalerService, err := helpers.GetServiceFromPlanNameAndValues(
		a.ctx,
		application.HelmAppClusterAutoscaler,
		a.kubeCluster,
		"kube-system",
		map[string]interface{}{
			"cloud-provider": "aws",
			"image": map[string]interface{}{
				"tag": "v1.16.4",
			},
			"awsRegion": *awsClient.Sess.Config.Region,
			"rbac": map[string]interface{}{
				"create": true,
				"serviceAccount": map[string]interface{}{
					"name": "cluster-autoscaler",
				},
				"serviceAccountAnnotations": map[string]interface{}{
					"eks.amazonaws.com/role-arn": *asgRole.Arn,
				},
			},
			"autoscalingGroups": asgValues,
		})
	if err != nil {
		return reconcile.Result{}, err
	}

	return helpers.EnsureService(a.ctx, awsAutoScalerService, a.kubeCluster, a.kubeCluster.Status.Components)
}

func (a *awsAutoScaler) updateCluster() error {
	key := types.NamespacedName{
		Namespace: a.eks.Spec.Cluster.Namespace,
		Name:      a.eks.Spec.Cluster.Name,
	}
	cluster := &clustersv1.Cluster{}
	if err := a.ctx.Client().Get(a.ctx, key, cluster); err != nil {
		return err
	}
	a.cluster = cluster
	return nil
}

func (a *awsAutoScaler) updateAutoscalingNodeGroups() error {
	if a.cluster == nil {
		if err := a.updateCluster(); err != nil {
			return err
		}
	}
	ngs := []eks.EKSNodeGroup{}
	for _, c := range a.cluster.Status.Components {
		// TODO reflect the Kind from the expected type
		if c.Resource.Kind == "EKSNodeGroup" {
			// Get the object from reference:
			ng := &eks.EKSNodeGroup{}
			if err := a.ctx.Client().Get(a.ctx, c.Resource.NamespacedName(), ng); err != nil {
				return err
			}
			if ng.Spec.EnableAutoscaler {
				ngs = append(ngs, *ng)
			}
		}
	}
	a.nodeGroups = &ngs
	return nil
}
