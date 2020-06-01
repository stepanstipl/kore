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

package configuration

import (
	"context"
	"encoding/json"
	"fmt"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type cachedSecrets map[types.NamespacedName]*configv1.Secret

func (c cachedSecrets) getIfExists(ctx context.Context, client client.Client, nsName types.NamespacedName) (*configv1.Secret, error) {
	if secret, cached := c[nsName]; cached {
		return secret, nil
	}

	secret := &configv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nsName.Name,
			Namespace: nsName.Namespace,
		},
	}

	exists, err := kubernetes.GetIfExists(ctx, client, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to load secret %s/%s: %w", nsName.Namespace, nsName.Name, err)
	}

	if !exists {
		c[nsName] = nil
		return nil, nil
	}

	if err := secret.Decode(); err != nil {
		return nil, fmt.Errorf("failed to decode secret %s/%s: %w", nsName.Namespace, nsName.Name, err)
	}

	c[nsName] = secret

	return secret, nil
}

type ObjectWithConfiguration interface {
	metav1.Object
	GetConfiguration() *apiextv1.JSON
}

type ObjectWithConfigurationFrom interface {
	metav1.Object
	GetConfigurationFrom() []corev1.ConfigurationFromSource
}

func ParseObjectConfiguration(ctx context.Context, client client.Client, o ObjectWithConfiguration, v interface{}) error {
	var configFromSource []corev1.ConfigurationFromSource
	if o, ok := o.(ObjectWithConfigurationFrom); ok {
		configFromSource = o.GetConfigurationFrom()
	}
	return ParseConfiguration(ctx, client, o.GetNamespace(), o.GetConfiguration(), configFromSource, v)
}

func ParseConfiguration(
	ctx context.Context,
	client client.Client,
	namespace string,
	config *apiextv1.JSON,
	configFromSource []corev1.ConfigurationFromSource,
	v interface{},
) error {
	if config == nil || len(config.Raw) == 0 {
		config = &apiextv1.JSON{Raw: []byte(`{}`)}
	}

	configData := map[string]interface{}{}

	if err := json.Unmarshal(config.Raw, &configData); err != nil {
		return fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	cachedSecrets := cachedSecrets{}

	for _, cfs := range configFromSource {
		switch {
		case cfs.SecretKeyRef != nil:
			secretNsName := types.NamespacedName{
				Namespace: cfs.SecretKeyRef.Namespace,
				Name:      cfs.SecretKeyRef.Name,
			}
			if secretNsName.Namespace == "" {
				secretNsName.Namespace = namespace
			}

			secret, err := cachedSecrets.getIfExists(ctx, client, secretNsName)
			if err != nil {
				return err
			}

			if secret == nil {
				if !cfs.SecretKeyRef.Optional {
					return fmt.Errorf("failed to load secret %s/%s: does not exist", secretNsName.Namespace, secretNsName.Name)
				} else {
					continue
				}
			}

			value, ok := secret.Spec.Data[cfs.SecretKeyRef.Key]
			if !ok {
				if !cfs.SecretKeyRef.Optional {
					return fmt.Errorf("key %q does not exist in secret %s/%s", cfs.SecretKeyRef.Key, secret.Namespace, secret.Name)
				} else {
					continue
				}
			}

			configData[cfs.Name] = string(value)
		default:
			return fmt.Errorf("configuration source definition is invalid, reference is missing for %s", cfs.Name)
		}
	}

	finalConfigBytes, err := json.Marshal(configData)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	if err := json.Unmarshal(finalConfigBytes, v); err != nil {
		return fmt.Errorf("failed to unmarshal the configuration: %w", err)
	}
	return nil
}
