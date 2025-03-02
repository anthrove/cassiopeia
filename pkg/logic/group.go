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

func (is IdentityService) CreateGroup(ctx context.Context, tenantID string, createGroup object.CreateGroup) (object.Group, error) {
	err := validate.Struct(createGroup)

	if err != nil {

		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.Group{}, errors.Join(fmt.Errorf("problem while validating create group data"), util.ConvertValidationError(validateErrs))
		}
	}

	return repository.CreateGroup(ctx, is.db, tenantID, createGroup)
}

func (is IdentityService) UpdateGroup(ctx context.Context, tenantID string, groupID string, updateGroup object.UpdateGroup) error {
	if len(tenantID) == 0 {
		return errors.New("tenantID is required")
	}

	err := validate.Struct(updateGroup)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating create tenant data"), validateErrs)
		}
	}

	return repository.UpdateGroup(ctx, is.db, tenantID, groupID, updateGroup)
}

func (is IdentityService) KillGroup(ctx context.Context, tenantID string, groupID string) error {
	return repository.KillGroup(ctx, is.db, tenantID, groupID)
}

func (is IdentityService) FindGroup(ctx context.Context, tenantID string, groupID string) (object.Group, error) {
	return repository.FindGroup(ctx, is.db, tenantID, groupID)
}

func (is IdentityService) FindGroups(ctx context.Context, tenantID string, pagination object.Pagination) ([]object.Group, error) {
	return repository.FindGroups(ctx, is.db, tenantID, pagination)
}
