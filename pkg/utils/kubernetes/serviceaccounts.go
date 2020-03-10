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

package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	k8s "k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CreateOrUpdateServiceAccount is used to ensure a service account
func CreateOrUpdateServiceAccount(ctx context.Context, cc client.Client, account *corev1.ServiceAccount) (*corev1.ServiceAccount, error) {
	if err := cc.Create(ctx, account); err != nil {
		if !kerrors.IsAlreadyExists(err) {
			return nil, err
		}
		// @step: we need to retrieve the current one
		original := account.DeepCopy()

		if err := cc.Get(ctx, types.NamespacedName{
			Name:      account.Name,
			Namespace: account.Namespace,
		}, original); err != nil {
			return nil, err
		}

		return account, cc.Patch(ctx, account, client.MergeFrom(original))
	}

	return account, nil
}

// WaitForServiceAccountToken waits for the service account to get a token
func WaitForServiceAccountToken(client k8s.Interface, namespace, name string, timeout time.Duration) (*corev1.Secret, error) {
	var secret *corev1.Secret
	doneCh := make(chan struct{})

	// @step: create a context for the go rountine to run under
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			if err := func() error {
				sa, err := client.CoreV1().ServiceAccounts(namespace).Get(name, metav1.GetOptions{})
				if err != nil {
					return err
				}
				if len(sa.Secrets) <= 0 {
					return errors.New("no secrets")
				}
				secret, err = client.CoreV1().Secrets(namespace).Get(sa.Secrets[0].Name, metav1.GetOptions{})
				if err != nil {
					return err
				}

				return nil
			}(); err == nil {
				doneCh <- struct{}{}
				return
			}
		}
	}()

	select {
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout after: %s waiting for service account: %s/%s secrets", timeout, namespace, name)
	case <-doneCh:
	}

	return secret, nil
}
