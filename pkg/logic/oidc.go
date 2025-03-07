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

package logic

import (
	"context"
	"github.com/anthrove/identity/pkg/repository"
	"github.com/go-jose/go-jose/v4"
)

func (is IdentityService) GetJWKs(ctx context.Context) (jose.JSONWebKeySet, error) {
	certs, err := repository.FindAllCertificates(ctx, is.db)

	if err != nil {
		return jose.JSONWebKeySet{}, err
	}

	jwks := jose.JSONWebKeySet{}
	for _, cert := range certs {
		jwk, err := cert.ToJWK()

		if err != nil {
			return jose.JSONWebKeySet{}, err
		}

		jwks.Keys = append(jwks.Keys, jwk)
	}

	return jwks, nil
}
