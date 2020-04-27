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

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Security represents the interface to the top-level Kore Security service.
type Security interface {
	ScanPlan(ctx context.Context, plan *configv1.Plan) error
	ScanCluster(ctx context.Context, cluster *clustersv1.Cluster) error
	ListScans(ctx context.Context, latestOnly bool) (*securityv1.ScanResultList, error)
	ScanHistoryForResource(ctx context.Context, typ metav1.TypeMeta, obj metav1.ObjectMeta) (*securityv1.ScanResultList, error)
	GetCurrentScanForResource(ctx context.Context, typ metav1.TypeMeta, obj metav1.ObjectMeta) (*securityv1.ScanResult, error)
	GetScan(ctx context.Context, id uint64) (*securityv1.ScanResult, error)
}

type securityImpl struct {
	scanner         security.Scanner
	securityPersist persistence.Security
}

func (s *securityImpl) ScanPlan(ctx context.Context, plan *configv1.Plan) error {
	scanResult := s.scanner.ScanPlan(plan)
	return s.persistScan(ctx, scanResult)
}

func (s *securityImpl) ScanCluster(ctx context.Context, cluster *clustersv1.Cluster) error {
	scanResult := s.scanner.ScanCluster(cluster)
	return s.persistScan(ctx, scanResult)
}

func (s *securityImpl) persistScan(ctx context.Context, scanResult *securityv1.ScanResult) error {
	scanResultDB := DefaultConvertor.ToSecurityScanResult(scanResult)
	err := s.securityPersist.StoreScan(ctx, &scanResultDB)
	if err != nil {
		log.WithError(err).Error("trying to persist security security scan")
		return err
	}
	return nil
}

func (s *securityImpl) ListScans(ctx context.Context, latestOnly bool) (*securityv1.ScanResultList, error) {
	user := authentication.MustGetIdentity(ctx)
	if !user.IsGlobalAdmin() {
		log.WithField(
			"username", user.Username(),
		).Warn("user trying to access the security logs")

		return nil, NewErrNotAllowed("Must be global admin")
	}

	res, err := s.securityPersist.ListScans(ctx, latestOnly)
	if err != nil {
		return nil, err
	}
	result := &securityv1.ScanResultList{}
	result.Items = make([]securityv1.ScanResult, len(res))
	for i, r := range res {
		result.Items[i] = DefaultConvertor.FromSecurityScanResult(r)
	}
	return result, nil
}

func (s *securityImpl) ScanHistoryForResource(ctx context.Context, typ metav1.TypeMeta, obj metav1.ObjectMeta) (*securityv1.ScanResultList, error) {
	user := authentication.MustGetIdentity(ctx)
	if !user.IsGlobalAdmin() {
		log.WithField(
			"username", user.Username(),
		).Warn("user trying to access the security logs")

		return nil, NewErrNotAllowed("Must be global admin")
	}

	res, err := s.securityPersist.ListResourceScanHistory(ctx, typ.APIVersion, typ.Kind, obj.Namespace, obj.Name)
	if err != nil {
		return nil, err
	}
	result := &securityv1.ScanResultList{}
	result.Items = make([]securityv1.ScanResult, len(res))
	for i, r := range res {
		result.Items[i] = DefaultConvertor.FromSecurityScanResult(r)
	}
	return result, nil
}

func (s *securityImpl) GetCurrentScanForResource(ctx context.Context, typ metav1.TypeMeta, obj metav1.ObjectMeta) (*securityv1.ScanResult, error) {
	user := authentication.MustGetIdentity(ctx)
	if !user.IsGlobalAdmin() {
		log.WithField(
			"username", user.Username(),
		).Warn("user trying to access the security logs")

		return nil, NewErrNotAllowed("Must be global admin")
	}

	res, err := s.securityPersist.GetLatestResourceScan(ctx, typ.APIVersion, typ.Kind, obj.Namespace, obj.Name)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	conv := DefaultConvertor.FromSecurityScanResult(res)
	return &conv, nil
}

func (s *securityImpl) GetScan(ctx context.Context, id uint64) (*securityv1.ScanResult, error) {
	user := authentication.MustGetIdentity(ctx)
	if !user.IsGlobalAdmin() {
		log.WithField(
			"username", user.Username(),
		).Warn("user trying to access the security logs")

		return nil, NewErrNotAllowed("Must be global admin")
	}

	res, err := s.securityPersist.GetScan(ctx, id)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	conv := DefaultConvertor.FromSecurityScanResult(res)
	return &conv, nil
}
