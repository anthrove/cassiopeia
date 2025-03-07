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

package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/anthrove/identity/pkg/object"
	"github.com/anthrove/identity/pkg/provider/storage"
	"github.com/anthrove/identity/pkg/repository"
	"github.com/anthrove/identity/pkg/util"
	"github.com/go-playground/validator/v10"
	"io"
)

// CreateResource creates a new resource within a specified tenant.
// It validates the input data using the validator package and returns an error if validation fails.
// If validation passes, it calls the repository to create the resource in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the resource belongs.
//   - createResource: object containing the details of the resource to be created.
//
// Returns:
//   - Resource object if creation is successful.
//   - Error if there is any issue during validation or creation.
func (is IdentityService) CreateResource(ctx context.Context, tenantId string, createResource object.CreateResource, file io.Reader) (object.Resource, error) {
	err := validate.Struct(createResource)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.Resource{}, errors.Join(fmt.Errorf("problem while validating create resource data"), util.ConvertValidationError(validateErrs))
		}
	}

	provider, err := is.FindProvider(ctx, tenantId, createResource.ProviderID)
	if err != nil {
		return object.Resource{}, err
	}

	fileProvider, err := storage.GetStorageProvider(provider)
	if err != nil {
		return object.Resource{}, err
	}

	randomPrefix, err := util.RandomString(10)
	if err != nil {
		return object.Resource{}, err
	}

	filenameWithPrefix := fmt.Sprintf("%s_%s", randomPrefix, createResource.FileName)

	resourcePath := fmt.Sprintf("%s/%s", createResource.Tag, filenameWithPrefix)

	hash, err := util.HashFileMD5(file)
	if err != nil {
		return object.Resource{}, err
	}

	resourceObject, err := fileProvider.Put(resourcePath, file)
	if err != nil {
		return object.Resource{}, err
	}

	// TODO: the URL is not the full URL of the file, including the gin path
	resourceURL, err := resourceObject.StorageInterface.GetURL(resourceObject.Path)
	if err != nil {
		return object.Resource{}, err
	}

	bucketName, err := fileProvider.GetBucketName()
	if err != nil {
		return object.Resource{}, err
	}

	if resourceURL == resourceObject.Path {
		resourceURL = fmt.Sprintf("%s/%s/%s", resourceObject.StorageInterface.GetEndpoint(), bucketName, resourceObject.Path)
	}

	return repository.CreateResource(ctx, is.db, tenantId, createResource, resourcePath, resourceURL, hash)
}

// KillResource deletes an existing resource within a specified tenant.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the resource belongs.
//   - resourceID: unique identifier of the resource to be deleted.
//
// Returns:
//   - Error if there is any issue during deletion.
func (is IdentityService) KillResource(ctx context.Context, tenantID string, resourceID string) error {
	return repository.KillResource(ctx, is.db, tenantID, resourceID)
}

// FindResource retrieves a specific resource within a specified tenant.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the resource belongs.
//   - resourceID: unique identifier of the resource to be retrieved.
//
// Returns:
//   - Resource object if retrieval is successful.
//   - Error if there is any issue during retrieval.
func (is IdentityService) FindResource(ctx context.Context, tenantID string, resourceID string) (object.Resource, error) {
	return repository.FindResource(ctx, is.db, tenantID, resourceID)
}

// FindResources retrieves a list of resources within a specified tenant, with pagination support.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the resources belong.
//   - pagination: object containing pagination details (limit and page).
//
// Returns:
//   - Slice of Resource objects if retrieval is successful.
//   - Error if there is any issue during retrieval.
func (is IdentityService) FindResources(ctx context.Context, tenantID string, pagination object.Pagination) ([]object.Resource, error) {
	return repository.FindResources(ctx, is.db, tenantID, pagination)
}

// FindResourceURL retrieves just the URL of resources within a specified tenant.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the resources belong.
//
// Returns:
//   - URl of Resource objects if retrieval is successful.
//   - Error if there is any issue during retrieval.
func (is IdentityService) FindResourceURL(ctx context.Context, tenantID string, resourceID string) (string, error) {
	return repository.FindResourceURL(ctx, is.db, tenantID, resourceID)
}
