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

package repository

import (
	"context"
	"github.com/anthrove/identity/pkg/object"
	"gorm.io/gorm"
)

// CreateUser creates a new user within a specified tenant in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - tenantId: unique identifier of the tenant to which the user belongs.
//   - createUser: object containing the details of the user to be created.
//
// Returns:
//   - User object if creation is successful.
//   - Error if there is any issue during creation.
func CreateUser(ctx context.Context, db *gorm.DB, tenantId string, user object.User) (object.User, error) {
	err := db.WithContext(ctx).Model(&object.User{}).Create(&user).Error

	return user, err
}

// UpdateUser updates an existing user's information within a specified tenant in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - tenantID: unique identifier of the tenant to which the user belongs.
//   - userId: unique identifier of the user to be updated.
//   - updateUser: object containing the updated details of the user.
//
// Returns:
//   - Error if there is any issue during updating.
func UpdateUser(ctx context.Context, db *gorm.DB, tenantID string, userId string, updateUser object.UpdateUser) error {
	user := object.User{
		ID:          userId,
		TenantID:    tenantID,
		DisplayName: updateUser.DisplayName,
	}

	err := db.WithContext(ctx).Model(&object.User{}).Where("id = ? AND tenant_id = ?", userId, tenantID).Updates(&user).Error

	return err
}

// KillUser deletes an existing user within a specified tenant from the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - tenantID: unique identifier of the tenant to which the user belongs.
//   - userID: unique identifier of the user to be deleted.
//
// Returns:
//   - Error if there is any issue during deletion.
func KillUser(ctx context.Context, db *gorm.DB, tenantID string, userID string) error {
	return db.WithContext(ctx).Delete(&object.User{}, "id = ? AND tenant_id = ?", userID, tenantID).Error
}

// FindUser retrieves a specific user within a specified tenant from the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - tenantID: unique identifier of the tenant to which the user belongs.
//   - userID: unique identifier of the user to be retrieved.
//
// Returns:
//   - User object if retrieval is successful.
//   - Error if there is any issue during retrieval.
func FindUser(ctx context.Context, db *gorm.DB, tenantID string, userID string) (object.User, error) {
	var user object.User
	err := db.WithContext(ctx).Take(&user, "id = ? AND tenant_id = ?", userID, tenantID).Error
	return user, err
}

// FindUsers retrieves a list of users within a specified tenant from the database, with pagination support.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - tenantID: unique identifier of the tenant to which the users belong.
//   - pagination: object containing pagination details (limit and page).
//
// Returns:
//   - Slice of User objects if retrieval is successful.
//   - Error if there is any issue during retrieval.
func FindUsers(ctx context.Context, db *gorm.DB, tenantID string, pagination object.Pagination) ([]object.User, error) {
	var data []object.User
	err := db.WithContext(ctx).Scopes(Pagination(pagination)).Where("tenant_id = ?", tenantID).Find(&data).Error
	return data, err
}
