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

func TestGetConfigsBad(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	_, err := store.Configs().Get(context.Background(), "not_there")
	assert.Error(t, err)
}

func TestGetConfigsOK(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	v, err := store.Configs().Get(context.Background(), "example1")
	assert.NoError(t, err)
	assert.NotNil(t, v)
}

func TestListConfig(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	v, err := store.Configs().List(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.NotEmpty(t, v)
	assert.Equal(t, "example1", v[0].Name)
}

func TestCreateConfig(t *testing.T) {
	store := makeTestStore(t)
	ctx := context.Background()
	defer store.Stop()

	v := &model.Config{Name: "na"}
	require.NoError(t, store.Configs().Update(ctx, v))
	require.NoError(t, store.Configs().Update(ctx, v))

	conf, err := store.Configs().Get(ctx, "na")
	require.NoError(t, err)
	require.NotNil(t, conf)

	_, err = store.Configs().Delete(ctx, conf)

	require.NoError(t, err)
}

func TestDeleteConfig(t *testing.T) {
	store := makeTestStore(t)
	ctx := context.Background()
	defer store.Stop()

	name := "delete_me"

	conf := &model.Config{Name: name}
	require.NoError(t, store.Configs().Update(ctx, conf))

	found, err := makeTestStore(t).Configs().Exists(ctx, name)
	require.NoError(t, err)
	require.True(t, found)

	c, err := store.Configs().Delete(ctx, conf)
	require.NoError(t, err)
	require.NotNil(t, c)

	found, err = makeTestStore(t).Configs().Exists(ctx, name)
	require.NoError(t, err)
	require.False(t, found)
}
