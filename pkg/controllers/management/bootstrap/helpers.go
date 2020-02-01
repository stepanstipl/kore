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
	"text/template"
	"time"

	"github.com/appvia/kore/pkg/utils/kubernetes"

	"github.com/Masterminds/sprig"
	yaml "github.com/ghodss/yaml"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// JobConfigExists checks if the configuration exists already
func JobConfigExists(ctx context.Context, cc client.Client) (bool, error) {
	return HasConfigMap(ctx, cc, jobName)
}

// JobOLMConfigExists check if the olm configuration exists
func JobOLMConfigExists(ctx context.Context, cc client.Client) (bool, error) {
	return HasConfigMap(ctx, cc, jobOLMConfig)
}

// HasConfigMap checks if a configmap exists
func HasConfigMap(ctx context.Context, cc client.Client, name string) (bool, error) {
	return kubernetes.CheckIfExists(ctx, cc, &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: jobNamespace,
		},
	})
}

// JobExists checks if the bootstrap job there
func JobExists(ctx context.Context, cc client.Client) (bool, error) {
	return kubernetes.CheckIfExists(ctx, cc, &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: jobNamespace,
		},
	})
}

// GetBatchJob returns a job from the api
func GetBatchJob(ctx context.Context, cc client.Client, name string) (*batchv1.Job, error) {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: jobNamespace,
		},
	}

	found, err := kubernetes.GetIfExists(ctx, cc, job)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("resource not found")
	}

	return job, nil
}

// WaitOnJob checks the status of the job and if not successful returns the error
func WaitOnJob(ctx context.Context, cc client.Client) error {
	for {
		select {
		case <-ctx.Done():
			return errors.New("context has been cancelled")
		default:
		}

		err := func() error {
			job, err := GetBatchJob(ctx, cc, jobName)
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
