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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/appvia/kore/pkg/utils/kubernetes"
	"github.com/appvia/kore/pkg/utils/validation"

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils"

	log "github.com/sirupsen/logrus"
)

// Teams is the kore api teams interface
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Teams
type Teams interface {
	// CheckDelete verifies whether the team can be deleted
	CheckDelete(context.Context, *orgv1.Team, ...DeleteOptionFunc) error
	// Delete removes the team from the kore
	Delete(context.Context, string, ...DeleteOptionFunc) error
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
	// GenerateTeamIdentifier creates and records a new team identity
	GenerateTeamIdentifier(ctx context.Context, teamName string) (string, error)
}

// teamsImpl provides the teams implementation
type teamsImpl struct {
	*hubImpl
}

// CheckDelete verifies whether the team can be deleted
func (t *teamsImpl) CheckDelete(ctx context.Context, team *orgv1.Team, o ...DeleteOptionFunc) error {
	opts := ResolveDeleteOptions(o)

	// @step: ensure the admin team can never be delete - oh my gosh
	if team.Name == HubAdminTeam || team.Name == HubDefaultTeam {
		return ErrNotAllowed{message: team.Name + " team cannot be deleted"}
	}

	if !opts.Cascade {
		var dependents []kubernetes.DependentReference

		clusters, err := t.Team(team.Name).Clusters().List(ctx)
		if err != nil {
			return fmt.Errorf("failed to list services: %w", err)
		}
		for _, item := range clusters.Items {
			dependents = append(dependents, kubernetes.DependentReferenceFromObject(&item))
		}

		if len(dependents) > 0 {
			return validation.ErrDependencyViolation{
				Message:    "the following objects need to be deleted first",
				Dependents: dependents,
			}
		}
	}

	return nil
}

// Delete removes the team from the kore
func (t *teamsImpl) Delete(ctx context.Context, name string, o ...DeleteOptionFunc) error {
	opts := ResolveDeleteOptions(o)

	// @step: retrieve the team
	teamRecord, err := t.persistenceMgr.Teams().Get(ctx, name)
	if err != nil {
		return err
	}

	team := &orgv1.Team{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: HubNamespace,
		},
	}

	if err := opts.Check(team, func(o ...DeleteOptionFunc) error { return t.CheckDelete(ctx, team, o...) }); err != nil {
		return err
	}

	if opts.Cascade {
		clusters, err := t.Team(name).Clusters().List(ctx)
		if err != nil {
			return fmt.Errorf("failed to list clusters: %w", err)
		}

		for _, cluster := range clusters.Items {
			if _, err := t.Team(name).Clusters().Delete(ctx, cluster.Name, DeleteOptionCascade(true)); err != nil {
				return err
			}
		}
	}

	log.WithFields(log.Fields{
		"team": name,
	}).Info("deleting the team from the kore")

	// @step: delete in the db
	if err := t.persistenceMgr.Teams().Delete(ctx, teamRecord); err != nil {
		log.WithError(err).Error("trying to delete the team in kore")

		return err
	}

	if teamRecord.Identifier != "" {
		if err := t.persistenceMgr.TeamAssets().MarkTeamIdentityDeleted(ctx, teamRecord.Identifier); err != nil {
			log.WithError(err).Error("trying to mark team identity as deleted in kore")
			return fmt.Errorf("error marking team identity as deleted: %w", err)
		}
	}

	return t.Store().Client().Delete(ctx, store.DeleteOptions.From(team))
}

// Get returns the team from the kore
func (t *teamsImpl) Get(ctx context.Context, name string) (*orgv1.Team, error) {
	model, err := t.persistenceMgr.Teams().Get(ctx, name)
	if err != nil {
		log.WithField("name", name).WithError(err).Error("trying to retrieve the team")

		return nil, err
	}

	return DefaultConvertor.FromTeamModel(model), err
}

