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
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"

	kschema "github.com/appvia/kore/pkg/schema"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// New returns or creates a default client
func New() (k8s.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return k8s.NewForConfig(config)
}

// NewClientFromSecret creates a client from the secret
func NewClientFromSecret(secret *core.Secret) (k8s.Interface, error) {
	endpoint := string(secret.Data["endpoint"])
	token := string(secret.Data["token"])
	ca := string(secret.Data["ca.crt"])

	return NewFromToken(endpoint, token, ca)
}

// NewClientFromConfigSecret creates a client from the secret
func NewClientFromConfigSecret(secret *configv1.Secret) (k8s.Interface, error) {
	endpoint := secret.Spec.Data["endpoint"]
	token := secret.Spec.Data["token"]
	ca := secret.Spec.Data["ca.crt"]

	return NewFromToken(endpoint, token, ca)
}

// NewRuntimeClientFromConfigSecret creates and returns a runtime client from configv1.Secret
func NewRuntimeClientFromConfigSecret(secret *configv1.Secret) (client.Client, error) {
	cfg := &rest.Config{
		Host:        secret.Spec.Data["endpoint"],
		BearerToken: secret.Spec.Data["token"],
		TLSClientConfig: rest.TLSClientConfig{
			CAData: []byte(secret.Spec.Data["ca.crt"]),
		},
	}

	return NewRuntimeClient(cfg)
}

// NewRuntimeClientFromSecret creates and returns a runtime client
func NewRuntimeClientFromSecret(secret *core.Secret) (client.Client, error) {
	cfg := &rest.Config{
		Host:        string(secret.Data["endpoint"]),
		BearerToken: string(secret.Data["token"]),
		TLSClientConfig: rest.TLSClientConfig{
			CAData: secret.Data["ca.crt"],
		},
	}

	return NewRuntimeClient(cfg)
}

// NewRuntimeClientForAPI will return a new controller-runtime client for the given Kubernetes API
func NewRuntimeClientForAPI(k KubernetesAPI) (client.Client, error) {
	cfg, err := MakeKubernetesConfig(k)
	if err != nil {
		return nil, fmt.Errorf("failed creating kubernetes config: %s", err)
	}

	return NewRuntimeClient(cfg)
}

func NewRuntimeClient(cfg *rest.Config) (client.Client, error) {
	mapper, err := apiutil.NewDynamicRESTMapper(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create mapper for client: %w", err)
	}

	c, err := client.New(cfg, client.Options{
		Scheme: kschema.GetScheme(),
		Mapper: mapper,
	})
	if err != nil {
		return nil, fmt.Errorf("failed creating kubernetes runtime client: %s", err)
	}

	return c, nil
}

// NewFromToken creates a kubernetes client from a endpoint and token
func NewFromToken(endpoint, token, ca string) (k8s.Interface, error) {
	return k8s.NewForConfig(&rest.Config{
		BearerToken: token,
		Host:        endpoint,
		TLSClientConfig: rest.TLSClientConfig{
			CAData:   []byte(ca),
			Insecure: (len(ca) <= 0),
		},
	})
}

// NewGKEClient returns a kube api client for gke clusters
func NewGKEClient(account, endpoint string) (k8s.Interface, error) {
	scopes := []string{
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/userinfo.email",
	}

	token, err := google.JWTConfigFromJSON([]byte(account), scopes...)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(endpoint, "https") {
		endpoint = "https://" + endpoint
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return k8s.NewForConfig(&rest.Config{
		Host: endpoint,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
		WrapTransport: func(rt http.RoundTripper) http.RoundTripper {
			return &oauth2.Transport{
				Source: token.TokenSource(context.Background()),
				Base:   tr,
			}
		},
	})
}

// WaitOnKubeAPI waits for the kube-apiserver to become available
func WaitOnKubeAPI(ctx context.Context, client k8s.Interface, interval, timeout time.Duration) error {
	return wait.PollImmediate(interval, timeout, func() (bool, error) {
		healthStatus := 0
		client.Discovery().RESTClient().Get().AbsPath("/healthz").Do(ctx).StatusCode(&healthStatus)

		if healthStatus != http.StatusOK {
			return false, nil
		}

		return true, nil
	})
}

// HasGroup checks if an api group exists
func HasGroup(version metav1.GroupVersionResource) (bool, error) {
	client, err := New()
	if err != nil {
		return false, err
	}

	list, err := client.Discovery().ServerGroups()
	if err != nil {
		return false, err
	}

	for _, x := range list.Groups {
		if len(x.Versions) <= 0 {
			continue
		}
		for _, v := range x.Versions {
			if v.GroupVersion == version.String() {
				return true, nil
			}
		}
	}

	return false, nil
}

// MakeKubernetesConfig returns a rest.Config from the options
func MakeKubernetesConfig(config KubernetesAPI) (*rest.Config, error) {
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
