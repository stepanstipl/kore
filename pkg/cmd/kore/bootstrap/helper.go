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

package bootstrap

import (
	"context"
	"fmt"
	"io"
)

const (
	markSuccess = "✅"
	markFailed  = "❌"
)

// TaskFunc is a action to implement
type TaskFunc func(context.Context) error

// Task is just a wrapper for output
type Task struct {
	// Description is a summary for the task
	Description string
	// Header providers a short predescription
	Header string
	// Handler is the action handler
	Handler TaskFunc
}

// Run is called to perform the task
func (t *Task) Run(ctx context.Context, w io.Writer) error {
	if t.Header != "" {
		fmt.Fprintf(w, "%s %s\n", markSuccess, t.Header)
	}

	err := func() error {
		if t.Handler == nil {
			return nil
		}

		return t.Handler(ctx)
	}()
	if err != nil {
		fmt.Fprintf(w, "%s %s\n", markFailed, t.Description)
	} else {
		if t.Description != "" {
			fmt.Fprintf(w, "%s %s\n", markSuccess, t.Description)
		}
	}

	return err
}
