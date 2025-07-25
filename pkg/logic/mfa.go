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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anthrove/identity/pkg/object"
	"github.com/anthrove/identity/pkg/provider/mfa"
	"github.com/anthrove/identity/pkg/repository"
	"github.com/anthrove/identity/pkg/util"
	"github.com/go-playground/validator/v10"
	"slices"
)

// CreateMFA creates a new MFA for a specific User within a specified tenant.
// It validates the input data using the validator package and returns an error if validation fails.
// If validation passes, it calls the repository to create the MFA in the database.
//
// Returns:
//   - MFA object if creation is successful.
//   - Error if there is any issue during validation or creation.
func (is IdentityService) CreateMFA(ctx context.Context, tenantID string, userID string, createMFA object.CreateMFA) (object.MFACreationResponse, error) {
	err := validate.Struct(createMFA)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.MFACreationResponse{}, errors.Join(fmt.Errorf("problem while validating create MFA data"), util.ConvertValidationError(validateErrs))
		}
	}

	provider, err := is.FindProvider(ctx, tenantID, createMFA.ProviderID)
	if err != nil {
		return object.MFACreationResponse{}, err
	}

	// I dont see that we needs this for this provider, but i keep it here for now
	// var parameters map[string]string
	// err = json.Unmarshal(provider.Parameter, &parameters)
	// if err != nil {
	// 	return object.MFA{}, err
	// }

	mfaProvider, err := mfa.GetMFAProvider(provider)
	if err != nil {
		return object.MFACreationResponse{}, err
	}

	mfaData, err := mfaProvider.GenerateUserConfig(userID)
	if err != nil {
		return object.MFACreationResponse{}, err
	}

	var recoveryCodes []string

	// This could be configurable, but i don't see the reason why. So I left this note here for a future dev to maybe implement.
	for range 6 {
		phrase, err := util.RandomPassPhrase(3, "-")
		if err != nil {
			return object.MFACreationResponse{}, err
		}
		recoveryCodes = append(recoveryCodes, phrase)

	}

	createMFA.RecoveryCodes = recoveryCodes
	createMFA.Properties = mfaData.Properties

	createdMFA, err := repository.CreateMFA(ctx, is.db, tenantID, userID, createMFA)
	if err != nil {
		return object.MFACreationResponse{}, err
	}

	return object.MFACreationResponse{
		ID:          createdMFA.ID,
		UserID:      createdMFA.UserID,
		ProviderID:  createdMFA.ProviderID,
		CreatedAt:   createdMFA.CreatedAt,
		UpdatedAt:   createdMFA.UpdatedAt,
		DisplayName: createdMFA.DisplayName,
		Type:        createdMFA.Type,
		Priority:    createdMFA.Priority,
		Verified:    createdMFA.Verified,
		Properties:  createdMFA.Properties,
	}, nil
}

// UpdateMFA updates an existing MFA's information within a specified tenant.
// It validates the input data using the validator package and returns an error if validation fails.
// If validation passes, it calls the repository to update the MFA in the database.
//
// Returns:
//   - Error if there is any issue during validation or updating.
func (is IdentityService) UpdateMFA(ctx context.Context, tenantID string, userID string, mfaID string, updateMFA object.UpdateMFA) error {
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

	return repository.UpdateMFA(ctx, is.db, tenantID, userID, mfaID, updateMFA)
}

