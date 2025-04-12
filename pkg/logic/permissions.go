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

func (is IdentityService) CreatePermission(ctx context.Context, tenantID string, createPermission object.CreatePermission) (object.Permission, error) {
	err := validate.Struct(createPermission)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.Permission{}, errors.Join(fmt.Errorf("problem while validating create permission data"), util.ConvertValidationError(validateErrs))
		}
	}

	permission, err := repository.CreatePermission(ctx, is.db, tenantID, createPermission)

	if err != nil {
		return object.Permission{}, err
	}

	err = is.SyncCasbinPermissions(ctx, tenantID, permission.EnforcerID)

	return permission, err
}

func (is IdentityService) UpdatePermission(ctx context.Context, tenantID string, permissionID string, updatePermission object.UpdatePermission) error {
	if len(tenantID) == 0 {
		return errors.New("tenantID is required")
	}

	err := validate.Struct(updatePermission)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating create tenant data"), validateErrs)
		}
	}

	err = repository.UpdatePermission(ctx, is.db, tenantID, permissionID, updatePermission)

	if err != nil {
		return err
	}

	err = is.SyncCasbinPermissions(ctx, tenantID, updatePermission.EnforcerID)

	return err
}

func (is IdentityService) KillPermission(ctx context.Context, tenantID string, permissionID string) error {
	return repository.KillPermission(ctx, is.db, tenantID, permissionID)
}

func (is IdentityService) FindPermission(ctx context.Context, tenantID string, permissionID string) (object.Permission, error) {
	return repository.FindPermission(ctx, is.db, tenantID, permissionID)
}

func (is IdentityService) FindPermissions(ctx context.Context, tenantID string, pagination object.Pagination) ([]object.Permission, error) {
	return repository.FindPermissions(ctx, is.db, tenantID, pagination)
}

func (is IdentityService) FindPermissionsByEnforcer(ctx context.Context, tenantID string, enforcerID string) ([]object.Permission, error) {
	return repository.FindPermissionsByEnforcer(ctx, is.db, tenantID, enforcerID)
}
