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
	User *User `gorm:"foreignkey:UserID"`
	// UserID a reference to the user this belong to
	UserID uint64 `gorm:"primary_key;auto_increment:false"`
}
