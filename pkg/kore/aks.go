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

	aksv1alpha1 "github.com/appvia/kore/pkg/apis/aks/v1alpha1"

	"github.com/appvia/kore/pkg/store"
	kerrors "k8s.io/apimachinery/pkg/api/errors"

	log "github.com/sirupsen/logrus"
)

// AKS is the interface for managing Azure AKS clusters
type AKS interface {
	// Delete is responsible for deleting an AKS cluster
	Delete(context.Context, string) error
	// Get return the EKS cluster
	Get(context.Context, string) (*aksv1alpha1.AKS, error)
	// List returns all the AKS clusters in the team
	List(context.Context) (*aksv1alpha1.AKSList, error)
	// Update is used to update the AKS cluster
	Update(context.Context, *aksv1alpha1.AKS) (*aksv1alpha1.AKS, error)
}

type aksImpl struct {
	*cloudImpl
	// team is the request team
	team string
}

// Delete is responsible for deleting an AKS cluster
func (a *aksImpl) Delete(ctx context.Context, name string) error {
	logger := log.WithFields(log.Fields{
		"name": name,
		"team": a.team,
	})

	aks := &aksv1alpha1.AKS{}
	err := a.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(a.team),
		store.GetOptions.InTo(aks),
		store.GetOptions.WithName(name),
	)
	if err != nil {
		if !kerrors.IsNotFound(err) {
			logger.WithError(err).Error("failed to retrieve the AKS cluster from api")
		}

		return err
	}

	return a.Store().Client().Delete(ctx, store.DeleteOptions.From(aks))
}

// Get return the EKS cluster
func (a *aksImpl) Get(ctx context.Context, name string) (*aksv1alpha1.AKS, error) {
	logger := log.WithFields(log.Fields{
		"name": name,
		"team": a.team,
	})

	aks := &aksv1alpha1.AKS{}

	err := a.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(a.team),
		store.GetOptions.InTo(aks),
		store.GetOptions.WithName(name),
	)
	if err != nil {
		logger.Error("failed to retrieve the AKS cluster")

		return nil, err
	}

	return aks, nil
}

// List returns all the AKS clusters in the team
func (a *aksImpl) List(ctx context.Context) (*aksv1alpha1.AKSList, error) {
	list := &aksv1alpha1.AKSList{}

	return list, a.Store().Client().List(ctx,
		store.ListOptions.InNamespace(a.team),
		store.ListOptions.InTo(list),
	)
}

// Update is used to update the AKS cluster
func (a *aksImpl) Update(ctx context.Context, aks *aksv1alpha1.AKS) (*aksv1alpha1.AKS, error) {
	logger := log.WithFields(log.Fields{
		"name": aks.Name,
		"team": a.team,
	})

	aks.Namespace = a.team

	permitted, err := a.Teams().Team(a.team).Allocations().IsPermitted(ctx, aks.Spec.Credentials)
	if err != nil {
		logger.WithError(err).Error("failed to check for credentials allocation for the AKS cluster")

		return nil, err
	}
	if !permitted {
		return nil, NewErrNotAllowed("the requested credentials have not been allocated to you")
	}

	_, err = a.Get(ctx, aks.Name)
	if err != nil {
		if !kerrors.IsNotFound(err) {
			logger.WithError(err).Error("failed to retrieve the AKS cluster")

			return nil, err
		}
	}

	if err := a.Store().Client().Update(ctx,
		store.UpdateOptions.To(aks),
		store.UpdateOptions.WithCreate(true),
		store.UpdateOptions.WithForce(true),
		store.UpdateOptions.WithPatch(true),
	); err != nil {
		logger.WithError(err).Error("failed to update the AKS cluster")

		return nil, err
	}

	return aks, nil
}
