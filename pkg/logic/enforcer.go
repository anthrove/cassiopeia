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

func (is IdentityService) CreateEnforcer(ctx context.Context, tenantID string, createEnforcer object.CreateEnforcer) (object.Enforcer, error) {
	dbConn, _ := is.getDBConn(ctx)

	err := validate.Struct(createEnforcer)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.Enforcer{}, errors.Join(fmt.Errorf("problem while validating create enforcer data"), util.ConvertValidationError(validateErrs))
		}
	}

	return repository.CreateEnforcer(ctx, dbConn, tenantID, createEnforcer)
}

func (is IdentityService) UpdateEnforcer(ctx context.Context, tenantID string, enforcerID string, updateEnforcer object.UpdateEnforcer) error {
	dbConn, _ := is.getDBConn(ctx)

	if len(tenantID) == 0 {
		return errors.New("tenantID is required")
	}

	err := validate.Struct(updateEnforcer)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating create tenant data"), validateErrs)
		}
	}

	return repository.UpdateEnforcer(ctx, dbConn, tenantID, enforcerID, updateEnforcer)
}

func (is IdentityService) KillEnforcer(ctx context.Context, tenantID string, enforcerID string) error {
	dbConn, _ := is.getDBConn(ctx)

	return repository.KillEnforcer(ctx, dbConn, tenantID, enforcerID)
}

func (is IdentityService) FindEnforcer(ctx context.Context, tenantID string, enforcerID string) (object.Enforcer, error) {
	dbConn, _ := is.getDBConn(ctx)

	return repository.FindEnforcer(ctx, dbConn, tenantID, enforcerID)
}

func (is IdentityService) FindEnforcers(ctx context.Context, tenantID string, pagination object.Pagination) ([]object.Enforcer, error) {
	dbConn, _ := is.getDBConn(ctx)

	return repository.FindEnforcers(ctx, dbConn, tenantID, pagination)
}

func (is IdentityService) Enforce(ctx context.Context, tenantID string, enforcerID string, request []any) (bool, error) {
	enforcer, err := is.FindEnforcer(ctx, tenantID, enforcerID)

	if err != nil {
		return false, err
	}

	casbinEnforcer, err := is.GetCasbinEnforcer(ctx, tenantID, enforcer.ModelID, enforcer.AdapterID)

	if err != nil {
		return false, err
	}

	success, err := casbinEnforcer.Enforce(request...)

	if err != nil {
		return false, err
	}

	return success, nil
}
