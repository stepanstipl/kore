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

	accountv1beta1 "github.com/appvia/kore/pkg/apis/accounts/v1beta1"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

// Accounts is the interface to the account accounts
type Accounts interface {
	// Delete is used to delete a account in the kore
	Delete(context.Context, string) (*accountv1beta1.AccountManagement, error)
	// Get returns the account from the kore
	Get(context.Context, string) (*accountv1beta1.AccountManagement, error)
	// List returns a list of accountes
	List(context.Context) (*accountv1beta1.AccountManagementList, error)
	// Has checks if a resource exists within an available account in the scope
	Has(context.Context, string) (bool, error)
	// Update is responsible for update a account in the kore
	Update(context.Context, *accountv1beta1.AccountManagement) error
}

type accountsImpl struct {
	Interface
}

// SupportedAccountProviders returns a list of supported provides
func (a accountsImpl) SupportedAccountProviders() []string {
	return []string{"GKE", "EKS", "AKS"}
}

// Update is responsible for update a account in the kore
func (a accountsImpl) Update(ctx context.Context, account *accountv1beta1.AccountManagement) error {
	account.Namespace = HubAdminTeam

	// @step: only admins can update the accounts
	user := authentication.MustGetIdentity(ctx)
	if !user.IsGlobalAdmin() {
		log.WithField("user", user.Username()).Warn("trying to update a account without permissions")

		return ErrUnauthorized
	}

	// @check: the provider is valid
	if !utils.Contains(account.Spec.Provider, a.SupportedAccountProviders()) {
		return fmt.Errorf("unsupported provider: %q, permitted: %s",
			account.Spec.Provider,
			strings.Join(a.SupportedAccountProviders(), ","),
		)
	}

	// @check: the plans exist
	for _, x := range account.Spec.Rules {
		for _, p := range x.Plans {
			found, err := a.Plans().Has(ctx, p)
			if err != nil {
				return err
			}
			if !found {
				return fmt.Errorf("plan %q does not exist", p)
			}
		}
	}

	err := a.Store().Client().Update(ctx,
		store.UpdateOptions.To(account),
		store.UpdateOptions.WithCreate(true),
	)
	if err != nil {
		log.WithError(err).Error("trying to update a account in the kore")

		return err
	}

	return nil
}

// Delete is used to delete a account in the kore
func (a accountsImpl) Delete(ctx context.Context, name string) (*accountv1beta1.AccountManagement, error) {
	user := authentication.MustGetIdentity(ctx)
	if !user.IsGlobalAdmin() {
		log.WithField("user", user.Username()).Warn("trying to delete a account without permission")

		return nil, ErrUnauthorized
	}

	account := &accountv1beta1.AccountManagement{}
	err := a.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(HubAdminTeam),
		store.GetOptions.InTo(account),
		store.GetOptions.WithName(name),
	)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		log.WithError(err).Error("trying to retrieve account in the kore")

		return nil, err
	}

	if err := a.Store().Client().Delete(ctx, store.DeleteOptions.From(account)); err != nil {
		log.WithError(err).Error("trying to delete the account from kore")

		return nil, err
	}

	return account, nil
}

// Get returns the account from the kore
func (a accountsImpl) Get(ctx context.Context, name string) (*accountv1beta1.AccountManagement, error) {
	account := &accountv1beta1.AccountManagement{}

	return account, a.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(HubAdminTeam),
		store.GetOptions.WithName(name),
		store.GetOptions.InTo(account),
	)
}

// List returns a list of accounts
func (a accountsImpl) List(ctx context.Context) (*accountv1beta1.AccountManagementList, error) {
	accounts := &accountv1beta1.AccountManagementList{}

	return accounts, a.Store().Client().List(ctx,
		store.ListOptions.InNamespace(HubAdminTeam),
		store.ListOptions.InTo(accounts),
	)
}

// Has checks if a resource exists within an available account in the scope
func (a accountsImpl) Has(ctx context.Context, name string) (bool, error) {
	return a.Store().Client().Has(ctx,
		store.HasOptions.InNamespace(HubAdminTeam),
		store.HasOptions.From(&accountv1beta1.AccountManagement{}),
		store.HasOptions.WithName(name),
	)
}
