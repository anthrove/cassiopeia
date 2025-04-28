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
	"database/sql"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
	"time"
)

type Token struct {
	ID            string         `json:"id" gorm:"primaryKey;type:char(25)" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	TenantID      string         `json:"tenant_id" maxLength:"25" minLength:"25" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	ApplicationID string         `json:"application_id" maxLength:"25" minLength:"25" gorm:"type:char(25)" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	UserID        sql.NullString `json:"user_id" maxLength:"25" gorm:"type:char(25)" `

	CreatedAt time.Time `json:"created_at" format:"date-time"`
	ExpiredAt time.Time `json:"expired_at" format:"date-time"`

	RefreshTokenID string `json:"refresh_token"`

	Scope    string `json:"scope"`
	Audience string `json:"audience"`
}

func (base *Token) BeforeCreate(db *gorm.DB) error {
	if base.ID == "" {
		id, err := gonanoid.New(25)
		if err != nil {
			return err
		}

		base.ID = id
	}

	return nil
}

type CreateToken struct {
	ApplicationID string    `json:"application_id" maxLength:"25" minLength:"25"`
	UserID        string    `json:"user_id"`
	Scope         string    `json:"scope"`
	Audience      string    `json:"audience"`
	ExpiredAt     time.Time `json:"expired_at" format:"date-time"`
}
