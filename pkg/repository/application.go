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

func CreateApplication(ctx context.Context, db *gorm.DB, tenantId string, createApplication object.CreateApplication) (object.Application, error) {
	application := object.Application{
		TenantID:      tenantId,
		CertificateID: createApplication.CertificateID,
		DisplayName:   createApplication.DisplayName,
		Logo:          createApplication.Logo,
		SignInURL:     createApplication.SignInURL,
		SignUpURL:     createApplication.SignUpURL,
		ForgetURL:     createApplication.ForgetURL,
		TermsURL:      createApplication.TermsURL,
		RedirectURLs:  createApplication.RedirectURLs,
	}

	err := db.WithContext(ctx).Model(&object.Application{}).Create(&application).Error

	return application, err
}

func UpdateApplication(ctx context.Context, db *gorm.DB, tenantID string, applicationID string, updateApplication object.UpdateApplication) error {
	application := object.Application{
		TenantID:      tenantID,
		CertificateID: updateApplication.CertificateID,
		DisplayName:   updateApplication.DisplayName,
		Logo:          updateApplication.Logo,
		SignInURL:     updateApplication.SignInURL,
		SignUpURL:     updateApplication.SignUpURL,
		ForgetURL:     updateApplication.ForgetURL,
		TermsURL:      updateApplication.TermsURL,
		RedirectURLs:  updateApplication.RedirectURLs,
	}

	err := db.WithContext(ctx).Model(&object.Application{}).Where("id = ? AND tenant_id = ?", applicationID, tenantID).Updates(&application).Error

	return err
}

func KillApplication(ctx context.Context, db *gorm.DB, tenantID string, applicationID string) error {
	return db.WithContext(ctx).Delete(&object.Application{}, "id = ? AND tenant_id = ?", applicationID, tenantID).Error
}

func FindApplication(ctx context.Context, db *gorm.DB, tenantID string, applicationID string) (object.Application, error) {
	var group object.Application
	err := db.WithContext(ctx).Take(&group, "id = ? AND tenant_id = ?", applicationID, tenantID).Error
	return group, err
}

func FindApplications(ctx context.Context, db *gorm.DB, tenantID string, pagination object.Pagination) ([]object.Application, error) {
	var data []object.Application
	err := db.WithContext(ctx).Scopes(Pagination(pagination)).Where("tenant_id = ?", tenantID).Find(&data).Error
	return data, err
}
