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

package awsservicebroker

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
)

func getServiceAccountToken(ctx context.Context, client client.Client, namespace, name string) (_ *corev1.Secret, _ error) {
	sa := &corev1.ServiceAccount{}
	err := client.Get(ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, sa)
	if err != nil {
		return nil, fmt.Errorf("failed to get serviceaccount %q: %w", name, err)
	}
	if len(sa.Secrets) <= 0 {
		return nil, fmt.Errorf("no secrets found in serviceaccount %q", name)
	}

	return getSecret(ctx, client, namespace, sa.Secrets[0].Name)
}

func getSecret(ctx context.Context, client client.Client, namespace, name string) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	err := client.Get(ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret %q: %w", name, err)
	}

	return secret, nil
}
