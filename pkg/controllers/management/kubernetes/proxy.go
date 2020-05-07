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
	"fmt"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// KubeProxyNamespace is the namespace the proxy
	KubeProxyNamespace = "kore"
	// KubeProxySecret is the secret for the proxy config
	KubeProxySecret = "config"
	// KubeProxyTLSSecret is the secret containing the tls
	KubeProxyTLSSecret = "tls"
)

// APIHostname is the hostname of the kube api proxy
func (a k8sCtrl) APIHostname(cluster *clustersv1.Kubernetes) string {
	return fmt.Sprintf("api.%s.%s.%s",
		cluster.Name,
		cluster.Namespace,
		cluster.Spec.Domain,
	)
}

// IsProxyProtocolRequired checks if we should enabled proxy protocol
func IsProxyProtocolRequired(cluster *clustersv1.Kubernetes) bool {
	return cluster.Spec.Provider.Kind == "EKS"
}

// EnsureAPIService is responsible for ensuring the api proxy is provisioned
func (a k8sCtrl) EnsureAPIService(ctx context.Context, cc client.Client, cluster *clustersv1.Kubernetes) error {
	logger := log.WithFields(log.Fields{
		"cluster":   cluster.Name,
		"namespace": cluster.Namespace,
		"team":      cluster.Namespace,
	})
	logger.Debug("ensuring the kube api service is provisioned")

	// @step: ensure the namespace is there
	if err := kubernetes.EnsureNamespace(ctx, cc, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: KubeProxyNamespace,
			Labels: map[string]string{
				kore.Label("owned"): "true",
			},
		},
	}); err != nil {
		logger.WithError(err).Error("trying to provision the namespace")

		return err
	}

	// @step: build the parameters
	var parameters struct {
		Domain        string
		Hostname      string
		Image         string
		Name          string
		ProxyProtocol bool
		Team          string
	}
	parameters.Team = cluster.Namespace
	parameters.Name = cluster.Name
	parameters.Domain = cluster.Spec.Domain
	parameters.Hostname = a.APIHostname(cluster)
	parameters.Image = a.Config().AuthProxyImage
	parameters.ProxyProtocol = IsProxyProtocolRequired(cluster)
	if cluster.Spec.AuthProxyImage != "" {
		parameters.Image = cluster.Spec.AuthProxyImage
	}

	// @step: ensure the tls secret is provisioned and configured
	found, err := kubernetes.HasSecret(ctx, cc, cluster.Namespace, KubeProxyTLSSecret)
	if err != nil {
		logger.WithError(err).Error("trying to check for proxy tls secret")

		return err
	}

	// @TODO need to add a check for expiration?
	if !found {
		// @step: we need to generate a server certificate for the api
		cert, key, err := a.SignedServerCertificate([]string{parameters.Hostname, "localhost", "127.0.0.1"}, 24*365*time.Hour)
		if err != nil {
			logger.WithError(err).Error("generate the server certificate")

			return err
		}

		if _, err := kubernetes.CreateOrUpdateSecret(ctx, cc, &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      KubeProxyTLSSecret,
				Namespace: KubeProxyNamespace,
			},
			Data: map[string][]byte{"tls.crt": cert, "tls.key": key},
		}); err != nil {
			logger.WithError(err).Error("trying to create the tls secret")

			return err
		}
	}

	// @step: ensure the service account
	if _, err := kubernetes.CreateOrUpdateServiceAccount(ctx, cc, &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "proxy",
			Namespace: KubeProxyNamespace,
			Labels: map[string]string{
				kore.Label("owner"): "true",
			},
		},
	}); err != nil {
		logger.WithError(err).Error("trying to create the service account")

		return err
	}

	// @step: ensure the cluster role and binding exist
	if _, err := kubernetes.CreateOrUpdateManagedClusterRole(ctx, cc, &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kore:oidc:proxy",
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"users", "groups", "serviceaccount"},
				Verbs:     []string{"impersonate"},
			},
			{
				APIGroups: []string{"authentication.k8s.io"},
				Resources: []string{"userextras/scopes", "tokenreviews"},
				Verbs:     []string{"create", "impersonate"},
			},
		},
	}); err != nil {
		logger.WithError(err).Error("trying to ensyre the cluster role for proxy")

		return err
	}
	if _, err := kubernetes.CreateOrUpdateManagedClusterRoleBinding(ctx, cc, &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kore:oidc:proxy",
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.SchemeGroupVersion.Group,
			Kind:     "ClusterRole",
			Name:     "kore:oidc:proxy",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Namespace: KubeProxyNamespace,
				Name:      "proxy",
			},
		},
	}); err != nil {
		logger.WithError(err).Error("trying to ensyre the cluster role binding for proxy")

		return err
	}

	name := "oidc-proxy"

	// @step: ensure the kubernetes service
	if err := cc.Get(context.Background(), types.NamespacedName{
		Namespace: KubeProxyNamespace,
		Name:      "proxy",
	}, &v1.Service{}); err != nil {
		if !kerrors.IsNotFound(err) {
			logger.WithError(err).Error("trying to create the service for proxy")

			return err
		}

		annotations := map[string]string{
			"external-dns.alpha.kubernetes.io/hostname": parameters.Hostname,
		}
		if parameters.ProxyProtocol {
			annotations["service.beta.kubernetes.io/aws-load-balancer-proxy-protocol"] = "*"
		}

		if _, err := kubernetes.CreateOrUpdate(ctx, cc, &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "proxy",
				Namespace:   KubeProxyNamespace,
				Annotations: annotations,
			},
			Spec: v1.ServiceSpec{
				Type:                  v1.ServiceTypeLoadBalancer,
				ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeLocal,
				Ports: []v1.ServicePort{
					{
						Port:       443,
						TargetPort: intstr.FromInt(10443),
						Protocol:   v1.ProtocolTCP,
						Name:       "https",
					},
				},
				Selector: map[string]string{
					"name": name,
				},
			},
		}); err != nil {
			logger.WithError(err).Error("trying to create the service for proxy")

			return err
		}
	}

	replicas := int32(2)

	args := []string{
		"--idp-client-id=" + a.Config().IDPClientID,
		"--idp-server-url=" + a.Config().IDPServerURL,
		"--tls-cert=/tls/tls.crt",
		"--tls-key=/tls/tls.key",
	}

	// @step: construct the readiness probe
	readiness := &v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Path:   "/ready",
				Port:   intstr.FromInt(10443),
				Scheme: "HTTPS",
			},
		},
		PeriodSeconds: 10,
	}

	for _, x := range a.Config().IDPUserClaims {
		args = append(args, fmt.Sprintf("--idp-user-claims=%s", x))
	}

	for _, allowedIP := range cluster.Spec.AuthProxyAllowedIPs {
		args = append(args, fmt.Sprintf("--allowed-ips=%s", allowedIP))
	}

	if parameters.ProxyProtocol {
		logger.Debug("enabling proxy protocol readiness check for auth-proxy")
		readiness = &v1.Probe{
			Handler: v1.Handler{
				TCPSocket: &v1.TCPSocketAction{Port: intstr.FromInt(10443)},
			},
			PeriodSeconds: 10,
		}
		args = append(args, "--enable-proxy-protocol="+fmt.Sprintf("%t", parameters.ProxyProtocol))
	}

	// @step: create the readiness probe for the proxy - this has to change to
	// tcp if we are using proxy protocol

	// @step: ensure the deployment is there
	if _, err := kubernetes.CreateOrUpdate(ctx, cc, &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: KubeProxyNamespace,
			Labels: map[string]string{
				"name": name,
			},
			Annotations: map[string]string{
				"prometheus.io/port":   "8080",
				"prometheus.io/scrape": "true",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": name,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": name,
					},
				},
				Spec: v1.PodSpec{
					ServiceAccountName: "proxy",
					Containers: []v1.Container{
						{
							Name:  name,
							Image: parameters.Image,
							Ports: []v1.ContainerPort{
								{ContainerPort: 10443},
								{ContainerPort: 8080},
							},
							ReadinessProbe: readiness,
							Args:           args,
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      "tls",
									MountPath: "/tls",
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: "tls",
							VolumeSource: v1.VolumeSource{
								Secret: &v1.SecretVolumeSource{
									SecretName: "tls",
								},
							},
						},
					},
				},
			},
		},
	}); err != nil {
		logger.WithError(err).Error("trying to create the deployment")

		return err
	}

	// @step: wait for the service to have an endpoint
	timeout, cancel := context.WithTimeout(ctx, 240*time.Second)
	defer cancel()

	endpoint, err := kubernetes.WaitForServiceEndpoint(timeout, cc, KubeProxyNamespace, "proxy")
	if err != nil {
		logger.WithError(err).Error("trying to wait for service endpont")

		return err
	}
	cluster.Status.Endpoint = endpoint

	return nil
}
