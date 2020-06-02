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
	"text/template"

	"github.com/appvia/kore/pkg/kore"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	log "github.com/sirupsen/logrus"

	"github.com/Masterminds/sprig"
	yaml "github.com/ghodss/yaml"
	apps "k8s.io/api/apps/v1"
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

	var res []clustersv1.Kubernetes
	for _, item := range list.Items {
		if item.Annotations[kore.AnnotationSystem] == kore.AnnotationValueTrue {
			continue
		}
		res = append(res, item)
	}

	return ClustersToRequests(res), nil
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

// DeploymentExists checks if the bootstrap job there
func DeploymentExists(ctx context.Context, cc client.Client, name, namespace string) (bool, error) {
	return kubernetes.CheckIfExists(ctx, cc, &apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
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
