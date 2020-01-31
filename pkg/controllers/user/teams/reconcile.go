/*
 * Copyright (C) 2019 Appvia Ltd <info@appvia.io>
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

package teams

import (
	"context"
	"errors"

	core "github.com/appvia/kore/pkg/apis/core/v1"
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/hub"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "teams"
)

// Reconcile is the entrypoint for the reconcilation logic
func (t teamController) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
		"team":      request.NamespacedName.Name,
	})

	// @step: retrieve the team from the api
	team := &orgv1.Team{}
	if err := t.mgr.GetClient().Get(context.TODO(), request.NamespacedName, team); err != nil {
		if !kerrors.IsNotFound(err) {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	// @step: are we deleting the deleting team?
	finalizer := kubernetes.NewFinalizer(t.mgr.GetClient(), finalizerName)
	if finalizer.IsDeletionCandidate(team) {
		return t.Delete(ctx, request, team, finalizer)
	}

	err := func() error {
		// @step: check if the team name is protected - i.e. cannot be used
		if IsProtected(team.Name) {
			logger.Warn("team name is protected and cannot be used")

			team.Status.Status = core.FailureStatus
			team.Status.Conditions = []core.Condition{{
				Message: "team name is protected, you must rename the team",
			}}

			return errors.New("team name is protected, must change the name")
		}

		// @step: we first check if the team namespace is there
		err := t.mgr.GetClient().Get(ctx, types.NamespacedName{Name: team.Name}, &corev1.Namespace{})
		if err != nil {
			if !kerrors.IsNotFound(err) {
				return err
			}
			// @step: else we need to create the namespace
			if err := t.mgr.GetClient().Create(ctx, &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: team.Name,
					Labels: map[string]string{
						hub.Label("team"): team.Name,
					},
				},
			}); err != nil {
				return err
			}
		}

		// @step: update the status of the team as success
		team.Status.Conditions = []core.Condition{}
		team.Status.Status = core.SuccessStatus

		return nil
	}()
	if err == nil {
		if finalizer.NeedToAdd(team) {
			logger.Info("adding our finalizer to the team resource")

			if err := finalizer.Add(team); err != nil {
				return reconcile.Result{}, err
			}
		}

		return reconcile.Result{Requeue: true}, nil
	}
	if err != nil {
		logger.WithError(err).Error("failed to reconcile the team resource")
	}

	// @step: update the status of the team as failed
	if err := t.mgr.GetClient().Status().Update(ctx, team); err != nil {
		logger.WithError(err).Error("failed to update resource status")

		return reconcile.Result{}, err
	}

	return reconcile.Result{}, err
}
