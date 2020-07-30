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
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	core "github.com/appvia/kore/pkg/apis/core/v1"
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/kore/assets"
	"github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/persistence/model"
	"github.com/appvia/kore/pkg/store"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Setup is called one on initialization and used to provision and empty kore
func (h hubImpl) Setup(ctx context.Context) error {
	log.Info("initializing the kore")

	// @step: ensure the kore namespaces are there
	for _, x := range []string{HubNamespace, HubAdminTeam, HubDefaultTeam, HubSystem, HubOperators} {
		if err := h.ensureNamespace(ctx, x); err != nil {
			return err
		}
	}

	// @step: ensure the default user is there
	for _, x := range []string{HubAdminUser} {
		if err := h.ensureHubAdminUser(ctx, x, "admin@local"); err != nil {
			return err
		}
	}

	// @step: ensure the kore admin team exists
	for _, x := range []string{HubAdminTeam, HubDefaultTeam} {
		if err := h.ensureHubTeam(ctx, x, "Team for "+x); err != nil {
			return err
		}
	}

	// @step: ensure the kore admin user
	for _, x := range []string{HubAdminUser} {
		if err := h.ensureHubAdminMembership(ctx, x, HubAdminTeam); err != nil {
			return err
		}
	}

	// @step: ensure an OIDC client is created in IDP broker
	if h.Config().DEX.EnabledDex {
		if err := h.ensureHubIDPClientExists(ctx); err != nil {
			return err
		}
	}

	for _, clusterProvider := range ClusterProviders() {
		for _, plan := range clusterProvider.DefaultPlans() {
			if err := h.Plans().Update(getAdminContext(ctx), &plan, true); err != nil {
				return err
			}
		}

		if planPolicy := clusterProvider.DefaultPlanPolicy(); planPolicy != nil {
			if err := h.PlanPolicies().Update(getAdminContext(ctx), planPolicy, true); err != nil {
				return err
			}

			allocation := planPolicy.CreateAllocation([]string{"*"})
			if err := h.Teams().Team(HubAdminTeam).Allocations().Update(getAdminContext(ctx), allocation, true); err != nil {
				return err
			}
		}
	}

	for _, x := range assets.GetDefaultClusterRoles() {
		x.Namespace = HubAdminTeam

		found, err := h.Store().Client().Has(ctx,
			store.HasOptions.From(&clustersv1.ManagedClusterRole{}),
			store.HasOptions.InNamespace(HubAdminTeam),
			store.HasOptions.WithName(x.Name),
		)
		if err != nil {
			return err
		}
		if !found {
			if err := h.Store().Client().Create(ctx, store.CreateOptions.From(&x)); err != nil {
				return err
			}
		}
	}

	// @step: ensure we have the Kore cluster definition
	cluster := &clustersv1.Cluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Cluster",
			APIVersion: clustersv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kore",
			Namespace: "kore-admin",
			Annotations: map[string]string{
				AnnotationSystem:   AnnotationValueTrue,
				AnnotationReadOnly: AnnotationValueTrue,
			},
		},
		Spec: clustersv1.ClusterSpec{
			Kind:          "Kore",
			Plan:          "kore",
			Configuration: apiextv1.JSON{Raw: []byte(`{}`)},
		},
	}
	if err := h.Store().Client().Create(ctx, store.CreateOptions.From(cluster)); err != nil {
		if !kerrors.IsAlreadyExists(err) {
			return fmt.Errorf("failed to create Kore cluster object: %w", err)
		}
	}

	kubernetes := &clustersv1.Kubernetes{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Kubernetes",
			APIVersion: clustersv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kore",
			Namespace: "kore-admin",
			Annotations: map[string]string{
				AnnotationSystem: AnnotationValueTrue,
			},
		},
		Spec: clustersv1.KubernetesSpec{
			Cluster: cluster.Ownership(),
		},
	}
	if err := h.Store().Client().Create(ctx, store.CreateOptions.From(kubernetes)); err != nil {
		if !kerrors.IsAlreadyExists(err) {
			return fmt.Errorf("failed to create Kore kubernetes object: %w", err)
		}
	}

	for _, providerFactory := range ServiceProviderFactories() {
		for _, provider := range providerFactory.DefaultProviders() {
			provider.Namespace = HubNamespace
			if provider.Annotations == nil {
				provider.Annotations = map[string]string{}
			}
			provider.Annotations[AnnotationSystem] = AnnotationValueTrue
			provider.Annotations[AnnotationReadOnly] = AnnotationValueTrue
			if err := h.ServiceProviders().Update(getAdminContext(ctx), &provider); err != nil {
				return err
			}
		}
	}

	// migrate the allocation names to be prefixed with the spec.resource name eg. planpoliy-default-gke
	allocationList, err := h.Teams().Team(HubAdminTeam).Allocations().List(getAdminContext(ctx))
	if err != nil {
		return err
	}

	for _, allocation := range allocationList.Items {
		if err != nil {
			return err
		}
		logger := log.WithFields(log.Fields{
			"name": allocation.ObjectMeta.Name,
		})

		if !strings.HasPrefix(allocation.ObjectMeta.Name, strings.ToLower(allocation.Spec.Resource.Kind)) {
			newName := strings.ToLower(allocation.Spec.Resource.Kind) + "-" + allocation.ObjectMeta.Name

			_, err := h.Teams().Team(HubAdminTeam).Allocations().Get(getAdminContext(ctx), newName)
			if err == nil {
				logger.Infof("allocation with name %s already exists, don't migrate to the new name, just delete the old one", newName)

				if _, err := h.Teams().Team(HubAdminTeam).Allocations().Delete(getAdminContext(ctx), allocation.ObjectMeta.Name, true); err != nil {
					logger.Errorf("error deleting allocation: %v", err)
					return err
				}
			}
			if err != nil {
				logger.Infof("migrating allocation to new name: %s", newName)
				newAllocation := allocation.DeepCopy()
				newAllocation.ObjectMeta.Name = newName
				newAllocation.ObjectMeta.SetResourceVersion("")

				if err := h.Teams().Team(HubAdminTeam).Allocations().Update(getAdminContext(ctx), newAllocation, true); err != nil {
					logger.Errorf("error creating allocation: %v", err)
					return err
				}

				if _, err := h.Teams().Team(HubAdminTeam).Allocations().Delete(getAdminContext(ctx), allocation.ObjectMeta.Name, true); err != nil {
					logger.Errorf("error deleting allocation: %v", err)
					return err
				}
			}
		}
	}

	// @migration: ensure we move sso user to have an identity
	if err := h.ensureUserMigration(ctx); err != nil {
		log.WithError(err).Error("trying to migrate sso users")

		return err
	}

	if err := h.ensureTeamsAndClustersIdentified(ctx); err != nil {
		return err
	}

	return nil
}

