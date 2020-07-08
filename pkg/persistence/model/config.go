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

// Config defines a config in the kore
type Config struct {
	ID uint64 `gorm:"primary_key"`
	// Name is the unique record
	Name string
	// Items is the complimented value for a name
	Items []ConfigItems `gorm:"foreignkey:ItemID"`
}

type ConfigItems struct {
	ID     uint64 `gorm:"primary_key"`
	ItemID uint64
	Key    string
	Value  string
}
