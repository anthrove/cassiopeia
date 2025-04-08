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
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/go-playground/validator/v10"
)

func (is IdentityService) CreateEnforcer(ctx context.Context, tenantID string, createEnforcer object.CreateEnforcer) (object.Enforcer, error) {
	err := validate.Struct(createEnforcer)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.Enforcer{}, errors.Join(fmt.Errorf("problem while validating create enforcer data"), util.ConvertValidationError(validateErrs))
		}
	}

	return repository.CreateEnforcer(ctx, is.db, tenantID, createEnforcer)
}

func (is IdentityService) UpdateEnforcer(ctx context.Context, tenantID string, enforcerID string, updateEnforcer object.UpdateEnforcer) error {
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

	return repository.UpdateEnforcer(ctx, is.db, tenantID, enforcerID, updateEnforcer)
}

func (is IdentityService) KillEnforcer(ctx context.Context, tenantID string, enforcerID string) error {
	return repository.KillEnforcer(ctx, is.db, tenantID, enforcerID)
}

func (is IdentityService) FindEnforcer(ctx context.Context, tenantID string, enforcerID string) (object.Enforcer, error) {
	return repository.FindEnforcer(ctx, is.db, tenantID, enforcerID)
}

func (is IdentityService) FindEnforcers(ctx context.Context, tenantID string, pagination object.Pagination) ([]object.Enforcer, error) {
	return repository.FindEnforcers(ctx, is.db, tenantID, pagination)
}

func (is IdentityService) Enforce(ctx context.Context, tenantID string, enforcerID string, request []any) (bool, error) {
	enforcer, err := is.FindEnforcer(ctx, tenantID, enforcerID)

	if err != nil {
		return false, err
	}

	dbModel, err := is.FindModel(ctx, tenantID, enforcer.ModelID)

	if err != nil {
		return false, err
	}

	adapter, err := is.FindAdapter(ctx, tenantID, enforcer.ModelID)
	if err != nil {
		return false, err
	}

	var casAdapter *gormadapter.Adapter

	if adapter.ExternalDB {
		casAdapter, err = gormadapter.NewAdapter(adapter.Driver, "mysql_username:mysql_password@tcp(127.0.0.1:3306)/", adapter.TableName)
	} else {
		casAdapter, err = gormadapter.NewAdapterByDBUseTableName(is.db, "", adapter.TableName)
	}

	casModel, err := model.NewModelFromString(dbModel.Model)

	if err != nil {
		return false, err
	}

	casEnforcer, err := casbin.NewEnforcer(casModel, casAdapter)

	if err != nil {
		return false, err
	}

	success, err := casEnforcer.Enforce(request...)

	if err != nil {
		return false, err
	}

	return success, nil
}
