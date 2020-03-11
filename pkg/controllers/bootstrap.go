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

package controllers

import (
	"context"
	"time"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	kube "github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8s "k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BootCluster provides exposes generic cluster bootstrap data
type BootCluster struct {
	Name      string
	Namespace string
	Endpoint  string
	Client    k8s.Interface
}

// Bootstrap provides the shared interface for specific providers to interact with
type Bootstrap interface {
	Run(context.Context, client.Client) error
	GetCluster() *BootCluster
	GetClusterObj() runtime.Object
	GetLogger() *log.Entry
}

// bootstrap represents a remote cluster to bootstrap
type bootstrapImpl struct {
	// provider allows acces to the specific cloud entry point(s)
	provider Bootstrap
}

var _ Bootstrap = &bootstrapImpl{}

// NewBootstrap wrapps cloud specific calls
// TODO: move into cloud.go and create THE cloud interface
func NewBootstrap(p Bootstrap) Bootstrap {
	b := &bootstrapImpl{
		provider: p,
	}
	return b
}

// GetCluster will get the cloud providers cluster details and client
func (b *bootstrapImpl) GetCluster() *BootCluster {
	return b.provider.GetCluster()
}

func (b *bootstrapImpl) GetClusterObj() runtime.Object {
	return b.provider.GetClusterObj()
}

// GetLogger gets the configured logger for bootstrapping
func (b *bootstrapImpl) GetLogger() *log.Entry {
	cluster := b.GetCluster()
	logger := log.WithFields(log.Fields{
		"endpoint":  cluster.Endpoint,
		"name":      cluster.Name,
		"namespace": cluster.Namespace,
	})
	return logger
}

// Run will start the bootstrap of a remote cluser and update access in kore cluster
func (b *bootstrapImpl) Run(ctx context.Context, cc client.Client) error {
	cluster := b.GetCluster()
	logger := b.GetLogger()
	logger.Infof("waiting for the kube-apiserver to become available at %s", cluster.Endpoint)

	// @step: wait for the kubernetes api to become available
	timeout := 5 * time.Minute

	if err := kube.WaitOnKubeAPI(ctx, cluster.Client, time.Duration(10)*time.Second, timeout); err != nil {
		logger.WithError(err).Error("timeout waiting for kube-apiserver to become available")

		return err
	}
	logger.Debug("cluster kubeapi is available now, continuing bootstrapping")

	logger.Info("creating cloud specifics...")

	if err := b.provider.Run(ctx, cc); err != nil {
		logger.WithError(err).Error("failed to create the cloud spcific entries")

		return err
	}

	logger.Info("creating the kore-admin service account for cluster")
	// @step: create or retrieve the kore-sysadmin secret token
	creds, err := CreateSysadminCredential(cluster.Client)
	if err != nil {
		logger.WithError(err).Error("creating kore admin service account")

		return err
	}
	secret := NewEmptySecret().
		Description("Kubernetes Cluster credentials for " + cluster.Name).
		Name(cluster.Name).
		Namespace(cluster.Namespace).
		Type(configv1.KubernetesSecret).
		Values(map[string]string{
			"ca.crt":   string(creds.Data["ca.crt"]),
			"endpoint": cluster.Endpoint,
			"token":    string(creds.Data["token"]),
		}).Secret()

	if err := CreateManagedSecret(ctx, b.GetClusterObj(), cc, secret.Encode()); err != nil {
		logger.WithError(err).Error("trying to create sysadmin secret")

		return err
	}

	logger.Info("successfully bootstrapped the cluster")

	return nil
}

// CreateSysadminCredential will create a service account in remote cluster
func CreateSysadminCredential(rc k8s.Interface) (*corev1.Secret, error) {
	// @step: check if the service account already exists
	name := "kore-admin"
	namespace := "kube-system"

	_, err := rc.CoreV1().ServiceAccounts(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if !kerrors.IsNotFound(err) {
			return nil, err
		}

		// @step: create the service account
		if _, err := rc.CoreV1().ServiceAccounts(namespace).Create(&corev1.ServiceAccount{
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
	if _, err := rc.RbacV1().ClusterRoleBindings().Get(name, metav1.GetOptions{}); err != nil {
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

		if _, err := rc.RbacV1().ClusterRoleBindings().Create(binding); err != nil {
			return nil, err
		}
	}

	return kube.WaitForServiceAccountToken(rc, namespace, name, time.Duration(1*time.Minute))
}
