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
	"time"

	"github.com/appvia/kore/pkg/services/users"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/certificates"

	log "github.com/sirupsen/logrus"
)

// hubImpl is the implementation for the kore api
type hubImpl struct {
	// config is the configuration of the kore
	config *Config
	// store is the access layer / kubernetes api
	store store.Store
	// idp is the idp implementation
	idp *idpImpl
	// invitations handles generated links
	invitations Invitations
	// teams is the teams implementation
	teams Teams
	// users is the users implementation
	users Users
	// plans
	plans Plans
	// usermgr is the user manager
	usermgr users.Interface
	// signer is used to sign off client certs
	signer certificates.Signer
}

// New returns a new instance of the kore bridge
func New(sc store.Store, usermgr users.Interface, config Config) (Interface, error) {
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

	// @step: create a signer for client certificates
	signer, err := certificates.NewSignerFromFiles(
		config.CertificateAuthority,
		config.CertificateAuthorityKey,
	)
	if err != nil {
		log.WithError(err).Error("trying to create certificate signer")

		return nil, err
	}

	// @step: enable the open

	h := &hubImpl{
		config: &config,
		store:  sc,
		signer: signer,
	}
	h.idp = &idpImpl{Interface: h}
	h.invitations = &ivImpl{Interface: h}
	h.plans = &plansImpl{Interface: h}
	h.teams = &teamsImpl{hubImpl: h}
	h.usermgr = usermgr
	h.users = &usersImpl{hubImpl: h}

	// @step: call the setup code for the kore
	if err := h.Setup(context.Background()); err != nil {
		return nil, err
	}

	return h, nil
}

// SignedClientCertificate is used to generate a client certificate
func (h hubImpl) SignedClientCertificate(identity, team string) ([]byte, []byte, error) {
	logger := log.WithFields(log.Fields{
		"identity": identity,
		"team":     team,
	})
	logger.Info("generating a client certificate for remote cluster")

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
	logger.Info("generating a server certificate")

	cert, key, err := h.signer.GenerateServer(hosts, duration)
	if err != nil {
		logger.WithError(err).Error("trying to generate server certificate")

		return []byte{}, []byte{}, err
	}

	return cert, key, nil
}

// Audit returns the auditor
func (h *hubImpl) Audit() users.Audit {
	return h.usermgr.Audit()
}

// Users returns the user implementation
func (h hubImpl) Users() Users {
	return h.users
}

// Plans returns a plans interface
func (h hubImpl) Plans() Plans {
	return h.plans
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
