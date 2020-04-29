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

package persistence

import (
	"context"
	"testing"

	"github.com/appvia/kore/pkg/persistence/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListUserTeamsWithoutPreloadOK(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	v, err := store.Members().List(context.Background(), List.WithUser("rohith"))

	require.NoError(t, err)
	require.NotNil(t, v)
	require.Equal(t, 2, len(v))
	assert.NotNil(t, v[0].Team)
	assert.NotNil(t, v[1].Team)
}

func TestListUserTeamsWithNoTeams(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	v, err := store.Members().List(context.Background(), List.WithTeam("no_teams"))
	require.NoError(t, err)
	require.NotNil(t, v)
	require.Equal(t, 0, len(v))
}

func TestListTeamsInvalidUser(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	v, err := store.Members().List(context.Background(), List.WithUser("does_not_exist"))
	require.NoError(t, err)
	require.Empty(t, v)
}

func TestListMembersOK(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	v, err := store.Members().Preload("User").List(context.Background(), List.WithTeam("frontend"))
	require.NoError(t, err)
	require.NotNil(t, v)
	require.NotEmpty(t, v)
	require.Equal(t, 2, len(v))

	assert.Equal(t, "rohith", v[0].User.Username)
	assert.Equal(t, "test", v[1].User.Username)
}

func TestListMembersWithNoMembers(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	v, err := store.Members().List(context.Background(), List.WithTeam("no_members"))
	require.NoError(t, err)
	require.NotNil(t, v)
	require.Equal(t, 0, len(v))
}

func TestListMembersWithInvalidTeam(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	v, err := store.Members().List(context.Background(), List.WithTeam("not_there"))
	require.NoError(t, err)
	require.Empty(t, v)
}

func TestUpdateMember(t *testing.T) {
	store := makeTestStore(t)
	ctx := context.Background()
	defer store.Stop()

	name := "member_test"

	err := store.Users().Update(ctx, &model.User{Username: name})
	require.NoError(t, err)

	u, err := store.Users().Get(ctx, name)
	require.NoError(t, err)
	require.NotNil(t, u)

	tm, err := store.Teams().Get(ctx, "devs")
	require.NoError(t, err)
	require.NotNil(t, tm)

	teams, err := store.Members().List(ctx, List.WithUser(name))
	require.NoError(t, err)
	require.NotNil(t, teams)
	require.Empty(t, teams)

	err = store.Members().Update(ctx, &model.Member{UserID: u.ID, TeamID: tm.ID})
	require.NoError(t, err)

	teams, err = store.Members().List(ctx, List.WithUser(name))
	require.NoError(t, err)
	require.NotNil(t, teams)
	assert.Equal(t, 1, len(teams))

	u, err = store.Users().Delete(ctx, &model.User{Username: name})
	require.NoError(t, err)
	require.NotNil(t, u)
}

func TestMembersRoles(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	user := &model.User{Username: "member_role_test", Email: "member_role_test"}
	require.NoError(t, store.Users().Update(context.Background(), user))
	team := &model.Team{Name: "memeber_role_test", Description: "member_role_test"}
	require.NoError(t, store.Teams().Update(context.Background(), team))
	member := &model.Member{UserID: user.ID, TeamID: team.ID, Roles: []string{"admin", "cluster-admin"}}
	require.NoError(t, store.Members().Update(context.Background(), member))

	m, err := store.Members().List(context.Background(), Filter.WithUser("member_role_test"))
	require.NoError(t, err)
	require.NotNil(t, m)
	require.Equal(t, 1, len(m))
	assert.Equal(t, []string{"admin", "cluster-admin"}, m[0].Roles)
}

func TestMembersRoleNone(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	m, err := store.Members().List(context.Background(),
		Filter.WithUser("rohith"),
		Filter.WithTeam("devs"),
	)
	require.NoError(t, err)
	require.NotNil(t, m)
	require.Equal(t, 1, len(m))
	assert.Equal(t, []string{}, m[0].Roles)
}
