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

package server

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/appvia/kore/pkg/persistence"

	v1 "k8s.io/api/core/v1"

	"github.com/appvia/kore/pkg/kore"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/serviceproviders/application"
	korek "github.com/appvia/kore/pkg/utils/kubernetes"
	"github.com/appvia/kore/pkg/version"

	// controller imports
	_ "github.com/appvia/kore/pkg/controllers/register"

	// service provider imports
	_ "github.com/appvia/kore/pkg/serviceproviders/register"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type bootstrapServerImpl struct {
	// config is the server configuration
	config Config
	// server is the real server
	server Interface
}

// NewBootstrap is responsible for creating the bootstrap server
func NewBootstrap(ctx context.Context, config Config) (Interface, error) {
	log.SetReportCaller(false)
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.Println("Starting bootstrap process")

	config.APIServer.Enabled = false
	config.Kore.RunSetup = false
	config.Kore.Controllers = []string{"services", "serviceproviders"}
	config.Kore.CertificateAuthority = "hack/ca/ca.pem"
	config.Kore.CertificateAuthorityKey = "hack/ca/ca-key.pem"
	config.Kore.FeatureGates = map[string]bool{"services": true}

	server, err := New(ctx, config, bootstrapPersistenceManager{})
	if err != nil {
		return nil, err
	}

	return &bootstrapServerImpl{
		config: config,
		server: server,
	}, nil
}

// Run is responsible for starting the services
func (s bootstrapServerImpl) Run(ctx context.Context) error {
	client, err := korek.NewRuntimeClientForAPI(s.config.Kubernetes)
	if err != nil {
		return fmt.Errorf("failed creating runtime client: %s", err)
	}

	if err := korek.EnsureNamespace(ctx, client, &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "kore"}}); err != nil {
		return err
	}

	if err := korek.EnsureNamespace(ctx, client, &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "kore-admin"}}); err != nil {
		return err
	}

	go func() {
		log.Println("starting the kore API server")
		err := s.server.Run(ctx)
		if err != nil {
			log.Errorf("the kore API server failed to start: %s", err.Error())
		}
	}()

	for _, providerFactory := range kore.ServiceProviderFactories() {
		for _, provider := range providerFactory.DefaultProviders() {
			provider.Namespace = kore.HubNamespace

			if _, err := korek.CreateOrUpdate(ctx, client, &provider); err != nil {
				return err
			}
		}
	}

	koreService := &servicesv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: servicesv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kore",
			Namespace: "kore-admin",
		},
		Spec: servicesv1.ServiceSpec{
			Kind: "helm-app",
			Plan: "helm-app-kore",
			Cluster: corev1.Ownership{
				Group:     clustersv1.ClusterGroupVersionKind.Group,
				Version:   clustersv1.ClusterGroupVersionKind.Version,
				Kind:      clustersv1.ClusterGroupVersionKind.Kind,
				Namespace: "kore-admin",
				Name:      "kore",
			},
			ClusterNamespace: "kore",
			Configuration:    nil,
		},
		Status: servicesv1.ServiceStatus{},
	}

	serviceConfig := application.HelmAppConfiguration{
		Source: application.HelmAppSource{
			GitRepository: &application.GitRepository{
				URL:  "https://github.com/appvia/kore",
				Ref:  version.GitSHA,
				Path: "charts/kore",
			},
		},
		ResourceKinds: []metav1.GroupKind{
			{
				Group: "",
				Kind:  "Service",
			},
			{
				Group: "apps",
				Kind:  "Deployment",
			},
		},
	}

	if err := koreService.Spec.SetConfiguration(serviceConfig); err != nil {
		return fmt.Errorf("failed to generate kore service manifest: %w", err)
	}

	exists, err := korek.CheckIfExists(ctx, client, koreService)
	if err != nil {
		return fmt.Errorf("failed to check existing service: %w", err)
	}
	if exists {
		log.Println("kore service already exists")
	} else {
		if err := client.Create(ctx, koreService); err != nil {
			return fmt.Errorf("failed to create Kore service: %w", err)
		}
		log.Println("kore service has been created")
	}

	log.Println("Waiting for kore service to be created")

	timer := time.NewTicker(5 * time.Second)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return nil
	case <-timer.C:
		service := koreService.DeepCopy()
		exists, err := korek.GetIfExists(ctx, client, service)
		if err != nil {
			log.Warnf("failed to check kore service status: %s", err)
		}
		if !exists {
			return fmt.Errorf("kore service has gone away, exiting")
		}

		if service.Status.Status == corev1.SuccessStatus {
			log.Println("kore service has been successfully created")
			return s.Stop(ctx)
		}

		if service.Status.Status == corev1.FailureStatus {
			return fmt.Errorf("kore service has failed: %s", service.Status.Message)
		}

		statusMessage := string(service.Status.Status)
		if service.Status.Message != "" {
			statusMessage = statusMessage + " - " + service.Status.Message
		}
		log.Printf("kore service status is %s\n", statusMessage)
	}

	return nil
}

// Stop is responsible for trying to stop services
func (s bootstrapServerImpl) Stop(ctx context.Context) error {
	return s.server.Stop(ctx)
}

type bootstrapPersistenceManager struct {
}

func (b bootstrapPersistenceManager) Audit() persistence.Audit {
	return nil
}

func (b bootstrapPersistenceManager) Identities() persistence.Identities {
	return nil
}

func (b bootstrapPersistenceManager) Invitations() persistence.Invitations {
	return nil
}

func (b bootstrapPersistenceManager) Members() persistence.Members {
	return nil
}

func (b bootstrapPersistenceManager) Security() persistence.Security {
	return nil
}

func (b bootstrapPersistenceManager) Stop() error {
	return nil
}

func (b bootstrapPersistenceManager) Teams() persistence.Teams {
	return nil
}

func (b bootstrapPersistenceManager) Users() persistence.Users {
	return nil
}
