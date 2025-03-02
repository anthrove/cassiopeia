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

func CreateGroup(ctx context.Context, db *gorm.DB, tenantId string, createGroup object.CreateGroup) (object.Group, error) {
	group := object.Group{
		TenantID:      tenantId,
		DisplayName:   createGroup.DisplayName,
		ParentGroupID: createGroup.ParentGroupID,
	}

	err := db.WithContext(ctx).Model(&object.Group{}).Create(&group).Error

	return group, err
}

func UpdateGroup(ctx context.Context, db *gorm.DB, tenantID string, groupId string, updateGroup object.UpdateGroup) error {
	group := object.Group{
		TenantID:      tenantID,
		DisplayName:   updateGroup.DisplayName,
		ParentGroupID: updateGroup.ParentGroupID,
	}

	err := db.WithContext(ctx).Model(&object.Group{}).Where("id = ? AND tenant_id = ?", groupId, tenantID).Updates(&group).Error

	return err
}

func KillGroup(ctx context.Context, db *gorm.DB, tenantID string, groupID string) error {
	return db.WithContext(ctx).Delete(&object.Group{}, "id = ? AND tenant_id = ?", groupID, tenantID).Error
}

func FindGroup(ctx context.Context, db *gorm.DB, tenantID string, groupID string) (object.Group, error) {
	var group object.Group
	err := db.WithContext(ctx).Take(&group, "id = ? AND tenant_id = ?", groupID, tenantID).Error
	return group, err
}

func FindGroups(ctx context.Context, db *gorm.DB, tenantID string, pagination object.Pagination) ([]object.Group, error) {
	var data []object.Group
	err := db.WithContext(ctx).Scopes(Pagination(pagination)).Where("tenant_id = ?", tenantID).Find(&data).Error
	return data, err
}
