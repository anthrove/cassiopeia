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

func CreateAdapter(ctx context.Context, db *gorm.DB, tenantId string, createAdapter object.CreateAdapter) (object.Adapter, error) {
	adapter := object.Adapter{
		TenantID:     tenantId,
		Name:         createAdapter.Name,
		TableName:    createAdapter.TableName,
		ExternalDB:   createAdapter.ExternalDB,
		Driver:       createAdapter.Driver,
		Host:         createAdapter.Host,
		Port:         createAdapter.Port,
		Username:     createAdapter.Username,
		Password:     createAdapter.Password,
		DatabaseName: createAdapter.DatabaseName,
	}

	err := db.WithContext(ctx).Model(&object.Adapter{}).Create(&adapter).Error

	return adapter, err
}

func UpdateAdapter(ctx context.Context, db *gorm.DB, tenantID string, adapterId string, updateAdapter object.UpdateAdapter) error {
	adapter := object.Adapter{
		TenantID:     tenantID,
		Name:         updateAdapter.Name,
		TableName:    updateAdapter.TableName,
		ExternalDB:   updateAdapter.ExternalDB,
		Driver:       updateAdapter.Driver,
		Host:         updateAdapter.Host,
		Port:         updateAdapter.Port,
		Username:     updateAdapter.Username,
		Password:     updateAdapter.Password,
		DatabaseName: updateAdapter.DatabaseName,
	}

	err := db.WithContext(ctx).Model(&object.Adapter{}).Where("id = ? AND tenant_id = ?", adapterId, tenantID).Updates(&adapter).Error

	return err
}

func KillAdapter(ctx context.Context, db *gorm.DB, tenantID string, adapterID string) error {
	return db.WithContext(ctx).Delete(&object.Adapter{}, "id = ? AND tenant_id = ?", adapterID, tenantID).Error
}

func FindAdapter(ctx context.Context, db *gorm.DB, tenantID string, adapterID string) (object.Adapter, error) {
	var adapter object.Adapter
	err := db.WithContext(ctx).Take(&adapter, "id = ? AND tenant_id = ?", adapterID, tenantID).Error
	return adapter, err
}

func FindAdapters(ctx context.Context, db *gorm.DB, tenantID string, pagination object.Pagination) ([]object.Adapter, error) {
	var data []object.Adapter
	err := db.WithContext(ctx).Scopes(Pagination(pagination)).Where("tenant_id = ?", tenantID).Find(&data).Error
	return data, err
}
