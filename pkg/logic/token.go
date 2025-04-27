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
	"errors"
	"fmt"
	"github.com/anthrove/identity/pkg/object"
	"github.com/anthrove/identity/pkg/repository"
	"github.com/anthrove/identity/pkg/util"
	"github.com/go-playground/validator/v10"
)

func (is IdentityService) CreateToken(ctx context.Context, tenantID string, createToken object.CreateToken) (object.Token, error) {
	dbConn, _ := is.getDBConn(ctx)

	err := validate.Struct(createToken)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.Token{}, errors.Join(fmt.Errorf("problem while validating create token data"), util.ConvertValidationError(validateErrs))
		}
	}

	return repository.CreateToken(ctx, dbConn, tenantID, createToken)
}

func (is IdentityService) KillToken(ctx context.Context, tenantID string, tokenID string) error {
	dbConn, _ := is.getDBConn(ctx)

	return repository.KillToken(ctx, dbConn, tenantID, tokenID)
}

func (is IdentityService) KillTokens(ctx context.Context, tenantID string, tokenIDs []string) error {
	dbConn, _ := is.getDBConn(ctx)

	return repository.KillTokens(ctx, dbConn, tenantID, tokenIDs)
}

func (is IdentityService) FindToken(ctx context.Context, tenantID string, tokenID string) (object.Token, error) {
	dbConn, _ := is.getDBConn(ctx)

	return repository.FindToken(ctx, dbConn, tenantID, tokenID)
}

func (is IdentityService) FindTokenByRefresh(ctx context.Context, tenantID string, refreshToken string) (object.Token, error) {
	dbConn, _ := is.getDBConn(ctx)

	return repository.FindTokenByRefresh(ctx, dbConn, tenantID, refreshToken)
}

func (is IdentityService) FindTokens(ctx context.Context, tenantID string, pagination object.Pagination) ([]object.Token, error) {
	dbConn, _ := is.getDBConn(ctx)

	return repository.FindTokens(ctx, dbConn, tenantID, pagination)
}

func (is IdentityService) FindUserTokens(ctx context.Context, tenantID string, applicationID string, userID string) ([]object.Token, error) {
	dbConn, _ := is.getDBConn(ctx)

	return repository.FindUserTokens(ctx, dbConn, tenantID, applicationID, userID)
}
