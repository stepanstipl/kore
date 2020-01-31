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

import (
	"time"
)

// Team defines the hub team model
type Team struct {
	// ID is the unique record id
	ID uint64 `gorm:"primary_key"`
	// CreatedAt is the timestamp of record creation
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp"`
	// Name is the name of the team
	Name string `gorm:"primary_key"`
	// Description is a description for the team
	Description string `gorm:"not null"`
	// Summary is a short liner for the team
	Summary string
}
