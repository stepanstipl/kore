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

package clusterbindings

import (
	"context"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type crCtrl struct {
	kore.Interface
	// mgr is the manager
	mgr manager.Manager
	// stopCh is the stop channel
	stopCh chan struct{}
}

func init() {
	controllers.Register(&crCtrl{})
}

// Name returns the name of the controller
func (a crCtrl) Name() string {
	return "cluster-bindings"
}

// Run is called when the controller is started
func (a *crCtrl) Run(ctx context.Context, cfg *rest.Config, hi kore.Interface) error {
	a.Interface = hi

	// @step: create the manager for the controller
	mgr, err := manager.New(cfg, controllers.DefaultManagerOptions(a))
	if err != nil {
		log.WithError(err).Error("failed to create the manager")

		return err
	}
	a.mgr = mgr

	// @step: create the controller
	ctrl, err := controller.New(a.Name(), mgr, controllers.DefaultControllerOptions(a))
	if err != nil {
		log.WithError(err).Error("trying to create the controller")

		return err
	}

	// @step: setup watches for the resources
	if err := ctrl.Watch(&source.Kind{Type: &clustersv1.ManagedClusterRoleBinding{}},
		&handler.EnqueueRequestForObject{},
		&predicate.GenerationChangedPredicate{},
	); err != nil {

		log.WithField("error", err.Error()).Error("failed to create watcher on resource")

		return err
	}
	go func() {
		log.Info("starting the controller loop")

		for {
			a.stopCh = make(chan struct{})

			if err := mgr.Start(a.stopCh); err != nil {
				log.WithError(err).Error("failed to start the controller")
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// @step: use a routine to catch the stop channel
	go func() {
		<-ctx.Done()
		log.WithFields(log.Fields{
			"controller": a.Name(),
		}).Info("stopping the controller")

		close(a.stopCh)
	}()

	return nil
}

// Stop is responsible for calling a halt on the controller
func (a *crCtrl) Stop(context.Context) error {
	log.WithFields(log.Fields{
		"controller": a.Name(),
	}).Info("attempting to stop the controller")

	return nil
}

// FilterClustersBySource returns a list of kubenetes cluster in the kore - if the
// namespace is global we retrieve all clusters, else just the local teams
func (a *crCtrl) FilterClustersBySource(ctx context.Context,
	clusters []corev1.Ownership,
	teams []string,
	namespace string) (*clustersv1.KubernetesList, error) {

	list := &clustersv1.KubernetesList{}

	// @step: is the role targeting a specific cluster
	if len(clusters) > 0 {
		item := &clustersv1.Kubernetes{}
		for _, x := range clusters {
			if err := a.mgr.GetClient().Get(ctx, types.NamespacedName{
				Name:      x.Name,
				Namespace: x.Namespace,
			}, item); err != nil {
				if !kerrors.IsNotFound(err) {
					return list, err
				}

				continue
			}

			list.Items = append(list.Items, *item)
		}

		return list, nil
	}

	// @step: check if it's filter down to teams
	if len(teams) > 0 {
		for _, x := range teams {
			clusters := &clustersv1.KubernetesList{}

			if err := a.mgr.GetClient().List(ctx, clusters, client.InNamespace(x)); err != nil {
				return list, err
			}
			list.Items = append(list.Items, clusters.Items...)
		}

		return list, nil
	}

	if kore.IsGlobalTeam(namespace) {
		return list, a.mgr.GetClient().List(ctx, list, client.InNamespace(""))
	}

	return list, a.mgr.GetClient().List(ctx, list, client.InNamespace(namespace))
}
