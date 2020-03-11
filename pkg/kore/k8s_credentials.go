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

package kore

import (
	"context"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/store"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KubernetesCredentials is the interface to the kubernetes credentials
type KubernetesCredentials interface {
	// Delete is used to delete a kubernetes credentials in the kore
	Delete(context.Context, string) (*clustersv1.KubernetesCredentials, error)
	// Get returns the class from the kore
	Get(context.Context, string) (*clustersv1.KubernetesCredentials, error)
	// List returns a list of classes
	List(context.Context) (*clustersv1.KubernetesCredentialsList, error)
	// Has checks if a resource exists
	Has(context.Context, string) (bool, error)
	// Update is responsible for update a kubernetes credentials in the kore
	Update(context.Context, *clustersv1.KubernetesCredentials) (*clustersv1.KubernetesCredentials, error)
}

type kcImpl struct {
	*hubImpl
	// team is the team
	team string
}

// Delete is used to delete a kubernetes credentials in the kore
func (n *kcImpl) Delete(ctx context.Context, name string) (*clustersv1.KubernetesCredentials, error) {
	original, err := n.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	if err := n.Store().Client().Delete(ctx, store.DeleteOptions.From(original)); err != nil {
		log.WithError(err).Error("trying to delete the kubernetes credentials")

		return nil, err
	}

	return original, nil
}

// Get returns the class from the kore
func (n *kcImpl) Get(ctx context.Context, name string) (*clustersv1.KubernetesCredentials, error) {
	cred := &v1.Secret{}

	if err := n.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(n.team),
		store.GetOptions.InTo(cred),
		store.GetOptions.WithName(name),
	); err != nil {
		log.WithError(err).Error("trying to retrieve credentials from the api")

		return nil, err
	}

	return fromSecretToKubernetesCredential(cred), nil
}

// List returns a list of classes
func (n *kcImpl) List(ctx context.Context) (*clustersv1.KubernetesCredentialsList, error) {
	list := &v1.SecretList{}
	if err := n.Store().Client().List(ctx,
		store.ListOptions.InNamespace(n.team),
		store.ListOptions.InTo(list),
		store.ListOptions.WithLabel(Label("type"), "kubernetescredentials"),
	); err != nil {
		log.WithError(err).Error("trying to retrieve the list of credentials")

		return nil, err
	}

	items := &clustersv1.KubernetesCredentialsList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: clustersv1.GroupVersion.String(),
			Kind:       "KubernetesCredentialsList",
		},
		Items: make([]clustersv1.KubernetesCredentials, len(list.Items)),
	}

	for i := 0; i < len(list.Items); i++ {
		items.Items[i] = *fromSecretToKubernetesCredential(&list.Items[i])
	}

	return items, nil
}

// Has checks if a resource exists
func (n *kcImpl) Has(ctx context.Context, name string) (bool, error) {
	if _, err := n.Get(ctx, name); err != nil {
		if kerrors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// Update is responsible for update a kubernetes credentials in the kore
func (n *kcImpl) Update(ctx context.Context, credentials *clustersv1.KubernetesCredentials) (*clustersv1.KubernetesCredentials, error) {
	return credentials, n.Store().Client().Update(ctx,
		store.UpdateOptions.To(&v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      credentials.Name,
				Namespace: n.team,
				Labels: map[string]string{
					Label("type"): "kubernetescredentials",
				},
			},
			Data: map[string][]byte{
				"ca.crt":   []byte(credentials.Spec.CaCertificate),
				"endpoint": []byte(credentials.Spec.Endpoint),
				"token":    []byte(credentials.Spec.Token),
			},
		}),
		store.UpdateOptions.WithCreate(true),
	)
}

func fromSecretToKubernetesCredential(secret *v1.Secret) *clustersv1.KubernetesCredentials {
	return &clustersv1.KubernetesCredentials{
		TypeMeta: metav1.TypeMeta{
			APIVersion: clustersv1.GroupVersion.String(),
			Kind:       "KubernetesCredentials",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secret.Name,
			Namespace: secret.Namespace,
		},
		Spec: clustersv1.KubernetesCredentialsSpec{
			CaCertificate: string(secret.Data["ca.crt"]),
			Endpoint:      string(secret.Data["endpoint"]),
			Token:         string(secret.Data["token"]),
		},
	}
}
