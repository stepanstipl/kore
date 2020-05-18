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

package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/client"
	kutils "github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// WaitForDeletion is used to wait for a resource to disappear
func (f *factory) WaitForDeletion(request client.RestInterface, name string, nowait bool) error {
	if err := request.Delete().Error(); err != nil {
		return err
	}

	if nowait {
		f.Println("Successfully requested resource %q to be deleted", name)

		return nil
	}

	f.Println("Waiting for resource %q to delete", name)

	// @step: setup the signalling
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// @step: we need to handle the user cancelling the blocking
	go func() {
		<-interrupt
		cancel()
	}()

	// @step: wait for the resource to disappear
	err := kutils.WaitUntilComplete(ctx, 20*time.Minute, 5*time.Second, func() (bool, error) {
		if err := request.Get().Error(); err != nil {
			if client.IsNotAllowed(err) {
				return false, err
			}
			if client.IsNotAuthorized(err) {
				return false, err
			}
			if client.IsNotFound(err) {
				return true, nil
			}

			return false, nil
		}

		return false, nil
	})
	if err != nil {
		if err == kutils.ErrCancelled {
			return nil
		}

		return err
	}
	f.Println("Resource %q has been deleted", name)

	return nil
}

// WaitForCreation is used to wait for the resource to be created
// @TODO clean this method up
func (f *factory) WaitForCreation(request client.RestInterface, nowait bool) error {
	// @step: retrieve the metav1.Object
	mo, err := kutils.GetMetaObject(request.GetPayload())
	if err != nil {
		return err
	}
	name := mo.GetName()

	// @step: retrieve the runtime payload as runtime.Object
	obj, err := kutils.GetRuntimeObject(request.GetPayload())
	if err != nil {
		return err
	}
	gvk := obj.GetObjectKind().GroupVersionKind()
	kind := gvk.Kind

	// @step: check if the object already exists
	found, err := request.Exists()
	if err != nil {
		return err
	}
	if found {
		selflink, err := kubernetes.GetRuntimeSelfLink(obj)
		if err != nil {
			return err
		}

		return fmt.Errorf("Resource %q already exists, please edit instead", selflink)
	}

	// @step: attempt to create the resource in kore
	if err := request.Update().Error(); err != nil {
		return err
	}

	// we are return straight if not waiting
	if nowait {
		f.Println("Successfully requested the resource %q to provision", name)

		return nil
	}

	// maxFailure is the max number of requests where the status
	// is failed we are willing to accept
	maxAttempts := 5
	// attempts is the above we have reached
	var attempts int

	// @step: setup the signalling
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// @step: create a cancellable context to operate within
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	// @step: we need to handle the user cancelling the blocking
	go func() {
		<-interrupt
		cancel()
	}()

	f.Println("Waiting for resource %q to provision (you can background with ctrl-c)", name)

	u := &unstructured.Unstructured{}

	// @step: craft the status from the resource type - used later
	status := fmt.Sprintf("kore get %s", strings.ToLower(kind))
	if v, found := request.HasParamater("name"); found {
		status = fmt.Sprintf("%s %s", status, v)
	}
	if v, found := request.HasParamater("team"); found {
		status = fmt.Sprintf("%s -t %s", status, v)
	}

	err = kutils.WaitUntilComplete(ctx, 20*time.Minute, 5*time.Second, func() (bool, error) {
		if err := request.Payload(nil).Result(u).Get().Error(); err != nil {
			if client.IsNotAllowed(err) || client.IsMethodNotAllowed(err) {
				return false, err
			}
			if client.IsNotAuthorized(err) {
				return false, err
			}

			return false, nil
		}

		// @step: check the status of the resource
		status, ok := u.Object["status"].(map[string]interface{})
		if !ok {
			return false, nil
		}
		state, ok := status["status"].(string)
		if !ok {
			return false, nil
		}

		switch state {
		case string(corev1.FailureStatus):
			if attempts > maxAttempts {
				return false, errors.New("resource has failed to provision")
			}
			attempts++
		case string(corev1.SuccessStatus):
			return true, nil
		}

		return false, nil
	})

	if err != nil {
		if err == kutils.ErrCancelled {
			f.Println("\nOperation will background, get status via $ %s", status)

			return nil
		}

		return fmt.Errorf("Unable to provision resource: %q, check status via: %s", name, status)
	}

	f.Println("Successfully provisioned the resource: %q", name)

	return nil
}
