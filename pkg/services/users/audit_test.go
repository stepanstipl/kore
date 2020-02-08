/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuditRecord(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	store.Audit().Record(context.Background(),
		User("test"),
		Team("test"),
		Type(AuditUpdate),
	).Event("test message")
}

func TestAuditFind(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	store.Audit().Record(context.Background(),
		User("no_one"),
		Team("no_one"),
		Type(AuditUpdate),
	).Event("test message")

	list, err := store.Audit().Find(context.Background(),
		Filter.WithUser("no_one"),
		Filter.WithTeam("no_one"),
	).Do()
	require.NoError(t, err)
	require.NotNil(t, list)
	assert.Equal(t, 1, len(list))
}
