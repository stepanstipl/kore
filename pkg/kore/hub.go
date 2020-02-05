/*
 * Copyright (C) 2019  Appvia Ltd <info@appvia.io>
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
	"fmt"
	"time"

	"github.com/appvia/kore/pkg/services/audit"
	"github.com/appvia/kore/pkg/services/users"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/certificates"

	log "github.com/sirupsen/logrus"
)

// hubImpl is the implementation for the kore api
type hubImpl struct {
	// auditor is the audit interface
	auditor audit.Interface
	// config is the configuration of the kore
	config *Config
	// store is the access layer / kubernetes api
	store store.Store
	// idp is the idp implimentation
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
func New(sc store.Store, usermgr users.Interface, auditor audit.Interface, config Config) (Interface, error) {
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
		auditor: auditor,
		config:  &config,
		store:   sc,
		signer:  signer,
	}
	h.invitations = &ivImpl{Interface: h}
	h.plans = &plansImpl{Interface: h}
	h.teams = &teamsImpl{hubImpl: h}
	h.users = &usersImpl{hubImpl: h}
	h.usermgr = usermgr
	h.idp = &idpImpl{Interface: h}

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
func (h *hubImpl) Audit() audit.Interface {
	return h.auditor
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
