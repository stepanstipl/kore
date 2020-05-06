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

package controllers

import (
	"context"
	"fmt"

	config "github.com/appvia/kore/pkg/apis/config/v1"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetDecodedSecret accepts a credentialsRef and passes back a reference to the decoded Secret
func GetDecodedSecret(ctx context.Context, cc client.Client, credentialsRef *v1.SecretReference) (*config.Secret, error) {
	// @step: we need to grab the secret
	secret := &config.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      credentialsRef.Name,
			Namespace: credentialsRef.Namespace,
		},
	}

	found, err := kubernetes.GetIfExists(ctx, cc, secret)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("Secret for credentialsRef: (%s/%s) not found", credentialsRef.Namespace, credentialsRef.Name)
	}

	// @step: ensure the secret is decoded before using
	if err := secret.Decode(); err != nil {
		return nil, err
	}

	return secret, nil
}
