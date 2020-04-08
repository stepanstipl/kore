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
	"encoding/base64"
	"fmt"
	"time"

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/services/users"
	"github.com/appvia/kore/pkg/services/users/model"
	"github.com/appvia/kore/pkg/store"

	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

// TeamMembers returns the members interface for a team
type TeamMembers interface {
	// Add is responsible for adding a member to team
	Add(context.Context, string) error
	// Delete is responsible for removing an member from the team
	Delete(context.Context, string) error
	// DeleteInvitation deletes an invitation to the team
	DeleteInvitation(context.Context, string) error
	// Exists check if a member exists
	Exists(context.Context, string) (bool, error)
	// GenerateLink is called to generate a invitation link
	GenerateLink(context.Context, GenerateLinkOptions) (string, error)
	// Invite is responsible for creating an invitation to the team
	Invite(context.Context, string, InvitationOptions) error
	// List returns a list of memberships for the team
	List(context.Context) (*orgv1.UserList, error)
	// ListInvitations returns a list of invitations for the team
	ListInvitations(context.Context) (*orgv1.TeamInvitationList, error)
}

// GenerateLinkOptions are the options for the link
type GenerateLinkOptions struct {
	// Duration is the time the invition will take open
	Duration time.Duration
	// User is a specific username its being generated for
	User string
}

// InvitationOptions are the options for an invitation
type InvitationOptions struct {
	// Duration is the time the invition will take open
	Duration time.Duration
}

type tmsImpl struct {
	*hubImpl
	// team is the team name
	team string
}

// Add is responsible for adding a member to team
// @TODO we need to add a role to this later
func (t tmsImpl) Add(ctx context.Context, name string) error {
	// @step: construct the membership claim
	logger := log.WithFields(log.Fields{
		"team":     t.team,
		"username": name,
	})
	logger.Info("attempting user membership to the team")

	// @step: check the user exists
	user, err := t.usermgr.Users().Get(ctx, name)
	if err != nil {
		if !users.IsNotFound(err) {
			logger.WithError(err).Error("trying to retrieve the user from service")

			return err
		}

		return NewErrNotAllowed("user does not exist in the kore")
	}

	// @TODO we need to check if the caller is an member (and later admin) to
	// to the team

	// @step: we need to retrieve the team
	team, err := t.usermgr.Teams().Get(ctx, t.team)
	if err != nil {
		logger.WithError(err).Error("trying to retrieve the team")

		return err
	}

	// @step: add the user as a member of the team
	if err := t.usermgr.Members().Update(ctx, &model.Member{
		UserID: user.ID,
		TeamID: team.ID,
		Roles:  []string{"member"},
	}); err != nil {
		logger.WithError(err).Error("trying to add the user as a member of team")

		return err
	}

	if err := t.UpdateTeam(ctx); err != nil {
		logger.WithError(err).Error("trying to update the team in the api")

		return err
	}

	// @step: ensure any invitations for the user is removed
	return t.DeleteInvitation(ctx, name)
}

// Delete is responsible for removing an member from the team
func (t tmsImpl) Delete(ctx context.Context, name string) error {
	logger := log.WithFields(log.Fields{
		"team":     t.team,
		"username": name,
	})
	logger.Info("attempting to delete the user from the team")

	// @step: check for membership
	list, err := t.usermgr.Members().List(ctx,
		users.Filter.WithUser(name),
		users.Filter.WithTeam(t.team),
	)
	if err != nil {
		logger.WithError(err).Error("trying to check for user membership")

		return err
	}
	if len(list) <= 0 {
		return nil
	}

	// @check: is this the admin team?
	if t.team == HubAdminTeam {
		if name == HubAdminUser {
			return ErrNotAllowed{message: "you cannot remove " + HubAdminUser + " user from this team"}
		}
		if len(list) == 1 {
			return ErrNotAllowed{message: "you cannot remove all administrators from kore"}
		}
	}

	// @step: we can delete all the membership
	if err := t.usermgr.Members().DeleteBy(ctx,
		users.Filter.WithUser(name),
		users.Filter.WithTeam(t.team),
	); err != nil {
		logger.WithError(err).Error("trying to delete user membership from team")

		return err
	}

	// @TODO: we need to tap the team object in the api to get the controller
	// to reconcile the team
	if err := t.UpdateTeam(ctx); err != nil {
		logger.WithError(err).Error("trying to update the team in the api")

		return err
	}

	return nil
}

