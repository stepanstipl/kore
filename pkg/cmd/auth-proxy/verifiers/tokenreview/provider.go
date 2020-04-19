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

package tokenreview

import (
	"net/http"
	"time"

	"github.com/appvia/kore/pkg/cmd/auth-proxy/verifiers"
	"github.com/appvia/kore/pkg/utils"

	pcache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	authentication "k8s.io/api/authentication/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type kvImpl struct {
	Options
	client kubernetes.Interface
	cache  *pcache.Cache
}

// Options are options of the verifier
type Options struct {
	// Audiences are placed in the token review
	Audiences []string
	// CacheFailure is the duration to cache failures
	CacheFailure time.Duration
	// CacheSuccess is the duration for caching successes
	CacheSuccess time.Duration
}

// NewFromClient creates and returns a tokenreview from a k8s client
func NewFromClient(client kubernetes.Interface, options Options) (verifiers.Interface, error) {
	log.WithField(
		"audiences", options.Audiences,
	).Debug("using the following audiences")

	return &kvImpl{
		Options: options,
		client:  client,
		cache:   pcache.New(10*time.Minute, 20*time.Minute),
	}, nil
}

// New creates and returns a verifier for the k8s cluster
func New(config *rest.Config, options Options) (verifiers.Interface, error) {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return NewFromClient(client, options)

}

// Admit checks the token is valid
func (k *kvImpl) Admit(request *http.Request) (bool, error) {
	// @step: extract the token from the request
	bearer, found := utils.GetBearerToken(request.Header.Get("Authorization"))
	if !found {
		return false, nil
	}
	log.Debug("checking against the tokenreview verifier")

	// @step: check the cache for the token
	review, err := func() (*authentication.TokenReview, error) {
		review, found := k.cache.Get(bearer)
		if found {
			log.Debug("found the token review in the cache")

			return review.(*authentication.TokenReview), nil
		}

		resp, err := k.client.AuthenticationV1().TokenReviews().Create(
			&authentication.TokenReview{
				Spec: authentication.TokenReviewSpec{
					Token:     bearer,
					Audiences: k.Audiences,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}()
	if err != nil {
		log.WithError(err).Error("trying to retrieve review of token")

		return false, err
	}

	// @step: check the token review
	if len(review.Status.Error) > 0 {
		log.Debug("token failed review against the api")

		k.cache.Set(bearer, review, k.CacheFailure)

		return false, nil
	}
	log.Debug("successfully authentication against tokenreview")

	k.cache.Set(bearer, review, k.CacheSuccess)

	return true, nil
}
