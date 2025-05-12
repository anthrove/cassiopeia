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
	"encoding/json"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
	"time"
)

type MFA struct {
	ID         string `json:"id" gorm:"primaryKey;type:char(25)" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	UserID     string `json:"user_id" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	ProviderID string `json:"provider_id" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`

	CreatedAt time.Time `json:"createdAt" format:"date-time" example:"2025-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updatedAt" format:"date-time" example:"2025-01-01T00:00:00Z"`

	DisplayName   string          `json:"display_name" validate:"required,max=100" maxLength:"100" example:"Authenticator App"`
	Type          string          `json:"type" validate:"required,max=100" maxLength:"100" example:"totp"`
	Priority      int             `json:"priority" validate:"required" example:"1"`
	Verified      bool            `json:"verified" example:"true"`
	RecoveryCodes []string        `json:"-" swaggerignore:"true" gorm:"type:text[]; serializer:json"`
	Properties    json.RawMessage `json:"properties" validate:"required"`
}

func (base *MFA) BeforeCreate(db *gorm.DB) error {
	if base.ID == "" {
		id, err := gonanoid.New(25)
		if err != nil {
			return err
		}

		base.ID = id
	}

	return nil
}

type CreateMFA struct {
	ProviderID    string          `json:"provider_id" validate:"required" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	DisplayName   string          `json:"display_name" validate:"required,max=100" maxLength:"100" example:"Authenticator App"`
	Type          string          `json:"type" validate:"required,max=100" maxLength:"100" example:"totp"`
	Priority      int             `json:"priority" validate:"required" example:"1"`
	RecoveryCodes []string        `json:"-" swaggerignore:"true"`
	Properties    json.RawMessage `json:"-" swaggerignore:"true"`
}

type UpdateMFA struct {
	DisplayName string `json:"display_name" validate:"required,max=100" maxLength:"100" example:"Authenticator App"`
	Priority    int    `json:"priority" validate:"required" example:"1"`
}

type MFAProviderData struct {
	Properties json.RawMessage `json:"secret" validate:"required"`
}