// GenerateLink is called to generate a invitation link
func (t tmsImpl) GenerateLink(ctx context.Context, options GenerateLinkOptions) (string, error) {
	logger := log.WithFields(log.Fields{
		"duration": options.Duration.String(),
		"team":     t.team,
		"user":     options.User,
	})
	if t.Config().HMAC == "" {
		return "", fmt.Errorf("kore has no hmac token for signature")
	}
	claims := jwt.MapClaims{
		"type": "invitation",
		"exp":  time.Now().Add(options.Duration).Unix(),
		"iat":  time.Now().Unix(),
		"iss":  "kore",
		"nbf":  time.Now().Unix(),
		"team": t.team,
	}
	if options.User != "" {
		claims["user"] = options.User
	}
	// @step: create a payload for the invitation
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	encoded, err := token.SignedString([]byte(t.Config().HMAC))
	if err != nil {
		logger.WithError(err).Error("failed to generate invitation link for team")

		return "", err
	}

	return base64.StdEncoding.EncodeToString([]byte(encoded)), nil
}

// Invite is responsible for creating an invitation to the team
func (t tmsImpl) Invite(ctx context.Context, name string, options InvitationOptions) error {
	// @step: check the options
	if options.Duration <= 0 {
		return ErrNotAllowed{message: "invalid expiry on invitation, must be greater than zero"}
	}
	logger := log.WithFields(log.Fields{
		"team":     t.team,
		"username": name,
	})
	logger.Info("adding user invitation to the team")

	// @step: check the user and team
	user, err := t.usermgr.Users().Get(ctx, name)
	if err != nil {
		if users.IsNotFound(err) {
			return ErrNotAllowed{message: "user does not exist in the appvia kore"}
		}
		logger.WithError(err).Error("trying to retrieve user")

		return err
	}
	team, err := t.usermgr.Teams().Get(ctx, t.team)
	if err != nil {
		if users.IsNotFound(err) {
			return ErrNotAllowed{message: "team does not exist in the appvia kore"}
		}
		logger.WithError(err).Error("trying to retrieve team from user service")

		return err
	}

	// @TODO need to add a audit entry about the invitation

	// @step: invite the user into the team - @TODO need to add a role as well
	if err := t.usermgr.Invitations().Update(ctx,
		&model.Invitation{
			Expires: time.Now().Add(options.Duration),
			TeamID:  team.ID,
			UserID:  user.ID,
		},
	); err != nil {
		logger.WithError(err).Error("trying to create an invitation for the team")

		return err
	}

	return nil
}

// DeleteInvitation deletes an user invitation in a team
func (t tmsImpl) DeleteInvitation(ctx context.Context, name string) error {
	logger := log.WithFields(log.Fields{
		"username": name,
	})
	logger.Debug("removing any invitations for user")

	err := t.usermgr.Invitations().DeleteBy(ctx,
		users.Filter.WithUser(name),
		users.Filter.WithTeam(t.team),
	)
	if err != nil {
		logger.WithError(err).Error("trying to delete invitations for team")

		return err
	}

	return nil
}

// Exists check if a member exists in the team
func (t tmsImpl) Exists(ctx context.Context, name string) (bool, error) {
	list, err := t.usermgr.Members().List(ctx,
		users.Filter.WithUser(name),
		users.Filter.WithTeam(t.team),
	)
	if err != nil {
		log.WithError(err).Error("trying to retrieve list of memberships")
	}

	if len(list) > 0 {
		return true, nil
	}

	return false, nil
}

// List returns a list of memberships for the team
func (t tmsImpl) List(ctx context.Context) (*orgv1.UserList, error) {
	list, err := t.usermgr.Members().Preload("User").List(ctx,
		users.Filter.WithTeam(t.team),
	)
	if err != nil {
		log.WithError(err).Error("trying to retrieve list of memberships")

		return nil, err
	}

	return DefaultConvertor.FromMembersToUserList(list), nil
}

// ListInvitations returns a list of invitations for the team
func (t tmsImpl) ListInvitations(ctx context.Context) (*orgv1.TeamInvitationList, error) {
	list, err := t.usermgr.Invitations().List(ctx,
		users.Filter.WithTeam(t.team),
	)
	if err != nil {
		log.WithError(err).Error("trying to retrieve the users invitations for team")

		return nil, err
	}

	return DefaultConvertor.FromUserInvitationModelList(list), nil
}

func (t tmsImpl) UpdateTeam(ctx context.Context) error {
	tm := &orgv1.Team{}
	if err := t.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(HubNamespace),
		store.GetOptions.InTo(tm),
		store.GetOptions.WithName(t.team),
	); err != nil {
		log.WithError(err).Error("trying to retrieve the team from api")

		return err
	}
	tm.Annotations = map[string]string{
		Label("changed"): fmt.Sprintf("%d", time.Now().Unix()),
	}

	if err := t.Store().Client().Update(ctx,
		store.UpdateOptions.To(tm),
		store.UpdateOptions.WithForce(true),
	); err != nil {
		log.WithError(err).Error("trying to update the team from api")

		return err
	}

	return nil
}
