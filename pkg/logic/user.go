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
	"github.com/anthrove/identity/pkg/provider/auth"
	"github.com/anthrove/identity/pkg/repository"
	"github.com/anthrove/identity/pkg/util"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"math"
	"strconv"
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
func (is IdentityService) CreateUser(ctx context.Context, tenantID string, createUser object.CreateUser, opt ...string) (object.User, error) {
	dbConn, nested := is.getDBConn(ctx)

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

	// If user already exists, no error returned!
	user, err := is.FindUserByUsername(ctx, tenantID, createUser.Username)
	if err == nil {
		return object.User{}, errors.New("username already exists")
	}
	// ==============================================

	// Check if email already exists while verified.
	// We want to make sure nobody is blocking a email while not verifying it!
	users, err := is.FindUsersByEmail(ctx, tenantID, createUser.Email)
	if err != nil {
		return object.User{}, err
	}

	for _, user := range users {
		if user.EmailVerified == true {
			return object.User{}, errors.New("email is already verified")
		}
	}

	userTenant, err := is.FindTenant(ctx, tenantID)
	if err != nil {
		return object.User{}, err
	}

	tenantProviders, err := is.FindProviders(ctx, tenantID, object.Pagination{
		Limit: math.MaxInt,
		Page:  0,
	})

	if err != nil {
		return object.User{}, err
	}

	emailVerificationToken := util.RandomNumber(6)

	var tx *gorm.DB
	if !nested {
		tx = dbConn.Begin()
		ctx = saveDBConn(ctx, tx)
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
	}

	user, err = repository.CreateUser(ctx, dbConn, tenantID, createUser, opt...)

	if err != nil {
		if !nested {
			tx.Rollback()
		}
		return object.User{}, err
	}

	err = is.UpdateUserEmail(ctx, tenantID, user.ID, object.UpdateEmail{
		Email:                  user.Email,
		EmailVerified:          false,
		EmailVerificationToken: strconv.Itoa(emailVerificationToken),
	})

	if err != nil {
		if !nested {
			tx.Rollback()
		}
		return object.User{}, err
	}

	var passProvider object.Provider
	for _, provider := range tenantProviders {
		if provider.ProviderType == "password" {
			passProvider = provider
		}
	}

	if len(passProvider.ID) == 0 {
		if !nested {
			tx.Rollback()
		}
		return object.User{}, errors.New("no password provider configured in tenant")
	}

	passwordProvider, err := auth.GetAuthProvider(passProvider)

	if err != nil {
		return object.User{}, err
	}

	metadata, err := passwordProvider.Configure(ctx, auth.ProviderContext{
		Tenant:     userTenant,
		User:       user,
		Credential: object.Credentials{},
		SendMail:   nil,
	}, map[string]any{
		"password": createUser.Password,
	})

	if err != nil {
		if !nested {
			tx.Rollback()
		}
		return object.User{}, err
	}

	_, err = is.CreateCredential(ctx, tenantID, object.CreateCredential{
		UserID:   user.ID,
		Type:     "password",
		Metadata: metadata,
		Enabled:  true,
	})

	if err != nil {
		if !nested {
			tx.Rollback()
		}
		return object.User{}, err
	}

	if !nested {
		err = tx.Commit().Error
		if err != nil {
			return object.User{}, err
		}
	}
	return user, nil
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
	dbConn, _ := is.getDBConn(ctx)

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

	return repository.UpdateUser(ctx, dbConn, tenantID, userID, updateUser)
}

func (is IdentityService) UpdateUserEmail(ctx context.Context, tenantID string, userID string, updateUserEmail object.UpdateEmail) error {
	dbConn, _ := is.getDBConn(ctx)

	if len(tenantID) == 0 {
		return errors.New("tenantID is required")
	}

	if len(userID) == 0 {
		return errors.New("userID is required")
	}

	err := validate.Struct(updateUserEmail)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating create tenant data"), validateErrs)
		}
	}

	return repository.UpdateUserEmail(ctx, dbConn, tenantID, userID, updateUserEmail)
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
	dbConn, _ := is.getDBConn(ctx)

	if len(tenantID) == 0 {
		return errors.New("tenantID is required")
	}

	if len(userID) == 0 {
		return errors.New("userID is required")
	}

	return repository.KillUser(ctx, dbConn, tenantID, userID)
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
	dbConn, _ := is.getDBConn(ctx)

	if len(tenantID) == 0 {
		return object.User{}, errors.New("tenantID is required")
	}

	if len(userID) == 0 {
		return object.User{}, errors.New("userID is required")
	}

	return repository.FindUser(ctx, dbConn, tenantID, userID)
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
	dbConn, _ := is.getDBConn(ctx)

	if len(tenantID) == 0 {
		return nil, errors.New("tenantID is required")
	}

	return repository.FindUsers(ctx, dbConn, tenantID, pagination)
}

// FindUserByUsername retrieves a user within a specified tenant based on their username.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the user belongs.
//   - userName: the username of the user to be retrieved.
//
// Returns:
//   - User object if retrieval is successful.
//   - Error if there is any issue during retrieval.
func (is IdentityService) FindUserByUsername(ctx context.Context, tenantID string, userName string) (object.User, error) {
	dbConn, _ := is.getDBConn(ctx)

	if len(tenantID) == 0 {
		return object.User{}, errors.New("tenantID is required")
	}

	if len(userName) == 0 {
		return object.User{}, errors.New("userName is required")
	}

	return repository.FindUserByUsername(ctx, dbConn, tenantID, userName)
}

// FindUsersByEmail retrieves a user within a specified tenant based on their email.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the user belongs.
//   - email: the email address of the user to be retrieved.
//
// Returns:
//   - User object if retrieval is successful.
//   - Error if there is any issue during retrieval.
func (is IdentityService) FindUsersByEmail(ctx context.Context, tenantID string, email string) ([]object.User, error) {
	dbConn, _ := is.getDBConn(ctx)

	if len(tenantID) == 0 {
		return nil, errors.New("tenantID is required")
	}

	if len(email) == 0 {
		return nil, errors.New("email is required")
	}

	return repository.FindUsersByEmail(ctx, dbConn, tenantID, email)
}
