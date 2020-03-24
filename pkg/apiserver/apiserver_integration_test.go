// +build integration

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

package apiserver_test

import (
	"testing"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/appvia/kore/pkg/apiclient"
	"github.com/appvia/kore/pkg/apiclient/models"
	"github.com/appvia/kore/pkg/apiclient/operations"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAPI(t *testing.T) {
	//generateJWTKeys()
	RegisterFailHandler(Fail)
	BeforeSuite(func() {
		loadTestData()
		setupJWT()
		setupTeamsAndUsers()
	})
	RunSpecs(t, "API")
}

// These reference users defined in testdata/integration_test_data.json
type testUsers string

const (
	TestUserAdmin     testUsers = "admin"
	TestUserTeam1     testUsers = "team1user"
	TestUserTeam2     testUsers = "team2user"
	TestUserMultiTeam testUsers = "multiTeamUser"
)

// These reference teams defined in testdata/integration_test_data.json
type testTeams string

const (
	TestTeam1 testTeams = "team1"
	TestTeam2 testTeams = "team2"
)

// intTestData contains the data used by tests in this suite, loaded from testdata/integration_test_data.json
var intTestData testData

// getApi returns an instance of the API client.
func getApi() *apiclient.AppviaKore {
	uri := os.Getenv("KORE_API_HOST")
	if uri == "" {
		uri = "localhost"
	}
	transport := httptransport.New(uri+":10080", "", nil)
	return apiclient.New(transport, strfmt.Default)
}

// getAuthBuiltInAdmin gets a token for accessing the API as the built-in admin user.
func getAuthBuiltInAdmin() runtime.ClientAuthInfoWriter {
	return httptransport.BearerToken("password")
}

// getAuthAnon accesses the API as a non-logged-in user.
func getAuthAnon() runtime.ClientAuthInfoWriter {
	return httptransport.PassThroughAuth
}

// getAuth gets a token for accessing the API as the specified user.
func getAuth(user testUsers) runtime.ClientAuthInfoWriter {
	return httptransport.BearerToken(getJWT(intTestData.Users[string(user)]))
}

type testData struct {
	Auth struct {
		PubKey  string `json:"pubKey"`
		PrivKey string `json:"privKey"`
	} `json:"auth"`
	Teams map[string]*testTeam `json:"teams"`
	Users map[string]*testUser `json:"users"`
}
type testUser struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Teams    []string `json:"teams"`
}
type testTeam struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Summary     string `json:"summary"`
}

var jwtPubKey *rsa.PublicKey
var jwtPrivKey *rsa.PrivateKey

// loadTestData loads the test data from testdata/integration_test_data.json into intTestData.
func loadTestData() {
	file, err := ioutil.ReadFile("testdata/integration_test_data.json")
	if err != nil {
		Expect(err).ToNot(HaveOccurred())
	}
	err = json.Unmarshal(file, &intTestData)
	if err != nil {
		Expect(err).ToNot(HaveOccurred())
	}
}

// setupTeamsAndUsers ensures that the teams and users used by this test suite all exist in the API.
func setupTeamsAndUsers() {
	for team := range intTestData.Teams {
		if err := ensureTeamExists(intTestData.Teams[team]); err != nil {
			Expect(err).ToNot(HaveOccurred())
		}
	}
	for user := range intTestData.Users {
		if err := ensureUserExists(intTestData.Users[user]); err != nil {
			Expect(err).ToNot(HaveOccurred())
		}
		for _, userTeam := range intTestData.Users[user].Teams {
			if err := ensureUserInTeam(userTeam, intTestData.Users[user].Username); err != nil {
				Expect(err).ToNot(HaveOccurred())
			}
		}
	}
}

// setupJWT sets up the tokens for JWT-based authentication.
func setupJWT() {
	var err error
	pubKeyBytes, err := base64.StdEncoding.DecodeString(intTestData.Auth.PubKey)
	Expect(err).ToNot(HaveOccurred())
	pubKey, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	Expect(err).ToNot(HaveOccurred())
	jwtPubKey = pubKey.(*rsa.PublicKey)
	// Priv key must match the above pub key:
	privKeyBytes, err := base64.StdEncoding.DecodeString(intTestData.Auth.PrivKey)
	Expect(err).ToNot(HaveOccurred())
	jwtPrivKey, err = x509.ParsePKCS1PrivateKey(privKeyBytes)
	Expect(err).ToNot(HaveOccurred())
}

// generateJWTKeys generates a new RSA keypair and prints it to the console.
func generateJWTKeys() error {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	privKeyBytes := x509.MarshalPKCS1PrivateKey(key)
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return err
	}
	fmt.Println("Private key: ", base64.StdEncoding.EncodeToString(privKeyBytes))
	fmt.Println("Public key: ", base64.StdEncoding.EncodeToString(pubKeyBytes))
	return nil
}

func ensureTeamExists(testTeam *testTeam) error {
	a := getApi()
	_, err := a.Operations.GetTeam(operations.NewGetTeamParams().WithTeam(testTeam.Name), getAuthBuiltInAdmin())
	if err != nil {
		switch err := err.(type) {
		case *operations.GetTeamNotFound:
			team := &models.V1Team{
				Metadata: &models.V1ObjectMeta{
					Name: testTeam.Name,
				},
				Spec: &models.V1TeamSpec{
					Description: &testTeam.Description,
					Summary:     &testTeam.Summary,
				},
			}
			if _, err := a.Operations.UpdateTeam(operations.NewUpdateTeamParams().WithTeam(testTeam.Name).WithBody(team), getAuthBuiltInAdmin()); err != nil {
				return err
			}
		default:
			return err
		}
	}
	return nil
}

func ensureUserExists(testUser *testUser) error {
	a := getApi()
	disabled := false
	_, err := a.Operations.GetUser(operations.NewGetUserParams().WithUser(testUser.Username), getAuthBuiltInAdmin())
	if err != nil {
		switch err := err.(type) {
		case *operations.GetUserNotFound:
			user := &models.V1User{
				Metadata: &models.V1ObjectMeta{
					Name: testUser.Username,
				},
				Spec: &models.V1UserSpec{
					Disabled: &disabled,
					Username: &testUser.Username,
					Email:    &testUser.Email,
				},
			}
			if _, err := a.Operations.UpdateUser(operations.NewUpdateUserParams().WithUser(testUser.Username).WithBody(user), getAuthBuiltInAdmin()); err != nil {
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
		WithBody(&models.V1TeamMember{
			Spec: &models.V1TeamMemberSpec{
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

func getJWT(user *testUser) string {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		jwt.StandardClaims
	}{
		user.Email,
		user.Username,
		jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := token.SignedString(jwtPrivKey)
	if err != nil {
		Expect(err).ToNot(HaveOccurred())
	}
	return tokenStr
}
