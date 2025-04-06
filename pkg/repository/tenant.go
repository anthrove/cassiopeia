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

// CreateTenant creates a new tenant in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - createTenant: object containing the details of the tenant to be created.
//
// Returns:
//   - Tenant object if creation is successful.
//   - Error if there is any issue during creation.
func CreateTenant(ctx context.Context, db *gorm.DB, createTenant object.CreateTenant) (object.Tenant, error) {
	tenant := object.Tenant{
		DisplayName:  createTenant.DisplayName,
		PasswordType: createTenant.PasswordType,
	}

	err := db.WithContext(ctx).Model(&object.Tenant{}).Create(&tenant).Error

	return tenant, err
}

// UpdateTenant updates an existing tenant's information in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - tenantID: unique identifier of the tenant to be updated.
//   - updateTenant: object containing the updated details of the tenant.
//
// Returns:
//   - Error if there is any issue during updating.
func UpdateTenant(ctx context.Context, db *gorm.DB, tenantID string, updateTenant object.UpdateTenant) error {
	tenant := object.Tenant{
		DisplayName:          updateTenant.DisplayName,
		PasswordType:         updateTenant.PasswordType,
		SigningCertificateID: &updateTenant.SigningCertificateID,
	}

	err := db.WithContext(ctx).Model(&object.Tenant{
		ID: tenantID,
	}).Updates(&tenant).Error

	return err
}

// KillTenant deletes an existing tenant from the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - tenantID: unique identifier of the tenant to be deleted.
//
// Returns:
//   - Error if there is any issue during deletion.
func KillTenant(ctx context.Context, db *gorm.DB, tenantID string) error {
	return db.WithContext(ctx).Delete(&object.Tenant{}, "id = ?", tenantID).Error
}

// FindTenant retrieves a specific tenant from the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - tenantID: unique identifier of the tenant to be retrieved.
//
// Returns:
//   - Tenant object if retrieval is successful.
//   - Error if there is any issue during retrieval.
func FindTenant(ctx context.Context, db *gorm.DB, tenantID string) (object.Tenant, error) {
	var tenant object.Tenant
	err := db.WithContext(ctx).Take(&tenant, "id = ?", tenantID).Error
	return tenant, err
}

// FindTenants retrieves a list of tenants from the database, with pagination support.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - pagination: object containing pagination details (limit and page).
//
// Returns:
//   - Slice of Tenant objects if retrieval is successful.
//   - Error if there is any issue during retrieval.
func FindTenants(ctx context.Context, db *gorm.DB, pagination object.Pagination) ([]object.Tenant, error) {
	var data []object.Tenant
	err := db.WithContext(ctx).Scopes(Pagination(pagination)).Find(&data).Error
	return data, err
}
