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

// Reconcile is the entrypoint for the reconciliation logic
func (c *Controller) Reconcile(ctx kore.Context, request reconcile.Request) (reconcile.Result, error) {
	ctx.Logger().Debug("attempting to cleanup the gcp projects")

	result, err := func() (reconcile.Result, error) {
		return controllers.DefaultEnsureHandler.Run(ctx,
			[]controllers.EnsureFunc{
				c.EnsureRemoval(),
			},
		)
	}()

	if err != nil {
		ctx.Logger().WithError(err).Error("failed to reconcile the project sweeper")
	}

	return result, nil
}

// EnsureRemoval is responsible for deleting any project no longer claims
func (c *Controller) EnsureRemoval() controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		// @logic
		// - we retrieve a list of projects
		// - we retrieve a list of claims
		// - we find any projects whom no longer have a claim for x period of time and delete it

		// @step: ensure no claim exists outside of the team
		projects := &gcp.ProjectList{}
		if err := ctx.Client().List(ctx, projects, client.InNamespace("")); err != nil {
			ctx.Logger().WithError(err).Error("trying to retrieve all the projects")

			return reconcile.Result{}, err
		}

		claims := &gcp.ProjectClaimList{}
		if err := ctx.Client().List(ctx, claims, client.InNamespace("")); err != nil {
			ctx.Logger().WithError(err).Error("trying to retrieve all the projects")

			return reconcile.Result{}, err
		}

		ctx.Logger().WithField(
			"size", len(projects.Items),
		).Debug("found the following gcp projects")

		// @step: iterate the projects and looks for claims
		for i := 0; i < len(projects.Items); i++ {
			name := projects.Items[i].Spec.ProjectName

			switch projects.Items[i].Status.Status {
			case corev1.PendingStatus, corev1.DeletingStatus, "":
				ctx.Logger().WithField(
					"name", name,
				).Debug("skipping the gcp project due to state")

				continue
			}

			// has we need this project before?
			v, found := c.cache.Get(name)
			if !found {
				ctx.Logger().WithField(
					"name", name,
				).Debug("first time seeing the project, holding off for now")

				c.cache.Set(name, time.Now(), cache.DefaultExpiration)

				continue
			}
			timer := v.(time.Time)

			// @step: double check no cluster is using this as a projectclaim
			logger := ctx.Logger().WithFields(log.Fields{
				"age":       time.Since(timer).String(),
				"name":      projects.Items[i].Name,
				"namespace": projects.Items[i].Namespace,
				"project":   projects.Items[i].Spec.ProjectName,
			})
			logger.Debug("checking if the project requires deletion")

			claimed, err := c.isClaimed(&projects.Items[i], claims)
			if err != nil {
				c.cache.Set(name, time.Now(), cache.DefaultExpiration)

				logger.WithError(err).Error("trying to check if project is claimed")

				return reconcile.Result{}, err
			}

			expiration := time.Duration(10 * time.Minute)

			switch claimed {
			case true:
				logger.Debug("project is still referenced, keeping for now")

				c.cache.Set(name, time.Now(), cache.DefaultExpiration)

			default:
				if time.Since(timer) > expiration {
					logger.Info("attempting to delete the unclaimed gcp project")

					if err := kubernetes.DeleteIfExists(ctx, ctx.Client(), &projects.Items[i]); err != nil {
						logger.WithError(err).Error("trying to delete the gcp project")

						return reconcile.Result{}, err
					}
					c.cache.Delete(name)

				} else {
					logger.WithField(
						"expires", time.Duration(time.Since(timer)-expiration).String(),
					).Debug("project has the following period before deletion")
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
