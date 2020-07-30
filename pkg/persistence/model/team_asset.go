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

package model

import (
	"time"
)

// TeamAssetType defines the type of a team asset
type TeamAssetType string

// @note: IMPORTANT, ensure you add any new asset types to the SQL enum definition in the
// TeamAsset struct below as well as here:
const (
	// TeamAssetTypeCluster identifies a cluster asset
	TeamAssetTypeCluster TeamAssetType = "Cluster"
	// TeamAssetTypeNodePool identifies a node pool asset
	TeamAssetTypeNodePool TeamAssetType = "NodePool"
	// TeamAssetTypeNamespace identifies a namespace asset
	TeamAssetTypeNamespace TeamAssetType = "Namespace"
	// TeamAssetTypeCloudService identifies a cloud service (e.g. S3 bucket, RDS instance) asset
	TeamAssetTypeCloudService TeamAssetType = "CloudService"
)

// TeamAsset defines the relationship between a team and an asset (cluster, node pool,
// cloud service, etc) owned by that team - can be used for (e.g.) cost tracking etc.
type TeamAsset struct {
	// AssetIdentifier is the identity of the asset in question
	AssetIdentifier string `sql:"type:char(20)" gorm:"primary_key"`
	// TeamIdentifier is the identity of the owning team
	TeamIdentifier string `sql:"type:char(20)"`
	// TeamAssetType is the type of the asset (e.g. Cluster, CloudService, etc)
	// @note: IMPORTANT - ensure you add any new asset types to the const above as well
	// as here:
	AssetType TeamAssetType `sql:"type:enum('Cluster','NodePool','Namespace','CloudService')"`
	// AssetName is the name of the asset at the point it was created, for reference only
	AssetName string
	// Provider identifies the name of the cloud provider who provides this asset
	Provider string
	// Account identifies the account/project/subscription within the cloud provider for this asset
	Account string
	// CreatedAt is the timestamp of record creation
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp"`
	// DeletedAt is the timestamp the asset was deleted from Kore, null if the asset still exists
	DeletedAt *time.Time `sql:"DEFAULT:null"`
	// AssetCosts represents the costs associated with this asset
	AssetCosts []TeamAssetCost `gorm:"foreignkey:AssetIdentifier"`
}
