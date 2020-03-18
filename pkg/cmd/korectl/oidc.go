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

package korectl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Refresh updates our oidc config (in memory)
func (o *OIDC) Refresh() error {
	resp, err := o.refreshRequest()
	if err != nil {
		return err
	}

	token := &RefreshResponse{}

	if err := json.NewDecoder(resp.Body).Decode(token); err != nil {
		return err
	}

	o.IDToken = token.IDToken
	o.AccessToken = token.AccessToken

	return nil
}

// refreshRequestValues builds the values structure for a refresh request
func (o *OIDC) refreshRequestValues() url.Values {
	return url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {o.RefreshToken},
		"client_id":     {o.ClientID},
		"client_secret": {o.ClientSecret},
	}
}

// refreshRequest makes a refresh request
func (o *OIDC) refreshRequest() (*http.Response, error) {
	resp, err := http.PostForm(o.TokenURL, o.refreshRequestValues())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Refresh request failed with code %d", resp.StatusCode)
	}
	return resp, nil
}
