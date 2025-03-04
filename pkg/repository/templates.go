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

// CreateMessageTemplate creates a new template in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - createMessageTemplate: object containing the details of the template to be created.
//
// Returns:
//   - Template object if creation is successful.
//   - Error if there is any issue during creation.
func CreateMessageTemplate(ctx context.Context, db *gorm.DB, tenantId string, createMessageTemplate object.CreateMessageTemplate) (object.MessageTemplate, error) {
	template := object.MessageTemplate{
		TenantID:     tenantId,
		DisplayName:  createMessageTemplate.DisplayName,
		TemplateType: createMessageTemplate.TemplateType,
		Template:     createMessageTemplate.Template,
	}

	err := db.WithContext(ctx).Model(&object.MessageTemplate{}).Create(&template).Error

	return template, err
}

// UpdateMessageTemplate updates an existing template's information in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - templateID: unique identifier of the template to be updated.
//   - updateMessageTemplate: object containing the updated details of the template.
//
// Returns:
//   - Error if there is any issue during updating.
func UpdateMessageTemplate(ctx context.Context, db *gorm.DB, tenantID string, templateID string, updateMessageTemplate object.UpdateMessageTemplate) error {
	template := object.UpdateMessageTemplate{
		DisplayName: updateMessageTemplate.DisplayName,
		Template:    updateMessageTemplate.Template,
	}

	err := db.WithContext(ctx).Model(&object.MessageTemplate{ID: templateID, TenantID: tenantID}).Updates(&template).Error

	return err
}

// KillMessageTemplate deletes an existing template from the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - templateID: unique identifier of the template to be deleted.
//
// Returns:
//   - Error if there is any issue during deletion.
func KillMessageTemplate(ctx context.Context, db *gorm.DB, tenantID string, templateID string) error {
	return db.WithContext(ctx).Delete(&object.MessageTemplate{}, "tenant_id = ? AND id = ?", tenantID, templateID).Error
}

// FindMessageTemplate retrieves a specific template from the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - templateID: unique identifier of the template to be retrieved.
//
// Returns:
//   - Template object if retrieval is successful.
//   - Error if there is any issue during retrieval.
func FindMessageTemplate(ctx context.Context, db *gorm.DB, tenantID string, templateID string) (object.MessageTemplate, error) {
	var template object.MessageTemplate
	err := db.WithContext(ctx).Take(&template, "tenant_id = ? AND id = ?", tenantID, templateID).Error
	return template, err
}

// FindMessageTemplates retrieves a list of templates from the database, with pagination support.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - pagination: object containing pagination details (limit and page).
//
// Returns:
//   - Slice of Template objects if retrieval is successful.
//   - Error if there is any issue during retrieval.
func FindMessageTemplates(ctx context.Context, db *gorm.DB, tenantID string, pagination object.Pagination) ([]object.MessageTemplate, error) {
	var data []object.MessageTemplate
	err := db.WithContext(ctx).Scopes(Pagination(pagination)).Find(&data).Where("tenant_id = ?", tenantID).Error
	return data, err
}
