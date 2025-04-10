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

func CreatePermission(ctx context.Context, db *gorm.DB, tenantId string, createPermission object.CreatePermission) (object.Permission, error) {
	permission := object.Permission{
		TenantID:    tenantId,
		AdapterID:   createPermission.AdapterID,
		Name:        createPermission.Name,
		Description: createPermission.Description,
		Users:       createPermission.Users,
		Groups:      createPermission.Groups,
		V1:          createPermission.V1,
		V2:          createPermission.V2,
		V3:          createPermission.V3,
		V4:          createPermission.V4,
		V5:          createPermission.V5,
	}

	err := db.WithContext(ctx).Model(&object.Permission{}).Create(&permission).Error

	return permission, err
}

func UpdatePermission(ctx context.Context, db *gorm.DB, tenantID string, permissionId string, updatePermission object.UpdatePermission) error {
	permission := object.Permission{
		TenantID:    tenantID,
		AdapterID:   updatePermission.AdapterID,
		Name:        updatePermission.Name,
		Description: updatePermission.Description,
		Users:       updatePermission.Users,
		Groups:      updatePermission.Groups,
		V1:          updatePermission.V1,
		V2:          updatePermission.V2,
		V3:          updatePermission.V3,
		V4:          updatePermission.V4,
		V5:          updatePermission.V5,
	}

	err := db.WithContext(ctx).Model(&object.Permission{}).Where("id = ? AND tenant_id = ?", permissionId, tenantID).Updates(&permission).Error

	return err
}

func KillPermission(ctx context.Context, db *gorm.DB, tenantID string, permissionID string) error {
	return db.WithContext(ctx).Delete(&object.Permission{}, "id = ? AND tenant_id = ?", permissionID, tenantID).Error
}

func FindPermission(ctx context.Context, db *gorm.DB, tenantID string, permissionID string) (object.Permission, error) {
	var permission object.Permission
	err := db.WithContext(ctx).Take(&permission, "id = ? AND tenant_id = ?", permissionID, tenantID).Error
	return permission, err
}

func FindPermissions(ctx context.Context, db *gorm.DB, tenantID string, pagination object.Pagination) ([]object.Permission, error) {
	var data []object.Permission
	err := db.WithContext(ctx).Scopes(Pagination(pagination)).Where("tenant_id = ?", tenantID).Find(&data).Error
	return data, err
}
