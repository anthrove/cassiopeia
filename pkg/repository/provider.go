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

func CreateProvider(ctx context.Context, db *gorm.DB, tenantId string, createProvider object.CreateProvider) (object.Provider, error) {
	provider := object.Provider{
		TenantID:     tenantId,
		DisplayName:  createProvider.DisplayName,
		Category:     createProvider.Category,
		ProviderType: createProvider.ProviderType,
		Parameter:    createProvider.Parameter,
	}

	err := db.WithContext(ctx).Model(&object.Provider{}).Create(&provider).Error

	return provider, err
}

func UpdateProvider(ctx context.Context, db *gorm.DB, tenantID string, providerID string, updateProvider object.UpdateProvider) error {
	provider := object.Provider{
		TenantID:    tenantID,
		DisplayName: updateProvider.DisplayName,
		Parameter:   updateProvider.Parameter,
	}

	err := db.WithContext(ctx).Model(&object.Provider{}).Where("id = ? AND tenant_id = ?", providerID, tenantID).Updates(&provider).Error

	return err
}

func KillProvider(ctx context.Context, db *gorm.DB, tenantID string, providerID string) error {
	return db.WithContext(ctx).Delete(&object.Provider{}, "id = ? AND tenant_id = ?", providerID, tenantID).Error
}

func FindProvider(ctx context.Context, db *gorm.DB, tenantID string, groupID string) (object.Provider, error) {
	var group object.Provider
	err := db.WithContext(ctx).Take(&group, "id = ? AND tenant_id = ?", groupID, tenantID).Error
	return group, err
}

func FindProviders(ctx context.Context, db *gorm.DB, tenantID string, pagination object.Pagination) ([]object.Provider, error) {
	var data []object.Provider
	err := db.WithContext(ctx).Scopes(Pagination(pagination)).Where("tenant_id = ?", tenantID).Find(&data).Error
	return data, err
}
