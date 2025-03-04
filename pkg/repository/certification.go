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

func CreateCertificate(ctx context.Context, db *gorm.DB, tenantId string, createCertificate object.Certificate) (object.Certificate, error) {
	err := db.WithContext(ctx).Model(&object.Certificate{}).Create(&createCertificate).Error

	return createCertificate, err
}

func UpdateCertificate(ctx context.Context, db *gorm.DB, tenantID string, certificateId string, updateCertificate object.UpdateCertificate) error {
	certificate := object.Certificate{
		TenantID:    tenantID,
		DisplayName: updateCertificate.DisplayName,
	}

	err := db.WithContext(ctx).Model(&object.Certificate{}).Where("id = ? AND tenant_id = ?", certificateId, tenantID).Updates(&certificate).Error

	return err
}

func KillCertificate(ctx context.Context, db *gorm.DB, tenantID string, certificateID string) error {
	return db.WithContext(ctx).Delete(&object.Certificate{}, "id = ? AND tenant_id = ?", certificateID, tenantID).Error
}

func FindCertificate(ctx context.Context, db *gorm.DB, tenantID string, certificateID string) (object.Certificate, error) {
	var certificate object.Certificate
	err := db.WithContext(ctx).Take(&certificate, "id = ? AND tenant_id = ?", certificateID, tenantID).Error
	return certificate, err
}

func FindCertificates(ctx context.Context, db *gorm.DB, tenantID string, pagination object.Pagination) ([]object.Certificate, error) {
	var data []object.Certificate
	err := db.WithContext(ctx).Scopes(Pagination(pagination)).Where("tenant_id = ?", tenantID).Find(&data).Error
	return data, err
}
