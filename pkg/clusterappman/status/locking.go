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

package status

import (
	"fmt"
	"sync"

	kcore "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/utils/kubernetes"
)

// Keep track of the status of the clsuterappman for a given cluster i.e.:
// - the health of application controller
// - the health of the helm controller
type appmanStatus struct {
	// is the app controller for this apps cluster working
	appControllerStatus bool
	// koreClusterAppManComponents are the dependencies of all cluster apps
	clusterAppManComponents kcore.Components
}

var (
	mu = &sync.Mutex{}
	// clusterAppControl the health by specific cluster
	clusterAppManStatusByCluster = map[kubernetes.KubernetesAPI]appmanStatus{}
	// koreClusterKubeAPICfg KubernetesAPI details as provided as flags when starting the kore-apiserver
	koreClusterKubeAPICfg *kubernetes.KubernetesAPI
)

// GetLocalKoreClusterAPI will get the API server details used by the kore-apiserver
func GetLocalKoreClusterAPI() *kubernetes.KubernetesAPI {
	mu.Lock()
	defer mu.Unlock()
	return koreClusterKubeAPICfg
}

// SetLocalKoreClusterAPI sets the API details used for clusterapps in the Kore cluster
func SetLocalKoreClusterAPI(apiCfg kubernetes.KubernetesAPI) {
	mu.Lock()
	defer mu.Unlock()
	koreClusterKubeAPICfg = &apiCfg
}

// SetLocalKoreClusterComponents saves the clusterappman status for the local cluster
func SetLocalKoreClusterComponents(components kcore.Components) error {
	kubeCfg := GetLocalKoreClusterAPI()
	if kubeCfg == nil {
		return fmt.Errorf("the local kubernetes cluster config has not been set yet cannot set the status for unknown cluster api")
	}
	SetAppManComponents(components, *kubeCfg)

	return nil
}

// SetAppManComponents will set the clusterappman component status for a given cluster
func SetAppManComponents(clusterAppManConponents kcore.Components, kubeCfg kubernetes.KubernetesAPI) {
	defer mu.Unlock()
	myAppmanStatus, ok := clusterAppManStatusByCluster[kubeCfg]
	if !ok {
		myAppmanStatus = appmanStatus{
			clusterAppManComponents: clusterAppManConponents,
		}
	}
	clusterAppManStatusByCluster[kubeCfg] = myAppmanStatus
}

// GetAppControllerStatus gets the application controller status for a given cluster
func GetAppControllerStatus(kubeCfg kubernetes.KubernetesAPI) bool {
	mu.Lock()
	defer mu.Unlock()
	myAppmanStatus, ok := clusterAppManStatusByCluster[kubeCfg]
	if !ok {
		return false
	}
	return myAppmanStatus.appControllerStatus
}

// SetAppControllerStatus will update the application controller status for a given cluster
func SetAppControllerStatus(s bool, kubeCfg kubernetes.KubernetesAPI) {
	mu.Lock()
	defer mu.Unlock()
	myAppmanStatus, ok := clusterAppManStatusByCluster[kubeCfg]
	if !ok {
		myAppmanStatus = appmanStatus{
			appControllerStatus: s,
		}
	}
	clusterAppManStatusByCluster[kubeCfg] = myAppmanStatus
}
