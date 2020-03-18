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
	"time"

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/services/users"
	"github.com/appvia/kore/pkg/services/users/model"

	log "github.com/sirupsen/logrus"
)

// Users is the kore api users interface
type Users interface {
	// EnableUser is used to create an user in the kore
	EnableUser(context.Context, string, string) error
	// Delete removes the user from the kore
	Delete(context.Context, string) (*orgv1.User, error)
	// Exist checks if the user exists
	Exists(context.Context, string) (bool, error)
	// Get returns the user from the kore
	Get(context.Context, string) (*orgv1.User, error)
	// List returns a list of users
	List(context.Context) (*orgv1.UserList, error)
	// ListInvitations returns a list of invitations for a user
	ListInvitations(context.Context, string) (*orgv1.TeamInvitationList, error)
	// ListTeams returns the teams the user is in
	ListTeams(context.Context, string) (*orgv1.TeamList, error)
	// Update is responsible for updating the user
	Update(context.Context, *orgv1.User) (*orgv1.User, error)
}

// usersImpl provides the user implementation
type usersImpl struct {
	*hubImpl
}

// EnableUser is used to create an user in the kore
func (h *usersImpl) EnableUser(ctx context.Context, username, email string) error {
	logger := log.WithFields(log.Fields{
		"email":    email,
		"username": username,
	})
	logger.Info("enabling the user in the kore")

	found, err := h.Users().Exists(ctx, username)
	if err != nil {
		logger.WithError(err).Error("trying to check for user")

		return err
	}
	if found {
		logger.Debug("user already exists, no need to continue")

		return nil
	}

	if !found {
		logger.Debug("provisioning the user in the kore")

		if err := h.usermgr.Users().Update(ctx, &model.User{
			Username: username,
			Email:    email,
		}); err != nil {
			logger.WithError(err).Error("trying to create the user in the kore")

			return err
		}

		// @step: check for the user count - if this is the first user (minus admin)
		// they should be placed into the admin group
		count, err := h.usermgr.Users().Size(ctx)
		if err != nil {
			log.WithError(err).Error("trying to get a count on the kore users")

			return err
		}
		logger.WithField("count", count).Debug("we have x users already in the kore")

		isAdmin := count == 2
		roles := []string{"members"}
		if isAdmin {
			logger.Info("enabling the first user in the kore and providing admin access")

			// Add a custom audit for this special operation:
			start := time.Now()
			responseCode := 500
			defer func() {
				finish := time.Now()
				h.Audit().Record(ctx,
					users.Resource("/users"),
					users.ResourceURI("/users/"+username),
					users.Verb("PUT"),
					users.Operation("InitialiseFirstUserAsAdmin"),
					users.User(username),
					users.StartedAt(start),
					users.CompletedAt(finish),
					users.ResponseCode(responseCode),
				).Event("InitialiseFirstUserAsAdmin: Adding first user as administrator")
			}()

			if err := h.usermgr.Members().AddUser(ctx, username, HubAdminTeam, roles); err != nil {
				logger.WithError(err).Error("trying to add user to admin team")

				return err
			}
			responseCode = 200
		} else {
			logger.Info("adding the user into the kore")

			if err := h.usermgr.Teams().AddUser(ctx, username, HubDefaultTeam, roles); err != nil {
				logger.WithError(err).Error("trying to add user to default team")

				return err
			}
		}
	}

	return nil
}

// Get returns the user from the kore
func (h *usersImpl) Get(ctx context.Context, username string) (*orgv1.User, error) {
	user, err := h.usermgr.Users().Get(ctx, username)
	if err != nil {
		if users.IsNotFound(err) {
			return nil, ErrNotFound
		}
		log.WithError(err).Error("trying to retrieve the user")

		return nil, err
	}

	return DefaultConvertor.FromUserModel(user), nil
}

// List returns a list of users
func (h *usersImpl) List(ctx context.Context) (*orgv1.UserList, error) {
	list, err := h.usermgr.Users().List(ctx)
	if err != nil {
		log.WithError(err).Error("trying to retrieve a list of users")

		return nil, err
	}

	return DefaultConvertor.FromUsersModelList(list), err
}

// ListInvitations returns a list of team memberships for a user
func (h *usersImpl) ListInvitations(ctx context.Context, username string) (*orgv1.TeamInvitationList, error) {
	// @step: check the user exists
	if found, err := h.Exists(ctx, username); err != nil {
		return nil, err
	} else if !found {
		return nil, ErrNotFound
	}

	list, err := h.usermgr.Invitations().List(ctx,
		users.Filter.WithUser(username),
	)
	if err != nil {
		log.WithError(err).Error("trying to list the invitations for user")

		return nil, err
	}

	return DefaultConvertor.FromInvitationModelList(list), nil
}

// Delete removes the user from the kore
func (h *usersImpl) Delete(ctx context.Context, username string) (*orgv1.User, error) {
	// @step: check the user exists
	u, err := h.usermgr.Users().Get(ctx, username)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"username": u.Username,
	}).Info("deleting the user from the kore")

	teams, err := h.Users().ListTeams(ctx, username)
	if err != nil {
		return nil, err
	}

	for _, x := range teams.Items {
		team := x.Name
		if err := h.Teams().Team(team).Members().Delete(ctx, username); err != nil {
			return nil, fmt.Errorf("failed to delete team membership: %s", err)
		}
	}

	if _, err := h.usermgr.Users().Delete(ctx, u); err != nil {
		log.WithError(err).Error("trying to remove user from kore")

		return nil, err
	}

	// @TODO add an entry into the audit log

	return DefaultConvertor.FromUserModel(u), nil
}

// Update is responsible for updating the user
func (h *usersImpl) Update(ctx context.Context, user *orgv1.User) (*orgv1.User, error) {
	user.Namespace = HubNamespace

	// @step: we need to validate the user
	if user.Spec.Username == "" || user.Name == "" {
		return nil, NewErrNotAllowed("user must have a username")
	}
	if user.Spec.Email == "" {
		return nil, NewErrNotAllowed("user must have a email address")
	}

	// @TODO add an entry into the audit

	// @step: update the user in the user management service
	if err := h.usermgr.Users().Update(ctx, DefaultConvertor.ToUserModel(user)); err != nil {
		log.WithError(err).Error("trying to update the user in the kore")

		return nil, err
	}

	return user, nil
}

// ListTeams return a list of teams the user is in
func (h *usersImpl) ListTeams(ctx context.Context, username string) (*orgv1.TeamList, error) {
	list, err := h.usermgr.Members().List(ctx,
		users.Filter.WithUser(username),
	)
	if err != nil {
		log.WithError(err).Error("trying to list the teams the user is in")

		return nil, err
	}

	return DefaultConvertor.FromMembersToTeamList(list), nil
}

// Exists checks if the user exists
func (h usersImpl) Exists(ctx context.Context, name string) (bool, error) {
	return h.usermgr.Users().Exists(ctx, name)
}
