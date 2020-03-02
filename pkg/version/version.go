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

package version

import (
	"fmt"
	"strconv"
	"time"
)

var (
	// Prog is the name of the product - changes so often best to use a global var
	Prog = "Appvia Kore"
	// Email is the default email
	Email = "info@appvia.io"
	// version in computed version
	version = ""
	// Compiled in the time it was compiling
	Compiled = "0"
	// GitSHA is the sha this was built off
	GitSHA = "no gitsha provided"
	// Release is the releasing version
	Release = "latest"
)

// Version returns the proxy version
func Version() string {
	if version == "" {
		tm, err := strconv.ParseInt(Compiled, 10, 64)
		if err != nil {
			return "unable to parse compiled time"
		}
		version = fmt.Sprintf("%s (git+sha: %s, built: %s)", Release, GitSHA, time.Unix(tm, 0).Format("02-01-2006"))
	}

	return version
}
