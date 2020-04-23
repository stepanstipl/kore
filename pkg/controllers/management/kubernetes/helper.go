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
	"bytes"
	"context"
	"errors"
	"fmt"
	"text/template"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	clusterappman "github.com/appvia/kore/pkg/clusterappman"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	log "github.com/sirupsen/logrus"

	"github.com/Masterminds/sprig"
	yaml "github.com/ghodss/yaml"
	apps "k8s.io/api/apps/v1"
	appsv1 "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ClusterUserRolesToMap iterates the clusters users and dedups them
func ClusterUserRolesToMap(users []clustersv1.ClusterUser) map[string][]string {
	roles := make(map[string][]string)

	for _, user := range users {
		for _, role := range user.Roles {
			if list, found := roles[role]; found {
				list = append(list, user.Username)
				roles[role] = list
			} else {
				roles[role] = []string{user.Username}
			}
		}
	}

	return roles
}

// ReconcileClusterRequests builds a list of requests based on a team change
func ReconcileClusterRequests(ctx context.Context, cc client.Client, team string) ([]reconcile.Request, error) {
	logger := log.WithFields(log.Fields{
		"team": team,
	})
	logger.Info("triggering a cluster reconcilation based on upstream trigger")

	list, err := ListAllClustersInTeam(ctx, cc, team)
	if err != nil {
		logger.WithError(err).Error("trying to list teams in clusters")

		// @TODO we need way to surface these to the users
		return []reconcile.Request{}, err
	}

	return ClustersToRequests(list.Items), nil
}

// ClustersToRequests converts a collection of claims to requests
func ClustersToRequests(items []clustersv1.Kubernetes) []reconcile.Request {
	requests := make([]reconcile.Request, len(items))

	// @step: trigger the namespaceclaims to reconcile
	for i := 0; i < len(items); i++ {
		requests[i] = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      items[i].Name,
				Namespace: items[i].Namespace,
			},
		}
	}

	return requests
}

// ListAllClustersInTeam does what it says
func ListAllClustersInTeam(ctx context.Context, cc client.Client, namespace string) (*clustersv1.KubernetesList, error) {
	list := &clustersv1.KubernetesList{}

	return list, cc.List(ctx, list, client.InNamespace(namespace))

}

// WaitOnStatus will wait until the status object exists
// TODO: define a status object suitabvle for overall cluster status
func WaitOnStatus(ctx context.Context, cc client.Client) error {
	// WaitOnStatus checks the status of the job and if not successful returns the error
	for {
		select {
		case <-ctx.Done():
			return errors.New("context has been cancelled")
		default:
		}

		err := func() error {
			exists, err := StatusExists(ctx, cc)
			if err != nil {
				return err
			}
			if exists {
				return nil
			}

			return errors.New("Kore cluster manager has not reported status yet")
		}()
		if err == nil {
			return nil
		}
		time.Sleep(10 * time.Second)
	}
}

// StatusExists checks if the status exists already
func StatusExists(ctx context.Context, cc client.Client) (bool, error) {
	return HasConfigMap(ctx, cc, clusterappman.StatusCongigMap)
}

// CreateConfig creates a configmap for configuring the kore cluster manager
func CreateConfig(ctx context.Context, cc client.Client, params Parameters) error {
	b, err := yaml.Marshal(params)
	if err != nil {
		return errors.New("can not marshall params into yaml")
	}
	// Specify the parameters in the config map
	cm := &core.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterappmanConfig,
			Namespace: clusterappmanNamespace,
		},
		Data: (map[string]string{"clusterconfig": string(b)}),
	}
	if _, err := kubernetes.CreateOrUpdate(ctx, cc, cm); err != nil {
		return err
	}
	return nil
}

