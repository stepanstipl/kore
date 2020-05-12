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

package clusterappman

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	kcore "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/clusterapp"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	rc "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// KoreNamespace is the namespace where the clusterappmanager and repo are deployed
	KoreNamespace string = "kore"
	// ParamsConfigMap provides the customisations for deplpoyments carried out from here
	ParamsConfigMap string = "kore-cluster-config"
	// ParamsConfigKey is the configmap key used to store the parameters
	ParamsConfigKey string = "clusterconfig"
	// StatusCongigMap is the name of the configmap object used to store kore cluster status
	StatusCongigMap string = "kore-cluster-status"
	// StatusConfigMapComponentsKey is the key in the configmap used for all conditions
	StatusConfigMapComponentsKey string = "components"
	appControlNamespsace         string = "application-system"
	fluxNamespace                string = "flux"
	// DefaultAppTimeout is the amount of time allowed if none set
	DefaultAppTimeout = 3 * time.Minute
)

type clusterappmanImpl struct {
	// client is the kubernetes client to use
	Client        kubernetes.Interface
	RuntimeClient rc.Client
	ClusterApps   []clusterapp.Instance
	cfg           *rest.Config
	kubeAPIConfig clusterapp.KubernetesAPI
}

// manifest defines the data types that are required to initialise a clusterapp from
// embebed manifest data.
type manifest struct {
	EmededManifests    []string
	Name               string
	Namespace          string
	EnsureNamespace    bool
	DeployTimeOut      time.Duration
	PreDeleteManifests []string
}

var (
	// cas is the map of clusterapp instances
	cas map[string]*clusterapp.Instance = make(map[string]*clusterapp.Instance)

	// mm defines all the embedded manifests and data required to initialise clusterappman
	mm []manifest = []manifest{
		{
			Name: clusterapp.ClusterAppControllerComponentName,
			EmededManifests: []string{
				"application-controller/application-all.yaml",
			},
			Namespace:       appControlNamespsace,
			EnsureNamespace: false,
			DeployTimeOut:   3 * time.Minute,
			PreDeleteManifests: []string{
				"application-controller/pre-delete.yaml",
			},
		},
		{
			Name: "Helm Chart Operator",
			EmededManifests: []string{
				"flux/application.yaml",
				"flux/crds.yaml",
				"flux/rbac.yaml",
				"flux/deployment.yaml",
			},
			Namespace:       fluxNamespace,
			EnsureNamespace: true,
			DeployTimeOut:   10 * time.Minute, // Allow for time for application kind to be ready
			PreDeleteManifests: []string{
				"captain/delete-captain.yaml",
			},
		},
		{
			Name: "Kore Helm Repository",
			EmededManifests: []string{
				"kore-helm-repo/application-all.yaml",
			},
			Namespace:       KoreNamespace,
			EnsureNamespace: false,
			DeployTimeOut:   10 * time.Minute,
		},
	}
)

// New is responsible for creating the clusterappman server
func New(config Config) (Interface, error) {
	if err := config.IsValid(); err != nil {
		return nil, err
	}
	var client kubernetes.Interface
	cc, cfg, err := clusterapp.GetKubeCfgAndControllerClient(config.Kubernetes)
	if err != nil {
		return nil, fmt.Errorf("failed creating controller-runtime client or config: %s", err)
	}
	client, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed creating kubernetes client: %s", err)
	}
	return &clusterappmanImpl{
		Client:        client,
		kubeAPIConfig: config.Kubernetes,
		RuntimeClient: cc,
	}, nil
}

// Run is responsible for starting the services
func (s clusterappmanImpl) Run(ctx context.Context) error {
	logger := log.WithFields(log.Fields{
		"service": "clusterappman",
	})
	logger.Info("attempting to reconcile the applications incluster")
	// initialise clusterapps and parse all the manifests
	logger.Info("loading manifests")
	if err := LoadAllManifests(s.RuntimeClient, s.kubeAPIConfig); err != nil {
		return fmt.Errorf("failed loading manifests - %s", err)
	}
	ticker := time.NewTicker(45 * time.Second)
	for {
		select {
		case <-ctx.Done():
			logger.Print("exiting as requested")
			return nil
		case <-ticker.C:
			err := s.Deploy(ctx, logger)
			if err != nil {
				logger.Errorf("error deploying clusterapp dependencies - %s", err)
			}
		}
	}
}

func (s clusterappmanImpl) Deploy(ctx context.Context, logger *log.Entry) error {
	// deploy / reconcile the Application kind controller
	// deploy / reconcile the Helm operator
	// deploy / reconcile the Appvia helm repo

	var wg sync.WaitGroup
	// component updates channel
	ch := make(chan *kcore.Component, len(mm))
	var cs = make([]*kcore.Component, len(mm))
	var components kcore.Components = cs
	for i, m := range mm {
		// Get the pre-loaded cluster application
		ca, ok := cas[m.Name]
		if !ok {
			return fmt.Errorf("failed creating all manifests as could not find manifest %s", m.Name)
		}
		deployCtx, cancel := context.WithTimeout(ctx, m.DeployTimeOut)
		// make sure we run this cancel whatevs
		defer cancel()
		wg.Add(1)
		// start a deployment thread:
		logger.Infof("starting to wait for '%s' to become ready", ca.Component.Name)
		go ca.WaitForReadyOrTimeout(deployCtx, ch, &wg)
		// Ensure we capture initial status
		components[i] = ca.Component
	}
	// report initial component status' ahead of https://github.com/appvia/kore/issues/89
	logger.Logger.Infof("saving initial status")
	if err := createStatusConfig(ctx, s.RuntimeClient, components); err != nil {
		return fmt.Errorf("error reporting status: %s", err)
	}
	logger.Infof("waiting for %d cluster apps", len(components))
	wg.Wait()
	logger.Debug("finished waiting for all cluster apps")
	defer close(ch)

	// now gather up all the component slices from channels now the routines are complete...
	for i := range components {
		// get the first component "status"
		logger.Debugf("getting cluster app status %d of %d", i, len(components))
		c := <-ch
		logger.Debugf("received channel update %d for %s", i, c.Name)
		// get the corresponding app
		ca, ok := cas[c.Name]
		// overwrite the component with the received updated "status"
		ca.Component = c
		if !ok {
			return fmt.Errorf("missing application from component -  %s", c.Name)
		}
		if ca.Component.Status == kcore.SuccessStatus {
			logger.Infof("cluster app '%s' is ready", ca.Component.Name)
		} else {
			logger.Errorf(
				"cluster app '%s' has timed out or failed and has status [%s] due to [%s]",
				ca.Component.Name,
				ca.Component.Status,
				ca.Component.Detail,
			)
		}
		cs[i] = ca.Component
	}
	logger.Logger.Infof("saving final status")
	if err := createStatusConfig(ctx, s.RuntimeClient, components); err != nil {
		return fmt.Errorf("error reporting status: %s", err)
	}
	return nil
}

// Stop is responsible for trying to stop services
func (s clusterappmanImpl) Stop(context.Context) error {
	return nil
}

func (s clusterappmanImpl) UpgradeClient(options rc.Options) (err error) {
	// TODO may need to make this thread safe!
	s.RuntimeClient, err = rc.New(s.cfg, options)
	if err != nil {
		return err
	}
	return nil
}
