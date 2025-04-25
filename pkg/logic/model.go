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

func (is IdentityService) CreateModel(ctx context.Context, tenantID string, createModel object.CreateModel) (object.Model, error) {
	dbConn := is.getDBConn(ctx)

	err := validate.Struct(createModel)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.Model{}, errors.Join(fmt.Errorf("problem while validating create model data"), util.ConvertValidationError(validateErrs))
		}
	}

	return repository.CreateModel(ctx, dbConn, tenantID, createModel)
}

func (is IdentityService) UpdateModel(ctx context.Context, tenantID string, modelID string, updateModel object.UpdateModel) error {
	dbConn := is.getDBConn(ctx)

	if len(tenantID) == 0 {
		return errors.New("tenantID is required")
	}

	err := validate.Struct(updateModel)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating create tenant data"), validateErrs)
		}
	}

	return repository.UpdateModel(ctx, dbConn, tenantID, modelID, updateModel)
}

func (is IdentityService) KillModel(ctx context.Context, tenantID string, modelID string) error {
	dbConn := is.getDBConn(ctx)

	return repository.KillModel(ctx, dbConn, tenantID, modelID)
}

func (is IdentityService) FindModel(ctx context.Context, tenantID string, modelID string) (object.Model, error) {
	return repository.FindModel(ctx, is.db, tenantID, modelID)
}

func (is IdentityService) FindModels(ctx context.Context, tenantID string, pagination object.Pagination) ([]object.Model, error) {
	return repository.FindModels(ctx, is.db, tenantID, pagination)
}
