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

// CreateMessageTemplate creates a new messageTemplate in the system.
// It validates the input data using the validator package and returns an error if validation fails.
// If validation passes, it calls the repository to create the messageTemplate in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - createMessageTemplate: object containing the details of the messageTemplate to be created.
//
// Returns:
//   - MessageTemplate object if creation is successful.
//   - Error if there is any issue during validation or creation.
func (is IdentityService) CreateMessageTemplate(ctx context.Context, tenantID string, createMessageTemplate object.CreateMessageTemplate) (object.MessageTemplate, error) {
	dbConn, _ := is.getDBConn(ctx)

	err := validate.Struct(createMessageTemplate)

	if err != nil {

		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.MessageTemplate{}, errors.Join(fmt.Errorf("problem while validating create messageTemplate data"), util.ConvertValidationError(validateErrs))
		}
	}

	return repository.CreateMessageTemplate(ctx, dbConn, tenantID, createMessageTemplate)
}

// UpdateMessageTemplate updates an existing messageTemplate's information in the system.
// It validates the input data using the validator package and returns an error if validation fails.
// If validation passes, it calls the repository to update the messageTemplate in the database.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - messageTemplateID: unique identifier of the messageTemplate to be updated.
//   - updateMessageTemplate: object containing the updated details of the messageTemplate.
//
// Returns:
//   - Error if there is any issue during validation or updating.
func (is IdentityService) UpdateMessageTemplate(ctx context.Context, tenantID string, messageTemplateID string, updateMessageTemplate object.UpdateMessageTemplate) error {
	dbConn, _ := is.getDBConn(ctx)

	if len(messageTemplateID) == 0 {
		return errors.New("messageTemplateID is required")
	}

	err := validate.Struct(updateMessageTemplate)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating create messageTemplate data"), validateErrs)
		}
	}

	return repository.UpdateMessageTemplate(ctx, dbConn, tenantID, messageTemplateID, updateMessageTemplate)
}

// KillMessageTemplate deletes an existing messageTemplate from the system.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - messageTemplateID: unique identifier of the messageTemplate to be deleted.
//
// Returns:
//   - Error if there is any issue during deletion.
func (is IdentityService) KillMessageTemplate(ctx context.Context, tenantID string, messageTemplateID string) error {
	dbConn, _ := is.getDBConn(ctx)

	return repository.KillMessageTemplate(ctx, dbConn, tenantID, messageTemplateID)
}

// FindMessageTemplate retrieves a specific messageTemplate from the system.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - messageTemplateID: unique identifier of the messageTemplate to be retrieved.
//
// Returns:
//   - MessageTemplate object if retrieval is successful.
//   - Error if there is any issue during retrieval.
func (is IdentityService) FindMessageTemplate(ctx context.Context, tenantID string, messageTemplateID string) (object.MessageTemplate, error) {
	dbConn, _ := is.getDBConn(ctx)

	return repository.FindMessageTemplate(ctx, dbConn, tenantID, messageTemplateID)
}

// FindMessageTemplates retrieves a list of messageTemplates from the system, with pagination support.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - pagination: object containing pagination details (limit and page).
//
// Returns:
//   - Slice of MessageTemplate objects if retrieval is successful.
//   - Error if there is any issue during retrieval.
func (is IdentityService) FindMessageTemplates(ctx context.Context, tenantID string, pagination object.Pagination) ([]object.MessageTemplate, error) {
	dbConn, _ := is.getDBConn(ctx)

	return repository.FindMessageTemplates(ctx, dbConn, tenantID, pagination)
}
