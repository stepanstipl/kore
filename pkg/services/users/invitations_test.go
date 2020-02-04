/**
	"github.com/appvia/kore/pkg/services/users/model"
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

package users

import (
	"context"
	"testing"
	"time"

	"github.com/appvia/kore/pkg/services/users/model"

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
