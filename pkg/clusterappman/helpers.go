/*
 * Copyright (C) 2019 Appvia Ltd <info@appvia.io>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package clusterappman

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	kcore "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/clusterapp"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	yaml "github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// makeKubernetesConfig returns a rest.Config from the options
func makeKubernetesConfig(config KubernetesAPI) (*rest.Config, error) {
	// @step: are we creating an in-cluster kubernetes client
	if config.InCluster {
		return rest.InClusterConfig()
	}

	if config.KubeConfig != "" {
		return clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	}

	return &rest.Config{
		Host:        config.MasterAPIURL,
		BearerToken: config.Token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: config.SkipTLSVerify,
		},
	}, nil
}

// LoadAllManifests will load all the manifests defined here
// This provides a simple testable entrypoint
func LoadAllManifests(cc client.Client) error {
	for _, m := range mm {
		ca, err := getClusterAppFromEmbeddedManifests(m, cc)
		log.Infof("loading manifest for cluster app - %s", ca.Component.Name)
		if err != nil {
			return fmt.Errorf("failed to load %s manifests: %s", m.Name, err)
		}
		log.Debugf("manifests loaded for cluster app - %s", ca.Component.Name)
		cas[ca.Component.Name] = &ca
	}
	return nil
}

func getClusterAppFromEmbeddedManifests(m manifest, cc client.Client) (clusterapp.Instance, error) {
	// for all the embeded paths specified...
	resfiles := make([]http.File, 0)
	for _, manifestFile := range m.EmededManifests {
		file, err := Manifests.Open(manifestFile)
		if err != nil {
			return clusterapp.Instance{}, err
		}
		resfiles = append(resfiles, file)
	}
	return clusterapp.NewAppFromManifestFiles(cc, m.Name, resfiles)
}

func ensureNamespace(ctx context.Context, cc client.Client, name string) error {
	return kubernetes.EnsureNamespace(ctx, cc, &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	})
}

// GetStatus returns the status of all compoents deployed by ClusterAppMan
func GetStatus(ctx context.Context, cc client.Client) (components *kcore.Components, err error) {
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      StatusCongigMap,
			Namespace: KoreNamespace,
		},
	}
	exists, err := kubernetes.GetIfExists(ctx, cc, cm)
	if err != nil {
		return nil, fmt.Errorf(
			"error obtaining config for %s/%s - %s",
			KoreNamespace,
			StatusCongigMap,
			err,
		)
	}
	if exists {
		components := &kcore.Components{}
		// get the data from the configmap data key
		b := []byte(cm.Data[StatusConfigMapComponentsKey])
		if err := yaml.Unmarshal(b, components); err != nil {
			return nil, fmt.Errorf(
				"unable to deserialize yaml data from config map: %s/%s, key: %s - %s",
				KoreNamespace,
				StatusCongigMap,
				StatusConfigMapComponentsKey,
				err,
			)
		}
		return components, nil
	}
	return nil, fmt.Errorf(
		"missing configmap %s in namespace %s",
		StatusCongigMap,
		KoreNamespace,
	)
}

// createStatusConfig creates a configmap for configuring the kore cluster manager
func createStatusConfig(ctx context.Context, cc client.Client, componants kcore.Components) error {
	b, err := yaml.Marshal(componants)
	if err != nil {
		return errors.New("can not marshall componants into yaml")
	}
	// Specify the parameters in the config map
	cm := &core.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      StatusCongigMap,
			Namespace: KoreNamespace,
		},
		Data: (map[string]string{StatusConfigMapComponentsKey: string(b)}),
	}
	if _, err := kubernetes.CreateOrUpdateConfigMap(ctx, cc, cm); err != nil {
		return err
	}
	return nil
}
