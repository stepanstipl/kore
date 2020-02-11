/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package kore

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/services/users"
	"github.com/appvia/kore/pkg/store"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	log "github.com/sirupsen/logrus"
)

// Teams is the kore api teams interface
type Teams interface {
	// AuditEvents returns a stream of events in relation to the teams since x
	AuditEvents(context.Context, time.Duration) (*orgv1.AuditEventList, error)
	// Delete removes the team from the kore
	Delete(context.Context, string) error
	// Exists checks if the team exists
	Exists(context.Context, string) (bool, error)
	// Get returns the team from the kore
	Get(context.Context, string) (*orgv1.Team, error)
	// List returns a list of teams
	List(context.Context) (*orgv1.TeamList, error)
	// Team returns a team interface
	Team(string) Team
	// Update is responsible for creating / updating a team
	Update(context.Context, *orgv1.Team) (*orgv1.Team, error)
}

// teamsImpl provides the teams implementation
type teamsImpl struct {
	*hubImpl
}

// Delete removes the team from the kore
func (t *teamsImpl) Delete(ctx context.Context, name string) error {
	// @step: retrieve the user from context
	user := authentication.MustGetIdentity(ctx)
	if !user.IsGlobalAdmin() {
		log.Warn("non admin user attempting to delete a team")

		return ErrUnauthorized
	}

	// @step: retrieve the team
	team, err := t.usermgr.Teams().Get(ctx, name)
	if err != nil {
		return err
	}
	// we need to check if the team has any resources under it
	// if so we thrown back a error saying these must be deleted
	// else we can go ahead and remove the team

	// @step: check if the team has any resources
	resources, err := t.Team(name).Clusters().List(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"team": name,
		}).Warn("failed checking for team clusters")

		return fmt.Errorf("failed checking for team clusters: %s", err)
	}

	if len(resources.Items) > 0 {
		log.WithFields(log.Fields{
			"team": name,
		}).Warn("attempting to delete a team whom has cluster")

		return ErrNotAllowed{}
	}

	// @step: check if the team has any allocations to other teams
	log.WithFields(log.Fields{
		"team": name,
	}).Info("deleting the team from the kore")

	// @step: delete in the db
	if err := t.usermgr.Teams().Delete(ctx, team); err != nil {
		log.WithError(err).Error("trying to delete the team in kore")

		return err
	}

	tm := &orgv1.Team{
		ObjectMeta: metav1.ObjectMeta{
			Name:      team.Name,
			Namespace: HubNamespace,
		},
	}

	return t.Store().Client().Delete(ctx, store.DeleteOptions.From(tm))
}

// Get returns the team from the kore
func (t *teamsImpl) Get(ctx context.Context, name string) (*orgv1.Team, error) {
	model, err := t.usermgr.Teams().Get(ctx, name)
	if err != nil {
		log.WithField("name", name).WithError(err).Error("trying to retrieve the user")

		return nil, err
	}

	return DefaultConvertor.FromTeamModel(model), err
}

// List returns a list of teams
func (t *teamsImpl) List(ctx context.Context) (*orgv1.TeamList, error) {
	model, err := t.usermgr.Teams().List(ctx)
	if err != nil {
		log.WithError(err).Error("trying to retrieve list of team in kore")

		return nil, err
	}

	return DefaultConvertor.FromTeamsModelList(model), nil
}

// Exists checks if the team exists in the kore
func (t *teamsImpl) Exists(ctx context.Context, name string) (bool, error) {
	// @step: we check the user management service for teams
	return t.usermgr.Teams().Exists(ctx, name)
}

// Update is responsible for updating the team
func (t *teamsImpl) Update(ctx context.Context, team *orgv1.Team) (*orgv1.Team, error) {
	team.Namespace = HubNamespace
	team.Annotations = map[string]string{
		Label("changed"): fmt.Sprintf("%d", time.Now().Unix()),
	}

	// @step: check the team is ok
	if team.Name == "" {
		return nil, errors.New("no name for team defined")
	}

	// @step: retrieve the user from context
	user := authentication.MustGetIdentity(ctx)
	if !user.IsGlobalAdmin() {
		log.Warn("non admin user attempting to create a team")

		return nil, ErrUnauthorized
	}

	if !user.IsGlobalAdmin() && !t.IsValidTeamName(team.Name) {
		return nil, ErrNotAllowed{message: "name: " + team.Name + " cannot be used in the kore"}
	}

	logger := log.WithFields(log.Fields{
		"team": team.Name,
		"user": user.Username,
	})
	logger.Info("attempting to update or create team in kore")

	model := DefaultConvertor.ToTeamModel(team)

	// @step: update or create in the kore
	if err := t.usermgr.Teams().Update(ctx, model); err != nil {
		log.WithError(err).Error("trying to update a team in user management")

		return nil, err
	}

	// @step: add the user whom created is a admin member
	if err := t.usermgr.Members().AddUser(ctx,
		user.Username(), team.Name, []string{"members", "admin"}); err != nil {

		logger.WithError(err).Error("trying to add the user a admin user on the team")

		return nil, err
	}

	// @step: check if the team is in the api
	if found, err := t.Store().Client().Has(ctx, store.HasOptions.From(team)); err != nil {
		log.WithError(err).Error("trying to check for team in the api")

		return nil, err
	} else if !found {
		if err := t.Store().Client().Update(ctx,
			store.UpdateOptions.To(team),
			store.UpdateOptions.WithCreate(true),
			store.UpdateOptions.WithForce(true),
		); err != nil {
			log.WithError(err).Error("trying to create / update a team in kore")

			return nil, err
		}
	}

	// @step: add the entry into the audit
	t.Audit().Record(ctx,
		users.Resource(team.Name),
		users.ResourceUID(string(team.UID)),
		users.Type(users.AuditUpdate),
		users.User(user.Username()),
	).Event("team has been update or created in the kore")

	return DefaultConvertor.FromTeamModel(model), nil
}

// IsValidTeamName checks the team name is ok to use
func (t *teamsImpl) IsValidTeamName(name string) bool {
	if name == "kore-admin" {
		return true
	}
	// @step: ensure the team name does not infridge on policy
	for _, x := range []string{"kube-", "prometheus", "kore", "istio", "grafana", "olm", "default"} {
		if strings.HasPrefix(name, x) {
			return false
		}
	}

	return true
}

// AuditEvents returns a stream of events in relation to the teams since x
func (t *teamsImpl) AuditEvents(ctx context.Context, since time.Duration) (*orgv1.AuditEventList, error) {
	// @step: must be a admin user
	user := authentication.MustGetIdentity(ctx)
	if !user.IsGlobalAdmin() {
		log.WithField(
			"username", user.Username(),
		).Warn("user trying to access the audit logs")

		return nil, ErrUnauthorized
	}

	// @step: retrieve a list of audit events across all teams
	list, err := t.Audit().Find(ctx,
		users.Filter.WithDuration(since),
		users.Filter.WithTeamNotNull(),
	).Do()
	if err != nil {
		log.WithError(err).Error("trying to retrieve audit logs for teams")

		return nil, err
	}

	return DefaultConvertor.FromAuditModelList(list), nil
}

// Team returns the team interface
func (t *teamsImpl) Team(team string) Team {
	return &tmImpl{
		hubImpl: t.hubImpl,
		team:    team,
	}
}
