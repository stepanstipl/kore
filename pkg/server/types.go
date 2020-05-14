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

package server

import (
	"context"

	api "github.com/appvia/kore/pkg/apiserver"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/utils/kubernetes"
)

// Interface is the contract to the server
type Interface interface {
	// Run is responsible for starting the services
	Run(context.Context) error
	// Stop is responsible for trying to stop services
	Stop(context.Context) error
}

// Config is the configuration of the various components
type Config struct {
	// APIServer is the config for the api server
	APIServer api.Config `json:"apiServer"`
	// Kubernetes is configuration for the api
	Kubernetes kubernetes.KubernetesAPI `json:"kubernetes"`
	// Kore is the configuration for the kore bridge
	Kore kore.Config `json:"kore"`
	// PersistenceMgr are the options for the persistence manager
	PersistenceMgr persistence.Config `json:"persistenceMgr"`
}
