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
