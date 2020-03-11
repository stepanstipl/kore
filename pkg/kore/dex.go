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

package kore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	"github.com/dexidp/dex/api"
	"github.com/dexidp/dex/connector/github"
	"github.com/dexidp/dex/connector/oidc"
	"github.com/dexidp/dex/connector/saml"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	status "google.golang.org/grpc/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// googleIssuer provides the known issuer address used for the google specific OIDC IDP
	googleIssuer = "https://accounts.google.com"
)

// ValidateDEXConfig will ensure the DEX configuration is valid
func ValidateDEXConfig(c DEX) error {
	if len(c.GRPCServer) <= 0 {
		return fmt.Errorf("no host configured - GRPC disabled")
	}
	if len(c.PublicURL) <= 0 {
		return fmt.Errorf("no endpoint configured for Public DEX Url")
	}
	return nil
}

func updateDEXConector(c DEX, idp *corev1.IDP) (err error) {
	// Validate IPD config for DEX...
	req, err := createDEXConnectorReq(c, idp)
	if err != nil {
		return fmt.Errorf("invalid idp as can not cnvert to DEX idp: %v", err)
	}
	conn, err := getDEXConnection(c)
	if err != nil {
		return err
	}
	dc := api.NewDexClient(conn)
	// Clean up when done
	defer conn.Close()

	// get existsing connectors so we know if we need to
	// use the create or update api...
	idps, err := getDEXConnectors(c)
	if err != nil {
		return err
	}
	var dexConnector *api.Connector
	for _, anIdp := range idps {
		if anIdp.Name == idp.Name {
			// The connector exists so Update:
			dexConnector, err = dc.UpdateConnector(context.Background(), req)
			if err != nil {
				log.Errorf("failed updating dex connector: %v", err)
				return err
			}
			// idp is a pointer so will be updated
			idp, err = getIDPfromDEXCRD(dexConnector)
			if err != nil {
				return err
			}
			// Required for SA4006: this value of `idp` is never used
			log.Debugf("idp from connector updated %s", idp.Name)
			return nil
		}
	}
	// If here we need to create...
	dexConnector, err = dc.CreateConnector(context.Background(), req)
	if err != nil {
		log.Errorf("failed creating dex connector: %v", err)
		return err
	}
	// idp is a pointer so will be updated
	idp, err = getIDPfromDEXCRD(dexConnector)
	if err != nil {
		return err
	}
	// Required for SA4006: this value of `idp` is never used
	log.Debugf("idp from connector updated %s", idp.Name)
	return nil
}

func updateDexUser(c DEX, username string, password string) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("cannot generate hash of admin password for username %s", username)
	}
	email := fmt.Sprintf("%s@dex.local", username)
	req := &api.UpdatePasswordReq{
		NewUsername: username,
		Email:       email,
		NewHash:     hash,
	}
	// get a DEX gRCP connection:
	conn, err := getDEXConnection(c)
	if err != nil {
		return err
	}
	// get a dex gRPC client
	dc := api.NewDexClient(conn)
	// Clean up when done
	defer conn.Close()
	res, err := dc.UpdatePassword(context.Background(), req)
	if err != nil {
		return fmt.Errorf("cannot update DEX user - %s - %s", username, err)
	}
	if res.NotFound {
		// we need to create the admin user for the first time!
		creq := &api.CreatePasswordReq{
			Password: &api.Password{
				Username: username,
				Email:    email,
				Hash:     hash,
				UserId:   username,
			},
		}
		_, err := dc.CreatePassword(context.Background(), creq)
		if err != nil {
			return fmt.Errorf("cannot create DEX user - %s - %s", username, err)
		}
	}
	return nil
}

