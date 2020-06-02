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

package teams

import (
	"context"
	"errors"

	core "github.com/appvia/kore/pkg/apis/core/v1"
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/kore"
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

// Reconcile is the entrypoint for the reconciliation logic
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
		if !kore.IsValidTeamName(team.Name) {
			logger.Warn("team name is protected and cannot be used")

			team.Status.Status = core.FailureStatus
			team.Status.Conditions = []core.Condition{{
				Message: "Team name is protected, you must rename the team",
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
						kore.Label("team"): team.Name,
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
