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

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	clusterappman "github.com/appvia/kore/pkg/clusterappman"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EnsureClusterman will ensure clusterappman is deployed
func (a k8sCtrl) EnsureClusterman(ctx context.Context, cc client.Client, image string) (*corev1.Components, error) {
	logger := log.WithFields(log.Fields{"controller": a.Name()})

	return clusterappman.Deploy(ctx, cc, logger, image)
}
