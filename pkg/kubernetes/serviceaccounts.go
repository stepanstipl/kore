/*
 * Copyright (C) 2019  Rohith Jayawardene <gambol99@gmail.com>
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

package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
)

// WaitForServiceAccountToken waits for the service account to get a token
func WaitForServiceAccountToken(client k8s.Interface, namespace, name string, timeout time.Duration) (*corev1.Secret, error) {
	var secret *corev1.Secret
	doneCh := make(chan struct{}, 0)

	// @step: create a context for the go rountine to run under
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		fmt.Println("dsjdksjds")
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			fmt.Println("checking")

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
