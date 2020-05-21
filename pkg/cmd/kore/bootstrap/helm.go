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
	"io/ioutil"
	"os"
	"time"

	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	"github.com/appvia/kore/pkg/version"

	"github.com/ghodss/yaml"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/kube"
)

// hasHelmRelease checks if the helm release exists
func hasHelmRelease(ctx context.Context, actionConfig *action.Configuration, releaseName string) (*chart.Chart, error) {
	list, err := action.NewList(actionConfig).Run()
	if err != nil {
		return nil, err
	}

	for _, x := range list {
		if x.Name == releaseName {
			return x.Chart, nil
		}
	}

	return nil, nil
}

// HelmDeploy is responsible for deploying
func HelmDeploy(ctx context.Context, chart *chart.Chart, values map[string]interface{}, name, namespace, context string) error {
	kubeconfig, err := kubernetes.GetOrCreateKubeConfig()
	if err != nil {
		return err
	}

	config := kube.GetConfig(kubeconfig, context, namespace)
	logger := func(format string, v ...interface{}) {
		//fmt.Printf(format, v...)
	}

	ac := &action.Configuration{}
	if err := ac.Init(config, name, os.Getenv("HELM_DRIVER"), logger); err != nil {
		return err
	}

	found, err := hasHelmRelease(ctx, ac, name)
	if err != nil {
		return err
	}

	switch found != nil {
	case true:
		action := action.NewUpgrade(ac)
		action.MaxHistory = 3
		action.Namespace = namespace
		action.Timeout = 5 * time.Minute
		action.Atomic = true
		action.Install = true
		action.Force = true
		action.Wait = true

		_, err = action.Run(name, chart, values)

	default:
		action := action.NewInstall(ac)
		action.CreateNamespace = true
		action.IsUpgrade = true
		action.Namespace = namespace
		action.ReleaseName = name
		action.Replace = true
		action.Timeout = 5 * time.Minute
		action.Wait = true

		_, err = action.Run(chart, values)
	}

	return err
}

// GetHelmValues returns returns or prompts for the values
func GetHelmValues(path string) (map[string]interface{}, error) {
	found, err := utils.FileExists(path)
	if err != nil {
		return nil, err
	}
	// @TODO we should probably check the params in the values file
	if !found {
		values := GetDefaultHelmValues()

		a := authInfoConfig{}

		if err := (&cmdutil.Prompts{
			&cmdutil.Prompt{Id: "Client ID", ErrMsg: "%s cannot be blank", Value: &a.ClientID},
			&cmdutil.Prompt{Id: "Client Secret", ErrMsg: "%s cannot be blank", Value: &a.ClientSecret},
			&cmdutil.Prompt{Id: "Authorization Endpoint", ErrMsg: "%s cannot be blank", Value: &a.AuthorizeURL},
		}).Collect(); err != nil {
			return nil, err
		}

		values["idp"] = map[string]interface{}{
			"client_id":     a.ClientID,
			"client_secret": a.ClientSecret,
			"server_url":    a.AuthorizeURL,
		}

		return values, nil
	}

	// @step: we read in the values.yml
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	values := make(map[string]interface{})
	if err := yaml.Unmarshal(content, &values); err != nil {
		return nil, err
	}

	return values, nil
}

// GetDefaultHelmValues returns the default values for the chart
func GetDefaultHelmValues() map[string]interface{} {
	return map[string]interface{}{
		"api": map[string]interface{}{
			"feature_gates": []string{"services=true"},
			"hostPort":      10080,
			"replicas":      1,
			"serviceType":   "NodePort",
			"version":       version.Release,
		},
		"ui": map[string]interface{}{
			"feature_gates": []string{"services=true"},
			"hostPort":      3000,
			"replicas":      1,
			"serviceType":   "NodePort",
			"version":       version.Release,
		},
	}
}
