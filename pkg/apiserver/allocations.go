/*
 * Copyright (C) 2019 Rohith Jayawardene <gambol99@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package apiserver

import (
	restful "github.com/emicklei/go-restful"
)

// listAllocations returns a list of the teams in the allocation
func (u teamHandler) listAllocations(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		return nil
	})
}

// setAllocation is responsible for updating the allocations
func (u teamHandler) setAllocation(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		return nil
	})
}

// updateAllocation updates the allocation to include the team if required
func (u teamHandler) updateAllocation(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		return nil
	})
}

// deleteAllocationsWithTeam removes any allocations from the team
func (u teamHandler) deleteAllocationsWithTeam(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		return nil
	})
}

// deleteAllocations removes a team from the binding allocation
func (u teamHandler) deleteAllocations(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		return nil
	})
}
