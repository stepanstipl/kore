/*
 * Copyright (C) 2019  Appvia Ltd <info@appvia.io>
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

package server

import (
	"context"

	api "github.com/appvia/kore/pkg/apiserver"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/services/users"
)

// Interface is the contract to the server
type Interface interface {
	// Run is responsible for starting the services
	Run(context.Context) error
	// Stop is responsible for trying to stop services
	Stop(context.Context) error
}

// KubernetesAPI is the configuration for the kubernetes api
type KubernetesAPI struct {
	// InCluster indicates we are running in cluster
	InCluster bool `json:"inCluster"`
	// MasterAPIURL specifies the kube-apiserver url
	MasterAPIURL string `json:"masterAPIUrl"`
	// Token is kubernetes token to authenticate to the api
	Token string `json:"token"`
	// KubeConfig is the kubeconfig path
	KubeConfig string
	// SkipTLSVerify indicates we skip tls
	SkipTLSVerify bool
}

// Config is the configuration of the various components
type Config struct {
	// APIServer is the config for the api server
	APIServer api.Config `json:"apiServer"`
	// Kubernetes is configuration for the api
	Kubernetes KubernetesAPI `json:"kubernetes"`
	// Kore is the configuration for the kore bridge
	Kore kore.Config `json:"kore"`
	// UsersMgr are the user management service options
	UsersMgr users.Config `json:"usersMgr"`
}
