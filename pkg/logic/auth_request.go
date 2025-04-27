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

func (is IdentityService) CreateAuthRequest(ctx context.Context, tenantID string, createAuthRequest object.CreateAuthRequest) (object.AuthRequest, error) {
	dbConn, _ := is.getDBConn(ctx)

	err := validate.Struct(createAuthRequest)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.AuthRequest{}, errors.Join(fmt.Errorf("problem while validating create authRequest data"), util.ConvertValidationError(validateErrs))
		}
	}

	return repository.CreateAuthRequest(ctx, dbConn, tenantID, createAuthRequest)
}

func (is IdentityService) UpdateAuthRequest(ctx context.Context, tenantID string, authRequestID string, updateAuthRequest object.UpdateAuthRequest) error {
	dbConn, _ := is.getDBConn(ctx)

	if len(tenantID) == 0 {
		return errors.New("tenantID is required")
	}

	err := validate.Struct(updateAuthRequest)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating create tenant data"), validateErrs)
		}
	}

	return repository.UpdateAuthRequest(ctx, dbConn, tenantID, authRequestID, updateAuthRequest)
}

func (is IdentityService) KillAuthRequest(ctx context.Context, tenantID string, authRequestID string) error {
	dbConn, _ := is.getDBConn(ctx)

	return repository.KillAuthRequest(ctx, dbConn, tenantID, authRequestID)
}

func (is IdentityService) FindAuthRequest(ctx context.Context, tenantID string, authRequestID string) (object.AuthRequest, error) {
	dbConn, _ := is.getDBConn(ctx)

	return repository.FindAuthRequest(ctx, dbConn, tenantID, authRequestID)
}

func (is IdentityService) FindAuthRequests(ctx context.Context, tenantID string, pagination object.Pagination) ([]object.AuthRequest, error) {
	dbConn, _ := is.getDBConn(ctx)

	return repository.FindAuthRequests(ctx, dbConn, tenantID, pagination)
}
