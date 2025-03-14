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
	"github.com/anthrove/identity/pkg/provider/mfa"
	"github.com/anthrove/identity/pkg/repository"
	"github.com/anthrove/identity/pkg/util"
	"github.com/go-playground/validator/v10"
)

// CreateMFA creates a new MFA for a specific User within a specified tenant.
// It validates the input data using the validator package and returns an error if validation fails.
// If validation passes, it calls the repository to create the MFA in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the MFA belongs.
//   - createMFA: object containing the details of the MFA to be created.
//
// Returns:
//   - MFA object if creation is successful.
//   - Error if there is any issue during validation or creation.
func (is IdentityService) CreateMFA(ctx context.Context, tenantID string, userID string, createMFA object.CreateMFA) (object.MFA, error) {
	err := validate.Struct(createMFA)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.MFA{}, errors.Join(fmt.Errorf("problem while validating create MFA data"), util.ConvertValidationError(validateErrs))
		}
	}

	provider, err := is.FindProvider(ctx, tenantID, createMFA.ProviderID)
	if err != nil {
		return object.MFA{}, err
	}

	// I dont see that we needs this for this provider, but i keep it here for now
	// var parameters map[string]string
	// err = json.Unmarshal(provider.Parameter, &parameters)
	// if err != nil {
	// 	return object.MFA{}, err
	// }

	mfaProvider, err := mfa.GetMFAProvider(provider)
	if err != nil {
		return object.MFA{}, err
	}

	mfaData, err := mfaProvider.Create(userID)
	if err != nil {
		return object.MFA{}, err
	}

	createMFA.Secret = mfaData.Secret
	createMFA.URI = mfaData.URI

	createdMFA, err := repository.CreateMFA(ctx, is.db, tenantID, createMFA)
	if err != nil {
		return object.MFA{}, err
	}

	return createdMFA, nil
}

// UpdateMFA updates an existing MFA's information within a specified tenant.
// It validates the input data using the validator package and returns an error if validation fails.
// If validation passes, it calls the repository to update the MFA in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the MFA belongs.
//   - mfaID: unique identifier of the MFA to be updated.
//   - updateMFA: object containing the updated details of the MFA.
//
// Returns:
//   - Error if there is any issue during validation or updating.
func (is IdentityService) UpdateMFA(ctx context.Context, tenantID string, mfaID string, updateMFA object.UpdateMFA) error {
	if len(tenantID) == 0 {
		return errors.New("tenantID is required")
	}

	if len(mfaID) == 0 {
		return errors.New("mfaID is required")
	}

	err := validate.Struct(updateMFA)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating create tenant data"), validateErrs)
		}
	}

	return repository.UpdateMFA(ctx, is.db, tenantID, mfaID, updateMFA)
}

// KillMFA deletes an existing MFA within a specified tenant.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the MFA belongs.
//   - mfaID: unique identifier of the MFA to be deleted.
//
// Returns:
//   - Error if there is any issue during deletion.
func (is IdentityService) KillMFA(ctx context.Context, tenantID string, mfaID string) error {
	if len(tenantID) == 0 {
		return errors.New("tenantID is required")
	}

	if len(mfaID) == 0 {
		return errors.New("mfaID is required")
	}

	return repository.KillMFA(ctx, is.db, tenantID, mfaID)
}

// FindMFA retrieves a specific MFA within a specified tenant.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the MFA belongs.
//   - mfaID: unique identifier of the MFA to be retrieved.
//
// Returns:
//   - MFA object if retrieval is successful.
//   - Error if there is any issue during retrieval.
func (is IdentityService) FindMFA(ctx context.Context, tenantID string, mfaID string) (object.MFA, error) {
	if len(tenantID) == 0 {
		return object.MFA{}, errors.New("tenantID is required")
	}

	if len(mfaID) == 0 {
		return object.MFA{}, errors.New("mfaID is required")
	}

	return repository.FindMFA(ctx, is.db, tenantID, mfaID)
}

// FindMFAs retrieves a list of MFAs within a specified tenant, with pagination support.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - userID: unique identifier of the user to which the MFAs belong
//   - pagination: object containing pagination details (limit and page).
//
// Returns:
//   - Slice of MFA objects if retrieval is successful.
//   - Error if there is any issue during retrieval.
func (is IdentityService) FindMFAs(ctx context.Context, userID string, pagination object.Pagination) ([]object.MFA, error) {
	if len(userID) == 0 {
		return nil, errors.New("userID is required")
	}

	return repository.FindMFAs(ctx, is.db, userID, pagination)
}

// VerifieMFA updates the verification status of an existing MFA within a specified tenant in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - mfaID: unique identifier of the MFA to be updated.
//   - userID: unique identifier of the user to which the MFAs belong.
//   - verified: boolean indicating the new verification status.
//
// Returns:
//   - Error if there is any issue during updating.
func (is IdentityService) VerifieMFA(ctx context.Context, mfaID string, userID string, verified bool) error {
	if len(mfaID) == 0 {
		return errors.New("mfaID is required")
	}

	if len(userID) == 0 {
		return errors.New("tenantID is required")
	}

	return repository.VerifieMFA(ctx, is.db, mfaID, userID, verified)
}

// UpdateMFARecoveryCodesKeys updates the recovery codes of an existing MFA within a specified tenant in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - mfaID: unique identifier of the MFA to be updated.
//   - userID: unique identifier of the user to which the MFAs belong.
//   - recoveryCodes: slice of strings containing the new recovery codes.
//
// Returns:
//   - Error if there is any issue during updating.
func (is IdentityService) UpdateMFARecoveryCodesKeys(ctx context.Context, mfaID string, userID string, recoveryCodes []string) error {
	if len(mfaID) == 0 {
		return errors.New("mfaID is required")
	}

	if len(userID) == 0 {
		return errors.New("tenantID is required")
	}

	if len(recoveryCodes) == 0 {
		return errors.New("recoveryCodes is required")
	}

	return repository.UpdateMFARecoveryCodes(ctx, is.db, mfaID, userID, recoveryCodes)
}
