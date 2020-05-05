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

	. "github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/persistence/model"

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
