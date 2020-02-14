/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
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
