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

package gke

import (
	"context"
	"time"

	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	kube "github.com/appvia/kore/pkg/utils/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	psp "k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8s "k8s.io/client-go/kubernetes"
)

type bootImpl struct {
	// credentials are the gke credentials
	credentials *credentials
	// cluster is the gke cluster
	cluster *gke.GKE
	// client is the k8s client
	client k8s.Interface
}

// newBootstrapClient returns a bootstrap client
func newBootstrapClient(cluster *gke.GKE, credentials *credentials) (*bootImpl, error) {
	// @step: retrieve the credentials for the cluster
	client, err := kube.NewGKEClient(
		credentials.key,
		cluster.Status.Endpoint,
	)
	if err != nil {
		log.WithError(err).Error("trying to create gke kubernetes client from credentials")

		return nil, err
	}

	return &bootImpl{
		credentials: credentials,
		cluster:     cluster,
		client:      client,
	}, nil
}

func (p *bootImpl) GetCluster() *controllers.BootCluster {
	return &controllers.BootCluster{
		Endpoint:  p.cluster.Status.Endpoint,
		Name:      p.cluster.Name,
		Namespace: p.cluster.Namespace,
		Client:    p.client,
	}
}

func (p *bootImpl) GetLogger() *log.Entry {
	cluster := p.GetCluster()
	logger := log.WithFields(log.Fields{
		"endpoint":  cluster.Endpoint,
		"name":      cluster.Name,
		"namespace": cluster.Namespace,
		"bootstrap": "gke",
	})
	return logger
}

// Cloud specific bootstrap for gke cluster
func (p *bootImpl) Run(ctx context.Context, client client.Client) error {
	logger := p.GetLogger()
	logger.Info("creating the pod security policies")

	if err := p.DeployPodSecurityPolicies(ctx, p.client); err != nil {
		logger.WithError(err).Error("failed to create the gke psp bindings")

		return err
	}
	return nil
}

func (p *bootImpl) GetClusterObj() runtime.Object {
	return p.cluster
}

// DeployPodSecurityPolicies is responsible deploying the PSP to te gke cluster
func (p *bootImpl) DeployPodSecurityPolicies(ctx context.Context, client k8s.Interface) error {
	isFalse := false

	psp := psp.PodSecurityPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kore.default",
			Annotations: map[string]string{
				"apparmor.security.beta.kubernetes.io/allowedProfileNames": "runtime/default",
				"apparmor.security.beta.kubernetes.io/defaultProfileName":  "runtime/default",
				"seccomp.security.alpha.kubernetes.io/allowedProfileNames": "runtime/default,docker/default",
				"seccomp.security.alpha.kubernetes.io/defaultProfileName":  "docker/default",
			},
		},
		Spec: psp.PodSecurityPolicySpec{
			AllowedCapabilities: []corev1.Capability{
				"AUDIT_WRITE", "CHOWN", "DAC_OVERRIDE", "FOWNER",
				"FSETID", "KILL", "MKNOD", "NET_BIND_SERVICE",
				"NET_RAW", "SETFCAP", "SETGID", "SETPCAP",
				"SETUID", "SYS_CHROOT",
			},
			AllowPrivilegeEscalation: &isFalse,
			FSGroup: psp.FSGroupStrategyOptions{
				Rule: psp.FSGroupStrategyRunAsAny,
			},
			RunAsUser: psp.RunAsUserStrategyOptions{
				Rule: psp.RunAsUserStrategyRunAsAny,
			},
			SELinux: psp.SELinuxStrategyOptions{
				Rule: psp.SELinuxStrategyRunAsAny,
			},
			SupplementalGroups: psp.SupplementalGroupsStrategyOptions{
				Rule: psp.SupplementalGroupsStrategyRunAsAny,
			},
			Volumes: []psp.FSType{
				"awsElasticBlockStore",
				"azureDisk",
				"azureFile",
				"cephFS",
				"configMap",
				"downwardAPI",
				"emptyDir",
				"gcePersistentDisk",
				"persistentVolumeClaim",
				"projected",
				"secret",
			},
		},
	}
	if err := p.CreateClusterPodSecurityPolicy(&psp); err != nil {
		return err
	}

	role := rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default:psp",
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups:     []string{"policy"},
				ResourceNames: []string{"kore.default"},
				Resources:     []string{"podsecuritypolicies"},
				Verbs:         []string{"use"},
			},
		},
	}
	if err := p.CreateClusterRole(&role); err != nil {
		return err
	}

	return p.CreateClusterRoleBinding(&rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default:psp",
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "ClusterRole",
			Name:     "default:psp",
		},
		Subjects: []rbacv1.Subject{
			{
				APIGroup: rbacv1.GroupName,
				Kind:     "Group",
				Name:     "system:authenticated",
			},
			{
				APIGroup: rbacv1.GroupName,
				Kind:     "Group",
				Name:     "system:serviceaccounts",
			},
		},
	})
}

// CreateClusterPodSecurityPolicy creates a psp in the cluster
func (p *bootImpl) CreateClusterPodSecurityPolicy(policy *psp.PodSecurityPolicy) error {
	if _, err := p.client.ExtensionsV1beta1().PodSecurityPolicies().Create(policy); err != nil {
		if kerrors.IsAlreadyExists(err) {
			return nil
		}

		return err
	}

	return nil
}

// makeClusterRole is responsible creating a cluster role
func (p *bootImpl) CreateClusterRole(role *rbacv1.ClusterRole) error {
	if _, err := p.client.RbacV1().ClusterRoles().Create(role); err != nil {
		if kerrors.IsAlreadyExists(err) {
			return nil
		}

		return err
	}

	return nil
}

// CreateClusterRoleBinding is responsible for cluster role binding
func (p *bootImpl) CreateClusterRoleBinding(binding *rbacv1.ClusterRoleBinding) error {
	if _, err := p.client.RbacV1().ClusterRoleBindings().Create(binding); err != nil {
		if kerrors.IsAlreadyExists(err) {
			return nil
		}

		return err
	}

	return nil
}

// CreateSysadminCredential is responsible for creating admin creds
func (p *bootImpl) CreateSysadminCredential() (*corev1.Secret, error) {
	// @step: check if the service account already exists
	name := "kore-admin"
	namespace := "kube-system"

	_, err := p.client.CoreV1().ServiceAccounts(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if !kerrors.IsNotFound(err) {
			return nil, err
		}

		// @step: create the service account
		if _, err := p.client.CoreV1().ServiceAccounts(namespace).Create(&corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
				Labels: map[string]string{
					"kore.appvia.io/owner": "true",
				},
			}}); err != nil {

			return nil, err
		}
	}

	// @step: create the binding for the cluster admin role
	if _, err := p.client.RbacV1().ClusterRoleBindings().Get(name, metav1.GetOptions{}); err != nil {
		if !kerrors.IsNotFound(err) {
			return nil, err
		}
		binding := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: rbacv1.GroupName,
				Kind:     "ClusterRole",
				Name:     "cluster-admin",
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      rbacv1.ServiceAccountKind,
					Namespace: namespace,
					Name:      name,
				},
			},
		}

		if _, err := p.client.RbacV1().ClusterRoleBindings().Create(binding); err != nil {
			return nil, err
		}
	}

	return kube.WaitForServiceAccountToken(p.client, namespace, name, time.Duration(1*time.Minute))
}
