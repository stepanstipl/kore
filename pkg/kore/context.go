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

package kore

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Context interface {
	context.Context
	Logger() log.Ext1FieldLogger
	Client() client.Client
	Kore() Interface
	WithLogger(log.Ext1FieldLogger) Context
}

type contextImpl struct {
	ctx    context.Context
	logger log.Ext1FieldLogger
	client client.Client
	kore   Interface
}

func (s contextImpl) Deadline() (deadline time.Time, ok bool) {
	return s.ctx.Deadline()
}

func (s contextImpl) Done() <-chan struct{} {
	return s.ctx.Done()
}

func (s contextImpl) Err() error {
	return s.ctx.Err()
}

func (s contextImpl) Value(key interface{}) interface{} {
	return s.ctx.Value(key)
}

func (s contextImpl) Logger() log.Ext1FieldLogger {
	return s.logger
}

func (s contextImpl) Client() client.Client {
	return s.client
}

func (s contextImpl) Kore() Interface {
	return s.kore
}

func (s contextImpl) WithLogger(logger log.Ext1FieldLogger) Context {
	return contextImpl{
		ctx:    s.ctx,
		logger: logger,
		client: s.client,
		kore:   s.kore,
	}
}

func NewContext(
	ctx context.Context,
	logger log.Ext1FieldLogger,
	client client.Client,
	kore Interface,
) Context {
	return contextImpl{
		ctx:    ctx,
		logger: logger,
		client: client,
		kore:   kore,
	}
}
