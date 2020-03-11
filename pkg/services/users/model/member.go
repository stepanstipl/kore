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
	"bytes"
	"encoding/json"
	"strings"

	"github.com/jinzhu/gorm"
)

// Member defines a member of a team
type Member struct {
	// UserID is the reference to the user
	UserID uint64 `gorm:"primary_key;auto_increment:false"`
	// User is the user member
	User *User `gorm:"foreignkey:UserID"`
	// TeamID is the team i
	TeamID uint64 `gorm:"primary_key;auto_increment:false" json:"-"`
	// Team is the team the user is a member of
	Team *Team `gorm:"foreignkey:TeamID" json:"-"`
	// Roles is a collection of roles for this user
	Roles []string `gorm:"-" json:"roles"`
	// UserRoles is the json column for the roles
	UserRoles string `json:"-"`
}

// AfterFind
func (m *Member) AfterFind() error {
	if m.UserRoles == "" {
		m.Roles = []string{}

		return nil
	}

	// @step: decode the users roles into the roles
	if err := json.NewDecoder(strings.NewReader(m.UserRoles)).Decode(&m.Roles); err != nil {
		return err
	}

	return nil
}

// BeforeCreate is called before saving
func (m *Member) BeforeCreate(scope *gorm.Scope) error {
	encoded := &bytes.Buffer{}

	if err := json.NewEncoder(encoded).Encode(&m.Roles); err != nil {
		return err
	}
	m.UserRoles = encoded.String()

	return nil
}
