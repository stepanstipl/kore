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
	"context"
	"time"

	"github.com/appvia/kore/pkg/kore"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Context interface {
	context.Context
	Logger() log.FieldLogger
	Client() client.Client
	Kore() kore.Interface
}

type contextImpl struct {
	ctx    context.Context
	logger log.FieldLogger
	client client.Client
	kore   kore.Interface
}

func NewContext(ctx context.Context, logger log.FieldLogger, client client.Client, kore kore.Interface) Context {
	return contextImpl{
		ctx:    ctx,
		logger: logger,
		client: client,
		kore:   kore,
	}
}

func (c contextImpl) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c contextImpl) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c contextImpl) Err() error {
	return c.ctx.Err()
}

func (c contextImpl) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

func (c contextImpl) Logger() log.FieldLogger {
	return c.logger
}

func (c contextImpl) Client() client.Client {
	return c.client
}

func (c contextImpl) Kore() kore.Interface {
	return c.kore
}
