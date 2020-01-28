/*
 * Copyright (C) 2019 Appvia Ltd. <gambol99@gmail.com>
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

package hub

import (
	"context"
	"fmt"
	"time"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// IDPClientMaxRetries specifies the maximum times to try connecting to an IDP provider
	IDPClientMaxRetries int = 50
	// IDPClientBackOff is the time to wait before retrying
	IDPClientBackOff time.Duration = 5 * time.Second
	// DefaultIDP is the name to use for the IDP that is used as "default"
	DefaultIDP = "default"
)

// IDP is the hub api idp interface
type IDP interface {
	// Get returns the idp from the hub
	Get(context.Context, string) (*corev1.IDP, error)
	// Default gets the default identity provider
	Default(context.Context) (*corev1.IDP, error)
	// Exists checks if the idp provider exists
	Exists(context.Context, string) (bool, error)
	// List returns a list of configured idps
	List(context.Context) (*corev1.IDPList, error)
	// ConfigTypes returns the idp config types available
	ConfigTypes(context.Context) ([]*corev1.IDPConfig, error)
	// Update is responsible for updating / creating a confgured IDP
	Update(context.Context, *corev1.IDP) error
	// UpdateClient is responsible for updating / creating an idp client
	UpdateClient(context.Context, *corev1.IDPClient) error
}

// authImpl provides the implementation for Identity providers
type idpImpl struct {
	Interface
}

// Default returns the default identity configuration
func (a idpImpl) Default(ctx context.Context) (idp *corev1.IDP, err error) {
	return a.Get(ctx, DefaultIDP)
}

// Exists checks if an identity provider exists
func (a idpImpl) Exists(ctx context.Context, name string) (bool, error) {
	if _, err := a.Get(ctx, name); err != nil {
		if err == ErrNotFound {
			return false, nil
		}
		return false, nil
	}
	return true, nil
}

// Get returns the identity configurations in the hub
func (a idpImpl) Get(ctx context.Context, name string) (*corev1.IDP, error) {
	// Get from DEX...
	idp, err := getDEXConnector(a.Config().DEX, name)
	if err != nil {
		return nil, err
	}
	return idp, nil
}

// List returns a list of identity providers
func (a idpImpl) List(ctx context.Context) (*corev1.IDPList, error) {
	items, err := getDEXConnectors(a.Config().DEX)
	if err != nil {
		return nil, err
	}
	return &corev1.IDPList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "IDPList",
		},
		Items: items,
	}, nil
}

// ConfigTypes provides sample identity configurations
// should swagger be used by UX to know this?
func (a idpImpl) ConfigTypes(ctx context.Context) ([]*corev1.IDPConfig, error) {
	ci := []*corev1.IDPConfig{
		{
			Google: &corev1.GoogleIDP{},
		},
		{
			SAML: &corev1.SAMLIDP{},
		},
		{
			OIDC: &corev1.OIDCIDP{},
		},
		{
			Github: &corev1.GithubIDP{},
		},
	}
	return ci, nil

	// TODO: reflect the different types of configuration with sample values
	//v := reflect.ValueOf(corev1.IDPConfig{}).Elem()
	//for i := 0; i < v.NumField(); i++ {
	// Get an empty instance of the config
	//	c := reflect.Zero(v)
	// now assign an empty copy of the config to it:
	//}
}

// Update is responsible for updating / creating a identity providers
func (a idpImpl) Update(ctx context.Context, idp *corev1.IDP) error {
	// Update DEX
	err := updateDEXConector(a.Config().DEX, idp)
	if err != nil {
		return err
	}
	return nil
}

// UpdateClient is responsible for updating / creating a idp clients
func (a idpImpl) UpdateClient(ctx context.Context, c *corev1.IDPClient) error {
	// Update DEX
	if err := updateDEXClient(a.Config().DEX, c); err != nil {
		return err
	}
	return nil
}

// UpdateUser will create / update a static user in IDP
func (a idpImpl) UpdateUser(ctx context.Context, username string, password string) error {
	if len(password) <= 0 {
		return fmt.Errorf("must set a non 0 length password for user %s", username)
	}
	if err := updateDexUser(a.Config().DEX, username, password); err != nil {
		return err
	}
	return nil
}
