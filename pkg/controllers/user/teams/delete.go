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
	"fmt"

	core "github.com/appvia/kore/pkg/apis/core/v1"
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is responsible for deleting the resource
func (t teamController) Delete(ctx context.Context, request reconcile.Request, team *orgv1.Team, finalizer *kubernetes.Finalizer) (reconcile.Result, error) {
	fields := log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
		"team":      request.NamespacedName.Name,
	}

	log.WithFields(fields).Info("handling the deletion of the team resource")

	err := func() error {
		// @step: we first check if the team namespace is there
		if err := t.mgr.GetClient().Get(ctx, types.NamespacedName{Name: team.Name}, &corev1.Namespace{}); err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
			log.WithFields(fields).Debug("team namespace does not exist")

			return nil
		}

		// @step: delete the team namespace
		return t.mgr.GetClient().Delete(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: team.Name}})
	}()
	if err != nil {
		fields["error"] = err.Error()
		log.WithFields(fields).Error("failed to reconcile deletion of team")

		return reconcile.Result{}, err
	}

	log.Debug("attempting to remove the teams-controller finalizer")
	// @step: removed the finalizer and allow the resource to be deleted
	if err := finalizer.Remove(team); err != nil {
		fields["error"] = err.Error()
		log.WithFields(fields).Error("failed to reconcile deletion of team")

		team.Status.Status = core.FailureStatus
		team.Status.Conditions = []core.Condition{{
			Detail:  err.Error(),
			Message: fmt.Sprintf("Failed to remove the team due as all resources not removed"),
		}}

		return reconcile.Result{}, t.mgr.GetClient().Status().Update(ctx, team)
	}
	log.Info("teams resource has been successfully reconciled")

	return reconcile.Result{}, nil
}
