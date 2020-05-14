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

package projectsweeper

import (
	"context"
	"time"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	gcp "github.com/appvia/kore/pkg/apis/gcp/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	"github.com/patrickmn/go-cache"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "projectsweeper.gcp.compute.kore.appvia.io"
)

// Reconcile is the entrypoint for the reconciliation logic
func (c *Controller) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.TODO()

	logger := c.logger.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to reconcile the projects")

	result, err := func() (reconcile.Result, error) {
		return controllers.DefaultEnsureHandler.Run(ctx,
			[]controllers.EnsureFunc{
				c.EnsureRemoval(),
			},
		)
	}()

	if err != nil {
		logger.WithError(err).Error("failed to reconcile the project sweeper")
	}

	return result, nil
}

// EnsureRemoval is responsible for deleting any project no longer claims
func (c *Controller) EnsureRemoval() controllers.EnsureFunc {
	cc := c.mgr.GetClient()

	return func(ctx context.Context) (reconcile.Result, error) {
		// @logic
		// - we retrieve a list of projects
		// - we retrieve a list of claims
		// - we find any projects whom no longer have a claim for x period of time and delete it

		// @step: ensure no claim exists outside of the team
		projects := &gcp.ProjectList{}
		if err := cc.List(ctx, projects, client.InNamespace("")); err != nil {
			c.logger.WithError(err).Error("trying to retrieve all the projects")

			return reconcile.Result{}, err
		}

		claims := &gcp.ProjectClaimList{}
		if err := cc.List(ctx, claims, client.InNamespace("")); err != nil {
			c.logger.WithError(err).Error("trying to retrieve all the projects")

			return reconcile.Result{}, err
		}

		// @step: iterate the projects and looks for claims
		for i := 0; i < len(projects.Items); i++ {
			switch projects.Items[i].Status.Status {
			case corev1.PendingStatus, corev1.DeletingStatus, "":
				continue
			}

			name := projects.Items[i].Spec.ProjectName

			// has we need this project before?
			v, found := c.cache.Get(name)
			if !found {
				c.cache.Set(name, time.Now(), cache.DefaultExpiration)

				continue
			}
			timer := v.(time.Time)

			// @step: double check no cluster is using this as a projectclaim
			logger := c.logger.WithFields(log.Fields{
				"project":   projects.Items[i].Spec.ProjectName,
				"name":      projects.Items[i].Name,
				"namespace": projects.Items[i].Namespace,
			})

			claimed, err := c.isClaimed(&projects.Items[i], claims)
			if err != nil {
				c.cache.Set(name, time.Now(), cache.DefaultExpiration)

				logger.WithError(err).Error("trying to check if project is claimed")

				return reconcile.Result{}, err
			}

			switch claimed {
			case true:
				logger.Debug("project is still referenced, keeping for now")

				c.cache.Set(name, time.Now(), cache.DefaultExpiration)

			default:
				if time.Since(timer) > 10*time.Minute {
					logger.Info("attempting to delete the unclaimed gcp project")

					if err := kubernetes.DeleteIfExists(ctx, cc, &projects.Items[i]); err != nil {
						logger.WithError(err).Error("trying to delete the gcp project")

						return reconcile.Result{}, err
					}
				}
			}
		}

		return reconcile.Result{}, nil
	}
}

// isClaimed checks if the project has a matching claim
func (c *Controller) isClaimed(project *gcp.Project, claims *gcp.ProjectClaimList) (bool, error) {
	for i := 0; i < len(claims.Items); i++ {
		match, err := kore.IsResourceOwner(project, claims.Items[i].Status.ProjectRef)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}

	return false, nil
}
