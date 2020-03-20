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

package apitest

import (
	"github.com/appvia/kore/pkg/testing/apiclient/operations"
	"github.com/appvia/kore/pkg/testing/apimodels"
	. "github.com/onsi/gomega"
)

// setupTeamsAndUsers ensures that the teams and users used by this test suite all exist in the API.
func setupTeamsAndUsers() {
	if err := ensureUserExists(testUserAdmin, testUserAdmin+emailSuffix); err != nil {
		Expect(err).ToNot(HaveOccurred())
	}
	if err := ensureUserInTeam("kore-admin", testUserAdmin); err != nil {
		Expect(err).ToNot(HaveOccurred())
	}

	if err := ensureUserExists(testUserTeam1, testUserTeam1+emailSuffix); err != nil {
		Expect(err).ToNot(HaveOccurred())
	}
	if err := ensureUserExists(testUserTeam2, testUserTeam2+emailSuffix); err != nil {
		Expect(err).ToNot(HaveOccurred())
	}
	if err := ensureUserExists(testUserMultiTeam, testUserMultiTeam+emailSuffix); err != nil {
		Expect(err).ToNot(HaveOccurred())
	}
	if err := ensureTeamExists(testTeam1, "Test team 1 for API testing", "API Test Team 1"); err != nil {
		Expect(err).ToNot(HaveOccurred())
	}
	if err := ensureTeamExists(testTeam2, "Test team 2 for API testing", "API Test Team 2"); err != nil {
		Expect(err).ToNot(HaveOccurred())
	}
	if err := ensureUserInTeam(testTeam1, testUserTeam1); err != nil {
		Expect(err).ToNot(HaveOccurred())
	}
	if err := ensureUserInTeam(testTeam2, testUserTeam2); err != nil {
		Expect(err).ToNot(HaveOccurred())
	}
	if err := ensureUserInTeam(testTeam1, testUserMultiTeam); err != nil {
		Expect(err).ToNot(HaveOccurred())
	}
	if err := ensureUserInTeam(testTeam2, testUserMultiTeam); err != nil {
		Expect(err).ToNot(HaveOccurred())
	}
}

func ensureTeamExists(name string, description string, summary string) error {
	a := getApi()
	_, err := a.Operations.GetTeam(operations.NewGetTeamParams().WithTeam(name), getAuthBuiltInAdmin())
	if err != nil {
		switch err := err.(type) {
		case *operations.GetTeamNotFound:
			team := &apimodels.V1Team{
				Metadata: &apimodels.V1ObjectMeta{
					Name: name,
				},
				Spec: &apimodels.V1TeamSpec{
					Description: &description,
					Summary:     &summary,
				},
			}
			if _, err := a.Operations.UpdateTeam(operations.NewUpdateTeamParams().WithTeam(name).WithBody(team), getAuthBuiltInAdmin()); err != nil {
				return err
			}
		default:
			return err
		}
	}
	return nil
}

func ensureUserExists(username string, email string) error {
	a := getApi()
	disabled := false
	_, err := a.Operations.GetUser(operations.NewGetUserParams().WithUser(username), getAuthBuiltInAdmin())
	if err != nil {
		switch err := err.(type) {
		case *operations.GetUserNotFound:
			user := &apimodels.V1User{
				Metadata: &apimodels.V1ObjectMeta{
					Name: username,
				},
				Spec: &apimodels.V1UserSpec{
					Disabled: &disabled,
					Username: &username,
					Email:    &email,
				},
			}
			if _, err := a.Operations.UpdateUser(operations.NewUpdateUserParams().WithUser(username).WithBody(user), getAuthBuiltInAdmin()); err != nil {
				return err
			}
		default:
			return err
		}
	}
	return nil

}

func ensureUserInTeam(team string, username string) error {
	a := getApi()
	resp, err := a.Operations.ListUserTeams(operations.NewListUserTeamsParams().WithUser(username), getAuthBuiltInAdmin())
	if err != nil {
		return err
	}
	// See if user already in team
	for _, ut := range resp.Payload.Items {
		if ut.Metadata.Name == team {
			return nil
		}
	}
	addTeamMemberParams := operations.NewAddTeamMemberParams().
		WithTeam(team).
		WithUser(username).
		WithBody(&apimodels.V1TeamMember{
			Spec: &apimodels.V1TeamMemberSpec{
				Team:     &team,
				Username: &username,
				Roles:    []string{},
			},
		})

	if _, err := a.Operations.AddTeamMember(addTeamMemberParams, getAuthBuiltInAdmin()); err != nil {
		return err
	}
	return nil
}
