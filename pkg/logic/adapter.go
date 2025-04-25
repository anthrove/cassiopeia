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

func (is IdentityService) CreateAdapter(ctx context.Context, tenantID string, createAdapter object.CreateAdapter) (object.Adapter, error) {
	dbConn := is.getDBConn(ctx)

	err := validate.Struct(createAdapter)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.Adapter{}, errors.Join(fmt.Errorf("problem while validating create adapter data"), util.ConvertValidationError(validateErrs))
		}
	}

	return repository.CreateAdapter(ctx, dbConn, tenantID, createAdapter)
}

func (is IdentityService) UpdateAdapter(ctx context.Context, tenantID string, adapterID string, updateAdapter object.UpdateAdapter) error {
	dbConn := is.getDBConn(ctx)

	if len(tenantID) == 0 {
		return errors.New("tenantID is required")
	}

	err := validate.Struct(updateAdapter)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating create tenant data"), validateErrs)
		}
	}

	return repository.UpdateAdapter(ctx, dbConn, tenantID, adapterID, updateAdapter)
}

func (is IdentityService) KillAdapter(ctx context.Context, tenantID string, adapterID string) error {
	dbConn := is.getDBConn(ctx)

	return repository.KillAdapter(ctx, dbConn, tenantID, adapterID)
}

func (is IdentityService) FindAdapter(ctx context.Context, tenantID string, adapterID string) (object.Adapter, error) {
	return repository.FindAdapter(ctx, is.db, tenantID, adapterID)
}

func (is IdentityService) FindAdapters(ctx context.Context, tenantID string, pagination object.Pagination) ([]object.Adapter, error) {
	return repository.FindAdapters(ctx, is.db, tenantID, pagination)
}
