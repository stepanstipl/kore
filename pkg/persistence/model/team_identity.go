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

// TeamIdentity records persistently the existence of a team. This continues to exist
// even if the team is subsequently deleted.
type TeamIdentity struct {
	// TeamIdentifier is a globally-unique immutable identifier for a team
	TeamIdentifier string `sql:"type:char(20)" gorm:"PRIMARY_KEY"`
	// TeamName is the name of the team at the point it was created, for reference
	TeamName string
	// CreatedAt is the timestamp of record creation
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp"`
	// DeletedAt is the timestamp the team was deleted from Kore, null if the team still exists
	DeletedAt *time.Time `sql:"DEFAULT:null"`
	// Assets represents the assets associated with this team identity
	Assets []TeamAsset `gorm:"foreignkey:TeamIdentifier"`
	// Costs represents the costs associated with this team identity
	Costs []TeamAssetCost `gorm:"foreignkey:TeamIdentifier"`
}
