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

package jsonschema

import (
	"net"
	"time"

	"github.com/xeipuuv/gojsonschema"
)

func init() {
	gojsonschema.FormatCheckers.Add("hh:mm", hourMinuteChecker{})
	gojsonschema.FormatCheckers.Add("1.2.3.4/16", cidrChecker{})
}

type hourMinuteChecker struct{}

// IsFormat returns true if the input is a time in "hh:mm" format
func (f hourMinuteChecker) IsFormat(input interface{}) bool {
	val, ok := input.(string)
	if !ok {
		return false
	}

	_, err := time.Parse("03:04", val)
	return err == nil
}

type cidrChecker struct{}

// IsFormat returns true if the input is a CIDR notation
func (f cidrChecker) IsFormat(input interface{}) bool {
	val, ok := input.(string)
	if !ok {
		return false
	}

	_, _, err := net.ParseCIDR(val)
	return err == nil
}
