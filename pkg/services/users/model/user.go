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

// User defines a user in the kore
type User struct {
	// ID is the unique record id
	ID uint64 `gorm:"primary_key"`
	// CreatedAt is the timestamp of record creation
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp"`
	// Username is the unique record id
	Username string `gorm:"primary_key"`
	// Email is the user email address
	Email string `json:"email,omitempty"`
	// Disabled indicates the user is disabled
	Disabled bool `sql:"DEFAULT:false"`
}
