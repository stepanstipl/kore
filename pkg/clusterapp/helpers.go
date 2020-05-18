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

package clusterapp

import (
	"context"
	"errors"

	"github.com/appvia/kore/pkg/utils/kubernetes"

	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// HelmKeyForSecrets is the key to use in secret data that contains a yaml file data
	HelmKeyForSecrets = "values.yaml"
)

func setMissingNamespace(namespace string, obj runtime.Object) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		log.Debugf("no setting namespace here ->%v<- - %s", obj, err)

		return err
	}
	if accessor.GetNamespace() == "" {
		accessor.SetNamespace(namespace)
		log.Debugf(
			"updated namespace to %s on %s",
			accessor.GetNamespace(),
			accessor.GetName(),
		)
	}

	return nil
}

func ensureNamespace(ctx context.Context, cc client.Client, name string) error {
	return kubernetes.EnsureNamespace(ctx, cc, &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	})
}

func getObjMetaAndSetDefaultNamespace(obj runtime.Object, defaultNamepsace string) metav1.ObjectMeta {
	objMeta, _ := kubernetes.GetMeta(obj)
	if err := setMissingNamespace(defaultNamepsace, obj); err != nil {
		log.Debugf("error setting namespace for %v - %s", obj, err)
	}
	return objMeta
}

// createHelmSecrets creates a configmap for configuring the kore cluster manager
func getHelmSecret(chartApp ChartApp) (*HelmSecret, error) {
	hs := HelmSecret{
		Name:      chartApp.ReleaseName,
		Namespace: chartApp.DefaultNamespace,
		ValuesRef: map[string]interface{}{
			"secretKeyRef": map[string]interface{}{
				"name":       chartApp.ReleaseName,
				"namespace:": chartApp.DefaultNamespace,
				"key":        HelmKeyForSecrets,
				"optional":   false,
			},
		},
	}
	b, err := yaml.Marshal(chartApp.SecretValues)
	if err != nil {
		return nil, errors.New("can not marshall secret values into yaml")
	}
	// Specify the parameters in the secret
	hs.Secret = &core.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      chartApp.ReleaseName,
			Namespace: chartApp.DefaultNamespace,
		},
		Data: (map[string][]byte{HelmKeyForSecrets: b}),
	}
	return &hs, nil
}
