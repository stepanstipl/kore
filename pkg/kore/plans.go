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

	"github.com/appvia/kore/pkg/kore/assets"
	"github.com/appvia/kore/pkg/utils/jsonschema"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/store"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

// Plans is the interface to the class plans
type Plans interface {
	// Delete is used to delete a plan in the kore
	Delete(context.Context, string) (*configv1.Plan, error)
	// Get returns the class from the kore
	Get(context.Context, string) (*configv1.Plan, error)
	// List returns a list of classes
	List(context.Context) (*configv1.PlanList, error)
	// Has checks if a resource exists within an available class in the scope
	Has(context.Context, string) (bool, error)
	// Update is responsible for update a plan in the kore
	Update(context.Context, *configv1.Plan) error
}

type plansImpl struct {
	Interface
}

// Update is responsible for update a plan in the kore
func (p plansImpl) Update(ctx context.Context, plan *configv1.Plan) error {
	plan.Namespace = HubNamespace

	user := authentication.MustGetIdentity(ctx)
	if !user.IsGlobalAdmin() {
		log.WithField("user", user.Username()).Warn("trying to update a plan without permissions")

		return ErrUnauthorized
	}

	switch plan.Spec.Kind {
	case "GKE":
		if err := jsonschema.Validate(assets.GKEPlanSchema, "plan", plan.Spec.Values.Raw); err != nil {
			return err
		}
	case "EKS":
		// TODO: add the EKS Plan schema and validate the plan parameters
	}

	err := p.Store().Client().Update(ctx,
		store.UpdateOptions.To(plan),
		store.UpdateOptions.WithCreate(true),
		store.UpdateOptions.WithForce(true),
	)
	if err != nil {
		log.WithError(err).Error("trying to update a plan in the kore")

		return err
	}

	return nil
}

// Delete is used to delete a plan in the kore
func (p plansImpl) Delete(ctx context.Context, name string) (*configv1.Plan, error) {
	user := authentication.MustGetIdentity(ctx)
	if !user.IsGlobalAdmin() {
		log.WithField("user", user.Username()).Warn("trying to delete a plan without permission")

		return nil, ErrUnauthorized
	}

	plan := &configv1.Plan{}
	err := p.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(HubNamespace),
		store.GetOptions.InTo(plan),
		store.GetOptions.WithName(name),
	)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		log.WithError(err).Error("trying to retrieve plan in the kore")

		return nil, err
	}

	if err := p.Store().Client().Delete(ctx, store.DeleteOptions.From(plan)); err != nil {
		log.WithError(err).Error("trying to delete the plan from kore")

		return nil, err
	}

	return plan, nil
}

// Get returns the class from the kore
func (p plansImpl) Get(ctx context.Context, name string) (*configv1.Plan, error) {
	plan := &configv1.Plan{}

	if found, err := p.Has(ctx, name); err != nil {
		return nil, err
	} else if !found {
		return nil, ErrNotFound
	}

	return plan, p.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(HubNamespace),
		store.GetOptions.WithName(name),
		store.GetOptions.InTo(plan),
	)
}

// List returns a list of classes
func (p plansImpl) List(ctx context.Context) (*configv1.PlanList, error) {
	plans := &configv1.PlanList{}

	return plans, p.Store().Client().List(ctx,
		store.ListOptions.InNamespace(HubNamespace),
		store.ListOptions.InTo(plans),
	)
}

// Has checks if a resource exists within an available class in the scope
func (p plansImpl) Has(ctx context.Context, name string) (bool, error) {
	return p.Store().Client().Has(ctx,
		store.HasOptions.InNamespace(HubNamespace),
		store.HasOptions.From(&configv1.Plan{}),
		store.HasOptions.WithName(name),
	)
}
