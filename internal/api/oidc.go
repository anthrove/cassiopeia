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

package api

import "github.com/zitadel/oidc/v3/pkg/op"

func test()  {
	issuer := "https://issuer.example.com"

	op.NewProvider(&op.Config{
		CryptoKey:                         [32]byte{},
		DefaultLogoutRedirectURI:          "",
		CodeMethodS256:                    false,
		AuthMethodPost:                    false,
		AuthMethodPrivateKeyJWT:           false,
		GrantTypeRefreshToken:             false,
		RequestObjectSupported:            false,
		SupportedUILocales:                nil,
		SupportedClaims:                   nil,
		SupportedScopes:                   nil,
		DeviceAuthorization:               op.DeviceAuthorizationConfig{},
		BackChannelLogoutSupported:        false,
		BackChannelLogoutSessionSupported: false,
	}, op.Storage(), func(insecure bool) (op.IssuerFromRequest, error) {
		
	}, op.)
}