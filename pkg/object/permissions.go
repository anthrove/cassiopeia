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
)

type Permission struct {
	ID         string `json:"id" gorm:"primaryKey;type:char(25)" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	TenantID   string `json:"tenant_id" gorm:"type:char(25)" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	EnforcerID string `json:"enforcer_id" gorm:"type:char(25)" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Users  []string `json:"users" gorm:"serializer:json"`
	Groups []string `json:"groups" gorm:"serializer:json"`

	V1 []string `json:"v1" gorm:"serializer:json"`
	V2 []string `json:"v2" gorm:"serializer:json"`
	V3 []string `json:"v3" gorm:"serializer:json"`
	V4 []string `json:"v4" gorm:"serializer:json"`
	V5 []string `json:"v5" gorm:"serializer:json"`
}

func (base *Permission) BeforeCreate(db *gorm.DB) error {
	if base.ID == "" {
		id, err := gonanoid.New(25)
		if err != nil {
			return err
		}

		base.ID = id
	}

	return nil
}

type CreatePermission struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`

	EnforcerID string `json:"enforcer_id" validate:"required"`

	Users  []string `json:"users"`
	Groups []string `json:"groups"`

	V1 []string `json:"v1"`
	V2 []string `json:"v2"`
	V3 []string `json:"v3"`
	V4 []string `json:"v4"`
	V5 []string `json:"v5"`
}

type UpdatePermission struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`

	EnforcerID string `json:"enforcer_id" validate:"required"`

	Users  []string `json:"users"`
	Groups []string `json:"groups"`

	V1 []string `json:"v1"`
	V2 []string `json:"v2"`
	V3 []string `json:"v3"`
	V4 []string `json:"v4"`
	V5 []string `json:"v5"`
}
