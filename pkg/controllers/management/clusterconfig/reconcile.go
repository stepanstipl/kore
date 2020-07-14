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
	"bytes"
	"context"
	"time"

	clusterv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconcile ensure the cluste has it's configuration
// @TODO need to convert this controller over to using ensurefunc and the new controller interface
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

	client := t.mgr.GetClient()
	koreURL := "http://kore-apiserver.svc.kore.cluster.local:10080"
	portalURL := "http://kore-portal.svc.kore.cluster.local:3000"

	if request.Namespace != kore.HubAdminTeam {
		koreURL = t.Config().PublicAPIURL
		portalURL = t.Config().PublicHubURL

		credentials, err := controllers.GetClusterCredentialsSecret(ctx,
			t.mgr.GetClient(),
			cluster.Namespace,
			cluster.Name)

		if err != nil {
			if !kerrors.IsNotFound(err) {
				return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
			}

			return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
		}

		// @step: create a kubernetes client for the cluster
		client, err = kubernetes.NewRuntimeClientFromConfigSecret(credentials)
		if err != nil {
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

			return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
		}
	}

	// @step: ensure the namespace is there
	found, err := kubernetes.CheckIfExists(ctx, client, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: kore.HubNamespace},
	})
	if err != nil {
		return reconcile.Result{}, err
	}
	if !found {
		return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
	}

	logger.Debug("creating the configuration secret in the remote cluster")

	config := &core.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kore-config",
			Namespace: kore.HubNamespace,
		},
		Data: map[string][]byte{
			"ca":            []byte(t.CertificateAuthority()),
			"client_id":     []byte(t.Config().IDPClientID),
			"cluster":       []byte(cluster.Name),
			"discovery-url": []byte(t.Config().IDPServerURL),
			"domain":        []byte(cluster.Spec.Domain),
			"kore-url":      []byte(koreURL),
			"portal-url":    []byte(portalURL),
			"provider":      []byte(cluster.Spec.Provider.Kind),
			"team":          []byte(cluster.Namespace),
		},
	}

	// @step: ensure there is a client certificate for us
	current := &core.Secret{}
	current.Namespace = kore.HubNamespace
	current.Name = "kore-config"

	createConfigSecret := false

	found, err = kubernetes.GetIfExists(ctx, client, current)
	if err != nil {
		logger.WithError(err).Error("trying to check for kore configuration secret")

		return reconcile.Result{}, err
	}
	if !found {
		createConfigSecret = true
	}
	// @TODO we need to check if the secret exists, then check the client certificate
	// and if near expiration we need to rotate it.
	if found {
		logger.Debug("skipping kore configuration as cluster already has configuration")

		// @step: check we have all the fields
		asExpected := func() bool {
			for name, expected := range config.Data {
				value, found := current.Data[name]
				if !found || !bytes.Equal(value, expected) {
					return false
				}
			}

			return true
		}()
		if !asExpected {
			createConfigSecret = true
		}
	}

	if createConfigSecret {
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
			logger.WithError(err).Error("trying to place the cluster configuration")

			return reconcile.Result{}, err
		}
		logger.Debug("successfully added the cluster client configuration to cluster")
	}

	// @step: ensure the tls for services - we can remove this once we've sorted our the
	// dns for clusters and we can default to using cert-manager
	found, err = kubernetes.CheckIfNamespaceExists(ctx, client, kore.HubSystem)
	if err != nil {
		logger.WithError(err).Error("trying to check for namespace")

		return reconcile.Result{}, err
	}
	if !found {
		logger.Warn("namespace is not available yet")

		return reconcile.Result{RequeueAfter: 20 * time.Second}, nil
	}

	serverTLS := &v1.Secret{}
	serverTLS.Namespace = kore.HubSystem
	serverTLS.Name = "kore-server-tls"

	if exists, err := kubernetes.CheckIfExists(ctx, client, serverTLS); err != nil {
		logger.WithError(err).Error("trying to check for the server tls certificate")

		return reconcile.Result{}, err
	} else if !exists {
		cert, key, err := t.SignedClientCertificate("server-tls", cluster.Namespace)
		if err != nil {
			logger.WithError(err).Error("generating a server certificate for cluster")

			return reconcile.Result{}, err
		}
		serverTLS.Data = make(map[string][]byte)
		serverTLS.Data["tls.crt"] = []byte(string(cert))
		serverTLS.Data["tls.key"] = []byte(string(key))

		if _, err := kubernetes.CreateOrUpdate(ctx, client, serverTLS); err != nil {
			logger.WithError(err).Error("trying to create the server tls certificate")

			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}
