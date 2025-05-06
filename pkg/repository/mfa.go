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

// CreateMFA creates a new MFA for a specific User within a specified tenant in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - createMFA: object containing the details of the MFA to be created.
//
// Returns:
//   - MFA object if creation is successful.
//   - Error if there is any issue during creation.
func CreateMFA(ctx context.Context, db *gorm.DB, tenantID string, userID string, createMFA object.CreateMFA) (object.MFA, error) {
	mfa := object.MFA{
		UserID:        userID,
		DisplayName:   createMFA.DisplayName,
		Type:          createMFA.Type,
		Priority:      createMFA.Priority,
		Verified:      false,
		RecoveryCodes: createMFA.RecoveryCodes,
		Properties:    createMFA.Properties,
	}

	err := db.WithContext(ctx).Model(&object.MFA{}).Create(&mfa).Error

	return mfa, err
}

// UpdateMFA updates an existing MFA's information within a specified tenant in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - userID: unique identifier of the user to which the MFAs belong
//   - mfaID: unique identifier of the MFA to be updated.
//   - updateMFA: object containing the updated details of the MFA.
//
// Returns:
//   - Error if there is any issue during updating.
func UpdateMFA(ctx context.Context, db *gorm.DB, tenantID string, userID string, mfaID string, updateMFA object.UpdateMFA) error {
	mfa := object.MFA{
		ID:          mfaID,
		DisplayName: updateMFA.DisplayName,
		Priority:    updateMFA.Priority,
	}

	err := db.WithContext(ctx).Model(&object.MFA{}).Where("id = ? AND user_id = ?", mfaID, userID).Updates(&mfa).Error

	return err
}

// KillMFA deletes an existing MFA within a specified tenant from the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - userID: unique identifier of the user to which the MFAs belong
//   - mfaID: unique identifier of the MFA to be deleted.
//
// Returns:
//   - Error if there is any issue during deletion.
func KillMFA(ctx context.Context, db *gorm.DB, tenantID string, userID string, mfaID string) error {
	return db.WithContext(ctx).Delete(&object.MFA{}, "id = ? AND user_id = ?", mfaID, userID).Error
}

// FindMFA retrieves a specific MFA within a specified tenant from the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - userID: unique identifier of the user to which the MFAs belong
//   - mfaID: unique identifier of the MFA to be retrieved.
//
// Returns:
//   - MFA object if retrieval is successful.
//   - Error if there is any issue during retrieval.
func FindMFA(ctx context.Context, db *gorm.DB, tenantID string, userID string, mfaID string) (object.MFA, error) {
	var mfa object.MFA
	err := db.WithContext(ctx).Take(&mfa, "id = ? AND user_id = ?", mfaID, userID).Error
	return mfa, err
}

// FindMFAs retrieves a list of MFAs within a specified tenant from the database, with pagination support.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - userID: unique identifier of the user to which the MFAs belong
//   - pagination: object containing pagination details (limit and page).
//
// Returns:
//   - Slice of MFA objects if retrieval is successful.
//   - Error if there is any issue during retrieval.
func FindMFAs(ctx context.Context, db *gorm.DB, tenantID string, userID string, pagination object.Pagination) ([]object.MFA, error) {
	var data []object.MFA
	err := db.WithContext(ctx).Scopes(Pagination(pagination)).Where("user_id = ?", userID).Find(&data).Error
	return data, err
}

// VerifieMFA updates the verification status of an existing MFA within a specified tenant in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - mfaID: unique identifier of the MFA to be updated.
//   - userID: unique identifier of the user to which the MFAs belong.
//   - verified: boolean indicating the new verification status.
//
// Returns:
//   - Error if there is any issue during updating.
func VerifieMFA(ctx context.Context, db *gorm.DB, tenantID string, userID string, mfaID string, verified bool) error {
	mfa := object.MFA{
		ID:       mfaID,
		Verified: verified,
	}

	err := db.WithContext(ctx).Model(&object.MFA{}).Where("id = ? AND user_id = ?", mfaID, userID).Updates(&mfa).Error

	return err
}

// UpdateMFARecoveryCodes updates the recovery codes of an existing MFA within a specified tenant in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - mfaID: unique identifier of the MFA to be updated.
//   - userID: unique identifier of the user to which the MFAs belong.
//   - recoveryCodes: slice of strings containing the new recovery codes.
//
// Returns:
//   - Error if there is any issue during updating.
func UpdateMFARecoveryCodes(ctx context.Context, db *gorm.DB, tenantID string, userID string, mfaID string, recoveryCodes []string) error {
	mfa := object.MFA{
		ID:            mfaID,
		RecoveryCodes: recoveryCodes,
	}

	err := db.WithContext(ctx).Model(&object.MFA{}).Where("id = ? AND user_id = ?", mfaID, userID).Updates(&mfa).Error

	return err
}
