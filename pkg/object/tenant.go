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

// Tenant represents a tenant entity in the system.
// It contains information about the tenant such as its ID, timestamps, display name, password type, and associated groups.
type Tenant struct {
	ID string `json:"id" gorm:"primaryKey;type:char(25)" maxLength:"25" minLength:"25" example:"BsOOa4igppKxYwhAQQrD3GCRZ"`

	CreatedAt time.Time `json:"created_at" format:"date-time"`
	UpdatedAt time.Time `json:"updated_at" format:"date-time"`

	DisplayName  string `json:"display_name" gorm:"type:varchar(100)" maxLength:"100" example:"Tenant Title"`
	PasswordType string `json:"password_type" gorm:"type:varchar(100)" maxLength:"100" example:"bcrypt"`

	SigningCertificateID *string `json:"signing_certificate_id" gorm:"type:char(25)" maxLength:"25" minLength:"25" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`

	ProfileFields []ProfileField `json:"profile_fields" gorm:"serializer:json"`

	Groups       []Group           `json:"-" swaggerignore:"true"`
	Providers    []Provider        `json:"-" swaggerignore:"true"`
	Templates    []MessageTemplate `json:"-" swaggerignore:"true"`
	Users        []User            `json:"-" swaggerignore:"true"`
	Application  []Application     `json:"-" swaggerignore:"true"`
	Certificates []Certificate     `json:"-" swaggerignore:"true"`
	Tokens       []Token           `json:"-" swaggerignore:"true"`
}

// BeforeCreate is a GORM hook that is called before a new tenant record is inserted into the database.
// It generates a unique ID for the tenant if it is not already set.
//
// Parameters:
//   - db: a gorm.DB instance representing the database connection.
//
// Returns:
//   - An error if there is any issue generating the unique ID.
func (base *Tenant) BeforeCreate(db *gorm.DB) error {
	if base.ID == "" {
		id, err := gonanoid.New(25)
		if err != nil {
			return err
		}

		base.ID = id
	}

	return nil
}

// CreateTenant represents the data required to create a new tenant.
// It includes the display name and password type, both of which are required and have a maximum length of 100 characters.
type CreateTenant struct {
	DisplayName   string         `json:"display_name" validate:"required,max=100" maxLength:"100"`
	PasswordType  string         `json:"password_type" validate:"required,max=100" maxLength:"100"`
	ProfileFields []ProfileField `json:"profile_fields" validate:"required"`
}

// UpdateTenant represents the data required to update an existing tenant.
// It includes the display name and password type, both of which are required and have a maximum length of 100 characters.
type UpdateTenant struct {
	DisplayName          string         `json:"display_name" validate:"required,max=100" maxLength:"100"`
	PasswordType         string         `json:"password_type" validate:"required,max=100" maxLength:"100"`
	SigningCertificateID string         `json:"signing_certificate_id" validate:"required,max=25" maxLength:"25"`
	ProfileFields        []ProfileField `json:"profile_fields" validate:"required"`
}

type ProfileField struct {
	Identifier  string     `json:"identifier" validate:"required,max=100" maxLength:"100"`
	DisplayName string     `json:"display_name"`
	Regex       string     `json:"regex"`
	Required    bool       `json:"required"`
	ModifyBy    ModifyType `json:"modify_by"`
	ViewBy      AllowView  `json:"view_by"`
}

type ModifyType string

const (
	ModifyTypeImmutable ModifyType = "immutable"
	ModifyTypeSelf      ModifyType = "self"
)

type AllowView string

const (
	AllowViewPublic AllowView = "public"
	AllowViewSelf   AllowView = "self"
)
