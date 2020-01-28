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
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/appvia/kore/pkg/hub"
	"github.com/appvia/kore/pkg/utils"

	clusterv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	"github.com/gambol99/hub-utils/pkg/finalizers"
	kutils "github.com/gambol99/hub-utils/pkg/kubernetes"
	log "github.com/sirupsen/logrus"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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
	requeue := time.Duration(time.Minute * 1)

	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	ctx := context.Background()

	// @step: retrieve the resource from the api
	cluster := &clusterv1.Kubernetes{}
	if err := t.mgr.GetClient().Get(ctx, request.NamespacedName, cluster); err != nil {
		if !kerrors.IsNotFound(err) {
			return reconcile.Result{}, err
		}

		return reconcile.Result{RequeueAfter: requeue}, nil
	}

	finalizer := finalizers.NewFinalizer(t.mgr.GetClient(), finalizerName)
	if finalizer.IsDeletionCandidate(cluster) {
		return t.Delete(ctx, request)
	}

	logger.Info("reconciling the kubernetes cluster")

	// @step: ensure the phase of the resource
	if cluster.Status.Status != corev1.SuccessStatus {
		return reconcile.Result{RequeueAfter: requeue}, nil
	}

	// @step: ignore the resource if already bootstrapped
	if cluster.Status.Phase == "Installed" {
		return reconcile.Result{RequeueAfter: requeue}, nil
	}

	// @step: ensure we have a finalizer on this resource
	if finalizer.NeedToAdd(cluster) {
		logger.Info("adding our finalizer on the resource")

		if err := finalizer.Add(cluster); err != nil {
			logger.WithError(err).Error("failed to add finalizer on resource")

			return reconcile.Result{RequeueAfter: requeue}, err
		}

		// the resource generation has changed we need to requeue
		return reconcile.Result{Requeue: true}, nil
	}

	// @step: ensure the status is updated
	cluster.Status.Phase = "Installing"

	if err := t.mgr.GetClient().Status().Update(ctx, cluster); err != nil {
		logger.WithError(err).Error("failed to update resource status")

		return reconcile.Result{RequeueAfter: requeue}, err
	}

	// @step: lets be positive
	cluster.Status.Conditions = []corev1.Condition{}

	// @logic
	// - create a kubernetes client to the remote cluster
	// - retrieve the credentials for the broker from the cluster provider
	// - wait for kube api to be ready

	err := func() error {
		// @step: retrieve the provider instances
		pi, err := t.GetCloudInstance(ctx, cluster)
		if err != nil {
			logger.WithError(err).Error("failed to retrieve instance resource")

			return err
		}
		kind := pi.GetObjectKind().GroupVersionKind().Kind
		provider := strings.ToLower(kind)

		// @step: retrieve the credentials for the provider
		// kind of a noop but will be used in after this is defined.
		pc, err := t.GetCloudCredentials(ctx, pi)
		if err != nil {
			logger.WithError(err).Error("failed to retrieve cloud credentials")

			return err
		}

		// @step: ensure we can make a k8s client
		client, err := makeKubernetesClient(cluster)
		if err != nil {
			logger.WithError(err).Error("failed to create kubernetes client")

			return err
		}

		// ensure we can access the api
		if err := kutils.WaitOnKubeAPI(ctx, client, 5*time.Second, 60*time.Second); err != nil {
			logger.WithError(err).Error("failed to access to the kubernetes api")

			return err
		}

		// @step: push the namespace admin
		if found, err := ClusterRoleExists(ctx, client, "hub:system:ns-admin"); err != nil {
			logger.WithError(err).Error("failed to check for namespace admin")
		} else if !found {
			c, err := MakeTemplate(NamespaceAdminClusterRole, map[string]string{})
			if err != nil {
				logger.WithError(err).Error("failed to generate namespace admin template")

				return err
			}
			cm := &rbac.ClusterRole{}
			if err := DecodeInTo(c, cm); err != nil {
				logger.WithError(err).Error("failed to decode namespace admin")

				return err
			}

			if _, err := client.RbacV1().ClusterRoles().Create(cm); err != nil {
				logger.WithError(err).Error("failed to create namespace admin")

				return err
			}
		}

		// @step: check if the bundles exists			logger.Info("kube api ready for new cluster")
		if found, err := JobConfigExists(ctx, client); err != nil {
			logger.WithError(err).Error("failed to check for configuration configmap")

			return err
		} else if !found {
			c, err := MakeTemplate(BootstrapJobConfigmap, map[string]string{})
			if err != nil {
				logger.WithError(err).Error("failed to generate bootstrap template")

				return err
			}

			cm := &core.ConfigMap{}
			if err := DecodeInTo(c, cm); err != nil {
				logger.WithError(err).Error("failed to decode bootstrap configuration in configmap")

				return err
			}

			if _, err := client.CoreV1().ConfigMaps(jobNamespace).Create(cm); err != nil {
				logger.WithError(err).Error("failed to create the bootstrap configuration configmap")

				return err
			}
		}

		// @step: build the parameters for the job
		params, err := t.GetClusterConfiguration(ctx, provider)
		if err != nil {
			logger.WithError(err).Error("failed to generate the parameters")

			return err
		}
		params.Credentials = pc

		// @step: we check if the job configuration is there already and if not we make it
		if found, err := JobOLMConfigExists(ctx, client); err != nil {
			logger.WithError(err).Error("failed to check for olm job configuration")

			return err
		} else if !found {
			c, err := MakeTemplate(BootstrapJobOLMConfig, params)
			if err != nil {
				logger.WithError(err).Error("failed to generate the bootstrap olm template")

				return err
			}
			cm := &core.ConfigMap{}
			if err := DecodeInTo(c, cm); err != nil {
				logger.WithError(err).Error("failed to decode olm bootstrap configuration in configmap")

				return err
			}

			if _, err := client.CoreV1().ConfigMaps(jobNamespace).Create(cm); err != nil {
				logger.WithError(err).Error("failed to create the olm bootstrap configuration configmap")

				return err
			}
		}

		// @step: ensure the bootstrap job is there
		if found, err := JobExists(ctx, client); err != nil {
			logger.WithError(err).Error("failed to check for bootstrap job")

			return err
		} else if !found {
			c, err := MakeTemplate(BootstrapJobTemplate, params)
			if err != nil {
				logger.WithError(err).Error("failed to render the bootstrap job")

				return err
			}
			job := &batch.Job{}
			if err := DecodeInTo(c, job); err != nil {
				logger.WithError(err).Error("failed to decode the job")

				return err
			}
			if _, err := client.BatchV1().Jobs(jobNamespace).Create(job); err != nil {
				logger.WithError(err).Error("failed to create the bootstrap job")

				return err
			}
		}
		logger.Info("waiting for bootstrap job to finish")

		nctx, cancel := context.WithTimeout(ctx, 20*time.Minute)
		defer cancel()

		// @step: wait for the bootstrap job to complete
		if err := WaitOnJob(nctx, client); err != nil {
			logger.WithError(err).Error("failed waiting for bootstrap to complete")

			return err
		}
		logger.Info("waiting for bootstrap has completed")

		logger.Info("kube api ready for new cluster")

		// @step: else we can set the job as complete
		cluster.Status.Phase = "Installed"

		return nil
	}()
	if err != nil {
		cluster.Status.Status = corev1.FailureStatus
		cluster.Status.Phase = "Failed"
		cluster.Status.Conditions = []corev1.Condition{{
			Detail:  err.Error(),
			Message: "failed to bootstrap the kubernetes cluster",
		}}

		logger.WithError(err).Error("failed to reconcile the cluster")
	}

	// @step: update the status of the resource
	if err := t.mgr.GetClient().Status().Update(ctx, cluster); err != nil {
		logger.WithError(err).Error("failed to update status of resource")

		return reconcile.Result{}, err
	}

	return reconcile.Result{RequeueAfter: requeue}, nil
}

