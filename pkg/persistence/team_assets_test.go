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

package persistence_test

import (
	"context"
	"time"

	"github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/persistence/model"
	"github.com/appvia/kore/pkg/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Team Assets", func() {
	var store persistence.Interface
	var teamIdentifier string
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
		store = getTestStore()
		teamIdentifier = utils.GenerateIdentifier()
	})

	AfterEach(func() {
		store.Stop()
	})

	Describe("Team Identity", func() {

		Describe("RecordTeamIdentity", func() {
			When("called", func() {
				It("should persist a new team identity", func() {
					err := store.TeamAssets().RecordTeamIdentity(ctx, teamIdentifier, "test team")
					Expect(err).ToNot(HaveOccurred())
					name, err := store.TeamAssets().GetTeamNameForIdentity(ctx, teamIdentifier)
					Expect(err).ToNot(HaveOccurred())
					Expect(name).To(Equal("test team"))
				})
			})
		})

		Describe("MarkTeamIdentityDeleted", func() {
			When("called", func() {
				It("should mark an existing team identity as deleted", func() {
					err := store.TeamAssets().RecordTeamIdentity(ctx, teamIdentifier, "test team")
					Expect(err).ToNot(HaveOccurred())

					// Ensure we do actually have it recorded before deleting:
					name, err := store.TeamAssets().GetTeamNameForIdentity(ctx, teamIdentifier)
					Expect(name).To(Equal("test team"))
					err = store.TeamAssets().MarkTeamIdentityDeleted(ctx, teamIdentifier)
					Expect(err).ToNot(HaveOccurred())

					_, err = store.TeamAssets().GetTeamNameForIdentity(ctx, teamIdentifier)
					Expect(err).To(HaveOccurred())
					Expect(persistence.IsNotFound(err)).To(Equal(true))
				})
			})
		})
	})

	Describe("Assets", func() {
		BeforeEach(func() {
			err := store.TeamAssets().RecordTeamIdentity(ctx, teamIdentifier, "test team")
			Expect(err).ToNot(HaveOccurred())
		})

		Describe("RecordAsset", func() {
			It("should record a new asset", func() {
				assetIdent := utils.GenerateIdentifier()
				err := store.TeamAssets().RecordAsset(ctx, teamIdentifier, assetIdent, model.TeamAssetTypeCluster, "test asset", "gcp")
				Expect(err).ToNot(HaveOccurred())

				asset, err := store.TeamAssets().GetAsset(ctx, teamIdentifier, assetIdent)
				Expect(err).ToNot(HaveOccurred())
				Expect(asset.TeamIdentifier).To(Equal(teamIdentifier))
				Expect(asset.AssetIdentifier).To(Equal(assetIdent))
				Expect(asset.AssetName).To(Equal("test asset"))
				Expect(asset.AssetType).To(Equal(model.TeamAssetTypeCluster))
				Expect(asset.Provider).To(Equal("gcp"))
				Expect(asset.DeletedAt).To(BeNil())
			})
		})

		Describe("MarkAssetDeleted", func() {
			It("should mark an asset as deleted", func() {
				assetIdent := utils.GenerateIdentifier()
				err := store.TeamAssets().RecordAsset(ctx, teamIdentifier, assetIdent, model.TeamAssetTypeCluster, "test asset", "gcp")
				Expect(err).ToNot(HaveOccurred())

				err = store.TeamAssets().MarkAssetDeleted(ctx, teamIdentifier, assetIdent)
				Expect(err).ToNot(HaveOccurred())

				asset, err := store.TeamAssets().GetAsset(ctx, teamIdentifier, assetIdent)
				Expect(err).ToNot(HaveOccurred())
				// DeletedAt should be set:
				Expect(asset.DeletedAt).ToNot(BeNil())
			})
		})

		Describe("MarkAssetUndeleted", func() {
			It("should mark a deleted asset as undeleted", func() {
				assetIdent := utils.GenerateIdentifier()
				err := store.TeamAssets().RecordAsset(ctx, teamIdentifier, assetIdent, model.TeamAssetTypeCluster, "test asset", "gcp")
				Expect(err).ToNot(HaveOccurred())

				err = store.TeamAssets().MarkAssetDeleted(ctx, teamIdentifier, assetIdent)
				Expect(err).ToNot(HaveOccurred())

				err = store.TeamAssets().MarkAssetUndeleted(ctx, teamIdentifier, assetIdent, "new name", "aws")

				asset, err := store.TeamAssets().GetAsset(ctx, teamIdentifier, assetIdent)
				Expect(err).ToNot(HaveOccurred())
				Expect(asset.AssetName).To(Equal("new name"))
				Expect(asset.Provider).To(Equal("aws"))
				Expect(asset.DeletedAt).To(BeNil())
			})
		})

		Describe("ListAssets", func() {
			It("should list all non-deleted assets of a team", func() {
				deletedAssetId := utils.GenerateIdentifier()
				_ = store.TeamAssets().RecordAsset(ctx, teamIdentifier, utils.GenerateIdentifier(), model.TeamAssetTypeCluster, "test asset 1", "gcp")
				_ = store.TeamAssets().RecordAsset(ctx, teamIdentifier, utils.GenerateIdentifier(), model.TeamAssetTypeCluster, "test asset 2", "gcp")
				_ = store.TeamAssets().RecordAsset(ctx, teamIdentifier, deletedAssetId, model.TeamAssetTypeCluster, "test asset 3", "gcp")
				_ = store.TeamAssets().RecordAsset(ctx, teamIdentifier, utils.GenerateIdentifier(), model.TeamAssetTypeCluster, "test asset 4", "gcp")
				_ = store.TeamAssets().MarkAssetDeleted(ctx, teamIdentifier, deletedAssetId)

				assets, err := store.TeamAssets().ListAssets(ctx,
					persistence.TeamAssetFilters.WithTeam(teamIdentifier))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(assets)).To(Equal(3))
			})

			It("should include deleted assets of a team when called with WithDeleted filter", func() {
				deletedAssetId := utils.GenerateIdentifier()
				_ = store.TeamAssets().RecordAsset(ctx, teamIdentifier, utils.GenerateIdentifier(), model.TeamAssetTypeCluster, "test asset 1", "gcp")
				_ = store.TeamAssets().RecordAsset(ctx, teamIdentifier, utils.GenerateIdentifier(), model.TeamAssetTypeCluster, "test asset 2", "gcp")
				_ = store.TeamAssets().RecordAsset(ctx, teamIdentifier, deletedAssetId, model.TeamAssetTypeCluster, "test asset 3", "gcp")
				_ = store.TeamAssets().RecordAsset(ctx, teamIdentifier, utils.GenerateIdentifier(), model.TeamAssetTypeCluster, "test asset 4", "gcp")
				_ = store.TeamAssets().MarkAssetDeleted(ctx, teamIdentifier, deletedAssetId)

				assets, err := store.TeamAssets().ListAssets(ctx,
					persistence.TeamAssetFilters.WithTeam(teamIdentifier),
					persistence.TeamAssetFilters.WithDeleted())
				Expect(err).ToNot(HaveOccurred())
				Expect(len(assets)).To(Equal(4))
			})

			It("should filter by provider", func() {
				_ = store.TeamAssets().RecordAsset(ctx, teamIdentifier, utils.GenerateIdentifier(), model.TeamAssetTypeCluster, "test asset 1", "gcp")
				_ = store.TeamAssets().RecordAsset(ctx, teamIdentifier, utils.GenerateIdentifier(), model.TeamAssetTypeCluster, "test asset 2", "gcp")
				_ = store.TeamAssets().RecordAsset(ctx, teamIdentifier, utils.GenerateIdentifier(), model.TeamAssetTypeCluster, "test asset 3", "aws")
				_ = store.TeamAssets().RecordAsset(ctx, teamIdentifier, utils.GenerateIdentifier(), model.TeamAssetTypeCluster, "test asset 4", "gcp")

				assets, err := store.TeamAssets().ListAssets(ctx,
					persistence.TeamAssetFilters.WithTeam(teamIdentifier),
					persistence.TeamAssetFilters.WithProvider("aws"))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(assets)).To(Equal(1))
				Expect(assets[0].AssetName).To(Equal("test asset 3"))
			})

			It("should filter by asset", func() {
				assetId := utils.GenerateIdentifier()
				_ = store.TeamAssets().RecordAsset(ctx, teamIdentifier, utils.GenerateIdentifier(), model.TeamAssetTypeCluster, "test asset 1", "gcp")
				_ = store.TeamAssets().RecordAsset(ctx, teamIdentifier, utils.GenerateIdentifier(), model.TeamAssetTypeCluster, "test asset 2", "gcp")
				_ = store.TeamAssets().RecordAsset(ctx, teamIdentifier, utils.GenerateIdentifier(), model.TeamAssetTypeCluster, "test asset 3", "aws")
				_ = store.TeamAssets().RecordAsset(ctx, teamIdentifier, assetId, model.TeamAssetTypeCluster, "test asset 4", "gcp")

				assets, err := store.TeamAssets().ListAssets(ctx,
					persistence.TeamAssetFilters.WithTeam(teamIdentifier),
					persistence.TeamAssetFilters.WithAsset(assetId),
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(assets)).To(Equal(1))
				Expect(assets[0].AssetName).To(Equal("test asset 4"))
			})
		})
	})

	Describe("Asset Costs", func() {
		var assetIdentifier string
		BeforeEach(func() {
			err := store.TeamAssets().RecordTeamIdentity(ctx, teamIdentifier, "test team")
			Expect(err).ToNot(HaveOccurred())
			assetIdentifier = utils.GenerateIdentifier()
			err = store.TeamAssets().RecordAsset(ctx, teamIdentifier, assetIdentifier, model.TeamAssetTypeCluster, "test asset 1", "aws")
			Expect(err).ToNot(HaveOccurred())
		})

		Describe("StoreAssetCost", func() {
			It("should store an asset cost", func() {
				err := store.TeamAssets().StoreAssetCost(ctx, &model.TeamAssetCost{
					TeamIdentifier:  teamIdentifier,
					AssetIdentifier: assetIdentifier,
					UsageStartTime:  time.Now(),
					UsageEndTime:    time.Now(),
					UsageType:       "EUW2-EC2-Usage",
					UsageAmount:     1.001,
					UsageUnit:       "hour",
					Cost:            13441,
					Provider:        "aws",
					Account:         "123456123456",
				})
				Expect(err).ToNot(HaveOccurred())

				costs, err := store.TeamAssets().ListCosts(ctx, persistence.TeamAssetFilters.WithAsset(assetIdentifier))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(costs)).To(Equal(1))
				Expect(costs[0].TeamIdentifier).To(Equal(teamIdentifier))
				Expect(costs[0].AssetIdentifier).To(Equal(assetIdentifier))
				Expect(costs[0].UsageType).To(Equal("EUW2-EC2-Usage"))
				Expect(costs[0].UsageAmount).To(Equal(1.001))
				Expect(costs[0].UsageUnit).To(Equal("hour"))
				Expect(costs[0].Cost).To(Equal(int64(13441)))
				Expect(costs[0].Provider).To(Equal("aws"))
				Expect(costs[0].Account).To(Equal("123456123456"))
			})
		})
	})
})
