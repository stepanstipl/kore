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
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/equality"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils/validation"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

// ServiceKinds is the interface to manage service kinds
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ServiceKinds
type ServiceKinds interface {
	// Delete is used to delete a service kind in the kore
	Delete(context.Context, string, ...DeleteOptionFunc) (*servicesv1.ServiceKind, error)
	// Get returns the service kind
	Get(context.Context, string) (*servicesv1.ServiceKind, error)
	// List returns the existing service kinds
	// The optional filter functions can be used to include items only for which all functions return true
	List(context.Context, ...func(servicesv1.ServiceKind) bool) (*servicesv1.ServiceKindList, error)
	// Has checks if a service kind exists
	Has(context.Context, string) (bool, error)
	// Update is responsible for updating a service kind
	Update(context.Context, *servicesv1.ServiceKind) error
}

type serviceKindsImpl struct {
	Interface
}

// Update is responsible for updating a service kind
func (p serviceKindsImpl) Update(ctx context.Context, kind *servicesv1.ServiceKind) error {
	if err := IsValidResourceName("service kind", kind.Name); err != nil {
		return err
	}

	if kind.Namespace != HubNamespace {
		return validation.NewError("%q failed validation", kind.Name).
			WithFieldErrorf("namespace", validation.InvalidValue, "must be %q", HubNamespace)
	}

	existing, err := p.Get(ctx, kind.Name)
	if err != nil && err != ErrNotFound {
		return err
	}

	if existing == nil {
		return ErrNotAllowed{message: "creating a new service kind is not allowed"}
	}

	existing.Spec.Enabled = kind.Spec.Enabled

	if !equality.Semantic.DeepEqual(existing.Spec, kind.Spec) {
		return validation.NewError("%q failed validation", kind.Name).
			WithFieldErrorf(validation.FieldRoot, validation.NotAllowed, "only the enabled field can be modified")
	}

	err = p.Store().Client().Update(ctx,
		store.UpdateOptions.To(kind),
		store.UpdateOptions.WithCreate(true),
		store.UpdateOptions.WithForce(true),
	)
	if err != nil {
		log.WithError(err).Error("failed to update a service kind")

		return err
	}

	return nil
}

// Delete is used to delete a service kind in the kore
func (p serviceKindsImpl) Delete(ctx context.Context, name string, o ...DeleteOptionFunc) (*servicesv1.ServiceKind, error) {
	opts := ResolveDeleteOptions(o)

	kind := &servicesv1.ServiceKind{}
	err := p.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(HubNamespace),
		store.GetOptions.InTo(kind),
		store.GetOptions.WithName(name),
	)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		log.WithError(err).Error("failed to retrieve the service kind")

		return nil, err
	}

	servicesWithKind, err := p.getServicesWithKind(ctx, name)
	if err != nil {
		return nil, err
	}
	if len(servicesWithKind) > 0 {
		if len(servicesWithKind) <= 5 {
			return nil, fmt.Errorf(
				"the service kind can not be deleted as there are %d services using it: %s",
				len(servicesWithKind),
				strings.Join(servicesWithKind, ", "),
			)
		}
		return nil, fmt.Errorf(
			"the service kind can not be deleted as there are %d services using it",
			len(servicesWithKind),
		)
	}

	if err := p.Store().Client().Delete(ctx, append(opts.StoreOptions(), store.DeleteOptions.From(kind))...); err != nil {
		log.WithError(err).Error("failed to delete the service kind")

		return nil, err
	}

	return kind, nil
}

// Get returns the service kind
func (p serviceKindsImpl) Get(ctx context.Context, name string) (*servicesv1.ServiceKind, error) {
	kind := &servicesv1.ServiceKind{}

	if found, err := p.Has(ctx, name); err != nil {
		return nil, err
	} else if !found {
		return nil, ErrNotFound
	}

	return kind, p.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(HubNamespace),
		store.GetOptions.WithName(name),
		store.GetOptions.InTo(kind),
	)
}

// List returns the existing service kinds
func (p serviceKindsImpl) List(ctx context.Context, filters ...func(servicesv1.ServiceKind) bool) (*servicesv1.ServiceKindList, error) {
	list := &servicesv1.ServiceKindList{}

	err := p.Store().Client().List(ctx,
		store.ListOptions.InNamespace(HubNamespace),
		store.ListOptions.InTo(list),
	)
	if err != nil {
		return nil, err
	}

	if len(filters) == 0 {
		return list, nil
	}

	res := []servicesv1.ServiceKind{}
	for _, item := range list.Items {
		if func() bool {
			for _, filter := range filters {
				if !filter(item) {
					return false
				}
			}
			return true
		}() {
			res = append(res, item)
		}
	}
	list.Items = res

	return list, nil
}

// Has checks if a service kind exists
func (p serviceKindsImpl) Has(ctx context.Context, name string) (bool, error) {
	return p.Store().Client().Has(ctx,
		store.HasOptions.InNamespace(HubNamespace),
		store.HasOptions.From(&servicesv1.ServiceKind{}),
		store.HasOptions.WithName(name),
	)
}

func (p serviceKindsImpl) getServicesWithKind(ctx context.Context, kind string) ([]string, error) {
	var res []string

	teamList, err := p.Teams().List(ctx)
	if err != nil {
		return nil, err
	}

	for _, team := range teamList.Items {
		servicesList, err := p.Teams().Team(team.Name).Services().List(ctx)
		if err != nil {
			return nil, err
		}
		for _, service := range servicesList.Items {
			if service.Spec.Kind == kind {
				res = append(res, fmt.Sprintf("%s/%s", team.Name, service.Name))
			}
		}
	}

	return res, nil
}
