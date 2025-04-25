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

	passProvider, err := is.ImportProvider(ctx, tenant.ID, object.ImportProvider{
		DisplayName: "Password Provider",
		Category:    "auth",
		Type:        "password",
		Parameter:   []byte(`{"min_password_length":4,"max_password_length":50,"min_lowercase_length":0,"min_uppercase_letter":0,"min_digit_letter":0,"min_special_letter":0}`),
	})

	if err != nil {
		return object.Tenant{}, err
	}

	adminGroup, err := is.ImportGroup(ctx, tenant.ID, object.ImportGroup{
		ID:          "admin_____group",
		DisplayName: "Admin Group",
	})

	if err != nil {
		return object.Tenant{}, err
	}

	adminApplication, err := is.ImportApplication(ctx, tenant.ID, object.ImportApplication{
		ID:           "admin_____appli",
		DisplayName:  "Admin Application",
		Logo:         "", // we need a default logo
		SignInURL:    "https://localhost:8080/web/login",
		SignUpURL:    "",
		ForgetURL:    "",
		TermsURL:     "",
		RedirectURLs: []string{"https://localhost:8080/admin"},
	})

	if err != nil {
		return object.Tenant{}, err
	}

	err = is.AppendAuthProviderToApplication(ctx, tenant.ID, adminApplication.ID, passProvider.ID)

	if err != nil {
		return object.Tenant{}, err
	}

	adminUser, err := is.ImportUser(ctx, tenant.ID, object.ImportUser{
		ID:          "admin_____user1",
		Username:    "admin",
		Email:       "admin@admin.intern",
		DisplayName: "Admin User",
		Password:    "admin",
	})

	if err != nil {
		return object.Tenant{}, err
	}

	err = is.AppendUserToGroup(ctx, tenant.ID, adminUser.ID, adminGroup.ID)

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
			}, importTenant.ID)
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

func (is IdentityService) ImportProvider(ctx context.Context, tenantID string, importProvider object.ImportProvider) (object.Provider, error) {
	if importProvider.ID == "" {
		return is.CreateProvider(ctx, tenantID, object.CreateProvider{
			DisplayName:  importProvider.DisplayName,
			Category:     importProvider.Category,
			ProviderType: importProvider.Type,
			Parameter:    importProvider.Parameter,
		})
	}

	provider, err := is.FindProvider(ctx, tenantID, importProvider.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return is.CreateProvider(ctx, tenantID, object.CreateProvider{
				DisplayName:  importProvider.DisplayName,
				Category:     importProvider.Category,
				ProviderType: importProvider.Type,
				Parameter:    importProvider.Parameter,
			}, importProvider.ID)
		}

		return object.Provider{}, err
	}

	err = is.UpdateProvider(ctx, tenantID, provider.ID, object.UpdateProvider{
		DisplayName: provider.DisplayName,
		Parameter:   provider.Parameter,
	})

	if err != nil {
		return object.Provider{}, err
	}

	provider.DisplayName = importProvider.DisplayName
	provider.Parameter = importProvider.Parameter
	return provider, nil
}

func (is IdentityService) ImportGroup(ctx context.Context, tenantID string, importGroup object.ImportGroup) (object.Group, error) {
	if importGroup.ID == "" {
		return is.CreateGroup(ctx, tenantID, object.CreateGroup{
			DisplayName:   importGroup.DisplayName,
			ParentGroupID: importGroup.ParentGroupID,
		})
	}

	group, err := is.FindGroup(ctx, tenantID, importGroup.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return is.CreateGroup(ctx, tenantID, object.CreateGroup{
				DisplayName:   importGroup.DisplayName,
				ParentGroupID: importGroup.ParentGroupID,
			}, importGroup.ID)
		}

		return object.Group{}, err
	}

	err = is.UpdateGroup(ctx, tenantID, importGroup.ID, object.UpdateGroup{
		DisplayName:   importGroup.DisplayName,
		ParentGroupID: importGroup.ParentGroupID,
	})

	if err != nil {
		return object.Group{}, err
	}

	group.DisplayName = importGroup.DisplayName
	group.ParentGroupID = importGroup.ParentGroupID
	return group, nil
}

