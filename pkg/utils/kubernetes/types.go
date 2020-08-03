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

package kubernetes

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	koreschema "github.com/appvia/kore/pkg/schema"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

// Object is a Kubernetes object
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Object
type Object interface {
	runtime.Object
	metav1.Object
}

// ObjectWithStatus is a Kubernetes object where you can set/get the status and manage the status components
type ObjectWithStatus interface {
	Object
	GetStatus() (status corev1.Status, message string)
	SetStatus(status corev1.Status, message string)
}

type ObjectWithStatusComponents interface {
	Object
	StatusComponents() *corev1.Components
}

// DependentReference is an object reference to a dependent object in the same namespace
type DependentReference struct {
	// API version of the dependent
	APIVersion string `json:"apiVersion"`
	// Kind of the dependent
	Kind string `json:"kind"`
	// Name of the dependent
	Name string `json:"name"`
}

func (d DependentReference) String() string {
	return fmt.Sprintf("%s/%s/%s", d.APIVersion, d.Kind, d.Name)
}

// DependentReferenceFromObject creates a DependentReference from the given object
func DependentReferenceFromObject(o Object) DependentReference {
	return DependentReference{
		APIVersion: o.GetObjectKind().GroupVersionKind().GroupVersion().String(),
		Kind:       o.GetObjectKind().GroupVersionKind().Kind,
		Name:       o.GetName(),
	}
}

// KubernetesAPI is the configuration for the kubernetes api
type KubernetesAPI struct {
	// InCluster indicates we are running in cluster
	InCluster bool
	// MasterAPIURL specifies the kube-apiserver url
	MasterAPIURL string
	// Token is kubernetes token to authenticate to the api
	Token string
	// KubeConfig is the kubeconfig path
	KubeConfig string
	// SkipTLSVerify indicates we skip tls
	SkipTLSVerify bool
}

type Objects []runtime.Object

// MarshalYAML encodes the objects as a list of YAML objects
func (r Objects) MarshalYAML() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 16384))
	for _, obj := range r {
		yamlData, err := yaml.Marshal(obj)
		if err != nil {
			return nil, err
		}
		buf.WriteString("---\n")
		buf.Write(yamlData)
		buf.WriteRune('\n')
	}
	return buf.Bytes(), nil
}

// UnmarshalYAML decodes a YAML manifest into one or more runtime objects
func (r *Objects) UnmarshalYAML(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	documents := regexp.MustCompile("(?m)^---\n").Split(string(data), -1)

	var objects []runtime.Object

	for _, document := range documents {
		if strings.TrimSpace(document) == "" {
			continue
		}

		obj, err := koreschema.DecodeYAML([]byte(document), nil)
		if err != nil {
			return err
		}

		objects = append(objects, obj)
	}

	*r = objects
	return nil
}
