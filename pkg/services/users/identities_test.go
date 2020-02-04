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

func TestIdentities(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	assert.NotNil(t, store.Identities())
}

func TestIdentitiesList(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	v, err := store.Identities().List(context.Background())
	require.NoError(t, err)
	require.NotNil(t, v)
	require.NotEmpty(t, v)
}

func BenchmarkListByUsername(b *testing.B) {
	store := makeTestStore(b)
	defer store.Stop()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		store.Identities().List(context.Background(), List.WithUser("rohith"))
	}
}

func TestIdentitiesListByUsername(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	v, err := store.Identities().List(context.Background(), List.WithUser("rohith"))
	require.NoError(t, err)
	require.NotNil(t, v)
	require.Equal(t, 2, len(v))
}

func TestIdentitiesListByProviderName(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	v, err := store.Identities().List(context.Background(),
		List.WithUser("rohith"),
		List.WithProvider("api_token"),
	)
	require.NoError(t, err)
	require.NotNil(t, v)
	require.Equal(t, 1, len(v))
}

func TestIdentitiesUpdateBad(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	i := &model.Identity{UserID: 1}
	err := store.Identities().Update(context.Background(), i)
	assert.Error(t, err)
}

func TestIdentitiesUpdateOKEnsureNoDuplicates(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	count := func() int {
		v, err := store.Identities().List(context.Background(), List.WithProvider("test"))
		require.NoError(t, err)
		require.NotNil(t, v)
		return len(v)
	}
	require.Equal(t, 0, count())

	i := &model.Identity{UserID: 1, Provider: "test"}
	err := store.Identities().Update(context.Background(), i)
	require.NoError(t, err)
	require.Equal(t, 1, count())

	i = &model.Identity{UserID: 1, Provider: "test"}
	err = store.Identities().Update(context.Background(), i)
	require.NoError(t, err)
	require.Equal(t, 1, count())

	require.NoError(t, store.Identities().Delete(context.Background(), i))
}

func TestIdentitiesDelete(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()
	name := "test"

	count := func() int {
		v, err := store.Identities().List(context.Background(), List.WithProvider(name))
		require.NoError(t, err)
		require.NotNil(t, v)
		return len(v)
	}
	require.Equal(t, 0, count())

	i := &model.Identity{UserID: 1, Provider: name}
	require.NoError(t, store.Identities().Update(context.Background(), i))
	require.Equal(t, 1, count())

	require.NoError(t, store.Identities().Delete(context.Background(), i))
	require.Equal(t, 0, count())
}
