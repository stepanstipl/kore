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

import "time"

// Invitation defines an team invitation in the kore
type Invitation struct {
	Model
	// Expires the is time the invitation expires
	Expires time.Time `gorm:"not null"`
	// Team is the team the user is a member of
	Team *Team
	// TeamID is the remote key to the teams table
	TeamID uint64 `gorm:"unique_index:idx_invitation_user_team"`
	// User is
	User *User `gorm:"foreignkey:UserID"`
	// UserID is the remote key to the users table
	UserID uint64 `gorm:"unique_index:idx_invitation_user_team"`
}