// updateDEXClient will try and create or update a client (if the client exists)
func updateDEXClient(c DEX, idpC *corev1.IDPClient) error {
	// get a DEX gRCP connection:
	conn, err := getDEXConnection(c)
	if err != nil {
		return err
	}
	// get a dex gRPC client
	dc := api.NewDexClient(conn)
	// Clean up when done
	defer conn.Close()

	client := &api.Client{
		Secret:       idpC.Spec.Secret,
		RedirectUris: idpC.Spec.RedirectURIs,
		Id:           idpC.Spec.ID,
		Name:         idpC.Spec.DisplayName,
	}
	req := &api.CreateClientReq{
		Client: client,
	}
	// Make grpc call to create a client as there is no List and we can't check
	res, err := dc.CreateClient(context.Background(), req)
	if err != nil {
		return getHubErrors(err)
	}
	if res.AlreadyExists {
		log.Infof("re-creating client %s", idpC.Spec.ID)
		// No way to update a client secret so...
		dReq := &api.DeleteClientReq{
			Id: idpC.Spec.ID,
		}
		_, err := dc.DeleteClient(context.Background(), dReq)
		if err != nil {
			log.Errorf("failed re-creating dex client - delete: %v", err)
			return getHubErrors(err)
		}
		_, err = dc.CreateClient(context.Background(), req)
		if err != nil {
			log.Errorf("failed re-creating dex client: %v", err)
			return getHubErrors(err)
		}
	}
	return nil
}

func getHubErrors(err error) error {
	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.Unavailable:
			return ErrServerNotAvailable
		default:
		}
	}
	return err
}

func getDEXConnectors(c DEX) ([]corev1.IDP, error) {
	dexConnectorsList, err := getRawDEXConnectors(c)
	if err != nil {
		return nil, err
	}
	items := make([]corev1.IDP, 0)
	for _, aCon := range dexConnectorsList {
		log.Debugf("connector - %v", aCon)
		aIdp, err := getIDPfromDEXCRD(aCon)
		if err != nil {
			log.Warnf("problem getting detail from ipd error %v", err)
			// Skip this item
		} else {
			items = append(items, *aIdp)
		}
	}
	return items, nil
}

// getRawDEXConnectors will connect to DEX and return all the connectors without rendering
func getRawDEXConnectors(c DEX) ([]*api.Connector, error) {
	// get a DEX gRCP connection:
	conn, err := getDEXConnection(c)
	if err != nil {
		return nil, err
	}
	// get a dex gRPC client
	dc := api.NewDexClient(conn)
	// Clean up when done
	defer conn.Close()

	ver, err := dc.GetVersion(context.Background(), &api.VersionReq{})
	log.Debugf("dex version - %s", ver.Server)
	if err != nil {
		return nil, fmt.Errorf("failed getting dex version: %s, %v", c.GRPCServer, err)
	}

	// Make grpc call
	dexConnectors, err := dc.ListConnector(context.Background(), &api.ListConnectorReq{})
	log.Debugf("connectors len - %v", len(dexConnectors.Connectors))
	if err != nil {
		return nil, fmt.Errorf("failed getting dex connectors: %s, %v", c.GRPCServer, err)
	}
	return dexConnectors.Connectors, nil
}

func getDEXConnector(c DEX, name string) (*corev1.IDP, error) {
	// DEX doesn't support retrieving one specific IDP provider
	dexConnectors, err := getDEXConnectors(c)
	if err != nil {
		return nil, err
	}
	for _, aCon := range dexConnectors {
		if aCon.Name == name {
			return &aCon, nil
		}
	}
	return nil, ErrNotFound
}

func getIDPfromDEXCRD(c *api.Connector) (*corev1.IDP, error) {
	conf, err := fromDEXConnectorConfig(c)
	if err != nil {
		return nil, fmt.Errorf("error deserializing dex config %v", err)
	}
	name, _ := getNameAndTypeFromConnector(c)
	return &corev1.IDP{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "IDP",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: HubNamespace,
		},
		Spec: corev1.IDPSpec{
			DisplayName: c.Name,
			Config:      *conf,
		},
	}, nil
}

