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

package v1

// TeamAssetType defines the type of a team asset
type TeamAssetType string

// @note: IMPORTANT - ensure any updates here are also reflected in the persistence model
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

func (t TeamAssetType) String() string {
	return string(t)
}
