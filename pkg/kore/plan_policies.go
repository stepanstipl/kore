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

	"github.com/appvia/kore/pkg/utils/validation"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/store"

	log "github.com/sirupsen/logrus"
)

// PlanPolicies is the interface to the plan policies
type PlanPolicies interface {
	// Delete is used to delete a plan policy
	Delete(context.Context, string) (*configv1.PlanPolicy, error)
	// Get returns the plan policy
	Get(context.Context, string) (*configv1.PlanPolicy, error)
	// List returns a list of plan policies
	List(context.Context) (*configv1.PlanPolicyList, error)
	// Has checks if a plan policy
	Has(context.Context, string) (bool, error)
	// Update is responsible for updating a plan policy
	Update(ctx context.Context, planPolicy *configv1.PlanPolicy, ignoreReadonly bool) error
}

type planPoliciesImpl struct {
	Interface
}

// Update is responsible for updating a plan policy
func (p planPoliciesImpl) Update(ctx context.Context, planPolicy *configv1.PlanPolicy, ignoreReadonly bool) error {
	planPolicy.Namespace = HubAdminTeam

	// @TODO: check the user is admin or has kore permissions
	user := authentication.MustGetIdentity(ctx)
	if !user.IsGlobalAdmin() {
		log.WithField("user", user.Username()).Warn("failed to update the plan policy policy without permissions")

		return ErrUnauthorized
	}

	if !ignoreReadonly {
		original, err := p.Get(ctx, planPolicy.Name)
		if err != nil && err != ErrNotFound {
			return err
		}

		if original != nil && original.Annotations[AnnotationReadOnly] == AnnotationValueTrue {
			return validation.NewError("the plan policy can not be updated").WithFieldError(validation.FieldRoot, validation.ReadOnly, "policy is read-only")
		}
		if planPolicy.Annotations[AnnotationReadOnly] == AnnotationValueTrue {
			return validation.NewError("the plan policy can not be updated").WithFieldError(validation.FieldRoot, validation.ReadOnly, "read-only flag can not be set")
		}
	}

	err := p.Store().Client().Update(ctx,
		store.UpdateOptions.To(planPolicy),
		store.UpdateOptions.WithCreate(true),
		store.UpdateOptions.WithForce(true),
	)
	if err != nil {
		log.WithError(err).Error("failed to update the plan policy in the kore")

		return err
	}

	return nil
}

// Delete is used to delete a plan policy
func (p planPoliciesImpl) Delete(ctx context.Context, name string) (*configv1.PlanPolicy, error) {
	user := authentication.MustGetIdentity(ctx)
	if !user.IsGlobalAdmin() {
		log.WithField("user", user.Username()).Warn("failed to delete a plan policy without permission")

		return nil, ErrUnauthorized
	}

	planPolicy, err := p.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	if planPolicy.Annotations[AnnotationReadOnly] == AnnotationValueTrue {
		return nil, validation.NewError("the plan policy can not be deleted").WithFieldError(validation.FieldRoot, validation.ReadOnly, "policy is read-only")
	}

	if err := p.Store().Client().Delete(ctx, store.DeleteOptions.From(planPolicy)); err != nil {
		log.WithError(err).Error("failed to delete the plan policy")

		return nil, err
	}

	return planPolicy, nil
}

// Get returns the plan policy
func (p planPoliciesImpl) Get(ctx context.Context, name string) (*configv1.PlanPolicy, error) {
	planPolicy := &configv1.PlanPolicy{}

	if found, err := p.Has(ctx, name); err != nil {
		return nil, err
	} else if !found {
		return nil, ErrNotFound
	}

	return planPolicy, p.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(HubAdminTeam),
		store.GetOptions.WithName(name),
		store.GetOptions.InTo(planPolicy),
	)
}

// List returns a list of plan policies
func (p planPoliciesImpl) List(ctx context.Context) (*configv1.PlanPolicyList, error) {
	planPolicies := &configv1.PlanPolicyList{}

	return planPolicies, p.Store().Client().List(ctx,
		store.ListOptions.InNamespace(HubAdminTeam),
		store.ListOptions.InTo(planPolicies),
	)
}

// Has checks if a plan policy
func (p planPoliciesImpl) Has(ctx context.Context, name string) (bool, error) {
	return p.Store().Client().Has(ctx,
		store.HasOptions.InNamespace(HubAdminTeam),
		store.HasOptions.From(&configv1.PlanPolicy{}),
		store.HasOptions.WithName(name),
	)
}
