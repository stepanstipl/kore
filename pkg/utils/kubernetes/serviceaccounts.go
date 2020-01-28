/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of hub-apiserver.
 *
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
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
