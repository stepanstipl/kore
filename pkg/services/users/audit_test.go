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

package users

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuditRecord(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	store.Audit().Record(context.Background(),
		User("test"),
		Team("test"),
		Verb(AuditUpdate),
	).Event("test message")
}

func TestAuditFind(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	store.Audit().Record(context.Background(),
		User("no_one"),
		Team("no_one"),
		Verb(AuditUpdate),
	).Event("test message")

	list, err := store.Audit().Find(context.Background(),
		Filter.WithUser("no_one"),
		Filter.WithTeam("no_one"),
	).Do()
	require.NoError(t, err)
	require.NotNil(t, list)
	assert.Equal(t, 1, len(list))
}

func TestAuditFindByDuration(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	store.Audit().Record(context.Background(),
		User("no_one"),
		Team("no_one"),
		Verb(AuditUpdate),
	).Event("test message")

	list, err := store.Audit().Find(context.Background(),
		Filter.WithUser("no_one"),
		Filter.WithTeam("no_one"),
		Filter.WithDuration(2*time.Minute),
	).Do()
	require.NoError(t, err)
	require.NotNil(t, list)
	assert.Equal(t, 2, len(list))
}
