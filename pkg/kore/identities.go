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
	"time"

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/persistence/model"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/validation"
	jwt "github.com/dgrijalva/jwt-go"

	log "github.com/sirupsen/logrus"
)

var (
	// AccountLocal is a local basic auth account
	AccountLocal = "basicauth"
	// AccountToken is a api token account
	AccountToken = "token"
	// AccountSSO is a openid account
	AccountSSO = "sso"
	// SupportedAccounts is a list of supported accounts
	SupportedAccounts = []string{
		AccountLocal,
		AccountToken,
		AccountSSO,
	}
)

// IdentitiesListOptions are search options for listing
type IdentitiesListOptions struct {
	// IdentityTypes is a collection of type to search for
	IdentityTypes []string
	// User is a specific user to user to look for
	User string
}

// Identities is the contract to interact with user identities
type Identities interface {
	// AssociateIDPUser is used to associate an internal user to an idp user
	AssociateIDPUser(ctx context.Context, update *orgv1.UpdateIDPIdentity) error
	// Delete is called to delete an associated identity of a user
	Delete(ctx context.Context, user string, identity string) error
	// IssueToken is used to issue a token for a identity in kore
	IssueToken(ctx context.Context, audience string, scopes []string) ([]byte, error)
	// List returns a list of all the identities managed in kore
	List(ctx context.Context, options IdentitiesListOptions) (*orgv1.IdentityList, error)
	// UpdateUserBasicAuth is used to update a basic auth profile in kore
	UpdateUserBasicAuth(ctx context.Context, update *orgv1.UpdateBasicAuthIdentity) error
}

type idImpl struct {
	*hubImpl
}

// AssociateIDPUser is used to associate an internal user to an idp user
func (h *idImpl) AssociateIDPUser(ctx context.Context, update *orgv1.UpdateIDPIdentity) error {

	return nil
}

// IssueToken is used to issue a token for a identity in kore
func (h *idImpl) IssueToken(ctx context.Context, audience string, scopes []string) ([]byte, error) {
	user := authentication.MustGetIdentity(ctx)

	if user.AuthMethod() != AccountLocal {
		return nil, NewErrNotAllowed("only basicauth identities can be issues at present")
	}

	usercl := "preferred_username"
	claims := utils.NewClaims(jwt.MapClaims{
		"aud":    audience,
		"email":  user.Email(),
		"exp":    float64(time.Now().UTC().Add(24 * time.Hour).Unix()),
		"iss":    h.Config().PublicAPIURL,
		"nbf":    time.Now().UTC().Add(-60 * time.Second).Unix(),
		"scopes": scopes,
		usercl:   user.Username(),
	})

	minted, err := claims.Sign(h.CertificateAuthorityKey())
	if err != nil {
		log.WithField("user", user.Username()).WithError(err).Error("trying to mint local token")

		return nil, err
	}
	log.WithFields(log.Fields{
		"email":    user.Email(),
		"username": user.Username(),
	}).Info("successfully minted a token to local user")

	return minted, nil
}

// Delete is called to delete an associated identity of a user
func (h *idImpl) Delete(ctx context.Context, username string, identity string) error {
	user := authentication.MustGetIdentity(ctx)

	// @step: you must be the user or an admin to perform this
	if !user.IsGlobalAdmin() && user.Username() != username {
		return NewErrNotAllowed("must be administrator or the user to delete credential")
	}

	// @step: check the identity type and username
	if !utils.Contains(identity, SupportedAccounts) {
		return validation.NewError("invalid identity").
			WithFieldError("identity", validation.InvalidValue, "identity type does not exist")
	}
	if !UsernameRegex.MatchString(username) {
		return validation.NewError("invalid username").
			WithFieldError("username", validation.InvalidValue, "username is invalid")
	}

	// @step: check the user exists
	_, err := h.Persist().Users().Get(ctx, username)
	if err != nil {
		if persistence.IsNotFound(err) {
			return ErrNotFound
		}
		log.WithError(err).Error("trying to check if user exists")

		return err
	}

	// @step: retrieve the identity if any
	ident, err := h.Persist().Identities().Get(ctx,
		persistence.Filter.WithUser(username),
		persistence.Filter.WithProvider(identity),
	)
	if err != nil {
		if persistence.IsNotFound(err) {
			return NewErrNotAllowed("user does not have this identity type")
		}
	}

	return h.Persist().Identities().Delete(ctx, ident)
}

// List returns a list of all the identities managed in kore
func (h *idImpl) List(ctx context.Context, options IdentitiesListOptions) (*orgv1.IdentityList, error) {
	user := authentication.MustGetIdentity(ctx)
	var filters []persistence.ListFunc

	// @step: validate inputs
	for _, x := range options.IdentityTypes {
		if !utils.Contains(x, SupportedAccounts) {
			return nil, validation.NewError("invalid identity type").
				WithFieldError("type", validation.InvalidValue, "must be a valid identity type")
		}
		filters = append(filters, persistence.Filter.WithProvider(x))
	}

	if options.User != "" {
		if !UsernameRegex.MatchString(options.User) {
			return nil, validation.NewError("invalid username").
				WithFieldError("username", validation.InvalidValue, "username is invalid")
		}
		filters = append(filters, persistence.Filter.WithUser(options.User))
	}

	if options.User == "" && !user.IsGlobalAdmin() {
		return nil, ErrUnauthorized
	}

	list, err := h.Persist().Identities().List(ctx, filters...)
	if err != nil {
		return nil, err
	}

	return DefaultConvertor.FromIdentityModelList(list), err
}

// UpdateUserBasicAuth is used to update a basic auth profile in kore
func (h *idImpl) UpdateUserBasicAuth(ctx context.Context, update *orgv1.UpdateBasicAuthIdentity) error {
	user := authentication.MustGetIdentity(ctx)

	if !user.IsGlobalAdmin() && user.Username() != update.Username {
		return NewErrNotAllowed("must be administrator or the user to update credential")
	}

	logger := log.WithFields(log.Fields{
		"username": update.Username,
	})

	// @step: check the user exists
	u, err := h.Persist().Users().Get(ctx, update.Username)
	if err != nil {
		if persistence.IsNotFound(err) {
			return ErrNotFound
		}
		logger.WithError(err).Error("trying to check if user exists")

		return err
	}

	// @step: if the user has zero identites only the admin can it up
	identity, err := h.Persist().Identities().Get(ctx,
		persistence.Filter.WithUser(update.Username),
		persistence.Filter.WithProvider(AccountLocal),
	)
	if err != nil {
		if !persistence.IsNotFound(err) {
			logger.WithError(err).Error("trying to retrieve the identity")

			return err
		}
		logger.Info("setting up basicauth identity for user")
	}

	if identity == nil {
		identity = &model.Identity{
			Provider: AccountLocal,
			UserID:   u.ID,
		}
	}
	identity.ProviderToken = update.Password

	// @step: update the credentials
	if err := h.Persist().Identities().Update(ctx, identity); err != nil {
		logger.WithError(err).Error("trying to update the credential")

		return err
	}
	logger.Info("updated the basicauth credential for user")

	return nil
}
