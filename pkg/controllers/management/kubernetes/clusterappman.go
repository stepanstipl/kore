/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package kubernetes

import (
	"context"
	"strings"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	clusterappman "github.com/appvia/kore/pkg/clusterappman"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	"github.com/appvia/kore/pkg/version"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	// KoreImage is the container image for Kore
	KoreImage = "quay.io/appvia/kore-apiserver:" + version.Release
)

const (
	// clusterappmanNamespace is the namespace the clusterappmanager runs in
	clusterappmanNamespace = clusterappman.KoreNamespace
	// clusterappmanConfig is the name of the ConfigMap configuration required for kore cluster manager
	clusterappmanConfig = clusterappman.ParamsConfigMap
	// clusterappmanConfigKey is the configmap Key to store the cluster data
	clusterappmanConfigKey = clusterappman.ParamsConfigKey
	// clusterappmanDeployment
	clusterappmanDeployment = "kore-clusterappman"
)

// EnsureClusterman will ensure clusterappman is deployed
func (a k8sCtrl) EnsureClusterman(ctx context.Context, cc client.Client, cluster *clustersv1.Kubernetes) error {
	logger := log.WithFields(log.Fields{"controller": a.Name()})

	provider := strings.ToLower(cluster.Spec.Provider.Kind)

	params, err := a.GetClusterConfiguration(ctx, cluster, provider)
	if err != nil {

		return err
	}

	// @step: check if the cluster manager namespace exists and create it if not
	if err := EnsureNamespace(ctx, cc, clusterappmanNamespace); err != nil {
		logger.WithError(err).Errorf("trying to create the kore cluster-manager namespace %s", clusterappmanNamespace)

		return err
	}

	// @step: check if the cluster config exists
	found, err := ConfigExists(ctx, cc)
	if err != nil {
		logger.WithError(err).Error("failed to check for kore clusterappman config")

		return err

	}
	if !found {
		if err := CreateConfig(ctx, cc, params); err != nil {
			logger.WithError(err).Error("trying to create the kore cluster-manager configuration configmap")

			return err
		}
	}

	// @step: ensure the service account
	if _, err := kubernetes.CreateOrUpdateServiceAccount(ctx, cc, &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "clusterappman",
			Namespace: clusterappmanNamespace,
			Labels: map[string]string{
				kore.Label("owner"): "true",
			},
		},
	}); err != nil {
		logger.WithError(err).Error("trying to create the clusterappman service account")

		return err
	}
	// @step setup correct permissions for deployment
	if err := CreateClusterManClusterRoleBinding(ctx, cc); err != nil {
		logger.WithError(err).Error("can not create cluster-manager clusterrole")

		return err
	}

	// @step: check if the kore cluster manager deployment exists
	logger.Debugf("deploying clusterappman using image %s", KoreImage)
	if err := CreateOrUpdateClusterAppManDeployment(ctx, cc); err != nil {
		logger.WithError(err).Error("trying to create the cluster manager deployment")

		return err
	}
	logger.Debug("waiting for kore cluster manager deployment status to appear")

	nctx, cancel := context.WithTimeout(ctx, 4*time.Minute)
	defer cancel()

	logger.Info("waiting for kore cluster manager to complete")

	// @step: wait for the clusterappman deployment to complete
	if err := WaitOnStatus(nctx, cc); err != nil {
		logger.WithError(err).Error("failed waiting for kore cluster manager status to complete")

		return err
	}

	logger.Info("kube clusterappman ready for new cluster")

	return nil
}

// GetClusterConfiguration is responsible for generate the parameters for the cluster
func (a k8sCtrl) GetClusterConfiguration(ctx context.Context, cluster *clustersv1.Kubernetes, provider string) (Parameters, error) {
	params := Parameters{
		Domain:       cluster.Spec.Domain,
		Provider:     provider,
		StorageClass: "default",
	}
	switch provider {
	case "gke":
		params.StorageClass = "standard"
	}

	return params, nil
}