func getDEXConnection(c DEX) (conn *grpc.ClientConn, err error) {
	err = ValidateDEXConfig(c)
	if err != nil {
		log.Errorf("dex grpc client disabled %v", err)
		return nil, err
	}
	hostAndPort := fmt.Sprintf("%s:%d", c.GRPCServer, c.GRPCPort)
	if len(c.GRPCCaCrt) <= 0 {
		// Do not use TLS
		conn, err = grpc.Dial(hostAndPort, grpc.WithInsecure())
		if err != nil {
			return nil, fmt.Errorf("insecure connection dial %s error %v", hostAndPort, err)
		}
	} else {
		// Create a TLS client
		creds, err := credentials.NewClientTLSFromFile(c.GRPCCaCrt, "")
		if err != nil {
			return nil, fmt.Errorf("failure loading dex cert: %v", err)
		}
		conn, err = grpc.Dial(hostAndPort, grpc.WithTransportCredentials(creds))
		if err != nil {
			return nil, fmt.Errorf("TLS connection dial %s error %v", hostAndPort, err)
		}
	}
	return conn, nil
}

// createDEXConnectorConfig just makes this logic more readable above
func createDEXConnectorReq(dexConf DEX, idp *corev1.IDP) (c *api.Connector, err error) {
	c = &api.Connector{}
	c.Name = idp.Spec.DisplayName
	if idp.Spec.Config.Github != nil {
		c.Type = "github"
		c.Id = getConnectorID(idp)
		// Construct a SAML specific DEX config object
		conf := github.Config{
			ClientSecret: idp.Spec.Config.Github.ClientSecret,
			ClientID:     idp.Spec.Config.Github.ClientID,
			RedirectURI:  getDEXCallback(dexConf),
		}
		// TODO: add orgs support here!
		c.Config, err = json.Marshal(conf)
		if err != nil {
			return nil, fmt.Errorf("Error marshaling to connector json: %v", err)
		}
	}
	if idp.Spec.Config.Google != nil {
		// Error if config has already been provided
		if err = getErrIfConfSet(c); err != nil {
			return nil, err
		}
		// Use the generic OIDC provider for now
		// DEX is preparing a Google specific OIDC provider - this could change
		c.Type = "oidc"
		c.Id = getConnectorID(idp)

		// Make the correct config object from the right DEX fields:
		c.Config, err = getOIDCConfig(
			dexConf,
			googleIssuer,
			idp.Spec.Config.Google.ClientID,
			idp.Spec.Config.Google.ClientSecret,
		)
		if err != nil {
			return nil, err
		}
	}
	if idp.Spec.Config.OIDC != nil {
		// Error if config has already been provided
		if err = getErrIfConfSet(c); err != nil {
			return nil, err
		}
		c.Type = "oidc"
		c.Id = getConnectorID(idp)
		c.Config, err = getOIDCConfig(
			dexConf,
			idp.Spec.Config.OIDC.Issuer,
			idp.Spec.Config.Google.ClientID,
			idp.Spec.Config.Google.ClientSecret,
		)
		if err != nil {
			return nil, err
		}
	}
	if idp.Spec.Config.SAML != nil {
		if err = getErrIfConfSet(c); err != nil {
			return nil, err
		}
		c.Type = "saml"
		c.Id = getConnectorID(idp)
		// Construct a SAML specific DEX config object
		conf := saml.Config{
			SSOIssuer:     idp.Spec.Config.SAML.SSOURL,
			CAData:        idp.Spec.Config.SAML.CAData,
			RedirectURI:   getDEXCallback(dexConf),
			UsernameAttr:  idp.Spec.Config.SAML.UsernameAttr,
			EmailAttr:     idp.Spec.Config.SAML.EmailAttr,
			GroupsAttr:    idp.Spec.Config.SAML.GroupsAttr,
			AllowedGroups: idp.Spec.Config.SAML.AllowedGroups,
			GroupsDelim:   idp.Spec.Config.SAML.GroupsDelim,
		}
		c.Config, err = json.Marshal(conf)
		if err != nil {
			return nil, fmt.Errorf("Error marshaling to connector json: %v", err)
		}
	}
	return c, nil
}

