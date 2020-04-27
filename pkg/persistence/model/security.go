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

// SecurityScanResult stores the result of a security scan
type SecurityScanResult struct {
	// ID is the unique record id
	ID uint64 `gorm:"primary_key"`
	// ResourceAPIVersion is the group and version of the resource scanned by this scan
	ResourceAPIVersion string
	// ResourceKind is the kind of the resource scanned by this scan
	ResourceKind string
	// ResourceNamespace is the namespace of the resource scanned by this scan
	ResourceNamespace string
	// ResourceName is the name of the resource scanned by this scan
	ResourceName string
	// OwningTeam is the name of the Kore team that owns this resource, will be empty if it is a non-team resource.
	OwningTeam string
	// CheckedAt is the timestamp this scan was performed
	CheckedAt time.Time `sql:"DEFAULT:current_timestamp"`
	// ArchivedAt is the timestamp this scan was superceded by a new scan
	ArchivedAt time.Time `sql:"DEFAULT:null"`
	// OverallStatus is the overall status of the scan
	OverallStatus string
	// Results is the set of security rule results for this scan
	Results []SecurityRuleResult `gorm:"foreignkey:ScanID"`
}

// SecurityRuleResult stores the result of a specific rule when applied during a security scan
type SecurityRuleResult struct {
	// ID is the unique record id
	ID uint64 `gorm:"primary_key"`
	// ScanID is the ID of the scan which contains this result
	ScanID uint64
	// RuleCode identifies the rule to which this result relates
	RuleCode string
	// Status is the compliance of the target with this rule
	Status string
	// Message is any additional information about this result
	Message string
	// CheckedAt is the timestamp this scan was performed
	CheckedAt time.Time `sql:"DEFAULT:current_timestamp"`
}