// ensureUserMigration is called to ensure the users have a sso identity
func (h hubImpl) ensureUserMigration(ctx context.Context) error {
	// @step: we retrieve a list of users in system and ensure they have a sso user identity
	list, err := h.persistenceMgr.Users().List(ctx)
	if err != nil {
		return err
	}
	for _, x := range list {
		if x.Username == HubAdminUser {
			continue
		}
		identities, err := h.persistenceMgr.Identities().List(ctx,
			persistence.Filter.WithUser(x.Username),
		)
		if err != nil {
			if !persistence.IsNotFound(err) {
				return err
			}
		}

		if identities == nil || len(identities) <= 0 {
			// @logic
			// - if no identity has been found on the user we assume this
			//   was a sso user and needs to be copied in
			err := h.persistenceMgr.Identities().Update(ctx, &model.Identity{
				Provider:      IdentitySSO,
				ProviderEmail: x.Email,
				User:          x,
				UserID:        x.ID,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ensureHubAdminMembership ensures the user is there
func (h hubImpl) ensureHubAdminMembership(ctx context.Context, name, team string) error {
	return h.persistenceMgr.Teams().AddUser(ctx, name, team, []string{"member", "admin"})
}

// ensureHubAdminUser ensures the user is there
func (h hubImpl) ensureHubAdminUser(ctx context.Context, name, email string) error {
	logger := log.WithFields(log.Fields{
		"username": name,
	})

	found, err := h.Users().Exists(ctx, name)
	if err != nil {
		return err
	}

	if !found {
		logger.Info("provisioning the default kore team in api")

		err := h.persistenceMgr.Users().Update(ctx, &model.User{Email: email, Username: name})
		if err != nil {
			logger.WithError(err).Error("trying to create admin user")

			return err
		}
	}
	// Add or update user to IDP broker:
	if h.Config().DEX.EnabledDex {
		if err = h.idp.UpdateUser(ctx, name, h.Config().AdminPass); err != nil {
			logger.WithError(err).Error("trying to update idp password")

			return err
		}
	}

	if h.Config().AdminPass != "" {
		user, err := h.persistenceMgr.Users().Get(ctx, name)
		if err != nil {
			return err
		}

		return h.persistenceMgr.Identities().Update(ctx, &model.Identity{
			Provider:      "basicauth",
			ProviderEmail: email,
			ProviderToken: h.Config().AdminPass,
			UserID:        user.ID,
		})
	}

	return nil
}

// ensureHubTeam ensure a kore team exists in kore
func (h hubImpl) ensureHubTeam(ctx context.Context, name, description string) error {
	nc := getAdminContext(ctx)

	log.WithFields(log.Fields{
		"team": name,
	}).Info("provisioning the default kore team in api")

	_, err := h.Teams().Update(nc, &orgv1.Team{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: HubNamespace,
		},
		Spec: orgv1.TeamSpec{
			Description: description,
			Summary:     description,
		},
	})

	return err
}

// ensureNamespace ensures the namespace exists in kore
func (h hubImpl) ensureNamespace(ctx context.Context, namespace string) error {
	found, err := h.Store().Client().Has(ctx,
		store.HasOptions.From(&corev1.Namespace{}),
		store.HasOptions.InNamespace(HubNamespace),
		store.HasOptions.WithName(namespace),
	)
	if err != nil || found {
		return err
	}

	log.WithFields(log.Fields{
		"namespace": namespace,
	}).Info("provisioning the namespace in api")

	// @step: we need to create it
	return h.Store().Client().Create(ctx,
		store.CreateOptions.From(&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}),
	)
}

func (h hubImpl) ensureHubIDPClientExists(ctx context.Context) error {
	for i := 0; i < IDPClientMaxRetries; i++ {
		// Ensure there is a client created in DEX for the API server
		err := h.idp.UpdateClient(ctx, &core.IDPClient{
			Spec: core.IDPClientSpec{
				DisplayName: "The API server OIDC client",
				ID:          h.Config().IDPClientID,
				Secret:      h.Config().IDPClientSecret,
				RedirectURIs: []string{
					h.Config().PublicAPIURL + "/oauth/callback",
				},
			},
		})
		if err != nil {
			if err == ErrServerNotAvailable {
				// loop for now
				time.Sleep(IDPClientBackOff)
				log.Warn("IDP broker not available so waiting")
				continue
			} else {
				return fmt.Errorf("error configuring IDP client for IDP broker")
			}
		}
	}
	log.Info("API server OIDC client configured in IDP broker")

	return nil
}

func (h hubImpl) ensureTeamsAndClustersIdentified(ctx context.Context) error {
	teams, err := h.Teams().List(getAdminContext(ctx))
	if err != nil {
		log.Errorf("error getting team list: %v", err)
		return err
	}

	for _, team := range teams.Items {
		// @step: check team has identifier, assigning if needed
		teamlog := log.WithField("team", team.Name)
		if team.Labels[LabelTeamIdentifier] == "" {
			teamlog.Info("assigning new identifier to team")
			team.Labels[LabelTeamIdentifier], err = h.Teams().Team(team.Name).Assets().EnsureTeamIdentifier(getAdminContext(ctx))
			if err != nil {
				teamlog.Errorf("error assigning identifier to team: %v", err)
				return err
			}
		}

		// @step: check clusters owned by team all have identifiers
		teamClusters, err := h.Teams().Team(team.Name).Clusters().List(getAdminContext(ctx))
		if err != nil {
			teamlog.Errorf("error getting cluster list for team: %v", err)
			return err
		}
		for _, teamCluster := range teamClusters.Items {
			cluster := (&teamCluster).DeepCopy()
			updated := false
			logger := teamlog.WithField("cluster", cluster.Name)

			if cluster.Labels[LabelTeamIdentifier] == "" {
				logger.Info("setting team identifier for cluster")
				if cluster.Labels == nil {
					cluster.Labels = map[string]string{}
				}
				cluster.Labels[LabelTeamIdentifier] = team.Labels[LabelTeamIdentifier]
				updated = true
			}

			if cluster.Labels[LabelClusterIdentifier] == "" {
				logger.Info("assigning cluster identifier")

				if cluster.Labels == nil {
					cluster.Labels = map[string]string{}
				}
				cluster.Labels[LabelClusterIdentifier], err = h.Teams().Team(team.Name).Assets().GenerateAssetIdentifier(ctx, orgv1.TeamAssetTypeCluster, cluster.Name)
				if err != nil {
					logger.Errorf("error generating identifier for team cluster: %v", err)
					return err
				}
				updated = true
			}

			if updated {
				logger.Debugf("persisting cluster after identifiers assigned - team: %s cluster: %s", cluster.Labels[LabelTeamIdentifier], cluster.Labels[LabelClusterIdentifier])
				err = h.Store().Client().Update(
					getAdminContext(ctx),
					store.UpdateOptions.To(cluster),
				)
				if err != nil {
					logger.Errorf("error updating team cluster: %v", err)
					return err
				}
			}
		}
	}
	return nil
}
