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
	"regexp"
	"strings"
)

func (is IdentityService) CreateProfilePage(ctx context.Context, tenantID string, userID string, createProfilePage object.CreateProfilePage) (object.ProfilePage, error) {
	dbConn, _ := is.getDBConn(ctx)

	if len(tenantID) == 0 {
		return object.ProfilePage{}, errors.New("tenantID is required")
	}

	err := validate.Struct(createProfilePage)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.ProfilePage{}, errors.Join(fmt.Errorf("problem while validating create user data"), util.ConvertValidationError(validateErrs))
		}
	}

	tenant, err := is.FindTenant(ctx, tenantID)
	if err != nil {
		return object.ProfilePage{}, err
	}

	profileErrors := validateProfilePageFields(tenant.ProfileFields, createProfilePage.Fields)
	if len(profileErrors) > 0 {
		return object.ProfilePage{}, errors.New("multiple field errors: " + strings.Join(profileErrors, ","))
	}

	profilePage, err := repository.CreateUserProfile(ctx, dbConn, tenantID, userID, createProfilePage)

	if err != nil {
		return object.ProfilePage{}, err
	}

	return profilePage, nil
}

func (is IdentityService) UpdateProfilePage(ctx context.Context, tenantID string, userID string, profilePage object.UpdateProfilePage) error {
	dbConn, _ := is.getDBConn(ctx)

	if len(tenantID) == 0 {
		return errors.New("tenantID is required")
	}

	err := validate.Struct(profilePage)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating create user data"), util.ConvertValidationError(validateErrs))
		}
	}

	tenant, err := is.FindTenant(ctx, tenantID)
	if err != nil {
		return err
	}

	profileErrors := validateProfilePageFields(tenant.ProfileFields, profilePage.Fields)
	if len(profileErrors) > 0 {
		return errors.New("multiple field errors: " + strings.Join(profileErrors, ","))
	}

	err = repository.UpdatePageProfile(ctx, dbConn, tenantID, userID, profilePage)

	if err != nil {
		return err
	}

	return nil
}

func (is IdentityService) KillProfilePage(ctx context.Context, tenantID string, userID string) error {
	dbConn, _ := is.getDBConn(ctx)

	if len(tenantID) == 0 {
		return errors.New("tenantID is required")
	}

	err := repository.KillProfilePage(ctx, dbConn, tenantID, userID)

	if err != nil {
		return err
	}

	return nil
}

func (is IdentityService) FindProfilePage(ctx context.Context, tenantID string, userID string) (object.ProfilePage, error) {
	dbConn, _ := is.getDBConn(ctx)

	if len(tenantID) == 0 {
		return object.ProfilePage{}, errors.New("tenantID is required")
	}

	if len(userID) == 0 {
		return object.ProfilePage{}, errors.New("userID is required")
	}

	return repository.FindProfilePage(ctx, dbConn, tenantID, userID)
}

func validateProfilePageFields(tenantFields []object.ProfileField, profileFields []object.ProfilePageField) []string {
	fieldErrors := make([]string, 0)

out:
	for _, tenantField := range tenantFields {
		if !tenantField.Required && len(tenantField.Regex) == 0 {
			continue
		}

		for _, profileField := range profileFields {
			if tenantField.Identifier == profileField.Identifier {
				if tenantField.Required && profileField.Value == nil {
					fieldErrors = append(fieldErrors, fmt.Sprintf("field %s is required", profileField.Identifier))
					continue out
				}

				if len(tenantField.Regex) > 0 {
					r, err := regexp.Compile(tenantField.Regex)

					if err != nil {
						fieldErrors = append(fieldErrors, fmt.Sprintf("regex from field %s is invalid. please contact an administrator", tenantField.Identifier))
						continue out
					}

					strValue, ok := profileField.Value.(string)

					if !ok {
						fieldErrors = append(fieldErrors, fmt.Sprintf("field %s needs to be a string to be validated by a regex", profileField.Identifier))
						continue out
					}

					success := r.MatchString(strValue)

					if !success {
						fieldErrors = append(fieldErrors, fmt.Sprintf("field %s is invalid", profileField.Identifier))
						continue out
					}
				}

				continue out
			}
		}

		if tenantField.Required {
			fieldErrors = append(fieldErrors, fmt.Sprintf("field %s is required", tenantField.Identifier))
			continue out
		}
	}

	return fieldErrors
}
