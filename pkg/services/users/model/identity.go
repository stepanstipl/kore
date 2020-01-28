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

const (
	// APIToken provides an internal api token for the user
	APIToken = "api_token"
)

// Identity defines a user identity in the
type Identity struct {
	// Provider is the provider name
	Provider string `gorm:"primary_key;auto_increment:false"`
	// Extras provides additional storage from config
	Extras string `sql:"type:text"`
	// ProviderUsername is the name of this user in the provider
	ProviderUsername string
	// ProviderEmail is the email of this user in the provider
	ProviderEmail string
	// ProviderToken is the token of the user in the provider
	ProviderToken string
	// ProviderUID is the uid of the user in the provider
	ProviderUID string
	// User is the user the identity belong
	User User `gorm:"foreignkey:UserID"`
	// UserID a reference to the user this belong to
	UserID uint64 `gorm:"primary_key;auto_increment:false"`
}
