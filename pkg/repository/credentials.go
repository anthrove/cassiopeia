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

func CreateCredential(ctx context.Context, db *gorm.DB, tenantId string, createCredentials object.CreateCredential) (object.Credentials, error) {
	credentials := object.Credentials{
		TenantID: tenantId,
		UserID:   createCredentials.UserID,
		Type:     createCredentials.Type,
		Metadata: createCredentials.Metadata,
		Enabled:  createCredentials.Enabled,
	}

	err := db.WithContext(ctx).Model(&object.Credentials{}).Create(&credentials).Error

	return credentials, err
}

func UpdateCredential(ctx context.Context, db *gorm.DB, tenantID string, credentialsId string, updateCredentials object.UpdateCredential) error {
	credentials := object.Credentials{
		ID:       credentialsId,
		TenantID: tenantID,
		Metadata: updateCredentials.Metadata,
		Enabled:  updateCredentials.Enabled,
	}

	err := db.WithContext(ctx).Model(&object.Credentials{}).Where("id = ? AND tenant_id = ?", credentialsId, tenantID).Updates(&credentials).Error

	return err
}

func KillCredential(ctx context.Context, db *gorm.DB, tenantID string, credentialsID string) error {
	return db.WithContext(ctx).Delete(&object.Credentials{}, "id = ? AND tenant_id = ?", credentialsID, tenantID).Error
}

func FindCredential(ctx context.Context, db *gorm.DB, tenantID string, credentialsID string) (object.Credentials, error) {
	var credentials object.Credentials
	err := db.WithContext(ctx).Take(&credentials, "id = ? AND tenant_id = ?", credentialsID, tenantID).Error
	return credentials, err
}

func FindCredentials(ctx context.Context, db *gorm.DB, tenantID string, pagination object.Pagination) ([]object.Credentials, error) {
	var data []object.Credentials
	err := db.WithContext(ctx).Scopes(Pagination(pagination)).Where("tenant_id = ?", tenantID).Find(&data).Error
	return data, err
}

func FindCredentialsByUser(ctx context.Context, db *gorm.DB, tenantID string, userID string) ([]object.Credentials, error) {
	var data []object.Credentials
	err := db.WithContext(ctx).Where("tenant_id = ? AND user_id = ?", tenantID, userID).Find(&data).Error
	return data, err
}
