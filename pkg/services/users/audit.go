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

import "context"

// Find is used to return results from the log
func (a *storeImpl) Find(ctx context.Context, filters ...ListFunc) Find {
	return newQuery(ctx, a.dbc, filters...)
}

// Record is responsible for adding an entry into the log
func (a *storeImpl) Record(ctx context.Context, fields ...AuditFunc) Log {
	return newEntry(ctx, a.dbc, fields...)
}
