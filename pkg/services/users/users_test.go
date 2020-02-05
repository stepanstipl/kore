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

package users

import (
	"context"
	"testing"

	"github.com/appvia/kore/pkg/services/users/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUsersSize(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	size, err := store.Users().Size(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(5), size)
}

func TestGetUserBad(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	user, err := store.Users().Get(context.TODO(), "not_there")
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestGetUserOK(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	u, err := store.Users().Get(context.TODO(), "rohith")
	require.NoError(t, err)
	require.NotNil(t, u)

	assert.NotEqual(t, 0, u.ID)
	assert.Equal(t, "rohith", u.Username)
	assert.False(t, u.Disabled)
}

func TestUserUpdate(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	user := &model.User{Username: "henry"}
	err := store.Users().Update(context.TODO(), user)
	assert.NoError(t, err)

	u, err := store.Users().Get(context.TODO(), "henry")
	require.NoError(t, err)
	require.NotNil(t, u)
}

func TestUserDelete(t *testing.T) {
	store := makeTestStore(t)
	ctx := context.Background()
	defer store.Stop()

	name := "delete_me"

	user := &model.User{Username: name}
	require.NoError(t, store.Users().Update(ctx, user))

	found, err := makeTestStore(t).Users().Exists(ctx, name)
	require.NoError(t, err)
	require.True(t, found)

	u, err := store.Users().Delete(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, u)

	found, err = makeTestStore(t).Users().Exists(ctx, name)
	require.NoError(t, err)
	require.False(t, found)
}

func TestUsersNoDups(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	name := "test_dups"

	user := &model.User{Username: name, Email: "email@email.com"}
	require.NoError(t, store.Users().Update(context.Background(), user))

	list, err := store.Users().List(context.Background(), Filter.WithUser(name))
	require.NoError(t, err)
	require.NotNil(t, list)
	require.Equal(t, 1, len(list))

	user = &model.User{Username: name, Email: "email@email.com"}
	require.NoError(t, store.Users().Update(context.Background(), user))
	user = &model.User{Username: name, Email: "email@email.com"}
	require.NoError(t, store.Users().Update(context.Background(), user))

	list, err = store.Users().List(context.Background(), Filter.WithUser(name))
	require.NoError(t, err)
	require.NotNil(t, list)
	require.Equal(t, 1, len(list))
}

func TestUserExists(t *testing.T) {
	found, err := makeTestStore(t).Users().Exists(context.Background(), "not_there")
	assert.NoError(t, err)
	assert.False(t, found)
}

func TestUserList(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	users, err := store.Users().List(context.Background())
	require.NoError(t, err)
	require.NotNil(t, users)
	require.NotEmpty(t, users)
	assert.Equal(t, "rohith", users[0].Username)
}
