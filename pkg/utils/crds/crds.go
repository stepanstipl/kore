/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 * 
 * This file is part of hub-apiserver.
 * 
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * 
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * 
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */


package crds

import (
	"context"
	"fmt"
	"time"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

// NewExtentionsAPIClient returns an extensions api client
func NewExtentionsAPIClient(cfg *rest.Config) (client.Interface, error) {
	return client.NewForConfig(cfg)
}

// ApplyCustomResourceDefinitions s responsible for applying a collection of CRDs
func ApplyCustomResourceDefinitions(c client.Interface, list []*apiextensions.CustomResourceDefinition) error {
	for _, crd := range list {
		if err := ApplyCustomResourceDefinition(c, crd); err != nil {
			return err
		}
	}

	return nil
}

// ApplyCustomResourceDefinition is responsible for applying the CRD to the cluster
func ApplyCustomResourceDefinition(c client.Interface, crd *apiextensions.CustomResourceDefinition) error {
	// @step: retrieve the current if there
	err := func() error {
		current, err := c.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crd.Name, metav1.GetOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				return err
			}

			_, err := c.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)

			return err
		}

		crd.SetGeneration(current.GetGeneration())
		crd.SetResourceVersion(current.GetResourceVersion())

		_, err = c.ApiextensionsV1beta1().CustomResourceDefinitions().Update(crd)

		return err
	}()
	if err != nil {
		return err
	}

	return CheckCustomResourceDefinition(c, crd)
}

// CheckCustomResourceDefinition ensures the CRD is ok to go
func CheckCustomResourceDefinition(c client.Interface, crd *apiextensions.CustomResourceDefinition) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	doneCh := make(chan struct{}, 0)

	go func() {
		for {
			select {
			case _ = <-ctx.Done():
				return
			default:
			}

			err := func() error {
				// @step: ensure the crd has been registered
				obj, err := c.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crd.Name, metav1.GetOptions{})
				if err != nil {
					return err
				}
				if len(obj.Status.Conditions) < 2 {
					return fmt.Errorf("waiting for crd conditions to reach 2")
				}
				for _, x := range obj.Status.Conditions {
					if x.Status != "True" {
						return fmt.Errorf("condition not met, reason: %s", x.Reason)
					}
				}
				time.Sleep(100 * time.Millisecond)

				return nil
			}()
			if err == nil {
				doneCh <- struct{}{}
				return
			}
		}
	}()

	select {
	case <-doneCh:
		return nil
	case <-time.After(10 * time.Second):
		return fmt.Errorf("failed to register the crd")
	}
}
