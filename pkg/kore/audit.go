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

package kore

import (
	"context"
	"time"

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/services/users"
	log "github.com/sirupsen/logrus"
)

// Audit represents the interface to the top-level Kore Audit service.
type Audit interface {
	// AuditEvents returns a stream of events across teams since x
	AuditEvents(context.Context, time.Duration) (*orgv1.AuditEventList, error)
	// AuditEventsTeam returns a stream of events in relation to a specific team
	AuditEventsTeam(ctx context.Context, team string, since time.Duration) (*orgv1.AuditEventList, error)
	// Record stores an event in the underlying audit service.
	Record(ctx context.Context, fields ...users.AuditFunc) users.Log
}

type auditImpl struct {
	*hubImpl
}

func (a *auditImpl) Record(ctx context.Context, fields ...users.AuditFunc) users.Log {
	return a.usermgr.Audit().Record(ctx, fields[:]...)
}

// AuditEvents returns all events since the specified duration before the current time
func (a *auditImpl) AuditEvents(ctx context.Context, since time.Duration) (*orgv1.AuditEventList, error) {
	// @step: must be a admin user
	user := authentication.MustGetIdentity(ctx)
	if !user.IsGlobalAdmin() {
		log.WithField(
			"username", user.Username(),
		).Warn("user trying to access the audit logs")

		return nil, NewErrNotAllowed("Must be global admin")
	}

	// @step: retrieve a list of audit events across all teams
	list, err := a.usermgr.Audit().Find(ctx,
		users.Filter.WithDuration(since),
	).Do()
	if err != nil {
		log.WithError(err).Error("trying to retrieve audit logs for teams")

		return nil, err
	}

	return DefaultConvertor.FromAuditModelList(list), nil
}

// AuditEventsTeam returns events in relation to a specific team since the specified duration before the current time
func (a *auditImpl) AuditEventsTeam(ctx context.Context, team string, since time.Duration) (*orgv1.AuditEventList, error) {
	// @step: Check user in the team requested or a global admin
	user := authentication.MustGetIdentity(ctx)
	userInTeam := false
	for _, t := range user.Teams() {
		if t == team {
			userInTeam = true
		}
	}
	if !userInTeam && !user.IsGlobalAdmin() {
		log.WithFields(log.Fields{
			"username": user.Username(),
			"team":     team,
		}).Warn("user trying to access the audit logs for team they're not a member of")

		return nil, NewErrNotAllowed("Must be global admin or a team member")
	}

	list, err := a.usermgr.Audit().Find(ctx,
		users.Filter.WithTeam(team),
		users.Filter.WithDuration(since),
	).Do()
	if err != nil {
		log.WithError(err).Error("trying to retrieve audit logs for team")

		return nil, err
	}

	return DefaultConvertor.FromAuditModelList(list), nil
}
