/*
 * Copyright (C) 2019 Rohith Jayawardene <gambol99@gmail.com>
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
	"text/template"
	"time"

	clusterv1 "github.com/appvia/kore/pkg/apis/clusters/v1"

	"github.com/Masterminds/sprig"
	kutils "github.com/gambol99/hub-utils/pkg/kubernetes"
	yaml "github.com/ghodss/yaml"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ClusterRoleExists checks if a cluster role exists
func ClusterRoleExists(ctx context.Context, client kubernetes.Interface, name string) (bool, error) {
	_, err := client.RbacV1().ClusterRoles().Get(name, metav1.GetOptions{})
	if err != nil {
		if !kerrors.IsNotFound(err) {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

// JobConfigExists checks if the configuration exists already
func JobConfigExists(ctx context.Context, client kubernetes.Interface) (bool, error) {
	return HasConfigMap(ctx, client, jobName)
}

// JobOLMConfigExists check if the olm configuration exists
func JobOLMConfigExists(ctx context.Context, client kubernetes.Interface) (bool, error) {
	return HasConfigMap(ctx, client, jobOLMConfig)
}

// HasConfigMap checks if a configmap exists
func HasConfigMap(ctx context.Context, client kubernetes.Interface, name string) (bool, error) {
	if _, err := client.CoreV1().ConfigMaps(jobNamespace).Get(name, metav1.GetOptions{}); err != nil {
		if !kerrors.IsNotFound(err) {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

// JobExists checks if the bootstrap job there
func JobExists(ctx context.Context, client kubernetes.Interface) (bool, error) {
	_, err := client.BatchV1().Jobs(jobNamespace).Get(jobName, metav1.GetOptions{})
	if err != nil {
		if !kerrors.IsNotFound(err) {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

// WaitOnJob checks the status of the job and if not successful returns the error
func WaitOnJob(ctx context.Context, client kubernetes.Interface) error {
	for {
		select {
		case <-ctx.Done():
			return errors.New("context has been cancelled")
		default:
		}

		err := func() error {
			job, err := client.BatchV1().Jobs(jobNamespace).Get(jobName, metav1.GetOptions{})
			if err != nil {
				return err
			}
			if job.Status.Succeeded > 0 {
				return nil
			}

			return errors.New("job no completed yet")
		}()
		if err == nil {
			return nil
		}
		time.Sleep(10 * time.Second)
	}
}

// makeKubernetesClient returns a client from a cluster
func makeKubernetesClient(obj *clusterv1.Kubernetes) (kubernetes.Interface, error) {
	return kutils.NewFromToken(
		obj.Spec.Endpoint,
		obj.Spec.Token,
		obj.Spec.CaCertificate,
	)
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
