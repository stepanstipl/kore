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
	"fmt"
	"io/ioutil"
	"time"

	"github.com/appvia/kore/pkg/costs"
	"github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/security"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/certificates"

	log "github.com/sirupsen/logrus"
)

// hubImpl is the implementation for the kore api
type hubImpl struct {
	// caAuthority is the certificate authority for kore
	caAuthority []byte
	// caKey is the certificate authority key
	caKey []byte
	// config is the configuration of the kore
	config *Config
	// store is the access layer / kubernetes api
	store store.Store
	// accounts implementation
	accounts Accounts
	// idp is the idp implementation
	idp *idpImpl
	// alerting is the alerts implementation
	alerting AlertRules
	// invitations handles generated links
	invitations Invitations
	// teams is the teams implementation
	teams Teams
	// users is the users implementation
	users Users
	// plans
	plans Plans
	// plan policies
	planPolicies PlanPolicies
	// persistenceMgr is the persistence manager
	persistenceMgr persistence.Interface
	// signer is used to sign off client certs
	signer certificates.Signer
	// audit is the audit implementation
	audit Audit
	// serviceplans is the ServicePlans implementation
	servicePlans ServicePlans
	// servicekinds is the ServiceKinds implementation
	serviceKinds ServiceKinds
	// serviceProviders is the ServiceProviders implementation
	serviceProviders ServiceProviders
	// security provides the ability to scan kore objects for security compliance
	security Security
	// configs provides the ability to store key value pairs
	configs Configs
	// costs is the costs implementation
	costs costs.Costs
	// features is the features implementation
	features KoreFeatures
}

// New returns a new instance of the kore bridge
func New(sc store.Store, persistenceMgr persistence.Interface, config Config) (Interface, error) {
	log.Info("initializing the kore api bridge")

	// @step: check the options
	if err := config.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid options: %s", err)
	}

	// @step: ensure we have a hmax for the signing of things
	if !config.HasHMAC() {
		log.Warn("no hmac for kore was provided, generating a random one (this has consequences!)")
		config.HMAC = utils.Random(32)
	}

	// @step: read the public and private keys
	authority, err := ioutil.ReadFile(config.CertificateAuthority)
	if err != nil {
		return nil, err
	}
	key, err := ioutil.ReadFile(config.CertificateAuthorityKey)
	if err != nil {
		return nil, err
	}

	// @step: create a signer for client certificates
	signer, err := certificates.NewSignerFromFiles(
		config.CertificateAuthority,
		config.CertificateAuthorityKey,
	)
	if err != nil {
		log.WithError(err).Error("trying to create certificate signer")

		return nil, err
	}

	h := &hubImpl{
		caAuthority: authority,
		caKey:       key,
		config:      &config,
		store:       sc,
		signer:      signer,
	}
	h.accounts = &accountsImpl{Interface: h}
	h.alerting = &alertsImpl{Interface: h}
	h.idp = &idpImpl{Interface: h}
	h.invitations = &ivImpl{Interface: h}
	h.plans = &plansImpl{Interface: h}
	h.planPolicies = &planPoliciesImpl{Interface: h}
	h.teams = &teamsImpl{hubImpl: h}
	h.persistenceMgr = persistenceMgr
	h.users = &usersImpl{hubImpl: h}
	h.audit = &auditImpl{auditPersist: persistenceMgr.Audit()}
	h.servicePlans = &servicePlansImpl{Interface: h}
	h.serviceKinds = &serviceKindsImpl{Interface: h}
	h.serviceProviders = &serviceProvidersImpl{Interface: h}
	h.security = &securityImpl{
		scanner:         security.New(),
		securityPersist: persistenceMgr.Security(),
	}
	h.configs = &configImpl{hubImpl: h}
	h.costs = costs.New(&config.Costs)
	h.features = &koreFeaturesImpl{store: h.store}

	// @step: call the setup code for the kore
	if err := h.Setup(context.Background()); err != nil {
		return nil, err
	}

	return h, nil
}

// CertificateAuthority returns the ca pem authority
func (h *hubImpl) CertificateAuthority() []byte {
	return h.caAuthority
}

// CertificateAuthority returns the ca pem authority key
func (h *hubImpl) CertificateAuthorityKey() []byte {
	return h.caKey
}

// SignedClientCertificate is used to generate a client certificate
func (h hubImpl) SignedClientCertificate(identity, team string) ([]byte, []byte, error) {
	logger := log.WithFields(log.Fields{
		"identity": identity,
		"team":     team,
	})
	logger.Debug("generating a client certificate for remote cluster")

	cert, key, err := h.signer.GenerateClient(identity, team, 24*365*time.Hour)
	if err != nil {
		logger.WithError(err).Error("trying to generate client certificate")

		return []byte{}, []byte{}, err
	}

	return cert, key, nil
}

// SignedServerCertificate is used to generate a server certificate
func (h hubImpl) SignedServerCertificate(hosts []string, duration time.Duration) ([]byte, []byte, error) {
	logger := log.WithFields(log.Fields{
		"duration": duration.String(),
		"hosts":    hosts,
	})
	logger.Debug("generating a server certificate")

	cert, key, err := h.signer.GenerateServer(hosts, duration)
	if err != nil {
		logger.WithError(err).Error("trying to generate server certificate")

		return []byte{}, []byte{}, err
	}

	return cert, key, nil
}

// Accounts return the account implementation
func (h *hubImpl) Accounts() Accounts {
	return h.accounts
}

// AlertRules returns the alerting implementation
func (h *hubImpl) AlertRules() AlertRules {
	return h.alerting
}

// Audit returns the auditor
func (h *hubImpl) Audit() Audit {
	return h.audit
}

// Users returns the user implementation
func (h hubImpl) Users() Users {
	return h.users
}

// Plans returns a plans interface
func (h hubImpl) Plans() Plans {
	return h.plans
}

// PlanPolicies returns a plan policies interface
func (h hubImpl) PlanPolicies() PlanPolicies {
	return h.planPolicies
}

// Invitations returns the invitations implementation
func (h hubImpl) Invitations() Invitations {
	return h.invitations
}

// Teams returns the team implementation
func (h hubImpl) Teams() Teams {
	return h.teams
}

// Auth returns the authentication interface
func (h *hubImpl) IDP() IDP {
	return h.idp
}

// Config returns the store configuration
func (h hubImpl) Config() *Config {
	return h.config
}

// Store returns underlying data layer
func (h hubImpl) Store() store.Store {
	return h.store
}

// ServicePlans returns the serviceplans interface
func (h hubImpl) ServicePlans() ServicePlans {
	return h.servicePlans
}

// ServiceKinds returns the service interface
func (h hubImpl) ServiceKinds() ServiceKinds {
	return h.serviceKinds
}

func (h hubImpl) ServiceProviders() ServiceProviders {
	return h.serviceProviders
}

func (h hubImpl) Security() Security {
	return h.security
}

func (h hubImpl) Persist() persistence.Interface {
	return h.persistenceMgr
}

func (h hubImpl) Configs() Configs {
	return h.configs
}

func (h hubImpl) Costs() costs.Costs {
	return h.costs
}

func (h hubImpl) Features() KoreFeatures {
	return h.features
}
