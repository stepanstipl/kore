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
	"time"
)

// AlertRule defines the structure for the rule
type AlertRule struct {
	ResourceReference
	// ID is the unique record id
	ID uint64 `gorm:"primary_key"`
	// Team is the team for this alert
	Team *Team `gorm:"foreignkey:TeamID"`
	// TeamID is the remote key to the teams table
	TeamID uint64
	// CreatedAt is the timestamp this scan was performed
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp"`
	// Name is the name of the rule
	Name string `sql:"type:varchar(255);DEFAULT:''"`
	// Severity is the importance of the rule
	Severity string `sql:"type:varchar(32);DEFAULT:''"`
	// Summary provides a short summary of what the rule is checking
	Summary string `sql:"type:varchar(2048);DEFAULT:''"`
	// Source is producer of the rule
	Source string
	// Alerts is a collection of alerts for this rule
	Alerts []Alert `gorm:"foreignkey:RuleID"`
	// Raw holds the raw payload from the alerting event
	RawRule string `sql:"type:varchar(8192);DEFAULT:''"`
}
