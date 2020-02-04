/*
 * Copyright (C) 2019  Rohith Jayawardene <info@appvia.io>
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

package kubernetes

import (
	"context"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	client k8s.Interface
)

// New returns or creates a default client
func New() (k8s.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return k8s.NewForConfig(config)
}

// NewFromToken creates a kubernetes client from a endpoint and token
func NewFromToken(endpoint, token, ca string) (k8s.Interface, error) {
	return k8s.NewForConfig(&rest.Config{
		BearerToken: token,
		Host:        endpoint,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
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

	return k8s.NewForConfig(&rest.Config{
		Host:            endpoint,
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
		WrapTransport: func(rt http.RoundTripper) http.RoundTripper {
			return &oauth2.Transport{
				Source: token.TokenSource(context.Background()),
			}
		},
	})
}

// WaitOnKubeAPI waits for the kube-apiserver to become available
func WaitOnKubeAPI(ctx context.Context, client k8s.Interface, interval, timeout time.Duration) error {
	return wait.PollImmediate(interval, timeout, func() (bool, error) {
		healthStatus := 0
		client.Discovery().RESTClient().Get().AbsPath("/healthz").Do().StatusCode(&healthStatus)

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
