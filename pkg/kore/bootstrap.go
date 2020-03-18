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

package kore

import (
	"context"
	"time"

	"github.com/appvia/kore/pkg/kore"
	kube "github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Cluster represents a remote cluster to bootstrap
type Cluster struct {
	Name      string
	Namespace string
	Endpoint  string
	client    k8s.Interface
}

// Bootstrap a remote cluser (rc) and update access in kore cluster (kc)
func Bootstrap(ctx context.Context, cluster Cluster, client client.Client, logger *log.Entry) error {
	logger.Info("creating the kore-admin service account for cluster")
	// @step: create or retrieve the kore-sysadmin secret token
	secret, err := CreateSysadminCredential(cluster.client)
	if err != nil {
		logger.WithError(err).Error("creating kore admin service account")

		return err
	}

	_, err = kube.CreateOrUpdateSecret(ctx, client, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.Name,
			Namespace: cluster.Namespace,
			Labels: map[string]string{
				kore.Label("type"): "kubernetescredentials",
			},
		},
		Data: map[string][]byte{
			"ca.crt":   secret.Data["ca.crt"],
			"endpoint": []byte(cluster.Endpoint),
			"token":    secret.Data["token"],
		},
	})
	if err != nil {
		logger.WithError(err).Error("trying to create sysadmin token")

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
