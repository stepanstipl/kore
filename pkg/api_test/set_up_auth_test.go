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

package api_test

import (
	"encoding/base64"
	"fmt"
	"time"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	. "github.com/onsi/gomega"
)

var jwtPubKey *rsa.PublicKey
var jwtPrivKey *rsa.PrivateKey

func getAuthBuiltInAdmin() runtime.ClientAuthInfoWriter {
	return httptransport.BearerToken("password")
}

func getAuthAnon() runtime.ClientAuthInfoWriter {
	return httptransport.PassThroughAuth
}

func getAuthAdmin() runtime.ClientAuthInfoWriter {
	return httptransport.BearerToken(getJWT(testUserAdmin, testUserAdmin+emailSuffix))
}

func getAuthTeam1Member() runtime.ClientAuthInfoWriter {
	return httptransport.BearerToken(getJWT(testUserTeam1, testUserTeam1+emailSuffix))
}

func getAuthTeam2Member() runtime.ClientAuthInfoWriter {
	return httptransport.BearerToken(getJWT(testUserTeam2, testUserTeam2+emailSuffix))
}

func getAuthMultiTeamMember() runtime.ClientAuthInfoWriter {
	return httptransport.BearerToken(getJWT(testUserMultiTeam, testUserMultiTeam+emailSuffix))
}

func getJWT(username string, email string) string {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		jwt.StandardClaims
	}{
		email,
		username,
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

func setupJWT() {
	var err error
	pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKey)
	Expect(err).ToNot(HaveOccurred())
	pubKey, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	Expect(err).ToNot(HaveOccurred())
	jwtPubKey = pubKey.(*rsa.PublicKey)
	// Priv key must match the above pub key:
	privKeyBytes, err := base64.StdEncoding.DecodeString(privKey)
	Expect(err).ToNot(HaveOccurred())
	jwtPrivKey, err = x509.ParsePKCS1PrivateKey(privKeyBytes)
	Expect(err).ToNot(HaveOccurred())
}

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
