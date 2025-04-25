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
	"github.com/anthrove/identity/pkg/crypto"
	"github.com/anthrove/identity/pkg/object"
	"github.com/anthrove/identity/pkg/repository"
	"github.com/anthrove/identity/pkg/util"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-playground/validator/v10"
	"time"
)

// CreateTenant creates a new tenant in the system.
// It validates the input data using the validator package and returns an error if validation fails.
// If validation passes, it calls the repository to create the tenant in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - createTenant: object containing the details of the tenant to be created.
//
// Returns:
//   - Tenant object if creation is successful.
//   - Error if there is any issue during validation or creation.
func (is IdentityService) CreateTenant(ctx context.Context, createTenant object.CreateTenant, opt ...string) (object.Tenant, error) {
	dbConn := is.getDBConn(ctx)

	err := validate.Struct(createTenant)

	if err != nil {

		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.Tenant{}, errors.Join(fmt.Errorf("problem while validating create tenant data"), util.ConvertValidationError(validateErrs))
		}
	}

	_, err = crypto.GetPasswordHasher(createTenant.PasswordType)
	if err != nil {
		return object.Tenant{}, errors.New("password type does not match any known types")
	}

	tenant, err := repository.CreateTenant(ctx, dbConn, createTenant, opt...)

	if err != nil {
		return object.Tenant{}, err
	}

	certificate, err := is.CreateCertificate(ctx, tenant.ID, object.CreateCertificate{
		DisplayName: fmt.Sprintf("Signing Key for %s", tenant.DisplayName),
		Algorithm:   string(jose.RS256),
		BitSize:     2048,
		ExpiredAt:   time.Now().Add(time.Hour * 24 * 365 * 1),
	})

	if err != nil {
		return object.Tenant{}, err
	}

	err = is.UpdateTenant(ctx, tenant.ID, object.UpdateTenant{
		DisplayName:          tenant.DisplayName,
		PasswordType:         tenant.PasswordType,
		SigningCertificateID: certificate.ID(),
	})

	if err != nil {
		return object.Tenant{}, err
	}

	tenant.SigningCertificateID = &certificate.IDs
	return tenant, nil
}

// UpdateTenant updates an existing tenant's information in the system.
// It validates the input data using the validator package and returns an error if validation fails.
// If validation passes, it calls the repository to update the tenant in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to be updated.
//   - updateTenant: object containing the updated details of the tenant.
//
// Returns:
//   - Error if there is any issue during validation or updating.
func (is IdentityService) UpdateTenant(ctx context.Context, tenantID string, updateTenant object.UpdateTenant) error {
	dbConn := is.getDBConn(ctx)

	if len(tenantID) == 0 {
		return errors.New("tenantID is required")
	}

	err := validate.Struct(updateTenant)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating update tenant data"), validateErrs)
		}
	}

	_, err = crypto.GetPasswordHasher(updateTenant.PasswordType)
	if err != nil {
		return errors.New("password type does not match any known types")
	}

	return repository.UpdateTenant(ctx, dbConn, tenantID, updateTenant)
}

// KillTenant deletes an existing tenant from the system.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to be deleted.
//
// Returns:
//   - Error if there is any issue during deletion.
func (is IdentityService) KillTenant(ctx context.Context, tenantID string) error {
	return repository.KillTenant(ctx, is.db, tenantID)
}

// FindTenant retrieves a specific tenant from the system.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to be retrieved.
//
// Returns:
//   - Tenant object if retrieval is successful.
//   - Error if there is any issue during retrieval.
func (is IdentityService) FindTenant(ctx context.Context, tenantID string) (object.Tenant, error) {
	return repository.FindTenant(ctx, is.db, tenantID)
}

// FindTenants retrieves a list of tenants from the system, with pagination support.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - pagination: object containing pagination details (limit and page).
//
// Returns:
//   - Slice of Tenant objects if retrieval is successful.
//   - Error if there is any issue during retrieval.
func (is IdentityService) FindTenants(ctx context.Context, pagination object.Pagination) ([]object.Tenant, error) {
	return repository.FindTenants(ctx, is.db, pagination)
}
