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
	"bytes"
	"context"
	"fmt"
	"text/template"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/schema"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	"github.com/Masterminds/sprig"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EnsureAPIService is responsible for ensuring the api proxy is provisioned
func (a k8sCtrl) EnsureAPIService(
	ctx context.Context,
	cc client.Client,
	cluster *clustersv1.Kubernetes) error {

	logger := log.WithFields(log.Fields{
		"cluster":   cluster.Name,
		"namespace": cluster.Namespace,
		"team":      cluster.Namespace,
	})
	logger.Debug("ensuring the kube api service is provisioned")

	// @step: ensure the namespace is there
	if err := kubernetes.EnsureNamespace(ctx, cc, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   kore.HubNamespace,
			Labels: map[string]string{kore.Label("owned"): "true"},
		},
	}); err != nil {
		logger.WithError(err).Error("trying to provision the namespace")

		return err
	}

	// @step: build the context for the template
	params := map[string]interface{}{
		"AllowedIPs":  cluster.Spec.AuthProxyAllowedIPs,
		"CACert":      a.CertificateAuthority(),
		"ClientID":    a.Config().IDPClientID,
		"Deployment":  "oidc-proxy",
		"Domain":      cluster.Spec.Domain,
		"Hostname":    a.APIHostname(cluster),
		"Image":       a.Config().AuthProxyImage,
		"MaxReplicas": 10,
		"Name":        cluster.Name,
		"Namespace":   kore.HubNamespace,
		"OpenID":      a.Config().HasOpenID(),
		"Provider":    cluster.Spec.Provider.Kind,
		"Replicas":    2,
		"ServerURL":   a.Config().IDPServerURL,
		"Team":        cluster.Namespace,
		"TLSKey":      "",
		"TLSCert":     "",
		"UserClaims":  a.Config().IDPUserClaims,
	}
	if cluster.Spec.AuthProxyImage != "" {
		params["Image"] = cluster.Spec.AuthProxyImage
	}

	// @step: does the tls secret already exist
	found, err := kubernetes.HasSecret(ctx, cc, kore.HubNamespace, "tls")
	if err != nil {
		logger.WithError(err).Error("trying to check for tls secret")

		return err
	}
	if !found {
		hosts := []string{
			a.APIHostname(cluster),
			"localhost",
			"127.0.0.1",
		}
		cert, key, err := a.SignedServerCertificate(hosts, (24*365*time.Hour)*10)
		if err != nil {
			logger.WithError(err).Error("trying to generate the auth proxy certificate")

			return err
		}
		params["TLSCert"] = cert
		params["TLSKey"] = key
	}

	// @step: generate the content
	tmpl, err := template.New("main").Funcs(sprig.TxtFuncMap()).Parse(AuthProxyDeployment)
	if err != nil {
		logger.WithError(err).Error("trying to generate the auth proxy template")

		return err
	}
	generated := &bytes.Buffer{}
	if err := tmpl.Execute(generated, &params); err != nil {
		logger.WithError(err).Error("trying to generate deployment")

		return err
	}
	resources, err := utils.YAMLDocuments(generated)
	if err != nil {
		logger.WithError(err).Error("trying to split the documents")

		return err
	}
	for _, x := range resources {
		object, err := schema.DecodeYAML([]byte(x))
		if err != nil {
			logger.WithError(err).Error("trying to decode the document")

			return err
		}
		kind := object.GetObjectKind().GroupVersionKind().Kind
		ignore := []string{
			"ClusterRole",
			"ClusterRoleBinding",
			"Service",
			"ServiceAccount",
		}
		if utils.Contains(kind, ignore) {
			if found, err := kubernetes.CheckIfExists(ctx, cc, object); err != nil {
				return err
			} else if found {
				continue
			}
		}

		if _, err := kubernetes.CreateOrUpdate(ctx, cc, object); err != nil {
			logger.WithError(err).Error("trying to create or update the resource")

			return err
		}
	}

	// @step: wait for the service to have an endpoint
	timeout, cancel := context.WithTimeout(ctx, 240*time.Second)
	defer cancel()

	// @TODO we need to move this is a non-blocking method
	endpoint, err := kubernetes.WaitForServiceEndpoint(timeout, cc, kore.HubNamespace, "proxy")
	if err != nil {
		logger.WithError(err).Error("trying to wait for service endpont")

		return err
	}
	cluster.Status.Endpoint = endpoint

	return nil
}

// APIHostname is the hostname of the kube api proxy
func (a k8sCtrl) APIHostname(cluster *clustersv1.Kubernetes) string {
	return fmt.Sprintf("api.%s.%s.%s",
		cluster.Name,
		cluster.Namespace,
		cluster.Spec.Domain,
	)
}
