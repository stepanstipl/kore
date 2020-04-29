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
	"time"

	"github.com/appvia/kore/pkg/persistence/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListInvitationsOK(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	v, err := store.Invitations().List(context.Background())
	require.NoError(t, err)
	require.NotNil(t, v)
	assert.Empty(t, v)
}

func TestUpdateInvitationOK(t *testing.T) {
	store := makeTestStore(t)
	ctx := context.Background()
	defer store.Stop()

	team, err := store.Teams().Get(ctx, "devs")
	require.NoError(t, err)
	require.NotNil(t, team)

	name := "iv_delete_me"

	err = store.Users().Update(ctx, &model.User{Username: name})
	require.NoError(t, err)

	user, err := store.Users().Get(ctx, name)
	require.NoError(t, err)
	require.NotNil(t, user)

	err = store.Invitations().Update(ctx, &model.Invitation{
		User:    user,
		Team:    team,
		Expires: time.Now().Add(1 * time.Hour),
	})
	require.NoError(t, err)
}

func TestUpdateInvitationBad(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()
	ctx := context.Background()

	err := store.Invitations().Update(ctx, &model.Invitation{
		UserID:  99999,
		Expires: time.Now().Add(1 * time.Hour),
	})
	require.Error(t, err)
}

func TestInvitationGoesWithUser(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()
	name := "iv_delete_me"

	v, err := store.Invitations().List(context.Background(), List.WithUser(name))
	require.NoError(t, err)
	require.NotNil(t, v)
	require.Equal(t, 1, len(v))

	_, err = store.Users().Delete(context.Background(), &model.User{Username: name})
	require.NoError(t, err)

	v, err = store.Invitations().List(context.Background(), List.WithUser(name))
	require.NoError(t, err)
	require.NotNil(t, v)
	require.Equal(t, 0, len(v))
}

func TestDeleteInvitation(t *testing.T) {
	store := makeTestStore(t)
	ctx := context.Background()
	defer store.Stop()

	team, err := store.Teams().Get(ctx, "devs")
	require.NoError(t, err)
	require.NotNil(t, team)

	name := "iv_delete_me"

	err = store.Users().Update(ctx, &model.User{Username: name})
	require.NoError(t, err)

	user, err := store.Users().Get(ctx, name)
	require.NoError(t, err)
	require.NotNil(t, user)

	err = store.Invitations().Update(ctx, &model.Invitation{
		User:    user,
		Team:    team,
		Expires: time.Now().Add(1 * time.Hour),
	})
	require.NoError(t, err)

	v, err := store.Invitations().List(context.Background(), List.WithUser(name))
	require.NoError(t, err)
	require.NotNil(t, v)
	require.Equal(t, 1, len(v))

	err = store.Invitations().Delete(ctx, &model.Invitation{
		TeamID: v[0].TeamID,
		UserID: v[0].UserID,
	})
	require.NoError(t, err)

	v, err = store.Invitations().List(context.Background(), List.WithUser(name))
	require.NoError(t, err)
	require.NotNil(t, v)
	require.Equal(t, 0, len(v))

	_, err = store.Users().Delete(ctx, &model.User{Username: name})
	require.NoError(t, err)
}
