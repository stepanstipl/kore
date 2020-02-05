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

package fake

import (
	"context"

	"github.com/appvia/kore/pkg/services/audit"
	"github.com/appvia/kore/pkg/services/audit/model"
)

type fakeImpl struct{}

func NewFakeAudit(config audit.Config) audit.Interface {
	return &fakeImpl{}
}

// Find is used to retrieve records from the log
func (f *fakeImpl) Find(context.Context, ...audit.ListFunc) audit.Find {
	return f
}

// Record records an event in the audit log
func (f *fakeImpl) Record(context.Context, ...audit.AuditFunc) audit.Log {
	return f
}

// Event records the entry into the audit log
func (f *fakeImpl) Event(string) error {
	return nil
}

func (f *fakeImpl) Do() ([]*model.AuditEvent, error) {
	return []*model.AuditEvent{}, nil
}

// Stop stops the service and releases resources
func (f *fakeImpl) Stop() error {
	return nil
}
