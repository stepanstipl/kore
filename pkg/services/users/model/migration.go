/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package model

import (
	"github.com/jinzhu/gorm"
)

// Migrations performs the migrations
func Migrations(db *gorm.DB) error {
	db.AutoMigrate(AuditEvent{})
	db.AutoMigrate(Identity{})
	db.AutoMigrate(Invitation{})
	db.AutoMigrate(Member{})
	db.AutoMigrate(Team{})
	db.AutoMigrate(User{})

	db.Model(&User{}).
		AddIndex("idx_users_name", "username")

	db.Model(&Team{}).
		AddIndex("idx_teams_name", "name")

	db.Model(&Member{}).
		AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT").
		AddForeignKey("team_id", "teams(id)", "CASCADE", "RESTRICT")

	db.Model(&Identity{}).
		AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT").
		AddIndex("idx_identity_provider")

	db.AutoMigrate(Invitation{}).
		AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT").
		AddForeignKey("team_id", "teams(id)", "CASCADE", "RESTRICT")

	return nil
}
