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

package controllers

import (
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Component interface {
	Reconcile(ctx kore.Context) (reconcile.Result, error)
	Delete(ctx kore.Context) (reconcile.Result, error)
}

type ComponentWithStatus interface {
	Component
	ComponentName() string
	SetComponent(*corev1.Component)
}

type ComponentWithDeletionStatus interface {
	Component
	IsDeleted(kore.Context) (bool, error)
}

type Components []Component

func (c Components) Reconcile(ctx kore.Context, object kubernetes.ObjectWithStatus) (reconcile.Result, error) {
	var res reconcile.Result
	var err error

	if object.GetDeletionTimestamp().IsZero() {
		res, err = c.reconcile(ctx, object)
	} else {
		res, err = c.delete(ctx, object)
	}

	if err != nil {
		if IsCriticalError(err) {
			object.SetStatus(corev1.FailureStatus, err.Error())
			return reconcile.Result{}, nil
		}

		object.SetStatus(corev1.ErrorStatus, err.Error())
		return reconcile.Result{}, err
	}

	return res, nil
}

func (c Components) reconcile(ctx kore.Context, object kubernetes.ObjectWithStatus) (reconcile.Result, error) {
	object.SetStatus(corev1.PendingStatus, "")

	for _, comp := range c {
		c.initComponent(comp, object.StatusComponents(), corev1.PendingStatus)

		res, err := comp.Reconcile(ctx)

		c.setComponentStatusFromResult(object, comp, corev1.SuccessStatus, res, err)

		if err != nil || res.Requeue || res.RequeueAfter > 0 {
			return res, err
		}
	}

	object.SetStatus(corev1.SuccessStatus, "")

	return reconcile.Result{}, nil
}

func (c Components) delete(ctx kore.Context, object kubernetes.ObjectWithStatus) (reconcile.Result, error) {
	object.SetStatus(corev1.DeletingStatus, "")

	startDeleteFrom := -1
	for i, comp := range c {
		if isDeleted, err := c.isDeleted(ctx, object, comp); err != nil {
			return reconcile.Result{}, err
		} else if isDeleted {
			break
		}
		startDeleteFrom = i
	}

	for i := startDeleteFrom; i >= 0; i-- {
		comp := c[i]

		c.initComponent(comp, object.StatusComponents(), corev1.DeletingStatus)

		res, err := comp.Delete(ctx)

		c.setComponentStatusFromResult(object, comp, corev1.DeletedStatus, res, err)

		if err != nil || res.Requeue || res.RequeueAfter > 0 {
			return res, err
		}
	}

	object.SetStatus(corev1.DeletedStatus, "")

	return reconcile.Result{}, nil
}

func (c Components) initComponent(comp Component, statuses *corev1.Components, defaultStatus corev1.Status) {
	if compWithStatus, ok := comp.(ComponentWithStatus); ok {
		if _, exists := statuses.GetComponent(compWithStatus.ComponentName()); !exists {
			statuses.SetCondition(corev1.Component{
				Name: compWithStatus.ComponentName(),
			})
		}

		compRef, _ := statuses.GetComponent(compWithStatus.ComponentName())
		compRef.Status = defaultStatus
		compRef.Message = ""
		compRef.Detail = ""

		compWithStatus.SetComponent(compRef)
	}
}

func (c Components) setComponentStatusFromResult(
	object kubernetes.ObjectWithStatus, comp Component, successStatus corev1.Status, res reconcile.Result, err error,
) {
	if compWithStatus, ok := comp.(ComponentWithStatus); ok {
		if err != nil {
			if IsCriticalError(err) {
				object.StatusComponents().SetStatus(
					compWithStatus.ComponentName(), corev1.FailureStatus, "failed to reconcile", err.Error(),
				)
				return
			}

			object.StatusComponents().SetStatus(
				compWithStatus.ComponentName(), corev1.ErrorStatus, "failed to reconcile", err.Error(),
			)
			return
		}

		if !res.Requeue && res.RequeueAfter == 0 {
			object.StatusComponents().SetStatus(
				compWithStatus.ComponentName(), successStatus, "", "",
			)
		}
	}
}

func (c Components) isDeleted(ctx kore.Context, object kubernetes.ObjectWithStatus, comp Component) (bool, error) {
	if compWithStatus, ok := comp.(ComponentWithStatus); ok {
		status, exists := object.StatusComponents().GetComponent(compWithStatus.ComponentName())
		if !exists {
			return false, nil
		}
		return status.Status == corev1.DeletedStatus, nil
	}

	if compWithDeletionStatus, ok := comp.(ComponentWithDeletionStatus); ok {
		return compWithDeletionStatus.IsDeleted(ctx)
	}

	return false, nil
}
