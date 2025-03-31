/*
 * Copyright (C) 2025 Anthrove
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package oidc

import (
	"crypto/sha256"
	"fmt"
	"github.com/zitadel/oidc/v3/pkg/op"
	"golang.org/x/text/language"
	"net/http"
)

func NewProvider(storage op.Storage, tenantID string) (*op.Provider, error) {

	//the OpenID Provider requires a 32-byte key for (token) encryption
	//be sure to create a proper crypto random key and manage it securely!
	// TODO only for dev
	key := sha256.Sum256([]byte("test"))

	return op.NewProvider(&op.Config{
		CryptoKey: key,

		//enables code_challenge_method S256 for PKCE (and therefore PKCE in general)
		CodeMethodS256: true,

		//enables additional client_id/client_secret authentication by form post (not only HTTP Basic Auth)
		AuthMethodPost: true,

		//enables additional authentication by using private_key_jwt
		AuthMethodPrivateKeyJWT: true,

		//enables refresh_token grant use
		GrantTypeRefreshToken: true,

		//enables use of the `request` Object parameter
		RequestObjectSupported: true,

		//this example has only static texts (in English), so we'll set the here accordingly
		SupportedUILocales: []language.Tag{language.English},
	}, storage, func(insecure bool) (op.IssuerFromRequest, error) {
		return func(r *http.Request) string {
			scheme := r.URL.Scheme
			if len(scheme) == 0 {
				scheme = "http"
			}

			return fmt.Sprintf("%s://%s/%s", scheme, r.Host, tenantID)
		}, nil
	})

}
