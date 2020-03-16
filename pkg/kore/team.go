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
	// KubernetesCredentials returns the k8s credentials
	KubernetesCredentials() KubernetesCredentials
	// Members returns the team members interface
	Members() TeamMembers
	// NamespaceClaims returns the the interface
	NamespaceClaims() NamespaceClaims
	// Secrets returns the secret interface
	Secrets() Secrets
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

// Secrets returns a secrets interface
func (t *tmImpl) Secrets() Secrets {
	return &secretImpl{
		hubImpl: t.hubImpl,
		team:    t.team,
	}
}

// KubernetesCredentials returns the k8s credentials
func (t *tmImpl) KubernetesCredentials() KubernetesCredentials {
	return &kcImpl{
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
