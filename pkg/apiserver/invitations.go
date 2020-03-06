/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package apiserver

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/appvia/kore/pkg/kore"
	restful "github.com/emicklei/go-restful"
)

// invitationSubmit is called to handle the submission of a generated link from the UI
func (u teamHandler) invitationSubmit(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()
		token := req.PathParameter("token")

		if err := u.Invitations().HandleGenerateLink(ctx, token); err != nil {
			return err
		}
		resp.WriteHeader(http.StatusOK)

		return nil
	})
}

// inviteLinkByUser is responsible for generating a link for a specific user
func (u teamHandler) inviteLinkByUser(req *restful.Request, resp *restful.Response) {
	u.inviteLink(req, resp)
}

// inviteLink is responsible for generating a link
func (u teamHandler) inviteLink(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()
		team := req.PathParameter("team")
		user := req.PathParameter("user")

		duration, err := parseInvitationExpiry(req)
		if err != nil {
			return err
		}
		options := kore.GenerateLinkOptions{Duration: duration, User: user}

		if u.Config().PublicHubURL == "" {
			return errors.New("An invitation URL can not be generated, as the Kore UI public URL is not set (ui-public-url)")
		}

		token, err := u.Teams().Team(team).Members().GenerateLink(ctx, options)
		if err != nil {
			return err
		}
		uri := fmt.Sprintf("%s/teams/invitation/%s", u.BaseURI(), token)

		return resp.WriteHeaderAndEntity(http.StatusOK, uri)
	})
}

// inviteUser create an team invitation
func (u teamHandler) inviteUser(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()
		team := req.PathParameter("team")
		user := req.PathParameter("user")

		duration, err := parseInvitationExpiry(req)
		if err != nil {
			return err
		}
		options := kore.InvitationOptions{Duration: duration}

		if err := u.Teams().Team(team).Members().Invite(ctx, user, options); err != nil {
			return err
		}
		resp.WriteHeader(http.StatusOK)

		return nil
	})
}

// removeInvite removes a user invitation for the team
func (u teamHandler) removeInvite(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()
		team := req.PathParameter("team")
		user := req.PathParameter("user")

		if err := u.Teams().Team(team).Members().DeleteInvitation(ctx, user); err != nil {
			return err
		}
		resp.WriteHeader(http.StatusOK)

		return nil
	})
}

// listInvites returns a list of invitations
func (u teamHandler) listInvites(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Members().ListInvitations(ctx)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

func parseInvitationExpiry(req *restful.Request) (time.Duration, error) {
	// @TODO this is a bit shit, as the default value should be passed through
	expiry := req.QueryParameter("expire")
	if expiry == "" {
		expiry = "1h"
	}

	duration, err := time.ParseDuration(expiry)
	if err != nil {
		return 0, fmt.Errorf("invalid expire: '%s' for invitation", duration)
	}

	return duration, nil
}
