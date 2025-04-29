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

type Provider struct {
	ID       string `json:"id" gorm:"primaryKey;type:char(25)" `
	TenantID string `json:"tenant_id"`

	CreatedAt time.Time `json:"created_at" format:"date-time"`
	UpdatedAt time.Time `json:"updated_at" format:"date-time"`

	DisplayName  string          `json:"display_name" validate:"required,max=100" maxLength:"100"`
	Category     string          `json:"category" validate:"required,max=100" maxLength:"100"`
	ProviderType string          `json:"provider_type" validate:"required,max=100" maxLength:"100"`
	Parameter    json.RawMessage `json:"parameter"`
}

func (base *Provider) BeforeCreate(db *gorm.DB) error {
	if base.ID == "" {
		id, err := gonanoid.New(25)
		if err != nil {
			return err
		}

		base.ID = id
	}

	return nil
}

type CreateProvider struct {
	DisplayName  string          `json:"display_name" validate:"required,max=100" maxLength:"100"`
	Category     string          `json:"category" validate:"required,max=100" maxLength:"100"`
	ProviderType string          `json:"provider_type" validate:"required,max=100" maxLength:"100"`
	Parameter    json.RawMessage `json:"parameter"`
}

type UpdateProvider struct {
	DisplayName string          `json:"display_name" validate:"required,max=100" maxLength:"100"`
	Parameter   json.RawMessage `json:"parameter"`
}

type ProviderConfigurationField struct {
	FieldKey  string `json:"field_key"`
	FieldType string `json:"field_type"`
}

type ProviderConfigurationFieldValue struct {
	FieldKey   string `json:"field_key"`
	FieldValue string `json:"field_value"`
}
