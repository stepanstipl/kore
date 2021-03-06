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
	"encoding/json"
	"fmt"
	"strings"

	awsv1alpha1 "github.com/appvia/kore/pkg/apis/aws/v1alpha1"
	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eksv1alpha1 "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/clusterproviders"
	"github.com/appvia/kore/pkg/controllers"
	awsc "github.com/appvia/kore/pkg/controllers/cloud/aws"
	"github.com/appvia/kore/pkg/controllers/helpers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/serviceproviders/application"
	awsutils "github.com/appvia/kore/pkg/utils/cloud/aws"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SetComponents adds all provider-specific cluster components and updates dependencies if required
func (p Provider) SetComponents(ctx kore.Context, cluster *clustersv1.Cluster, components *kore.ClusterComponents) error {
	kubernetesObj := components.Find(func(comp kore.ClusterComponent) bool {
		_, ok := comp.Object.(*clustersv1.Kubernetes)
		return ok
	})

	meta := metav1.ObjectMeta{
		Name:      cluster.Name,
		Namespace: cluster.Namespace,
	}

	eksVPC := &eksv1alpha1.EKSVPC{ObjectMeta: meta}
	eks := &eksv1alpha1.EKS{ObjectMeta: meta}
	eksVPCDependencies := []kubernetes.Object{}
	if clusterproviders.IsAccountManaged(cluster.Spec.Credentials) {
		accountClaim := &awsv1alpha1.AWSAccountClaim{ObjectMeta: meta}
		components.Add(accountClaim)
		eksVPCDependencies = append(eksVPCDependencies, accountClaim)
	}

	components.Add(eksVPC, eksVPCDependencies...)
	components.AddComponent(&kore.ClusterComponent{
		Object: eks,
		Dependencies: []kubernetes.Object{
			eksVPC,
		},
		IsProvider: true,
	})

	config := string(cluster.Spec.Configuration.Raw)

	enableAutoscaler := false

	if groups := gjson.Get(config, "nodeGroups"); groups.Exists() && groups.IsArray() {
		var err error
		groups.ForEach(func(key, value gjson.Result) bool {
			if name := value.Get("name"); name.Exists() {
				groupName := cluster.Name + "-" + name.String()

				eksNodeGroup := &eksv1alpha1.EKSNodeGroup{
					ObjectMeta: metav1.ObjectMeta{
						Name:      groupName,
						Namespace: cluster.Namespace,
					},
				}

				if err = kubernetes.PatchSpec(eksNodeGroup, []byte(value.Raw)); err != nil {
					return false
				}

				if eksNodeGroup.Spec.EnableAutoscaler {
					enableAutoscaler = true
				}

				components.Add(eksNodeGroup, eks)

				kubernetesObj.Dependencies = append(kubernetesObj.Dependencies, eksNodeGroup)
			}

			return true
		})
		if err != nil {
			return err
		}
	}

	helmOperatorName := cluster.Name + "-" + kore.AppHelmOperator
	helmOperatorService := components.Find(func(comp kore.ClusterComponent) bool {
		if service, ok := comp.Object.(*servicesv1.Service); ok {
			return service.Name == helmOperatorName
		}
		return false
	})
	if helmOperatorService == nil {
		return fmt.Errorf("%q service can not be found", helmOperatorName)
	}

	autoscalerService := &servicesv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.Name + "-autoscaler",
			Namespace: cluster.Namespace,
		},
	}

	components.AddComponent(&kore.ClusterComponent{
		Object:       autoscalerService,
		Dependencies: []kubernetes.Object{helmOperatorService},
		Absent:       !enableAutoscaler,
		AfterDelete:  p.deleteAutoScalerRole,
	})

	return nil
}

