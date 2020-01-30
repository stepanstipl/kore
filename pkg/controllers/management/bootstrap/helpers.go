/*
 * Copyright (C) 2019 Appvia Ltd <info@appvia.io>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package bootstrap

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"text/template"
	"time"

	"github.com/appvia/kore/pkg/clusterman"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	"github.com/Masterminds/sprig"
	yaml "github.com/ghodss/yaml"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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
	return HasConfigMap(ctx, cc, clusterman.StatusCongigMap)
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
			Name:      clustermanConfig,
			Namespace: clustermanNamespace,
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
			Name:      clustermanConfig,
			Namespace: clustermanNamespace,
		},
	}
	exists, err := kubernetes.GetIfExists(ctx, cc, cm)
	if err != nil {
		return p, fmt.Errorf(
			"error obtaining config for %s/%s - %s",
			clustermanNamespace,
			clustermanConfig,
			err,
		)
	}
	if exists {
		// get the data from the configmap data key
		b := []byte(cm.Data[clustermanConfigKey])
		if err := yaml.Unmarshal(b, p); err != nil {
			return p, fmt.Errorf(
				"unable to deserialize yaml data from config map: %s/%s, key: %s",
				clustermanNamespace,
				clustermanConfig,
				clustermanConfigKey,
			)
		}
		return p, nil
	}
	return p, fmt.Errorf(
		"missing configmap %s in namespace %s",
		clustermanConfig,
		clustermanNamespace,
	)
}

// ConfigExists check if the cluster configuration exists
func ConfigExists(ctx context.Context, cc client.Client) (bool, error) {
	return HasConfigMap(ctx, cc, clustermanConfig)
}

// HasConfigMap checks if a configmap exists
func HasConfigMap(ctx context.Context, cc client.Client, name string) (bool, error) {
	return kubernetes.CheckIfExists(ctx, cc, &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: clustermanNamespace,
		},
	})
}

// NamespaceExists checks if the bootstrap job there
func NamespaceExists(ctx context.Context, cc client.Client) (bool, error) {
	return kubernetes.CheckIfExists(ctx, cc, &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: clustermanNamespace,
		},
	})
}

// DeploymentExists checks if the bootstrap job there
func DeploymentExists(ctx context.Context, cc client.Client) (bool, error) {
	return kubernetes.CheckIfExists(ctx, cc, &apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clustermanDeployment,
			Namespace: clustermanNamespace,
		},
	})
}

// EnsureNamespace creates a namespace for the clustermanager if required
func EnsureNamespace(ctx context.Context, cc client.Client) error {
	return kubernetes.EnsureNamespace(ctx, cc, &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: clustermanNamespace,
		},
	})
}

// GetDeployment returns a deployment from the api
func GetDeployment(ctx context.Context, cc client.Client, name string) (*apps.Deployment, error) {
	deployment := &apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: clustermanNamespace,
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

// CreateClusterRoleBinding creates (or updates) the cluster role binding required for the clustermanager
func CreateClusterRoleBinding(ctx context.Context, cc client.Client) error {
	if _, err := kubernetes.CreateOrUpdateManagedClusterRoleBinding(ctx, cc, &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kore-clusterman",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "default",
				Namespace: clustermanNamespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
	}); err != nil {
		return fmt.Errorf("error tying to apply kore clusterman clusterrole %q", err)
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
