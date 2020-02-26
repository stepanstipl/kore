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

package openid

import (
	"context"
	"errors"
	"sync"
	"time"

	oidc "github.com/coreos/go-oidc"
	log "github.com/sirupsen/logrus"
)

type authImpl struct {
	sync.RWMutex
	// config is the configuration
	Config
	// provider is the oidc provider
	provider *oidc.Provider
	// verifier is the oidc token config
	verifier *oidc.IDTokenVerifier
}

// New creates and returns a authenticator
func New(config Config) (Authenticator, error) {
	if err := config.IsValid(); err != nil {
		return nil, err
	}

	return &authImpl{Config: config}, nil
}

// Provider returns the oidc provider
func (a *authImpl) Provider() *oidc.Provider {
	return a.provider
}

// Verify is called to verify the token
func (a *authImpl) Verify(ctx context.Context, token string) (*oidc.IDToken, error) {
	// @step: we lock the struct and check if the verifier has been configured yet
	verifier := a.GetVerifier()
	if verifier == nil {
		log.Info("openid has not been configured yet")

		return nil, errors.New("unable to verify token")
	}

	return verifier.Verify(ctx, token)
}

// RunWithSync waits for the discovery process to occur
func (a *authImpl) RunWithSync(ctx context.Context) error {
	if err := a.Run(ctx); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)

	nctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	for {
		if a.provider != nil && a.verifier != nil {
			return nil
		}
		select {
		case <-nctx.Done():
			return errors.New("context has been cancelled")
		default:
		}

		time.Sleep(5 * time.Second)
	}
}

// Run starts the request to the discovery url
func (a *authImpl) Run(ctx context.Context) error {
	go func() {
		for {
			err := func() error {
				log.WithFields(log.Fields{
					"discovery-url": a.DiscoveryURL,
				}).Info("attempting to retrieve provider details via discovery url")

				// @step: attempt to retrieve the details for the discovery url
				provider, err := oidc.NewProvider(ctx, a.DiscoveryURL)
				if err != nil {
					return err
				}

				verifier := provider.Verifier(&oidc.Config{
					ClientID:          a.ClientID,
					SkipClientIDCheck: a.SkipClientIDCheck,
					SkipExpiryCheck:   false,
				})

				log.Info("openid authentication plugin successfully retrieve configuration")

				a.Lock()
				defer a.Unlock()

				a.provider = provider
				a.verifier = verifier

				return nil
			}()
			if err != nil {
				log.WithError(err).Error("failed to retrieve provider configuration")

				time.Sleep(10 * time.Second)
				continue
			}

			return
		}
	}()

	return nil
}

// GetVerifier returns the internal token verifier
func (a *authImpl) GetVerifier() *oidc.IDTokenVerifier {
	a.RLock()
	defer a.RUnlock()

	return a.verifier
}
