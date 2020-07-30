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

	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/persistence/model"
	"github.com/appvia/kore/pkg/utils"
	log "github.com/sirupsen/logrus"
)

//
const (
	// IdentityBasicAuth is basicauth identity
	IdentityBasicAuth = "basicauth"
	// SSO is a single-sign-on identity
	IdentitySSO = "sso"
)

type identImpl struct {
	metadata  map[string]string
	user      *model.User
	teams     []*model.Member
	teamNames []string
}

func (i identImpl) Username() string {
	return i.user.Username
}

func (i identImpl) Email() string {
	return i.user.Email
}

func (i identImpl) Disabled() bool {
	return i.user.Disabled
}

func (i identImpl) Teams() []string {
	return i.teamNames
}

func (i identImpl) AuthMethod() string {
	return i.metadata["method"]
}

func (i identImpl) IsGlobalAdmin() bool {
	if i.user.Username == "admin" {
		return true
	}

	return utils.Contains(HubAdminTeam, i.teamNames)
}

// MetaFunc defined an meta func
type MetaFunc func(*identImpl)

// WithAuthMethod sets the auth method
func WithAuthMethod(v string) MetaFunc {
	return func(m *identImpl) {
		m.metadata["method"] = v
	}
}

// IsMember checks if the user is a member of the team
func (i identImpl) IsMember(name string) bool {
	return utils.Contains(name, i.Teams())
}

// GetUserIdentity queries the user services for the identity
func (h *hubImpl) GetUserIdentity(ctx context.Context, username string, meta ...MetaFunc) (authentication.Identity, bool, error) {
	// @step: retrieve the user from the service
	user, err := h.persistenceMgr.Users().Get(ctx, username)
	if err != nil {
		if !persistence.IsNotFound(err) {
			return nil, false, err
		}

		return nil, false, nil
	}

	filters := []persistence.ListFunc{
		persistence.Filter.WithUser(username),
	}

	// @step: check if the user is disabled
	if user.Disabled {
		log.WithField("username", username).Warn("disabled user attempting to login")

		return nil, false, nil
	}

	// @step: retrieve the teams the user is in
	teams, err := h.persistenceMgr.Members().List(ctx, filters...)
	if err != nil {
		return nil, false, err
	}

	list := make([]string, len(teams))
	for i := 0; i < len(teams); i++ {
		list[i] = teams[i].Team.Name
	}

	ident := &identImpl{
		metadata:  make(map[string]string),
		user:      user,
		teams:     teams,
		teamNames: list,
	}

	// @step: apply the metadata
	for _, x := range meta {
		x(ident)
	}

	return ident, true, nil
}

// GetUserIdentityByProvider returns the user model by proviser if any
func (h *hubImpl) GetUserIdentityByProvider(ctx context.Context, username, provider string) (*model.Identity, bool, error) {
	id, err := h.persistenceMgr.Identities().Get(ctx,
		persistence.Filter.WithUser(username),
		persistence.Filter.WithProvider(provider),
	)
	if err != nil {
		if !persistence.IsNotFound(err) {
			return nil, false, err
		}

		return nil, false, nil
	}

	return id, true, nil
}

func getAdminContext(ctx context.Context) context.Context {
	ident := &identImpl{
		user: &model.User{
			Username: "admin",
		},
	}

	return context.WithValue(ctx, authentication.ContextKey{}, ident)
}
