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

/*
import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// WaitForServiceEndpoint is used to wait for the service to generate an endpoint
func WaitForServiceEndpoint(ctx context.Context, cc client.Client, namespace, name string) (string, error) {
	logger := log.WithFields(log.Fields{
		"namespace": namespace,
		"name":      name,
	})
	logger.Debug("checking if the service has and endpoint yet")
	var endpoint string

	for {
		err := func() error {
			select {
			case <-ctx.Done():
			default:
			}
			service := &v1.Service{}
			if err := cc.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, service); err != nil {
				return err
			}

		}()
		if err != nil {
			logger.WithError(err).Debug("unable to retrieve a service endpoint, will retry")
		}
		if err == nil {
			logger

		}

		time.Sleep(10 * time.Second)
	}
}
*/