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

package object

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
	"time"
)

type Application struct {
	ID            string `json:"id" gorm:"primaryKey;type:char(25)" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	TenantID      string `json:"tenant_id" maxLength:"25" minLength:"25" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	CertificateID string `json:"certificate_id" gorm:"type:char(25)" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`

	CreatedAt time.Time `json:"createdAt" format:"date-time"`
	UpdatedAt time.Time `json:"updatedAt" format:"date-time"`

	DisplayName string `json:"display_name" gorm:"type:varchar(100)" example:"Frontend Application"`
	Logo        string `json:"logo" gorm:"type:varchar(255)" example:"https://domain.tld/files/logo.png"`

	SignInURL string `json:"sign_in_url" gorm:"type:varchar(255)"`
	SignUpURL string `json:"sign_up_url" gorm:"type:varchar(255)"`
	ForgetURL string `json:"forget_url" gorm:"type:varchar(255)"`
	TermsURL  string `json:"terms_url" gorm:"type:varchar(255)"`
}

// BeforeCreate is a GORM hook that is called before a new group record is inserted into the database.
// It generates a unique ID for the group if it is not already set.
//
// Parameters:
//   - db: a gorm.DB instance representing the database connection.
//
// Returns:
//   - An error if there is any issue generating the unique ID.
func (base *Application) BeforeCreate(db *gorm.DB) error {
	if base.ID == "" {
		id, err := gonanoid.New(25)
		if err != nil {
			return err
		}

		base.ID = id
	}

	return nil
}

type CreateApplication struct {
	CertificateID string `json:"certificate_id" validate:"required,len=25" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`

	DisplayName string `json:"display_name" validate:"required,max=100" example:"Frontend Application"`
	Logo        string `json:"logo"  example:"https://domain.tld/files/logo.png"`

	SignInURL string `json:"sign_in_url"`
	SignUpURL string `json:"sign_up_url"`
	ForgetURL string `json:"forget_url"`
	TermsURL  string `json:"terms_url"`
}

type UpdateApplication struct {
	CertificateID string `json:"certificate_id" validate:"required,len=25" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`

	DisplayName string `json:"display_name" validate:"required,max=100" example:"Frontend Application"`
	Logo        string `json:"logo"  example:"https://domain.tld/files/logo.png"`

	SignInURL string `json:"sign_in_url"`
	SignUpURL string `json:"sign_up_url"`
	ForgetURL string `json:"forget_url"`
	TermsURL  string `json:"terms_url"`
}
