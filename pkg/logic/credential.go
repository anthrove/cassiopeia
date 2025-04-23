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

func (is IdentityService) CreateCredential(ctx context.Context, tenantID string, createCredential object.CreateCredential) (object.Credentials, error) {
	err := validate.Struct(createCredential)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.Credentials{}, errors.Join(fmt.Errorf("problem while validating create credential data"), util.ConvertValidationError(validateErrs))
		}
	}

	return repository.CreateCredential(ctx, is.db, tenantID, createCredential)
}

func (is IdentityService) UpdateCredential(ctx context.Context, tenantID string, credentialID string, updateCredential object.UpdateCredential) error {
	if len(tenantID) == 0 {
		return errors.New("tenantID is required")
	}

	err := validate.Struct(updateCredential)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating create tenant data"), validateErrs)
		}
	}

	return repository.UpdateCredential(ctx, is.db, tenantID, credentialID, updateCredential)
}

func (is IdentityService) KillCredential(ctx context.Context, tenantID string, credentialID string) error {
	return repository.KillCredential(ctx, is.db, tenantID, credentialID)
}

func (is IdentityService) FindCredential(ctx context.Context, tenantID string, credentialID string) (object.Credentials, error) {
	return repository.FindCredential(ctx, is.db, tenantID, credentialID)
}

func (is IdentityService) FindCredentials(ctx context.Context, tenantID string, pagination object.Pagination) ([]object.Credentials, error) {
	return repository.FindCredentials(ctx, is.db, tenantID, pagination)
}

func (is IdentityService) FindCredentialsByUser(ctx context.Context, tenantID string, userID string) ([]object.Credentials, error) {
	return repository.FindCredentialsByUser(ctx, is.db, tenantID, userID)
}