// List returns a list of teams
func (t *teamsImpl) List(ctx context.Context) (*orgv1.TeamList, error) {
	model, err := t.persistenceMgr.Teams().List(ctx)
	if err != nil {
		log.WithError(err).Error("trying to retrieve list of team in kore")

		return nil, err
	}

	return DefaultConvertor.FromTeamsModelList(model), nil
}

// Exists checks if the team exists in kore
func (t *teamsImpl) Exists(ctx context.Context, name string) (bool, error) {
	// @step: we check the user management service for teams
	return t.persistenceMgr.Teams().Exists(ctx, name)
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

	found := true
	existingTeam, err := t.persistenceMgr.Teams().Get(ctx, team.Name)
	if err != nil {
		if !t.persistenceMgr.IsNotFound(err) {
			logger.WithError(err).Error("trying to check if the team exists")

			return nil, err
		}
		found = false
	}

	if !user.IsGlobalAdmin() {
		if found {
			// ensure the user is a member of the team
			if !utils.Contains(team.Name, user.Teams()) {
				logger.Warn("trying to update a team they do not belong to")

				return nil, ErrUnauthorized
			}
		}
	}

	if !IsValidTeamName(team.Name) {
		return nil, ErrNotAllowed{message: "name: " + team.Name + " cannot be used in Kore"}
	}

	// @step: ensure the default and apply a timestamp for triggers
	team.Annotations = map[string]string{
		Label("changed"): fmt.Sprintf("%d", time.Now().Unix()),
	}

	// @step: ensure team has a globally-unique identifier and assign one if not
	if found && team.Labels[LabelTeamIdentifier] == "" && existingTeam.Identifier != "" {
		// Populate with current identifier if empty identifier specified.
		if team.Labels == nil {
			team.Labels = map[string]string{}
		}
		team.Labels[LabelTeamIdentifier] = existingTeam.Identifier
	}
	if found && team.Labels[LabelTeamIdentifier] != existingTeam.Identifier {
		return nil, ErrNotAllowed{message: "Identifier is assigned by Kore and cannot be changed"}
	}
	// Still empty? Need to generate and assign an identifier to this team
	if team.Labels[LabelTeamIdentifier] == "" {
		if team.Labels == nil {
			team.Labels = map[string]string{}
		}
		team.Labels[LabelTeamIdentifier], err = t.GenerateTeamIdentifier(ctx, team.Name)
		if err != nil {
			logger.WithError(err).Error("trying to assign team identifier")

			return nil, err
		}
	}

	// @step: convert the team to a team model
	model := DefaultConvertor.ToTeamModel(team)

	// @step: update the team in the users store
	if err := t.persistenceMgr.Teams().Update(ctx, model); err != nil {
		log.WithError(err).Error("trying to update a team in user management")

		return nil, err
	}

	// @step: add the user whom created is a admin member if the team never existed
	if !found {
		roles := []string{"admin", "member"}

		logger.Info("adding the user as the admin of team")

		if err := t.persistenceMgr.Members().AddUser(ctx, user.Username(), team.Name, roles); err != nil {
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

func (t *teamsImpl) GenerateTeamIdentifier(ctx context.Context, teamName string) (string, error) {
	identifier := utils.GenerateIdentifier()
	err := t.hubImpl.persistenceMgr.TeamAssets().RecordTeamIdentity(ctx, identifier, teamName)
	if err != nil {
		return "", fmt.Errorf("Failed to persist identifier for team %s: %w", teamName, err)
	}
	return identifier, nil
}

// Team returns the team interface
func (t *teamsImpl) Team(team string) Team {
	return &tmImpl{
		hubImpl: t.hubImpl,
		team:    team,
	}
}

// IsValidTeamName checks the team name is ok to use
func IsValidTeamName(name string) bool {
	if name == "kore-admin" || name == "kore-default" {
		return true
	}

	return !IsNamespaceNameProtected(name)
}