// GetConfig will retrieve the cluster configuration values
func GetConfig(ctx context.Context, cc client.Client) (params *Parameters, err error) {
	p := &Parameters{}
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterappmanConfig,
			Namespace: clusterappmanNamespace,
		},
	}
	exists, err := kubernetes.GetIfExists(ctx, cc, cm)
	if err != nil {
		return p, fmt.Errorf(
			"error obtaining config for %s/%s - %s",
			clusterappmanNamespace,
			clusterappmanConfig,
			err,
		)
	}
	if exists {
		// get the data from the configmap data key
		b := []byte(cm.Data[clusterappmanConfigKey])
		if err := yaml.Unmarshal(b, p); err != nil {
			return p, fmt.Errorf(
				"unable to deserialize yaml data from config map: %s/%s, key: %s",
				clusterappmanNamespace,
				clusterappmanConfig,
				clusterappmanConfigKey,
			)
		}
		return p, nil
	}
	return p, fmt.Errorf(
		"missing configmap %s in namespace %s",
		clusterappmanConfig,
		clusterappmanNamespace,
	)
}

// ConfigExists check if the cluster configuration exists
func ConfigExists(ctx context.Context, cc client.Client) (bool, error) {
	return HasConfigMap(ctx, cc, clusterappmanConfig)
}

// HasConfigMap checks if a configmap exists
func HasConfigMap(ctx context.Context, cc client.Client, name string) (bool, error) {
	return kubernetes.CheckIfExists(ctx, cc, &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: clusterappmanNamespace,
		},
	})
}

// NamespaceExists checks if the bootstrap job there
func NamespaceExists(ctx context.Context, cc client.Client) (bool, error) {
	return kubernetes.CheckIfExists(ctx, cc, &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterappmanNamespace,
		},
	})
}

// DeploymentExists checks if the bootstrap job there
func DeploymentExists(ctx context.Context, cc client.Client, name, namespace string) (bool, error) {
	return kubernetes.CheckIfExists(ctx, cc, &apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	})
}

// EnsureNamespace creates a namespace for the clusterappmanager if required
func EnsureNamespace(ctx context.Context, cc client.Client, namespace string) error {
	return kubernetes.EnsureNamespace(ctx, cc, &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
			Labels: map[string]string{
				kore.Label("owned"): "true",
			},
		},
	})
}

// GetDeployment returns a deployment from the api
func GetDeployment(ctx context.Context, cc client.Client, name, namespace string) (*apps.Deployment, error) {
	deployment := &apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	found, err := kubernetes.GetIfExists(ctx, cc, deployment)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("resource not found")
	}

	return deployment, nil
}

// CreateOrUpdateClusterAppManDeployment will reconcile the clusterappman deployment
func CreateOrUpdateClusterAppManDeployment(ctx context.Context, cc client.Client, image string) error {
	name := clusterappmanDeployment
	if _, err := kubernetes.CreateOrUpdate(ctx, cc, &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterappmanDeployment,
			Namespace: clusterappmanNamespace,
			Labels: map[string]string{
				"name": name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": name,
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": name,
					},
				},
				Spec: v1.PodSpec{
					ServiceAccountName: "clusterappman",
					Containers: []v1.Container{
						{
							Name:  name,
							Image: image,
							Env: []v1.EnvVar{
								{
									Name:  "IN_CLUSTER",
									Value: "true",
								},
							},
							Command: []string{
								"/bin/kore-clusterappman",
							},
						},
					},
				},
			},
		},
	}); err != nil {
		return err
	}
	return nil
}

// CreateClusterManClusterRoleBinding creates (or updates) the cluster role binding required for the clusterappman
func CreateClusterManClusterRoleBinding(ctx context.Context, cc client.Client) error {
	if _, err := kubernetes.CreateOrUpdateManagedClusterRoleBinding(ctx, cc, &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kore:clusterappman",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "clusterappman",
				Namespace: clusterappman.KoreNamespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
	}); err != nil {
		return fmt.Errorf("error tying to apply kore clusterappman clusterrole %q", err)
	}
	return nil
}

// DecodeInTo decodes the template into the thing
func DecodeInTo(in string, out interface{}) error {
	return yaml.Unmarshal([]byte(in), out)
}

// MakeTemplate is responsible for generating the template
func MakeTemplate(content string, params interface{}) (string, error) {
	tpl, err := template.New("main").Funcs(sprig.TxtFuncMap()).Parse(content)
	if err != nil {
		return "", err
	}

	b := &bytes.Buffer{}
	if err := tpl.Execute(b, params); err != nil {
		return "", err
	}

	return b.String(), nil
}
