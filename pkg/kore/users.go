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
	"regexp"
	"time"

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/persistence/model"
	"github.com/appvia/kore/pkg/utils/validation"

	log "github.com/sirupsen/logrus"
)

var (
	// UsernameRegex is a filter for the username
	UsernameRegex = regexp.MustCompile(`^[a-z-A-Z0-9\@\-\.]{3,63}$`)
)

// Users is the kore api users interface
type Users interface {
	// EnableUser is used to create an user in kore
	EnableUser(context.Context, string, string) error
	// Delete removes the user from the kore
	Delete(context.Context, string) (*orgv1.User, error)
	// Exist checks if the user exists
	Exists(context.Context, string) (bool, error)
	// Get returns the user from the kore
	Get(context.Context, string) (*orgv1.User, error)
	// Identities returns the identities interface
	Identities() Identities
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

// Identities returns the identities interface
func (h *usersImpl) Identities() Identities {
	return &idImpl{hubImpl: h.hubImpl}
}

// EnableUser is used to create an user in the kore
func (h *usersImpl) EnableUser(ctx context.Context, username, email string) error {
	logger := log.WithFields(log.Fields{
		"email":    email,
		"username": username,
	})
	logger.Info("enabling the user in kore")

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
		logger.Debug("provisioning the user in kore")

		if err := h.persistenceMgr.Users().Update(ctx, &model.User{
			Username: username,
			Email:    email,
		}); err != nil {
			logger.WithError(err).Error("trying to create the user in kore")

			return err
		}

		// @step: check for the user count - if this is the first user (minus admin)
		// they should be placed into the admin group
		count, err := h.persistenceMgr.Users().Size(ctx)
		if err != nil {
			log.WithError(err).Error("trying to get a count on the kore users")

			return err
		}
		logger.WithField("count", count).Debug("we have x users already in kore")

		isAdmin := count == 2
		roles := []string{"members"}
		if isAdmin {
			logger.Info("enabling the first user in kore and providing admin access")

			// Add a custom audit for this special operation:
			start := time.Now()
			responseCode := 500
			defer func() {
				finish := time.Now()
				h.Audit().Record(ctx,
					persistence.Resource("/users"),
					persistence.ResourceURI("/users/"+username),
					persistence.Verb("PUT"),
					persistence.Operation("InitialiseFirstUserAsAdmin"),
					persistence.User(username),
					persistence.StartedAt(start),
					persistence.CompletedAt(finish),
					persistence.ResponseCode(responseCode),
				).Event("InitialiseFirstUserAsAdmin: Adding first user as administrator")
			}()

			if err := h.persistenceMgr.Members().AddUser(ctx, username, HubAdminTeam, roles); err != nil {
				logger.WithError(err).Error("trying to add user to admin team")

				return err
			}
			responseCode = 200
		} else {
			logger.Info("adding the user into the kore")

			if err := h.persistenceMgr.Teams().AddUser(ctx, username, HubDefaultTeam, roles); err != nil {
				logger.WithError(err).Error("trying to add user to default team")

				return err
			}
		}
	}

	return nil
}

// Get returns the user from the kore
func (h *usersImpl) Get(ctx context.Context, username string) (*orgv1.User, error) {
	user, err := h.persistenceMgr.Users().Get(ctx, username)
	if err != nil {
		if persistence.IsNotFound(err) {
			return nil, ErrNotFound
		}
		log.WithError(err).Error("trying to retrieve the user")

		return nil, err
	}

	return DefaultConvertor.FromUserModel(user), nil
}

// List returns a list of users
func (h *usersImpl) List(ctx context.Context) (*orgv1.UserList, error) {
	list, err := h.persistenceMgr.Users().List(ctx)
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

	list, err := h.persistenceMgr.Invitations().List(ctx,
		persistence.Filter.WithUser(username),
	)
	if err != nil {
		log.WithError(err).Error("trying to list the invitations for user")

		return nil, err
	}

	return DefaultConvertor.FromInvitationModelList(list), nil
}

// Delete removes the user from the kore
func (h *usersImpl) Delete(ctx context.Context, username string) (*orgv1.User, error) {
	if !authentication.MustGetIdentity(ctx).IsGlobalAdmin() {
		return nil, ErrUnauthorized
	}

	// @step: you should be permitted to delete the admin user
	if username == "admin" {
		return nil, NewErrNotAllowed("you are not permitted to delete the admin user")
	}

	// @step: check the user exists
	u, err := h.persistenceMgr.Users().Get(ctx, username)
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

	if _, err := h.persistenceMgr.Users().Delete(ctx, u); err != nil {
		log.WithError(err).Error("trying to remove user from kore")

		return nil, err
	}

	// @TODO add an entry into the audit log

	return DefaultConvertor.FromUserModel(u), nil
}

// Update is responsible for updating the user
func (h *usersImpl) Update(ctx context.Context, user *orgv1.User) (*orgv1.User, error) {
	if !authentication.MustGetIdentity(ctx).IsGlobalAdmin() {
		return nil, ErrUnauthorized
	}

	user.Namespace = HubNamespace

	// @step: we need to validate the user
	valErr := validation.NewError("user has failed validation")
	if user.Name == "" {
		valErr.AddFieldError("metadata.name", validation.Required, "can not be empty")
	}
	if user.Spec.Username == "" {
		valErr.AddFieldError("spec.username", validation.Required, "can not be empty")
	}
	if user.Spec.Email == "" {
		valErr.AddFieldError("spec.email", validation.Required, "can not be empty")
	}
	if !UsernameRegex.MatchString(user.Spec.Username) {
		valErr.AddFieldError("spec.username", validation.InvalidValue, "invalid username")
	}
	if valErr.HasErrors() {
		return nil, valErr
	}

	// @step: update the user in the user management service
	if err := h.persistenceMgr.Users().Update(ctx, DefaultConvertor.ToUserModel(user)); err != nil {
		log.WithError(err).Error("trying to update the user in kore")

		return nil, err
	}

	return user, nil
}

// ListTeams return a list of teams the user is in
func (h *usersImpl) ListTeams(ctx context.Context, username string) (*orgv1.TeamList, error) {
	list, err := h.persistenceMgr.Members().List(ctx,
		persistence.Filter.WithUser(username),
	)
	if err != nil {
		log.WithError(err).Error("trying to list the teams the user is in")

		return nil, err
	}

	return DefaultConvertor.FromMembersToTeamList(list), nil
}

// Exists checks if the user exists
func (h usersImpl) Exists(ctx context.Context, name string) (bool, error) {
	return h.persistenceMgr.Users().Exists(ctx, name)
}
