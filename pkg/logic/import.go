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
	"github.com/anthrove/identity/pkg/object"
	"gorm.io/gorm"
)

func (is IdentityService) SetupAdminTenant(ctx context.Context) (object.Tenant, error) {
	adminTenant, err := is.FindTenant(ctx, "_____admin_____")

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return object.Tenant{}, err
		}
	} else {
		return adminTenant, nil
	}

	tenant, err := is.ImportTenant(ctx, object.ImportTenant{
		ID:           "_____admin_____",
		DisplayName:  "Admin Tenant",
		PasswordType: "bcrypt",
	})

	if err != nil {
		return object.Tenant{}, err
	}

	return tenant, err
}

func (is IdentityService) ImportTenant(ctx context.Context, importTenant object.ImportTenant) (object.Tenant, error) {
	if importTenant.ID == "" {
		return is.CreateTenant(ctx, object.CreateTenant{
			DisplayName:  importTenant.DisplayName,
			PasswordType: importTenant.PasswordType,
		})
	}

	tenant, err := is.FindTenant(ctx, importTenant.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return is.CreateTenant(ctx, object.CreateTenant{
				DisplayName:  importTenant.DisplayName,
				PasswordType: importTenant.PasswordType,
			})
		}

		return object.Tenant{}, err
	}

	err = is.UpdateTenant(ctx, tenant.ID, object.UpdateTenant{
		DisplayName:          importTenant.DisplayName,
		PasswordType:         importTenant.PasswordType,
		SigningCertificateID: *importTenant.SigningCertificateID,
	})

	if err != nil {
		return object.Tenant{}, err
	}

	tenant.DisplayName = importTenant.DisplayName
	tenant.PasswordType = importTenant.PasswordType
	tenant.SigningCertificateID = importTenant.SigningCertificateID
	return tenant, nil
}
