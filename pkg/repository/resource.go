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
	"fmt"
	"github.com/anthrove/identity/pkg/object"
	"github.com/anthrove/identity/pkg/provider/storage"
	"gorm.io/gorm"
	"io"
	"path/filepath"
)

// CreateResource creates a new resource within a specified tenant in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - tenantId: unique identifier of the tenant to which the resource belongs.
//   - createResource: object containing the details of the resource to be created.
//
// Returns:
//   - Resource object if creation is successful.
//   - Error if there is any issue during creation.
func CreateResource(ctx context.Context, db *gorm.DB, tenantId string, createResource object.CreateResource, mimeType string, fileSize int64, file io.Reader) (object.Resource, error) {
	provider, err := FindProvider(ctx, db, tenantId, createResource.ProviderID)
	if err != nil {
		return object.Resource{}, err
	}

	fileProvider, err := storage.GetStorageProvider(provider)
	if err != nil {
		return object.Resource{}, err
	}

	filePath := fmt.Sprintf("%s", createResource.Tag)

	resourceObject, err := fileProvider.Put(filePath, file)
	if err != nil {
		return object.Resource{}, err
	}

	resourceURL, err := resourceObject.StorageInterface.GetURL(resourceObject.Path)
	if err != nil {
		return object.Resource{}, err
	}

	resource := object.Resource{
		TenantID:   tenantId,
		ProviderID: createResource.ProviderID,
		Tag:        createResource.Tag,
		Type:       mimeType,
		FilePath:   resourceObject.Path,
		FileSize:   fileSize,
		Format:     filepath.Ext(resourceObject.Path),
		Url:        resourceURL,
	}

	err = db.WithContext(ctx).Model(&object.Resource{}).Create(&resource).Error

	return resource, err
}

// KillResource deletes an existing resource within a specified tenant from the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - tenantID: unique identifier of the tenant to which the resource belongs.
//   - resourceID: unique identifier of the resource to be deleted.
//
// Returns:
//   - Error if there is any issue during deletion.
func KillResource(ctx context.Context, db *gorm.DB, tenantID string, resourceID string) error {
	return db.WithContext(ctx).Delete(&object.Resource{}, "id = ? AND tenant_id = ?", resourceID, tenantID).Error
}

// FindResource retrieves a specific resource within a specified tenant from the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - tenantID: unique identifier of the tenant to which the resource belongs.
//   - resourceID: unique identifier of the resource to be retrieved.
//
// Returns:
//   - Resourceobject if retrieval is successful.
//   - Error if there is any issue during retrieval.
func FindResource(ctx context.Context, db *gorm.DB, tenantID string, resourceID string) (object.Resource, error) {
	var resource object.Resource
	err := db.WithContext(ctx).Take(&resource, "id = ? AND tenant_id = ?", resourceID, tenantID).Error
	return resource, err
}

// FindResources retrieves a list of resources within a specified tenant from the database, with pagination support.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - tenantID: unique identifier of the tenant to which the resources belong.
//   - pagination: object containing pagination details (limit and page).
//
// Returns:
//   - Slice of Resourceobjects if retrieval is successful.
//   - Error if there is any issue during retrieval.
func FindResources(ctx context.Context, db *gorm.DB, tenantID string, pagination object.Pagination) ([]object.Resource, error) {
	var data []object.Resource
	err := db.WithContext(ctx).Scopes(Pagination(pagination)).Where("tenant_id = ?", tenantID).Find(&data).Error
	return data, err
}

// FindResourceURL retrieves the URL of a specific resource within a specified tenant from the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - db: a gorm.DB instance representing the database connection.
//   - tenantID: unique identifier of the tenant to which the resource belongs.
//   - resourceID: unique identifier of the resource to be retrieved.
//
// Returns:
//   - Resourceobject if retrieval is successful.
//   - Error if there is any issue during retrieval.
func FindResourceURL(ctx context.Context, db *gorm.DB, tenantID string, resourceID string) (string, error) {
	var resource object.Resource
	err := db.WithContext(ctx).Take(&resource, "id = ? AND tenant_id = ?", resourceID, tenantID).Error
	return resource.Url, err
}
