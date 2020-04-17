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
	"errors"
	"net/http"

	"github.com/appvia/kore/pkg/utils"
	"github.com/coreos/go-oidc"
	log "github.com/sirupsen/logrus"
)

type kvImpl struct {
	verifier *oidc.IDTokenVerifier
}

// NewKubeVerifier creates and returns a verifier for the k8s cluster
func NewKubeVerifier(caPath string) (Verifier, error) {
	log.WithField(
		"signing_ca", caPath,
	).Info("using the signing certificate to verify the requests")

	options := &oidc.Config{
		ClientID:          "none",
		SkipClientIDCheck: true,
		SkipExpiryCheck:   false,
	}

	keyset, err := newStaticKeySet(caPath)
	if err != nil {
		return nil, err
	}

	return &kvImpl{verifier: oidc.NewVerifier("kubernetes/serviceaccount", keyset, options)}, nil
}

// Admit checks the token is valid
func (k *kvImpl) Admit(request *http.Request) (bool, error) {
	// @step: extract the token from the request
	bearer, found := utils.GetBearerToken(request.Header.Get("Authorization"))
	if !found {
		return false, errors.New("no authorization token")
	}

	// @step: parse and extract the identity
	_, err := k.verifier.Verify(request.Context(), bearer)
	if err != nil {
		return false, err
	}

	return true, nil
}
