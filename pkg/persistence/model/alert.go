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
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

const (
	// AlertStatusActive indicates the alert is active
	AlertStatusActive = "Active"
	// AlertStatusOK indicates status is fine
	AlertStatusOK = "OK"
	// AlertStatusSilenced indicates an silenced status
	AlertStatusSilenced = "Silenced"
)

var (
	// AlertStatuses is a full list of possible statuses
	AlertStatuses = []string{
		AlertStatusActive,
		AlertStatusOK,
		AlertStatusSilenced,
	}
)

// Alert defines the structure for the alert
type Alert struct {
	// UID is unique id for this alert
	UID string `gorm:"not null"`
	// ID is the unique record id
	ID uint64 `gorm:"primary_key"`
	// CreatedAt is the timestamp this alert was performed
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp"`
	// ArchivedAt is the timestamp this alert was superceded by a new alert
	ArchivedAt *time.Time `sql:"DEFAULT:null"`
	// Labels is a collection of labels on the alert
	Labels []AlertLabel `gorm:"foreignkey:AlertID"`
	// Expiration is the experation time of the expiry
	Expiration *time.Time
	// Fingerprint is a unique code for this alert from producer
	Fingerprint string
	// Rule is the rule the alert is associated to
	Rule *AlertRule `gorm:"foreignkey:RuleID"`
	// RuleID is the id of the rule which triggered the event
	RuleID uint64
	// RawAlert holds the raw payload from the alerting event
	RawAlert string `sql:"type:varchar(8192);DEFAULT:''"`
	// Summary is a summary for the alert - this the generate alert
	Summary string `sql:"type:varchar(2048);DEFAULT:''"`
	// Status is the overall status of the alert
	Status string `gorm:"not null"`
	// StatusMessage provides a location to place status related message
	// i.e. if the alert has been silenced a user message will be placed here
	StatusMessage string `sql:"DEFAULT:''"`
}

// GetLabels returns the labels in an array
func (a *Alert) GetLabels() []string {
	var list []string
	for _, x := range a.Labels {
		list = append(list, fmt.Sprintf("%s=%s", x.Name, x.Value))
	}

	return list
}

// HasLabels checks all the labels exist
func (a *Alert) HasLabels(labels map[string]string) bool {
	for k, v := range labels {
		if !a.HasLabel(k, v) {
			return false
		}
	}

	return true
}

// HasLabel checks if the label exists
func (a *Alert) HasLabel(key, value string) bool {
	for _, x := range a.Labels {
		if x.Name == key && x.Value == value {
			return true
		}
	}

	return false
}

// BeforeCreate is called pre creation
func (a *Alert) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("UID", uuid.NewV4().String())

	if a.Status == "" {
		scope.SetColumn("Status", AlertStatusOK)
	}

	return nil
}
