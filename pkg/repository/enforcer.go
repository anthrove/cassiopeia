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

func CreateEnforcer(ctx context.Context, db *gorm.DB, tenantId string, createEnforcer object.CreateEnforcer) (object.Enforcer, error) {
	enforcer := object.Enforcer{
		TenantID:    tenantId,
		Name:        createEnforcer.Name,
		Description: createEnforcer.Description,
		ModelID:     createEnforcer.ModelID,
		AdapterID:   createEnforcer.AdapterID,
	}

	err := db.WithContext(ctx).Model(&object.Enforcer{}).Create(&enforcer).Error

	return enforcer, err
}

func UpdateEnforcer(ctx context.Context, db *gorm.DB, tenantID string, enforcerId string, updateEnforcer object.UpdateEnforcer) error {
	enforcer := object.Enforcer{
		TenantID:    tenantID,
		Name:        updateEnforcer.Name,
		Description: updateEnforcer.Description,
		ModelID:     updateEnforcer.ModelID,
		AdapterID:   updateEnforcer.AdapterID,
	}

	err := db.WithContext(ctx).Model(&object.Enforcer{}).Where("id = ? AND tenant_id = ?", enforcerId, tenantID).Updates(&enforcer).Error

	return err
}

func KillEnforcer(ctx context.Context, db *gorm.DB, tenantID string, enforcerID string) error {
	return db.WithContext(ctx).Delete(&object.Enforcer{}, "id = ? AND tenant_id = ?", enforcerID, tenantID).Error
}

func FindEnforcer(ctx context.Context, db *gorm.DB, tenantID string, enforcerID string) (object.Enforcer, error) {
	var enforcer object.Enforcer
	err := db.WithContext(ctx).Take(&enforcer, "id = ? AND tenant_id = ?", enforcerID, tenantID).Error
	return enforcer, err
}

func FindEnforcers(ctx context.Context, db *gorm.DB, tenantID string, pagination object.Pagination) ([]object.Enforcer, error) {
	var data []object.Enforcer
	err := db.WithContext(ctx).Scopes(Pagination(pagination)).Where("tenant_id = ?", tenantID).Find(&data).Error
	return data, err
}
