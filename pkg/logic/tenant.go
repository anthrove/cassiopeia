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
	"github.com/anthrove/identity/pkg/repository"
	"github.com/anthrove/identity/pkg/util"
	"github.com/go-playground/validator/v10"
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
func (is IdentityService) CreateTenant(ctx context.Context, createTenant object.CreateTenant) (object.Tenant, error) {
	err := validate.Struct(createTenant)

	if err != nil {

		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.Tenant{}, errors.Join(fmt.Errorf("problem while validating create tenant data"), util.ConvertValidationError(validateErrs))
		}
	}

	return repository.CreateTenant(ctx, is.db, createTenant)
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
	if len(tenantID) == 0 {
		return errors.New("tenantID is required")
	}

	err := validate.Struct(updateTenant)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating create tenant data"), validateErrs)
		}
	}

	return repository.UpdateTenant(ctx, is.db, tenantID, updateTenant)
}