// BeforeComponentsUpdate runs after the components are loaded but before updated
// The cluster components will be provided in dependency order
func (p Provider) BeforeComponentsUpdate(ctx kore.Context, cluster *clustersv1.Cluster, components *kore.ClusterComponents) error {
	vpcComponent := components.Find(func(comp kore.ClusterComponent) bool {
		_, ok := comp.Object.(*eksv1alpha1.EKSVPC)
		return ok
	})
	if vpcComponent == nil {
		panic("EKS VPC object not found in cluster components")
	}
	vpc := vpcComponent.Object.(*eksv1alpha1.EKSVPC)

	eksComponent := components.Find(func(comp kore.ClusterComponent) bool {
		_, ok := comp.Object.(*eksv1alpha1.EKS)
		return ok
	})
	if eksComponent == nil {
		panic("EKS object not found in cluster components")
	}
	eks := eksComponent.Object.(*eksv1alpha1.EKS)

	config := cluster.Spec.Configuration.Raw

	var autoScalingNodeGroups []map[string]interface{}

	for _, c := range *components {
		switch o := c.Object.(type) {
		case *eksv1alpha1.EKSVPC:
			if err := kubernetes.PatchSpec(o, config); err != nil {
				return err
			}
			o.Spec.Cluster = cluster.Ownership()
			o.Spec.Credentials = p.getNewCredentialFromCluster(cluster)
			clusterproviders.SetCommonTags(&o.Spec.Tags, cluster)

		case *eksv1alpha1.EKS:
			if err := kubernetes.PatchSpec(o, config); err != nil {
				return err
			}

			o.Spec.Cluster = cluster.Ownership()
			o.Spec.Credentials = p.getNewCredentialFromCluster(cluster)
			o.Spec.Region = vpc.Spec.Region
			o.Spec.SecurityGroupIDs = vpc.Status.Infra.SecurityGroupIDs
			o.Spec.SubnetIDs = vpc.Status.Infra.PrivateSubnetIDs
			o.Spec.SubnetIDs = append(o.Spec.SubnetIDs, vpc.Status.Infra.PublicSubnetIDs...)
			clusterproviders.SetCommonTags(&o.Spec.Tags, cluster)

		case *eksv1alpha1.EKSNodeGroup:
			groupName := strings.TrimPrefix(o.Name, cluster.Name+"-")

			if groups := gjson.Get(string(config), "nodeGroups"); groups.Exists() && groups.IsArray() {
				var err error

				groups.ForEach(func(key, value gjson.Result) bool {
					if name := value.Get("name"); name.Exists() && groupName == name.String() {
						if err = kubernetes.PatchSpec(o, []byte(value.Raw)); err != nil {
							return false
						}
					}

					return true
				})
				if err != nil {
					return err
				}
			}

			if o.Spec.EnableAutoscaler {
				for _, name := range o.Status.AutoScalingGroupNames {
					autoScalingNodeGroups = append(autoScalingNodeGroups, map[string]interface{}{
						"minSize": o.Spec.MinSize,
						"maxSize": o.Spec.MaxSize,
						"name":    name,
					})
				}
			}

			o.Spec.Cluster = cluster.Ownership()
			o.Spec.Credentials = p.getNewCredentialFromCluster(cluster)
			o.Spec.Region = vpc.Spec.Region
			o.Spec.Subnets = vpc.Status.Infra.PrivateSubnetIDs
			clusterproviders.SetCommonTags(&o.Spec.Tags, cluster)

		case *awsv1alpha1.AWSAccountClaim:
			// @step: we never touch the aws account object claim under these circumstances
			if c.Exists() {
				continue
			}

			// @step: we find the matching account rule
			mgmt, err := clusterproviders.FindAccountManagement(ctx, cluster.Spec.Credentials)
			if err != nil {
				return err
			}
			o.Spec.Organization = corev1.Ownership{
				Group:     awsv1alpha1.GroupVersion.Group,
				Kind:      "AWSOrganization",
				Name:      mgmt.Spec.Organization.Name,
				Namespace: mgmt.Spec.Organization.Namespace,
				Version:   awsv1alpha1.GroupVersion.Version,
			}

			rule, found := clusterproviders.FindAccountingRule(mgmt, cluster.Spec.Plan)
			if !found {
				return controllers.NewCriticalError(
					fmt.Errorf("no account rule matching plan: %q exist", cluster.Spec.Plan),
				)
			}

			// @step: we derive the account name from the rule
			name := cluster.Namespace
			if rule.Suffix != "" {
				name = fmt.Sprintf("%s-%s", name, rule.Suffix)
			}
			if rule.Prefix != "" {
				name = fmt.Sprintf("%s-%s", rule.Prefix, name)
			}

			o.Spec.AccountName = name

		case *servicesv1.Service:
			switch o.Name {
			case cluster.Name + "-autoscaler":
				if eks.Status.Status != corev1.SuccessStatus {

					// for now...
					return nil
				}
				creds, err := awsc.GetCredentials(ctx, vpc.Namespace, eks.Spec.Credentials)
				if err != nil {

					return err
				}
				if creds == nil {

					return fmt.Errorf("aws account credentials not available after eks is created")
				}

				iam := awsutils.NewIamClient(*creds)
				role, err := iam.EnsureClusterAutoscalerRole(cluster.Name, eks.Status.OIDCProviderURL)
				if err != nil {

					return err
				}

				o.Spec.Cluster = cluster.Ownership()
				o.Spec.ClusterNamespace = "kube-system"

				if o.Annotations == nil {
					o.Annotations = map[string]string{}
				}
				o.Annotations[kore.AnnotationSystem] = kore.AnnotationValueTrue

				values := map[string]interface{}{
					"cloud-provider": "aws",
					"image": map[string]interface{}{
						"tag": "v1.16.4",
					},
					"awsRegion": eks.Spec.Region,
					"rbac": map[string]interface{}{
						"create": true,
						"serviceAccount": map[string]interface{}{
							"name": "cluster-autoscaler",
						},
						"serviceAccountAnnotations": map[string]interface{}{
							"eks.amazonaws.com/role-arn": *role.Arn,
						},
					},
					"autoscalingGroups": autoScalingNodeGroups,
				}

				if err := helpers.ApplyServicePlanToAppService(ctx, o, application.HelmAppClusterAutoscaler, values); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// SetProviderData saves the provider data on the cluster
// The cluster components will be provided in dependency order
func (p Provider) SetProviderData(ctx kore.Context, cluster *clustersv1.Cluster, components *kore.ClusterComponents) error {
	providerData := map[string]interface{}{}
	if err := cluster.Status.GetProviderData(&providerData); err != nil {

		return err
	}
	creds, err := awsc.GetCredentials(ctx, cluster.Namespace, p.getNewCredentialFromCluster(cluster))
	if err != nil {

		return err
	}
	if creds == nil {

		// Account creds are not ready, we need to wait...
		return nil
	}

	providerData["awsAccountID"] = creds.AccountID

	for _, c := range *components {
		switch o := c.Object.(type) {
		case *eksv1alpha1.EKSVPC:
			vpcJSON, err := json.Marshal(o.Status.Infra)
			if err != nil {
				return err
			}
			vpcData := map[string]interface{}{}
			if err := json.Unmarshal(vpcJSON, &vpcData); err != nil {
				return err
			}

			providerData["vpc"] = vpcData
		case *eksv1alpha1.EKS:
			eksData := map[string]interface{}{
				"arn":             o.Status.ARN,
				"roleARN":         o.Status.RoleARN,
				"oidcProviderURL": o.Status.OIDCProviderURL,
			}
			providerData["eks"] = eksData
		}
	}

	return cluster.Status.SetProviderData(providerData)
}

func (p Provider) deleteAutoScalerRole(ctx kore.Context, cluster *clustersv1.Cluster, _ *kore.ClusterComponent, components *kore.ClusterComponents) error {
	eksComponent := components.Find(func(comp kore.ClusterComponent) bool {
		_, ok := comp.Object.(*eksv1alpha1.EKS)
		return ok
	})
	if eksComponent == nil {
		panic("EKS object not found in cluster components")
	}
	eks := eksComponent.Object.(*eksv1alpha1.EKS)
	creds, err := awsc.GetCredentials(ctx, cluster.Namespace, eks.Spec.Credentials)
	if err != nil {
		return err
	}
	if creds == nil {

		// Account credentials not available - nothing we can do
		log.Warnf("unable to delete cluster autoscaler roles as credential is no longer available")
		return nil
	}

	iam := awsutils.NewIamClient(*creds)
	if err := iam.DeleteClusterAutoscalerRole(cluster.Name); err != nil {
		return err
	}

	return nil
}

func (p Provider) getNewCredentialFromCluster(cluster *clustersv1.Cluster) corev1.Ownership {
	switch clusterproviders.IsAccountManaged(cluster.Spec.Credentials) {
	case true:
		return corev1.Ownership{
			Group:     awsv1alpha1.GroupVersion.Group,
			Version:   awsv1alpha1.GroupVersion.Version,
			Kind:      "AWSAccountClaim",
			Namespace: cluster.Namespace,
			Name:      cluster.Name,
		}
	default:
		return cluster.Spec.Credentials
	}

}
