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

package persistence_test

import (
	"context"
	"testing"

	"github.com/appvia/kore/pkg/persistence/model"

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
	assert.Equal(t, "All", v[0].Name)
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
