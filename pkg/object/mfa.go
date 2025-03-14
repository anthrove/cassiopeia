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

type MFA struct {
	ID     string `json:"id" gorm:"primaryKey;type:char(25)" `
	UserID string `json:"user_id"`

	CreatedAt time.Time `json:"createdAt" format:"date-time"`
	UpdatedAt time.Time `json:"updatedAt" format:"date-time"`

	DisplayName   string   `json:"display_name" validate:"required,max=100" maxLength:"100"`
	Type          string   `json:"type" validate:"required,max=100" maxLength:"100"`
	Priority      int      `json:"priority" validate:"required"`
	Verified      bool     `json:"verified"`
	Secret        string   `json:"secret" validate:"required"`
	RecoveryCodes []string `json:"recovery_codes" validate:"required" gorm:"type:text[]"`
	URI           string   `json:"uri"`
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
	ProviderID  string `json:"provider_id" validate:"required" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	DisplayName string `json:"display_name" validate:"required,max=100" maxLength:"100"`
	Type        string `json:"type" validate:"required,max=100" maxLength:"100"`
	Priority    int    `json:"priority" validate:"required"`
	Secret      string `json:"-" swaggerignore:"true"`
	URI         string `json:"-" swaggerignore:"true"`
}

type UpdateMFA struct {
	DisplayName string `json:"display_name" validate:"required,max=100" maxLength:"100"`
	Priority    int    `json:"priority" validate:"required"`
}

type MFAProviderData struct {
	Secret string `json:"secret" validate:"required"`
	URI    string `json:"uri"`
}
