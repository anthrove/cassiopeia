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

func CreateTenant(ctx context.Context, db *gorm.DB, createTenant object.CreateTenant) (object.Tenant, error) {
	tenant := object.Tenant{
		DisplayName:  createTenant.DisplayName,
		PasswordType: createTenant.PasswordType,
	}

	err := db.WithContext(ctx).Model(&object.Tenant{}).Create(&tenant).Error

	return tenant, err
}

func UpdateTenant(ctx context.Context, db *gorm.DB, tenantID string, updateTenant object.UpdateTenant) (object.Tenant, error) {
	tenant := object.Tenant{
		ID:           tenantID,
		DisplayName:  updateTenant.DisplayName,
		PasswordType: updateTenant.PasswordType,
	}

	err := db.WithContext(ctx).Model(&object.Tenant{}).Create(&tenant).Error

	return tenant, err
}
