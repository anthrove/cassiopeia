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
	"github.com/anthrove/identity/pkg/provider"
	"github.com/anthrove/identity/pkg/provider/auth"
	"github.com/anthrove/identity/pkg/provider/email"
	"github.com/anthrove/identity/pkg/provider/mfa"
	"github.com/anthrove/identity/pkg/provider/storage"
	"github.com/anthrove/identity/pkg/repository"
	"github.com/anthrove/identity/pkg/util"
	"github.com/go-playground/validator/v10"
)

func (is IdentityService) CreateProvider(ctx context.Context, tenantID string, createProvider object.CreateProvider, opt ...string) (object.Provider, error) {
	dbConn, _ := is.getDBConn(ctx)

	err := validate.Struct(createProvider)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.Provider{}, errors.Join(fmt.Errorf("problem while validating create provider data"), util.ConvertValidationError(validateErrs))
		}
	}

	err = validateProvider(object.Provider{
		TenantID:     tenantID,
		Category:     createProvider.Category,
		ProviderType: createProvider.ProviderType,
		Parameter:    createProvider.Parameter,
	})

	if err != nil {
		return object.Provider{}, err
	}

	return repository.CreateProvider(ctx, dbConn, tenantID, createProvider, opt...)
}

func (is IdentityService) UpdateProvider(ctx context.Context, tenantID string, providerID string, updateProvider object.UpdateProvider) error {
	dbConn, _ := is.getDBConn(ctx)

	if len(tenantID) == 0 {
		return errors.New("tenantID is required")
	}

	err := validate.Struct(updateProvider)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating create tenant data"), validateErrs)
		}
	}

	oldProviderObj, err := repository.FindProvider(ctx, dbConn, tenantID, providerID)

	if err != nil {
		return err
	}

	err = validateProvider(object.Provider{
		TenantID:     tenantID,
		Category:     oldProviderObj.Category,
		ProviderType: oldProviderObj.ProviderType,
		Parameter:    updateProvider.Parameter,
	})

	if err != nil {
		return err
	}

	return repository.UpdateProvider(ctx, is.db, tenantID, providerID, updateProvider)
}

func (is IdentityService) KillProvider(ctx context.Context, tenantID string, providerID string) error {
	dbConn, _ := is.getDBConn(ctx)

	return repository.KillProvider(ctx, dbConn, tenantID, providerID)
}

func (is IdentityService) FindProvider(ctx context.Context, tenantID string, providerID string) (object.Provider, error) {
	dbConn, _ := is.getDBConn(ctx)

	return repository.FindProvider(ctx, dbConn, tenantID, providerID)
}

func (is IdentityService) FindProviders(ctx context.Context, tenantID string, pagination object.Pagination) ([]object.Provider, error) {
	dbConn, _ := is.getDBConn(ctx)

	return repository.FindProviders(ctx, dbConn, tenantID, pagination)
}

func (is IdentityService) FindProviderCategories(ctx context.Context, tenantID string) ([]string, error) {
	return []string{"email", "storage", "auth", "mfa"}, nil
}

func (is IdentityService) FindProviderTypes(ctx context.Context, tenantID string, category string) []string {
	return provider.Types(category)
}

func (is IdentityService) FindProviderConfiguration(ctx context.Context, tenantID string, category string, providerType string) []object.ProviderConfigurationField {
	return provider.ConfigurationFields(category, providerType)
}

func validateProvider(providerObj object.Provider) error {
	var provider interface {
		ValidateConfigurationFields() error
	}
	var err error

	switch providerObj.Category {
	case "email":
		provider, err = email.GetEMailProvider(providerObj)
	case "storage":
		provider, err = storage.GetStorageProvider(providerObj)
	case "mfa":
		provider, err = mfa.GetMFAProvider(providerObj)
	case "auth":
		provider, err = auth.GetAuthProvider(providerObj)
	default:
		return errors.New("invalid provider category")
	}

	if err != nil {
		return errors.Join(fmt.Errorf("problem validate provider parameter"), err)
	}

	err = provider.ValidateConfigurationFields()
	if err != nil {
		return errors.Join(fmt.Errorf("problem validate provider parameter"), err)
	}

	return nil
}
