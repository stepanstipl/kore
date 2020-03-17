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

package clusterconfig

import (
	"context"
	"time"

	clusterv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	core "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconcile ensure the cluste has it's configuration
func (t ccCtrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to reconcile kubernetes cluster configuration")

	ctx := context.Background()

	// @step: retrieve the resource from the api
	cluster := &clusterv1.Kubernetes{}
	if err := t.mgr.GetClient().Get(ctx, request.NamespacedName, cluster); err != nil {
		if !kerrors.IsNotFound(err) {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	// @step: ensure the cluster has a configuration
	logger.Debug("checking for the cluster credentials secret")

	credentials, err := controllers.GetClusterCredentialsSecret(ctx,
		t.mgr.GetClient(),
		cluster.Namespace,
		cluster.Name)

	if err != nil {
		if !kerrors.IsNotFound(err) {
			logger.WithError(err).Error("trying retrieve the cluster credentials")

			return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
		}
		logger.Debug("no credentials secret found, perhaps not ready yet")

		return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	// @step: create a kubernetes client for the cluster
	client, err := kubernetes.NewRuntimeClientFromConfigSecret(credentials)
	if err != nil {
		logger.WithError(err).Error("trying to create a cluster client")

		return reconcile.Result{}, err
	}

	// @step: check if the api is available
	kc, err := kubernetes.NewClientFromConfigSecret(credentials)
	if err != nil {
		logger.WithError(err).Error("trying to create kubernetes client for cluster")

		return reconcile.Result{RequeueAfter: 2 * time.Minute}, nil
	}

	logger.Debug("checking if the kubernetes api is available yet, else waiting")

	// @step wait for the api to become available
	if err := kubernetes.WaitOnKubeAPI(ctx, kc, 5*time.Second, 20*time.Second); err != nil {
		logger.Debug("kubernetes api for cluster have't come up yet, forging into background")

		return reconcile.Result{RequeueAfter: 2 * time.Minute}, nil
	}

	logger.Debug("checking if the kore namespace exists in the remote cluster")

	// @step: ensure the namespace is there
	for _, namespace := range []string{kore.HubNamespace} {
		if err := client.Create(context.Background(), &core.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
				Labels: map[string]string{
					kore.Label("owned"): "true",
				},
			},
		}); err != nil {
			logger.WithField("namespace", namespace).Debug("kore namespace already exists, skipping")

			if !kerrors.IsAlreadyExists(err) {
				logger.WithError(err).Error("trying to create the kore namespace")
			}
		}
	}

	logger.Debug("creating the configuration secret in the remote cluster")

	// @step: ensure there is a client certificate for us
	secretName := "kore-config"
	found, err := kubernetes.HasSecret(ctx, client, kore.HubNamespace, secretName)
	if err != nil {
		logger.WithError(err).Error("trying to check for kore configuration secret")

		return reconcile.Result{}, err
	}
	// @TODO we need to check if the secret exists, then check the client certificate
	// and if near expiration we need to rotate it.
	if found {
		logger.Debug("skipping kore configuration as cluster already has configuration")

		return reconcile.Result{}, nil
	}

	config := &core.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kore-config",
			Namespace: kore.HubNamespace,
		},
		Data: map[string][]byte{
			"api-url":        []byte(t.Config().PublicAPIURL),
			"idp-server-url": []byte(t.Config().IDPServerURL),
			"domain":         []byte(cluster.Spec.Domain),
			"hub-url":        []byte(t.Config().PublicHubURL),
		},
	}
	logger.Debug("adding the client configuration to the cluster")

	// @step: create a client certificate for the cluster to call back with
	if t.Config().HasCertificateAuthority() {
		cert, key, err := t.SignedClientCertificate(cluster.Name, cluster.Namespace)
		if err != nil {
			logger.WithError(err).Error("generating a client certificate for cluster")

			return reconcile.Result{}, err
		}

		config.Data["tls.crt"] = []byte(string(cert))
		config.Data["tls.key"] = []byte(string(key))
	}

	if _, err := kubernetes.CreateOrUpdateSecret(context.Background(), client, config); err != nil {
		logger.WithError(err).Error("trying to place the cluster configutation")

		return reconcile.Result{}, err
	}

	logger.Debug("sucessfully added the cluster client configuration to cluster")

	return reconcile.Result{}, nil
}
