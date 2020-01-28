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

package users

import (
	"context"
	"testing"
	
	"github.com/appvia/kore/pkg/services/users/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTeamBad(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	_, err := store.Teams().Get(context.Background(), "not_there")
	assert.Error(t, err)
}

func TestGetTeamsOK(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	v, err := store.Teams().Get(context.Background(), "devs")
	assert.NoError(t, err)
	assert.NotNil(t, v)
}

func TestTeamList(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	v, err := store.Teams().List(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.NotEmpty(t, v)
	assert.Equal(t, "devs", v[0].Name)
}

func TestTeamCreateOK(t *testing.T) {
	store := makeTestStore(t)
	ctx := context.Background()
	defer store.Stop()

	v := &model.Team{Name: "na"}
	require.NoError(t, store.Teams().Update(ctx, v))
	require.NoError(t, store.Teams().Update(ctx, v))

	team, err := store.Teams().Get(ctx, "na")
	require.NoError(t, err)
	require.NotNil(t, team)

	require.NoError(t, store.Teams().Delete(ctx, team))
}

func TestTeamsDelete(t *testing.T) {
	store := makeTestStore(t)
	ctx := context.Background()
	defer store.Stop()

	name := "delete_me"

	v := &model.Team{Name: name}
	require.NoError(t, store.Teams().Update(ctx, v))

	found, err := makeTestStore(t).Teams().Exists(ctx, name)
	require.NoError(t, err)
	require.True(t, found)

	require.NoError(t, store.Teams().Delete(ctx, v))

	found, err = makeTestStore(t).Teams().Exists(ctx, name)
	require.NoError(t, err)
	require.False(t, found)
}