// GetClusterConfiguration is responsible for generate the parameters for the cluster
func (t bsCtrl) GetClusterConfiguration(ctx context.Context, provider string) (Parameters, error) {
	params := Parameters{
		BootImage: "quay.io/appvia/hub-bootstrap:v0.2.0",
		Broker: BrokerOptions{
			Username: "broker",
			Password: utils.Random(12),
			Database: DatabaseOptions{
				Name:     "broker",
				Password: utils.Random(12),
			},
		},
		Catalog:             CatalogOptions{Image: "v0.0.1"},
		EnableIstio:         false,
		EnableKiali:         false,
		EnableServiceBroker: true,
		Kiali: KialiOptions{
			Password: utils.Random(12),
		},
		Provider:   provider,
		OLMVersion: "0.11.0",
		Namespaces: []NamespaceOptions{
			{Name: "brokers", EnableIstio: false},
			{Name: "grafana", EnableIstio: false},
			{Name: "kube-dns", EnableIstio: false},
			{Name: "logging", EnableIstio: false},
			{Name: "prometheus", EnableIstio: false},
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

	if params.EnableServiceBroker {
		switch params.Provider {
		case "gke":
			params.Operators = append(params.Operators, OperatorOptions{
				Package:   "gcp-service-broker-operator",
				Channel:   "stable",
				Label:     "name=gcp-service-broker-operator",
				Namespace: "brokers",
			})
			params.Operators = append(params.Operators, OperatorOptions{
				Package:   "mariadb-operator",
				Channel:   "stable",
				Label:     "name=mariadb-operator",
				Namespace: "brokers",
			})
		case "eks":
			params.Operators = append(params.Operators, OperatorOptions{
				Package:   "aws-service-broker-operator",
				Channel:   "stable",
				Label:     "name=aws-service-broker-operator",
				Namespace: "brokers",
			})
		}
	}

	if params.EnableIstio {
		params.OperatorGroups = []string{"istio-system"}
		params.EnableKiali = true
		params.Operators = append(params.Operators, OperatorOptions{
			Catalog:   "operatorhubio-catalog",
			Package:   "aws-service-broker-operator",
			Channel:   "stable",
			Label:     "app=kiali-operator",
			Namespace: "istio-system",
		})
	}

	return params, nil
}

// GetCloudCredentials returns the cloud credentials for this provider
func (t bsCtrl) GetCloudCredentials(ctx context.Context, obj runtime.Object) (Credentials, error) {
	// @step: retrieve the instance the resource was created from
	m, ok := obj.(metav1.Object)
	if !ok {
		return Credentials{}, errors.New("resource does not implement metav1.Object interface")
	}

	name, found := m.GetLabels()[hub.Label("binding")]
	if !found {
		return Credentials{}, errors.New("resource does not have a binding label")
	}

	team, found := m.GetLabels()[hub.Label("team")]
	if !found {
		return Credentials{}, errors.New("resource does not have a team label")
	}

	binding, err := t.Teams().Team(team).Bindings().Get(ctx, name)
	if err != nil {
		return Credentials{}, err
	}

	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Kind:    binding.Spec.Ref.Kind,
		Group:   binding.Spec.Ref.Group,
		Version: binding.Spec.Ref.Version,
	})
	if err := t.mgr.GetClient().Get(ctx, types.NamespacedName{
		Namespace: binding.Spec.Ref.Namespace,
		Name:      binding.Spec.Ref.Name}, u); err != nil {
		return Credentials{}, err
	}

	spec, found := u.Object["spec"].(map[string]interface{})
	if !found {
		return Credentials{}, errors.New("resource does not have a spec")
	}

	c := Credentials{}

	// the hardcoded part
	switch cp := obj.GetObjectKind().GroupVersionKind().Kind; cp {
	case "GKE":
		c.GKE.Account = fmt.Sprintf("%s", spec["account"])
	case "EKS":
		c.AWS.AccessKey = fmt.Sprintf("%s", spec["accessKey"])
		c.AWS.AccountID = fmt.Sprintf("%s", spec["accountID"])
		c.AWS.Region = fmt.Sprintf("%s", spec["region"])
		c.AWS.SecretKey = fmt.Sprintf("%s", spec["secretKey"])
	default:
		return c, fmt.Errorf("unkwown cloud provider type: %s", cp)
	}

	return c, nil
}

// GetCloudInstance returns the cloud instance
func (t bsCtrl) GetCloudInstance(ctx context.Context, c *clusterv1.Kubernetes) (runtime.Object, error) {
	// @step: ensure we have a provider and retrieve the instance
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   c.Spec.Use.Group,
		Version: c.Spec.Use.Version,
		Kind:    c.Spec.Use.Kind,
	})

	return u, t.mgr.GetClient().Get(ctx, types.NamespacedName{
		Name:      c.Spec.Use.Name,
		Namespace: c.Spec.Use.Namespace,
	}, u)
}
