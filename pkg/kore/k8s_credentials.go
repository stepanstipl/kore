/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 * 
 * This file is part of kore.
 * 
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * 
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * 
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package kore 

import (
	"context"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/store"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

// KubernetesCredentials is the interface to the kubernetes credentials
type KubernetesCredentials interface {
	// Delete is used to delete a kubernetes credentials in the kore
	Delete(context.Context, string) (*clustersv1.KubernetesCredentials error)
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

	return ns, n.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(n.team),
		store.GetOptions.InTo(cred),
		store.GetOptions.WithName(name),
	)
}

// List returns a list of classes
func (n *kcImpl) List(ctx context.Context) (*clustersv1.KubernetesCredentialsList, error) {
	return nil, nil
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
func (n *kcImpl) Update(ctx context.Context, credentials *clustersv1.KubernetesCredentials) (*clustersv1.NamespaceClaim, error) {
	return namespace, n.Store().Client().Update(ctx,
		store.UpdateOptions.To(&v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: credentials.Name,
				Namespace: n.team,
			},
			Data: map[string][]byte{
				"ca.crt": []byte(credentials.Spec.CaCertificate),
				"endpoint": []byte(credentials.Spec.Endpoint),
				"token": []byte(credentials.Spec.Token),
			},
		}),
		store.UpdateOptions.WithCreate(true),
	)
}
