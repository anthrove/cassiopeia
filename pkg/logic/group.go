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

// CreateGroup creates a new group within a specified tenant.
// It validates the input data using the validator package and returns an error if validation fails.
// If validation passes, it calls the repository to create the group in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the group belongs.
//   - createGroup: object containing the details of the group to be created.
//
// Returns:
//   - Group object if creation is successful.
//   - Error if there is any issue during validation or creation.
func (is IdentityService) CreateGroup(ctx context.Context, tenantID string, createGroup object.CreateGroup, opt ...string) (object.Group, error) {
	dbConn := is.getDBConn(ctx)

	err := validate.Struct(createGroup)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.Group{}, errors.Join(fmt.Errorf("problem while validating create group data"), util.ConvertValidationError(validateErrs))
		}
	}

	return repository.CreateGroup(ctx, dbConn, tenantID, createGroup, opt...)
}

// UpdateGroup updates an existing group's information within a specified tenant.
// It validates the input data using the validator package and returns an error if validation fails.
// If validation passes, it calls the repository to update the group in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the group belongs.
//   - groupID: unique identifier of the group to be updated.
//   - updateGroup: object containing the updated details of the group.
//
// Returns:
//   - Error if there is any issue during validation or updating.
func (is IdentityService) UpdateGroup(ctx context.Context, tenantID string, groupID string, updateGroup object.UpdateGroup) error {
	dbConn := is.getDBConn(ctx)

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

	return repository.UpdateGroup(ctx, dbConn, tenantID, groupID, updateGroup)
}

// KillGroup deletes an existing group within a specified tenant.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the group belongs.
//   - groupID: unique identifier of the group to be deleted.
//
// Returns:
//   - Error if there is any issue during deletion.
func (is IdentityService) KillGroup(ctx context.Context, tenantID string, groupID string) error {
	dbConn := is.getDBConn(ctx)

	return repository.KillGroup(ctx, dbConn, tenantID, groupID)
}

// FindGroup retrieves a specific group within a specified tenant.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the group belongs.
//   - groupID: unique identifier of the group to be retrieved.
//
// Returns:
//   - Group object if retrieval is successful.
//   - Error if there is any issue during retrieval.
func (is IdentityService) FindGroup(ctx context.Context, tenantID string, groupID string) (object.Group, error) {
	return repository.FindGroup(ctx, is.db, tenantID, groupID)
}

// FindGroups retrieves a list of groups within a specified tenant, with pagination support.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the groups belong.
//   - pagination: object containing pagination details (limit and page).
//
// Returns:
//   - Slice of Group objects if retrieval is successful.
//   - Error if there is any issue during retrieval.
func (is IdentityService) FindGroups(ctx context.Context, tenantID string, pagination object.Pagination) ([]object.Group, error) {
	return repository.FindGroups(ctx, is.db, tenantID, pagination)
}

func (is IdentityService) FindGroupsByParentID(ctx context.Context, tenantID string, parentGroupID string) ([]object.Group, error) {
	return repository.FindGroupsByParentID(ctx, is.db, tenantID, parentGroupID)
}

func (is IdentityService) AppendUserToGroup(ctx context.Context, tenantID string, userID string, groupID string) error {
	dbConn := is.getDBConn(ctx)

	return repository.AppendUserToGroup(ctx, dbConn, tenantID, userID, groupID)
}

func (is IdentityService) RemoveUserFromGroup(ctx context.Context, tenantID string, userID string, groupID string) error {
	dbConn := is.getDBConn(ctx)

	return repository.RemoveUserFromGroup(ctx, dbConn, tenantID, userID, groupID)
}

func (is IdentityService) FindUsersInGroup(ctx context.Context, tenantID string, groupID string) ([]object.User, error) {
	return repository.FindUsersInGroup(ctx, is.db, tenantID, groupID)
}
