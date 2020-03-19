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

package authproxy

import (
	"context"
	"errors"

	"github.com/appvia/kore/pkg/utils/openid"

	"github.com/coreos/go-oidc"
	log "github.com/sirupsen/logrus"
)

type verifierWrapper struct {
	verifier *oidc.IDTokenVerifier
}

func (v verifierWrapper) Verify(ctx context.Context, rawIDToken string) (openid.IDToken, error) {
	return v.verifier.Verify(ctx, rawIDToken)
}

func CreateVerifier(config Config) (openid.Verifier, error) {
	options := &oidc.Config{
		ClientID:          config.IDPClientID,
		SkipClientIDCheck: true,
		SkipExpiryCheck:   false,
	}

	switch {
	case config.SigningCA != "":
		log.WithField(
			"signing_ca", config.SigningCA,
		).Info("using the signing certificate to verify the requests")

		keyset, err := newStaticKeySet(config.SigningCA)
		if err != nil {
			return nil, err
		}

		return verifierWrapper{
			verifier: oidc.NewVerifier(config.IDPClientID, keyset, options),
		}, nil
	case config.IDPServerURL != "":
		log.WithField(
			"idp-server-url", config.IDPServerURL,
		).Info("using the IDP server to verify the requests")

		provider, err := oidc.NewProvider(context.Background(), config.IDPServerURL)
		if err != nil {
			log.WithError(err).Error("trying to retrieve provider details")

			return nil, err
		}

		return verifierWrapper{
			verifier: provider.Verifier(options),
		}, nil
	default:
		panic(errors.New("unable to create verifier from the given configuration"))
	}
}