// getIDPTypeFromConfig returns the IDP type from the idp config
// note - will not check invalid from here
func getIDPTypeFromConfig(c *corev1.IDPConfig) string {
	if c.Github != nil {
		return "github"
	}
	if c.Google != nil {
		return "google"
	}
	if c.OIDC != nil {
		return "oidc"
	}
	if c.SAML != nil {
		return "saml"
	}
	return ""
}

func fromDEXConnectorConfig(dexC *api.Connector) (*corev1.IDPConfig, error) {
	c := &corev1.IDPConfig{}
	switch dexC.Type {
	case "github":
		ghc := github.Config{}
		if err := json.Unmarshal(dexC.Config, &ghc); err != nil {
			return nil, fmt.Errorf("bad json data for github connector")
		}
		c.Github = &corev1.GithubIDP{
			ClientID:     ghc.ClientID,
			ClientSecret: ghc.ClientSecret,
		}
		// TODO: add orgs support here
	case "oidc":
		oidcConf := oidc.Config{}
		if err := json.Unmarshal(dexC.Config, &oidcConf); err != nil {
			return nil, fmt.Errorf("bad json data for oidc connector")
		}
		if oidcConf.Issuer == googleIssuer {
			c.Google = &corev1.GoogleIDP{
				ClientID:     oidcConf.ClientID,
				ClientSecret: oidcConf.ClientSecret,
				Domains:      oidcConf.HostedDomains,
			}
		} else {
			c.OIDC = &corev1.OIDCIDP{
				ClientID:     oidcConf.ClientID,
				ClientSecret: oidcConf.ClientSecret,
				Issuer:       oidcConf.Issuer,
			}
		}
	case "saml":
		samlConf := saml.Config{}
		if err := json.Unmarshal(dexC.Config, &samlConf); err != nil {
			return nil, fmt.Errorf("bad json data for saml connector")
		}
		c.SAML = &corev1.SAMLIDP{
			SSOURL:        samlConf.SSOURL,
			CAData:        samlConf.CAData,
			UsernameAttr:  samlConf.UsernameAttr,
			EmailAttr:     samlConf.EmailAttr,
			GroupsAttr:    samlConf.GroupsAttr,
			GroupsDelim:   samlConf.GroupsDelim,
			AllowedGroups: samlConf.AllowedGroups,
		}
	default:
		return nil, fmt.Errorf("bad json data - unknown connector type")
	}
	return c, nil
}

func getConnectorID(idp *corev1.IDP) string {
	return fmt.Sprintf("%s-%s", idp.Name, getIDPTypeFromConfig(&idp.Spec.Config))
}

func getNameAndTypeFromConnector(c *api.Connector) (name string, t string) {
	s := strings.Split(c.Id, "-")
	return s[0], s[1]
}

func getErrIfConfSet(c *api.Connector) error {
	if c.Type != "" {
		return errors.New("Abiguouse IDP configurations - only one is supported")
	}
	return nil
}

func getOIDCConfig(
	dexConf DEX,
	issuer string,
	clientID string,
	clientSecret string) ([]byte, error) {
	conf := oidc.Config{
		Issuer:       issuer,
		ClientSecret: clientSecret,
		ClientID:     clientID,
		RedirectURI:  getDEXCallback(dexConf),
	}
	// now create serialized json string
	b, err := json.Marshal(conf)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling to connector json: %v", err)
	}
	return b, nil
}

func getDEXCallback(dexConf DEX) string {
	return fmt.Sprintf("%s/%s", strings.TrimRight(dexConf.PublicURL, "/"), "callback")
}
