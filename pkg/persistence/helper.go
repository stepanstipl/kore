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

package persistence

import "github.com/jinzhu/gorm"

// IsNotFound checks if the error is an not found error
func IsNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}

// Preload applys the preloading to the query
func Preload(load []string, db *gorm.DB) *gorm.DB {
	for _, x := range load {
		db = db.Preload(x)
	}

	return db
}
