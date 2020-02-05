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

package audit

import (
	"context"
	"errors"

	"github.com/appvia/kore/pkg/services/audit/model"

	"github.com/jinzhu/gorm"
)

// entryImpl implements the recoreentry interface
type entryImpl struct {
	filter *ListOptions
	// ctx is the context
	ctx context.Context
	// event is the event being added
	event *model.AuditEvent
	// db is the database connection
	db *gorm.DB
}

// newQuery creates and returns a query
func newQuery(ctx context.Context, db *gorm.DB, opts ...ListFunc) *entryImpl {
	return &entryImpl{
		ctx:    ctx,
		db:     db,
		filter: ApplyListOptions(opts...),
	}
}

// newEntry creates a new entry for use
func newEntry(ctx context.Context, db *gorm.DB, fields ...AuditFunc) *entryImpl {
	m := &model.AuditEvent{}
	for _, method := range fields {
		method(m)
	}

	return &entryImpl{ctx: ctx, db: db, event: m}
}

// Find is used to return a search from the audit log
func (e *entryImpl) Do() ([]*model.AuditEvent, error) {
	// @step: construct the query
	q := e.db.
		Model(&model.AuditEvent{}).
		Select("q.*").
		Table("audit_events q")

	if e.filter.HasType() {
		q = q.Where("q.type = ?", e.filter.GetType())
	}
	if e.filter.HasTeam() {
		q = q.Where("q.team = ?", e.filter.GetTeam())
	}
	if e.filter.HasTeams() {
		q = q.Where("q.team IN (?)", e.filter.GetTeams())
	}
	if e.filter.HasUser() {
		q = q.Where("q.user = ?", e.filter.GetUser())
	}
	list := []*model.AuditEvent{}

	return list, q.Find(&list).Error
}

// Event records the entry into the audit log
func (e *entryImpl) Event(message string) error {
	if e.event.Type == "" {
		return errors.New("no event type defined")
	}
	if message == "" {
		return errors.New("no message defined")
	}
	e.event.Message = message

	return e.db.Model(&model.AuditEvent{}).Save(e.event).Error
}
