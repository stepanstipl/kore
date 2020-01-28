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

package model

import "time"

// Invitation defines an team invitation in the hub
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