// KillMFA deletes an existing MFA within a specified tenant
//
// Returns:
//   - Error if there is any issue during deletion.
func (is IdentityService) KillMFA(ctx context.Context, tenantID string, userID string, mfaID string) error {
	if len(tenantID) == 0 {
		return errors.New("tenantID is required")
	}

	if len(mfaID) == 0 {
		return errors.New("mfaID is required")
	}

	return repository.KillMFA(ctx, is.db, tenantID, userID, mfaID)
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
func (is IdentityService) FindMFA(ctx context.Context, tenantID string, userID string, mfaID string) (object.MFA, error) {
	if len(userID) == 0 {
		return object.MFA{}, errors.New("userID is required")
	}

	if len(mfaID) == 0 {
		return object.MFA{}, errors.New("mfaID is required")
	}

	return repository.FindMFA(ctx, is.db, tenantID, userID, mfaID)
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
func (is IdentityService) FindMFAs(ctx context.Context, tenantID string, userID string, pagination object.Pagination) ([]object.MFA, error) {
	if len(userID) == 0 {
		return nil, errors.New("userID is required")
	}

	return repository.FindMFAs(ctx, is.db, tenantID, userID, pagination)
}

// VerifyMFA updates the verification status of an existing MFA within a specified tenant in the database.
//
// Returns:
//   - Error if there is any issue during updating.
func (is IdentityService) VerifyMFA(ctx context.Context, tenantID string, userID string, mfaID string, body map[string]any) error {
	if len(mfaID) == 0 {
		return errors.New("mfaID is required")
	}

	if len(userID) == 0 {
		return errors.New("userID is required")
	}

	mfaObj, err := is.FindMFA(ctx, tenantID, userID, mfaID)

	if err != nil {
		return err
	}

	providerObj, err := is.FindProvider(ctx, tenantID, mfaObj.ProviderID)

	if err != nil {
		return err
	}

	mfaProvider, err := mfa.GetMFAProvider(providerObj)

	if err != nil {
		return errors.New(fmt.Sprintf("mfa provider not found: %s", mfaObj))
	}

	success, err := mfaProvider.ValidateDatFlow(mfaObj.Properties, body)

	if err != nil {
		return errors.New(fmt.Sprintf("mfa validation error: %s", mfaObj))
	}

	if !success {
		return errors.New("mfa validation failed")
	}

	return repository.VerifieMFA(ctx, is.db, tenantID, userID, mfaID, true)
}

func (is IdentityService) UseRecoveryCode(ctx context.Context, tenantID string, userID string, mfaID string, recoveryCode string) (bool, error) {
	userMFA, err := is.FindMFA(ctx, tenantID, userID, mfaID)
	if err != nil {
		return false, err
	}

	recoveryCodeIndex := slices.Index(userMFA.RecoveryCodes, recoveryCode)

	if recoveryCodeIndex == -1 {
		return false, errors.New("mfaID is invalid")
	}

	// This deletes the uses recovery code
	// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
	userMFA.RecoveryCodes[recoveryCodeIndex] = userMFA.RecoveryCodes[len(userMFA.RecoveryCodes)-1]
	userMFA.RecoveryCodes = userMFA.RecoveryCodes[:len(userMFA.RecoveryCodes)-1]

	err = repository.UpdateMFARecoveryCodes(ctx, is.db, tenantID, userID, mfaID, userMFA.RecoveryCodes)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (is IdentityService) MfaInitDataflow(ctx context.Context, tenantID string, userID string, mfaID string, body json.RawMessage) (map[string]any, error) {
	if len(tenantID) == 0 {
		return nil, errors.New("tenantID is required")
	}

	if len(mfaID) == 0 {
		return nil, errors.New("mfaID is required")
	}

	if len(userID) == 0 {
		return nil, errors.New("userID is required")
	}

	mfaObj, err := is.FindMFA(ctx, tenantID, userID, mfaID)
	if err != nil {
		return nil, err
	}

	providerObj, err := is.FindProvider(ctx, tenantID, mfaObj.ProviderID)
	if err != nil {
		return nil, err
	}

	mfaProvider, err := mfa.GetMFAProvider(providerObj)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("mfa provider not found: %s", mfaObj))
	}

	resp, err := mfaProvider.InitDataFlow(body)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("mfa init dataflow error: %s", mfaObj))
	}

	return resp, nil
}

func (is IdentityService) MFaVerifyDataFlow(ctx context.Context, tenantID string, userID string, mfaID string, body map[string]any) (bool, error) {
	if len(tenantID) == 0 {
		return false, errors.New("tenantID is required")
	}

	if len(mfaID) == 0 {
		return false, errors.New("mfaID is required")
	}

	if len(userID) == 0 {
		return false, errors.New("userID is required")
	}

	mfaObj, err := is.FindMFA(ctx, tenantID, userID, mfaID)

	if err != nil {
		return false, err
	}

	providerObj, err := is.FindProvider(ctx, tenantID, mfaObj.ProviderID)
	if err != nil {
		return false, err
	}

	mfaProvider, err := mfa.GetMFAProvider(providerObj)

	if err != nil {
		return false, errors.New(fmt.Sprintf("mfa provider not found: %s", mfaObj))
	}

	success, err := mfaProvider.ValidateDatFlow(mfaObj.Properties, body)

	if err != nil {
		return false, errors.New(fmt.Sprintf("mfa validation error: %s", mfaObj))
	}

	return success, nil
}