func (is IdentityService) ImportApplication(ctx context.Context, tenantID string, importApplication object.ImportApplication) (object.Application, error) {
	if importApplication.ID == "" {
		return is.CreateApplication(ctx, tenantID, object.CreateApplication{
			DisplayName:  importApplication.DisplayName,
			Logo:         importApplication.Logo,
			SignInURL:    importApplication.SignInURL,
			SignUpURL:    importApplication.SignUpURL,
			ForgetURL:    importApplication.ForgetURL,
			TermsURL:     importApplication.TermsURL,
			RedirectURLs: importApplication.RedirectURLs,
		})
	}

	application, err := is.FindApplication(ctx, tenantID, importApplication.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return is.CreateApplication(ctx, tenantID, object.CreateApplication{
				DisplayName:  importApplication.DisplayName,
				Logo:         importApplication.Logo,
				SignInURL:    importApplication.SignInURL,
				SignUpURL:    importApplication.SignUpURL,
				ForgetURL:    importApplication.ForgetURL,
				TermsURL:     importApplication.TermsURL,
				RedirectURLs: importApplication.RedirectURLs,
			}, importApplication.ID)
		}

		return object.Application{}, err
	}

	err = is.UpdateApplication(ctx, tenantID, application.ID, object.UpdateApplication{
		DisplayName:  importApplication.DisplayName,
		Logo:         importApplication.Logo,
		SignInURL:    importApplication.SignInURL,
		SignUpURL:    importApplication.SignUpURL,
		ForgetURL:    importApplication.ForgetURL,
		TermsURL:     importApplication.TermsURL,
		RedirectURLs: importApplication.RedirectURLs,
	})

	if err != nil {
		return object.Application{}, err
	}

	application.DisplayName = importApplication.DisplayName
	application.Logo = importApplication.Logo
	application.SignInURL = importApplication.SignInURL
	application.SignUpURL = importApplication.SignUpURL
	application.ForgetURL = importApplication.ForgetURL
	application.TermsURL = importApplication.TermsURL
	application.RedirectURLs = importApplication.RedirectURLs
	return application, nil
}

func (is IdentityService) ImportUser(ctx context.Context, tenantID string, importUser object.ImportUser) (object.User, error) {
	if importUser.ID == "" {
		user, err := is.CreateUser(ctx, tenantID, object.CreateUser{
			Username:    importUser.Username,
			DisplayName: importUser.DisplayName,
			Email:       importUser.Email,
			Password:    importUser.Password,
		})

		if err != nil {
			return object.User{}, err
		}

		err = is.UpdateUserEmail(ctx, tenantID, user.ID, object.UpdateEmail{
			Email:                  user.Email,
			EmailVerified:          true,
			EmailVerificationToken: "",
		})
		return user, err
	}

	user, err := is.FindUser(ctx, tenantID, importUser.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user, err = is.CreateUser(ctx, tenantID, object.CreateUser{
				Username:    importUser.Username,
				DisplayName: importUser.DisplayName,
				Email:       importUser.Email,
				Password:    importUser.Password,
			}, importUser.ID)

			if err != nil {
				return object.User{}, err
			}

			err = is.UpdateUserEmail(ctx, tenantID, user.ID, object.UpdateEmail{
				Email:                  user.Email,
				EmailVerified:          true,
				EmailVerificationToken: "",
			})
			return user, err
		}

		return object.User{}, err
	}

	err = is.UpdateUser(ctx, tenantID, user.ID, object.UpdateUser{
		DisplayName: importUser.DisplayName,
	})

	if err != nil {
		return object.User{}, err
	}

	err = is.UpdateUserEmail(ctx, tenantID, user.ID, object.UpdateEmail{
		Email:                  importUser.Email,
		EmailVerified:          true,
		EmailVerificationToken: "",
	})

	user.DisplayName = importUser.DisplayName
	user.Email = importUser.Email
	user.EmailVerified = true
	user.EmailVerificationToken = ""
	return user, nil
}
