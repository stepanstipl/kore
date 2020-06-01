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

package application

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/kore/assets"

	koreschema "github.com/appvia/kore/pkg/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	appv1beta1 "sigs.k8s.io/application/api/v1beta1"
)

func init() {
	kore.RegisterServiceProviderFactory(Factory{})
}

type Factory struct {
}

func (d Factory) Type() string {
	return Type
}

func (d Factory) JSONSchema() string {
	return ProviderSchema
}

func (d Factory) Create(ctx kore.Context, provider *servicesv1.ServiceProvider) (kore.ServiceProvider, error) {
	manifests, err := assets.Applications.Open("/")
	if err != nil {
		return nil, fmt.Errorf("failed to load application manifests: %w", err)
	}

	dirs, err := manifests.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to load application manifests: %w", err)
	}

	var plans []servicesv1.ServicePlan

	for _, dirInfo := range dirs {
		plan, err := d.createPlan(dirInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to load application %q: %w", dirInfo.Name(), err)
		}
		plans = append(plans, *plan)
	}

	return Provider{name: provider.Name, plans: plans}, nil
}

func (d Factory) SetUp(ctx kore.Context, provider *servicesv1.ServiceProvider) (complete bool, _ error) {
	return true, nil
}

func (d Factory) TearDown(ctx kore.Context, provider *servicesv1.ServiceProvider) (complete bool, _ error) {
	return true, nil
}

func (d Factory) DefaultProviders() []servicesv1.ServiceProvider {
	return []servicesv1.ServiceProvider{
		{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ServiceProvider",
				APIVersion: servicesv1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: Type,
				Annotations: map[string]string{
					kore.AnnotationSystem: "true",
				},
			},
			Spec: servicesv1.ServiceProviderSpec{
				Type:          Type,
				Summary:       "Kubernetes Application provider",
				Description:   "The service provider will deploy one or more Kubernetes resources and an Application type for monitoring purposes",
				Configuration: nil,
			},
		},
	}
}

func (d Factory) createPlan(info os.FileInfo) (*servicesv1.ServicePlan, error) {
	dir, err := assets.Applications.Open("/" + info.Name())
	if err != nil {
		return nil, err
	}

	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var resources []runtime.Object
	var app *appv1beta1.Application

	for _, fileInfo := range files {
		if !strings.HasSuffix(fileInfo.Name(), ".yaml") && !strings.HasSuffix(fileInfo.Name(), ".yml") {
			continue
		}

		file, err := assets.Applications.Open("/" + info.Name() + "/" + fileInfo.Name())
		if err != nil {
			return nil, err
		}

		content, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		splitter := regexp.MustCompile("(?m)^---\n")
		documents := splitter.Split(string(content), -1)

		for _, document := range documents {
			if strings.TrimSpace(document) == "" {
				continue
			}

			obj, err := koreschema.DecodeYAML([]byte(document))
			if err != nil {
				return nil, err
			}

			switch o := obj.(type) {
			case *appv1beta1.Application:
				if app != nil {
					return nil, fmt.Errorf("multiple applications found in %q application", info.Name())
				}
				app = o
			}

			resources = append(resources, obj)
		}
	}

	plan := &servicesv1.ServicePlan{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServicePlan",
			APIVersion: servicesv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceKindApp + "-" + info.Name(),
			Namespace: "kore",
			Annotations: map[string]string{
				kore.AnnotationSystem: "true",
			},
		},
		Spec: servicesv1.ServicePlanSpec{
			Kind:        ServiceKindApp,
			Labels:      nil,
			Description: fmt.Sprintf("%s application", info.Name()),
			Summary:     fmt.Sprintf("%s application", info.Name()),
		},
	}

	if err := plan.Spec.SetConfiguration(AppConfiguration{
		Resources: resources,
	}); err != nil {
		return nil, err
	}

	return plan, nil
}
