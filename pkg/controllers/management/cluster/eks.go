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

package cluster

import (
	"context"
	"strings"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	"github.com/tidwall/gjson"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type eksComponents struct {
	*Controller
}

func (e *eksComponents) Components(cluster *clustersv1.Cluster, components *Components) controllers.EnsureFunc {
	return func(ctx context.Context) (reconcile.Result, error) {
		meta := metav1.ObjectMeta{
			Name:      cluster.Name,
			Namespace: cluster.Namespace,
		}

		v := components.Add(&eks.EKSVPC{ObjectMeta: meta})
		c := components.Add(&eks.EKS{ObjectMeta: meta})
		components.Edge(v, c)

		config := string(cluster.Spec.Configuration.Raw)

		if groups := gjson.Get(config, "nodeGroups"); groups.Exists() && groups.IsArray() {
			groups.ForEach(func(key, value gjson.Result) bool {

				if name := value.Get("name"); name.Exists() {
					groupName := cluster.Name + "-" + name.String()

					v := components.Add(&eks.EKSNodeGroup{
						ObjectMeta: metav1.ObjectMeta{
							Name:      groupName,
							Namespace: cluster.Namespace,
						},
					})
					components.Edge(c, v)
				}

				return true
			})
		}

		return reconcile.Result{}, nil
	}
}

// CompleteClusterComponents is used to fill in the resources if required
func (e *eksComponents) Complete(cluster *clustersv1.Cluster, components *Components) controllers.EnsureFunc {
	return func(ctx context.Context) (reconcile.Result, error) {
		var vpc *eks.EKSVPC

		config := cluster.Spec.Configuration.Raw

		return reconcile.Result{}, components.WalkFunc(func(v *Vertex) (bool, error) {
			switch {
			case utils.IsEqualType(v.Object, &eks.EKSVPC{}):
				vpc = v.Object.(*eks.EKSVPC)
				if err := kubernetes.PatchSpec(vpc, config); err != nil {
					return false, err
				}
				vpc.Spec.Cluster = cluster.Ownership()
				vpc.Spec.Credentials = cluster.Spec.Credentials

			case utils.IsEqualType(v.Object, &eks.EKS{}):
				ek := v.Object.(*eks.EKS)
				if err := kubernetes.PatchSpec(ek, config); err != nil {
					return false, err
				}

				ek.Spec.Cluster = cluster.Ownership()
				ek.Spec.Credentials = cluster.Spec.Credentials
				ek.Spec.Region = vpc.Spec.Region
				ek.Spec.SecurityGroupIDs = vpc.Status.Infra.SecurityGroupIDs
				ek.Spec.SubnetIDs = append(ek.Spec.SubnetIDs, vpc.Status.Infra.PublicSubnetIDs...)
				ek.Spec.SubnetIDs = vpc.Status.Infra.PrivateSubnetIDs

			case utils.IsEqualType(v.Object, &eks.EKSNodeGroup{}):
				eg := v.Object.(*eks.EKSNodeGroup)
				groupName := strings.TrimPrefix(eg.Name, cluster.Name+"-")

				if groups := gjson.Get(string(config), "nodeGroups"); groups.Exists() && groups.IsArray() {
					var err error

					groups.ForEach(func(key, value gjson.Result) bool {
						if name := value.Get("name"); name.Exists() && groupName == name.String() {
							if err = kubernetes.PatchSpec(eg, []byte(value.Raw)); err != nil {
								return false
							}
						}

						return true
					})
					if err != nil {
						return false, err
					}
				}

				eg.Spec.Cluster = cluster.Ownership()
				eg.Spec.Credentials = cluster.Spec.Credentials
				eg.Spec.Region = vpc.Spec.Region
				eg.Spec.Subnets = vpc.Status.Infra.PrivateSubnetIDs
			}

			return true, nil
		})
	}
}
