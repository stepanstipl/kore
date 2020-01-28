/**
 * Copyright (C) 2020 Rohith Jayawardene <info@appvia.io>
 *
 * This file is part of hub-apiserver.
 *
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package version

import "github.com/urfave/cli"

var (
	// Authors is a collection of authors - should probably generate from git
	Authors = []cli.Author{
		{
			Name:  "Rohith Jayawardene",
			Email: "rohith.jayawardene@appvia.io",
		},
		{
			Name:  "Lewis Marshal",
			Email: "lewis.marshal@appvia.io",
		},
	}
)
