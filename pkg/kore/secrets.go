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
	"github.com/appvia/kore/pkg/kore/assets"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils"

	log "github.com/sirupsen/logrus"
)

// Secrets is the interface to the class plans
type Secrets interface {
	// Delete is used to delete a plan in kore
	Delete(context.Context, string) (*configv1.Secret, error)
	// Get returns the class from the kore
	Get(context.Context, string) (*configv1.Secret, error)
	// Has checks if a resource exists within an available class in the scope
	Has(context.Context, string) (bool, error)
	// List returns a list of classes
	List(context.Context) (*configv1.SecretList, error)
	// SupportedSecretTypes returns a list of supported types
	SupportedSecretTypes() []string
	// Update is responsible for update a plan in kore
	Update(context.Context, *configv1.Secret) error
}

type secretImpl struct {
	*hubImpl
	// team is the team we are residing
	team string
}

// SupportedSecretTypes returns a list of supported types
func (h *secretImpl) SupportedSecretTypes() []string {
	return assets.SupportedSecretTypesNames()
}

// Update is responsible for updating the resource
func (h *secretImpl) Update(ctx context.Context, secret *configv1.Secret) error {
	user := authentication.MustGetIdentity(ctx)
	logger := log.WithFields(log.Fields{
		"name": secret.Name,
		"team": h.team,
		"user": user.Username(),
	})
	logger.Info("attempting to update or create a secret in the team")

	if secret.Namespace == "" {
		secret.Namespace = h.team
	}

	if secret.Namespace != h.team {
		return ErrNotAllowed{message: "you cannot create a secret in another team"}
	}
	if secret.Spec.Type == "" {
		return ErrNotAllowed{message: "secret must have a type"}
	}
	if !utils.Contains(secret.Spec.Type, h.SupportedSecretTypes()) {
		return ErrNotAllowed{message: "secret type is unsupported"}
	}

	if !user.IsGlobalAdmin() {
		if secret.Status.SystemManaged != nil && *secret.Status.SystemManaged {
			return ErrNotAllowed{message: "managed secrets can only be changed by global admins"}
		}
	}

	h.Audit().Record(ctx,
		persistence.Resource("secrets/"+secret.Name),
		persistence.Team(h.team),
		persistence.User(user.Username()),
	).Event("user creating or updating secret in team")

	return h.Store().Client().Update(ctx,
		store.UpdateOptions.To(secret),
		store.UpdateOptions.WithCreate(true),
		store.UpdateOptions.WithForce(true),
	)
}

// Delete is responsible for deleting a gke environment
func (h *secretImpl) Delete(ctx context.Context, name string) (*configv1.Secret, error) {
	user := authentication.MustGetIdentity(ctx)

	logger := log.WithFields(log.Fields{
		"name": name,
		"team": h.team,
		"user": user.Username(),
	})
	logger.Info("attempting to delete the secret from the team")

	// @step: ensure the secret exists in the team
	secret := &configv1.Secret{}
	if err := h.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(h.team),
		store.GetOptions.InTo(secret),
		store.GetOptions.WithName(name)); err != nil {

		logger.WithError(err).Error("trying to check for secret in team")

		return nil, err
	}

	if !user.IsGlobalAdmin() {
		if secret.Status.SystemManaged != nil && *secret.Status.SystemManaged {
			return nil, ErrNotAllowed{message: "managed secrets can only be deleted by global admins"}
		}
	}

	// @TODO should this check if the secret is being used?
	return secret, h.Store().Client().Delete(ctx, store.DeleteOptions.From(secret))
}

// Get return the definition from the api
func (h *secretImpl) Get(ctx context.Context, name string) (*configv1.Secret, error) {
	user := authentication.MustGetIdentity(ctx)
	logger := log.WithFields(log.Fields{
		"name": name,
		"team": h.team,
		"user": user.Username(),
	})
	logger.Info("attempting to retrieve the secret from team")

	secret := &configv1.Secret{}

	err := h.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(h.team),
		store.GetOptions.InTo(secret),
		store.GetOptions.WithName(name),
	)
	if err != nil {
		log.WithError(err).Error("trying to retrieve secret from api")

		return nil, err
	}

	// @step: check the user have the permissions to retrieve the secret
	switch secret.Spec.Type {
	case configv1.KubernetesSecret:
		if !user.IsGlobalAdmin() {
			return nil, ErrNotAllowed{message: "permission denied, only global admin are retrieve these secrets"}
		}
	}

	return secret, nil
}

// Has checks if the resource exists
func (h *secretImpl) Has(ctx context.Context, name string) (bool, error) {
	return h.Store().Client().Has(ctx,
		store.HasOptions.InNamespace(h.team),
		store.HasOptions.From(&configv1.Secret{}),
		store.HasOptions.WithName(name),
	)
}

// List returns all the secrets in the team
func (h *secretImpl) List(ctx context.Context) (*configv1.SecretList, error) {
	list := &configv1.SecretList{}

	return list, h.Store().Client().List(ctx,
		store.ListOptions.InNamespace(h.team),
		store.ListOptions.InTo(list),
	)
}
