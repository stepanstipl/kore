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
	"errors"
	"fmt"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var ErrNotReady = NewReconcileError(errors.New("not ready yet"), false)

type ReconcileError struct {
	Err      error
	Critical bool
}

func NewReconcileError(err error, critical bool) *ReconcileError {
	return &ReconcileError{
		Err:      err,
		Critical: critical,
	}
}

func (r *ReconcileError) Error() string {
	return r.Err.Error()
}

func (r *ReconcileError) Result() (reconcile.Result, error) {
	if r == nil {
		return reconcile.Result{}, nil
	}

	if r.Critical {
		return reconcile.Result{}, r.Err
	}

	return reconcile.Result{Requeue: true}, nil
}

func (r *ReconcileError) ResultWithRequeueAfter(requeueAfter time.Duration) (reconcile.Result, error) {
	res, err := r.Result()
	if res.Requeue && requeueAfter > 0 {
		res.Requeue = false
		res.RequeueAfter = requeueAfter
	}
	return res, err
}

func (r *ReconcileError) Wrap(message string) *ReconcileError {
	r.Err = fmt.Errorf("%s: %w", message, r.Err)
	return r
}

func (r *ReconcileError) Wrapf(format string, args ...interface{}) *ReconcileError {
	r.Err = fmt.Errorf(format, append(args, r.Err)...)
	return r
}
