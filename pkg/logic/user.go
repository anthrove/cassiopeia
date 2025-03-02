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
	"github.com/anthrove/identity/pkg/crypto"
	"github.com/anthrove/identity/pkg/object"
	"github.com/anthrove/identity/pkg/repository"
	"github.com/anthrove/identity/pkg/util"
	"github.com/go-playground/validator/v10"
)

// CreateUser creates a new user within a specified tenant.
// It validates the input data using the validator package and returns an error if validation fails.
// If validation passes, it calls the repository to create the user in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the user belongs.
//   - createUser: object containing the details of the user to be created.
//
// Returns:
//   - User object if creation is successful.
//   - Error if there is any issue during validation or creation.
func (is IdentityService) CreateUser(ctx context.Context, tenantID string, createUser object.CreateUser) (object.User, error) {
	if len(tenantID) == 0 {
		return object.User{}, errors.New("tenantID is required")
	}

	err := validate.Struct(createUser)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.User{}, errors.Join(fmt.Errorf("problem while validating create user data"), util.ConvertValidationError(validateErrs))
		}
	}

	// TODO: Verify that the user dose not already exist in the tenant
	// TODO: If user email is not verified: error with "account with email already exists"

	userTenant, err := is.FindTenant(ctx, tenantID)
	if err != nil {
		return object.User{}, err
	}

	passwordSalt, err := util.RandomString(25)
	if err != nil {
		return object.User{}, err
	}

	passwordHasher, err := crypto.GetPasswordHasher(userTenant.PasswordType)
	if err != nil {
		return object.User{}, err
	}

	passwordHash, err := passwordHasher.HashPassword(createUser.Password, passwordSalt)
	if err != nil {
		return object.User{}, err
	}

	user := object.User{
		TenantID:     tenantID,
		Username:     createUser.Username,
		DisplayName:  createUser.DisplayName,
		Email:        createUser.Email,
		PasswordSalt: passwordSalt,
		PasswordType: userTenant.PasswordType,
		PasswordHash: passwordHash,
	}

	return repository.CreateUser(ctx, is.db, tenantID, user)
}

// UpdateUser updates an existing user's information within a specified tenant.
// It validates the input data using the validator package and returns an error if validation fails.
// If validation passes, it calls the repository to update the user in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the user belongs.
//   - userID: unique identifier of the user to be updated.
//   - updateUser: object containing the updated details of the user.
//
// Returns:
//   - Error if there is any issue during validation or updating.
func (is IdentityService) UpdateUser(ctx context.Context, tenantID string, userID string, updateUser object.UpdateUser) error {
	if len(tenantID) == 0 {
		return errors.New("tenantID is required")
	}

	if len(userID) == 0 {
		return errors.New("userID is required")
	}

	err := validate.Struct(updateUser)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating create tenant data"), validateErrs)
		}
	}

	return repository.UpdateUser(ctx, is.db, tenantID, userID, updateUser)
}

// KillUser deletes an existing user within a specified tenant.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the user belongs.
//   - userID: unique identifier of the user to be deleted.
//
// Returns:
//   - Error if there is any issue during deletion.
func (is IdentityService) KillUser(ctx context.Context, tenantID string, userID string) error {
	if len(tenantID) == 0 {
		return errors.New("tenantID is required")
	}

	if len(userID) == 0 {
		return errors.New("userID is required")
	}

	return repository.KillUser(ctx, is.db, tenantID, userID)
}

// FindUser retrieves a specific user within a specified tenant.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the user belongs.
//   - userID: unique identifier of the user to be retrieved.
//
// Returns:
//   - User object if retrieval is successful.
//   - Error if there is any issue during retrieval.
func (is IdentityService) FindUser(ctx context.Context, tenantID string, userID string) (object.User, error) {
	if len(tenantID) == 0 {
		return object.User{}, errors.New("tenantID is required")
	}

	if len(userID) == 0 {
		return object.User{}, errors.New("userID is required")
	}

	return repository.FindUser(ctx, is.db, tenantID, userID)
}

// FindUsers retrieves a list of users within a specified tenant, with pagination support.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the users belong.
//   - pagination: object containing pagination details (limit and page).
//
// Returns:
//   - Slice of User objects if retrieval is successful.
//   - Error if there is any issue during retrieval.
func (is IdentityService) FindUsers(ctx context.Context, tenantID string, pagination object.Pagination) ([]object.User, error) {
	if len(tenantID) == 0 {
		return nil, errors.New("tenantID is required")
	}

	return repository.FindUsers(ctx, is.db, tenantID, pagination)
}
