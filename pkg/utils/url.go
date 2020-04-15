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

package utils

import "regexp"

var (
	// urlRegex is the regex we use to validate endpoints
	urlRegex = regexp.MustCompile(`^https?://([0-9a-zA-Z\.]+)|([0-9]{1,3}\.){3,3}[0-9]{1,3}(:[0-9]+)?$`)
)

// URLRegex returns the above
func URLRegex() regexp.Regexp {
	return *urlRegex
}

// IsValidURL checks if the url is valid
func IsValidURL(v string) bool {
	return urlRegex.MatchString(v)
}
