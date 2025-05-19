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
	"github.com/anthrove/identity/pkg/util"
	"github.com/go-playground/validator/v10"
)

func (is IdentityService) CreateProfilePage(ctx context.Context, tenantID string, createProfilePage object.CreateProfilePage) (object.ProfilePage, error) {
	dbConn, nested := is.getDBConn(ctx)

	if len(tenantID) == 0 {
		return object.ProfilePage{}, errors.New("tenantID is required")
	}

	err := validate.Struct(createProfilePage)

	// TODO Validate each field with tenant configuration

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.ProfilePage{}, errors.Join(fmt.Errorf("problem while validating create user data"), util.ConvertValidationError(validateErrs))
		}
	}

}
