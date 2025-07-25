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

func (is IdentityService) CreateApplication(ctx context.Context, tenantID string, createApplication object.CreateApplication, opt ...string) (object.Application, error) {
	dbConn, _ := is.getDBConn(ctx)

	err := validate.Struct(createApplication)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.Application{}, errors.Join(fmt.Errorf("problem while validating create application data"), util.ConvertValidationError(validateErrs))
		}
	}

	return repository.CreateApplication(ctx, dbConn, tenantID, createApplication)
}

func (is IdentityService) UpdateApplication(ctx context.Context, tenantID string, applicationID string, updateApplication object.UpdateApplication) error {
	dbConn, _ := is.getDBConn(ctx)

	if len(tenantID) == 0 {
		return errors.New("tenantID is required")
	}

	err := validate.Struct(updateApplication)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating create tenant data"), validateErrs)
		}
	}

	return repository.UpdateApplication(ctx, dbConn, tenantID, applicationID, updateApplication)
}

func (is IdentityService) KillApplication(ctx context.Context, tenantID string, applicationID string) error {
	dbConn, _ := is.getDBConn(ctx)

	return repository.KillApplication(ctx, dbConn, tenantID, applicationID)
}

func (is IdentityService) FindApplication(ctx context.Context, tenantID string, applicationID string) (object.Application, error) {
	dbConn, _ := is.getDBConn(ctx)

	return repository.FindApplication(ctx, dbConn, tenantID, applicationID)
}

func (is IdentityService) FindApplications(ctx context.Context, tenantID string, pagination object.Pagination) ([]object.Application, error) {
	dbConn, _ := is.getDBConn(ctx)

	return repository.FindApplications(ctx, dbConn, tenantID, pagination)
}

func (is IdentityService) AppendAuthProviderToApplication(ctx context.Context, tenantID string, applicationID string, authProviderID string) error {
	dbConn, _ := is.getDBConn(ctx)

	return repository.AppendAuthProviderToApplication(ctx, dbConn, tenantID, applicationID, authProviderID)
}

func (is IdentityService) RemoveAuthProviderFromApplication(ctx context.Context, tenantID string, applicationID string, authProviderID string) error {
	dbConn, _ := is.getDBConn(ctx)

	return repository.RemoveAuthProviderFromApplication(ctx, dbConn, tenantID, applicationID, authProviderID)
}
