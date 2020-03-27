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
	"time"

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Teams is the kore api teams interface
type Teams interface {
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

	// @step: one has to be a member of the team to delete it
	if !user.IsMember(name) && !user.IsGlobalAdmin() {
		log.WithFields(log.Fields{
			"team": name,
			"user": user.Username(),
		}).Warn("user attempting to delete a delete not a member of")

		return ErrUnauthorized
	}

	// @step: ensure the admin team can never be delete - oh my gosh
	if name == HubAdminTeam {
		return ErrNotAllowed{message: HubAdminTeam + " team cannot be deleted"}
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

		return ErrNotAllowed{message: "all team resources must be deleted before team can be removed"}
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
	// @step: teams are slightly different as they are being placed in both a database
	// and a kubernetes. In order to ensure both instances comply with the naming conver
	if err := IsValidResourceName("team", team.Name); err != nil {
		return nil, err
	}

	team.Namespace = HubNamespace

	// @step: retrieve the user from context
	user := authentication.MustGetIdentity(ctx)

	logger := log.WithFields(log.Fields{
		"team":  team.Name,
		"teams": strings.Join(user.Teams(), ","),
		"user":  user.Username(),
	})
	logger.Info("attempting to update or create team in kore")

	// @logic
	// - bypass the check if the user is a admin
	// - check if the team already exists
	// - if they exist then only a member of that team can update it
	// - if the team does not exist then any user can claim the team and added as member
	// - ensure the namespace name is valid

	found, err := t.usermgr.Teams().Exists(ctx, team.Name)
	if err != nil {
		logger.WithError(err).Error("trying to check if the team exists")

		return nil, err
	}

	if !user.IsGlobalAdmin() {
		if found {
			// ensure the user is a member of the team
			if !utils.Contains(team.Name, user.Teams()) {
				logger.Warn("trying to update a team they do not belong to")

				return nil, ErrUnauthorized
			}
		}

		if !t.IsValidTeamName(team.Name) {
			return nil, ErrNotAllowed{message: "name: " + team.Name + " cannot be used in the kore"}
		}
	}

	// @step: ensure the default and apply a timestamp for triggers
	team.Annotations = map[string]string{
		Label("changed"): fmt.Sprintf("%d", time.Now().Unix()),
	}

	// @step: convert the team to a team model
	model := DefaultConvertor.ToTeamModel(team)

	// @step: update the team in the users store
	if err := t.usermgr.Teams().Update(ctx, model); err != nil {
		log.WithError(err).Error("trying to update a team in user management")

		return nil, err
	}

	// @step: add the user whom created is a admin member if the team never existed
	if !found {
		roles := []string{"admin", "member"}

		logger.Info("adding the user as the admin of team")

		if err := t.usermgr.Members().AddUser(ctx, user.Username(), team.Name, roles); err != nil {
			logger.WithError(err).Error("trying to add the user a admin user on the team")

			return nil, err
		}
	}

	// @step: check if the team is in the kube api
	found, err = t.Store().Client().Has(ctx, store.HasOptions.From(team))
	if err != nil {
		log.WithError(err).Error("trying to check for team in the api")

		return nil, err
	}
	if !found {
		if err := t.Store().Client().Update(ctx,
			store.UpdateOptions.To(team),
			store.UpdateOptions.WithCreate(true),
			store.UpdateOptions.WithForce(true),
		); err != nil {
			log.WithError(err).Error("trying to create / update a team in kore")

			return nil, err
		}
	}

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

// Team returns the team interface
func (t *teamsImpl) Team(team string) Team {
	return &tmImpl{
		hubImpl: t.hubImpl,
		team:    team,
	}
}
