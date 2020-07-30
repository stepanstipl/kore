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
	db.AutoMigrate(Config{})

	db.Model(&User{}).
		AddIndex("idx_users_name", "username")

	db.Model(&Team{}).
		AddIndex("idx_teams_name", "name")

	fields := []string{
		"resource_kind",
		"resource_group",
		"resource_version",
		"resource_namespace",
		"resource_name",
		"name",
	}

	db.AutoMigrate(&AlertRule{}).
		AddUniqueIndex("idx_alert_rules_unq", fields...).
		AddIndex("idx_alert_rules", fields...).
		AddForeignKey("team_id", "teams(id)", "CASCADE", "RESTRICT")

	db.AutoMigrate(&RuleLabel{}).
		AddForeignKey("rule_id", "alert_rules(id)", "CASCADE", "RESTRICT")

	db.AutoMigrate(&Alert{}).
		AddForeignKey("rule_id", "alert_rules(id)", "CASCADE", "RESTRICT")

	db.AutoMigrate(&AlertLabel{}).
		AddForeignKey("alert_id", "alerts(id)", "CASCADE", "RESTRICT")

	db.Model(&Member{}).
		AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT").
		AddForeignKey("team_id", "teams(id)", "CASCADE", "RESTRICT")

	db.Model(&Identity{}).
		AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")

	db.AutoMigrate(Invitation{}).
		AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT").
		AddForeignKey("team_id", "teams(id)", "CASCADE", "RESTRICT")

	db.AutoMigrate(SecurityScanResult{}).
		AddIndex("idx_scan_identity",
			"resource_kind",
			"resource_group",
			"resource_version",
			"resource_namespace",
			"resource_name",
			"archived_at").
		AddIndex("idx_scan_team",
			"owning_team",
			"archived_at")

	db.AutoMigrate(SecurityRuleResult{}).
		AddForeignKey("scan_id", "security_scan_results(id)", "CASCADE", "RESTRICT")

	db.AutoMigrate(ConfigItems{}).
		AddForeignKey("item_id", "configs(id)", "CASCADE", "RESTRICT")

	db.AutoMigrate(TeamIdentity{})

	db.AutoMigrate(TeamAsset{}).
		AddForeignKey("team_identifier", "team_identity(team_identifier)", "RESTRICT", "RESTRICT").
		AddUniqueIndex("idx_asset", "team_identifier", "asset_identifier").
		AddIndex("idx_assettype", "team_identifier", "asset_type")

	return nil
}
