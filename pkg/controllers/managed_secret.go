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

package controllers

import (
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ManagedSecret is used a convient wrapper for a secret
type ManagedSecret struct {
	secret *configv1.Secret
}

// NewEmptySecret returns a managed secret wrapper
func NewEmptySecret() *ManagedSecret {
	return &ManagedSecret{
		secret: &configv1.Secret{
			TypeMeta: metav1.TypeMeta{
				APIVersion: configv1.SchemeGroupVersion.String(),
				Kind:       "Secret",
			},
			Spec: configv1.SecretSpec{},
		},
	}
}

// Name sets the name of the secret
func (m *ManagedSecret) Name(v string) *ManagedSecret {
	m.secret.SetName(v)

	return m
}

// Namespace set the location or namespace of the secret
func (m *ManagedSecret) Namespace(v string) *ManagedSecret {
	m.secret.SetNamespace(v)

	return m
}

// Description sets the description
func (m *ManagedSecret) Description(v string) *ManagedSecret {
	m.secret.Spec.Description = v

	return m
}

// Type sets the secret type
func (m *ManagedSecret) Type(v string) *ManagedSecret {
	m.secret.Spec.Type = v

	return m
}

// Values sets the values
func (m *ManagedSecret) Values(v map[string]string) *ManagedSecret {
	m.secret.Spec.Data = v

	return m
}

// Secret returns the actual secret
func (m *ManagedSecret) Secret() *configv1.Secret {
	return m.secret
}
