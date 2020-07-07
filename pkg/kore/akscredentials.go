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

	aks "github.com/appvia/kore/pkg/apis/aks/v1alpha1"
	"github.com/appvia/kore/pkg/store"

	log "github.com/sirupsen/logrus"
)

// AKSCredentials is the gke interface
type AKSCredentials interface {
	// Delete is responsible for deleting the AKS credentials
	Delete(context.Context, string) error
	// Get returns the AKS credentials
	Get(context.Context, string) (*aks.AKSCredentials, error)
	// List returns all the AKS credentials in the team
	List(context.Context) (*aks.AKSCredentialsList, error)
	// Update is used to update the EKS credentials
	Update(context.Context, *aks.AKSCredentials) (*aks.AKSCredentials, error)
}

type aksCredsImpl struct {
	*cloudImpl
	// team is the request team
	team string
}

// Delete is responsible for deleting the AKS credentials
func (h *aksCredsImpl) Delete(ctx context.Context, name string) error {
	logger := log.WithFields(log.Fields{
		"name": name,
		"team": h.team,
	})

	aksCreds := &aks.AKSCredentials{}
	if err := h.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(h.team),
		store.GetOptions.InTo(aksCreds),
		store.GetOptions.WithName(name),
	); err != nil {
		logger.WithError(err).Error("failed to retrieve the credentials")

		return err
	}

	return h.Store().Client().Delete(ctx, store.DeleteOptions.From(aksCreds))
}

// Get returns the AKS credentials
func (h *aksCredsImpl) Get(ctx context.Context, name string) (*aks.AKSCredentials, error) {
	aksCreds := &aks.AKSCredentials{}

	return aksCreds, h.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(h.team),
		store.GetOptions.InTo(aksCreds),
		store.GetOptions.WithName(name),
	)
}

// List returns all the AKS credentials in the team
func (h *aksCredsImpl) List(ctx context.Context) (*aks.AKSCredentialsList, error) {
	list := &aks.AKSCredentialsList{}

	return list, h.Store().Client().List(ctx,
		store.ListOptions.InNamespace(h.team),
		store.ListOptions.InTo(list),
	)
}

// Update is used to update the EKS credentials
func (h *aksCredsImpl) Update(ctx context.Context, aksCreds *aks.AKSCredentials) (*aks.AKSCredentials, error) {
	logger := log.WithFields(log.Fields{
		"name": aksCreds.Name,
		"team": h.team,
	})

	// @step: update the resource in the api
	if err := h.Store().Client().Update(ctx,
		store.UpdateOptions.To(aksCreds),
		store.UpdateOptions.WithCreate(true),
		store.UpdateOptions.WithForce(true),
		store.UpdateOptions.WithPatch(true),
	); err != nil {
		logger.WithError(err).Error("failed to update the aks credentials")

		return nil, err
	}

	return aksCreds, nil
}
