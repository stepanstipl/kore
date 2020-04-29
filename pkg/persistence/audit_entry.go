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

package persistence

import (
	"context"

	"github.com/appvia/kore/pkg/persistence/model"

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

	if e.filter.HasVerb() {
		q = q.Where("q.type = ?", e.filter.GetVerb())
	}
	if e.filter.HasTeam() {
		q = q.Where("q.team = ?", e.filter.GetTeam())
	}
	if e.filter.HasTeams() {
		q = q.Where("q.team IN (?)", e.filter.GetTeams())
	}
	if e.filter.HasTeamsNotNull() {
		q = q.Where("q.team != \"\"")
	}
	if e.filter.HasUser() {
		q = q.Where("q.user = ?", e.filter.GetUser())
	}
	if e.filter.HasDuration() {
		q = q.Where("timestampdiff(minute, q.created_at, NOW()) < ?", int(e.filter.GetDuration().Minutes()))
	}

	list := []*model.AuditEvent{}

	return list, q.Find(&list).Error
}

// Event records the entry into the audit log
func (e *entryImpl) Event(message string) {
	if e.event.Verb == "" || message == "" {
		return
	}
	e.event.Message = message

	e.db.Model(&model.AuditEvent{}).Save(e.event)
}
