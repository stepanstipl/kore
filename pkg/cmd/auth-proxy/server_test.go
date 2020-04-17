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

package authproxy_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/appvia/kore/pkg/utils/openid/openidfakes"

	authproxyfakes "github.com/appvia/kore/pkg/cmd/auth-proxy/auth-proxyfakes"

	log "github.com/sirupsen/logrus"

	authproxy "github.com/appvia/kore/pkg/cmd/auth-proxy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Server", func() {
	var authProxy authproxy.Interface
	var verifier *authproxyfakes.FakeVerifier
	var config authproxy.Config
	var createErr, runErr error
	var k8sAPI, idpServer *ghttp.Server
	var upstreamAuthTokenFile *os.File
	var allowedIPs []string

	BeforeEach(func() {
		verifier = &authproxyfakes.FakeVerifier{}
		createErr = nil
		runErr = nil
		allowedIPs = []string{"0.0.0.0/0"}

		k8sAPI = ghttp.NewServer()
		k8sAPI.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/hello"),
				ghttp.RespondWith(200, "hello"),
			),
		)

		idpServer = ghttp.NewServer()
		idpDiscovery := map[string]string{
			"issuer": idpServer.URL(),
		}
		idpServer.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/.well-known/openid-configuration"),
				ghttp.RespondWithJSONEncoded(200, idpDiscovery),
			),
		)

		var err error
		upstreamAuthTokenFile, err = ioutil.TempFile(os.TempDir(), "kore-test-*")
		Expect(err).ToNot(HaveOccurred())
	})

	JustBeforeEach(func() {
		logger := log.New()
		logger.Out = GinkgoWriter

		config = authproxy.Config{
			Listen:                     "127.0.0.1:0",
			MetricsListen:              "127.0.0.1:0",
			IDPClientID:                "idp-client-id",
			IDPServerURL:               idpServer.URL(),
			IDPUserClaims:              []string{"name"},
			AllowedIPs:                 allowedIPs,
			UpstreamURL:                k8sAPI.URL(),
			UpstreamAuthorizationToken: upstreamAuthTokenFile.Name(),
		}

		authProxy, createErr = authproxy.New(logger, config, []authproxy.Verifier{verifier})
		if createErr == nil {
			runErr = authProxy.Run(context.Background())
		}
	})

	AfterEach(func() {
		if authProxy != nil {
			_ = authProxy.Stop()
		}
		k8sAPI.Close()
		idpServer.Close()

		_ = os.Remove(upstreamAuthTokenFile.Name())
	})

	Context("with invalid configuration", func() {
		When("allowedIPs is empty", func() {
			BeforeEach(func() {
				allowedIPs = []string{}
			})
			It("should return an error", func() {
				Expect(createErr).To(MatchError("allowed IPs must be set"))
			})
		})
		When("allowedIPs contains an invalid value", func() {
			BeforeEach(func() {
				allowedIPs = []string{"invalid value"}
			})
			It("should return an error", func() {
				Expect(createErr).To(MatchError("invalid CIDR notation: \"invalid value\""))
			})
		})
	})

	Context("with valid configuration", func() {
		BeforeEach(func() {
			idToken := &openidfakes.FakeIDToken{}
			idToken.ClaimsStub = func(v interface{}) error {
				return json.Unmarshal([]byte(`{"name":"testUser"}`), v)
			}
			verifier.AdmitReturns(true, nil)
		})
		JustBeforeEach(func() {
			Expect(createErr).ToNot(HaveOccurred())
			Expect(runErr).ToNot(HaveOccurred())
		})

		When("the client IP is allowed", func() {
			It("should allow it", func() {
				body, statusCode, err := makeGetRequest("http://" + authProxy.Addr() + "/hello")
				Expect(err).ToNot(HaveOccurred())
				Expect(statusCode).To(Equal(http.StatusOK))
				Expect(body).To(Equal("hello"))
			})
		})

		When("the client IP is disallowed", func() {
			BeforeEach(func() {
				allowedIPs = []string{"1.2.3.4/32"}
			})
			It("should return 403", func() {
				body, statusCode, err := makeGetRequest("http://" + authProxy.Addr() + "/hello")
				Expect(err).ToNot(HaveOccurred())
				Expect(statusCode).To(Equal(http.StatusForbidden))
				Expect(strings.TrimSpace(body)).To(Equal("Forbidden"))
			})
		})

		When("the verifier forbids the request", func() {
			BeforeEach(func() {
				verifier.AdmitReturns(false, nil)
			})
			It("should return 403", func() {
				_, statusCode, err := makeGetRequest("http://" + authProxy.Addr() + "/hello")
				Expect(err).ToNot(HaveOccurred())
				Expect(statusCode).To(Equal(http.StatusForbidden))
			})
		})
	})
})

func makeGetRequest(url string) (string, int, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Authorization", "Bearer test-token")
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	return string(body), resp.StatusCode, nil
}
