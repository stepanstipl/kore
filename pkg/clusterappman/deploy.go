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

package clusterappman

import (
	"context"
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/resource"

	kcore "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/clusterappman/status"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	rc "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	clusterappmanNamespace  = KoreNamespace
	clusterappmanDeployment = "kore-clusterappman"
	// DeployerServiceName is the name used for logging when deploying clusterappman
	DeployerServiceName = "clusterappman-deployer"
)

type deployerImpl struct {
	// ControllerClient is the controller runtime client for deploying clusterappman
	client rc.Client
	// ClusterAppManImage is the container image to use for clusterappman
	clusterAppManImage string
}

// NewLocalDeployer will enable a clusterappman deploy in the LOCAL kore kubernetes cluster
func NewLocalDeployer(client client.Client, clusterappmanImage string) Interface {
	return &deployerImpl{
		client:             client,
		clusterAppManImage: clusterappmanImage,
	}
}

// Run is responsible for starting the deployment services and keeping them running
func (d deployerImpl) Run(ctx context.Context) error {
	logger := log.WithFields(log.Fields{
		"service": DeployerServiceName,
	})

	logger.Debugf("clusterappman using image %s", d.clusterAppManImage)

	for {
		components, err := Deploy(ctx, d.client, logger, d.clusterAppManImage)
		if err != nil {
			logger.Errorf("error deploying clusterappman - %s", err)
		}
		if components != nil {
			status.SetAppManComponents(*components, d.client)
		}

		if utils.Sleep(ctx, 1*time.Minute) {
			logger.Info("clusterappman deployer stopped")
			return nil
		}
	}
}

// Stop will stop the services
func (d deployerImpl) Stop(ctx context.Context) error {
	return nil
}

// Deploy will install clusterappman in a cluster and return initial status
func Deploy(ctx context.Context, cc client.Client, logger *log.Entry, clusterAppManImage string) (*kcore.Components, error) {
	// @step: check if the cluster manager namespace exists and create it if not
	if err := EnsureNamespace(ctx, cc, clusterappmanNamespace); err != nil {
		logger.WithError(err).Errorf("trying to create the kore cluster-manager namespace %s", clusterappmanNamespace)

		return nil, err
	}

	// @step: ensure the service account
	if _, err := kubernetes.CreateOrUpdateServiceAccount(ctx, cc, &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "clusterappman",
			Namespace: clusterappmanNamespace,
			Labels: map[string]string{
				kore.Label("owner"): "true",
			},
		},
	}); err != nil {
		logger.WithError(err).Error("trying to create the clusterappman service account")

		return nil, err
	}
	// @step setup correct permissions for deployment
	if err := CreateClusterManClusterRoleBinding(ctx, cc); err != nil {
		logger.WithError(err).Error("can not create cluster-manager clusterrole")

		return nil, err
	}

	// @step: check if the kore cluster manager deployment exists
	if err := CreateOrUpdateClusterAppManDeployment(ctx, cc, clusterAppManImage); err != nil {
		logger.WithError(err).Error("trying to create the cluster manager deployment")

		return nil, err
	}

	return GetStatus(ctx, cc)
}

// CreateOrUpdateClusterAppManDeployment will reconcile the clusterappman deployment
func CreateOrUpdateClusterAppManDeployment(ctx context.Context, cc client.Client, image string) error {
	name := clusterappmanDeployment

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: clusterappmanNamespace,
			Labels: map[string]string{
				"name": name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": name,
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": name,
					},
				},
				Spec: v1.PodSpec{
					ServiceAccountName: "clusterappman",
					Containers: []v1.Container{
						{
							Name:  name,
							Image: image,
							Env: []v1.EnvVar{
								{
									Name:  "IN_CLUSTER",
									Value: "true",
								},
							},
							Command: []string{
								"/bin/kore-clusterappman",
							},
						},
					},
				},
			},
		},
	}

	if strings.HasSuffix(image, ":dev") {
		deployment.Spec.Template.Spec.Containers[0].Command = []string{
			"sh", "-c", "time make kore-clusterappman && bin/kore-clusterappman",
		}
		deployment.Spec.Template.Spec.Containers[0].VolumeMounts = []v1.VolumeMount{
			{
				Name:      "kore",
				MountPath: "/go/src/github.com/appvia/kore",
			},
			{
				Name:      "gocache",
				MountPath: "/root/.cache/go-build",
			},
		}

		deployment.Spec.Template.Spec.Volumes = []v1.Volume{
			{
				Name: "kore",
				VolumeSource: v1.VolumeSource{
					HostPath: &v1.HostPathVolumeSource{
						Path: "/go/src/github.com/appvia/kore",
						Type: (*v1.HostPathType)(utils.StringPtr(string(v1.HostPathDirectory))),
					},
				},
			},
			{
				Name: "gocache",
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: name + "-gocache",
					},
				},
			},
		}

		pvc := &v1.PersistentVolumeClaim{
			TypeMeta: metav1.TypeMeta{
				Kind:       "PersistentVolumeClaim",
				APIVersion: v1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name + "-gocache",
				Namespace: clusterappmanNamespace,
			},
			Spec: v1.PersistentVolumeClaimSpec{
				AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceStorage: resource.MustParse("1Gi"),
					},
				},
			},
		}
		exists, err := kubernetes.CheckIfExists(ctx, cc, pvc)
		if err != nil {
			return fmt.Errorf("failed to get persistent volume claim %q: %w", name+"-gocache", err)
		}
		if !exists {
			if err := cc.Create(ctx, pvc); err != nil {
				return fmt.Errorf("failed to create persistent volume claim %q: %w", name+"-gocache", err)
			}
		}
	}

	if _, err := kubernetes.CreateOrUpdate(ctx, cc, deployment); err != nil {
		return fmt.Errorf("failed to create deployment %q: %w", name, err)
	}

	return nil
}

// EnsureNamespace creates a namespace for the clusterappmanager if required
func EnsureNamespace(ctx context.Context, cc client.Client, namespace string) error {
	return kubernetes.EnsureNamespace(ctx, cc, &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
			Labels: map[string]string{
				kore.Label("owned"): "true",
			},
		},
	})
}

// CreateClusterManClusterRoleBinding creates (or updates) the cluster role binding required for the clusterappman
func CreateClusterManClusterRoleBinding(ctx context.Context, cc client.Client) error {
	if _, err := kubernetes.CreateOrUpdateManagedClusterRoleBinding(ctx, cc, &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kore:clusterappman",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "clusterappman",
				Namespace: KoreNamespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
	}); err != nil {
		return fmt.Errorf("error tying to apply kore clusterappman clusterrole %q", err)
	}
	return nil
}
