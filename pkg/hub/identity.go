/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of hub-apiserver.
 *
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package hub

import (
	"context"

	"github.com/appvia/kore/pkg/hub/authentication"
	"github.com/appvia/kore/pkg/services/users"
	"github.com/appvia/kore/pkg/services/users/model"
	"github.com/appvia/kore/pkg/utils"
)

type identImpl struct {
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

func (i identImpl) IsGlobalAdmin() bool {
	if i.user.Username == "admin" {
		return true
	}

	return utils.Contains(HubAdminTeam, i.teamNames)
}

// GetUserIdentity queries the user services for the identity
func (h *hubImpl) GetUserIdentity(ctx context.Context, username string) (authentication.Identity, bool, error) {
	// @step: retrieve the user from the service
	user, err := h.usermgr.Users().Get(ctx, username)
	if err != nil {
		if users.IsNotFound(err) {
			return nil, false, nil
		}

		return nil, false, err
	}

	// @step: retrieve the teams the user is in
	teams, err := h.usermgr.Members().List(ctx,
		users.Filter.WithUser(username),
	)
	if err != nil {
		return nil, false, err
	}

	list := make([]string, len(teams))
	for i := 0; i < len(teams); i++ {
		list[i] = teams[i].Team.Name
	}

	return &identImpl{
		user:      user,
		teams:     teams,
		teamNames: list,
	}, true, nil
}

func getAdminContext(ctx context.Context) context.Context {
	ident := &identImpl{
		user: &model.User{
			Username: "admin",
		},
	}

	return context.WithValue(ctx, authentication.ContextKey{}, ident)
}
