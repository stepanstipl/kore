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

package kore

import (
	"context"
	"fmt"
	"strings"

	v1 "github.com/appvia/kore/pkg/apis/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	"github.com/appvia/kore/pkg/utils/validation"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/store"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

// NamespaceClaims is the interface to the class namespace claims
type NamespaceClaims interface {
	// CheckDelete verifies whether the namespace claim can be deleted
	CheckDelete(context.Context, *clustersv1.NamespaceClaim, ...DeleteOptionFunc) error
	// Delete is used to delete a namespace claim in the kore
	Delete(context.Context, string, ...DeleteOptionFunc) (*clustersv1.NamespaceClaim, error)
	// Get returns the class from the kore
	Get(context.Context, string) (*clustersv1.NamespaceClaim, error)
	// List returns a list of classes
	// The optional filter functions can be used to include items only for which all functions return true
	List(context.Context, ...func(clustersv1.NamespaceClaim) bool) (*clustersv1.NamespaceClaimList, error)
	// Has checks if a resource exists
	Has(context.Context, string) (bool, error)
	// Update is responsible for update a namespace claim in the kore
	Update(context.Context, *clustersv1.NamespaceClaim) (*clustersv1.NamespaceClaim, error)
	// CreateForCluster creates a namespace claim for a cluster and namespace, if it doesn't already exist
	CreateForCluster(ctx context.Context, cluster v1.Ownership, clusterNamespace string) error
}

type nsImpl struct {
	*hubImpl
	// team is the team
	team string
}

// CheckDelete verifies whether the cluster can be deleted
func (n *nsImpl) CheckDelete(ctx context.Context, namespaceClaim *clustersv1.NamespaceClaim, o ...DeleteOptionFunc) error {
	opts := ResolveDeleteOptions(o)

	if IsNamespaceNameProtected(namespaceClaim.Spec.Name) {
		return NewErrNotAllowed("the namespace name is a protected or has a protected prefix")
	}

	exists, err := n.Teams().Exists(ctx, namespaceClaim.Spec.Name)
	if err != nil {
		return fmt.Errorf("failed to load team: %w", err)
	}

	if exists {
		return NewErrNotAllowed("it is a team namespace")
	}

	if !opts.Cascade {
		var dependents []kubernetes.DependentReference
		services, err := n.Teams().Team(n.team).Services().List(ctx, func(s servicesv1.Service) bool { return kubernetes.HasOwnerReference(&s, namespaceClaim) })
		if err != nil {
			return fmt.Errorf("failed to list services: %w", err)
		}
		for _, item := range services.Items {
			dependents = append(dependents, kubernetes.DependentReferenceFromObject(&item))
		}

		serviceCredentials, err := n.Teams().Team(n.team).ServiceCredentials().List(ctx, func(sc servicesv1.ServiceCredentials) bool { return kubernetes.HasOwnerReference(&sc, namespaceClaim) })
		if err != nil {
			return fmt.Errorf("failed to list service credentials: %w", err)
		}
		for _, item := range serviceCredentials.Items {
			dependents = append(dependents, kubernetes.DependentReferenceFromObject(&item))
		}

		if len(dependents) > 0 {
			return validation.ErrDependencyViolation{
				Message:    "the following objects need to be deleted first",
				Dependents: dependents,
			}
		}
	}

	return nil
}

// Delete is used to delete a namespace claim in the kore
func (n *nsImpl) Delete(ctx context.Context, name string, o ...DeleteOptionFunc) (*clustersv1.NamespaceClaim, error) {
	opts := ResolveDeleteOptions(o)

	original, err := n.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	if err := opts.Check(original, func(o ...DeleteOptionFunc) error { return n.CheckDelete(ctx, original, o...) }); err != nil {
		return nil, err
	}

	if err := n.Store().Client().Delete(ctx, append(opts.StoreOptions(), store.DeleteOptions.From(original))...); err != nil {
		log.WithError(err).Error("trying to delete the namespace claim")

		return nil, err
	}

	return original, nil
}

// Get returns the class from the kore
func (n *nsImpl) Get(ctx context.Context, name string) (*clustersv1.NamespaceClaim, error) {
	ns := &clustersv1.NamespaceClaim{}

	return ns, n.Store().Client().Get(ctx,
		store.GetOptions.InNamespace(n.team),
		store.GetOptions.InTo(ns),
		store.GetOptions.WithName(name),
	)
}

// List returns a list of classes
func (n *nsImpl) List(ctx context.Context, filters ...func(clustersv1.NamespaceClaim) bool) (*clustersv1.NamespaceClaimList, error) {
	list := &clustersv1.NamespaceClaimList{}

	err := n.Store().Client().List(ctx,
		store.ListOptions.InNamespace(n.team),
		store.ListOptions.InTo(list),
	)
	if err != nil {
		return nil, err
	}

	if len(filters) == 0 {
		return list, nil
	}

	res := []clustersv1.NamespaceClaim{}
	for _, item := range list.Items {
		if func() bool {
			for _, filter := range filters {
				if !filter(item) {
					return false
				}
			}
			return true
		}() {
			res = append(res, item)
		}
	}
	list.Items = res

	return list, nil
}

// Has checks if a resource exists
func (n *nsImpl) Has(ctx context.Context, name string) (bool, error) {
	if _, err := n.Get(ctx, name); err != nil {
		if kerrors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// Update is responsible for update a namespace claim in the kore
func (n *nsImpl) Update(ctx context.Context, namespace *clustersv1.NamespaceClaim) (*clustersv1.NamespaceClaim, error) {
	// @step: ensure it's for cluster we own
	if namespace.Spec.Cluster.Namespace != n.team {
		return nil, ErrNotAllowed{message: "namespace must exist in a cluster you own"}
	}
	namespace.Namespace = n.team

	if IsNamespaceNameProtected(namespace.Spec.Name) {
		return nil, validation.NewError("namespace can not be created").WithFieldError("name", validation.InvalidValue, "is a protected name or has a protected prefix")
	}

	exists, err := n.Teams().Exists(ctx, namespace.Spec.Name)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, validation.NewError("namespace can not be created").WithFieldError("name", validation.InvalidValue, "is a team namespace")
	}

	return namespace, n.Store().Client().Update(ctx,
		store.UpdateOptions.To(namespace),
		store.UpdateOptions.WithCreate(true),
	)
}

func IsNamespaceNameProtected(name string) bool {
	switch {
	case name == "kore" || name == "default":
		return true
	case strings.HasPrefix(name, "kore-"):
		return true
	case strings.HasPrefix(name, "kube-"):
		return true
	default:
		return false
	}
}

// CreateForCluster creates a namespace claim for a cluster and namespace, if it doesn't already exist
func (n *nsImpl) CreateForCluster(ctx context.Context, cluster v1.Ownership, clusterNamespace string) error {
	name := fmt.Sprintf("%s-%s", cluster.Name, clusterNamespace)

	exists, err := n.Has(ctx, name)
	if err != nil || exists {
		return err
	}

	if exists {
		return nil
	}

	namespaceClaim := &clustersv1.NamespaceClaim{
		TypeMeta: metav1.TypeMeta{
			APIVersion: clustersv1.GroupVersion.String(),
			Kind:       "NamespaceClaim",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: n.team,
		},
		Spec: clustersv1.NamespaceClaimSpec{
			Name:    clusterNamespace,
			Cluster: cluster,
		},
	}

	if _, err := n.Update(ctx, namespaceClaim); err != nil {
		return err
	}

	return nil
}
