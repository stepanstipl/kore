/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
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
	if err := controllers.Register(&crCtrl{}); err != nil {
		log.WithError(err).Fatal("failed to register cluster bindings controller")
	}
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
