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
	KubeProxyNamespace = "kube-proxy"
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
		Hostname string
		Domain   string
		Image    string
		Name     string
		Team     string
	}
	parameters.Team = cluster.Namespace
	parameters.Name = cluster.Name
	parameters.Domain = cluster.Spec.Domain
	parameters.Hostname = a.APIHostname(cluster)
	parameters.Image = "quay.io/jetstack/kube-oidc-proxy:v0.2.0"
	if cluster.Spec.ProxyImage != "" {
		parameters.Image = cluster.Spec.ProxyImage
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
		cert, key, err := a.SignedServerCertificate(
			[]string{
				parameters.Hostname,
				"localhost",
				"127.0.0.1",
			},
			24*365*time.Hour,
		)
		if err != nil {
			logger.WithError(err).Error("generate the server certificate")

			return err
		}

		if _, err := kubernetes.CreateOrUpdateSecret(ctx, cc, &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      KubeProxyTLSSecret,
				Namespace: KubeProxyNamespace,
			},
			Data: map[string][]byte{
				"tls.crt": cert,
				"tls.key": key,
			},
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

	// @step: ensure the kubernete service
	if err := cc.Get(context.Background(), types.NamespacedName{
		Namespace: KubeProxyNamespace,
		Name:      "proxy",
	}, &v1.Service{}); err != nil {
		if !kerrors.IsNotFound(err) {
			logger.WithError(err).Error("trying to create the service for proxy")

			return err
		}

		if _, err := kubernetes.CreateOrUpdate(ctx, cc, &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "proxy",
				Namespace: KubeProxyNamespace,
				Annotations: map[string]string{
					"external-dns.alpha.kubernetes.io/hostname": parameters.Hostname,
				},
			},
			Spec: v1.ServiceSpec{
				Type: v1.ServiceTypeLoadBalancer,
				Ports: []v1.ServicePort{
					{
						Port:       443,
						TargetPort: intstr.FromInt(443),
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

	// @step: ensure the deployment is there
	if _, err := kubernetes.CreateOrUpdate(ctx, cc, &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: KubeProxyNamespace,
			Labels: map[string]string{
				"name": name,
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
								{ContainerPort: 443},
								{ContainerPort: 8080},
							},
							ReadinessProbe: &v1.Probe{
								Handler: v1.Handler{
									HTTPGet: &v1.HTTPGetAction{
										Path: "/ready",
										Port: intstr.FromInt(8080),
									},
								},
								InitialDelaySeconds: 5,
								PeriodSeconds:       10,
							},
							Command: []string{"kube-oidc-proxy"},
							Args: []string{
								"--oidc-client-id=" + a.Config().ClientID,
								"--oidc-issuer-url=" + a.Config().DiscoveryURL,
								"--oidc-username-claim=" + a.Config().UserClaims[0],
								"--secure-port=443",
								"--tls-cert-file=/tls/tls.crt",
								"--tls-private-key-file=/tls/tls.key",
							},
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
	timeout, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	endpoint, err := kubernetes.WaitForServiceEndpoint(timeout, cc, KubeProxyNamespace, "proxy")
	if err != nil {
		logger.WithError(err).Error("trying to wait for service endpont")

		return err
	}
	cluster.Status.Endpoint = endpoint

	return nil
}
