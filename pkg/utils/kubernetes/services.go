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
	"time"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CreateOrUpdateService does what is says on the tin
func CreateOrUpdateService(ctx context.Context, cc client.Client, s *corev1.Service) (*corev1.Service, error) {
	if err := cc.Create(ctx, s); err != nil {
		if !kerrors.IsAlreadyExists(err) {
			return nil, err
		}

		key := types.NamespacedName{
			Namespace: s.Namespace,
			Name:      s.Name,
		}
		current := s.DeepCopy()
		if err := cc.Get(ctx, key, current); err != nil {
			return nil, err
		}

		s.SetResourceVersion(current.GetResourceVersion())
		s.SetGeneration(current.GetGeneration())

		return s, cc.Update(ctx, s)
	}

	return s, nil
}

// WaitForServiceEndpoint is used to wait for the service to generate an endpoint
func WaitForServiceEndpoint(ctx context.Context, cc client.Client, namespace, name string) (string, error) {
	logger := log.WithFields(log.Fields{
		"namespace": namespace,
		"name":      name,
	})
	logger.Debug("checking if the service has and endpoint yet")

	for {
		service := &v1.Service{}

		// @step: we break out or continue
		select {
		case <-ctx.Done():
			return "", errors.New("operation has been cancelled")
		default:
		}

		if err := cc.Get(ctx, types.NamespacedName{
			Namespace: namespace,
			Name:      name}, service); err != nil {

			logger.WithError(err).Debug("unable to retrieve a service endpoint, will retry")
		} else {
			if len(service.Status.LoadBalancer.Ingress) <= 0 {
				logger.Debug("loadbalancer does not have a status yet")
			} else {
				if service.Status.LoadBalancer.Ingress[0].IP != "" {
					logger.Debug("found an ip address for the service")

					return service.Status.LoadBalancer.Ingress[0].IP, nil
				}
			}
		}

		time.Sleep(15 * time.Second)
	}
}
