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

// CreateGroup creates a new group within a specified tenant in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - tenantId: unique identifier of the tenant to which the group belongs.
//   - createGroup: object containing the details of the group to be created.
//
// Returns:
//   - Group object if creation is successful.
//   - Error if there is any issue during creation.
func CreateGroup(ctx context.Context, db *gorm.DB, tenantId string, createGroup object.CreateGroup) (object.Group, error) {
	group := object.Group{
		TenantID:      tenantId,
		DisplayName:   createGroup.DisplayName,
		ParentGroupID: createGroup.ParentGroupID,
	}

	err := db.WithContext(ctx).Model(&object.Group{}).Create(&group).Error

	return group, err
}

// UpdateGroup updates an existing group's information within a specified tenant in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - tenantID: unique identifier of the tenant to which the group belongs.
//   - groupId: unique identifier of the group to be updated.
//   - updateGroup: object containing the updated details of the group.
//
// Returns:
//   - Error if there is any issue during updating.
func UpdateGroup(ctx context.Context, db *gorm.DB, tenantID string, groupId string, updateGroup object.UpdateGroup) error {
	group := object.Group{
		TenantID:      tenantID,
		DisplayName:   updateGroup.DisplayName,
		ParentGroupID: updateGroup.ParentGroupID,
	}

	err := db.WithContext(ctx).Model(&object.Group{}).Where("id = ? AND tenant_id = ?", groupId, tenantID).Updates(&group).Error

	return err
}

// KillGroup deletes an existing group within a specified tenant from the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - tenantID: unique identifier of the tenant to which the group belongs.
//   - groupID: unique identifier of the group to be deleted.
//
// Returns:
//   - Error if there is any issue during deletion.
func KillGroup(ctx context.Context, db *gorm.DB, tenantID string, groupID string) error {
	return db.WithContext(ctx).Delete(&object.Group{}, "id = ? AND tenant_id = ?", groupID, tenantID).Error
}

// FindGroup retrieves a specific group within a specified tenant from the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - tenantID: unique identifier of the tenant to which the group belongs.
//   - groupID: unique identifier of the group to be retrieved.
//
// Returns:
//   - Group object if retrieval is successful.
//   - Error if there is any issue during retrieval.
func FindGroup(ctx context.Context, db *gorm.DB, tenantID string, groupID string) (object.Group, error) {
	var group object.Group
	err := db.WithContext(ctx).Take(&group, "id = ? AND tenant_id = ?", groupID, tenantID).Error
	return group, err
}

// FindGroups retrieves a list of groups within a specified tenant from the database, with pagination support.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - tenantID: unique identifier of the tenant to which the groups belong.
//   - pagination: object containing pagination details (limit and page).
//
// Returns:
//   - Slice of Group objects if retrieval is successful.
//   - Error if there is any issue during retrieval.
func FindGroups(ctx context.Context, db *gorm.DB, tenantID string, pagination object.Pagination) ([]object.Group, error) {
	var data []object.Group
	err := db.WithContext(ctx).Scopes(Pagination(pagination)).Where("tenant_id = ?", tenantID).Find(&data).Error
	return data, err
}

func FindGroupsByParentID(ctx context.Context, db *gorm.DB, tenantID string, parentGroupID string) ([]object.Group, error) {
	var data []object.Group
	err := db.WithContext(ctx).Where("tenant_id = ? AND parent_group_id = ?", tenantID, parentGroupID).Find(&data).Error
	return data, err
}

func AppendUserToGroup(ctx context.Context, db *gorm.DB, tenantID string, userID string, groupID string) error {
	return db.WithContext(ctx).Model(object.Group{
		ID:       groupID,
		TenantID: tenantID,
	}).Association("Users").Append(&object.User{
		ID:       userID,
		TenantID: tenantID,
	})
}

func RemoveUserFromGroup(ctx context.Context, db *gorm.DB, tenantID string, userID string, groupID string) error {
	return db.WithContext(ctx).Model(object.Group{
		ID:       groupID,
		TenantID: tenantID,
	}).Association("Users").Delete(&object.User{
		ID:       userID,
		TenantID: tenantID,
	})
}

func FindUsersInGroup(ctx context.Context, db *gorm.DB, tenantID string, groupID string) ([]object.User, error) {
	var users []object.User
	err := db.WithContext(ctx).Model(object.Group{
		ID:       groupID,
		TenantID: tenantID,
	}).Association("Users").Find(&users)

	return users, err
}
