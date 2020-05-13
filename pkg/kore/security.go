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

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	securityv1 "github.com/appvia/kore/pkg/apis/security/v1"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/security"
	"sigs.k8s.io/controller-runtime/pkg/client"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Security represents the interface to the top-level Kore Security service.
type Security interface {
	ScanPlan(ctx context.Context, client client.Client, plan *configv1.Plan) error
	ScanCluster(ctx context.Context, client client.Client, cluster *clustersv1.Cluster) error
	ListScans(ctx context.Context, latestOnly bool) (*securityv1.SecurityScanResultList, error)
	ScanHistoryForResource(ctx context.Context, typ metav1.TypeMeta, obj metav1.ObjectMeta) (*securityv1.SecurityScanResultList, error)
	GetCurrentScanForResource(ctx context.Context, typ metav1.TypeMeta, obj metav1.ObjectMeta) (*securityv1.SecurityScanResult, error)
	GetScan(ctx context.Context, id uint64) (*securityv1.SecurityScanResult, error)
	ListRules(ctx context.Context) (*securityv1.SecurityRuleList, error)
	GetRule(ctx context.Context, code string) (*securityv1.SecurityRule, error)
	GetOverview(ctx context.Context) (*securityv1.SecurityOverview, error)
	GetTeamOverview(ctx context.Context, team string) (*securityv1.SecurityOverview, error)
	// ArchiveResourceScans sets all scans for a resource to archived, for when that resource is deleted.
	ArchiveResourceScans(ctx context.Context, typ metav1.TypeMeta, obj metav1.ObjectMeta) error
}

var _ Security = &securityImpl{}

type securityImpl struct {
	scanner         security.Scanner
	securityPersist persistence.Security
}

func (s *securityImpl) ScanPlan(ctx context.Context, client client.Client, plan *configv1.Plan) error {
	scanResult := s.scanner.ScanPlan(ctx, client, plan)

	return s.persistScan(ctx, scanResult)
}

func (s *securityImpl) ArchiveResourceScans(ctx context.Context, typ metav1.TypeMeta, obj metav1.ObjectMeta) error {
	gvk := typ.GroupVersionKind()
	return s.securityPersist.ArchiveResourceScans(ctx, gvk.Group, gvk.Version, gvk.Kind, obj.Namespace, obj.Name)
}

func (s *securityImpl) ScanCluster(ctx context.Context, client client.Client, cluster *clustersv1.Cluster) error {
	scanResult := s.scanner.ScanCluster(ctx, client, cluster)

	return s.persistScan(ctx, scanResult)
}

func (s *securityImpl) persistScan(ctx context.Context, scanResult *securityv1.SecurityScanResult) error {
	scanResultDB := DefaultConvertor.ToSecurityScanResult(scanResult)

	if err := s.securityPersist.StoreScan(ctx, &scanResultDB); err != nil {
		log.WithError(err).Error("trying to persist security security scan")

		return err
	}

	return nil
}

func (s *securityImpl) ListScans(ctx context.Context, latestOnly bool) (*securityv1.SecurityScanResultList, error) {
	err := s.userPermitted(ctx, "")
	if err != nil {
		return nil, err
	}

	res, err := s.securityPersist.ListScans(ctx, latestOnly)
	if err != nil {
		return nil, err
	}
	result := &securityv1.SecurityScanResultList{}
	result.Items = make([]securityv1.SecurityScanResult, len(res))
	for i, r := range res {
		result.Items[i] = DefaultConvertor.FromSecurityScanResult(r)
	}
	return result, nil
}

func (s *securityImpl) ScanHistoryForResource(ctx context.Context, typ metav1.TypeMeta, obj metav1.ObjectMeta) (*securityv1.SecurityScanResultList, error) {
	err := s.userPermitted(ctx, obj.Namespace)
	if err != nil {
		return nil, err
	}

	gvk := typ.GroupVersionKind()
	res, err := s.securityPersist.ListResourceScanHistory(ctx, gvk.Group, gvk.Version, gvk.Kind, obj.Namespace, obj.Name)
	if err != nil {
		return nil, err
	}
	result := &securityv1.SecurityScanResultList{}
	result.Items = make([]securityv1.SecurityScanResult, len(res))
	for i, r := range res {
		result.Items[i] = DefaultConvertor.FromSecurityScanResult(r)
	}
	return result, nil
}

func (s *securityImpl) GetCurrentScanForResource(ctx context.Context, typ metav1.TypeMeta, obj metav1.ObjectMeta) (*securityv1.SecurityScanResult, error) {
	gvk := typ.GroupVersionKind()
	res, err := s.securityPersist.GetLatestResourceScan(ctx, gvk.Group, gvk.Version, gvk.Kind, obj.Namespace, obj.Name)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	err = s.userPermitted(ctx, res.OwningTeam)
	if err != nil {
		return nil, err
	}

	conv := DefaultConvertor.FromSecurityScanResult(res)
	return &conv, nil
}

func (s *securityImpl) GetScan(ctx context.Context, id uint64) (*securityv1.SecurityScanResult, error) {
	res, err := s.securityPersist.GetScan(ctx, id)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	err = s.userPermitted(ctx, res.OwningTeam)
	if err != nil {
		return nil, err
	}

	conv := DefaultConvertor.FromSecurityScanResult(res)
	return &conv, nil
}

func (s *securityImpl) ListRules(ctx context.Context) (*securityv1.SecurityRuleList, error) {
	rules := s.scanner.GetRules()
	ruleList := DefaultConvertor.FromSecurityRuleList(rules)
	return &ruleList, nil
}

func (s *securityImpl) GetRule(ctx context.Context, code string) (*securityv1.SecurityRule, error) {
	rule := s.scanner.GetRule(code)
	if rule == nil {
		return nil, nil
	}

	r := DefaultConvertor.FromSecurityRule(rule)
	return &r, nil
}

func (s *securityImpl) GetOverview(ctx context.Context) (*securityv1.SecurityOverview, error) {
	err := s.userPermitted(ctx, "")
	if err != nil {
		return nil, err
	}

	overview, err := s.securityPersist.GetOverview(ctx)
	if err != nil {
		return nil, err
	}

	r := DefaultConvertor.FromSecurityOverview(overview)
	return &r, nil
}

func (s *securityImpl) GetTeamOverview(ctx context.Context, team string) (*securityv1.SecurityOverview, error) {
	err := s.userPermitted(ctx, team)
	if err != nil {
		return nil, err
	}

	overview, err := s.securityPersist.GetTeamOverview(ctx, team)
	if err != nil {
		return nil, err
	}

	r := DefaultConvertor.FromSecurityOverview(overview)
	return &r, nil
}

func (s *securityImpl) userPermitted(ctx context.Context, owningTeam string) error {
	user := authentication.MustGetIdentity(ctx)
	if user.IsGlobalAdmin() {
		// Global admin always allowed
		return nil
	}

	if owningTeam == "" {
		// Not a team-scoped check, so not allowed if not admin.
		log.WithField(
			"username", user.Username(),
		).Warn("user trying to access the security logs")
		return NewErrNotAllowed("Must be global admin")
	}

	// Check if user in team which owns this resource
	for _, t := range user.Teams() {
		if t == owningTeam {
			// Permitted as team member
			return nil
		}
	}

	log.WithField(
		"username", user.Username(),
	).Warn("user trying to access the security logs")
	return NewErrNotAllowed("Must be global admin or in team which owns this resource")
}
