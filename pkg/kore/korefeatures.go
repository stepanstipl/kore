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

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	"github.com/appvia/kore/pkg/store"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

// KoreFeatures describes the API for controlling kore features
type KoreFeatures interface {
	// List returns a list of features
	List(ctx context.Context) (*configv1.KoreFeatureList, error)
	// Get returns a specific feature
	Get(ctx context.Context, name string) (*configv1.KoreFeature, error)
	// Update changes or creates a feature
	Update(ctx context.Context, feature *configv1.KoreFeature) (*configv1.KoreFeature, error)
	// Delete removes a feature
	Delete(ctx context.Context, name string, o ...DeleteOptionFunc) (*configv1.KoreFeature, error)
	// CheckDelete verifies whether the feature can be deleted
	CheckDelete(ctx context.Context, feature *configv1.KoreFeature, o ...DeleteOptionFunc) error
}

type koreFeaturesImpl struct {
	store store.Store
}

var _ KoreFeatures = &koreFeaturesImpl{}

func (f *koreFeaturesImpl) List(ctx context.Context) (*configv1.KoreFeatureList, error) {
	list := &configv1.KoreFeatureList{}

	err := f.store.Client().List(ctx,
		store.ListOptions.InNamespace(HubNamespace),
		store.ListOptions.InTo(list),
	)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (f *koreFeaturesImpl) Get(ctx context.Context, name string) (*configv1.KoreFeature, error) {
	if found, err := f.store.Client().Has(ctx,
		store.HasOptions.InNamespace(HubNamespace),
		store.HasOptions.From(&configv1.KoreFeature{}),
		store.HasOptions.WithName(name)); err != nil {
		return nil, err
	} else if !found {
		return nil, ErrNotFound
	}

	feature := &configv1.KoreFeature{}
	return feature, f.store.Client().Get(ctx,
		store.GetOptions.InNamespace(HubNamespace),
		store.GetOptions.WithName(name),
		store.GetOptions.InTo(feature),
	)
}

func (f *koreFeaturesImpl) Update(ctx context.Context, feature *configv1.KoreFeature) (*configv1.KoreFeature, error) {
	// Force namespace to the kore namespace as these are not team scoped resources.
	feature.Namespace = HubNamespace

	err := f.store.Client().Update(ctx,
		store.UpdateOptions.To(feature),
		store.UpdateOptions.WithCreate(true),
		store.UpdateOptions.WithForce(true),
	)

	return feature, err
}

func (f *koreFeaturesImpl) CheckDelete(ctx context.Context, feature *configv1.KoreFeature, o ...DeleteOptionFunc) error {
	// opts := ResolveDeleteOptions(o)

	// if !opts.Cascade {
	// 	// @TODO: Consider checking if this feature has any dependencies which will prevent it being deleted
	// }

	return nil
}

func (f *koreFeaturesImpl) Delete(ctx context.Context, name string, o ...DeleteOptionFunc) (*configv1.KoreFeature, error) {
	opts := ResolveDeleteOptions(o)

	feature := &configv1.KoreFeature{}
	err := f.store.Client().Get(ctx,
		store.GetOptions.InNamespace(HubNamespace),
		store.GetOptions.InTo(feature),
		store.GetOptions.WithName(name),
	)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		log.WithError(err).Error("failed to retrieve the feature")

		return nil, err
	}

	if err := opts.Check(feature, func(o ...DeleteOptionFunc) error { return f.CheckDelete(ctx, feature, o...) }); err != nil {
		return nil, err
	}

	if err := f.store.Client().Delete(ctx, append(opts.StoreOptions(), store.DeleteOptions.From(feature))...); err != nil {
		log.WithError(err).Error("failed to delete the feature")

		return nil, err
	}

	return feature, nil
}
