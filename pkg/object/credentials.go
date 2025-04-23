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

type Credentials struct {
	ID       string `json:"id" gorm:"primaryKey;type:char(25)" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	TenantID string `json:"tenant_id" gorm:"type:char(25)" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	UserID   string `json:"user_id" gorm:"type:char(25)" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`

	CreatedAt time.Time      `json:"created_at" format:"date-time"`
	UpdatedAt time.Time      `json:"updated_at" format:"date-time"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" format:"date-time" gorm:"index"`

	Type     string         `json:"type"`
	Metadata map[string]any `json:"metadata"`
	Enabled  bool           `json:"enabled"`
}

func (base *Credentials) BeforeCreate(db *gorm.DB) error {
	if base.ID == "" {
		id, err := gonanoid.New(25)
		if err != nil {
			return err
		}

		base.ID = id
	}

	return nil
}

type CreateCredential struct {
	UserID string `json:"user_id" gorm:"type:char(25)"`

	Type     string         `json:"type"`
	Metadata map[string]any `json:"metadata"`
	Enabled  bool           `json:"enabled"`
}

type UpdateCredential struct {
	Metadata map[string]any `json:"metadata"`
	Enabled  bool           `json:"enabled"`
}
