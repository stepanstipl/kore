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

package clusterapp

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/appvia/kore/pkg/utils/kubernetes"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	applicationv1beta "sigs.k8s.io/application/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetKubeCfgAndControllerClient will return a new controller-runtime client and the cfg used to create it
func GetKubeCfgAndControllerClient(k KubernetesAPI) (client.Client, *rest.Config, error) {
	cfg, err := makeKubernetesConfig(k)
	if err != nil {
		return nil, cfg, fmt.Errorf("failed creating kubernetes config: %s", err)
	}

	options, err := GetClientOptions()
	if err != nil {
		return nil, cfg, fmt.Errorf("failed getting client options (schemes): %s", err)
	}
	cc, err := client.New(cfg, options)
	if err != nil {
		return nil, cfg, fmt.Errorf("failed creating kubernetes runtime client: %s", err)
	}
	return cc, cfg, nil
}

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

func setMissingNamespace(namespace string, obj runtime.Object) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		log.Debugf("no setting namespace here ->%v<- - %s", obj, err)

		return err
	}
	if accessor.GetNamespace() == "" {
		accessor.SetNamespace(namespace)
		log.Debugf(
			"updated namespace to %s on %s",
			accessor.GetNamespace(),
			accessor.GetName(),
		)
	}

	return nil
}

func ensureNamespace(ctx context.Context, cc client.Client, name string) error {
	return kubernetes.EnsureNamespace(ctx, cc, &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	})
}

func getObjMetaAndSetDefaultNamespace(obj runtime.Object, defaultNamepsace string) metav1.ObjectMeta {
	objMeta := getObjMeta(obj)
	if err := setMissingNamespace(defaultNamepsace, obj); err != nil {
		log.Debugf("error setting namespace for %v - %s", obj, err)
	}
	return objMeta
}

func getObjMeta(obj runtime.Object) metav1.ObjectMeta {
	metaObj := metav1.ObjectMeta{}
	accessor, err := meta.Accessor(obj)
	// TODO: error or not error
	if err != nil {
		if err != nil {
			log.Errorf("error getting metadata for %v - %s", obj, err)
		}
		log.Debugf(
			"got object %s/%s",
			metaObj.Namespace,
			metaObj.Name,
		)
	}
	// TODO: this should be a pointer to the origonal data?
	metaObj.Name = accessor.GetName()
	metaObj.Namespace = accessor.GetNamespace()
	metaObj.Labels = accessor.GetLabels()
	return metaObj
}

func toUnstructuredObj(obj runtime.Object) (*unstructured.Unstructured, error) {
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Version: obj.GetObjectKind().GroupVersionKind().Version,
		Kind:    obj.GetObjectKind().GroupVersionKind().Kind,
		Group:   obj.GetObjectKind().GroupVersionKind().Group,
	})
	objMeta := getObjMeta(obj)
	u.SetName(objMeta.Name)
	u.SetNamespace(objMeta.Namespace)
	u.SetLabels(objMeta.Labels)
	return u, nil
}

func fromUnstructuredApplication(us *unstructured.Unstructured) (*applicationv1beta.Application, error) {
	app := &applicationv1beta.Application{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(us.Object, app); err != nil {
		return nil, err
	}
	return app, nil
}
