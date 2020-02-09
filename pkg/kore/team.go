/*
 * Copyright (C) 2019 Appvia Ltd <info@appvia.io>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package kore

import (
	"context"
	"time"

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/services/users"
	log "github.com/sirupsen/logrus"
)

// Team is the contract to a team
type Team interface {
	// AuditEvents returns a collection of audit events for this team
	AuditEvents(context.Context, time.Duration) (*orgv1.AuditEventList, error)
	// Allocations returns the team allocation interface
	Allocations() Allocations
	// Cloud returns the cloud providers
	Cloud() Cloud
	// Clusters returns the teams clusters
	Clusters() Clusters
	// Members returns the team members interface
	Members() TeamMembers
	// NamespaceClaims returns the the interface
	NamespaceClaims() NamespaceClaims
}

// tmImpl is a team interface
type tmImpl struct {
	*hubImpl
	// team is the name of the team
	team string
}

// Allocations return an interface to the team allocations
func (t *tmImpl) Allocations() Allocations {
	return &acaImpl{
		hubImpl: t.hubImpl,
		team:    t.team,
	}
}

func (t *tmImpl) Cloud() Cloud {
	return &cloudImpl{
		hubImpl: t.hubImpl,
		team:    t.team,
	}
}

func (t *tmImpl) Clusters() Clusters {
	return &clsImpl{
		hubImpl: t.hubImpl,
		team:    t.team,
	}
}

// Members returns the team members interface
func (t *tmImpl) Members() TeamMembers {
	return &tmsImpl{
		hubImpl: t.hubImpl,
		team:    t.team,
	}
}

// NamespaceClaims returns a namespace claim interface
func (t *tmImpl) NamespaceClaims() NamespaceClaims {
	return &nsImpl{
		hubImpl: t.hubImpl,
		team:    t.team,
	}
}

// AuditEvents returns a stream of events in relation to the teams since x
func (t *tmImpl) AuditEvents(ctx context.Context, since time.Duration) (*orgv1.AuditEventList, error) {
	list, err := t.Audit().Find(ctx,
		users.Filter.WithTeam(t.team),
		users.Filter.WithDuration(since),
	).Do()
	if err != nil {
		log.WithError(err).Error("trying to retrieve audit logs for team")

		return nil, err
	}

	return DefaultConvertor.FromAuditModelList(list), nil
}
