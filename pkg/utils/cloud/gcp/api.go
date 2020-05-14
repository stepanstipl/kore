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

package gcp

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/appvia/kore/pkg/utils"

	servicemanagement "google.golang.org/api/servicemanagement/v1"
	serviceusage "google.golang.org/api/serviceusage/v1"
)

// EnableAPIs is responsible for enabling a collection of services in the project
func EnableAPIs(ctx context.Context, cc *servicemanagement.APIService, project string, services []string) error {
	for _, x := range services {
		if err := EnableAPI(ctx, cc, project, x); err != nil {
			return fmt.Errorf("trying to enable: %s, error: %s", x, err)
		}
	}

	return nil
}

// ListEnabledAPIs returns a list of enabled services
func ListEnabledAPIs(ctx context.Context, client *serviceusage.Service, project string) ([]string, error) {
	resp, err := client.Services.List("projects/" + project).Filter("state:ENABLED").Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	enabled := make([]string, len(resp.Services))

	for i := 0; i < len(resp.Services); i++ {
		parts := strings.Split(resp.Services[i].Name, "/")

		enabled[i] = parts[len(parts)-1]
	}

	return enabled, nil
}

// EnableAPI is used to enabled a api service in a project
func EnableAPI(ctx context.Context, client *servicemanagement.APIService, project, service string) error {
	request := &servicemanagement.EnableServiceRequest{
		ConsumerId: "project:" + project,
	}

	resp, err := client.Services.Enable(service, request).Context(ctx).Do()
	if err != nil {
		return err
	}

	return utils.WaitUntilComplete(ctx, 3*time.Minute, 5*time.Second, func() (bool, error) {
		status, err := client.Operations.Get(resp.Name).Context(ctx).Do()
		if err != nil {
			return false, nil
		}
		if !status.Done {
			return false, nil
		}
		if status.Error != nil {
			return false, errors.New(status.Error.Message)
		}

		return true, nil
	})
}
