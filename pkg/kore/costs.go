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

	costsv1 "github.com/appvia/kore/pkg/apis/costs/v1beta1"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils/validation"

	log "github.com/sirupsen/logrus"
)

// Costs represents the interface to the top-level Kore Costs service.
type Costs interface {
	List(ctx context.Context) (*costsv1.CostList, error)
	Get(ctx context.Context, name string) (*costsv1.Cost, error)
	Update(ctx context.Context, cost *costsv1.Cost, ignoreReadonly bool) error
	Delete(ctx context.Context, name string) (*costsv1.Cost, error)
}

var _ Costs = &costsImpl{}

type costsImpl struct {
	Interface
}

// Update is responsible for updating a cost
func (c costsImpl) Update(ctx context.Context, cost *costsv1.Cost, ignoreReadonly bool) error {
	cost.Namespace = HubAdminTeam

	// Check the user is admin or has kore permissions
	user := authentication.MustGetIdentity(ctx)
	if !user.IsGlobalAdmin() {
		log.WithField("user", user.Username()).Warn("failed to update the cost without permissions")

		return ErrUnauthorized
	}

	if !ignoreReadonly {
		original, err := c.Get(ctx, cost.Name)
		if err != nil && err != ErrNotFound {
			return err
		}

		if original != nil && original.Annotations[AnnotationReadOnly] == AnnotationValueTrue {
			return validation.NewError("the cost can not be updated").WithFieldError(validation.FieldRoot, validation.ReadOnly, "cost is read-only")
		}
		if cost.Annotations[AnnotationReadOnly] == AnnotationValueTrue {
			return validation.NewError("the cost can not be updated").WithFieldError(validation.FieldRoot, validation.ReadOnly, "read-only flag can not be set")
		}
	}

	err := c.Store().Client().Update(ctx,
		store.UpdateOptions.To(cost),
		store.UpdateOptions.WithCreate(true),
		store.UpdateOptions.WithForce(true),
	)
	if err != nil {
		log.WithError(err).Error("failed to update the cost in the kore")

		return err
	}

	return nil
}

// Delete is used to delete a cost
func (c costsImpl) Delete(ctx context.Context, name string) (*costsv1.Cost, error) {
	user := authentication.MustGetIdentity(ctx)
	if !user.IsGlobalAdmin() {
		log.WithField("user", user.Username()).Warn("failed to delete a cost without permission")

		return nil, ErrUnauthorized
	}

	cost, err := c.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	if cost.Annotations[AnnotationReadOnly] == AnnotationValueTrue {
		return nil, validation.NewError("the cost can not be deleted").WithFieldError(validation.FieldRoot, validation.ReadOnly, "cost is read-only")
	}

	if err := c.Store().Client().Delete(ctx, store.DeleteOptions.From(cost)); err != nil {
		log.WithError(err).Error("failed to delete the cost")

		return nil, err
	}

	return cost, nil
}

// Get returns the cost
func (c costsImpl) Get(ctx context.Context, name string) (*costsv1.Cost, error) {
	cost := &costsv1.Cost{}

	if found, err := c.Has(ctx, name); err != nil {
		return nil, err
	} else if !found {
		return nil, ErrNotFound
	}

	return cost, c.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(HubAdminTeam),
		store.GetOptions.WithName(name),
		store.GetOptions.InTo(cost),
	)
}

// List returns a list of costs
func (c costsImpl) List(ctx context.Context) (*costsv1.CostList, error) {
	costs := &costsv1.CostList{}

	return costs, c.Store().Client().List(ctx,
		store.ListOptions.InNamespace(HubAdminTeam),
		store.ListOptions.InTo(costs),
	)
}

// Has checks if a cost exists
func (c costsImpl) Has(ctx context.Context, name string) (bool, error) {
	return c.Store().Client().Has(ctx,
		store.HasOptions.InNamespace(HubAdminTeam),
		store.HasOptions.From(&costsv1.Cost{}),
		store.HasOptions.WithName(name),
	)
}
