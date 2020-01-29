/*
 * Copyright (C) 2019 Rohith Jayawardene <gambol99@gmail.com>
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

	aws "github.com/appvia/kore/pkg/apis/aws/v1alpha1"
	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	clusterv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	log "github.com/sirupsen/logrus"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	// convertor is the unstructured converter client
	converter = runtime.DefaultUnstructuredConverter
)

const (
	finalizerName = "bootstrap.compute.hub.appvia.io"
	// jobNamespace is the namespace the job runs in
	jobNamespace = "kube-system"
	// jobName is the name of the job
	jobName = "bootstrap"
	// jobOLMConfig is the configuration for the olm config
	jobOLMConfig = "bootstrap-olm"
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

	status, found = cluster.Status.Components.GetStatus("applications")
	if !found {
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
	result, err := func() (reconcile.Result, error) {
		// @step: we need to grab the credentials for the cloud provider
		us, err := controllers.GetCloudProviderCredentials(ctx, t.mgr.GetClient(), cluster)
		if err != nil {
			logger.WithError(err).Error("trying to retrieve cloud provider credentials")

			return reconcile.Result{}, err
		}

		var credentials Credentials
		provider := strings.ToLower(cluster.Spec.Provider.Kind)

		switch provider {
		case "gke":
			u := &gke.GKECredentials{}
			if err := converter.FromUnstructured(us.Object, u); err != nil {
				logger.WithError(err).Error("trying to convert object type")

				return reconcile.Result{}, err
			}
			credentials.GKE.Account = u.Spec.Account
		case "eks":
			u := &aws.AWSCredentials{}
			if err := converter.FromUnstructured(us.Object, u); err != nil {
				logger.WithError(err).Error("trying to convert object type")

				return reconcile.Result{}, err
			}
			credentials.AWS = AWSCredentials{
				AccessKey: u.Spec.AccessKeyID,
				AccountID: u.Spec.AccountID,
				Region:    "eu-west-2",
				SecretKey: u.Spec.SecretAccessKey,
			}
		default:
			logger.WithField("provider", provider).Warn("unknown cloud provider")
		}

		// @step: build a client from the cluster secret
		client, err := controllers.CreateClientFromSecret(ctx, t.mgr.GetClient(), cluster.Name, cluster.Namespace)
		if err != nil {
			logger.WithError(err).Error("trying to retrieve the credentials")

			return reconcile.Result{}, err
		}

		// @step: check if the job configmap exists
		found, err := JobConfigExists(ctx, client)
		if err != nil {
			logger.WithError(err).Error("failed to check for configuration configmap")

			// @step: we need to fix this as its a loop
			return reconcile.Result{RequeueAfter: 2 * time.Minute}, nil
		}
		if !found {
			c, err := MakeTemplate(BootstrapJobConfigmap, map[string]string{})
			if err != nil {
				logger.WithError(err).Error("failed to generate bootstrap template")

				return reconcile.Result{}, err
			}

			cm := &core.ConfigMap{}
			if err := DecodeInTo(c, cm); err != nil {
				logger.WithError(err).Error("failed to decode bootstrap configuration in configmap")

				return reconcile.Result{}, err
			}

			if _, err := kubernetes.CreateOrUpdate(ctx, client, cm); err != nil {
				logger.WithError(err).Error("failed to create the bootstrap configuration configmap")

				return reconcile.Result{}, err
			}
		}

		// @step: build the parameters for the job
		params, err := t.GetClusterConfiguration(ctx, cluster, strings.ToLower(cluster.Spec.Provider.Kind))
		if err != nil {
			logger.WithError(err).Error("failed to generate the parameters")

			return reconcile.Result{}, err
		}
		params.Credentials = credentials

		// @step: we check if the job configuration is there already and if not we make it
		found, err = JobOLMConfigExists(ctx, client)
		if err != nil {
			logger.WithError(err).Error("failed to check for olm job configuration")

			return reconcile.Result{}, err

		}
		if !found {
			c, err := MakeTemplate(BootstrapJobOLMConfig, params)
			if err != nil {
				logger.WithError(err).Error("failed to generate the bootstrap olm template")

				return reconcile.Result{}, err
			}
			cm := &core.ConfigMap{}
			if err := DecodeInTo(c, cm); err != nil {
				logger.WithError(err).Error("trying to decode olm bootstrap configuration in configmap")

				return reconcile.Result{}, err
			}

			if _, err := kubernetes.CreateOrUpdate(ctx, client, cm); err != nil {
				logger.WithError(err).Error("trying to create the olm bootstrap configuration configmap")

				return reconcile.Result{}, err
			}
		}

		// @step: ensure the bootstrap job is there
		found, err = JobExists(ctx, client)
		if err != nil {
			logger.WithError(err).Error("trying to check for bootstrap job")

			return reconcile.Result{}, err
		}
		if !found {
			c, err := MakeTemplate(BootstrapJobTemplate, params)
			if err != nil {
				logger.WithError(err).Error("trying to render the bootstrap job")

				return reconcile.Result{}, err
			}
			job := &batch.Job{}
			if err := DecodeInTo(c, job); err != nil {
				logger.WithError(err).Error("trying to decode the job")

				return reconcile.Result{}, err
			}
			if _, err := kubernetes.CreateOrUpdate(ctx, client, job); err != nil {
				logger.WithError(err).Error("trying to create the bootstrap job")

				return reconcile.Result{}, err
			}
		}
		logger.Debug("waiting for bootstrap job to finish")

		nctx, cancel := context.WithTimeout(ctx, 20*time.Minute)
		defer cancel()

		logger.Info("waiting for bootstrap has completed")

		// @step: wait for the bootstrap job to complete
		if err := WaitOnJob(nctx, client); err != nil {
			logger.WithError(err).Error("failed waiting for bootstrap to complete")

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
	params := Parameters{
		BootImage: "quay.io/appvia/hub-bootstrap:v0.2.0",
		Catalog:   CatalogOptions{Image: "v0.0.2"},
		Kiali: KialiOptions{
			Password: utils.Random(12),
		},
		Domain: cluster.Spec.Domain,
		Grafana: GrafanaOptions{
			AuthURL:      t.Config().DiscoveryURL + "/protocol/openid-connect/auth",
			ClientID:     t.Config().ClientID,
			ClientSecret: t.Config().ClientSecret,
			Hostname:     "grafana." + cluster.Name + "." + cluster.Namespace + "." + cluster.Spec.Domain,
			Password:     utils.Random(12),
			TokenURL:     t.Config().DiscoveryURL + "/openid-connect/token",
			Database: DatabaseOptions{
				Name:     "grafana",
				Password: utils.Random(12),
			},
		},
		Provider:   provider,
		OLMVersion: "0.11.0",
		Namespaces: []NamespaceOptions{
			{Name: "kube-dns"},
			{Name: "grafana"},
			{Name: "logging"},
			{Name: "prometheus"},
		},
		StorageClass: "default",
	}
	switch provider {
	case "gke":
		params.StorageClass = "standard"
	}

	// @step: ensure we have the operators
	params.Operators = []OperatorOptions{
		{
			Package:   "prometheus",
			Channel:   "beta",
			Label:     "k8s-app=prometheus-operator",
			Namespace: "prometheus",
		},
		{
			Package:   "grafana-operator",
			Channel:   "alpha",
			Label:     "app=grafana-operator",
			Namespace: "grafana",
		},
		{
			Package:   "loki-operator",
			Channel:   "stable",
			Label:     "name=loki-operator",
			Namespace: "logging",
		},
		{
			Package:   "metrics-operator",
			Channel:   "stable",
			Label:     "name=metrics-operator",
			Namespace: "prometheus",
		},
		{
			Package:   "mariadb-operator",
			Channel:   "stable",
			Label:     "name=mariadb-operator",
			Namespace: "grafana",
		},
		{
			Package:   "external-dns-operator",
			Channel:   "stable",
			Label:     "name=external-dns-operator",
			Namespace: "kube-dns",
		},
	}

	return params, nil
}
