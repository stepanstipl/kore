/*
 * Copyright (C) 2019 Appvia Ltd <info@appvia.io>
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

package bootstrap

import (
	"context"
	"strings"
	"time"

	clusterv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/clusterman"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	"github.com/appvia/kore/pkg/utils/openid"
	"github.com/appvia/kore/pkg/version"

	log "github.com/sirupsen/logrus"
	apps "k8s.io/api/apps/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	// convertor is the unstructured converter client
	converter = runtime.DefaultUnstructuredConverter
	// KoreImage is the container image for Kore
	KoreImage = "quay.io/appvia/kore-apiserver:" + version.Release
)

const (
	finalizerName = "bootstrap.compute.kore.appvia.io"
	// clustermanNamespace is the namespace the clustermanager runs in
	clustermanNamespace = clusterman.KoreNamespace
	// clustermanConfig is the name of the ConfigMap configuration required for kore cluster manager
	clustermanConfig = clusterman.ParamsConfigMap
	// clustermanConfigKey is the configmap Key to store the cluster data
	clustermanConfigKey = clusterman.ParamsConfigKey
	// clustermanDeployment
	clustermanDeployment = "kore-clusterman"
)

// Reconcile is the entrypoint for the reconcilation logic
// @QUESTION we could move this into a Kubernetes Job it self and allow that to update
// the resource status rather than using a controller?
func (t bsCtrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to reconcile the applications incluster")

	// @step: retrieve the resource from the api
	cluster := &clusterv1.Kubernetes{}
	if err := t.mgr.GetClient().Get(ctx, request.NamespacedName, cluster); err != nil {
		if !kerrors.IsNotFound(err) {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}
	original := cluster.DeepCopy()

	// @step: ignore the resource if already bootstrapped
	status, found := cluster.Status.Components.GetStatus("provision")
	if !found {
		logger.Warn("cluster does not have a status on the provisioning yet")

		return reconcile.Result{RequeueAfter: 2 * time.Minute}, nil
	}
	if status.Status != corev1.SuccessStatus {
		logger.Warn("cluster provision is not successfully yet, waiting")

		return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	if _, found = cluster.Status.Components.GetStatus("applications"); !found {
		cluster.Status.Components.SetCondition(corev1.Component{
			Name:    "applications",
			Status:  corev1.PendingStatus,
			Message: "attempting to install the applications",
		})

		if err := t.mgr.GetClient().Status().Patch(ctx, cluster, client.MergeFrom(original)); err != nil {
			log.WithError(err).Error("trying to update the resource status")

			return reconcile.Result{RequeueAfter: 1 * time.Minute}, err
		}

		return reconcile.Result{Requeue: true}, nil
	}

	// @logic
	// - create a kubernetes client to the remote cluster
	// - retrieve the credentials for the broker from the cluster provider
	// - wait for kube api to be ready
	// - deploy the kore-cluster-manager
	// - wait for critical kore-cluster componets to be ready
	result, err := func() (reconcile.Result, error) {
		var credentials Credentials
		var provider string

		if kore.IsProviderBacked(cluster) {
			// @step: we need to grab the credentials for the cloud provider
			us, err := controllers.GetCloudProviderCredentials(ctx, t.mgr.GetClient(), cluster)
			if err != nil {
				logger.WithError(err).Error("trying to retrieve cloud provider credentials")

				return reconcile.Result{}, err
			}

			provider = strings.ToLower(cluster.Spec.Provider.Kind)

			switch provider {
			case "gke":
				u := &gke.GKE{}
				if err := converter.FromUnstructured(us.Object, u); err != nil {
					logger.WithError(err).Error("trying to convert object type")

					return reconcile.Result{}, err
				}
				creds := &gke.GKECredentials{}
				creds.SetName(u.Spec.Credentials.Name)
				creds.SetNamespace(u.Spec.Credentials.Namespace)
				ref := types.NamespacedName{
					Namespace: u.Spec.Credentials.Namespace,
					Name:      u.Spec.Credentials.Name,
				}

				if err := t.mgr.GetClient().Get(ctx, ref, creds); err != nil {
					return reconcile.Result{}, err
				}

				credentials.GKE.Account = creds.Spec.Account
			case "eks":
			default:
				logger.WithField("provider", provider).Warn("unknown cloud provider")
			}
		}

		// @step: build a client from the cluster secret
		client, err := controllers.CreateClientFromSecret(ctx, t.mgr.GetClient(), cluster.Namespace, cluster.Name)
		if err != nil {
			logger.WithError(err).Error("trying to retrieve the credentials")

			return reconcile.Result{}, err
		}

		params, err := t.GetClusterConfiguration(ctx, cluster, strings.ToLower(cluster.Spec.Provider.Kind))
		if err != nil {
			logger.WithError(err).Error("failed to generate the parameters")

			return reconcile.Result{}, err
		}
		params.Credentials = credentials

		// @step: check if the cluster manager namespace exists and create it if not
		if err := EnsureNamespace(ctx, client); err != nil {
			logger.WithError(err).Errorf("trying to create the kore cluster-manager namespace %s", clustermanNamespace)

			return reconcile.Result{}, err
		}

		// @step: check if the cluster config exists
		found, err = ConfigExists(ctx, client)
		if err != nil {
			logger.WithError(err).Error("failed to check for kore clusterman config")

			return reconcile.Result{}, err

		}
		if !found {
			if err := CreateConfig(ctx, client, params); err != nil {
				logger.WithError(err).Error("trying to create the kore cluster-manager configuration configmap")

				return reconcile.Result{}, err
			}
		}

		// @step setup correct permissions for deployment
		if err := CreateClusterRoleBinding(ctx, client); err != nil {
			logger.WithError(err).Error("can not create cluster-manager clusterrole")

			return reconcile.Result{}, err
		}

		// @step: check if the kore cluster manager deployment exists
		found, err = DeploymentExists(ctx, client)
		if err != nil {
			logger.WithError(err).Error("trying to check for cluster-manager depoloyment")

			return reconcile.Result{}, err
		}
		if !found {
			c, err := MakeTemplate(clustermanDeployment, params)
			if err != nil {
				logger.WithError(err).Error("trying to render the cluster-manager deployment")

				return reconcile.Result{}, err
			}
			deployment := &apps.Deployment{}
			if err := DecodeInTo(c, deployment); err != nil {
				logger.WithError(err).Error("trying to decode the deployment")

				return reconcile.Result{}, err
			}
			if _, err := kubernetes.CreateOrUpdate(ctx, client, deployment); err != nil {
				logger.WithError(err).Error("trying to create the cluster manager deployment")

				return reconcile.Result{}, err
			}
		}
		logger.Debug("waiting for kore cluster manager deployment status to appear")

		nctx, cancel := context.WithTimeout(ctx, 4*time.Minute)
		defer cancel()

		logger.Info("waiting for kore cluster manager to complete")

		// @step: wait for the bootstrap job to complete
		if err := WaitOnStatus(nctx, client); err != nil {
			logger.WithError(err).Error("failed waiting for kore cluster manager status to complete")

			return reconcile.Result{}, err
		}

		logger.Info("kube api ready for new cluster")

		// @step: else we can set the job as complete
		cluster.Status.Components.SetCondition(corev1.Component{
			Name:    "applications",
			Status:  corev1.SuccessStatus,
			Message: "in cluster applications installed",
		})

		return reconcile.Result{}, nil
	}()
	if err != nil {
		cluster.Status.Components.SetCondition(corev1.Component{
			Name:    "applications",
			Status:  corev1.FailureStatus,
			Message: "attempting to install the applications",
			Detail:  err.Error(),
		})
	}

	if err := t.mgr.GetClient().Status().Patch(ctx, cluster, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("trying to update the cluster resource status")

		return reconcile.Result{}, err
	}

	return result, err
}

// GetClusterConfiguration is responsible for generate the parameters for the cluster
func (t bsCtrl) GetClusterConfiguration(ctx context.Context, cluster *clusterv1.Kubernetes, provider string) (Parameters, error) {
	// @step: retrieve the authentication endpoints
	var authURL, tokenURL string

	if t.Config().DiscoveryURL != "" {
		discovery, err := openid.New(openid.Config{
			DiscoveryURL: t.Config().DiscoveryURL,
			ClientID:     t.Config().ClientID,
		})
		if err != nil {
			log.WithError(err).Error("trying to create discovery client")

			return Parameters{}, err
		}

		if err := discovery.RunWithSync(ctx); err != nil {
			log.WithError(err).Error("trying to retrieve the discovery details")

			return Parameters{}, err
		}
		authURL = discovery.Provider().Endpoint().AuthURL
		tokenURL = discovery.Provider().Endpoint().TokenURL
	}

	params := Parameters{
		KoreImage: KoreImage,
		Kiali: KialiOptions{
			Password: utils.Random(12),
		},
		Domain: cluster.Spec.Domain,
		Grafana: GrafanaOptions{
			AuthURL:      authURL,
			ClientID:     t.Config().ClientID,
			ClientSecret: t.Config().ClientSecret,
			Hostname:     "grafana." + cluster.Name + "." + cluster.Namespace + "." + cluster.Spec.Domain,
			Password:     utils.Random(24),
			TokenURL:     tokenURL,
			Database: DatabaseOptions{
				Name:     "grafana",
				Password: utils.Random(12),
			},
		},
		Provider:     provider,
		StorageClass: "default",
	}
	if t.Config().DEX.EnabledDex {
		params.Grafana.AuthURL = authURL
		params.Grafana.TokenURL = tokenURL
		//params.Grafana.UserInfoURL = t.Config().DiscoveryURL + "/userinfo"
	}

	switch provider {
	case "gke":
		params.StorageClass = "standard"
	}
	return params, nil
}
