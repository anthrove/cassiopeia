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

func CreateModel(ctx context.Context, db *gorm.DB, tenantId string, createModel object.CreateModel) (object.Model, error) {
	model := object.Model{
		TenantID:    tenantId,
		Name:        createModel.Name,
		Description: createModel.Description,
		Model:       createModel.Model,
	}

	err := db.WithContext(ctx).Model(&object.Model{}).Create(&model).Error

	return model, err
}

func UpdateModel(ctx context.Context, db *gorm.DB, tenantID string, modelId string, updateModel object.UpdateModel) error {
	model := object.Model{
		TenantID:    tenantID,
		Name:        updateModel.Name,
		Description: updateModel.Description,
		Model:       updateModel.Model,
	}

	err := db.WithContext(ctx).Model(&object.Model{}).Where("id = ? AND tenant_id = ?", modelId, tenantID).Updates(&model).Error

	return err
}

func KillModel(ctx context.Context, db *gorm.DB, tenantID string, modelID string) error {
	return db.WithContext(ctx).Delete(&object.Model{}, "id = ? AND tenant_id = ?", modelID, tenantID).Error
}

func FindModel(ctx context.Context, db *gorm.DB, tenantID string, modelID string) (object.Model, error) {
	var model object.Model
	err := db.WithContext(ctx).Take(&model, "id = ? AND tenant_id = ?", modelID, tenantID).Error
	return model, err
}

func FindModels(ctx context.Context, db *gorm.DB, tenantID string, pagination object.Pagination) ([]object.Model, error) {
	var data []object.Model
	err := db.WithContext(ctx).Scopes(Pagination(pagination)).Where("tenant_id = ?", tenantID).Find(&data).Error
	return data, err
}
